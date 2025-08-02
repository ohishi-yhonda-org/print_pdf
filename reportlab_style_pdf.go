package main

import (
	"fmt"
	"os"
	"time"

	"github.com/jung-kurt/gofpdf"
)

// ReportLabStylePdfClient - ReportLabスタイルのPDF生成クライアント
type ReportLabStylePdfClient struct {
	pdf *gofpdf.Fpdf
	// ページサイズ定数 (A5横向き: 210mm x 148mm)
	pageSizeX float64
	pageSizeY float64
	// マージン
	lX, tY, rX, bY float64
	// テーブル位置とサイズ
	col1, col2, col3, col4, col5 float64
	rowH1                        float64
	tblTY                        float64
}

// NewReportLabStylePdfClient - ReportLabスタイルのPDFクライアントを作成
func NewReportLabStylePdfClient(data []Item) *ReportLabStylePdfClient {
	fmt.Println("Creating ReportLab Style PDF client...")

	// A5横向きでPDFを初期化 (210mm x 148mm)
	pdf := gofpdf.New("L", "mm", "A5", "")

	client := &ReportLabStylePdfClient{
		pdf:       pdf,
		pageSizeX: 210, // A5横向きの幅
		pageSizeY: 148, // A5横向きの高さ
	}

	// マージンを設定 (ReportLabの座標系に合わせて)
	client.lX = 10.0  // 左マージン
	client.tY = 138.0 // 上マージン
	client.rX = 200.0 // 右マージン
	client.bY = 10.0  // 下マージン

	// テーブルサイズ設定（mm単位）
	client.col1 = 30.0
	client.col2 = 25.0
	client.col3 = 28.75
	client.col4 = 30.0
	client.col5 = 30.0
	client.rowH1 = 15.0
	client.tblTY = 95.0 // 基本情報テーブルの位置

	// Windowsフォントを設定
	if err := client.setupWindowsFont(); err != nil {
		fmt.Printf("Windowsフォント設定エラー: %v\n", err)
		fmt.Println("標準フォントで継続...")
	}

	filePath := "travel_expense_reportlab_style.pdf"
	fmt.Println("ReportLab Style output file:", filePath)

	// 各アイテムを処理
	for index, item := range data {
		if index != 0 {
			pdf.AddPage()
		}
		pdf.AddPage()
		fmt.Printf("Processing ReportLab Style index: %d\n", index)

		client.drawLine()
		client.printItem(item)
	}

	// PDFを保存
	err := pdf.OutputFileAndClose(filePath)
	if err != nil {
		fmt.Printf("Error saving ReportLab Style PDF: %v\n", err)
		return nil
	} else {
		fmt.Println("ReportLab Style PDF saved successfully!")

		if fileInfo, err := os.Stat(filePath); err == nil {
			fmt.Printf("File size: %d bytes\n", fileInfo.Size())
			if fileInfo.Size() > 1000 {
				fmt.Println("PDF appears to have proper content!")
			}
		}
	}

	return client
}

// setupWindowsFont - Windowsフォントを設定
func (c *ReportLabStylePdfClient) setupWindowsFont() error {
	fmt.Println("Setting up Windows standard fonts for ReportLab style...")

	windowsFonts := []struct {
		name string
		path string
	}{
		{"yumin", "C:/Windows/Fonts/yumin.ttf"},
		{"yugothm", "C:/Windows/Fonts/yugothm.ttf"},
		{"meiryo", "C:/Windows/Fonts/meiryo.ttf"},
	}

	for _, font := range windowsFonts {
		if _, err := os.Stat(font.path); err == nil {
			fmt.Printf("Found Windows font: %s at %s\n", font.name, font.path)
			c.pdf.AddUTF8Font(font.name, "", font.path)
			fmt.Printf("Successfully added Windows font: %s\n", font.name)
			return nil
		}
	}

	fmt.Println("Windows標準フォントが見つからない、標準フォントを使用...")
	return fmt.Errorf("no Windows fonts found")
}

// truncateText - テキストを指定された文字数で切り詰める
func truncateText(text string, maxLength int) string {
	if len([]rune(text)) > maxLength {
		textRunes := []rune(text)
		return string(textRunes[:maxLength])
	}
	return text
}

