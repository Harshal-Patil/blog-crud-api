package main

import (
	"blog-crud-api/config"
	"blog-crud-api/handlers"
	"blog-crud-api/routes"
	"blog-crud-api/utils"
	"os"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
)

func main() {
	utils.InitLogger()
	log := utils.Log
	utils.Log.Info("Starting Blog API...")
	config.ConnectDatabase()

	app := fiber.New()

	// app.Get("/swagger/*", swagger.HandlerDefault)

	// Public route (optional login simulation)
	app.Post("/api/login", handlers.Login)
	app.Post("/api/register", handlers.Register)
	// Protect all /api/blog-post routes
	api := app.Group("/api")

	api.Use("/blog-post", jwtware.New(jwtware.Config{
		SigningKey: []byte("secret"), // replace with env key in production
		ContextKey: "user",           // this stores user info in ctx.Locals("user")
	}))

	routes.RegisterRoutes(api)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	utils.Log.Infof("Server listening on port %s", port)
	log.Fatal(app.Listen(":" + port))
}
