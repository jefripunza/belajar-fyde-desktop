package produk

import (
	"strconv"
	"strings"
)

// formatRupiah format integer rupiah ke tampilan "Rp 10.000"
func formatRupiah(n int) string {
	s := strconv.Itoa(n)
	// Tambah titik setiap 3 digit dari kanan
	var result []string
	for i, c := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			result = append(result, ".")
		}
		result = append(result, string(c))
	}
	return strings.Join(result, "")
}

func itoa(n int) string {
	return strconv.Itoa(n)
}

func parseHarga(s string) int {
	s = strings.ReplaceAll(s, ".", "")
	s = strings.ReplaceAll(s, ",", "")
	s = strings.TrimSpace(s)
	n, _ := strconv.Atoi(s)
	return n
}

func parseInt(s string) int {
	n, _ := strconv.Atoi(strings.TrimSpace(s))
	return n
}

func atoi(s string) int {
	n, _ := strconv.Atoi(strings.TrimSpace(s))
	return n
}

// FmtRupiah exported untuk digunakan dari package lain (struk, kasir)
func FmtRupiah(n int) string {
	if n < 0 {
		return "-Rp " + formatRupiah(-n)
	}
	return "Rp " + formatRupiah(n)
}

func Itoa(n int) string {
	return strconv.Itoa(n)
}

func Atoi(s string) int {
	n, _ := strconv.Atoi(strings.TrimSpace(s))
	return n
}
