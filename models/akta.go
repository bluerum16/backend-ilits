package models

import "time"

type AktaKelahiran struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	NomorAkta    string    `gorm:"uniqueIndex" json:"nomor_akta"`
	NamaAnak     string    `gorm:"not null" json:"nama_anak"`
	JenisKelamin string    `gorm:"size:1;not null" json:"jenis_kelamin"`
	TempatLahir  string    `gorm:"not null" json:"tempat_lahir"`
	TanggalLahir string    `gorm:"size:10;not null" json:"tanggal_lahir"`
	NamaAyah     string    `json:"nama_ayah"`
	NamaIbu      string    `gorm:"not null" json:"nama_ibu"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// AktaInput dipakai untuk create dan update, nomor akta tidak diisi client karena dibuat otomatis
type AktaInput struct {
	NamaAnak     string `json:"nama_anak" binding:"required"`
	JenisKelamin string `json:"jenis_kelamin" binding:"required,oneof=L P"`
	TempatLahir  string `json:"tempat_lahir" binding:"required"`
	TanggalLahir string `json:"tanggal_lahir" binding:"required,datetime=2006-01-02"`
	NamaAyah     string `json:"nama_ayah"`
	NamaIbu      string `json:"nama_ibu" binding:"required"`
}
