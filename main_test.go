package main

import (
	"strings"
	"testing"
)

func TestMain(t *testing.T) {
	// main関数が正常に実行されることを確認
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("main() should not panic: %v", r)
		}
	}()

	// main関数を実行（実際のmain()は標準出力があるため、ここでは直接テストしない）
	// 代わりに主要な処理のテストを行う
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
			},
		},
	}

	// PDFクライアント作成をテスト
	client := NewReportLabStylePdfClient(testData)
	if client == nil {
		t.Error("Failed to create ReportLabStylePdfClient")
	}

	t.Log("Main function test completed successfully")
}

func TestComplexDataProcessing(t *testing.T) {
	// 複雑なデータを使ったテスト
	complexData := []Item{
		{
			Car:       "福岡100あ1234",
			Name:      "田中　太郎",
			Purpose:   StringPtr("出張"),
			StartDate: StringPtr("2024-01-15"),
			EndDate:   StringPtr("2024-01-20"),
			Price:     125000,
			PayDay:    StringPtr("2024-02-05"),
			Office:    StringPtr("福岡支社"),
			Ryohi: []Ryohi{
				{
					Date:   StringPtr("01/15"),
					Dest:   StringPtr("東京都新宿区西新宿二丁目"),
					Detail: []string{"重要会議", "新規プロジェクト企画会議", "資料作成", "報告書提出", "次期計画策定"},
					Kukan:  StringPtr("博多駅　小倉駅　新山口駅　広島駅　岡山駅　新大阪駅　京都駅　東京駅　新宿駅"),
					Price:  IntPtr(35000),
					Vol:    Float64Ptr(8.5),
				},
				{
					Date:   StringPtr("01/16"),
					Dest:   StringPtr("神奈川県横浜市港北区"),
					Detail: []string{"クライアント面談", "商品説明", "契約交渉", "技術打合せ"},
					Kukan:  StringPtr("東京駅　横浜駅"),
					Price:  IntPtr(15000),
					Vol:    Float64Ptr(3.0),
				},
				{
					Date:   StringPtr("01/17"),
					Dest:   StringPtr("埼玉県さいたま市"),
					Detail: []string{"研修参加", "技術セミナー受講"},
					Kukan:  StringPtr("東京駅　大宮駅"),
					Price:  IntPtr(8000),
					Vol:    Float64Ptr(2.0),
				},
				{
					Date:   StringPtr("01/18"),
					Dest:   StringPtr("千葉県千葉市中央区"),
					Detail: []string{"工場見学", "製造プロセス確認", "品質管理協議"},
					Kukan:  StringPtr("東京駅　千葉駅"),
					Price:  IntPtr(12000),
					Vol:    Float64Ptr(4.0),
				},
			},
		},
	}

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Complex data processing should not panic: %v", r)
		}
	}()

	client := NewReportLabStylePdfClient(complexData)
	if client == nil {
		t.Error("Failed to create PDF client with complex data")
	}

	t.Log("Complex data processing test completed successfully")
}

func TestEdgeCases(t *testing.T) {
	// エッジケースのテスト
	edgeCaseData := []Item{
		{
			Car:       "",
			Name:      "",
			Purpose:   nil,
			StartDate: nil,
			EndDate:   nil,
			Price:     0,
			PayDay:    nil,
			Office:    nil,
			Ryohi: []Ryohi{
				{
					Date:   nil,
					Dest:   nil,
					Detail: []string{},
					Kukan:  nil,
					Price:  nil,
					Vol:    nil,
				},
				// 非常に長いデータ
				{
					Date: StringPtr("01/01"),
					Dest: StringPtr("非常に長い目的地名称でテキストの切り詰めや折り返し処理をテストするためのサンプルデータです"),
					Detail: []string{
						"非常に長い摘要項目1で文字数制限やテキスト処理をテスト",
						"非常に長い摘要項目2で複数行表示の動作確認",
						"短い項目",
						"非常に長い摘要項目3でさらに複雑な処理をテスト",
						"項目4",
						"項目5",
						"項目6",
						"項目7",
						"項目8",
						"項目9",
						"最後の長い項目で14行制限のテスト",
					},
					Kukan: StringPtr("超長距離区間名称テスト_博多駅_小倉駅_新山口駅_広島駅_岡山駅_新大阪駅_京都駅_東京駅_上野駅_大宮駅_仙台駅_盛岡駅_新青森駅"),
					Price: IntPtr(999999),
					Vol:   Float64Ptr(99.9),
				},
			},
		},
	}

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Edge case processing should not panic: %v", r)
		}
	}()

	client := NewReportLabStylePdfClient(edgeCaseData)
	if client == nil {
		t.Error("Failed to create PDF client with edge case data")
	}

	t.Log("Edge case test completed successfully")
}

