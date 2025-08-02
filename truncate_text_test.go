package main

import (
	"testing"
)

func TestTruncateText(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		maxLength int
		expected  string
	}{
		{
			name:      "短いテキスト - 切り詰め不要",
			input:     "短いテキスト",
			maxLength: 22,
			expected:  "短いテキスト",
		},
		{
			name:      "長いテキスト - 切り詰め必要",
			input:     "これは非常に長いテキストで切り詰めが必要です",
			maxLength: 10,
			expected:  "これは非常に長いテ...",
		},
		{
			name:      "英語テキスト - 切り詰め必要",
			input:     "This is a very long English text that needs truncation",
			maxLength: 15,
			expected:  "This is a very ...",
		},
		{
			name:      "空文字列",
			input:     "",
			maxLength: 10,
			expected:  "",
		},
		{
			name:      "maxLength が 0",
			input:     "テスト",
			maxLength: 0,
			expected:  "...",
		},
		{
			name:      "区間テスト - 実際のケース",
			input:     "東京駅→大阪駅→京都駅→神戸駅",
			maxLength: 22,
			expected:  "東京駅→大阪駅→京都駅→神戸駅",
		},
		{
			name:      "区間テスト - 長すぎるケース",
			input:     "東京駅→新横浜駅→名古屋駅→京都駅→新大阪駅→広島駅→博多駅",
			maxLength: 22,
			expected:  "東京駅→新横浜駅→名古屋駅→京都駅...",
		},
		{
			name:      "行先テスト - 短いケース",
			input:     "大阪",
			maxLength: 8,
			expected:  "大阪",
		},
		{
			name:      "行先テスト - 長いケース",
			input:     "東京都千代田区丸の内",
			maxLength: 8,
			expected:  "東京都千代田区...",
		},
		{
			name:      "摘要テスト - 複数項目",
			input:     "会議、研修、営業活動、資料作成、打ち合わせ",
			maxLength: 18,
			expected:  "会議、研修、営業活動、資料作成...",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := truncateText(tt.input, tt.maxLength)
			if result != tt.expected {
				t.Errorf("truncateText(%q, %d) = %q, expected %q",
					tt.input, tt.maxLength, result, tt.expected)
			}
		})
	}
}

func TestTruncateTextEdgeCases(t *testing.T) {
	// ユニコード文字のテスト
	t.Run("絵文字を含むテキスト", func(t *testing.T) {
		input := "会議🚀資料📝作成✨完了🎉"
		result := truncateText(input, 5)
		expected := "会議🚀資料📝..."
		if result != expected {
			t.Errorf("truncateText(%q, 5) = %q, expected %q", input, result, expected)
		}
	})

	// 境界値テスト
	t.Run("maxLengthと同じ長さ", func(t *testing.T) {
		input := "12345"
		result := truncateText(input, 5)
		expected := "12345"
		if result != expected {
			t.Errorf("truncateText(%q, 5) = %q, expected %q", input, result, expected)
		}
	})

	t.Run("maxLength + 1の長さ", func(t *testing.T) {
		input := "123456"
		result := truncateText(input, 5)
		expected := "12345..."
		if result != expected {
			t.Errorf("truncateText(%q, 5) = %q, expected %q", input, result, expected)
		}
	})
}

// ベンチマークテスト
func BenchmarkTruncateText(b *testing.B) {
	longText := "これは非常に長いテキストで、パフォーマンステストのために使用されます。日本語の文字数カウントは複雑な処理が必要です。"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		truncateText(longText, 22)
	}
}

func BenchmarkTruncateTextShort(b *testing.B) {
	shortText := "短いテキスト"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		truncateText(shortText, 22)
	}
}
