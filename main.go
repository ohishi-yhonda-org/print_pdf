package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

// イベントログ書き込み関数
func writeEventLog(level string, message string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logMessage := fmt.Sprintf("[%s] %s: %s", timestamp, level, message)

	// コンソール出力
	fmt.Println(logMessage)

	// 標準ログ出力（Windowsサービス時はイベントログに転送される）
	log.Printf("%s: %s", level, message)

	// ファイルログも出力（オプション）
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
		"timestamp": time.Now().Format(time.RFC3339),
	}
	json.NewEncoder(w).Encode(response)
}

func main() {
	writeEventLog("INFO", "PDF生成システム - Go版 (HTTPサーバーモード) 開始")
	writeEventLog("INFO", "Windowsフォント対応")

	// HTTPルートの設定
	http.HandleFunc("/generate-pdf", generatePDFHandler)
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		writeEventLog("INFO", fmt.Sprintf("ルートアクセス from %s", r.RemoteAddr))
		fmt.Fprintf(w, `
PDF Generator API Server

Available endpoints:
- POST /generate-pdf : Generate PDF from JSON data
- GET  /health       : Health check

Example request:
curl -X POST http://localhost:8081/generate-pdf \
  -H "Content-Type: application/json" \
  -d '[{"Car":"test","Name":"テスト","Ryohi":[]}]'
`)
	})

	// サーバー起動
	port := ":8081"
	writeEventLog("INFO", fmt.Sprintf("HTTPサーバーを起動中... http://localhost%s", port))
	writeEventLog("INFO", "PDF生成エンドポイント: POST /generate-pdf")
	writeEventLog("INFO", "ヘルスチェック: GET /health")
	writeEventLog("INFO", "サーバーを停止するには Ctrl+C を押してください")

	// サーバー起動（ブロッキング）
	if err := http.ListenAndServe(port, nil); err != nil {
		writeEventLog("FATAL", fmt.Sprintf("サーバー起動エラー: %v", err))
		log.Fatalf("サーバー起動エラー: %v", err)
	}
}
