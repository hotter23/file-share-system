package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/fileshare/internal/dto"
	"github.com/fileshare/internal/model"
	"github.com/fileshare/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func main() {
	dsn := "root:root123@tcp(localhost:3306)/file_share?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("连接数据库失败:", err)
	}

	r := gin.Default()

	// CORS 中间件
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-User-Id, X-Username, X-Role")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	// 路由
	api := r.Group("/api/file")
	api.Use(authMiddleware())
	{
		api.POST("/upload", uploadFile)
		api.POST("/upload/init", initChunkUpload)
		api.POST("/upload/chunk", uploadChunk)
		api.POST("/upload/complete", completeChunkUpload)
		api.GET("/list", getFileList)
		api.GET("/:fileUuid", getFileDetail)
		api.GET("/download/:fileUuid", downloadFile)
		api.DELETE("/:fileUuid", deleteFile)

		api.POST("/folder", createFolder)
		api.GET("/folder/list", getFolderList)
		api.DELETE("/folder/:folderUuid", deleteFolder)
	}

	log.Println("文件服务启动在 :8082")
	r.Run(":8082")
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.GetHeader("X-User-Id")
		if userId == "" {
			response.Unauthorized(c, "未认证")
			c.Abort()
			return
		}
		c.Set("userId", userId)
		c.Next()
	}
}

func getUserId(c *gin.Context) uint {
	userIdStr, _ := c.Get("userId")
	id, _ := strconv.ParseUint(userIdStr.(string), 10, 64)
	return uint(id)
}

// ========== 文件上传 ==========

func uploadFile(c *gin.Context) {
	userId := getUserId(c)

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		response.BadRequest(c, "请选择文件")
		return
	}
	defer file.Close()

	folderIdStr := c.PostForm("folderId")
	var folderId *uint
	if folderIdStr != "" {
		id, _ := strconv.ParseUint(folderIdStr, 10, 64)
		folderId = new(uint)
		*folderId = uint(id)
	}

	// 读取文件内容计算 MD5
	content, _ := io.ReadAll(file)
	md5Hash := md5.Sum(content)
	md5Str := hex.EncodeToString(md5Hash[:])

	// 生成文件 UUID
	fileUuid := uuid.New().String()

	// 存储路径
	datePath := time.Now().Format("2006/01")
	cosKey := fmt.Sprintf("users/%d/%s/%s/%s", userId, datePath, fileUuid, header.Filename)

	// 本地存储（简化版）
	localPath := fmt.Sprintf("/tmp/file-share/%d/%s", userId, fileUuid)
	os.MkdirAll(localPath, 0755)
	dstPath := filepath.Join(localPath, header.Filename)
	os.WriteFile(dstPath, content, 0644)

	// 保存到数据库
	fileRecord := model.File{
		FileUuid:    fileUuid,
		FileName:    header.Filename,
		FileSize:    uint64(header.Size),
		FileType:    header.Header.Get("Content-Type"),
		Md5:         md5Str,
		CosKey:      cosKey,
		BucketName:  "file-share-bucket",
		StorageType: "COS",
		FolderID:    folderId,
		UserID:      userId,
	}

	if err := db.Create(&fileRecord).Error; err != nil {
		response.ServerError(c, "保存文件记录失败")
		return
	}

	response.Success(c, gin.H{
		"fileUuid": fileUuid,
		"fileName": header.Filename,
		"fileSize": header.Size,
		"md5":      md5Str,
	})
}

