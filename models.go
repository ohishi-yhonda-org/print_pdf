package main

import (
	"time"
)

// Ryohi represents the expense data structure
type Ryohi struct {
	Date           *string   `json:"date"`
	DateAr         []string  `json:"dateAr"`
	Dest           *string   `json:"dest"`
	DestAr         []string  `json:"destAr"`
	Detail         []string  `json:"detail"`
	Kukan          *string   `json:"kukan"`
	KukanSprit     []string  `json:"kukanSprit"`
	Price          *int      `json:"price"`
	PriceAr        []int     `json:"priceAr"`
	Vol            *float64  `json:"vol"`
	VolAr          []float64 `json:"volAr"`
	PrintDetail    []string  `json:"printdetail"`
	PrintDetailRow int       `json:"printdetailRow"`
	PrintKukan     []string  `json:"printKukan"`
	PrintKukanRow  int       `json:"printKukanRow"`
	MaxRow         int       `json:"maxRow"`
	PageCount      int       `json:"pageCount"`
}

// Item represents the main item data structure
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

// PrintRequest represents the print request data structure
type PrintRequest struct {
	Items       []Item  `json:"items"`
	Print       bool    `json:"print,omitempty"`       // 印刷するかどうか
	PrinterName *string `json:"printerName,omitempty"` // 指定プリンター名（省略時はデフォルト）
}

// Helper functions
func StringPtr(s string) *string {
	return &s
}

func IntPtr(i int) *int {
	return &i
}

func Float64Ptr(f float64) *float64 {
	return &f
}

// FormatPrice formats price with comma separator
func FormatPrice(price int) string {
	if price == 0 {
		return ""
	}

	str := ""
	num := price
	if num < 0 {
		str = "-"
		num = -num
	}

	digits := []rune{}
	for num > 0 {
		digits = append([]rune{'0' + rune(num%10)}, digits...)
		num /= 10
	}

	if len(digits) == 0 {
		return "0"
	}

	result := ""
	for i, digit := range digits {
		if i > 0 && (len(digits)-i)%3 == 0 {
			result += ","
		}
		result += string(digit)
	}

	return str + result
}

// ParseDate parses date string and returns formatted date
func ParseDate(dateStr string) (string, error) {
	if dateStr == "" {
		return "", nil
	}

	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return "", err
	}

	return t.Format("01　02"), nil
}

// ParsePayDay parses pay day and returns formatted date
func ParsePayDay(dateStr string) (string, error) {
	if dateStr == "" {
		return "", nil
	}

	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return "", err
	}

	return t.Format("2006  01　02"), nil
}
