package DTO

// ShareCreateRequest represents the request body for creating a share
type ShareCreateRequest struct {
	StorageID    int  `json:"storage_id" binding:"required"`
	TargetUserID int  `json:"target_user_id" binding:"required"`
	CanShare     bool `json:"can_share"`
}
