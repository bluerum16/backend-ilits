package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"backend-ilits/models"
)

type AktaHandler struct {
	DB *gorm.DB
}

func NewAktaHandler(db *gorm.DB) *AktaHandler {
	return &AktaHandler{DB: db}
}

func (h *AktaHandler) GetAllAkta(c *gin.Context) {
	var aktaList []models.AktaKelahiran

	query := h.DB
	if nama := c.Query("nama"); nama != "" {
		query = query.Where("nama_anak LIKE ?", "%"+nama+"%")
	}
	if tahun := c.Query("tahun"); tahun != "" {
		query = query.Where("tanggal_lahir LIKE ?", tahun+"-%")
	}

	if err := query.Order("tanggal_lahir desc").Find(&aktaList).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "gagal mengambil data akta"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "berhasil", "data": aktaList})
}

func (h *AktaHandler) GetAkta(c *gin.Context) {
	akta, ok := h.findAkta(c)
	if !ok {
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "berhasil", "data": akta})
}

func (h *AktaHandler) CreateAkta(c *gin.Context) {
	input, ok := bindAktaInput(c)
	if !ok {
		return
	}

	akta := models.AktaKelahiran{}
	fillAkta(&akta, input)

	err := h.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&akta).Error; err != nil {
			return err
		}
		// nomor akta butuh id, jadi dibuat setelah data tersimpan
		akta.NomorAkta = generateNomorAkta(akta.TanggalLahir, akta.ID)
		return tx.Save(&akta).Error
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "gagal menyimpan akta"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "akta kelahiran berhasil dicatat", "data": akta})
}

func (h *AktaHandler) UpdateAkta(c *gin.Context) {
	akta, ok := h.findAkta(c)
	if !ok {
		return
	}

	input, ok := bindAktaInput(c)
	if !ok {
		return
	}

	fillAkta(&akta, input)

	if err := h.DB.Save(&akta).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "gagal mengubah akta"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "akta kelahiran berhasil diubah", "data": akta})
}

func (h *AktaHandler) DeleteAkta(c *gin.Context) {
	akta, ok := h.findAkta(c)
	if !ok {
		return
	}

	if err := h.DB.Delete(&akta).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "gagal menghapus akta"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "akta kelahiran berhasil dihapus"})
}

// findAkta mencari akta berdasarkan param id, langsung kirim response 404 kalau tidak ada
func (h *AktaHandler) findAkta(c *gin.Context) (models.AktaKelahiran, bool) {
	var akta models.AktaKelahiran

	err := h.DB.First(&akta, c.Param("id")).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"message": "akta kelahiran tidak ditemukan"})
		return akta, false
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "terjadi kesalahan"})
		return akta, false
	}

	return akta, true
}

func bindAktaInput(c *gin.Context) (models.AktaInput, bool) {
	var input models.AktaInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "input tidak valid", "error": err.Error()})
		return input, false
	}

	lahir, _ := time.Parse("2006-01-02", input.TanggalLahir)
	if lahir.After(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "tanggal lahir tidak boleh di masa depan"})
		return input, false
	}

	return input, true
}

func fillAkta(akta *models.AktaKelahiran, input models.AktaInput) {
	akta.NamaAnak = input.NamaAnak
	akta.JenisKelamin = input.JenisKelamin
	akta.TempatLahir = input.TempatLahir
	akta.TanggalLahir = input.TanggalLahir
	akta.NamaAyah = input.NamaAyah
	akta.NamaIbu = input.NamaIbu
}

// generateNomorAkta membuat nomor dengan format AKL-<tanggal lahir ddmmyyyy>-<id 4 digit>
func generateNomorAkta(tanggalLahir string, id uint) string {
	parts := strings.Split(tanggalLahir, "-")
	tanggal := parts[2] + parts[1] + parts[0]
	return fmt.Sprintf("AKL-%s-%04d", tanggal, id)
}
