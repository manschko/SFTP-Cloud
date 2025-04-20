package controllers

import (
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"manschko.com/cloud-storage/sftp"
)

// FileInfo represents file metadata
type FileInfo struct {
	Name         string `json:"name"`
	Size         int64  `json:"size"`
	IsDir        bool   `json:"is_dir"`
	ModTime      string `json:"mod_time"`
	Path         string `json:"path"`
	ThumbnailURL string `json:"thumbnail_url,omitempty"`
}

// MoveRequest represents a file move operation
type MoveRequest struct {
	Source      string `json:"source" binding:"required"`
	Destination string `json:"destination" binding:"required"`
}

// RenameRequest represents a file rename operation
type RenameRequest struct {
	Path    string `json:"path" binding:"required"`
	NewName string `json:"new_name" binding:"required"`
}

// FileController handles file operations
type FileController struct {
	SFTPHost string
	SFTPPort int
}

// NewFileController creates a new file controller
func NewFileController(sftpHost string, sftpPort int) *FileController {
	return &FileController{
		SFTPHost: sftpHost,
		SFTPPort: sftpPort,
	}
}

// getSFTPConnection creates an SFTP connection with the user's credentials
func (c *FileController) getSFTPConnection(ctx *gin.Context) (*sftp.Connection, error) {
	username, exists := ctx.Get("username")
	if !exists {
		return nil, fmt.Errorf("username not found in context")
	}

	// In a real application, you would retrieve the password from a secure store
	// For this example, we're assuming the password is available somehow
	// This is a placeholder - you should implement a proper credential management system
	password := "placeholder" // This would need to be securely stored/retrieved

	return sftp.NewConnection(c.SFTPHost, c.SFTPPort, username.(string), password), nil
}

// ListFiles lists files in the specified directory
func (c *FileController) ListFiles(ctx *gin.Context) {
	path := ctx.Param("path")
	if path == "" {
		path = "/"
	}

	// Get username from context (set by auth middleware)
	username, exists := ctx.Get("username")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "User not authenticated"})
		return
	}

	// In a real application, you would retrieve the password from a secure store
	// This is a placeholder
	password := "placeholder" // This would need to be securely stored/retrieved

	conn := sftp.NewConnection(c.SFTPHost, c.SFTPPort, username.(string), password)
	client, err := conn.Connect()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("SFTP connection error: %v", err)})
		return
	}
	defer client.Close()

	files, err := client.ReadDir(path)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to read directory: %v", err)})
		return
	}

	var fileInfos []FileInfo
	for _, file := range files {
		fileInfos = append(fileInfos, FileInfo{
			Name:    file.Name(),
			Size:    file.Size(),
			IsDir:   file.IsDir(),
			ModTime: file.ModTime().Format(time.RFC3339),
			Path:    filepath.Join(path, file.Name()),
		})
	}

	ctx.JSON(http.StatusOK, fileInfos)
}

// DownloadFile downloads a file from SFTP server
func (c *FileController) DownloadFile(ctx *gin.Context) {
	path := ctx.Param("path")

	// Get username from context
	username, exists := ctx.Get("username")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "User not authenticated"})
		return
	}

	// In a real application, you would retrieve the password from a secure store
	password := "placeholder" // This would need to be securely stored/retrieved

	conn := sftp.NewConnection(c.SFTPHost, c.SFTPPort, username.(string), password)
	client, err := conn.Connect()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("SFTP connection error: %v", err)})
		return
	}
	defer client.Close()

	file, err := client.Open(path)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to open file: %v", err)})
		return
	}
	defer file.Close()

	fileInfo, err := client.Stat(path)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to get file info: %v", err)})
		return
	}

	ctx.Header("Content-Description", "File Transfer")
	ctx.Header("Content-Transfer-Encoding", "binary")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filepath.Base(path)))
	ctx.Header("Content-Type", "application/octet-stream")
	ctx.Header("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))

	ctx.Status(http.StatusOK)
	io.Copy(ctx.Writer, file)
}

