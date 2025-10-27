package route

import (
	"database/sql"
	"prak4/app/repository"
	"prak4/app/service"
	"prak4/database"
	"prak4/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, db *sql.DB) {
	alumniRepo := repository.NewAlumniRepository(db)
	pgPekerjaanRepo := repository.NewPekerjaanRepository(db)
	pgUserRepo := repository.NewUserRepository(db)

	alumniService := service.NewAlumniService(alumniRepo)
	pgPekerjaanService := service.NewPekerjaanService(pgPekerjaanRepo)

	mongoPekerjaanRepo := repository.NewPekerjaanRepositoryMongo(database.Mongo)
	mongoUserRepo := repository.NewUserRepositoryMongo(database.Mongo)
	logRepo := repository.NewLogRepository(database.Mongo)

	mongoPekerjaanService := service.NewPekerjaanMongoService(mongoPekerjaanRepo)
	logService := service.NewLogService(logRepo)

	pgAuthService := service.NewAuthService(pgUserRepo, logRepo)
	mongoAuthService := service.NewAuthMongoService(mongoUserRepo, logRepo)

	// File Upload (Mongo)
	fileMongoRepo := repository.NewFileRepositoryMongo(database.Mongo)
	fileMongoService := service.NewFileMongoService(fileMongoRepo)

	api := app.Group("/prak4")

	api.Post("/login", pgAuthService.Login)
	api.Post("/login-mongo", mongoAuthService.LoginMongo)

	protected := api.Group("", middleware.AuthRequired())
	protected.Get("/logs", middleware.AdminOnly(), logService.GetAllLogs)
	protected.Post("/users", middleware.AdminOnly(), middleware.RequireIssuer("postgres"), pgAuthService.RegisterUser)

	// ALUMNI (Postgres)
	alumni := protected.Group("/alumni", middleware.RequireIssuer("postgres"))
	alumni.Get("/", alumniService.GetAllAlumni)
	alumni.Get("/angkatan/:angkatan", alumniService.GetAlumniByAngkatan)
	alumni.Get("/:id", alumniService.GetAlumniByID)
	alumni.Post("/", middleware.AdminOnly(), alumniService.CreateAlumni)
	alumni.Put("/:id", alumniService.UpdateAlumni)
	alumni.Delete("/:id", alumniService.DeleteAlumni)

	// PEKERJAAN (Postgres)
	pekerjaan := protected.Group("/pekerjaan", middleware.RequireIssuer("postgres"))
	pekerjaan.Get("/", pgPekerjaanService.GetAllPekerjaan)
	pekerjaan.Get("/trash", middleware.AdminOnly(), pgPekerjaanService.GetTrashedPekerjaan)
	pekerjaan.Get("/trash/me", pgPekerjaanService.GetMyTrashedPekerjaan)
	pekerjaan.Get("/:id", pgPekerjaanService.GetPekerjaanByID)
	pekerjaan.Post("/", middleware.AdminOnly(), pgPekerjaanService.CreatePekerjaan)
	pekerjaan.Put("/:id", middleware.AdminOnly(), pgPekerjaanService.UpdatePekerjaan)
	pekerjaan.Patch("/restore/:id", pgPekerjaanService.RestorePekerjaan)
	pekerjaan.Delete("/:id", pgPekerjaanService.DeletePekerjaan)
	pekerjaan.Delete("/force/:id", pgPekerjaanService.HardDeletePekerjaan)

	// PEKERJAAN (MongoDB)
	pMongo := protected.Group("/mongo/pekerjaan", middleware.RequireIssuer("mongo"))
	pMongo.Post("/", middleware.AdminOnly(), mongoPekerjaanService.CreatePekerjaan)
	pMongo.Get("/", mongoPekerjaanService.GetAllPekerjaan)
	pMongo.Get("/trash", mongoPekerjaanService.GetTrashed)
	pMongo.Get("/alumni/:alumni_id", middleware.AdminOnly(), mongoPekerjaanService.GetAllPekerjaanByAlumniID)
	pMongo.Get("/:id", mongoPekerjaanService.GetPekerjaanByID)
	pMongo.Put("/:id", middleware.AdminOnly(), mongoPekerjaanService.UpdatePekerjaan)
	pMongo.Patch("/restore/:id", mongoPekerjaanService.RestorePekerjaan)
	pMongo.Delete("/:id", mongoPekerjaanService.DeletePekerjaan)
	pMongo.Delete("/force/:id", mongoPekerjaanService.HardDeletePekerjaan)

	// FILE UPLOAD (Mongo) 
	filesGroup := protected.Group("/files", middleware.RequireIssuer("mongo"))

	// Upload menggunakan form-data
	filesGroup.Post("/upload/foto", fileMongoService.UploadFoto)
	filesGroup.Post("/upload/sertifikat", fileMongoService.UploadSertifikat)

	// Get & Delete file metadata
	filesGroup.Get("/", fileMongoService.GetAllFiles)
	filesGroup.Get("/:id", fileMongoService.GetFileByID)
	filesGroup.Delete("/:id", fileMongoService.DeleteFile)
}
