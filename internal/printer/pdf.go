package printer

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"rnd/internal/models"
	"runtime"
	"strconv"
	"strings"

	"github.com/jung-kurt/gofpdf"
)

// GeneratePDF buat file PDF struk dari transaksi.
// Return path file PDF yang dibuat.
func GeneratePDF(t *models.Transaksi) (string, error) {
	// Pastikan folder struk/ ada
	strukDir := "struk"
	if err := os.MkdirAll(strukDir, 0755); err != nil {
		return "", fmt.Errorf("gagal membuat folder struk: %w", err)
	}

	filePath := filepath.Join(strukDir, t.NoStruk+".pdf")

	// Buat PDF ukuran thermal 80mm (lebar 80mm, tinggi otomatis dari konten)
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(10, 10, 10)
	pdf.AddPage()

	// --- Header ---
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(0, 8, "TOKO POS", "", 1, "C", false, 0, "")

	pdf.SetFont("Arial", "", 9)
	pdf.CellFormat(0, 5, "Jl. Contoh No. 123", "", 1, "C", false, 0, "")
	pdf.CellFormat(0, 5, "Telp: 0812-3456-7890", "", 1, "C", false, 0, "")

	// Garis pemisah
	pdf.Line(10, pdf.GetY(), 200, pdf.GetY())
	pdf.Ln(2)

	// --- Info Struk ---
	pdf.SetFont("Arial", "", 9)
	pdf.CellFormat(0, 5, fmt.Sprintf("No. Struk : %s", t.NoStruk), "", 1, "L", false, 0, "")
	pdf.CellFormat(0, 5, fmt.Sprintf("Tanggal  : %s", t.CreatedAt.Format("02/01/2006 15:04:05")), "", 1, "L", false, 0, "")

	// Garis pemisah
	pdf.Line(10, pdf.GetY(), 200, pdf.GetY())
	pdf.Ln(2)

	// --- Header Tabel Item ---
	pdf.SetFont("Arial", "B", 9)
	pdf.CellFormat(75, 5, "Nama Barang", "", 0, "L", false, 0, "")
	pdf.CellFormat(20, 5, "Jml", "", 0, "C", false, 0, "")
	pdf.CellFormat(35, 5, "Harga", "", 0, "R", false, 0, "")
	pdf.CellFormat(40, 5, "Subtotal", "", 1, "R", false, 0, "")
	pdf.Line(10, pdf.GetY(), 200, pdf.GetY())
	pdf.Ln(1)

	// --- Item Barang ---
	pdf.SetFont("Arial", "", 9)
	for _, d := range t.Details {
		pdf.CellFormat(75, 5, d.NamaProduk, "", 0, "L", false, 0, "")
		pdf.CellFormat(20, 5, fmt.Sprintf("%dx", d.Jumlah), "", 0, "C", false, 0, "")
		pdf.CellFormat(35, 5, formatRupiah(d.HargaSatuan), "", 0, "R", false, 0, "")
		pdf.CellFormat(40, 5, formatRupiah(d.Subtotal), "", 1, "R", false, 0, "")
	}

	// Garis pemisah
	pdf.Line(10, pdf.GetY(), 200, pdf.GetY())
	pdf.Ln(2)

	// --- Footer: Total, Bayar, Kembali ---
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(110, 6, "TOTAL", "", 0, "R", false, 0, "")
	pdf.CellFormat(60, 6, formatRupiah(t.Total), "", 1, "R", false, 0, "")

	pdf.SetFont("Arial", "", 9)
	pdf.CellFormat(110, 5, "Bayar", "", 0, "R", false, 0, "")
	pdf.CellFormat(60, 5, formatRupiah(t.Bayar), "", 1, "R", false, 0, "")

	pdf.CellFormat(110, 5, "Kembali", "", 0, "R", false, 0, "")
	pdf.CellFormat(60, 5, formatRupiah(t.Kembali), "", 1, "R", false, 0, "")

	pdf.Ln(4)

	// --- Penutup ---
	pdf.Line(10, pdf.GetY(), 200, pdf.GetY())
	pdf.Ln(3)
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(0, 6, "TERIMA KASIH", "", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "", 8)
	pdf.CellFormat(0, 5, "Barang yang sudah dibeli tidak dapat dikembalikan", "", 1, "C", false, 0, "")

	err := pdf.OutputFileAndClose(filePath)
	if err != nil {
		return "", fmt.Errorf("gagal menyimpan PDF: %w", err)
	}

	return filePath, nil
}

// BukaFile buka file PDF menggunakan aplikasi default OS.
func BukaFile(filePath string) {
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return
	}

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", absPath)
	case "linux":
		cmd = exec.Command("xdg-open", absPath)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", absPath)
	default:
		return
	}
	cmd.Start()
}

// formatRupiah format integer menjadi string dengan titik setiap 3 digit.
func formatRupiah(n int) string {
	s := strconv.Itoa(n)
	var result []string
	for i, c := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			result = append(result, ".")
		}
		result = append(result, string(c))
	}
	return "Rp " + strings.Join(result, "")
}
