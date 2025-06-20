package handlers

import (
	"net/http"
	"path/filepath"
	"strings"
	"taskflow/internal/service"

	"github.com/gin-gonic/gin"
)

type FileHandler struct {
	fileUseCase service.FileService
}

func NewFileHandler(fileUseCase service.FileService) *FileHandler {
	return &FileHandler{
		fileUseCase: fileUseCase,
	}
}

func (h *FileHandler) UploadFile(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file provided"})
		return
	}
	defer file.Close()

	// Validate file type
	ext := strings.ToLower(filepath.Ext(header.Filename))
	allowedTypes := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true, ".gif": true,
		".txt": true, ".pdf": true, ".doc": true, ".docx": true,
	}

	if !allowedTypes[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File type not allowed"})
		return
	}

	// Validate file size (10MB limit)
	if header.Size > 10*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File too large"})
		return
	}

	fileID, err := h.fileUseCase.UploadFile(
		c.Request.Context(),
		header.Filename,
		file,
		header.Header.Get("Content-Type"),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"fileId": fileID})
}
