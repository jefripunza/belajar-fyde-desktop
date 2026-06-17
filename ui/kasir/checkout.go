package kasir

import (
	"rnd/internal/models"
	"rnd/internal/printer"
	"rnd/internal/services"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func tampilkanCheckout(window fyne.Window, cart []CartItem, total int, onSuccess func()) {
	// Tampilkan dialog checkout
	labelTotal := widget.NewLabelWithStyle("Total: Rp "+formatRupiah(total), fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	entryBayar := widget.NewEntry()
	entryBayar.SetPlaceHolder("Jumlah uang diterima...")

	labelKembali := widget.NewLabel("Kembali: Rp 0")

	entryBayar.OnChanged = func(s string) {
		bayar := atoi(s)
		kembali := bayar - total
		if kembali >= 0 {
			labelKembali.SetText("Kembali: Rp " + formatRupiah(kembali))
		} else {
			labelKembali.SetText("Uang kurang: Rp " + formatRupiah(-kembali))
		}
	}

	content := container.NewVBox(
		labelTotal,
		widget.NewSeparator(),
		widget.NewLabel("Uang Diterima:"),
		entryBayar,
		labelKembali,
	)

	dialog.ShowCustomConfirm("Checkout", "💳 Bayar", "Batal", content, func(bayarBtn bool) {
		if !bayarBtn {
			return
		}

		bayar := atoi(entryBayar.Text)
		kembali := bayar - total

		if bayar < total {
			dialog.ShowError(errUangKurang{short: bayar - total}, window)
			return
		}

		// Konversi CartItem ke services.DetailItem
		items := make([]services.DetailItem, len(cart))
		for i, ci := range cart {
			items[i] = services.DetailItem{
				ProdukID:    ci.ProdukID,
				NamaProduk:  ci.NamaProduk,
				HargaSatuan: ci.HargaSatuan,
				Jumlah:      ci.Jumlah,
				Subtotal:    ci.Subtotal,
			}
		}

		transaksi, err := services.Checkout(items, total, bayar, kembali)
		if err != nil {
			dialog.ShowError(err, window)
			return
		}

		// Tampilkan preview struk
		tampilkanPreviewStruk(window, transaksi)

		onSuccess()
	}, window)
}

type errUangKurang struct {
	short int
}

func (e errUangKurang) Error() string {
	return "Uang bayar kurang Rp " + formatRupiah(-e.short)
}

func tampilkanPreviewStruk(window fyne.Window, transaksi *models.Transaksi) {
	// Build konten preview menggunakan widget Fyne
	header := container.NewVBox(
		widget.NewLabelWithStyle("TOKO POS", fyne.TextAlignCenter, fyne.TextStyle{Bold: true, Monospace: true}),
		widget.NewLabelWithStyle("Jl. Contoh No. 123", fyne.TextAlignCenter, fyne.TextStyle{Monospace: true}),
		widget.NewSeparator(),
	)

	info := container.NewVBox(
		widget.NewLabelWithStyle("No: "+transaksi.NoStruk, fyne.TextAlignLeading, fyne.TextStyle{Monospace: true}),
		widget.NewLabelWithStyle("Tgl: "+transaksi.CreatedAt.Format("02/01/2006 15:04"), fyne.TextAlignLeading, fyne.TextStyle{Monospace: true}),
		widget.NewSeparator(),
	)

	var items []fyne.CanvasObject
	for _, d := range transaksi.Details {
		line := widget.NewLabelWithStyle(
			padRight(d.NamaProduk, 16, ' ')+padLeft(itoa(d.Jumlah), 2, ' ')+" x "+padLeft(formatRupiah(d.HargaSatuan), 10, ' '),
			fyne.TextAlignLeading,
			fyne.TextStyle{Monospace: true},
		)
		sub := widget.NewLabelWithStyle(
			padLeft(formatRupiah(d.Subtotal), 28, ' '),
			fyne.TextAlignTrailing,
			fyne.TextStyle{Monospace: true},
		)
		items = append(items, line, sub)
	}

	itemsBox := container.NewVBox(items...)

	footer := container.NewVBox(
		widget.NewSeparator(),
		widget.NewLabelWithStyle("Total:   "+padLeft(formatRupiah(transaksi.Total), 20, ' '), fyne.TextAlignLeading, fyne.TextStyle{Bold: true, Monospace: true}),
		widget.NewLabelWithStyle("Bayar:   "+padLeft(formatRupiah(transaksi.Bayar), 20, ' '), fyne.TextAlignLeading, fyne.TextStyle{Monospace: true}),
		widget.NewLabelWithStyle("Kembali: "+padLeft(formatRupiah(transaksi.Kembali), 20, ' '), fyne.TextAlignLeading, fyne.TextStyle{Monospace: true}),
		widget.NewSeparator(),
		widget.NewLabelWithStyle("Terima Kasih", fyne.TextAlignCenter, fyne.TextStyle{Monospace: true}),
	)

	// Tombol aksi
	btnCetak := widget.NewButton("🖨️ Cetak PDF", func() {
		filePath, err := printer.GeneratePDF(transaksi)
		if err != nil {
			dialog.ShowError(err, window)
			return
		}
		printer.BukaFile(filePath)
	})

	btnTutup := widget.NewButton("Tutup", func() {
		// Tidak perlu aksi — dialog akan tertutup otomatis
	})

	allContent := container.NewVBox(
		header,
		info,
		itemsBox,
		footer,
		container.NewHBox(btnCetak, btnTutup),
	)

	dialog.ShowCustom("Struk Pembayaran", "Tutup", allContent, window)
}

// Helper pad string
func padRight(s string, width int, pad rune) string {
	for len(s) < width {
		s += string(pad)
	}
	return s
}

func padLeft(s string, width int, pad rune) string {
	for len(s) < width {
		s = string(pad) + s
	}
	return s
}
