package routes

import (
	"github.com/gofiber/fiber/v2"

	"lgo/handlers"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/", handlers.Home)

	app.Get("/login", handlers.ShowLogin)
	app.Post("/login", handlers.Login)

	app.Get("/signup", handlers.ShowSignup)
	app.Post("/signup", handlers.Signup)

	app.Get("/logout", handlers.Logout)
}
