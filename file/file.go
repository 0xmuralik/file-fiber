package file

import (
	"fmt"

	"github.com/0xmuralik/file-share/database"
	fiber "github.com/gofiber/fiber/v2"
	"github.com/jinzhu/gorm"
)

type File struct {
	gorm.Model
	FileId int    `json:"file_id"`
	Name   string `json:"name"`
	Owner  string `json:"owner"`
	// Data   []byte `json:"data"`
}

func GetFiles(ctx *fiber.Ctx) error {
	db := database.FileDBConn
	var files []File

	owner := ctx.Params("owner")
	db.Find(&files, "owner=?", owner)
	ctx.JSON(files)

	return nil
}

func GetFileByName(ctx *fiber.Ctx) error {
	db := database.FileDBConn

	name := ctx.Params("name")

	var files []File
	db.Find(&files, "name=?", name)
	ctx.JSON(files)

	return nil
}

func GetFileById(ctx *fiber.Ctx) error {
	db := database.FileDBConn

	id := ctx.Params("file_id")

	var file File
	db.Find(&file, "file_id=?", id)
	ctx.JSON(file)

	return nil
}

func NewFile(ctx *fiber.Ctx) error {
	db := database.FileDBConn

	file := new(File)
	if err := ctx.BodyParser(file); err != nil {
		ctx.Status(503)
		return err
	}
	db.Create(&file)
	ctx.JSON(file)

	return nil
}

func DeleteFile(ctx *fiber.Ctx) error {
	db := database.FileDBConn

	id := ctx.Params("file_id")

	var file File
	db.First(&file, id)
	if file.Name == "" {
		ctx.Status(500)
		return fmt.Errorf("No file found")
	}
	db.Delete(&file)
	ctx.SendString("File deleted successfully")

	return nil
}
