package main

import (
	"github.com/gin-gonic/gin"
	"manschko.com/cloud-storage/controllers"
)

// Create a file controller
func getFileController() *controllers.FileController {
	return controllers.NewFileController(
		AppConfig.SFTPHost,
		AppConfig.SFTPPort,
	)
}

// File operation handlers
func listFiles(c *gin.Context) {
	getFileController().ListFiles(c)
}

func downloadFile(c *gin.Context) {
	getFileController().DownloadFile(c)
}

func uploadFile(c *gin.Context) {
	getFileController().UploadFile(c)
}

func deleteFile(c *gin.Context) {
	getFileController().DeleteFile(c)
}

func moveFile(c *gin.Context) {
	getFileController().MoveFile(c)
}

func renameFile(c *gin.Context) {
	getFileController().RenameFile(c)
}

func createDirectory(c *gin.Context) {
	getFileController().CreateDirectory(c)
}