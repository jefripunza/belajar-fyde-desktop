package main

import (
	"os"
	"path/filepath"
	"rnd/internal/db"
	"rnd/ui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	// Dapatkan path executable untuk menyimpan database
	execPath, _ := os.Executable()
	baseDir := filepath.Dir(execPath)
	dbPath := filepath.Join(baseDir, "pos.db")

	// Inisialisasi database
	db.Init(dbPath)

	// Buat aplikasi Fyne
	posApp := app.NewWithID("com.pos.kasir")
	posApp.SetIcon(loadIcon())

	window := posApp.NewWindow("Kasir POS — Toko")

	// Setup UI
	ui.Setup(window)

	window.Resize(fyne.NewSize(1024, 768))
	window.SetMaster()
	window.ShowAndRun()
}

func loadIcon() fyne.Resource {
	// Fallback: tidak ada icon khusus
	return nil
}
