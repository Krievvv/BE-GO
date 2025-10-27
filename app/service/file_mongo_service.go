package service

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"prak4/app/model"
	repoMongo "prak4/app/repository" 
	"prak4/helper"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const (
	uploadPathFoto       = "./uploads/foto"
	uploadPathSertifikat = "./uploads/sertifikat"
	maxSizeFoto          = 1 * 1024 * 1024 // 1 MB
	maxSizeSertifikat    = 2 * 1024 * 1024 // 2 MB
)

type FileMongoService struct {
	repo *repoMongo.FileRepositoryMongo
}

func NewFileMongoService(repo *repoMongo.FileRepositoryMongo) *FileMongoService {
	os.MkdirAll(uploadPathFoto, os.ModePerm)
	os.MkdirAll(uploadPathSertifikat, os.ModePerm)
	return &FileMongoService{repo: repo}
}

// Handler Upload
func (s *FileMongoService) UploadFoto(c *fiber.Ctx) error {
	allowedTypes := map[string]bool{"image/jpeg": true, "image/jpg": true, "image/png": true}
	return s.handleUpload(c, "foto", uploadPathFoto, maxSizeFoto, allowedTypes)
}

func (s *FileMongoService) UploadSertifikat(c *fiber.Ctx) error {
	allowedTypes := map[string]bool{"application/pdf": true}
	return s.handleUpload(c, "sertifikat", uploadPathSertifikat, maxSizeSertifikat, allowedTypes)
}

// Handler GET & DELETE
func (s *FileMongoService) GetAllFiles(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	files, err := s.repo.FindAll(ctx)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Gagal mengambil daftar file")
	}

	responses := []model.FileResponse{}
	for _, file := range files {
		responses = append(responses, s.toFileResponse(&file))
	}
	return helper.SuccessResponse(c, responses, "Daftar file berhasil diambil")
}

func (s *FileMongoService) GetFileByID(c *fiber.Ctx) error {
	id := c.Params("id")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	file, err := s.repo.FindByID(ctx, id)
	if err != nil { return helper.ErrorResponse(c, fiber.StatusInternalServerError, err.Error()) }
	if file == nil { return helper.ErrorResponse(c, fiber.StatusNotFound, "File tidak ditemukan") }

	return helper.SuccessResponse(c, s.toFileResponse(file), "File berhasil diambil")
}

func (s *FileMongoService) DeleteFile(c *fiber.Ctx) error {
	id := c.Params("id")
	loggedInUserID := c.Locals("user_id").(int)
	role := c.Locals("role").(string)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	file, err := s.repo.FindByID(ctx, id)
	if err != nil { return helper.ErrorResponse(c, fiber.StatusInternalServerError, err.Error()) }
	if file == nil { return helper.ErrorResponse(c, fiber.StatusNotFound, "File tidak ditemukan") }

	// Otorisasi
	if role != "admin" && file.UserID != loggedInUserID {
		return helper.ErrorResponse(c, fiber.StatusForbidden, "Akses ditolak")
	}

	// Hapus file fisik
	err = os.Remove(file.FilePath)
	if err != nil {
		fmt.Printf("Peringatan: Gagal menghapus file fisik %s: %v\n", file.FilePath, err)
	}

	// Hapus metadata dari DB
	_, err = s.repo.Delete(ctx, file.ID)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Gagal menghapus metadata file")
	}

	return c.JSON(fiber.Map{"success": true, "message": "File berhasil dihapus"})
}


// Fungsi Helper Internal

// Logika untuk upload
func (s *FileMongoService) handleUpload(c *fiber.Ctx, jenis string, uploadDir string, maxSize int64, allowedTypes map[string]bool) error {
	// Dapatkan file dari form
	fileHeader, err := c.FormFile("file") 
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Tidak ada file yang diupload atau key salah (gunakan 'file')")
	}

	// Tentukan User ID target
	loggedInUserID := c.Locals("user_id").(int)
	role := c.Locals("role").(string)
	targetUserIDStr := c.FormValue("user_id", "") 
	targetUserID := loggedInUserID 

	if role == "admin" && targetUserIDStr != "" {
		parsedID, err := strconv.Atoi(targetUserIDStr)
		if err == nil {
			targetUserID = parsedID // Admin bisa menentukan target user ID
		} else {
			return helper.ErrorResponse(c, fiber.StatusBadRequest, "user_id yang diberikan tidak valid")
		}
	} else if role != "admin" && targetUserIDStr != "" {
		// User biasa tidak boleh menentukan user_id
		return helper.ErrorResponse(c, fiber.StatusForbidden, "Anda hanya dapat mengupload file untuk diri sendiri.")
	}


	// Validasi Ukuran
	if fileHeader.Size > maxSize {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, fmt.Sprintf("Ukuran file melebihi batas maksimal (%.1f MB)", float64(maxSize)/(1024*1024)))
	}

	// Validasi Tipe
	contentType := fileHeader.Header.Get("Content-Type")
	if !allowedTypes[contentType] {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, fmt.Sprintf("Tipe file '%s' tidak diizinkan untuk %s", contentType, jenis))
	}

	// Generate nama file unik & path
	ext := filepath.Ext(fileHeader.Filename)
	newFileName := uuid.New().String() + ext
	filePath := filepath.Join(uploadDir, newFileName)

	// Simpan file fisik
	if err := c.SaveFile(fileHeader, filePath); err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Gagal menyimpan file: "+err.Error())
	}

	// Siapkan metadata
	fileModel := &model.File{
		UserID:       targetUserID,
		FileName:     newFileName,
		OriginalName: fileHeader.Filename,
		FilePath:     filePath,
		FileSize:     fileHeader.Size,
		FileType:     contentType,
		UploadedAt:   time.Now(),
		Jenis:        jenis,
	}

	// Simpan metadata ke DB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	createdFile, err := s.repo.Create(ctx, fileModel)
	if err != nil {
		// Jika gagal simpan DB, hapus file fisik yang sudah terlanjur disimpan
		os.Remove(filePath)
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Gagal menyimpan metadata file: "+err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": fmt.Sprintf("%s berhasil diupload", jenis),
		"data":    s.toFileResponse(createdFile),
	})
}


// Mengubah model DB ke model respons API
func (s *FileMongoService) toFileResponse(file *model.File) model.FileResponse {
	// Membuat URL relatif untuk akses file via static server
	urlPath := strings.Replace(file.FilePath, ".", "", 1) 

	return model.FileResponse{
		ID:           file.ID.Hex(),
		UserID:       file.UserID,
		FileName:     file.FileName,
		OriginalName: file.OriginalName,
		FilePath:     urlPath, 
		FileSize:     file.FileSize,
		FileType:     file.FileType,
		UploadedAt:   file.UploadedAt,
		Jenis:        file.Jenis,
	}
}