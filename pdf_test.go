package main

import (
	"testing"
)

func TestNewReportLabStylePdfClientCreation(t *testing.T) {
	// PDFクライアントの作成をテスト（空のItemスライスで）
	var items []Item
	client := NewReportLabStylePdfClient(items)

	if client == nil {
		t.Fatal("NewReportLabStylePdfClient() returned nil")
	}

	// 内部のPDFオブジェクトが作成されていることを確認
	// （プライベートフィールドなので直接テストはできないが、後続のメソッド呼び出しで検証）
}

func TestPrintRyohiItemsWithEmptyData(t *testing.T) {
	var items []Item
	client := NewReportLabStylePdfClient(items)

	// 空のデータでテスト
	var emptyItems []Ryohi

	// printRyohiItemsは戻り値がないので、パニックしないことを確認
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("printRyohiItems with empty data should not panic: %v", r)
		}
	}()

	client.printRyohiItems(emptyItems)
}

func TestPrintRyohiItemsWithSingleItem(t *testing.T) {
	var items []Item
	client := NewReportLabStylePdfClient(items)

	// 単一アイテムでテスト
	ryohiItems := []Ryohi{
		{
			Date:   StringPtr("01/15"),
			Dest:   StringPtr("東京"),
			Detail: []string{"会議"},
			Kukan:  StringPtr("東京駅　大阪駅"),
			Price:  IntPtr(5000),
			Vol:    Float64Ptr(1.5),
		},
	}

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("printRyohiItems with single item should not panic: %v", r)
		}
	}()

	client.printRyohiItems(ryohiItems)
}

func TestPrintRyohiItemsWithMultipleItems(t *testing.T) {
	var items []Item
	client := NewReportLabStylePdfClient(items)

	// 複数アイテムでテスト
	ryohiItems := []Ryohi{
		{
			Date:   StringPtr("01/15"),
			Dest:   StringPtr("東京"),
			Detail: []string{"会議", "資料作成"},
			Kukan:  StringPtr("東京駅　大阪駅"),
			Price:  IntPtr(5000),
			Vol:    Float64Ptr(1.5),
		},
		{
			Date:   StringPtr("01/16"),
			Dest:   StringPtr("大阪"),
			Detail: []string{"営業活動", "打ち合わせ"},
			Kukan:  StringPtr("大阪駅　京都駅"),
			Price:  IntPtr(8000),
			Vol:    Float64Ptr(2.0),
		},
	}

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("printRyohiItems with multiple items should not panic: %v", r)
		}
	}()

	client.printRyohiItems(ryohiItems)
}

func TestPrintRyohiItemsWithNilFields(t *testing.T) {
	var items []Item
	client := NewReportLabStylePdfClient(items)

	// nilフィールドを含むアイテムでテスト
	ryohiItems := []Ryohi{
		{
			Date:   nil,
			Dest:   nil,
			Detail: []string{},
			Kukan:  nil,
			Price:  nil,
			Vol:    nil,
		},
	}

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("printRyohiItems with nil fields should not panic: %v", r)
		}
	}()

	client.printRyohiItems(ryohiItems)
}

func TestPrintRyohiItemsWithLongData(t *testing.T) {
	var items []Item
	client := NewReportLabStylePdfClient(items)

	// 長いデータでテスト（14行の表示システムをテスト）
	ryohiItems := []Ryohi{
		{
			Date:   StringPtr("01/15"),
			Dest:   StringPtr("東京都千代田区丸の内一丁目"),
			Detail: []string{"重要会議", "新規プロジェクト打ち合わせ", "資料作成", "報告書提出", "クライアント面談", "進捗報告", "次回計画"},
			Kukan:  StringPtr("博多駅　小倉駅　新山口駅　広島駅　岡山駅　新大阪駅　京都駅　東京駅　上野駅　大宮駅"),
			Price:  IntPtr(25000),
			Vol:    Float64Ptr(5.5),
		},
	}

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("printRyohiItems with long data should not panic: %v", r)
		}
	}()

	client.printRyohiItems(ryohiItems)
}

func TestPrintRyohiItemsWithMaxItems(t *testing.T) {
	var items []Item
	client := NewReportLabStylePdfClient(items)

	// 14個のアイテム（最大表示数）でテスト
	var ryohiItems []Ryohi
	for i := 1; i <= 14; i++ {
		ryohiItems = append(ryohiItems, Ryohi{
			Date:   StringPtr("01/15"),
			Dest:   StringPtr("東京"),
			Detail: []string{"会議"},
			Kukan:  StringPtr("東京駅　大阪駅"),
			Price:  IntPtr(5000),
			Vol:    Float64Ptr(1.5),
		})
	}

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("printRyohiItems with 14 items should not panic: %v", r)
		}
	}()

	client.printRyohiItems(ryohiItems)
}

func TestPdfGenerationIntegration(t *testing.T) {
	// 統合テスト：実際のPDF生成プロセス全体をテスト
	// 実際のデータに近いテストケース
	items := []Item{
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
					Date:   StringPtr("01/15"),
					Dest:   StringPtr("東京都千代田区"),
					Detail: []string{"重要会議", "新規プロジェクト打ち合わせ", "資料作成"},
					Kukan:  StringPtr("博多駅　小倉駅　新山口駅　広島駅　岡山駅　新大阪駅　京都駅　東京駅"),
					Price:  IntPtr(15000),
					Vol:    Float64Ptr(3.5),
				},
				{
					Date:   StringPtr("01/16"),
					Dest:   StringPtr("大阪府大阪市"),
					Detail: []string{"クライアント面談", "営業活動"},
					Kukan:  StringPtr("東京駅　新大阪駅"),
					Price:  IntPtr(12000),
					Vol:    Float64Ptr(2.0),
				},
				{
					Date:   StringPtr("01/17"),
					Dest:   StringPtr("京都府京都市"),
					Detail: []string{"研修参加"},
					Kukan:  StringPtr("大阪駅　京都駅"),
					Price:  IntPtr(3000),
					Vol:    Float64Ptr(1.0),
				},
			},
		},
	}

	// PDF生成
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("PDF generation should not panic: %v", r)
		}
	}()

	client := NewReportLabStylePdfClient(items)
	if client == nil {
		t.Fatal("Failed to create PDF client")
	}

	t.Log("PDF generation integration test completed successfully")
}

// ベンチマークテスト
func BenchmarkNewReportLabStylePdfClient(b *testing.B) {
	var items []Item
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		client := NewReportLabStylePdfClient(items)
		_ = client
	}
}

func BenchmarkPrintRyohiItems(b *testing.B) {
	ryohiItems := []Ryohi{
		{
			Date:   StringPtr("01/15"),
			Dest:   StringPtr("東京"),
			Detail: []string{"会議", "資料作成"},
			Kukan:  StringPtr("東京駅　大阪駅"),
			Price:  IntPtr(5000),
			Vol:    Float64Ptr(1.5),
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var items []Item
		client := NewReportLabStylePdfClient(items)
		client.printRyohiItems(ryohiItems)
	}
}

func BenchmarkFullPdfGeneration(b *testing.B) {
	items := []Item{
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
					Date:   StringPtr("01/15"),
					Dest:   StringPtr("東京都千代田区"),
					Detail: []string{"重要会議", "新規プロジェクト打ち合わせ", "資料作成"},
					Kukan:  StringPtr("博多駅　小倉駅　新山口駅　広島駅　岡山駅　新大阪駅　京都駅　東京駅"),
					Price:  IntPtr(15000),
					Vol:    Float64Ptr(3.5),
				},
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		client := NewReportLabStylePdfClient(items)
		_ = client
	}
}
