package db

import (
	"log"
	"rnd/internal/models"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Init buka koneksi SQLite + jalankan AutoMigrate.
// Dipanggil dari main.go sebelum UI start.
func Init(dbPath string) {
	var err error
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatalf("Gagal koneksi database: %v", err)
	}

	err = DB.AutoMigrate(
		&models.Produk{},
		&models.Transaksi{},
		&models.DetailTransaksi{},
	)
	if err != nil {
		log.Fatalf("Gagal migrasi database: %v", err)
	}

	log.Println("Database terhubung & migrasi selesai")
}
