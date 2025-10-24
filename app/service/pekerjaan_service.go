package service

import (
	"database/sql"
	"prak4/app/model"
	"prak4/app/repository"
	"prak4/helper"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type PekerjaanService struct {
	Repo *repository.PekerjaanRepository
}

func NewPekerjaanService(repo *repository.PekerjaanRepository) *PekerjaanService {
	return &PekerjaanService{Repo: repo}
}

func (s *PekerjaanService) GetAllPekerjaan(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	sortBy := c.Query("sortBy", "id")
	order := c.Query("order", "asc")
	search := c.Query("search", "")

	role := c.Locals("role").(string)
	userID := c.Locals("user_id").(int)

	var pekerjaan []model.PekerjaanAlumni
	var total int
	var err error
	
	if role == "admin" {
		// Admin melihat semua data
		pekerjaan, err = s.Repo.GetAllPekerjaan(search, sortBy, order, limit, (page-1)*limit)
		if err == nil {
			total, err = s.Repo.CountPekerjaan(search)
		}
	} else {
		// User biasa hanya melihat datanya sendiri
		pekerjaan, err = s.Repo.GetAllPekerjaanForUser(userID, search, sortBy, order, limit, (page-1)*limit)
		if err == nil {
			total, err = s.Repo.CountPekerjaanForUser(userID, search)
		}
	}

	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Gagal mengambil data pekerjaan")
	}

	response := model.PekerjaanResponse{
		Data: pekerjaan,
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

func (s *PekerjaanService) GetPekerjaanByID(c *fiber.Ctx) error {
	pekerjaanID, _ := strconv.Atoi(c.Params("id"))
	role := c.Locals("role").(string)
	userID := c.Locals("user_id").(int)
	
	pekerjaan, err := s.Repo.GetPekerjaanByID(pekerjaanID)
	if err != nil {
		if err == sql.ErrNoRows {
			return helper.ErrorResponse(c, fiber.StatusNotFound, "Pekerjaan tidak ditemukan")
		}
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Gagal mengambil data pekerjaan")
	}
	
	// LOGIKA OTORISASI
	if role != "admin" && pekerjaan.AlumniID != userID {
		return helper.ErrorResponse(c, fiber.StatusForbidden, "Akses ditolak. Anda tidak memiliki izin untuk melihat data ini.")
	}

	return helper.SuccessResponse(c, pekerjaan, "Data pekerjaan berhasil diambil")
}

func (s *PekerjaanService) GetPekerjaanByAlumniID(c *fiber.Ctx) error {
	alumniID, err := strconv.Atoi(c.Params("alumni_id"))
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "ID alumni tidak valid")
	}

	pekerjaan, err := s.Repo.GetPekerjaanByAlumniID(alumniID)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Gagal mengambil data pekerjaan: "+err.Error())
	}
	return helper.SuccessResponse(c, pekerjaan, "Data pekerjaan untuk alumni berhasil diambil")
}

func (s *PekerjaanService) CreatePekerjaan(c *fiber.Ctx) error {
	var p model.PekerjaanAlumni
	if err := c.BodyParser(&p); err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Request body tidak valid")
	}

	if p.AlumniID == 0 || p.NamaPerusahaan == "" || p.PosisiJabatan == "" {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Field alumni_id, nama_perusahaan, dan posisi_jabatan wajib diisi")
	}

	newID, err := s.Repo.CreatePekerjaan(&p)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Gagal menambah pekerjaan: "+err.Error())
	}

	newPekerjaan, _ := s.Repo.GetPekerjaanByID(newID)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Pekerjaan berhasil ditambahkan",
		"data":    newPekerjaan,
	})
}

func (s *PekerjaanService) UpdatePekerjaan(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "ID tidak valid")
	}

	var p model.PekerjaanAlumni
	if err := c.BodyParser(&p); err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Request body tidak valid")
	}

	rowsAffected, err := s.Repo.UpdatePekerjaan(id, &p)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Gagal mengupdate pekerjaan: "+err.Error())
	}
	if rowsAffected == 0 {
		return helper.ErrorResponse(c, fiber.StatusNotFound, "Pekerjaan tidak ditemukan")
	}

	updatedPekerjaan, _ := s.Repo.GetPekerjaanByID(id)
	return helper.SuccessResponse(c, updatedPekerjaan, "Pekerjaan berhasil diupdate")
}

