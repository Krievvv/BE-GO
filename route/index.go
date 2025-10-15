package route

import (
	"database/sql"
	"prak4/app/repository"
	"prak4/app/service"
	"prak4/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, db *sql.DB) {
	// Initialize repositories
	alumniRepo := &repository.AlumniRepository{DB: db}
	pekerjaanRepo := &repository.PekerjaanRepository{DB: db}
	userRepo := &repository.UserRepository{DB: db} 

	// Initialize services
	alumniService := &service.AlumniService{Repo: alumniRepo}
	pekerjaanService := &service.PekerjaanService{Repo: pekerjaanRepo}
	authService := &service.AuthService{UserRepo: userRepo}

	// Grouping routes
	api := app.Group("/prak4")
	api.Post("/login", authService.Login)

	//rute yang dilindungi (membutuhkan login)
	protected := api.Group("", middleware.AuthRequired())
	
	// Alumni Routes [cite: 588]
	alumni := protected.Group("/alumni")
	alumni.Get("/", alumniService.GetAllAlumni)                                   
	alumni.Get("/:id", alumniService.GetAlumniByID)                               
	alumni.Post("/", middleware.AdminOnly(), alumniService.CreateAlumni)         
	alumni.Put("/:id", middleware.AdminOnly(), alumniService.UpdateAlumni)       
	alumni.Delete("/:id", middleware.AdminOnly(), alumniService.DeleteAlumni)    
	alumni.Get("/angkatan/:angkatan", alumniService.GetAlumniByAngkatan) 

	// Pekerjaan Alumni Routes [cite: 594]
	// pekerjaan := protected.Group("/pekerjaan")
	// pekerjaan.Get("/", pekerjaanService.GetAllPekerjaan)                                           
	// pekerjaan.Get("/:id", pekerjaanService.GetPekerjaanByID)                                       
	// pekerjaan.Get("/alumni/:alumni_id", middleware.AdminOnly(), pekerjaanService.GetPekerjaanByAlumniID) 
	// pekerjaan.Post("/", middleware.AdminOnly(), pekerjaanService.CreatePekerjaan)                  
	// pekerjaan.Put("/:id", middleware.AdminOnly(), pekerjaanService.UpdatePekerjaan)                
	// pekerjaan.Delete("/:id", pekerjaanService.DeletePekerjaan) 


	//INI UTS
	pekerjaan := protected.Group("/pekerjaan")
	pekerjaan.Get("/", pekerjaanService.GetAllPekerjaan)                                           
	pekerjaan.Get("/trash", middleware.AdminOnly(), pekerjaanService.GetTrashedPekerjaan)
	pekerjaan.Get("/trash/me", pekerjaanService.GetMyTrashedPekerjaan)
	pekerjaan.Get("/:id", pekerjaanService.GetPekerjaanByID)                                       
	pekerjaan.Get("/alumni/:alumni_id", middleware.AdminOnly(), pekerjaanService.GetPekerjaanByAlumniID) 
	pekerjaan.Post("/", middleware.AdminOnly(), pekerjaanService.CreatePekerjaan)                  
	pekerjaan.Put("/:id", middleware.AdminOnly(), pekerjaanService.UpdatePekerjaan)                
	pekerjaan.Patch("/restore/:id", pekerjaanService.RestorePekerjaan)            
	pekerjaan.Delete("/:id", pekerjaanService.DeletePekerjaan)         
	pekerjaan.Delete("/force/:id", pekerjaanService.HardDeletePekerjaan)
}