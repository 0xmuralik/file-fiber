package main

import (
	"fmt"
	"net/http"

	"github.com/0xmuralik/file-share/database"
	"github.com/0xmuralik/file-share/file"
	"github.com/0xmuralik/file-share/user"
	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/jinzhu/gorm"
)

func routes(app *fiber.App) {

	app.Get("/home", user.Home)

	userAPI := app.Group("/user")
	userAPI.Post("/login", user.LogIn)
	userAPI.Post("/logout", user.LogOut)
	userAPI.Post("/register", user.Register)

	fileAPI := app.Group("/file")
	fileAPI.Get("/id/:file_id", file.GetFileById)
	fileAPI.Get("/name/:name", file.GetFileByName)
	fileAPI.Get("/:owner", file.GetFiles)
	fileAPI.Post("/new", file.NewFile)
	fileAPI.Delete("/:file_id", file.DeleteFile)
}

func initDatabase() {
	var err error
	database.FileDBConn, err = gorm.Open("sqlite3", "files.db")
	if err != nil {
		panic("failed to connect file database")
	}
	fmt.Println("connected to file database")

	database.FileDBConn.AutoMigrate(&file.File{})
	fmt.Println("file database migrated")

	database.UserDBConn, err = gorm.Open("sqlite3", "users.db")
	if err != nil {
		panic("failed to connect user database")
	}
	fmt.Println("connected to user database")

	database.UserDBConn.AutoMigrate(&user.User{})
	fmt.Println("user database migrated")

	user.Store = session.New()
	fmt.Println("Initialized session store")
}

func main() {
	app := fiber.New()

	app.Use(filesystem.New(filesystem.Config{
		Root:   http.Dir("./templates"),
		Browse: true,
	}))

	initDatabase()
	defer database.FileDBConn.Close()

	routes(app)
	if err := app.Listen(":3000"); err != nil {
		panic(err)
	}
}