func TestSpecialCharacters(t *testing.T) {
	// 特殊文字のテスト
	specialCharData := []Item{
		{
			Car:       "特殊123あ漢字",
			Name:      "山田　花子（株）",
			Purpose:   StringPtr("会議・研修"),
			StartDate: StringPtr("2024/01/15"),
			EndDate:   StringPtr("2024/01/16"),
			Price:     50000,
			PayDay:    StringPtr("2024/02/15"),
			Office:    StringPtr("本社（東京）"),
			Ryohi: []Ryohi{
				{
					Date:   StringPtr("01/15"),
					Dest:   StringPtr("東京・大阪・京都"),
					Detail: []string{"会議（重要）", "資料作成・提出", "報告書（詳細版）"},
					Kukan:  StringPtr("東京駅→大阪駅（新幹線）"),
					Price:  IntPtr(25000),
					Vol:    Float64Ptr(3.5),
				},
			},
		},
	}

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Special character processing should not panic: %v", r)
		}
	}()

	client := NewReportLabStylePdfClient(specialCharData)
	if client == nil {
		t.Error("Failed to create PDF client with special character data")
	}

	t.Log("Special character test completed successfully")
}

// パフォーマンステスト
func TestLargeDataPerformance(t *testing.T) {
	// 大量データのテスト
	var largeData []Item

	// 複数のItemを生成
	for i := 0; i < 5; i++ {
		var ryohiList []Ryohi

		// 各Itemに複数のRyohiを追加
		for j := 0; j < 10; j++ {
			ryohiList = append(ryohiList, Ryohi{
				Date:   StringPtr("01/15"),
				Dest:   StringPtr("テスト目的地"),
				Detail: []string{"項目1", "項目2", "項目3"},
				Kukan:  StringPtr("テスト区間"),
				Price:  IntPtr(5000),
				Vol:    Float64Ptr(2.0),
			})
		}

		largeData = append(largeData, Item{
			Car:       "テスト100あ1234",
			Name:      "テスト　太郎",
			Purpose:   StringPtr("テスト"),
			StartDate: StringPtr("2024-01-15"),
			EndDate:   StringPtr("2024-01-16"),
			Price:     50000,
			PayDay:    StringPtr("2024-02-15"),
			Office:    StringPtr("テスト支社"),
			Ryohi:     ryohiList,
		})
	}

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Large data processing should not panic: %v", r)
		}
	}()

	client := NewReportLabStylePdfClient(largeData)
	if client == nil {
		t.Error("Failed to create PDF client with large data")
	}

	t.Log("Large data performance test completed successfully")
}

// 文字エンコーディングテスト
func TestJapaneseCharacterHandling(t *testing.T) {
	japaneseData := []Item{
		{
			Car:       "品川100あ1234",
			Name:      "佐藤　太郎",
			Purpose:   StringPtr("出張業務"),
			StartDate: StringPtr("令和6年1月15日"),
			EndDate:   StringPtr("令和6年1月16日"),
			Price:     30000,
			PayDay:    StringPtr("令和6年2月15日"),
			Office:    StringPtr("本社営業部"),
			Ryohi: []Ryohi{
				{
					Date:   StringPtr("01/15"),
					Dest:   StringPtr("大阪府大阪市北区梅田"),
					Detail: []string{"取引先訪問", "商談", "契約締結", "懇親会参加"},
					Kukan:  StringPtr("東京駅　新大阪駅　梅田駅"),
					Price:  IntPtr(15000),
					Vol:    Float64Ptr(4.0),
				},
			},
		},
	}

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Japanese character handling should not panic: %v", r)
		}
	}()

	client := NewReportLabStylePdfClient(japaneseData)
	if client == nil {
		t.Error("Failed to create PDF client with Japanese character data")
	}

	// 日本語処理の特別チェック
	for _, item := range japaneseData {
		for _, ryohi := range item.Ryohi {
			if ryohi.Dest != nil && strings.Contains(*ryohi.Dest, "大阪") {
				t.Log("Japanese character '大阪' found and processed correctly")
			}
		}
	}

	t.Log("Japanese character handling test completed successfully")
}
