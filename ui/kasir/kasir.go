package kasir

import (
	"fmt"
	"rnd/internal/models"
	"rnd/internal/services"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// CartItem merepresentasikan satu item dalam keranjang belanja.
type CartItem struct {
	ProdukID    uint
	NamaProduk  string
	HargaSatuan int
	Jumlah      int
	Subtotal    int
}

// BuildKasirTab membangun tab Kasir utama.
func BuildKasirTab(window fyne.Window) fyne.CanvasObject {
	// State keranjang
	var cart []CartItem
	var produkList []models.Produk

	// --- Komponen kiri: daftar produk ---
	entryCari := widget.NewEntry()
	entryCari.SetPlaceHolder("Cari produk...")

	listProduk := widget.NewList(
		func() int { return len(produkList) },
		func() fyne.CanvasObject {
			return widget.NewLabel("Produk Template")
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			if id < len(produkList) {
				p := produkList[id]
				obj.(*widget.Label).SetText(fmt.Sprintf("%s - Rp %s (Stok: %d)", p.Nama, formatRupiah(p.Harga), p.Stok))
			}
		},
	)

	refreshProduk := func(query string) {
		var err error
		if query == "" {
			produkList, err = services.GetAllProduk()
		} else {
			produkList, err = services.SearchProduk(query)
		}
		if err == nil {
			listProduk.Refresh()
		}
	}

	entryCari.OnChanged = refreshProduk
	refreshProduk("")

	listProduk.OnSelected = func(id widget.ListItemID) {
		if id >= len(produkList) {
			return
		}
		p := produkList[id]
		// Tambah ke keranjang (atau tambah jumlah jika sudah ada)
		found := false
		for i := range cart {
			if cart[i].ProdukID == p.ID {
				cart[i].Jumlah++
				cart[i].Subtotal = cart[i].Jumlah * cart[i].HargaSatuan
				found = true
				break
			}
		}
		if !found {
			cart = append(cart, CartItem{
				ProdukID:    p.ID,
				NamaProduk:  p.Nama,
				HargaSatuan: p.Harga,
				Jumlah:      1,
				Subtotal:    p.Harga,
			})
		}
		// Tidak perlu refresh — kita rebuild cart table
	}

	produkPanel := container.NewBorder(
		container.NewVBox(widget.NewLabel("Cari:"), entryCari, widget.NewSeparator()),
		nil, nil, nil,
		listProduk,
	)

	// --- Komponen kanan: keranjang ---
	cartTable := widget.NewTable(
		func() (int, int) { return len(cart), 4 },
		func() fyne.CanvasObject {
			return widget.NewLabel("Cell Template")
		},
		func(cell widget.TableCellID, obj fyne.CanvasObject) {
			lbl := obj.(*widget.Label)
			if cell.Row >= len(cart) {
				lbl.SetText("")
				return
			}
			item := cart[cell.Row]
			switch cell.Col {
			case 0:
				lbl.SetText(item.NamaProduk)
			case 1:
				lbl.SetText(itoa(item.Jumlah) + "x")
			case 2:
				lbl.SetText(formatRupiah(item.HargaSatuan))
			case 3:
				lbl.SetText(formatRupiah(item.Subtotal))
			}
		},
	)
	cartTable.SetColumnWidth(0, 160)
	cartTable.SetColumnWidth(1, 60)
	cartTable.SetColumnWidth(2, 100)
	cartTable.SetColumnWidth(3, 100)

	totalLabel := widget.NewLabelWithStyle("Total: Rp 0", fyne.TextAlignTrailing, fyne.TextStyle{Bold: true})

	updateTotal := func() {
		total := 0
		for _, item := range cart {
			total += item.Subtotal
		}
		totalLabel.SetText("Total: Rp " + formatRupiah(total))
	}

	btnHapusItem := widget.NewButton("Hapus Item", func() {
		// TODO: pilih item untuk dihapus
	})

	btnCheckout := widget.NewButton("💳 Bayar", func() {
		total := 0
		for _, item := range cart {
			total += item.Subtotal
		}
		if total == 0 {
			return
		}
		tampilkanCheckout(window, cart, total, func() {
			cart = nil
			cartTable.Refresh()
			updateTotal()
			refreshProduk("")
		})
	})

	cartPanel := container.NewBorder(
		container.NewVBox(
			widget.NewLabelWithStyle("🛒 Keranjang", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			widget.NewSeparator(),
		),
		container.NewVBox(
			widget.NewSeparator(),
			totalLabel,
			container.NewHBox(btnHapusItem, btnCheckout),
		),
		nil, nil,
		cartTable,
	)

	// Gabung layout: kiri produk (70%), kanan keranjang (30%)
	split := container.NewHSplit(produkPanel, cartPanel)
	split.SetOffset(0.65)

	return split
}
