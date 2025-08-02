package main

import (
	"testing"
)

func TestFormatPrice(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected string
	}{
		{
			name:     "正の数値",
			input:    5000,
			expected: "5,000",
		},
		{
			name:     "0",
			input:    0,
			expected: "0",
		},
		{
			name:     "1桁",
			input:    5,
			expected: "5",
		},
		{
			name:     "4桁",
			input:    1234,
			expected: "1,234",
		},
		{
			name:     "7桁",
			input:    1234567,
			expected: "1,234,567",
		},
		{
			name:     "負の数値",
			input:    -5000,
			expected: "-5,000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatPrice(tt.input)
			if result != tt.expected {
				t.Errorf("FormatPrice(%d) = %s, expected %s", tt.input, result, tt.expected)
			}
		})
	}
}

func TestStringPtr(t *testing.T) {
	s := "test"
	ptr := StringPtr(s)

	if ptr == nil {
		t.Error("StringPtr should not return nil")
	}

	if *ptr != s {
		t.Errorf("StringPtr returned pointer to %s, expected %s", *ptr, s)
	}
}

func TestIntPtr(t *testing.T) {
	i := 42
	ptr := IntPtr(i)

	if ptr == nil {
		t.Error("IntPtr should not return nil")
	}

	if *ptr != i {
		t.Errorf("IntPtr returned pointer to %d, expected %d", *ptr, i)
	}
}

func TestFloat64Ptr(t *testing.T) {
	f := 3.14
	ptr := Float64Ptr(f)

	if ptr == nil {
		t.Error("Float64Ptr should not return nil")
	}

	if *ptr != f {
		t.Errorf("Float64Ptr returned pointer to %f, expected %f", *ptr, f)
	}
}

func TestItemStruct(t *testing.T) {
	item := Item{
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
				Dest:   StringPtr("東京"),
				Detail: []string{"会議", "資料作成"},
				Kukan:  StringPtr("東京駅　大阪駅"),
				Price:  IntPtr(5000),
				Vol:    Float64Ptr(1.5),
			},
		},
	}

	// 基本的なフィールドアクセステスト
	if item.Car != "長崎100か4105" {
		t.Error("Car field not properly set")
	}

	if item.Name != "松本　俊之" {
		t.Error("Name field not properly set")
	}

	if item.Purpose == nil || *item.Purpose != "営業" {
		t.Error("Purpose field not properly set")
	}

	if item.Price != 86900 {
		t.Error("Price field not properly set")
	}

	if len(item.Ryohi) != 1 {
		t.Error("Ryohi array not properly set")
	}

	// Ryohiフィールドのテスト
	ryohi := item.Ryohi[0]
	if ryohi.Date == nil || *ryohi.Date != "01/15" {
		t.Error("Ryohi Date field not properly set")
	}

	if ryohi.Price == nil || *ryohi.Price != 5000 {
		t.Error("Ryohi Price field not properly set")
	}
}

func TestRyohiStruct(t *testing.T) {
	ryohi := Ryohi{
		Date:   StringPtr("01/15"),
		Dest:   StringPtr("東京"),
		Detail: []string{"会議", "資料作成"},
		Kukan:  StringPtr("東京駅　大阪駅"),
		Price:  IntPtr(5000),
		Vol:    Float64Ptr(1.5),
	}

	// Ryohiの各フィールドテスト
	if ryohi.Date == nil || *ryohi.Date != "01/15" {
		t.Error("Ryohi Date field not properly accessible")
	}

	if ryohi.Price == nil || *ryohi.Price != 5000 {
		t.Error("Ryohi Price field not properly accessible")
	}

	if ryohi.Dest == nil || *ryohi.Dest != "東京" {
		t.Error("Ryohi Dest field not properly accessible")
	}

	if len(ryohi.Detail) != 2 || ryohi.Detail[0] != "会議" {
		t.Error("Ryohi Detail field not properly accessible")
	}
}

// ベンチマークテスト
func BenchmarkFormatPrice(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FormatPrice(1234567)
	}
}

func BenchmarkWrapDetail(b *testing.B) {
	details := []string{"会議", "研修", "営業活動", "資料作成", "打ち合わせ", "報告書作成"}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		wrapDetail(details, 10)
	}
}

func BenchmarkWrapKukan(b *testing.B) {
	kukan := "博多駅　小倉駅　新山口駅　広島駅　岡山駅　新大阪駅　京都駅　東京駅"
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		wrapKukan(kukan, 20)
	}
}

func BenchmarkPrepareRyohiForPrint(b *testing.B) {
	ryohi := Ryohi{
		Date:   StringPtr("01/15"),
		Dest:   StringPtr("東京都千代田区"),
		Detail: []string{"重要会議", "新規プロジェクト打ち合わせ", "資料作成", "報告書提出", "クライアント面談"},
		Kukan:  StringPtr("博多駅　小倉駅　新山口駅　広島駅　岡山駅　新大阪駅　京都駅　東京駅"),
		Price:  IntPtr(15000),
		Vol:    Float64Ptr(3.5),
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		prepareRyohiForPrint(ryohi, 7, 20)
	}
}