func initChunkUpload(c *gin.Context) {
	userId := getUserId(c)

	var req dto.InitChunkUploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	fileUuid := uuid.New().String()
	chunkCount := int64(req.FileSize / req.ChunkSize)
	if req.FileSize%req.ChunkSize != 0 {
		chunkCount++
	}

	// 构建 COS 路径
	datePath := time.Now().Format("2006/01")
	cosKey := fmt.Sprintf("users/%d/%s/%s/%s", userId, datePath, fileUuid, req.FileName)

	// 保存文件记录
	fileRecord := model.File{
		FileUuid:    fileUuid,
		FileName:    req.FileName,
		FileSize:    uint64(req.FileSize),
		CosKey:      cosKey,
		BucketName:  "file-share-bucket",
		StorageType: "COS",
		FolderID:    req.FolderID,
		UserID:      userId,
	}

	if err := db.Create(&fileRecord).Error; err != nil {
		response.ServerError(c, "保存文件记录失败")
		return
	}

	// 生成分片上传 ID（这里简化处理）
	uploadId := uuid.New().String()

	// 创建分片记录
	for i := int64(0); i < chunkCount; i++ {
		db.Create(&model.FileChunk{
			FileUuid:   fileUuid,
			ChunkIndex: uint(i),
			UploadID:   uploadId,
			Status:     0,
		})
	}

	response.Success(c, gin.H{
		"fileUuid":   fileUuid,
		"uploadId":   uploadId,
		"chunkCount": chunkCount,
	})
}

func uploadChunk(c *gin.Context) {
	fileUuid := c.PostForm("fileUuid")
	uploadId := c.PostForm("uploadId")
	chunkIndexStr := c.PostForm("chunkIndex")

	chunkIndex, _ := strconv.ParseUint(chunkIndexStr, 10, 64)

	// 检查分片是否已上传
	var chunk model.FileChunk
	result := db.Where("file_uuid = ? AND chunk_index = ?", fileUuid, chunkIndex).First(&chunk)
	if result.Error == nil && chunk.Status == 1 {
		response.Success(c, gin.H{
			"chunkIndex": chunkIndex,
			"uploadId":   uploadId,
		})
		return
	}

	// 更新分片状态
	db.Model(&model.FileChunk{}).Where("file_uuid = ? AND chunk_index = ?", fileUuid, chunkIndex).
		Updates(map[string]interface{}{"status": 1, "upload_id": uploadId})

	response.Success(c, gin.H{
		"chunkIndex": chunkIndex,
		"uploadId":   uploadId,
	})
}

