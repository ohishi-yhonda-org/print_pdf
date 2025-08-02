package main

import (
	"fmt"
	"regexp"
	"strings"
)

// TextWrapResult - テキスト折り返し結果
type TextWrapResult struct {
	Lines    []string
	RowCount int
}

// wrapDetail - 摘要テキストを指定行数で折り返し
func wrapDetail(details []string, maxLen int) TextWrapResult {
	if len(details) == 0 {
		return TextWrapResult{Lines: []string{}, RowCount: 0}
	}

	var result []string
	currentLine := ""

	for _, detail := range details {
		// 区切り文字を考慮した新しい行の長さ
		separator := ""
		if currentLine != "" {
			separator = "、"
		}
		newLineLength := len([]rune(currentLine)) + len([]rune(separator)) + len([]rune(detail))

		if newLineLength <= maxLen {
			// 全体が収まる場合
			currentLine += separator + detail
		} else {
			// 収まらない場合、現在の行が空でなければ確定して次の行に移る
			if currentLine != "" {
				result = append(result, currentLine)
				currentLine = ""
			}

			// 新しい詳細項目を次の行に配置
			if len([]rune(detail)) > maxLen {
				// 詳細項目自体が最大長を超える場合は切り詰め
				currentLine = string([]rune(detail)[:maxLen])
			} else {
				currentLine = detail
			}
		}
	}

	// 最後の行を処理（空でない場合のみ）
	if currentLine != "" {
		result = append(result, currentLine)
	}

	// 空行を除去
	var filteredResult []string
	for _, line := range result {
		if strings.TrimSpace(line) != "" {
			filteredResult = append(filteredResult, line)
		}
	}

	return TextWrapResult{
		Lines:    filteredResult,
		RowCount: len(filteredResult),
	}
}

// wrapKukan - 区間テキストを指定行数で折り返し
func wrapKukan(kukan string, maxLen int) TextWrapResult {
	if kukan == "" {
		return TextWrapResult{Lines: []string{""}, RowCount: 1}
	}

	// 特殊な文字列を置換
	kukan = strings.ReplaceAll(kukan, "_九州外空車適用", "　九州外空車適用")
	kukan = strings.ReplaceAll(kukan, "適用*   追加", "適用*　追加")

	// 区切り文字で分割 (全角スペース、｜、半角スペース+|、など)
	re := regexp.MustCompile(`[　｜]| \||\|`)
	parts := re.Split(kukan, -1)

	var result []string
	currentLine := ""
	currentCount := 0

	for _, part := range parts {
		partLen := len([]rune(part))

		if currentCount != 0 && currentCount+partLen == maxLen {
			// ちょうど最大長になる場合
			result = append(result, currentLine+part)
			currentLine = ""
			currentCount = 0
		} else if partLen == maxLen && currentLine == "" {
			// 単体で最大長の場合
			result = append(result, part)
			currentCount = 0
		} else if partLen > maxLen {
			// 最大長を超える場合
			result = append(result, "exceed*")
			currentCount = 0
		} else if currentCount+partLen+1 > maxLen {
			// 現在行に追加すると最大長を超える場合
			if currentLine != "" {
				result = append(result, currentLine)
			}
			currentLine = part + "　"
			currentCount = partLen + 1
		} else {
			// 現在行に追加できる場合
			currentCount += partLen + 1
			currentLine += part + "　"
		}
	}

	// 最後の行を処理
	if currentCount != 0 {
		result = append(result, currentLine)
	}

	// 前後の全角スペースを削除
	for i, line := range result {
		line = strings.ReplaceAll(line, " ", "　")
		line = strings.TrimPrefix(line, "　")
		line = strings.TrimSuffix(line, "　")
		result[i] = line
	}

	return TextWrapResult{
		Lines:    result,
		RowCount: len(result),
	}
}

