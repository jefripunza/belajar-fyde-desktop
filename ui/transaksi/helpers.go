package transaksi

import (
	"strconv"
	"strings"
)

func fmtRupiah(n int) string {
	if n < 0 {
		return "-Rp " + formatRupiah(-n)
	}
	return "Rp " + formatRupiah(n)
}

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

func fmtInt(n int) string {
	return strconv.Itoa(n)
}
