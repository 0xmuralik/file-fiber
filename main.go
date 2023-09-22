package main

import (
	"fmt"

	"github.com/0xmuralik/file-share/database"
	"github.com/0xmuralik/file-share/file"
	"github.com/gofiber/fiber"
	"github.com/jinzhu/gorm"
)

func routes(app *fiber.App) {
	app.Get("/api/v1/file/id/:file_id", file.GetFileById)
	app.Get("/api/v1/file/name/:name", file.GetFileByName)
	app.Get("/api/v1/file/:owner", file.GetFiles)
	app.Post("/api/v1/file", file.NewFile)
	app.Delete("/api/v1/file/:id", file.DeleteFile)
}

func initDatabase() {
	var err error
	database.DBConn, err = gorm.Open("sqlite3", "files.db")
	if err != nil {
		panic("failed to connect database")
	}
	fmt.Println("connected to database")

	database.DBConn.AutoMigrate(&file.File{})
	fmt.Println("database migrated")
}

func main() {
	app := fiber.New()

	initDatabase()
	defer database.DBConn.Close()

	routes(app)
	app.Listen(3000)
}
