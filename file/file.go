package file

import (
	"github.com/0xmuralik/file-share/database"
	"github.com/gofiber/fiber"
	"github.com/jinzhu/gorm"
)

type File struct {
	gorm.Model
	FileId int    `json:"file_id"`
	Name   string `json:"name"`
	Owner  string `json:"owner"`
	// Data   []byte `json:"data"`
}

func GetFiles(ctx *fiber.Ctx) {
	db := database.DBConn
	var files []File

	owner := ctx.Params("owner")
	db.Find(&files, "owner=?", owner)
	ctx.JSON(files)
}
func GetFileByName(ctx *fiber.Ctx) {
	db := database.DBConn

	name := ctx.Params("name")

	var files []File
	db.Find(&files, "name=?", name)
	ctx.JSON(files)
}

func GetFileById(ctx *fiber.Ctx) {
	db := database.DBConn

	id := ctx.Params("file_id")

	var file File
	db.Find(&file, "file_id=?", id)
	ctx.JSON(file)
}

func NewFile(ctx *fiber.Ctx) {
	db := database.DBConn

	file := new(File)
	if err := ctx.BodyParser(file); err != nil {
		ctx.Status(503).Send(err)
		return
	}
	db.Create(&file)
	ctx.JSON(file)
}

func DeleteFile(ctx *fiber.Ctx) {
	db := database.DBConn

	id := ctx.Params("file_id")

	var file File
	db.First(&file, id)
	if file.Name == "" {
		ctx.Status(500).Send("No file found")
		return
	}
	db.Delete(&file)
	ctx.Send("File deleted successfully")
}
