package main

import (
	"reflect"
	"testing"
)

func TestWrapDetail(t *testing.T) {
	tests := []struct {
		name        string
		details     []string
		maxLen      int
		expectLines []string
		expectRows  int
	}{
		{
			name:        "短い摘要 - 1行",
			details:     []string{"会議", "資料作成"},
			maxLen:      10,
			expectLines: []string{"会議、資料作成"},
			expectRows:  1,
		},
		{
			name:        "長い摘要 - 複数行",
			details:     []string{"会議", "研修", "営業活動", "資料作成", "打ち合わせ"},
			maxLen:      7,
			expectLines: []string{"会議、研修", "営業活動", "資料作成", "打ち合わせ"},
			expectRows:  4,
		},
		{
			name:        "最大長ぴったり",
			details:     []string{"12345", "6789"},
			maxLen:      10,
			expectLines: []string{"12345、6789"},
			expectRows:  1,
		},
		{
			name:        "空の摘要",
			details:     []string{},
			maxLen:      10,
			expectLines: []string{},
			expectRows:  0,
		},
		{
			name:        "単一項目が最大長",
			details:     []string{"1234567890"},
			maxLen:      10,
			expectLines: []string{"1234567890"},
			expectRows:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := wrapDetail(tt.details, tt.maxLen)
			if !reflect.DeepEqual(result.Lines, tt.expectLines) {
				t.Errorf("wrapDetail() lines = %v, expected %v", result.Lines, tt.expectLines)
			}
			if result.RowCount != tt.expectRows {
				t.Errorf("wrapDetail() rowCount = %d, expected %d", result.RowCount, tt.expectRows)
			}
		})
	}
}

func TestWrapKukan(t *testing.T) {
	tests := []struct {
		name        string
		kukan       string
		maxLen      int
		expectLines []string
		expectRows  int
	}{
		{
			name:        "短い区間",
			kukan:       "東京　大阪",
			maxLen:      20,
			expectLines: []string{"東京　大阪"},
			expectRows:  1,
		},
		{
			name:        "長い区間 - 折り返し",
			kukan:       "東京駅　新横浜駅　名古屋駅　京都駅　新大阪駅",
			maxLen:      15,
			expectLines: []string{"東京駅　新横浜駅　名古屋駅", "京都駅　新大阪駅"},
			expectRows:  2,
		},
		{
			name:        "特殊文字列置換",
			kukan:       "東京_九州外空車適用　大阪",
			maxLen:      20,
			expectLines: []string{"東京　九州外空車適用　大阪"},
			expectRows:  1,
		},
		{
			name:        "空の区間",
			kukan:       "",
			maxLen:      10,
			expectLines: []string{""},
			expectRows:  1,
		},
		{
			name:        "｜区切り",
			kukan:       "東京｜大阪｜京都",
			maxLen:      10,
			expectLines: []string{"東京　大阪　京都"},
			expectRows:  1,
		},
		{
			name:        "最大長超過",
			kukan:       "非常に長い区間名称で最大長を超過するテスト",
			maxLen:      5,
			expectLines: []string{"exceed*"},
			expectRows:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := wrapKukan(tt.kukan, tt.maxLen)
			if !reflect.DeepEqual(result.Lines, tt.expectLines) {
				t.Errorf("wrapKukan() lines = %v, expected %v", result.Lines, tt.expectLines)
			}
			if result.RowCount != tt.expectRows {
				t.Errorf("wrapKukan() rowCount = %d, expected %d", result.RowCount, tt.expectRows)
			}
		})
	}
}

func TestAlignRows(t *testing.T) {
	tests := []struct {
		name        string
		date        *string
		dest        *string
		price       *int
		vol         *float64
		maxRows     int
		expectDate  []string
		expectDest  []string
		expectPrice []string
		expectVol   []string
	}{
		{
			name:        "基本的な行揃え",
			date:        StringPtr("2024-01-15"),
			dest:        StringPtr("東京"),
			price:       IntPtr(5000),
			vol:         Float64Ptr(1.5),
			maxRows:     3,
			expectDate:  []string{"01/15", "", ""},
			expectDest:  []string{"東京", "", ""},
			expectPrice: []string{"5,000", "", ""},
			expectVol:   []string{"1.5", "", ""},
		},
		{
			name:        "nilデータの処理",
			date:        nil,
			dest:        nil,
			price:       nil,
			vol:         nil,
			maxRows:     2,
			expectDate:  []string{"", ""},
			expectDest:  []string{"", ""},
			expectPrice: []string{"", ""},
			expectVol:   []string{"", ""},
		},
		{
			name:        "単一行",
			date:        StringPtr("2024-02-20"),
			dest:        StringPtr("大阪"),
			price:       IntPtr(10000),
			vol:         Float64Ptr(2.0),
			maxRows:     1,
			expectDate:  []string{"02/20"},
			expectDest:  []string{"大阪"},
			expectPrice: []string{"10,000"},
			expectVol:   []string{"2.0"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dateArr, destArr, priceArr, volArr := alignRows(tt.date, tt.dest, tt.price, tt.vol, tt.maxRows)

			if !reflect.DeepEqual(dateArr, tt.expectDate) {
				t.Errorf("alignRows() dateArr = %v, expected %v", dateArr, tt.expectDate)
			}
			if !reflect.DeepEqual(destArr, tt.expectDest) {
				t.Errorf("alignRows() destArr = %v, expected %v", destArr, tt.expectDest)
			}
			if !reflect.DeepEqual(priceArr, tt.expectPrice) {
				t.Errorf("alignRows() priceArr = %v, expected %v", priceArr, tt.expectPrice)
			}
			if !reflect.DeepEqual(volArr, tt.expectVol) {
				t.Errorf("alignRows() volArr = %v, expected %v", volArr, tt.expectVol)
			}
		})
	}
}

