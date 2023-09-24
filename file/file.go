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
	res := db.Find(&files, "owner=?", owner)
	if res.Error != nil {
		return ctx.Status(503).SendString(fmt.Sprintf("Cannot find files: %s", res.Error.Error()))
	}

	ctx.JSON(files)

	return nil
}

func GetFileByName(ctx *fiber.Ctx) error {
	db := database.FileDBConn

	name := ctx.Params("name")

	var files []File
	res := db.Find(&files, "name=?", name)
	if res.Error != nil {
		return ctx.Status(503).SendString(fmt.Sprintf("Cannot find files: %s", res.Error.Error()))
	}
	ctx.JSON(files)

	return nil
}

func GetFileById(ctx *fiber.Ctx) error {
	db := database.FileDBConn

	id := ctx.Params("file_id")

	var file File
	res := db.Find(&file, "file_id=?", id)
	if res.Error != nil {
		return ctx.Status(503).SendString(fmt.Sprintf("Cannot find file: %s", res.Error.Error()))
	}

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
	res := db.Create(&file)
	if res.Error != nil {
		return ctx.Status(503).SendString(fmt.Sprintf("Cannot create file: %s", res.Error.Error()))
	}
	ctx.JSON(file)

	return nil
}

func DeleteFile(ctx *fiber.Ctx) error {
	db := database.FileDBConn

	id := ctx.Params("file_id")

	var file File
	res := db.First(&file, id)
	if res.Error != nil {
		return ctx.Status(503).SendString(fmt.Sprintf("Cannot find file: %s", res.Error.Error()))
	}

	if file.Name == "" {
		ctx.Status(500)
		return fmt.Errorf("No file found")
	}

	res = db.Delete(&file)
	if res.Error != nil {
		return ctx.Status(503).SendString(fmt.Sprintf("Cannot delete file: %s", res.Error.Error()))
	}

	ctx.SendString("File deleted successfully")

	return nil
}
