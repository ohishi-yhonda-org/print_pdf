package main

import (
	"archive/zip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"  
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/debug"
	"golang.org/x/sys/windows/svc/eventlog"
)

// Version information (set during build with -ldflags)
var Version = "dev"

// グローバル変数
var (
	httpServer *http.Server
	elog       debug.Log
)

// Windowsサービスハンドラー
type service struct{}

func (m *service) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown | svc.AcceptPauseAndContinue
	changes <- svc.Status{State: svc.StartPending}
	
	// HTTPサーバーを起動
	go startHTTPServer()
	
	changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
	
	// サービス制御メッセージを待機
	for {
		select {
		case c := <-r:
			switch c.Cmd {
			case svc.Interrogate:
				changes <- c.CurrentStatus
			case svc.Stop, svc.Shutdown:
				writeEventLog("INFO", "サービス停止要求を受信")
				if httpServer != nil {
					ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
					defer cancel()
					httpServer.Shutdown(ctx)
				}
				changes <- svc.Status{State: svc.StopPending}
				return
			case svc.Pause:
				changes <- svc.Status{State: svc.Paused, Accepts: cmdsAccepted}
			case svc.Continue:
				changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
			default:
				elog.Error(1, fmt.Sprintf("unexpected control request #%d", c))
			}
		}
	}
}

// HTTPサーバー起動関数
func startHTTPServer() {
	writeEventLog("INFO", "PDF生成システム - Go版 (HTTPサーバーモード) 開始")
	writeEventLog("INFO", fmt.Sprintf("バージョン: %s", Version))
	writeEventLog("INFO", "Windowsフォント対応")

	// 起動時に自動アップデートをチェック（dev環境では無効）
	if Version != "dev" {
		go func() {
			// 少し待ってから実行（サーバー起動後）
			time.Sleep(5 * time.Second)
			checkForUpdates()
		}()
	} else {
		writeEventLog("INFO", "開発環境のため自動アップデートを無効にしています")
	}

	// HTTPルートの設定
	http.HandleFunc("/generate-pdf", generatePDFHandler)
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		writeEventLog("INFO", fmt.Sprintf("ルートアクセス from %s", r.RemoteAddr))
		fmt.Fprintf(w, `
PDF Generator API Server %s

Available endpoints:
- POST /generate-pdf : Generate PDF from JSON data
- GET  /health       : Health check

Example request:
curl -X POST http://localhost:8081/generate-pdf \
  -H "Content-Type: application/json" \
  -d '[{"car":"test","name":"テスト","ryohi":[]}]'

Version: %s
Platform: %s %s
`, Version, Version, runtime.GOOS, runtime.GOARCH)
	})

	// サーバー起動
	port := ":8081"
	writeEventLog("INFO", fmt.Sprintf("HTTPサーバーを起動中... http://localhost%s", port))
	writeEventLog("INFO", "PDF生成エンドポイント: POST /generate-pdf")
	writeEventLog("INFO", "ヘルスチェック: GET /health")

	httpServer = &http.Server{
		Addr: port,
	}

	// サーバー起動
	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		writeEventLog("FATAL", fmt.Sprintf("サーバー起動エラー: %v", err))
		log.Fatalf("サーバー起動エラー: %v", err)
	}
}