func completeChunkUpload(c *gin.Context) {
	var req struct {
		FileUuid string `json:"fileUuid" binding:"required"`
		UploadId string `json:"uploadId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 标记所有分片完成（这里简化处理）
	db.Model(&model.FileChunk{}).Where("file_uuid = ?", req.FileUuid).Update("status", 1)

	// 获取文件信息
	var file model.File
	if err := db.Where("file_uuid = ?", req.FileUuid).First(&file).Error; err != nil {
		response.NotFound(c, "文件不存在")
		return
	}

	response.Success(c, gin.H{
		"fileUuid": req.FileUuid,
		"fileName": file.FileName,
		"md5":      file.Md5,
	})
}

// ========== 文件管理 ==========

func getFileList(c *gin.Context) {
	userId := getUserId(c)

	folderIdStr := c.Query("folderId")
	keyword := c.Query("keyword")
	pageStr := c.DefaultQuery("page", "1")
	sizeStr := c.DefaultQuery("size", "20")

	page, _ := strconv.Atoi(pageStr)
	size, _ := strconv.Atoi(sizeStr)
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	query := db.Model(&model.File{}).Where("user_id = ?", userId)

	if folderIdStr != "" {
		folderId, _ := strconv.ParseUint(folderIdStr, 10, 64)
		query = query.Where("folder_id = ?", folderId)
	} else {
		query = query.Where("folder_id IS NULL")
	}

	if keyword != "" {
		query = query.Where("file_name LIKE ?", "%"+keyword+"%")
	}

	var total int64
	query.Count(&total)

	var files []model.File
	offset := (page - 1) * size
	if err := query.Order("created_at DESC").Offset(offset).Limit(size).Find(&files).Error; err != nil {
		response.ServerError(c, "查询失败")
		return
	}

	records := make([]gin.H, len(files))
	for i, f := range files {
		records[i] = gin.H{
			"fileUuid":   f.FileUuid,
			"fileName":   f.FileName,
			"fileSize":   f.FileSize,
			"fileType":   f.FileType,
			"createdAt":  f.CreatedAt,
		}
	}

	response.Success(c, gin.H{
		"records": records,
		"total":   total,
		"page":    page,
		"size":    size,
	})
}

func getFileDetail(c *gin.Context) {
	userId := getUserId(c)
	fileUuid := c.Param("fileUuid")

	var file model.File
	if err := db.Where("file_uuid = ? AND user_id = ?", fileUuid, userId).First(&file).Error; err != nil {
		response.NotFound(c, "文件不存在")
		return
	}

	response.Success(c, gin.H{
		"fileUuid":      file.FileUuid,
		"fileName":      file.FileName,
		"fileSize":      file.FileSize,
		"fileType":      file.FileType,
		"md5":           file.Md5,
		"downloadCount": file.DownloadCount,
		"createdAt":     file.CreatedAt,
	})
}

func downloadFile(c *gin.Context) {
	userId := getUserId(c)
	fileUuid := c.Param("fileUuid")

	var file model.File
	if err := db.Where("file_uuid = ? AND user_id = ?", fileUuid, userId).First(&file).Error; err != nil {
		response.NotFound(c, "文件不存在")
		return
	}

	// 增加下载次数
	db.Model(&file).Update("download_count", file.DownloadCount+1)

	// 构建本地文件路径
	localPath := fmt.Sprintf("/tmp/file-share/%d/%s/%s", userId, fileUuid, file.FileName)
	c.File(localPath)
}

func deleteFile(c *gin.Context) {
	userId := getUserId(c)
	fileUuid := c.Param("fileUuid")

	var file model.File
	if err := db.Where("file_uuid = ? AND user_id = ?", fileUuid, userId).First(&file).Error; err != nil {
		response.NotFound(c, "文件不存在")
		return
	}

	db.Delete(&file)
	response.SuccessWithMessage(c, "删除成功", nil)
}

// ========== 文件夹管理 ==========

func createFolder(c *gin.Context) {
	userId := getUserId(c)

	var req dto.CreateFolderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	folderUuid := uuid.New().String()

	folder := model.Folder{
		FolderUuid: folderUuid,
		FolderName: req.FolderName,
		ParentID:   req.ParentID,
		UserID:     userId,
	}

	if err := db.Create(&folder).Error; err != nil {
		response.ServerError(c, "创建文件夹失败")
		return
	}

	response.Success(c, gin.H{
		"folderUuid":  folderUuid,
		"folderName":  req.FolderName,
	})
}

func getFolderList(c *gin.Context) {
	userId := getUserId(c)
	parentIdStr := c.Query("parentId")

	query := db.Model(&model.Folder{}).Where("user_id = ?", userId)

	if parentIdStr != "" {
		parentId, _ := strconv.ParseUint(parentIdStr, 10, 64)
		query = query.Where("parent_id = ?", parentId)
	} else {
		query = query.Where("parent_id IS NULL")
	}

	var folders []model.Folder
	if err := query.Order("folder_name ASC").Find(&folders).Error; err != nil {
		response.ServerError(c, "查询失败")
		return
	}

	records := make([]gin.H, len(folders))
	for i, f := range folders {
		records[i] = gin.H{
			"folderUuid":  f.FolderUuid,
			"folderName":  f.FolderName,
			"parentId":    f.ParentID,
			"createdAt":   f.CreatedAt,
		}
	}

	response.Success(c, records)
}

func deleteFolder(c *gin.Context) {
	userId := getUserId(c)
	folderUuid := c.Param("folderUuid")

	var folder model.Folder
	if err := db.Where("folder_uuid = ? AND user_id = ?", folderUuid, userId).First(&folder).Error; err != nil {
		response.NotFound(c, "文件夹不存在")
		return
	}

	db.Delete(&folder)
	response.SuccessWithMessage(c, "删除成功", nil)
}