func TestExtendToMaxRows(t *testing.T) {
	tests := []struct {
		name     string
		lines    []string
		maxRows  int
		expected []string
	}{
		{
			name:     "拡張が必要",
			lines:    []string{"line1", "line2"},
			maxRows:  4,
			expected: []string{"line1", "line2", "", ""},
		},
		{
			name:     "拡張不要",
			lines:    []string{"line1", "line2", "line3"},
			maxRows:  3,
			expected: []string{"line1", "line2", "line3"},
		},
		{
			name:     "短縮（実際は短縮されない）",
			lines:    []string{"line1", "line2", "line3", "line4"},
			maxRows:  2,
			expected: []string{"line1", "line2"},
		},
		{
			name:     "空配列",
			lines:    []string{},
			maxRows:  3,
			expected: []string{"", "", ""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extendToMaxRows(tt.lines, tt.maxRows)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("extendToMaxRows() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestPrepareRyohiForPrint(t *testing.T) {
	tests := []struct {
		name            string
		ryohi           Ryohi
		maxDetailLen    int
		maxKukanLen     int
		expectedMaxRows int
	}{
		{
			name: "基本的な旅費データ",
			ryohi: Ryohi{
				Date:   StringPtr("01/15"),
				Dest:   StringPtr("東京"),
				Detail: []string{"会議", "資料作成"},
				Kukan:  StringPtr("東京駅　新横浜駅　名古屋駅"),
				Price:  IntPtr(5000),
				Vol:    Float64Ptr(1.5),
			},
			maxDetailLen:    10,
			maxKukanLen:     15,
			expectedMaxRows: 2, // 区間が2行になるため
		},
		{
			name: "摘要が長い場合",
			ryohi: Ryohi{
				Date:   StringPtr("01/20"),
				Dest:   StringPtr("大阪"),
				Detail: []string{"会議", "研修", "営業活動", "資料作成", "打ち合わせ", "報告書作成"},
				Kukan:  StringPtr("東京　大阪"),
				Price:  IntPtr(8000),
				Vol:    Float64Ptr(2.0),
			},
			maxDetailLen:    7,
			maxKukanLen:     20,
			expectedMaxRows: 4, // 摘要が4行になるため
		},
		{
			name: "空データ",
			ryohi: Ryohi{
				Date:   nil,
				Dest:   nil,
				Detail: []string{},
				Kukan:  nil,
				Price:  nil,
				Vol:    nil,
			},
			maxDetailLen:    10,
			maxKukanLen:     15,
			expectedMaxRows: 1, // 最低1行
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := prepareRyohiForPrint(tt.ryohi, tt.maxDetailLen, tt.maxKukanLen)

			if result.MaxRows != tt.expectedMaxRows {
				t.Errorf("prepareRyohiForPrint() MaxRows = %d, expected %d", result.MaxRows, tt.expectedMaxRows)
			}

			// すべての配列が同じ長さであることを確認
			if len(result.DateLines) != tt.expectedMaxRows {
				t.Errorf("DateLines length = %d, expected %d", len(result.DateLines), tt.expectedMaxRows)
			}
			if len(result.DestLines) != tt.expectedMaxRows {
				t.Errorf("DestLines length = %d, expected %d", len(result.DestLines), tt.expectedMaxRows)
			}
			if len(result.DetailLines) != tt.expectedMaxRows {
				t.Errorf("DetailLines length = %d, expected %d", len(result.DetailLines), tt.expectedMaxRows)
			}
			if len(result.KukanLines) != tt.expectedMaxRows {
				t.Errorf("KukanLines length = %d, expected %d", len(result.KukanLines), tt.expectedMaxRows)
			}
			if len(result.PriceLines) != tt.expectedMaxRows {
				t.Errorf("PriceLines length = %d, expected %d", len(result.PriceLines), tt.expectedMaxRows)
			}
			if len(result.VolLines) != tt.expectedMaxRows {
				t.Errorf("VolLines length = %d, expected %d", len(result.VolLines), tt.expectedMaxRows)
			}
		})
	}
}

// 実際のデータを使った統合テスト
func TestIntegrationRealData(t *testing.T) {
	ryohi := Ryohi{
		Date:   StringPtr("01/15"),
		Dest:   StringPtr("東京都千代田区"),
		Detail: []string{"重要会議", "新規プロジェクト打ち合わせ", "資料作成", "報告書提出", "クライアント面談"},
		Kukan:  StringPtr("博多駅　小倉駅　新山口駅　広島駅　岡山駅　新大阪駅　京都駅　東京駅"),
		Price:  IntPtr(15000),
		Vol:    Float64Ptr(3.5),
	}

	result := prepareRyohiForPrint(ryohi, 7, 20)

	t.Logf("MaxRows: %d", result.MaxRows)
	t.Logf("DateLines: %v", result.DateLines)
	t.Logf("DestLines: %v", result.DestLines)
	t.Logf("DetailLines: %v", result.DetailLines)
	t.Logf("KukanLines: %v", result.KukanLines)
	t.Logf("PriceLines: %v", result.PriceLines)
	t.Logf("VolLines: %v", result.VolLines)

	// 基本的な検証
	if result.MaxRows <= 0 {
		t.Error("MaxRows should be positive")
	}

	// 最初の行にデータが入っていることを確認
	if result.DateLines[0] != "01/15" {
		t.Errorf("First date line should be '01/15', got '%s'", result.DateLines[0])
	}
	if result.PriceLines[0] != "15,000" {
		t.Errorf("First price line should be '15,000', got '%s'", result.PriceLines[0])
	}
}
