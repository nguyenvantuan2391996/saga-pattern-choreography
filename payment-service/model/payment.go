package model

type Payment struct {
	ID      int32 `json:"id"`
	UserID  int32 `json:"user_id"`
	Balance int32 `json:"balance"`
}