// UploadFile uploads a file to the SFTP server
func (c *FileController) UploadFile(ctx *gin.Context) {
	path := ctx.Param("path")

	// Parse multipart form
	err := ctx.Request.ParseMultipartForm(32 << 20) // 32MB max
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form"})
		return
	}

	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get file from form"})
		return
	}
	defer file.Close()

	// Get username from context
	username, exists := ctx.Get("username")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "User not authenticated"})
		return
	}

	// In a real application, you would retrieve the password from a secure store
	password := "placeholder" // This would need to be securely stored/retrieved

	conn := sftp.NewConnection(c.SFTPHost, c.SFTPPort, username.(string), password)
	client, err := conn.Connect()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("SFTP connection error: %v", err)})
		return
	}
	defer client.Close()

	// Create destination path
	destPath := filepath.Join(path, header.Filename)

	// Create remote file
	dstFile, err := client.Create(destPath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create remote file: %v", err)})
		return
	}
	defer dstFile.Close()

	// Copy file content
	_, err = io.Copy(dstFile, file)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to copy file content: %v", err)})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully", "path": destPath})
}

// DeleteFile deletes a file or directory from the SFTP server
func (c *FileController) DeleteFile(ctx *gin.Context) {
	path := ctx.Param("path")

	// Get username from context
	username, exists := ctx.Get("username")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "User not authenticated"})
		return
	}

	// In a real application, you would retrieve the password from a secure store
	password := "placeholder" // This would need to be securely stored/retrieved

	conn := sftp.NewConnection(c.SFTPHost, c.SFTPPort, username.(string), password)
	client, err := conn.Connect()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("SFTP connection error: %v", err)})
		return
	}
	defer client.Close()

	// Check if it's a directory
	fileInfo, err := client.Stat(path)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to get file info: %v", err)})
		return
	}

	if fileInfo.IsDir() {
		// Remove directory recursively
		//err = c.removeDirectory(client, path)
	} else {
		// Remove file
		err = client.Remove(path)
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to delete: %v", err)})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Deleted successfully"})
}

// removeDirectory recursively removes a directory and its contents
/*func (c *FileController) removeDirectory(client *sftp.Client, path string) error {
	// List all items in directory
	files, err := client.ReadDir(path)
	if err != nil {
		return err
	}

	// Delete each item
	for _, file := range files {
		itemPath := filepath.Join(path, file.Name())
		if file.IsDir() {
			if err := c.removeDirectory(client, itemPath); err != nil {
				return err
			}
		} else {
			if err := client.Remove(itemPath); err != nil {
				return err
			}
		}
	}

	// Delete the now-empty directory
	return client.RemoveDirectory(path)
}*/

// MoveFile moves a file or directory to a new location
func (c *FileController) MoveFile(ctx *gin.Context) {
	var moveReq MoveRequest
	if err := ctx.ShouldBindJSON(&moveReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Get username from context
	username, exists := ctx.Get("username")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "User not authenticated"})
		return
	}

	// In a real application, you would retrieve the password from a secure store
	password := "placeholder" // This would need to be securely stored/retrieved

	conn := sftp.NewConnection(c.SFTPHost, c.SFTPPort, username.(string), password)
	client, err := conn.Connect()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("SFTP connection error: %v", err)})
		return
	}
	defer client.Close()

	err = client.Rename(moveReq.Source, moveReq.Destination)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to move file: %v", err)})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "File moved successfully"})
}

// RenameFile renames a file or directory
func (c *FileController) RenameFile(ctx *gin.Context) {
	var renameReq RenameRequest
	if err := ctx.ShouldBindJSON(&renameReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Get username from context
	username, exists := ctx.Get("username")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "User not authenticated"})
		return
	}

	// In a real application, you would retrieve the password from a secure store
	password := "placeholder" // This would need to be securely stored/retrieved

	conn := sftp.NewConnection(c.SFTPHost, c.SFTPPort, username.(string), password)
	client, err := conn.Connect()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("SFTP connection error: %v", err)})
		return
	}
	defer client.Close()

	dir := filepath.Dir(renameReq.Path)
	newPath := filepath.Join(dir, renameReq.NewName)

	err = client.Rename(renameReq.Path, newPath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to rename file: %v", err)})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "File renamed successfully", "new_path": newPath})
}

// CreateDirectory creates a new directory on the SFTP server
func (c *FileController) CreateDirectory(ctx *gin.Context) {
	path := ctx.Param("path")

	// Get username from context
	username, exists := ctx.Get("username")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "User not authenticated"})
		return
	}

	// In a real application, you would retrieve the password from a secure store
	password := "placeholder" // This would need to be securely stored/retrieved

	conn := sftp.NewConnection(c.SFTPHost, c.SFTPPort, username.(string), password)
	client, err := conn.Connect()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("SFTP connection error: %v", err)})
		return
	}
	defer client.Close()

	err = client.MkdirAll(path)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create directory: %v", err)})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Directory created successfully"})
}
