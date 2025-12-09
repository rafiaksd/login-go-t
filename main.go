package main

import (
	"fmt"
	"log"
	"golang.org/x/crypto/bcrypt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/gorilla/sessions"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var store = sessions.NewCookieStore([]byte("secret-key")) // Use a real secret key in production!

// Define User model
type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"unique;not null"`
	Email    string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
}

func main() {
	// Initialize Fiber app with template engine
	app := fiber.New(fiber.Config{
		Views: html.New("./views", ".html"),
	})

	// Initialize SQLite database
	db, err := gorm.Open("sqlite3", "./users.db")
	if err != nil {
		log.Fatalf("Error opening database: %v\n", err)
	}
	defer db.Close()

	// Migrate User table
	db.AutoMigrate(&User{})

	// Middleware to check for logged-in users
	app.Use(func(c *fiber.Ctx) error {
		session, err := store.Get(c.Request(), "session")
		if err != nil {
			return err
		}

		// Attach user info to the context if logged in
		if username, ok := session.Values["username"]; ok {
			c.Locals("user", username)
		}
		return c.Next()
	})

	// Show homepage
	app.Get("/", func(c *fiber.Ctx) error {
		username := c.Locals("user")
		return c.Render("index", fiber.Map{
			"User": username,
		})
	})

	// Route to display sign-up page
	app.Get("/signup", func(c *fiber.Ctx) error {
		return c.Render("signup", nil)
	})

	// Handle sign-up
	app.Post("/signup", func(c *fiber.Ctx) error {
		username := c.FormValue("username")
		email := c.FormValue("email")
		password := c.FormValue("password")

		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return c.Status(500).SendString("Error hashing password")
		}

		// Save new user in DB
		user := User{
			Username: username,
			Email:    email,
			Password: string(hashedPassword),
		}
		db.Create(&user)

		return c.Redirect("/login")
	})

	fmt.Println("Hi from login-go-t")
}
