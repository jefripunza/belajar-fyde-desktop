package services

import (
	"errors"
	"fmt"
	"rnd/internal/db"
	"rnd/internal/models"
	"time"

	"gorm.io/gorm"
)

// GetAllTransaksi ambil semua transaksi, diurut terbaru dulu.
func GetAllTransaksi() ([]models.Transaksi, error) {
	var transaksi []models.Transaksi
	result := db.DB.Order("created_at DESC").Find(&transaksi)
	return transaksi, result.Error
}

// GetTransaksiByID ambil satu transaksi beserta detail item-nya.
func GetTransaksiByID(id uint) (models.Transaksi, error) {
	var t models.Transaksi
	result := db.DB.Preload("Details").First(&t, id)
	return t, result.Error
}

// GetTransaksiByNoStruk ambil transaksi berdasarkan nomor struk.
func GetTransaksiByNoStruk(noStruk string) (models.Transaksi, error) {
	var t models.Transaksi
	result := db.DB.Preload("Details").Where("no_struk = ?", noStruk).First(&t)
	return t, result.Error
}

// GenerateNoStruk buat nomor struk format: INV-YYYYMMDD-NNNN
// Contoh: INV-20260617-0001
func GenerateNoStruk() (string, error) {
	var count int64
	today := time.Now().Format("20060102")

	// Hitung transaksi hari ini
	db.DB.Model(&models.Transaksi{}).
		Where("no_struk LIKE ?", "INV-"+today+"-%").
		Count(&count)

	// Format nomor urut 4 digit
	no := count + 1
	noStruk := "INV-" + today + "-" + fmt.Sprintf("%04d", no)
	return noStruk, nil
}

// Checkout memproses transaksi: validasi stok, kurangi stok, simpan transaksi.
// Semua dijalankan dalam SATU transaksi database.
func Checkout(items []DetailItem, total, bayar, kembali int) (*models.Transaksi, error) {
	if len(items) == 0 {
		return nil, errors.New("keranjang kosong")
	}
	if bayar < total {
		return nil, errors.New("uang bayar kurang")
	}

	noStruk, err := GenerateNoStruk()
	if err != nil {
		return nil, err
	}

	var transaksi *models.Transaksi

	err = db.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Validasi & kurangi stok setiap item
		for _, item := range items {
			var produk models.Produk
			if err := tx.First(&produk, item.ProdukID).Error; err != nil {
				return fmt.Errorf("produk dengan ID %d tidak ditemukan", item.ProdukID)
			}
			if produk.Stok < item.Jumlah {
				return fmt.Errorf("stok %s tidak mencukupi (tersedia: %d)", produk.Nama, produk.Stok)
			}
			// Kurangi stok
			if err := KurangiStok(tx, item.ProdukID, item.Jumlah); err != nil {
				return err
			}
		}

		// 2. Insert transaksi header
		t := models.Transaksi{
			NoStruk: noStruk,
			Total:   total,
			Bayar:   bayar,
			Kembali: kembali,
		}
		if err := tx.Create(&t).Error; err != nil {
			return err
		}

		// 3. Insert detail transaksi
		for _, item := range items {
			detail := models.DetailTransaksi{
				TransaksiID: t.ID,
				ProdukID:    item.ProdukID,
				NamaProduk:  item.NamaProduk,
				HargaSatuan: item.HargaSatuan,
				Jumlah:      item.Jumlah,
				Subtotal:    item.Subtotal,
			}
			if err := tx.Create(&detail).Error; err != nil {
				return err
			}
		}

		// 4. Preload details untuk return
		transaksi = &t
		return tx.Preload("Details").First(transaksi, t.ID).Error
	})

	return transaksi, err
}

// DetailItem adalah struct sementara untuk menampung item keranjang sebelum checkout.
type DetailItem struct {
	ProdukID    uint
	NamaProduk  string
	HargaSatuan int
	Jumlah      int
	Subtotal    int
}
