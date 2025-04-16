package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type StringArray []string

func (s StringArray) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *StringArray) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to convert value to []byte")
	}
	return json.Unmarshal(b, s)
}

type ChatFormRequest struct {
	UserName    string `json:"name" binding:"required" form:"user-name"`
	UserInfo    string `json:"phone_number" binding:"required" form:"user-info"`
	UserMessage string `json:"birth_date" binding:"required" form:"user-message"`
	ChatID      string `json:"chat_id" binding:"required" form:"chat-id"`
}

type Chat struct {
	ID          int         `json:"id" db:"id"`
	UserID      int         `json:"user_id" db:"user_id"`
	CharacterID int         `json:"character_id" form:"character-id" db:"character_id"`
	UserMessage StringArray `json:"user_message" form:"user_message" db:"user_message"`
	AiMessage   StringArray `json:"ai_message" form:"ai_message" db:"ai_message"`
	CreatedAt   time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at" db:"updated_at"`
}

func (Chat) TableName() string {
	return "chats"
}

type ChatSummary struct {
	ID        int    `db:"id" json:"id"`
	UserID    int    `db:"user_id" json:"user_id"`
	ChatID    int    `db:"chat_id" json:"chat_id"`
	Summary   string `db:"summary" json:"summary"`
	MessageID int    `db:"message_id" json:"message_id"`
}

func (ChatSummary) TableName() string {
	return "chat_summaries"
}
