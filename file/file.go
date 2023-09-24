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
		return ctx.Status(503).SendString(fmt.Sprintf("cannot find files, %s", res.Error.Error()))
	}

	return ctx.JSON(files)
}

func GetFileByName(ctx *fiber.Ctx) error {
	db := database.FileDBConn

	name := ctx.Params("name")

	var files []File
	res := db.Find(&files, "name=?", name)
	if res.Error != nil {
		return ctx.Status(503).SendString(fmt.Sprintf("cannot find files, %s", res.Error.Error()))
	}

	return ctx.JSON(files)
}

func GetFileById(ctx *fiber.Ctx) error {
	db := database.FileDBConn

	id := ctx.Params("file_id")

	var file File
	res := db.Find(&file, "file_id=?", id)
	if res.Error != nil {
		return ctx.Status(503).SendString(fmt.Sprintf("cannot find file, %s", res.Error.Error()))
	}

	return ctx.JSON(file)
}

func NewFile(ctx *fiber.Ctx) error {
	db := database.FileDBConn

	file := new(File)
	if err := ctx.BodyParser(file); err != nil {
		return ctx.Status(503).SendString(fmt.Sprintf("cannot parse file: %s", err))
	}

	res := db.Create(&file)
	if res.Error != nil {
		return ctx.Status(503).SendString(fmt.Sprintf("cannot create file: %s", res.Error.Error()))
	}

	return ctx.JSON(file)
}

func DeleteFile(ctx *fiber.Ctx) error {
	db := database.FileDBConn

	id := ctx.Params("file_id")

	var file File
	res := db.First(&file, "file_id=?", id)
	if res.Error != nil {
		return ctx.Status(503).SendString(fmt.Sprintf("cannot find file, %s", res.Error.Error()))
	}

	res = db.Delete(&file)
	if res.Error != nil {
		return ctx.Status(503).SendString(fmt.Sprintf("cannot delete file, %s", res.Error.Error()))
	}

	return ctx.SendString("File deleted successfully")
}
