package handlers

import (
	"blog-crud-api/config"
	"blog-crud-api/models"
	"blog-crud-api/services"
	"blog-crud-api/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"
)

type BlogHandler struct {
	Service services.BlogService
}

func NewBlogHandler(s services.BlogService) *BlogHandler {
	return &BlogHandler{Service: s}
}

// @Summary Create a new blog post
// @Tags Blog
// @Accept json
// @Produce json
// @Param post body models.BlogPost true "Blog Post"
// @Success 200 {object} models.BlogPost
// @Router /blog-post [post]
func (h *BlogHandler) Create(c *fiber.Ctx) error {
	var post models.BlogPost
	if err := c.BodyParser(&post); err != nil {
		utils.Log.WithError(err).Warn("Invalid input for blog post creation")
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Extract user ID from JWT
	user := c.Locals("user").(*jwt.Token)
	utils.Log.WithFields(logrus.Fields{
		"user_id": user,
		"title":   post.Title,
	}).Info("Creating blog post")
	claims := user.Claims.(jwt.MapClaims)
	userID := uint(claims["id"].(float64))
	post.UserID = userID

	// Check for duplicate
	var existing models.BlogPost
	err := config.DB.Where("title = ? AND description = ? AND body = ? AND user_id = ?", post.Title, post.Description, post.Body, userID).First(&existing).Error
	if err == nil {
		return c.Status(409).JSON(fiber.Map{"error": "Duplicate blog content"})
	}

	created, err := h.Service.Create(post)
	if err != nil {
		utils.Log.WithError(err).Error("Failed to create blog post")
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(created)
}

func (h *BlogHandler) GetAll(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	posts, err := h.Service.GetAllWithPagination(limit, offset)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(posts)
}

func (h *BlogHandler) GetByID(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	post, err := h.Service.GetByID(uint(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Not found"})
	}
	return c.JSON(post)
}
func (h *BlogHandler) Update(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var data models.BlogPost
	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}
	updated, err := h.Service.Update(uint(id), data)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(updated)
}
func (h *BlogHandler) Delete(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	if err := h.Service.Delete(uint(id)); err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Not found"})
	}
	return c.SendStatus(204)
}
func (h *BlogHandler) GetByUser(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := uint(claims["id"].(float64))

	posts, err := h.Service.GetByUser(int(userID))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(posts)
}
