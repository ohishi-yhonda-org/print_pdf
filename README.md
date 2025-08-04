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

### POST /print
封筒印刷専用エンドポイント（PHPからのマルチパート形式対応）

**リクエスト例 (PHP CakeHTTP):**
```php
$client = new \Cake\Http\Client();
$response = $client->post('http://172.18.21.233:8081/print', [
    'document' => fopen('/var/www/html/files/futo.pdf', 'r'),
    'printer' => 'LBP221-futo'
]);
```

**リクエスト例 (curl):**
```bash
curl -X POST http://localhost:8081/print \
  -F "document=@/path/to/envelope.pdf" \
  -F "printer=LBP221-futo"
```

**レスポンス（成功時）:**
```json
{
  "status": "success",
  "message": "封筒印刷が正常に完了しました",
  "filename": "futo.pdf",
  "printer": "LBP221-futo",
  "printed": true,
  "fileSize": 12345
}
```

**レスポンス（エラー時）:**
```json
{
  "status": "error",
  "message": "印刷エラー: プリンターが見つかりません",
  "filename": "futo.pdf",
  "printer": "LBP221-futo",
  "printed": false
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
# 1. 包括的テストを実行（91%+ カバレッジ）
# 2. 品質チェック（golangci-lint）
# 3. バージョン付きでビルド
# 4. zipアーカイブを作成
# 5. GitHub Releasesにアップロード
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

## リリース手順

### 自動リリース（推奨）

mainブランチにpushすると、GitHub Actionsが自動的にリリースを作成します：

1. **開発とテスト**
   ```bash
   # 開発環境での動作確認
   go run .
   # テスト実行
   go test -v ./...
   ```

2. **コミットとプッシュ**
   ```bash
   git add .
   git commit -m "機能追加: 新しい機能の説明"
   git push origin main
   ```

3. **自動リリース作成**
   - CIが最新のGitタグから次のバージョンを自動生成
   - 例：`v1.0.11` → `v1.0.12`
   - `print_pdf_vX.X.X.zip`が自動作成されます

### 手動リリース

特定のバージョンでリリースしたい場合：

1. **GitHub Actionsの手動実行**
   - リポジトリのActionsタブに移動
   - "CI"ワークフローを選択
   - "Run workflow"をクリック
   - パラメータを設定：
     - `create_release`: `true`
     - `version`: `v1.2.0`（任意のバージョン）

2. **直接タグ作成**
   ```bash
   git tag v1.2.0
   git push origin v1.2.0
   ```

### リリース内容

自動作成されるリリースには以下が含まれます：

- **バイナリファイル**: `print_pdf.exe`
- **サービス管理スクリプト**: `service_manager.bat`
- **ZIPアーカイブ**: `print_pdf_vX.X.X.zip`
- **自動生成されるリリースノート**：
  - セキュリティ最適化情報
  - API エンドポイント一覧
  - インストール手順
  - ウイルス対策ソフト対応ガイド

### 開発版の扱い

- **開発版タグ**: `v1.0.X-dev.YYYYMMDD.HHMM`形式のタグは自動削除される場合があります
- **dev版はリリース作成されません**: コードの`Version = "dev"`の場合、GitHub Releasesは作成されません
- **本番用ビルド**: CIが`-ldflags`でバージョンを自動設定します

### バージョニング規則

- **メジャーバージョン**: 破壊的変更
- **マイナーバージョン**: 新機能追加
- **パッチバージョン**: バグフィックス
- **開発版サフィックス**: `-dev.YYYYMMDD.HHMM`

### トラブルシューティング

**リリースが作成されない場合：**
1. CIログを確認
2. `-dev.`がバージョンに含まれていないか確認
3. GitHub Actionsの権限を確認

**古いdev版タグの削除：**
```bash
# ローカルのdevタグ削除
git tag -d v1.0.X-dev.YYYYMMDD.HHMM

# リモートのdevタグ削除
git push origin --delete v1.0.X-dev.YYYYMMDD.HHMM
```

## サポート

問題や質問がある場合は、GitHubのIssuesページで報告してください。