func (s *PekerjaanService) DeletePekerjaan(c *fiber.Ctx) error {
	pekerjaanID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "ID pekerjaan tidak valid")
	}

	role := c.Locals("role").(string)
	userID := c.Locals("user_id").(int)

	var rowsAffected int64

	// logika hak akses
	if role == "admin" {
		// Admin menggunakan fungsi hapus umum
		rowsAffected, err = s.Repo.SoftDeletePekerjaanByID(pekerjaanID)
	} else {
		// User biasa menggunakan fungsi hapus yang spesifik dengan ID mereka
		rowsAffected, err = s.Repo.SoftDeletePekerjaanForUser(pekerjaanID, userID)
	}

	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Gagal menghapus pekerjaan: "+err.Error())
	}

	if rowsAffected == 0 {
		return helper.ErrorResponse(c, fiber.StatusNotFound, "Pekerjaan tidak ditemukan atau Anda tidak memiliki hak akses.")
	}

	return c.JSON(fiber.Map{"success": true, "message": "Pekerjaan berhasil dihapus"})
}


//INI UTS

func (s *PekerjaanService) GetTrashedPekerjaan(c *fiber.Ctx) error {
	trashed, err := s.Repo.GetTrashedPekerjaan()
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Gagal mengambil data sampah")
	}
	return helper.SuccessResponse(c, trashed, "Data sampah berhasil diambil")
}

func (s *PekerjaanService) GetMyTrashedPekerjaan(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	trashed, err := s.Repo.GetTrashedPekerjaanForUser(userID)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Gagal mengambil data sampah Anda")
	}

	return helper.SuccessResponse(c, trashed, "Data sampah Anda berhasil diambil")
}

func (s *PekerjaanService) RestorePekerjaan(c *fiber.Ctx) error {
	pekerjaanID, _ := strconv.Atoi(c.Params("id"))
	role := c.Locals("role").(string)
	userID := c.Locals("user_id").(int)
	pekerjaan, err := s.Repo.GetPekerjaanByIDIncludeTrashed(pekerjaanID)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusNotFound, "Pekerjaan tidak ditemukan")
	}

	if pekerjaan.DeletedAt == nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Data pekerjaan ini aktif dan tidak perlu di-restore.")
	}
	
	// Otorisasi
	if role != "admin" && pekerjaan.AlumniID != userID {
		return helper.ErrorResponse(c, fiber.StatusForbidden, "Akses ditolak.")
	}

	_, err = s.Repo.RestorePekerjaanByID(pekerjaanID)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Gagal me-restore pekerjaan")
	}

	return c.JSON(fiber.Map{"success": true, "message": "Pekerjaan berhasil di-restore"})
}

func (s *PekerjaanService) HardDeletePekerjaan(c *fiber.Ctx) error {
	pekerjaanID, _ := strconv.Atoi(c.Params("id"))
	role := c.Locals("role").(string)
	userID := c.Locals("user_id").(int)

	pekerjaan, err := s.Repo.GetPekerjaanByIDIncludeTrashed(pekerjaanID)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusNotFound, "Pekerjaan tidak ditemukan di sampah")
	}

	if pekerjaan.DeletedAt == nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Data harus dihapus (soft delete) terlebih dahulu sebelum bisa dihapus permanen.")
	}

	if role != "admin" && pekerjaan.AlumniID != userID {
		return helper.ErrorResponse(c, fiber.StatusForbidden, "Akses ditolak.")
	}

	_, err = s.Repo.HardDeletePekerjaanByID(pekerjaanID)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Gagal menghapus permanen pekerjaan")
	}
	
	return c.JSON(fiber.Map{"success": true, "message": "Pekerjaan berhasil dihapus permanen"})
}