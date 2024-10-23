package controllers

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/gofiber/fiber/v3"
	"io"
	"mime/multipart"
	"net/http"
	"regexp"
	"storage-api/src/domain"
	"storage-api/src/infrastructure/services"
	"strings"
	"time"
)

const DefaultContentType = "application/octet-stream"

type FileInfo struct {
	Filename     string    `json:"filename"`
	Folder       string    `json:"folder"`
	Size         int64     `json:"size"`
	LastModified time.Time `json:"lastModified"`
	Url          string    `json:"url"`
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

	return ctx.Status(http.StatusOK).JSON(result)
}

func (c *ICloudflareController) GetFilesHandler(ctx fiber.Ctx) error {
	result := domain.ResultData[[]FileInfo]()

	fullPath := ctx.Params("*")

	segments := strings.Split(fullPath, "/")
	filename := segments[len(segments)-1]
	if strings.Contains(filename, ".") {
		result.AddError(http.StatusBadRequest, "File name is not allowed")
		return ctx.Status(http.StatusBadRequest).JSON(result)
	}

	rawFiles, err := c.storage.GetFiles(fullPath)
	if err != nil {
		result.AddError(http.StatusNotFound, err.Error())
		return ctx.Status(http.StatusNotFound).JSON(result)
	}

	var exclude = []string{
		"(/\\.|^\\.)",
	}

	var excludeFolders = make([]string, 0)
	if len(domain.CONFIG.ExcludeFolders) > 0 {
		for _, folder := range domain.CONFIG.ExcludeFolders {
			if folder != "" {
				excludeFolders = append(excludeFolders, folder)
			}
		}
	}

	var excludeFiles = make([]string, 0)
	if len(domain.CONFIG.ExcludeFiles) > 0 {
		for _, file := range domain.CONFIG.ExcludeFiles {
			if file != "" {
				excludeFiles = append(excludeFiles, file)
			}
		}
	}

	files := make([]FileInfo, 0, len(rawFiles))
	for _, rawFile := range rawFiles {
		filePath := *rawFile.Key

		if regexp.MustCompile(`(?i)` + strings.Join(exclude, "|")).MatchString(filePath) {
			continue
		}

		if len(excludeFolders) > 0 && regexp.MustCompile(`(?i)`+strings.Join(excludeFolders, "|")).MatchString(filePath) {
			continue
		}
		if len(excludeFiles) > 0 && regexp.MustCompile(`(?i)`+strings.Join(excludeFiles, "|")).MatchString(filePath) {
			continue
		}

		fileName := filePath[strings.LastIndex(filePath, "/")+1:]
		if !strings.Contains(fileName, ".") {
			continue
		}

		path := fmt.Sprintf("%s/file/%s", domain.CONFIG.ApiUrl, filePath)

		files = append(files, FileInfo{
			Filename:     fileName,
			Folder:       filePath[:strings.LastIndex(filePath, "/")],
			Url:          path,
			Size:         *rawFile.Size,
			LastModified: *rawFile.LastModified,
		})
	}

	result.AddData(files)
	return ctx.Status(http.StatusOK).JSON(result)
}

