package produk

import (
	"rnd/internal/models"
	"rnd/internal/services"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// BuildProdukTab membangun tab Manajemen Produk.
func BuildProdukTab(window fyne.Window) fyne.CanvasObject {
	listProduk := widget.NewList(
		func() int {
			// akan diisi ulang
			return 0
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Template Produk")
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {},
	)

	entryCari := widget.NewEntry()
	entryCari.SetPlaceHolder("Cari produk...")

	var produkList []models.Produk

	refreshList := func() {
		var err error
		produkList, err = services.GetAllProduk()
		if err != nil {
			dialog.ShowError(err, window)
			return
		}
		listProduk.Refresh()
	}

	search := func(query string) {
		var err error
		if query == "" {
			produkList, err = services.GetAllProduk()
		} else {
			produkList, err = services.SearchProduk(query)
		}
		if err != nil {
			dialog.ShowError(err, window)
			return
		}
		listProduk.Refresh()
	}

	entryCari.OnChanged = search

	// Override list render
	listProduk.Length = func() int {
		return len(produkList)
	}

	listProduk.CreateItem = func() fyne.CanvasObject {
		lblNama := widget.NewLabel("Nama Produk")
		lblHarga := widget.NewLabel("Rp 0")
		lblStok := widget.NewLabel("Stok: 0")
		return container.NewHBox(lblNama, widget.NewSeparator(), lblHarga, widget.NewSeparator(), lblStok)
	}

	listProduk.UpdateItem = func(id widget.ListItemID, obj fyne.CanvasObject) {
		if id >= len(produkList) {
			return
		}
		p := produkList[id]
		box := obj.(*fyne.Container)
		lblNama := box.Objects[0].(*widget.Label)
		lblHarga := box.Objects[2].(*widget.Label)
		lblStok := box.Objects[4].(*widget.Label)

		lblNama.SetText(p.Nama + " (" + p.Kode + ")")
		lblHarga.SetText("Rp " + formatRupiah(p.Harga))
		lblStok.SetText("Stok: " + itoa(p.Stok))
	}

	listProduk.OnSelected = func(id widget.ListItemID) {
		if id >= len(produkList) {
			return
		}
		p := produkList[id]
		tampilkanFormEdit(window, p, refreshList)
	}

	btnTambah := widget.NewButtonWithIcon("Tambah Produk", theme.ContentAddIcon(), func() {
		tampilkanFormTambah(window, refreshList)
	})

	btnRefresh := widget.NewButtonWithIcon("Refresh", theme.ViewRefreshIcon(), refreshList)

	toolbar := container.NewHBox(btnTambah, btnRefresh)
	searchBar := container.NewHBox(widget.NewLabel("Cari:"), entryCari)

	top := container.NewVBox(toolbar, searchBar, widget.NewSeparator())

	return container.NewBorder(top, nil, nil, nil, listProduk)
}

func tampilkanFormTambah(window fyne.Window, onSuccess func()) {
	entryKode := widget.NewEntry()
	entryKode.SetPlaceHolder("Kode/Barcode")
	entryNama := widget.NewEntry()
	entryNama.SetPlaceHolder("Nama Produk")
	entryHarga := widget.NewEntry()
	entryHarga.SetPlaceHolder("Harga (Rp)")
	entryStok := widget.NewEntry()
	entryStok.SetPlaceHolder("Stok Awal")

	form := container.NewVBox(
		widget.NewLabel("Tambah Produk Baru"),
		widget.NewForm(
			widget.NewFormItem("Kode", entryKode),
			widget.NewFormItem("Nama", entryNama),
			widget.NewFormItem("Harga", entryHarga),
			widget.NewFormItem("Stok", entryStok),
		),
	)

	dialog.ShowCustomConfirm("Tambah Produk", "Simpan", "Batal", form, func(simpan bool) {
		if !simpan {
			return
		}
		p := models.Produk{
			Kode:  entryKode.Text,
			Nama:  entryNama.Text,
			Harga: parseHarga(entryHarga.Text),
			Stok:  parseInt(entryStok.Text),
		}
		if err := services.CreateProduk(&p); err != nil {
			dialog.ShowError(err, window)
			return
		}
		onSuccess()
	}, window)
}

func tampilkanFormEdit(window fyne.Window, p models.Produk, onSuccess func()) {
	entryKode := widget.NewEntry()
	entryKode.SetText(p.Kode)
	entryNama := widget.NewEntry()
	entryNama.SetText(p.Nama)
	entryHarga := widget.NewEntry()
	entryHarga.SetText(itoa(p.Harga))
	entryStok := widget.NewEntry()
	entryStok.SetText(itoa(p.Stok))

	form := container.NewVBox(
		widget.NewLabel("Edit Produk"),
		widget.NewForm(
			widget.NewFormItem("Kode", entryKode),
			widget.NewFormItem("Nama", entryNama),
			widget.NewFormItem("Harga", entryHarga),
			widget.NewFormItem("Stok", entryStok),
		),
	)

	dialog.ShowCustomConfirm("Edit Produk", "Simpan", "Hapus", form, func(simpan bool) {
		if simpan {
			p.Kode = entryKode.Text
			p.Nama = entryNama.Text
			p.Harga = parseHarga(entryHarga.Text)
			p.Stok = parseInt(entryStok.Text)
			if err := services.UpdateProduk(&p); err != nil {
				dialog.ShowError(err, window)
				return
			}
			onSuccess()
			return
		}
		// Tombol "Hapus" ditekan
		dialog.ShowConfirm("Hapus Produk", "Yakin hapus '"+p.Nama+"'?", func(hapus bool) {
			if hapus {
				if err := services.DeleteProduk(p.ID); err != nil {
					dialog.ShowError(err, window)
					return
				}
				onSuccess()
			}
		}, window)
	}, window)
}
