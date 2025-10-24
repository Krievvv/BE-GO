package service

import (
	"context"
	"prak4/app/model"
	"prak4/app/repository"
	"prak4/helper"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

type PekerjaanMongoService struct {
	repo *repository.PekerjaanRepositoryMongo
}

func NewPekerjaanMongoService(repo *repository.PekerjaanRepositoryMongo) *PekerjaanMongoService {
	return &PekerjaanMongoService{repo: repo}
}

// === CREATE ===
func (s *PekerjaanMongoService) CreatePekerjaan(c *fiber.Ctx) error {
	var p model.PekerjaanAlumniMongo
	if err := c.BodyParser(&p); err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Request body tidak valid")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	created, err := s.repo.Create(ctx, &p)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Gagal membuat data: "+err.Error())
	}
	return c.Status(fiber.StatusCreated).JSON(created)
}

// === READ ===
func (s *PekerjaanMongoService) GetAllPekerjaan(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	sortBy := c.Query("sortBy", "_id")
	orderStr := c.Query("order", "asc")
	search := c.Query("search", "")

	order := 1 // asc
	if strings.ToLower(orderStr) == "desc" {
		order = -1
	}

	role := c.Locals("role").(string)
	userID := c.Locals("user_id").(int)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pekerjaan, total, err := s.repo.FindAll(ctx, userID, role == "admin", search, sortBy, order, limit, (page-1)*limit)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Gagal mengambil data pekerjaan")
	}

	response := model.PekerjaanMongoResponse{
		Data: pekerjaan,
		Meta: model.MetaInfo{
			Page: page, Limit: limit, Total: int(total), Pages: (int(total) + limit - 1) / limit,
			SortBy: sortBy, Order: orderStr, Search: search,
		},
	}
	return c.JSON(response)
}

func (s *PekerjaanMongoService) GetPekerjaanByID(c *fiber.Ctx) error {
	id := c.Params("id")
	role := c.Locals("role").(string)
	userID := c.Locals("user_id").(int)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pekerjaan, err := s.repo.FindByID(ctx, id, false)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}
	if pekerjaan == nil {
		return helper.ErrorResponse(c, fiber.StatusNotFound, "Pekerjaan tidak ditemukan")
	}
	if role != "admin" && pekerjaan.AlumniID != userID {
		return helper.ErrorResponse(c, fiber.StatusForbidden, "Akses ditolak")
	}

	return helper.SuccessResponse(c, pekerjaan, "Data berhasil diambil")
}

func (s *PekerjaanMongoService) GetAllPekerjaanByAlumniID(c *fiber.Ctx) error {
	alumniID, _ := strconv.Atoi(c.Params("alumni_id"))
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	pekerjaan, err := s.repo.FindAllByAlumniID(ctx, alumniID)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Gagal mengambil data")
	}
	return helper.SuccessResponse(c, pekerjaan, "Data berhasil diambil")
}

// === UPDATE ===
func (s *PekerjaanMongoService) UpdatePekerjaan(c *fiber.Ctx) error {
	id := c.Params("id")
	var p model.PekerjaanAlumniMongo
	if err := c.BodyParser(&p); err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Request body tidak valid")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := s.repo.Update(ctx, id, &p)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}
	if result.MatchedCount == 0 {
		return helper.ErrorResponse(c, fiber.StatusNotFound, "Pekerjaan tidak ditemukan")
	}

	return c.JSON(fiber.Map{"message": "Data pekerjaan berhasil diperbarui"})
}

func (s *PekerjaanMongoService) RestorePekerjaan(c *fiber.Ctx) error {
	id := c.Params("id")
	role := c.Locals("role").(string)
	userID := c.Locals("user_id").(int)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	p, err := s.repo.FindByID(ctx, id, true)
	if err != nil || p == nil {
		return helper.ErrorResponse(c, fiber.StatusNotFound, "Pekerjaan tidak ditemukan di sampah")
	}
	if p.DeletedAt == nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Pekerjaan tidak ada di sampah")
	}
	if role != "admin" && p.AlumniID != userID {
		return helper.ErrorResponse(c, fiber.StatusForbidden, "Akses ditolak")
	}

	_, err = s.repo.Restore(ctx, id)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Gagal me-restore pekerjaan")
	}

	return c.JSON(fiber.Map{"message": "Pekerjaan berhasil di-restore"})
}

// DELETE (SOFT & HARD)
func (s *PekerjaanMongoService) DeletePekerjaan(c *fiber.Ctx) error {
	return s.handleDelete(c, false)
}

func (s *PekerjaanMongoService) HardDeletePekerjaan(c *fiber.Ctx) error {
	return s.handleDelete(c, true)
}

func (s *PekerjaanMongoService) handleDelete(c *fiber.Ctx, isHardDelete bool) error {
	id := c.Params("id")
	role := c.Locals("role").(string)
	userID := c.Locals("user_id").(int)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	p, err := s.repo.FindByID(ctx, id, true) 
	if err != nil || p == nil {
		return helper.ErrorResponse(c, fiber.StatusNotFound, "Pekerjaan tidak ditemukan")
	}
	if role != "admin" && p.AlumniID != userID {
		return helper.ErrorResponse(c, fiber.StatusForbidden, "Akses ditolak")
	}

	if isHardDelete {
		if p.DeletedAt == nil {
			return helper.ErrorResponse(c, fiber.StatusBadRequest, "Data harus di-soft delete dulu")
		}
		_, err = s.repo.HardDelete(ctx, id)
	} else {
		if p.DeletedAt != nil {
			return helper.ErrorResponse(c, fiber.StatusBadRequest, "Data sudah ada di sampah")
		}
		_, err = s.repo.SoftDelete(ctx, id)
	}

	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Gagal menghapus pekerjaan")
	}

	msg := "Pekerjaan berhasil di-soft delete"
	if isHardDelete {
		msg = "Pekerjaan berhasil di-hard delete"
	}
	return c.JSON(fiber.Map{"message": msg})
}

// === TRASH ===
func (s *PekerjaanMongoService) GetTrashed(c *fiber.Ctx) error {
	role := c.Locals("role").(string)
	userID := c.Locals("user_id").(int)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	trash, err := s.repo.FindTrashed(ctx, userID, role == "admin")
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Gagal mengambil data sampah")
	}
	return helper.SuccessResponse(c, trash, "Data sampah berhasil diambil")
}