package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"golang.org/x/crypto/bcrypt"

	"lgo/models"
	"lgo/db"
)

var Store = session.New()

func ShowLogin(c *fiber.Ctx) error {
	return c.Render("login", nil)
}

func Login(c *fiber.Ctx) error {
	email := c.FormValue("email")
	password := c.FormValue("password")

	var user models.User
	if err := db.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return c.Render("login", fiber.Map{"Error": "Invalid email or password"})
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return c.Render("login", fiber.Map{"Error": "Invalid email or password"})
	}

	sess, _ := Store.Get(c)
	sess.Set("username", user.Username)
	sess.Save()

	return c.Redirect("/")
}

func ShowSignup(c *fiber.Ctx) error {
	return c.Render("signup", nil)
}

func Signup(c *fiber.Ctx) error {
	username := c.FormValue("username")
	email := c.FormValue("email")
	password := c.FormValue("password")

	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := models.User{
		Username: username,
		Email:    email,
		Password: string(hashed),
	}

	if err := db.DB.Create(&user).Error; err != nil {
		return c.Render("signup", fiber.Map{"Error": "Username or email already taken"})
	}

	return c.Redirect("/login")
}

func Logout(c *fiber.Ctx) error {
	sess, _ := Store.Get(c)
	sess.Destroy()
	return c.Redirect("/login")
}