func (c *ICloudflareController) GetFileHandler(ctx fiber.Ctx) error {
	result := domain.ResultData[FileInfo]()

	fullPath := ctx.Params("*")

	segments := strings.Split(fullPath, "/")
	filename := segments[len(segments)-1]
	if filename == "" {
		result.AddError(http.StatusBadRequest, "File name is missing")
		return ctx.Status(http.StatusBadRequest).JSON(result)
	}

	if !strings.Contains(filename, ".") {
		result.AddError(http.StatusBadRequest, "File name is not allowed")
		return ctx.Status(http.StatusBadRequest).JSON(result)
	}

	file, err := c.storage.GetFile(fullPath)
	if err != nil {
		result.AddError(http.StatusNotFound, err.Error())
		return ctx.Status(http.StatusNotFound).JSON(result)
	}

	if file == nil {
		result.AddError(http.StatusNotFound, "File not found")
		return ctx.Status(http.StatusNotFound).JSON(result)
	}

	if file.ContentType == nil {
		file.ContentType = aws.String(DefaultContentType)
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
		return ctx.Status(http.StatusBadRequest).JSON(result)
	}

	if !strings.Contains(filename, ".") {
		result.AddError(http.StatusBadRequest, "File name is not allowed")
		return ctx.Status(http.StatusBadRequest).JSON(result)
	}

	_, err := c.storage.GetFile(fullPath)
	if err != nil {
		result.AddError(http.StatusNotFound, err.Error())
		return ctx.Status(http.StatusNotFound).JSON(result)
	}

	_, errDelete := c.storage.DeleteFile(fullPath)
	if errDelete != nil {
		result.AddMessage("File could not be deleted")
		result.AddError(http.StatusInternalServerError, errDelete.Error())
		return ctx.Status(http.StatusInternalServerError).JSON(result)
	}

	result.AddMessage("File deleted successfully")

	return ctx.Status(http.StatusOK).JSON(result)
}

func (c *ICloudflareController) UploadFileHandler(ctx fiber.Ctx) error {
	result := domain.ResultData[[]FileInfo]()

	isOverwrite := ctx.Query("overwrite", "false")

	if !strings.Contains(ctx.Get("Content-Type"), "multipart/form-data") {
		result.AddError(http.StatusBadRequest, "Request is not a multipart/form-data")
		return ctx.Status(http.StatusBadRequest).JSON(result)
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		if err.Error() == "request Content-Type has bad boundary or is not multipart/form-data" {
			result.AddError(http.StatusBadRequest, "The request body is not a valid multipart/form-data")
		} else {
			result.AddError(http.StatusBadRequest, "Error retrieving form data: "+err.Error())
		}

		return ctx.Status(http.StatusBadRequest).JSON(result)
	}

	rawFolder := form.Value
	if len(rawFolder) < 0 || len(rawFolder["folder"]) < 0 || rawFolder["folder"][0] == "" {
		result.AddError(http.StatusBadRequest, "Folder is missing")
		return ctx.Status(http.StatusBadRequest).JSON(result)
	}

	folder := rawFolder["folder"][0]
	rawFiles := form.File["files"]
	if len(rawFiles) == 0 {
		result.AddError(http.StatusBadRequest, "File(s) is missing")
		return ctx.Status(http.StatusBadRequest).JSON(result)
	}

	var files []FileInfo
	for _, rawFile := range rawFiles {
		filename := rawFile.Filename
		contentType := rawFile.Header.Get("Content-Type")
		size := rawFile.Size

		path := fmt.Sprintf("/%s/%s", folder, filename)

		if isOverwrite == "false" {
			_, errFile := c.storage.GetFile(path)
			if errFile != nil {
				result.AddError(http.StatusConflict, "File already exists: "+rawFile.Filename)
				continue
			}
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

				domain.Logger.Error(errData.Error())
			}
		}(fileData)

		if contentType != DefaultContentType {
			contentType = DefaultContentType
		}

		_, errUpload := c.storage.UploadFile(fileData, folder, filename, contentType)
		if errUpload != nil {
			result.AddError(http.StatusBadRequest, "Error when uploading file: "+rawFile.Filename)

			domain.Logger.Error(errUpload.Error())

			continue
		}

		files = append(files, FileInfo{
			Filename:     filename,
			Folder:       folder,
			Size:         size,
			LastModified: time.Now(),
			Url:          domain.CONFIG.ApiUrl + "/file" + path,
		})
	}

	if len(files) > 0 {
		result.AddData(files)
		result.AddMessage(fmt.Sprintf("Files uploaded successfully: %d", len(files)))
		ctx.Status(http.StatusOK)
	} else {
		result.AddMessage("No files were uploaded successfully")
		ctx.Status(http.StatusBadRequest)
	}

	return ctx.JSON(result)
}
