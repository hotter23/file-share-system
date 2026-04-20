package model

import (
	"time"
)

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"uniqueIndex;size:64;not null" json:"username"`
	Password  string    `gorm:"size:255;not null" json:"-"`
	Email     string    `gorm:"size:128" json:"email"`
	Role      string    `gorm:"type:enum('USER','ADMIN');default:'USER'" json:"role"`
	Status    string    `gorm:"type:enum('ACTIVE','DISABLED');default:'ACTIVE'" json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (User) TableName() string {
	return "user"
}

type File struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	FileUuid      string    `gorm:"uniqueIndex;size:64;not null" json:"fileUuid"`
	FileName      string    `gorm:"size:255;not null" json:"fileName"`
	FileSize      uint64    `gorm:"not null;default:0" json:"fileSize"`
	FileType      string    `gorm:"size:64" json:"fileType"`
	Md5           string    `gorm:"size:64" json:"md5"`
	CosKey        string    `gorm:"size:512" json:"cosKey"`
	BucketName    string    `gorm:"size:128" json:"bucketName"`
	StorageType   string    `gorm:"type:enum('COS','LOCAL');default:'COS'" json:"storageType"`
	FolderID      *uint     `json:"folderId"`
	UserID        uint      `gorm:"not null" json:"userId"`
	DownloadCount uint      `gorm:"not null;default:0" json:"downloadCount"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

func (File) TableName() string {
	return "file"
}

type Folder struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	FolderUuid string    `gorm:"uniqueIndex;size:64;not null" json:"folderUuid"`
	FolderName string    `gorm:"size:255;not null" json:"folderName"`
	ParentID   *uint     `json:"parentId"`
	UserID     uint      `gorm:"not null" json:"userId"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

func (Folder) TableName() string {
	return "folder"
}

type Share struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	ShareUuid  string     `gorm:"uniqueIndex;size:64;not null" json:"shareUuid"`
	FileID     uint       `gorm:"not null" json:"fileId"`
	ShareCode  string     `gorm:"size:16" json:"shareCode"`
	Password   string     `gorm:"size:255" json:"-"`
	ExpireTime *time.Time `json:"expireTime"`
	ViewCount  uint       `gorm:"not null;default:0" json:"viewCount"`
	CreatedBy  uint       `gorm:"not null" json:"createdBy"`
	CreatedAt  time.Time  `json:"createdAt"`
}

func (Share) TableName() string {
	return "share"
}

type FileChunk struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	FileUuid   string    `gorm:"uniqueIndex:uk_file_chunk;size:64;not null" json:"fileUuid"`
	ChunkIndex uint      `gorm:"uniqueIndex:uk_file_chunk;not null" json:"chunkIndex"`
	ChunkSize  uint64    `gorm:"not null;default:0" json:"chunkSize"`
	UploadID   string    `gorm:"size:255;not null" json:"uploadId"`
	Status     uint      `gorm:"not null;default:0" json:"status"`
	CreatedAt  time.Time `json:"createdAt"`
}

func (FileChunk) TableName() string {
	return "file_chunk"
}
