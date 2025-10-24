package service

import (
	"database/sql"
	"prak4/app/model"
	"prak4/app/repository"
	"prak4/helper"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type AlumniService struct {
	Repo *repository.AlumniRepository
}

// NewAlumniService creates a new instance of AlumniService
func NewAlumniService(repo *repository.AlumniRepository) *AlumniService {
	return &AlumniService{Repo: repo}
}

func (s *AlumniService) GetAllAlumni(c *fiber.Ctx) error {
	// Ambil query parameter
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	sortBy := c.Query("sortBy", "id") 
	order := c.Query("order", "asc")   
	search := c.Query("search", "")

	// Validasi parameter
	sortByWhitelist := map[string]bool{"id": true, "nama": true, "nim": true, "angkatan": true, "jurusan": true}
	if !sortByWhitelist[sortBy] {
		sortBy = "id" 
	}
	if strings.ToLower(order) != "desc" {
		order = "asc"
	}

	// Hitung offset
	offset := (page - 1) * limit

	// Panggil repository
	alumni, err := s.Repo.GetAllAlumni(search, sortBy, order, limit, offset)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Gagal mengambil data alumni")
	}

	total, err := s.Repo.CountAlumni(search)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Gagal menghitung total data alumni")
	}

	response := model.AlumniResponse{
		Data: alumni,
		Meta: model.MetaInfo{
			Page:   page,
			Limit:  limit,
			Total:  total,
			Pages:  (total + limit - 1) / limit, 
			SortBy: sortBy,
			Order:  order,
			Search: search,
		},
	}

	return c.JSON(response) 
}

func (s *AlumniService) GetAlumniByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "ID tidak valid")
	}

	alumni, err := s.Repo.GetAlumniByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return helper.ErrorResponse(c, fiber.StatusNotFound, "Alumni tidak ditemukan")
		}
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Gagal mengambil data alumni: "+err.Error())
	}
	return helper.SuccessResponse(c, alumni, "Data alumni berhasil diambil")
}

func (s *AlumniService) GetAlumniByAngkatan(c *fiber.Ctx) error {
	angkatan, err := strconv.Atoi(c.Params("angkatan"))
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Parameter angkatan tidak valid")
	}

	result, err := s.Repo.GetAlumniByAngkatan(angkatan)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Gagal mengambil data alumni: "+err.Error())
	}

	message := "Data jumlah alumni untuk angkatan " + strconv.Itoa(angkatan) + " berhasil diambil"
	return helper.SuccessResponse(c, result, message)
}

func (s *AlumniService) CreateAlumni(c *fiber.Ctx) error {
	var alumni model.Alumni
	if err := c.BodyParser(&alumni); err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Request body tidak valid")
	}

	if alumni.NIM == "" || alumni.Nama == "" || alumni.Jurusan == "" || alumni.Email == "" {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Field NIM, Nama, Jurusan, dan Email wajib diisi")
	}

	newID, err := s.Repo.CreateAlumni(&alumni)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Gagal menambah alumni: "+err.Error())
	}

	newAlumni, _ := s.Repo.GetAlumniByID(newID)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Alumni berhasil ditambahkan",
		"data":    newAlumni,
	})
}

func (s *AlumniService) UpdateAlumni(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "ID tidak valid")
	}

	var alumni model.Alumni
	if err := c.BodyParser(&alumni); err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Request body tidak valid")
	}

	rowsAffected, err := s.Repo.UpdateAlumni(id, &alumni)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Gagal mengupdate alumni: "+err.Error())
	}
	if rowsAffected == 0 {
		return helper.ErrorResponse(c, fiber.StatusNotFound, "Alumni tidak ditemukan untuk diupdate")
	}

	updatedAlumni, _ := s.Repo.GetAlumniByID(id)
	return helper.SuccessResponse(c, updatedAlumni, "Alumni berhasil diupdate")
}

func (s *AlumniService) DeleteAlumni(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "ID tidak valid")
	}

	rowsAffected, err := s.Repo.DeleteAlumni(id)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Gagal menghapus alumni: "+err.Error())
	}
	if rowsAffected == 0 {
		return helper.ErrorResponse(c, fiber.StatusNotFound, "Alumni tidak ditemukan untuk dihapus")
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Alumni berhasil dihapus",
	})
}