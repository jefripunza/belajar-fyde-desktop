package ui

import (
	"rnd/ui/kasir"
	"rnd/ui/produk"
	"rnd/ui/transaksi"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// Setup menyusun semua tab dan menampilkan di window utama.
func Setup(window fyne.Window) {
	// Buat konten setiap tab
	tabKasir := kasir.BuildKasirTab(window)
	tabProduk := produk.BuildProdukTab(window)
	tabRiwayat := transaksi.BuildRiwayatTab(window)

	// AppTabs sebagai navigasi utama
	tabs := container.NewAppTabs(
		container.NewTabItem("Kasir", tabKasir),
		container.NewTabItem("Produk", tabProduk),
		container.NewTabItem("Riwayat", tabRiwayat),
	)

	// Label judul di atas tabs
	judul := widget.NewLabelWithStyle("💳 Kasir POS — Toko", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	mainContent := container.NewBorder(
		container.NewVBox(
			judul,
			widget.NewSeparator(),
		),
		nil, nil, nil,
		tabs,
	)

	window.SetContent(mainContent)
}