// GitHubリリース情報
type GitHubRelease struct {
	TagName string `json:"tag_name"`
	Name    string `json:"name"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

// イベントログ書き込み関数（Windowsサービス対応）
func writeEventLog(level string, message string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logMessage := fmt.Sprintf("[%s] %s: %s", timestamp, level, message)

	// コンソール出力（デバッグモード時のみ）
	if elog != nil {
		switch level {
		case "ERROR", "FATAL":
			elog.Error(1, message)
		case "WARN":
			elog.Warning(1, message)
		default:
			elog.Info(1, message)
		}
	} else {
		fmt.Println(logMessage)
	}

	// 標準ログ出力
	log.Printf("%s: %s", level, message)

	// ファイルログも出力
	writeToLogFile(logMessage)
}

// ログファイル書き込み関数
func writeToLogFile(message string) {
	logFile := "pdf_generator_service.log"
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("ログファイル書き込みエラー: %v", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(message + "\n")
	if err != nil {
		log.Printf("ログファイル書き込みエラー: %v", err)
	}
}

// HTTPからデータを取得する関数
func fetchDataFromAPI(url string) ([]Item, error) {
	writeEventLog("INFO", fmt.Sprintf("APIからデータを取得開始: %s", url))

	// HTTPクライアントの設定
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// HTTPリクエストを送信
	resp, err := client.Get(url)
	if err != nil {
		writeEventLog("ERROR", fmt.Sprintf("HTTP リクエストエラー: %v", err))
		return nil, fmt.Errorf("HTTP リクエストエラー: %v", err)
	}
	defer resp.Body.Close()

	// レスポンスステータスをチェック
	if resp.StatusCode != http.StatusOK {
		writeEventLog("ERROR", fmt.Sprintf("HTTP エラー: %d %s", resp.StatusCode, resp.Status))
		return nil, fmt.Errorf("HTTP エラー: %d %s", resp.StatusCode, resp.Status)
	}

	// レスポンスボディを読み取り
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		writeEventLog("ERROR", fmt.Sprintf("レスポンス読み取りエラー: %v", err))
		return nil, fmt.Errorf("レスポンス読み取りエラー: %v", err)
	}

	// JSONをパース
	var data []Item
	if err := json.Unmarshal(body, &data); err != nil {
		writeEventLog("ERROR", fmt.Sprintf("JSON パースエラー: %v", err))
		return nil, fmt.Errorf("JSON パースエラー: %v", err)
	}

	writeEventLog("INFO", fmt.Sprintf("データ取得完了: %d件", len(data)))
	return data, nil
}

// HTTPハンドラー: PDF生成エンドポイント
func generatePDFHandler(w http.ResponseWriter, r *http.Request) {
	// CORSヘッダーを設定
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// OPTIONSリクエストの処理
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// POSTメソッドのみ許可
	if r.Method != "POST" {
		writeEventLog("WARN", fmt.Sprintf("不正なメソッドでのアクセス: %s from %s", r.Method, r.RemoteAddr))
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	writeEventLog("INFO", fmt.Sprintf("PDF生成リクエストを受信 from %s", r.RemoteAddr))

	// リクエストボディを読み取り
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeEventLog("ERROR", fmt.Sprintf("リクエストボディ読み取りエラー: %v", err))
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// JSONをパース
	var requestData []Item
	if err := json.Unmarshal(body, &requestData); err != nil {
		writeEventLog("ERROR", fmt.Sprintf("JSON パースエラー: %v", err))
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	writeEventLog("INFO", fmt.Sprintf("受信データ: %d件のアイテム", len(requestData)))

	// PDF生成処理
	writeEventLog("INFO", "ReportLabスタイルPDF生成を開始")
	reportlabClient := NewReportLabStylePdfClient(requestData)

	if reportlabClient != nil {
		writeEventLog("INFO", "ReportLabスタイルPDF生成完了")

		// 成功レスポンス
		response := map[string]interface{}{
			"status":  "success",
			"message": "PDF generated successfully",
			"items":   len(requestData),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	} else {
		writeEventLog("ERROR", "ReportLabスタイルPDF生成に失敗")

		// エラーレスポンス
		response := map[string]interface{}{
			"status":  "error",
			"message": "Failed to generate PDF",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
	}
}

// ヘルスチェックエンドポイント
func healthHandler(w http.ResponseWriter, r *http.Request) {
	writeEventLog("INFO", fmt.Sprintf("ヘルスチェックアクセス from %s", r.RemoteAddr))
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"status":    "ok",
		"service":   "PDF Generator",
		"version":   Version,
		"timestamp": time.Now().Format(time.RFC3339),
	}
	json.NewEncoder(w).Encode(response)
}

// 自動アップデート機能
func checkForUpdates() {
	// dev環境では自動アップデートを実行しない
	if Version == "dev" {
		writeEventLog("INFO", "開発環境のため自動アップデートをスキップします")
		return
	}

	writeEventLog("INFO", fmt.Sprintf("現在のバージョン: %s", Version))
	writeEventLog("INFO", "GitHubリリースの最新バージョンをチェック中...")

	// GitHub APIから最新リリース情報を取得
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get("https://api.github.com/repos/ohishi-yhonda-org/print_pdf/releases/latest")
	if err != nil {
		writeEventLog("WARN", fmt.Sprintf("アップデートチェックに失敗: %v", err))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		writeEventLog("WARN", fmt.Sprintf("GitHub API エラー: %d", resp.StatusCode))
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		writeEventLog("WARN", fmt.Sprintf("レスポンス読み取りエラー: %v", err))
		return
	}

	var release GitHubRelease
	if err := json.Unmarshal(body, &release); err != nil {
		writeEventLog("WARN", fmt.Sprintf("JSON パースエラー: %v", err))
		return
	}

	latestVersion := release.TagName
	writeEventLog("INFO", fmt.Sprintf("最新バージョン: %s", latestVersion))

	// バージョン比較
	if Version == "dev" || Version != latestVersion {
		writeEventLog("INFO", fmt.Sprintf("新しいバージョンが利用可能: %s -> %s", Version, latestVersion))

		// 自動アップデートを実行
		if performUpdate(release) {
			writeEventLog("INFO", "アップデート完了。アプリケーションを再起動します...")
			os.Exit(0) // 正常終了してサービス管理で再起動される
		}
	} else {
		writeEventLog("INFO", "最新バージョンを使用中です")
	}
}

// アップデートを実行
func performUpdate(release GitHubRelease) bool {
	// Windows用のzipファイルを探す
	var downloadURL string
	for _, asset := range release.Assets {
		if strings.Contains(asset.Name, ".zip") && strings.Contains(asset.Name, release.TagName) {
			downloadURL = asset.BrowserDownloadURL
			break
		}
	}

	if downloadURL == "" {
		writeEventLog("ERROR", "アップデート用のzipファイルが見つかりません")
		return false
	}

	writeEventLog("INFO", fmt.Sprintf("アップデートファイルをダウンロード中: %s", downloadURL))

	// ダウンロード
	client := &http.Client{Timeout: 300 * time.Second} // 5分タイムアウト
	resp, err := client.Get(downloadURL)
	if err != nil {
		writeEventLog("ERROR", fmt.Sprintf("ダウンロードエラー: %v", err))
		return false
	}
	defer resp.Body.Close()

	// 一時ファイルに保存
	tempFile := "update_temp.zip"
	file, err := os.Create(tempFile)
	if err != nil {
		writeEventLog("ERROR", fmt.Sprintf("一時ファイル作成エラー: %v", err))
		return false
	}
	defer file.Close()
	defer os.Remove(tempFile)

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		writeEventLog("ERROR", fmt.Sprintf("ファイル書き込みエラー: %v", err))
		return false
	}
	file.Close()

	writeEventLog("INFO", "ダウンロード完了。アップデートを適用中...")

	// zipファイルを解凍
	if err := extractUpdate(tempFile); err != nil {
		writeEventLog("ERROR", fmt.Sprintf("アップデート適用エラー: %v", err))
		return false
	}

	writeEventLog("INFO", "アップデート適用完了")
	return true
}

// zipファイルを解凍してアップデートを適用
func extractUpdate(zipPath string) error {
	// 現在の実行ファイル名を取得
	currentExe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("実行ファイルパス取得エラー: %v", err)
	}

	// zipファイルを開く
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("zip読み込みエラー: %v", err)
	}
	defer reader.Close()

	// バックアップ作成
	backupPath := currentExe + ".backup"
	if err := copyFile(currentExe, backupPath); err != nil {
		return fmt.Errorf("バックアップ作成エラー: %v", err)
	}

	// 解凍処理
	for _, file := range reader.File {
		if file.Name == "print_pdf.exe" {
			// 実行ファイルを更新
			if err := extractFile(file, currentExe+".new"); err != nil {
				return fmt.Errorf("新しい実行ファイル展開エラー: %v", err)
			}

			// バッチファイルでファイル置換を実行（Windowsでは実行中のファイルを置換できないため）
			batchContent := fmt.Sprintf(`@echo off
timeout /t 2 /nobreak >nul
move "%s" "%s"
move "%s" "%s"
start "" "%s"
del "%%~f0"
`, currentExe+".new", currentExe, backupPath, currentExe+".old", currentExe)

			batchFile := "update_replace.bat"
			if err := os.WriteFile(batchFile, []byte(batchContent), 0755); err != nil {
				return fmt.Errorf("バッチファイル作成エラー: %v", err)
			}

			// バッチファイルを実行（非同期）
			cmd := exec.Command("cmd", "/c", batchFile)
			if err := cmd.Start(); err != nil {
				return fmt.Errorf("アップデートスクリプト実行エラー: %v", err)
			}

			return nil
		}
	}

	return fmt.Errorf("アップデート用実行ファイルが見つかりません")
}

// ファイルをコピー
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

// zipファイルから個別ファイルを展開
func extractFile(file *zip.File, destPath string) error {
	reader, err := file.Open()
	if err != nil {
		return err
	}
	defer reader.Close()

	// ディレクトリを作成
	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return err
	}

	// ファイルを作成
	writer, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer writer.Close()

	_, err = io.Copy(writer, reader)
	return err
}

func main() {
	// Windowsサービスとして実行されているかチェック
	isWindowsService, err := svc.IsWindowsService()
	if err != nil {
		log.Fatalf("サービス状態確認エラー: %v", err)
	}

	if isWindowsService {
		// Windowsサービスとして実行
		runWindowsService()
	} else {
		// コンソールアプリケーションとして実行
		runConsoleApp()
	}
}

// Windowsサービスとして実行
func runWindowsService() {
	var err error
	
	// イベントログを開く（失敗しても続行）
	elog, err = eventlog.Open("PDF Generator API Service")
	if err != nil {
		// イベントログが開けない場合はファイルログのみ使用
		elog = nil
		log.Printf("イベントログを開けませんでした: %v", err)
	}
	defer func() {
		if elog != nil {
			elog.Close()
		}
	}()

	writeEventLog("INFO", "PDF Generator API Service をWindowsサービスとして開始")

	err = svc.Run("PDF Generator API Service", &service{})
	if err != nil {
		writeEventLog("FATAL", fmt.Sprintf("サービス実行エラー: %v", err))
		log.Fatalf("サービス実行エラー: %v", err)
	}
}

// コンソールアプリケーションとして実行
func runConsoleApp() {
	writeEventLog("INFO", "PDF Generator をコンソールアプリケーションとして開始")
	writeEventLog("INFO", "サーバーを停止するには Ctrl+C を押してください")

	// シグナルハンドリング
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// HTTPサーバーを別goroutineで起動
	go startHTTPServer()

	// 終了シグナルを待機
	<-sigChan
	writeEventLog("INFO", "終了シグナルを受信。サーバーを停止中...")

	if httpServer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(ctx); err != nil {
			writeEventLog("ERROR", fmt.Sprintf("サーバー停止エラー: %v", err))
		}
	}

	writeEventLog("INFO", "サーバーが正常に停止しました")
}
