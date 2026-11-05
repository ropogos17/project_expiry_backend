package api

type createLogRequest struct {
	ActionTypeID int32  `json:"action_type_id"`
	Message      string `json:"message"`
	Data         struct{}
	UserID       int32 `json:"user_id"`
}
