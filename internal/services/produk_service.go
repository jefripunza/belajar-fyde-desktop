package services

import (
	"errors"
	"rnd/internal/db"
	"rnd/internal/models"

	"gorm.io/gorm"
)

// GetAllProduk ambil semua produk, diurut berdasarkan nama.
func GetAllProduk() ([]models.Produk, error) {
	var produk []models.Produk
	result := db.DB.Order("nama ASC").Find(&produk)
	return produk, result.Error
}

// SearchProduk cari produk berdasarkan kode atau nama (LIKE).
func SearchProduk(query string) ([]models.Produk, error) {
	var produk []models.Produk
	like := "%" + query + "%"
	result := db.DB.Where("kode LIKE ? OR nama LIKE ?", like, like).Order("nama ASC").Find(&produk)
	return produk, result.Error
}

// GetProdukByID ambil satu produk berdasarkan ID.
func GetProdukByID(id uint) (models.Produk, error) {
	var p models.Produk
	result := db.DB.First(&p, id)
	return p, result.Error
}

// CreateProduk tambah produk baru.
func CreateProduk(p *models.Produk) error {
	if p.Kode == "" || p.Nama == "" || p.Harga <= 0 {
		return errors.New("kode, nama, dan harga harus diisi dengan benar")
	}
	return db.DB.Create(p).Error
}

// UpdateProduk perbarui produk yang sudah ada.
func UpdateProduk(p *models.Produk) error {
	if p.ID == 0 {
		return errors.New("ID produk tidak valid")
	}
	// Hanya update kolom yang diizinkan
	return db.DB.Model(p).Select("Kode", "Nama", "Harga", "Stok").Updates(p).Error
}

// DeleteProduk hapus produk berdasarkan ID (soft delete via GORM).
func DeleteProduk(id uint) error {
	return db.DB.Delete(&models.Produk{}, id).Error
}

// KurangiStok kurangi stok produk sebanyak jumlah. Return error jika stok tidak cukup.
// Harus dipanggil dalam transaksi DB.
func KurangiStok(tx *gorm.DB, produkID uint, jumlah int) error {
	result := tx.Model(&models.Produk{}).
		Where("id = ? AND stok >= ?", produkID, jumlah).
		Update("stok", gorm.Expr("stok - ?", jumlah))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("stok tidak mencukupi atau produk tidak ditemukan")
	}
	return nil
}
