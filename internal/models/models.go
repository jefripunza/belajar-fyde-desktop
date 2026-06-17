package models

import (
	"time"

	"gorm.io/gorm"
)

// Produk merepresentasikan barang yang dijual.
type Produk struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Kode      string         `gorm:"uniqueIndex;not null" json:"kode"`
	Nama      string         `gorm:"not null" json:"nama"`
	Harga     int            `gorm:"not null" json:"harga"` // Harga jual dalam rupiah (integer)
	Stok      int            `gorm:"not null;default:0" json:"stok"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Transaksi merepresentasikan satu transaksi penjualan (header).
type Transaksi struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	NoStruk   string         `gorm:"uniqueIndex;not null" json:"no_struk"`
	Total     int            `gorm:"not null" json:"total"`   // Total belanja (rupiah)
	Bayar     int            `gorm:"not null" json:"bayar"`   // Uang diterima
	Kembali   int            `gorm:"not null" json:"kembali"` // Kembalian
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relasi: satu transaksi punya banyak detail
	Details []DetailTransaksi `gorm:"foreignKey:TransaksiID" json:"details,omitempty"`
}

// DetailTransaksi merepresentasikan satu item dalam transaksi.
type DetailTransaksi struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	TransaksiID uint      `gorm:"not null;index" json:"transaksi_id"`
	ProdukID    uint      `gorm:"not null" json:"produk_id"`
	NamaProduk  string    `gorm:"not null" json:"nama_produk"`  // Snapshot nama saat checkout
	HargaSatuan int       `gorm:"not null" json:"harga_satuan"` // Snapshot harga saat checkout
	Jumlah      int       `gorm:"not null" json:"jumlah"`
	Subtotal    int       `gorm:"not null" json:"subtotal"`
	CreatedAt   time.Time `json:"created_at"`

	// Relasi balik
	Transaksi Transaksi `gorm:"foreignKey:TransaksiID" json:"-"`
	Produk    Produk    `gorm:"foreignKey:ProdukID" json:"-"`
}
