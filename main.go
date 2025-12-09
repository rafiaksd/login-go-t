package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/template/html/v2"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"golang.org/x/crypto/bcrypt"
)

var store = session.New()

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
		session, err := store.Get(c)
		if err != nil {
			return err
		}

		if username := session.Get("username"); username != nil {
			c.Locals("user", username)
		}

		return c.Next()
	})

	// Show homepage
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{
			"User": c.Locals("user"),
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

	// Route to display login page
	app.Get("/login", func(c *fiber.Ctx) error {
		return c.Render("login", nil)
	})

	// Handle login
	app.Post("/login", func(c *fiber.Ctx) error {
		email := c.FormValue("email")
		password := c.FormValue("password")

		// Find user by email
		var user User
		if err := db.Where("email = ?", email).First(&user).Error; err != nil {
			return c.Status(401).SendString("Invalid credentials")
		}

		// Compare the password
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
			return c.Status(401).SendString("Invalid credentials")
		}

		sess, err := store.Get(c)
		if err != nil {
			return err
		}

		sess.Set("username", user.Username)
		if err := sess.Save(); err != nil {
			return err
		}

		// Redirect to homepage
		return c.Redirect("/")
	})

	// Handle logout
	app.Get("/logout", func(c *fiber.Ctx) error {
		sess, err := store.Get(c)
		if err != nil {
			return err
		}

		sess.Destroy()
		if err := sess.Save(); err != nil {
			return err
		}

		return c.Redirect("/login")
	})

	// Start the server
	log.Fatal(app.Listen(":3000"))
}
