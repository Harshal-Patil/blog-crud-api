package handlers

import (
	"blog-crud-api/config"
	"blog-crud-api/models"
	"blog-crud-api/services"
	"blog-crud-api/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type BlogHandler struct {
	Service services.BlogService
}

func NewBlogHandler(s services.BlogService) *BlogHandler {
	return &BlogHandler{Service: s}
}

// Create godoc
// @Summary Create a new blog post
// @Tags Blog
// @Accept json
// @Produce json
// @Param post body models.BlogPost true "Blog Post"
// @Success 201 {object} utils.APIResponse{data=models.BlogPost}
// @Failure 400 {object} utils.APIResponse
// @Failure 409 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Security ApiKeyAuth
// @Router /blog-post [post]
func (h *BlogHandler) Create(c *fiber.Ctx) error {
	var post models.BlogPost
	if err := c.BodyParser(&post); err != nil {
		utils.Log.WithError(err).Warn("Invalid input for blog post creation")
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse("Invalid input"))
	}

	// Extract user ID from JWT token
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := uint(claims["id"].(float64))
	post.UserID = userID

	// Check for duplicate post by same user
	var existing models.BlogPost
	err := config.DB.Where("title = ? AND description = ? AND body = ? AND user_id = ?", post.Title, post.Description, post.Body, userID).First(&existing).Error
	if err == nil {
		return c.Status(fiber.StatusConflict).JSON(utils.ErrorResponse("Duplicate blog content"))
	}

	created, err := h.Service.Create(post)
	if err != nil {
		utils.Log.WithError(err).Error("Failed to create blog post")
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse(err.Error()))
	}
	return c.Status(fiber.StatusCreated).JSON(utils.SuccessResponse(created))
}

// GetAll godoc
// @Summary Get all blog posts with pagination
// @Tags Blog
// @Produce json
// @Param limit query int false "Limit number" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} utils.APIResponse{data=[]models.BlogPost}
// @Failure 500 {object} utils.APIResponse
// @Router /blog-post [get]
func (h *BlogHandler) GetAll(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	posts, err := h.Service.GetAllWithPagination(limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse(err.Error()))
	}
	return c.JSON(utils.SuccessResponse(posts))
}

// GetByID godoc
// @Summary Get a blog post by ID
// @Tags Blog
// @Produce json
// @Param id path int true "Blog Post ID"
// @Success 200 {object} utils.APIResponse{data=models.BlogPost}
// @Failure 404 {object} utils.APIResponse
// @Router /blog-post/{id} [get]
func (h *BlogHandler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse("Invalid post ID"))
	}

	post, err := h.Service.GetByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(utils.ErrorResponse("Post not found"))
	}
	return c.JSON(utils.SuccessResponse(post))
}

// Update godoc
// @Summary Update a blog post by ID
// @Tags Blog
// @Accept json
// @Produce json
// @Param id path int true "Blog Post ID"
// @Param post body models.BlogPost true "Updated Blog Post"
// @Success 200 {object} utils.APIResponse{data=models.BlogPost}
// @Failure 400 {object} utils.APIResponse
// @Failure 404 {object} utils.APIResponse
// @Router /blog-post/{id} [patch]
// @Security ApiKeyAuth
func (h *BlogHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse("Invalid post ID"))
	}

	var data models.BlogPost
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse("Invalid input"))
	}

	updated, err := h.Service.Update(uint(id), data)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(utils.ErrorResponse(err.Error()))
	}
	return c.JSON(utils.SuccessResponse(updated))
}

// Delete godoc
// @Summary Delete a blog post by ID
// @Tags Blog
// @Produce json
// @Param id path int true "Blog Post ID"
// @Success 204
// @Failure 404 {object} utils.APIResponse
// @Router /blog-post/{id} [delete]
// @Security ApiKeyAuth
func (h *BlogHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse("Invalid post ID"))
	}

	if err := h.Service.Delete(uint(id)); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(utils.ErrorResponse("Post not found"))
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// GetByUser godoc
// @Summary Get blog posts by authenticated user
// @Tags Blog
// @Produce json
// @Success 200 {object} utils.APIResponse{data=[]models.BlogPost}
// @Failure 500 {object} utils.APIResponse
// @Security ApiKeyAuth
// @Router /blog-post/user [get]
func (h *BlogHandler) GetByUser(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := uint(claims["id"].(float64))

	posts, err := h.Service.GetByUser(int(userID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse(err.Error()))
	}
	return c.JSON(utils.SuccessResponse(posts))
}
