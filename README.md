# PDF出張旅費精算書生成システム（Go版）

このプロジェクトは、PythonのReportLabライブラリを使用したPDF生成システムをGoに移植したものです。

## 機能

- 出張旅費精算書のPDF生成
- 旅費データの自動レイアウト
- 複数ページ対応
- 日本語フォント対応（Arial使用）

## 必要な依存関係

- Go 1.24.4以上
- github.com/go-pdf/fpdf
- github.com/golang/freetype

## インストール

```bash
go mod tidy
go get github.com/go-pdf/fpdf
go get github.com/golang/freetype
```

## 使用方法

```bash
go build -o print_pdf.exe
.\print_pdf.exe
```

## ファイル構成

- `main.go` - メインプログラム
- `models.go` - データ構造とヘルパー関数
- `pdf_client.go` - PDF生成クライアント
- `tests.go` - テスト関数

## 出力

PDFファイルは実行ディレクトリ（current directory）に `sample.pdf` として出力されます。

## 主要な構造体

### Item
- `Car`: 車両番号
- `Name`: 氏名
- `Purpose`: 出張目的
- `StartDate`: 出発日
- `EndDate`: 帰着日
- `Price`: 合計金額
- `Ryohi`: 旅費明細（Ryohi配列）
- `Office`: 所属
- `PayDay`: 精算日

### Ryohi
- `Date`: 日付
- `Dest`: 行先
- `Detail`: 摘要詳細
- `Kukan`: 区間
- `Price`: 金額
- `Vol`: 数量

## テスト

プログラム実行時に自動的にテストが実行されます。以下のテストが含まれています：

- `testPdfClientCalDetailRow`: 摘要行計算のテスト
- `testPrintRyohiItemsOtherBase`: 特定ケースのテスト

## 注意事項

- 現在は基本的なレイアウトのみ実装
- 日本語フォントはビルトインのArialフォントを使用
- PDF生成時にコンソールにデバッグ情報が出力されます

## 元のPythonコードとの違い

- ReportLabの代わりにgofpdfを使用
- Pydanticの代わりにGo structsを使用
- 型安全性の向上
- メモリ効率の改善

## 今後の改善予定

- 日本語フォントファイルの組み込み
- より詳細なPDFレイアウト
- エラーハンドリングの強化
- テストケースの拡充
