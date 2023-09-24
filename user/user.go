package user

import (
	"fmt"

	"github.com/0xmuralik/file-share/database"
	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

var Store *session.Store

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Home(ctx *fiber.Ctx) error {
	sess, err := Store.Get(ctx)
	if err != nil {
		return ctx.Status(500).SendString(fmt.Sprintf("cannot get session: %v", err))
	}

	if auth := sess.Get("authenticated"); auth == nil || !auth.(bool) {
		return ctx.Status(503).SendString("This is home, User is not logged in")
	}

	username := sess.Get("user")
	return ctx.SendString(fmt.Sprintf("This is home, welcome %s", username))
}

func Register(ctx *fiber.Ctx) error {
	db := database.UserDBConn

	user := new(User)
	if err := ctx.BodyParser(user); err != nil {
		return ctx.Status(503).SendString(fmt.Sprintf("cannot register user, error: %s", err.Error()))
	}

	res := db.Create(&user)
	if res.Error != nil {
		return ctx.Status(503).SendString(fmt.Sprintf("cannot create user, %s", res.Error.Error()))
	}
	ctx.JSON(&user)
	return ctx.Status(200).SendString("User created successfully")
}

func Delete(ctx *fiber.Ctx) error {
	db := database.UserDBConn

	sess, err := Store.Get(ctx)
	if err != nil {
		return ctx.Status(500).SendString(fmt.Sprintf("cannot get session: %v", err))
	}

	if auth := sess.Get("authenticated"); !auth.(bool) {
		return ctx.Status(503).SendString("User is not logged in")
	}

	username := sess.Get("user")

	var user User
	res := db.First(&user, "username=?", username)
	if res.Error != nil {
		return ctx.Status(503).SendString(fmt.Sprintf("cannot find user: %s", res.Error.Error()))
	}

	if user.Username == "" {
		return ctx.Status(500).SendString("No file found")
	}

	res = db.Delete(&user)
	if res.Error != nil {
		return ctx.Status(503).SendString(fmt.Sprintf("Cannot delete user: %s", res.Error.Error()))
	}

	ctx.Status(200).SendString("User deleted successfully")
	return ctx.Redirect("/home")
}

func LogIn(ctx *fiber.Ctx) error {

	db := database.UserDBConn
	sess, err := Store.Get(ctx)
	if err != nil {
		return ctx.Status(500).SendString(fmt.Sprintf("cannot get session: %v", err))
	}
	req := new(User)
	if err := ctx.BodyParser(req); err != nil {
		return ctx.Status(503).SendString(fmt.Sprintf("invalid req: %s", err.Error()))
	}

	var user User
	res := db.Find(&user, "username=?", req.Username)
	if res.Error != nil {
		return ctx.Status(503).SendString(fmt.Sprintf("cannot find user: %s", res.Error.Error()))
	}

	if user.Username != "" && user.Password == req.Password {
		sess.Set("authenticated", true)
		sess.Set("user", user.Username)
		sess.Save()
		ctx.SendString(fmt.Sprintf("User logged in successfully. Username: %s", user.Username))
		return ctx.Redirect("/api/v1/user/home")
	}

	return ctx.Status(503).SendString("Invalid username or password")
}

func LogOut(ctx *fiber.Ctx) error {
	sess, err := Store.Get(ctx)
	if err != nil {
		return ctx.Status(500).SendString(fmt.Sprintf("cannot get session: %v", err))
	}

	sess.Delete("authenticated")
	sess.Save()
	ctx.SendString("User logged out")

	return ctx.Redirect("/api/v1/user/home")
}
