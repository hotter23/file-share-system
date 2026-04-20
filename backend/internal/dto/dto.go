package dto

type UserRegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=64"`
	Password string `json:"password" binding:"required,min=6,max=128"`
	Email    string `json:"email"`
}

type UserLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type InitChunkUploadRequest struct {
	FileName  string `json:"fileName" binding:"required"`
	FileSize  int64  `json:"fileSize" binding:"required"`
	ChunkSize int64  `json:"chunkSize" binding:"required"`
	FolderID  *uint  `json:"folderId"`
}

type CreateShareRequest struct {
	FileUuid   string `json:"fileUuid" binding:"required"`
	Password   string `json:"password"`
	ExpireType string `json:"expireType"` // 1H, 24H, 7D, PERMANENT
}

type CreateFolderRequest struct {
	FolderName string `json:"folderName" binding:"required"`
	ParentID   *uint  `json:"parentId"`
}
