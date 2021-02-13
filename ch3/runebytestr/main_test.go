package main

import "testing"

var str string = "1234567890"

func BenchmarkComma(b *testing.B) {
	for i := 0; i < b.N; i++ {
		comma(str)
	}
}

func BenchmarkCommaByByte(b *testing.B) {
	for i := 0; i < b.N; i++ {
		commaByByte(str)
	}
}

func BenchmarkCommaByRune(b *testing.B) {
	for i := 0; i < b.N; i++ {
		commaByRune(str)
	}
}

func TestCommaFull(t *testing.T) {
	givenWanted := map[string]string{
		"1":                     "1",
		"12":                    "12",
		"123":                   "123",
		"1234":                  "1,234",
		"12345":                 "12,345",
		"123456":                "123,456",
		"1234567":               "1,234,567",
		"12345678":              "12,345,678",
		"123456789":             "123,456,789",
		"1234567890":            "1,234,567,890",
		"1.1":                   "1.1",
		"12.12":                 "12.12",
		"123.123":               "123.123",
		"1234.1234":             "1,234.1234",
		"12345.12345":           "12,345.12345",
		"123456.123456":         "123,456.123456",
		"1234567.1234567":       "1,234,567.1234567",
		"12345678.12345678":     "12,345,678.12345678",
		"123456789.123456789":   "123,456,789.123456789",
		"1234567890.1234567890": "1,234,567,890.1234567890",
	}

	for _, sign := range []string{"", "-", "+"} {
		for given, wanted := range givenWanted {
			given = sign + given
			wanted = sign + wanted
			got := commaFull(given)
			if wanted != got {
				t.Errorf("%24q -> wanted %28q, got %28q", given, wanted, got)
			}
		}
	}
}

var anagramData = []string{
	"",
	"1",
	"12",
	"123",
	"1234",
	"12345",
	"123456",
	"1234567",
	"12345678",
	"123456789",
	"1234567890",
	"えええ",
}

var anagramData1 = []string{
	"",
	"1",
	"21",
	"321",
	"4321",
	"54321",
	"654321",
	"7654321",
	"87654321",
	"987654321",
	"0987654321",
	"えええ",
}
var n = len(anagramData)

func BenchmarkIsAnagram(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for j := 0; j < n; j++ {
			isAnagram(anagramData[j], anagramData1[j])
			isAnagram(anagramData[j], anagramData[j])
		}
	}
}

func BenchmarkIsAnagramByReverse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for j := 0; j < n; j++ {
			isAnagramByReverse(anagramData[j], anagramData1[j])
			isAnagramByReverse(anagramData[j], anagramData[j])
		}
	}
}