// drawLine - ReportLabスタイルのページ枠線とテーブルを描画
func (c *ReportLabStylePdfClient) drawLine() {
	startX := 10.0
	startY := 15.0
	endX := c.pageSizeX - 10.0
	endY := c.pageSizeY - 10.0
	// 外枠を描画
	c.pdf.SetLineWidth(0.5)
	c.pdf.Line(startX, startY, startX, endY) // left
	c.pdf.Line(startX, startY, endX, startY) // top
	c.pdf.Line(endX, startY, endX, endY)     // right
	c.pdf.Line(startX, endY, endX, endY)     // bottom

	// 承認テーブル（右上）
	c.drawApprovalTable()

	// 基本情報テーブル
	c.drawBasicInfoTable()

	// メインデータテーブル
	c.drawMainDataTable()

	// 備考・計テーブル
	c.drawSummaryTable()
}

// drawApprovalTable - 承認テーブルを描画
// drawApprovalTable - 承認テーブルを描画
func (c *ReportLabStylePdfClient) drawApprovalTable() {
	startX := 155.0 // 3列×15mm
	startY := 25.0  // 上から40mm下
	colWidth := 15.0
	rowHeight1 := 5.0
	rowHeight2 := 15.0

	// ヘッダー行
	c.pdf.SetFont("yumin", "", 9)
	headers := []string{"社　長", "会　計", "所　属"}
	for i, header := range headers {
		x := startX + float64(i)*colWidth
		c.pdf.Rect(x, startY, colWidth, rowHeight1, "D")
		textWidth := c.pdf.GetStringWidth(header)
		c.pdf.Text(x+(colWidth-textWidth)/2, startY+4, header)
	}

	// データ行（空）
	for i := 0; i < 3; i++ {
		x := startX + float64(i)*colWidth
		c.pdf.Rect(x, startY+rowHeight1, colWidth, rowHeight2, "D")
	}
}

func (c *ReportLabStylePdfClient) drawBasedata(item Item) {
	startX := 10.0
	startY := 15.0
	// タイトル
	c.pdf.SetFont("yumin", "", 14)
	sent := "出 張 旅 費 日 当 駐 車 料 込 精 算 書"
	c.pdf.Text(startX+13, startY+5, sent)
	// タイトル下線（2本）
	titleWidth := c.pdf.GetStringWidth(sent)
	c.pdf.Line(startX+13, startY+6, startX+15+titleWidth, startY+6)
	c.pdf.Line(startX+13, startY+7, startX+15+titleWidth, startY+7)

	// 精算日
	if item.PayDay != nil {
		c.pdf.SetFont("yumin", "", 9)
		if t, err := time.Parse("2006-01-02", *item.PayDay); err == nil {
			payDay := t.Format("清算日　2006年 01月 02日")
			c.pdf.Text(startX+100, startY+5, payDay)
		}
	}

	// 所属（右上）
	if item.Office != nil {
		c.pdf.SetFont("yumin", "", 10)
		textWidth := c.pdf.GetStringWidth(*item.Office)
		c.pdf.Text(startX+190-textWidth-2, startY+5, *item.Office)
	}

}

// drawBasicInfoTable - 基本情報テーブルを描画
func (c *ReportLabStylePdfClient) drawBasicInfoTable() {
	startX := 10.0
	startY := 30.0

	// 出発・帰着日（左側）
	rowHeight := 3.5
	diffStartY := 3.0
	c.pdf.SetFont("yumin", "", 9)
	c.pdf.Text(startX+1, startY+diffStartY, "出発")
	c.pdf.Text(startX+2, startY+diffStartY+rowHeight, "　　月　　日")
	c.pdf.Text(startX+1, startY+diffStartY+rowHeight*2, "帰着")
	c.pdf.Text(startX+2, startY+diffStartY+rowHeight*3, "　　月　　日")

	// テーブルヘッダー
	c.pdf.SetFont("yumin", "", 9)
	headers := []string{"", "出張目的", "車両No.", "氏　名", "サイン"}
	colWidths := []float64{31, 25, 28.75, 30, 30} // mm単位

	currentX := startX
	for i, header := range headers {
		c.pdf.Rect(currentX, startY, colWidths[i], 15, "D")
		if header != "" {
			// textWidth := c.pdf.GetStringWidth(header)
			c.pdf.Text(currentX+1, startY+4, header)
		}
		currentX += colWidths[i]
	}
}

