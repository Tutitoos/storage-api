package controllers

import (
	"fmt"
	"github.com/gofiber/fiber/v3"
	"io"
	"mime/multipart"
	"net/http"
	"storage-api/src/domain"
	"storage-api/src/infrastructure/services"
	"strings"
	"time"
)

type FileInfo struct {
	Id           string    `json:"id"`
	Url          string    `json:"url"`
	Size         int64     `json:"size"`
	LastModified time.Time `json:"lastModified"`
}

type ICloudflareController struct {
	storage *services.ICloudflareService
}

func CloudflareController() *ICloudflareController {
	return &ICloudflareController{
		storage: services.CloudflareService(),
	}
}

func (c *ICloudflareController) GetHomeHandler(ctx fiber.Ctx) error {
	result := domain.ResultData[string]()

	result.AddMessage("API is up and running!")

	return ctx.JSON(result)
}

func (c *ICloudflareController) GetFilesHandler(ctx fiber.Ctx) error {
	result := domain.ResultData[[]FileInfo]()

	fullPath := ctx.Params("*")

	segments := strings.Split(fullPath, "/")
	filename := segments[len(segments)-1]
	if strings.Contains(filename, ".") {
		result.AddError(http.StatusBadRequest, "File name is not allowed")
		return ctx.JSON(result)
	}

	rawFiles, err := c.storage.GetFiles(fullPath)
	if err != nil {
		result.AddError(http.StatusNotFound, err.Error())
		return ctx.JSON(result)
	}

	files := make([]FileInfo, 0, len(rawFiles))
	for _, rawFile := range rawFiles {
		filePath := *rawFile.Key
		fileId := filePath[strings.LastIndex(filePath, "/")+1:]
		path := fmt.Sprintf("%s/file/%s", domain.CONFIG.ApiUrl, filePath)

		files = append(files, FileInfo{
			Id:           fileId,
			Url:          path,
			Size:         *rawFile.Size,
			LastModified: *rawFile.LastModified,
		})
	}

	result.AddData(files)
	return ctx.JSON(result)
}

func (c *ICloudflareController) GetFileHandler(ctx fiber.Ctx) error {
	result := domain.ResultData[FileInfo]()

	fullPath := ctx.Params("*")

	segments := strings.Split(fullPath, "/")
	filename := segments[len(segments)-1]
	if filename == "" {
		result.AddError(http.StatusBadRequest, "File name is missing")
		return ctx.JSON(result)
	}

	if !strings.Contains(filename, ".") {
		result.AddError(http.StatusBadRequest, "File name is not allowed")
		return ctx.JSON(result)
	}

	file, err := c.storage.GetFile(fullPath)
	if err != nil {
		result.AddError(http.StatusNotFound, err.Error())
		return ctx.JSON(result)
	}

	ctx.Attachment(filename)
	ctx.Status(http.StatusOK)
	ctx.Set("Content-Type", *file.ContentType)
	return ctx.SendStream(io.NopCloser(file.Body))
}

func (c *ICloudflareController) DeleteFileHandler(ctx fiber.Ctx) error {
	result := domain.ResultData[string]()

	fullPath := ctx.Params("*")

	segments := strings.Split(fullPath, "/")
	filename := segments[len(segments)-1]
	if filename == "" {
		result.AddError(http.StatusBadRequest, "File name is missing")
		return ctx.JSON(result)
	}

	if !strings.Contains(filename, ".") {
		result.AddError(http.StatusBadRequest, "File name is not allowed")
		return ctx.JSON(result)
	}

	_, err := c.storage.GetFile(fullPath)
	if err != nil {
		result.AddError(http.StatusNotFound, err.Error())
		return ctx.JSON(result)
	}

	_, errDelete := c.storage.DeleteFile(fullPath)
	if errDelete != nil {
		result.AddMessage("File could not be deleted")
		result.AddError(http.StatusInternalServerError, errDelete.Error())
		return ctx.JSON(result)
	}

	result.AddMessage("File deleted successfully")

	return ctx.JSON(result)
}

func (c *ICloudflareController) UploadFileHandler(ctx fiber.Ctx) error {
	result := domain.ResultData[[]FileInfo]()

	query := ctx.Queries()
	folder := query["folder"]

	form, err := ctx.MultipartForm()
	if err != nil {
		result.AddError(http.StatusBadRequest, "Error retrieving form data: "+err.Error())
		return ctx.JSON(result)
	}

	rawFiles := form.File["files"]
	if len(rawFiles) == 0 {
		result.AddError(http.StatusBadRequest, "File(s) is missing")
		return ctx.JSON(result)
	}

	var files []FileInfo
	for _, rawFile := range rawFiles {
		filename := rawFile.Filename
		contentType := rawFile.Header.Get("Content-Type")
		size := rawFile.Size

		path := fmt.Sprintf("/%s/%s", folder, filename)

		_, errFile := c.storage.GetFile(path)
		if errFile != nil {
			result.AddError(http.StatusConflict, "File already exists: "+rawFile.Filename)
			continue
		}

		fileData, errFileData := rawFile.Open()
		if errFileData != nil {
			result.AddError(http.StatusBadRequest, "Invalid file: "+rawFile.Filename)
			continue
		}

		defer func(data multipart.File) {
			errData := data.Close()
			if errData != nil {
				result.AddError(http.StatusInternalServerError, "Error when closing file: "+rawFile.Filename)
			}
		}(fileData)

		_, errUpload := c.storage.UploadFile(fileData, folder, filename, contentType)
		if errUpload != nil {
			result.AddError(http.StatusBadRequest, "Error when uploading file: "+rawFile.Filename)
			continue
		}

		files = append(files, FileInfo{
			Id:           filename,
			Size:         size,
			LastModified: time.Now(),
			Url:          domain.CONFIG.ApiUrl + "/file" + path,
		})
	}

	if len(files) > 0 {
		result.AddData(files)
		result.AddMessage(fmt.Sprintf("Files uploaded successfully: %d", len(files)))
	} else {
		result.AddMessage("No files were uploaded successfully")
	}

	return ctx.JSON(result)
}
