package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"

	"lgo/db"
	"lgo/models"
	"lgo/routes"
	"lgo/handlers"
)

func main() {
	app := fiber.New(fiber.Config{
		Views: html.New("./views", ".html"),
	})

	// Connect to DB
	db.ConnectDatabase()

	// Do migrations
	db.DB.AutoMigrate(&models.User{})

	// Session middleware (set user in locals)
	app.Use(func(c *fiber.Ctx) error {
		sess, _ := handlers.Store.Get(c)
		user := sess.Get("username")
		if user != nil {
			c.Locals("user", user)
		}
		return c.Next()
	})

	// Register all routes
	routes.SetupRoutes(app)

	log.Fatal(app.Listen(":3000"))
}