// drawMainDataTable - メインデータテーブルを描画
func (c *ReportLabStylePdfClient) drawMainDataTable() {
	startX := 10.0
	startY := 45.0 // タイトルから85mm下

	// 正確な列幅（元の画像に合わせて）
	colWidths := []float64{10, 17, 40, 30, 15, 15, 15, 25, 23}
	rowHeight := 10.0
	headerHeight := 4.0

	// ヘッダー
	c.pdf.SetFont("yumin", "", 8)
	headers := []string{"日付", "行　先", "摘　　要", "区　　間", "交通機関", "運　賃", "特別料金", "旅費日当", "計"}

	currentX := startX
	for i, header := range headers {
		c.pdf.Rect(currentX, startY, colWidths[i], headerHeight, "D")
		// 中央揃え
		textWidth := c.pdf.GetStringWidth(header)
		c.pdf.Text(currentX+(colWidths[i]-textWidth)/2, startY+3, header)
		currentX += colWidths[i]
	}

	// データ行（7行）
	for row := 0; row < 7; row++ {
		currentX = startX
		currentY := startY + headerHeight + float64(row)*rowHeight

		for col := 0; col < len(colWidths); col++ {
			//摘要欄は上下の枠線を描画しない
			if col == 2 {
				// 摘要欄は左右の線のみ描画
				c.pdf.Line(currentX, currentY, currentX, currentY+rowHeight)                               // 左線
				c.pdf.Line(currentX+colWidths[col], currentY, currentX+colWidths[col], currentY+rowHeight) // 右線
			} else {
				c.pdf.Rect(currentX, currentY, colWidths[col], rowHeight, "D")
			}
			currentX += colWidths[col]
		}
	}
}

// drawSummaryTable - 備考・計テーブルを描画
func (c *ReportLabStylePdfClient) drawSummaryTable() {
	startX := 10.0
	startY := 119.0 // 底辺から18mm上

	colWidths := []float64{145, 45} // 備考欄と計欄
	rowHeight := 19.0

	c.pdf.SetFont("yumin", "", 8)
	headers := []string{"備考", "計"}

	currentX := startX
	for i, header := range headers {
		c.pdf.Rect(currentX, startY, colWidths[i], rowHeight, "D")
		c.pdf.Text(currentX+2, startY+4, header)
		currentX += colWidths[i]
	}
}

// printItem - アイテム情報を印刷
func (c *ReportLabStylePdfClient) printItem(item Item) {

	c.drawBasedata(item)
	// 出発日
	startX := 14.0
	startY := 36.8

	if item.StartDate != nil {
		c.pdf.SetFont("yumin", "", 10)
		if t, err := time.Parse("2006-01-02", *item.StartDate); err == nil {
			startDate := t.Format("01　 02")
			c.pdf.Text(startX, startY, startDate)
		}
	}

	// 帰着日
	if item.EndDate != nil {
		c.pdf.SetFont("yumin", "", 10)
		if t, err := time.Parse("2006-01-02", *item.EndDate); err == nil {
			endDate := t.Format("01　 02")
			c.pdf.Text(startX, startY+7, endDate)
		}
	}

	// 出張目的
	if item.Purpose != nil {
		c.pdf.SetFont("yumin", "", 10)
		c.pdf.Text(startX+32, startY+7, *item.Purpose)
	}

	// 車両
	if item.Car != "" {
		c.pdf.SetFont("yumin", "", 10)
		c.pdf.Text(startX+52, startY+7, item.Car)
	}

	// 氏名
	if item.Name != "" {
		c.pdf.SetFont("yumin", "", 10)
		c.pdf.Text(startX+85, startY+7, item.Name)
	}

	// 合計金額（上部の計欄）
	c.pdf.SetFont("yumin", "", 12)
	priceStr := FormatPrice(item.Price)
	textWidth := c.pdf.GetStringWidth(priceStr)
	c.pdf.Text(c.rX-textWidth-5, c.tY-12, priceStr)

	// 旅費データを処理
	c.printRyohiItems(item.Ryohi)
}

