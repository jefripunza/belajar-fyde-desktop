package transaksi

import (
	"rnd/internal/models"
	"rnd/internal/printer"
	"rnd/internal/services"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// BuildRiwayatTab membangun tab Riwayat Transaksi.
func BuildRiwayatTab(window fyne.Window) fyne.CanvasObject {
	var transaksiList []models.Transaksi

	tabelRiwayat := widget.NewTable(
		func() (int, int) { return len(transaksiList) + 1, 5 },
		func() fyne.CanvasObject {
			return widget.NewLabel("Cell")
		},
		func(cell widget.TableCellID, obj fyne.CanvasObject) {
			lbl := obj.(*widget.Label)
			// Header row
			if cell.Row == 0 {
				switch cell.Col {
				case 0:
					lbl.SetText("No. Struk")
					lbl.TextStyle = fyne.TextStyle{Bold: true}
				case 1:
					lbl.SetText("Tanggal")
					lbl.TextStyle = fyne.TextStyle{Bold: true}
				case 2:
					lbl.SetText("Total")
					lbl.TextStyle = fyne.TextStyle{Bold: true}
				case 3:
					lbl.SetText("Bayar")
					lbl.TextStyle = fyne.TextStyle{Bold: true}
				case 4:
					lbl.SetText("Kembali")
					lbl.TextStyle = fyne.TextStyle{Bold: true}
				}
				return
			}
			// Data rows
			idx := cell.Row - 1
			if idx >= len(transaksiList) {
				lbl.SetText("")
				return
			}
			t := transaksiList[idx]
			switch cell.Col {
			case 0:
				lbl.SetText(t.NoStruk)
			case 1:
				lbl.SetText(t.CreatedAt.Format("02/01/06 15:04"))
			case 2:
				lbl.SetText(fmtRupiah(t.Total))
			case 3:
				lbl.SetText(fmtRupiah(t.Bayar))
			case 4:
				lbl.SetText(fmtRupiah(t.Kembali))
			}
		},
	)

	tabelRiwayat.SetColumnWidth(0, 180)
	tabelRiwayat.SetColumnWidth(1, 130)
	tabelRiwayat.SetColumnWidth(2, 100)
	tabelRiwayat.SetColumnWidth(3, 100)
	tabelRiwayat.SetColumnWidth(4, 100)

	refreshTabel := func() {
		var err error
		transaksiList, err = services.GetAllTransaksi()
		if err != nil {
			dialog.ShowError(err, window)
			return
		}
		tabelRiwayat.Refresh()
	}

	btnRefresh := widget.NewButton("Refresh", refreshTabel)

	// Klik row untuk lihat detail
	tabelRiwayat.OnSelected = func(cell widget.TableCellID) {
		if cell.Row <= 0 || cell.Row-1 >= len(transaksiList) {
			return
		}
		t := transaksiList[cell.Row-1]
		tampilkanDetail(window, t.ID, refreshTabel)
	}

	refreshTabel()

	return container.NewBorder(
		container.NewHBox(widget.NewLabelWithStyle("Riwayat Transaksi", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), btnRefresh),
		nil, nil, nil,
		tabelRiwayat,
	)
}

func tampilkanDetail(window fyne.Window, transaksiID uint, onClose func()) {
	transaksi, err := services.GetTransaksiByID(transaksiID)
	if err != nil {
		dialog.ShowError(err, window)
		return
	}

	info := container.NewVBox(
		widget.NewLabel("No. Struk: "+transaksi.NoStruk),
		widget.NewLabel("Tanggal: "+transaksi.CreatedAt.Format("02/01/2006 15:04:05")),
		widget.NewLabel("Total: "+fmtRupiah(transaksi.Total)),
		widget.NewLabel("Bayar: "+fmtRupiah(transaksi.Bayar)),
		widget.NewLabel("Kembali: "+fmtRupiah(transaksi.Kembali)),
		widget.NewSeparator(),
	)

	// Tabel detail item
	tabelDetail := widget.NewTable(
		func() (int, int) { return len(transaksi.Details) + 1, 4 },
		func() fyne.CanvasObject {
			return widget.NewLabel("Cell")
		},
		func(cell widget.TableCellID, obj fyne.CanvasObject) {
			lbl := obj.(*widget.Label)
			if cell.Row == 0 {
				switch cell.Col {
				case 0:
					lbl.SetText("Produk")
					lbl.TextStyle = fyne.TextStyle{Bold: true}
				case 1:
					lbl.SetText("Harga")
					lbl.TextStyle = fyne.TextStyle{Bold: true}
				case 2:
					lbl.SetText("Jml")
					lbl.TextStyle = fyne.TextStyle{Bold: true}
				case 3:
					lbl.SetText("Subtotal")
					lbl.TextStyle = fyne.TextStyle{Bold: true}
				}
				return
			}
			idx := cell.Row - 1
			if idx >= len(transaksi.Details) {
				lbl.SetText("")
				return
			}
			d := transaksi.Details[idx]
			switch cell.Col {
			case 0:
				lbl.SetText(d.NamaProduk)
			case 1:
				lbl.SetText(fmtRupiah(d.HargaSatuan))
			case 2:
				lbl.SetText(fmtInt(d.Jumlah) + "x")
			case 3:
				lbl.SetText(fmtRupiah(d.Subtotal))
			}
		},
	)

	tabelDetail.SetColumnWidth(0, 160)
	tabelDetail.SetColumnWidth(1, 100)
	tabelDetail.SetColumnWidth(2, 60)
	tabelDetail.SetColumnWidth(3, 100)

	btnCetakUlang := widget.NewButton("Cetak Ulang PDF", func() {
		filePath, err := printer.GeneratePDF(&transaksi)
		if err != nil {
			dialog.ShowError(err, window)
			return
		}
		printer.BukaFile(filePath)
	})

	allContent := container.NewVBox(info, tabelDetail, btnCetakUlang)

	dialog.ShowCustom("Detail Transaksi", "Tutup", allContent, window)
}
