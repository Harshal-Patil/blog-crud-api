package routes

import (
	"blog-crud-api/config"
	"blog-crud-api/handlers"
	"blog-crud-api/repository"
	"blog-crud-api/services"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(router fiber.Router) {
	repo := repository.NewBlogRepository(config.DB)
	service := services.NewBlogService(repo)
	handler := handlers.NewBlogHandler(service)

	group := router.Group("/blog-post") // now /api/blog-post
	group.Post("/", handler.Create)
	group.Get("/", handler.GetAll)
	group.Get("/:id", handler.GetByID)
	group.Patch("/:id", handler.Update)
	group.Delete("/:id", handler.Delete)
}
