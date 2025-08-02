# PDF Generator - Go版

HTTP API サーバーとして動作するPDF生成システムです。日本語フォント対応、自動アップデート機能、Windowsサービス対応を備えています。

## 機能

- **HTTP API サーバー**: ポート8081でREST API提供
- **PDF生成**: ReportLabスタイルの日本語対応PDF生成
- **自動アップデート**: GitHub Releasesからの自動更新
- **イベントログ**: 包括的なログ記録（コンソール + ファイル + Windows Event Log）
- **Windowsサービス**: サービスとしての自動起動対応
- **CORS対応**: フロントエンドからのクロスオリジンアクセス対応

## API エンドポイント

### POST /generate-pdf
PDF生成エンドポイント

**リクエスト例:**
```json
[
  {
    "car": "Vehicle001",
    "name": "Test Item 1",
    "purpose": "Business Meeting",
    "startDate": "2025-08-01",
    "endDate": "2025-08-02",
    "price": 5000,
    "tax": 500.0,
    "description": "Test description",
    "ryohi": [],
    "office": "Tokyo Office",
    "payDay": "2025-08-15"
  }
]
```

**レスポンス:**
```json
{
  "status": "success",
  "message": "PDF generated successfully",
  "items": 1
}
```

### GET /health
ヘルスチェックエンドポイント

**レスポンス:**
```json
{
  "status": "ok",
  "service": "PDF Generator",
  "version": "v1.0.0",
  "timestamp": "2025-08-02T16:33:24+09:00"
}
```

## インストール方法

### 方法1: GitHub Releasesからダウンロード
1. [Releases](https://github.com/ohishi-yhonda-org/print_pdf/releases) から最新の `print_pdf_vX.X.X.zip` をダウンロード
2. 任意のフォルダに解凍
3. PowerShellでフォルダに移動
4. `service_manager.bat` を実行してWindowsサービスとして登録

### 方法2: ソースからビルド
```bash
git clone https://github.com/ohishi-yhonda-org/print_pdf.git
cd print_pdf
go build -ldflags "-X main.Version=v1.0.0" -o print_pdf.exe .
```

## 開発者向け

### ビルドとテスト
```bash
# テスト実行
go test -v ./...

# カバレッジ付きテスト
go test -v -race -coverprofile=coverage.out ./...

# リンター実行
golangci-lint run ./...

# ビルド（バージョン指定）
go build -ldflags "-X main.Version=v1.0.0" -o print_pdf.exe .
```

### リリース作成方法

#### 方法1: タグでリリース（推奨）
```bash
# 新しいバージョンタグを作成
git tag v1.0.1
git push origin v1.0.1

# GitHub Actionsが自動的に:
# 1. テストを実行
# 2. バージョン付きでビルド
# 3. zipアーカイブを作成
# 4. GitHub Releasesにアップロード
```

#### 方法2: 手動ワークフロー実行
1. GitHubリポジトリの「Actions」タブに移動
2. 「CI」ワークフローを選択
3. 「Run workflow」をクリック
4. 以下のパラメータを設定:
   - `create_release`: `true`
   - `version`: `v1.0.1` (任意のバージョン)
5. 「Run workflow」を実行

### 自動アップデート機能

アプリケーションは起動時に自動的に最新バージョンをチェックし、新しいバージョンが利用可能な場合は自動的にアップデートを実行します。

**アップデートプロセス:**
1. 起動5秒後にGitHub API (`/repos/ohishi-yhonda-org/print_pdf/releases/latest`) をチェック
2. 現在のバージョンと比較
3. 新しいバージョンが存在する場合、zipファイルをダウンロード
4. 現在の実行ファイルをバックアップ
5. 新しいファイルで置換
6. アプリケーションを自動再起動

## Windowsサービス管理

### service_manager.bat の使用方法
```cmd
# サービスをインストール
service_manager.bat install

# サービスを開始
service_manager.bat start

# サービスを停止
service_manager.bat stop

# サービスをアンインストール
service_manager.bat remove

# サービス状態を確認
service_manager.bat status

# ログを表示
service_manager.bat logs
```

## ログ

アプリケーションは以下の場所にログを出力します:
- **コンソール出力**: リアルタイムログ
- **ファイルログ**: `pdf_generator_service.log`
- **Windows Event Log**: Windowsサービス時

## 要件

- **Go**: 1.21以上
- **OS**: Windows (日本語フォント対応)
- **ポート**: 8081 (デフォルト)

## 技術スタック

- **言語**: Go 1.21+
- **PDF生成**: カスタムReportLabスタイルライブラリ
- **HTTP Framework**: 標準 `net/http`
- **CI/CD**: GitHub Actions
- **テスト**: 91%+ カバレッジ
- **フォント**: Windows標準日本語フォント (yumin.ttf)

## 主要な構造体

### Item
```go
type Item struct {
    Car         string   `json:"car"`
    Name        string   `json:"name"`
    Purpose     *string  `json:"purpose"`
    StartDate   *string  `json:"startDate"`
    EndDate     *string  `json:"endDate"`
    Price       int      `json:"price"`
    Tax         *float64 `json:"tax"`
    Description *string  `json:"description"`
    Ryohi       []Ryohi  `json:"ryohi"`
    Office      *string  `json:"office"`
    PayDay      *string  `json:"payDay"`
}
```

### Ryohi
```go
type Ryohi struct {
    Date   *string `json:"date"`
    Dest   *string `json:"dest"`
    Detail *string `json:"detail"`
    Kukan  *string `json:"kukan"`
    Price  int     `json:"price"`
    Vol    *int    `json:"vol"`
}
```

## ライセンス

このプロジェクトは組織内部での使用を想定しています。

## 貢献

1. このリポジトリをフォーク
2. フィーチャーブランチを作成 (`git checkout -b feature/amazing-feature`)
3. 変更をコミット (`git commit -m 'Add amazing feature'`)
4. ブランチにプッシュ (`git push origin feature/amazing-feature`)
5. プルリクエストを作成

## サポート

問題や質問がある場合は、GitHubのIssuesページで報告してください。
