package kasir

import (
	"strconv"
	"strings"
)

// Helpers untuk formatting rupiah yang digunakan di package kasir.
// Mirror dari ui/produk/helpers.go untuk menghindari circular import.

func formatRupiah(n int) string {
	s := strconv.Itoa(n)
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

func atoi(s string) int {
	n, _ := strconv.Atoi(strings.TrimSpace(s))
	return n
}