// printRyohiItems - 旅費データを印刷
func (c *ReportLabStylePdfClient) printRyohiItems(ryohiList []Ryohi) {
	startX := 10.0
	startY := 47.0 // メインテーブルのヘッダー下から開始
	colWidths := []float64{10, 17, 40, 30, 15, 15, 15, 25, 23}
	rowHeight := 10.0

	c.pdf.SetFont("yumin", "", 10)

	currentRow := 0
	for i, ryohi := range ryohiList {
		if currentRow >= 14 { // 最大14行まで表示
			break
		}

		// 旅費データを印刷用に準備（摘要10文字、区間22文字制限）
		printData := prepareRyohiForPrint(ryohi, 10, 22)

		// 残り行数をチェック（14行まで対応）
		remainingRows := 14 - currentRow
		actualRows := printData.MaxRows
		if actualRows > remainingRows {
			actualRows = remainingRows
		}

		// 実際に描画する行をカウント
		drawnRows := 0

		// 各行を印刷（コンテンツがある行のみ）
		for row := 0; row < actualRows && drawnRows < remainingRows; row++ {
			// この行にコンテンツがあるかチェック
			if !printData.hasContentInRow(row) {
				continue // 空行はスキップ
			}

			// 表の物理行を計算（14行を7行に配置）
			logicalRow := currentRow + drawnRows // 現在の論理行位置
			physicalRow := logicalRow / 2        // 実際のPDF上の表の行
			subRow := logicalRow % 2             // その行の上半分(0)か下半分(1)か
			yOffset := float64(subRow) * 5.0     // 上半分は+0mm、下半分は+5mm

			currentY := startY + float64(physicalRow)*rowHeight + yOffset
			currentX := startX

			// 日付
			if row < len(printData.DateLines) && printData.DateLines[row] != "" {
				date := printData.DateLines[row]
				c.pdf.SetFont("yumin", "", 10)
				textWidth := c.pdf.GetStringWidth(date)
				c.pdf.Text(currentX+(colWidths[0]-textWidth)/2, currentY+6, date)
			}
			currentX += colWidths[0]

			// 行先
			if row < len(printData.DestLines) && printData.DestLines[row] != "" {
				dest := printData.DestLines[row]
				textWidth := c.pdf.GetStringWidth(dest)
				c.pdf.Text(currentX+(colWidths[1]-textWidth)/2, currentY+6, dest)
			}
			currentX += colWidths[1]

			// 摘要
			if row < len(printData.DetailLines) && printData.DetailLines[row] != "" {
				detail := printData.DetailLines[row]
				c.pdf.Text(currentX+1, currentY+6, detail)
			}
			currentX += colWidths[2]

			// 区間
			if row < len(printData.KukanLines) && printData.KukanLines[row] != "" {
				kukan := printData.KukanLines[row]
				c.pdf.Text(currentX+1, currentY+6, kukan)
			}
			currentX += colWidths[3]

			// 交通機関（空）
			currentX += colWidths[4]

			// 運賃（空）
			currentX += colWidths[5]

			// 特別料金（空）
			currentX += colWidths[6]

			// 旅費日当
			if row < len(printData.PriceLines) && printData.PriceLines[row] != "" {
				priceStr := printData.PriceLines[row]
				textWidth := c.pdf.GetStringWidth(priceStr)
				c.pdf.Text(currentX+colWidths[7]-textWidth-1, currentY+6, priceStr)
			}
			currentX += colWidths[7]

			// 計
			if row < len(printData.VolLines) && printData.VolLines[row] != "" {
				volStr := printData.VolLines[row]
				textWidth := c.pdf.GetStringWidth(volStr)
				c.pdf.Text(currentX+colWidths[8]-textWidth-1, currentY+6, volStr)
			}

			drawnRows++ // 実際に描画した行数をインクリメント
		}

		currentRow += drawnRows // 実際に描画した行数分だけ進める

		// デバッグ情報
		fmt.Printf("旅費項目 %d: 最大行数=%d, 実際印刷行数=%d, 現在行=%d\n",
			i+1, printData.MaxRows, drawnRows, currentRow)
	}
}
