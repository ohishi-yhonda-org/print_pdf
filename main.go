package main

import (
	"fmt"
)

func main() {
	fmt.Println("PDF生成システム - Go版 (Windowsフォント対応)")

	// テストデータの作成
	testData := []Item{
		{
			Car:       "長崎100か4105",
			Name:      "松本　俊之",
			Purpose:   StringPtr("営業"),
			StartDate: StringPtr("2024-12-17"),
			EndDate:   StringPtr("2024-12-28"),
			Price:     86900,
			PayDay:    StringPtr("2025-01-06"),
			Office:    StringPtr("本社㈲"),
			Ryohi: []Ryohi{
				{
					Date:   StringPtr("2024-12-17"),
					Dest:   StringPtr("滋賀"),
					Detail: []string{"長崎", "大阪", "滋賀", "熊本"},
					Kukan:  StringPtr("大型車_長崎_滋賀_往復_荷配　"),
					Price:  IntPtr(21000),
					Vol:    Float64Ptr(7.0),
				},
				{
					Date:   StringPtr("2024-12-18"),
					Dest:   StringPtr("福岡"),
					Detail: []string{"長崎", "福岡"},
					Kukan:  StringPtr("大型車_長崎_福岡_片荷_荷配　"),
					Price:  IntPtr(2300),
					Vol:    Float64Ptr(2.0),
				},
				{
					Date:   StringPtr("2024-12-19"),
					Dest:   StringPtr("三重"),
					Detail: []string{"鹿児島", "三重", "滋賀", "長崎"},
					Kukan:  StringPtr("大型車_長崎_三重_往復_荷配 |加算額_大型車_片荷_A　"),
					Price:  IntPtr(26300),
					Vol:    Float64Ptr(7.0),
				},
				{
					Date:   StringPtr("2024-12-20"),
					Dest:   StringPtr("福岡"),
					Detail: []string{"長崎", "福岡"},
					Kukan:  StringPtr("大型車_長崎_福岡_片荷_荷配　"),
					Price:  IntPtr(5300),
					Vol:    Float64Ptr(3.0),
				},
				{
					Date:   StringPtr("2024-12-21"),
					Dest:   StringPtr("埼玉"),
					Detail: []string{"大阪", "埼玉", "茨城", "福岡"},
					Kukan:  StringPtr("大型車_福岡_埼玉_往復_荷配 |加算額_大型車_往復_A　"),
					Price:  IntPtr(32000),
					Vol:    Float64Ptr(7.0),
				},
			},
		},
	}

	fmt.Println("ReportLabスタイルPDF生成を開始...")

	// ReportLabスタイル版を実行
	reportlabClient := NewReportLabStylePdfClient(testData)
	if reportlabClient != nil {
		fmt.Println("ReportLabスタイルPDF生成完了！")
	} else {
		fmt.Println("ReportLabスタイルPDF生成に失敗しました")
	}

	fmt.Println("\n処理完了！")
}