// alignRows - 他のデータ項目を最大行数に合わせて配列を調整
func alignRows(date, dest *string, price *int, vol *float64, maxRows int) ([]string, []string, []string, []string) {
	dateArr := make([]string, maxRows)
	destArr := make([]string, maxRows)
	priceArr := make([]string, maxRows)
	volArr := make([]string, maxRows)

	// 最初の行に実際の値を設定
	if date != nil {
		// YYYY-MM-DD形式からMM/DD形式に変換
		dateStr := *date
		if len(dateStr) >= 10 && dateStr[4] == '-' && dateStr[7] == '-' {
			// YYYY-MM-DD形式の場合
			month := dateStr[5:7]
			day := dateStr[8:10]
			dateArr[0] = month + "/" + day
		} else {
			// そのまま使用
			dateArr[0] = dateStr
		}
	}
	if dest != nil {
		destArr[0] = *dest
	}
	if price != nil {
		priceArr[0] = FormatPrice(*price)
	}
	if vol != nil {
		volArr[0] = fmt.Sprintf("%.1f", *vol)
	}

	// 残りの行は空文字列で埋める
	for i := 1; i < maxRows; i++ {
		if dateArr[i] == "" {
			dateArr[i] = ""
		}
		if destArr[i] == "" {
			destArr[i] = ""
		}
		if priceArr[i] == "" {
			priceArr[i] = ""
		}
		if volArr[i] == "" {
			volArr[i] = ""
		}
	}

	return dateArr, destArr, priceArr, volArr
}

// extendToMaxRows - 配列を最大行数まで拡張（空行は追加しない）
func extendToMaxRows(lines []string, maxRows int) []string {
	// 空行を除去
	var filteredLines []string
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			filteredLines = append(filteredLines, line)
		}
	}

	// 必要に応じて最大行数まで拡張（ただし空文字列は追加しない）
	if len(filteredLines) < maxRows {
		result := make([]string, maxRows)
		copy(result, filteredLines)
		// 残りの行は空文字列で埋める（PDFで条件チェックされる）
		for i := len(filteredLines); i < maxRows; i++ {
			result[i] = ""
		}
		return result
	}

	return filteredLines
}

// RyohiPrintData - 旅費印刷用データ
type RyohiPrintData struct {
	DateLines   []string
	DestLines   []string
	DetailLines []string
	KukanLines  []string
	PriceLines  []string
	VolLines    []string
	MaxRows     int
}

// hasContentInRow - 指定した行にコンテンツがあるかチェック
func (r *RyohiPrintData) hasContentInRow(row int) bool {
	if row >= len(r.DateLines) && row >= len(r.DestLines) &&
		row >= len(r.DetailLines) && row >= len(r.KukanLines) &&
		row >= len(r.PriceLines) && row >= len(r.VolLines) {
		return false
	}

	// いずれかの列にコンテンツがあればtrue
	if row < len(r.DateLines) && strings.TrimSpace(r.DateLines[row]) != "" {
		return true
	}
	if row < len(r.DestLines) && strings.TrimSpace(r.DestLines[row]) != "" {
		return true
	}
	if row < len(r.DetailLines) && strings.TrimSpace(r.DetailLines[row]) != "" {
		return true
	}
	if row < len(r.KukanLines) && strings.TrimSpace(r.KukanLines[row]) != "" {
		return true
	}
	if row < len(r.PriceLines) && strings.TrimSpace(r.PriceLines[row]) != "" {
		return true
	}
	if row < len(r.VolLines) && strings.TrimSpace(r.VolLines[row]) != "" {
		return true
	}

	return false
}

// prepareRyohiForPrint - 旅費データを印刷用に準備
func prepareRyohiForPrint(ryohi Ryohi, maxDetailLen, maxKukanLen int) RyohiPrintData {
	// 摘要を折り返し
	detailResult := TextWrapResult{Lines: []string{""}, RowCount: 1}
	if len(ryohi.Detail) > 0 {
		detailResult = wrapDetail(ryohi.Detail, maxDetailLen)
	}

	// 区間を折り返し
	kukanResult := TextWrapResult{Lines: []string{""}, RowCount: 1}
	if ryohi.Kukan != nil {
		kukanResult = wrapKukan(*ryohi.Kukan, maxKukanLen)
	}

	// 最大行数を決定
	maxRows := detailResult.RowCount
	if kukanResult.RowCount > maxRows {
		maxRows = kukanResult.RowCount
	}

	// 他のデータを最大行数に合わせる
	dateLines, destLines, priceLines, volLines := alignRows(
		ryohi.Date, ryohi.Dest, ryohi.Price, ryohi.Vol, maxRows)

	// すべての配列を最大行数に拡張
	detailLines := extendToMaxRows(detailResult.Lines, maxRows)
	kukanLines := extendToMaxRows(kukanResult.Lines, maxRows)

	return RyohiPrintData{
		DateLines:   dateLines,
		DestLines:   destLines,
		DetailLines: detailLines,
		KukanLines:  kukanLines,
		PriceLines:  priceLines,
		VolLines:    volLines,
		MaxRows:     maxRows,
	}
}
