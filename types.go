package claudeapi

import (
	"errors"
	"fmt"
	"time"
)

type Settings struct {
	ClaudeConsolePrivacy string `json:"claude_console_privacy"`
}

type Organization struct {
	Uuid         string    `json:"uuid"`
	Name         string    `json:"name"`
	Settings     Settings  `json:"settings"`
	Capabilities []string  `json:"capabilities"`
	ActiveFlags  []string  `json:"active_flags"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Organizations []*Organization

type ConversationInfo struct {
	Uuid      string    `json:"uuid"`
	Name      string    `json:"name"`
	Summary   string    `json:"summary"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type ListConversations []*ConversationInfo

type Document struct {
	FileName         string `json:"file_name"`
	FileSize         int    `json:"file_size"`
	FileType         string `json:"file_type"`
	ExtractedContent string `json:"extracted_content"`
	TotalPages       int    `json:"totalPages"`
}

type ErrorRes struct {
	Error struct {
		Type    string `json:"type"`
		Message string `json:"message"`
		Code    string `json:"code"`
	} `json:"error"`
}

func (r *ErrorRes) String() string {
	return fmt.Sprintf("code: %v, message: %v, type: %v", r.Error.Code, r.Error.Message, r.Error.Type)
}

func (r *ErrorRes) Err() error {
	return errors.New(fmt.Sprintf("code: %v, message: %v, type: %v", r.Error.Code, r.Error.Message, r.Error.Type))
}

type Completion struct {
	Prompt   string `json:"prompt"`
	Timezone string `json:"timezone"`
	Model    string `json:"model"`
}

type Message struct {
	Completion       Completion  `json:"completion"`
	OrganizationUuid string      `json:"organization_uuid"`
	ConversationUuid string      `json:"conversation_uuid"`
	Text             string      `json:"text"`
	Attachments      []*Document `json:"attachments"`
}

type MessageRes struct {
	Completion   string `json:"completion"`
	StopReason   string `json:"stop_reason"`
	Model        string `json:"model"`
	Stop         string `json:"stop"`
	LogId        string `json:"log_id"`
	MessageLimit struct {
		Type string `json:"type"`
	} `json:"messageLimit"`
}

type Conversation struct {
	Uuid         string        `json:"uuid"`
	Name         string        `json:"name"`
	Summary      string        `json:"summary"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
	ChatMessages []ChatMessage `json:"chat_messages"`
}

type ChatMessage struct {
	Uuid         string       `json:"uuid"`
	Text         string       `json:"text"`
	Sender       string       `json:"sender"`
	Index        int          `json:"index"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
	EditedAt     time.Time    `json:"edited_at"`
	ChatFeedback ChatFeedback `json:"chat_feedback"`
	Attachments  []Attachment `json:"attachments"`
}

type ChatFeedback struct {
	Uuid      string    `json:"uuid"`
	Type      string    `json:"type"`
	Reason    string    `json:"reason"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Attachment struct {
	Id               string    `json:"id"`
	FileName         string    `json:"file_name"`
	FileSize         int       `json:"file_size"`
	FileType         string    `json:"file_type"`
	ExtractedContent string    `json:"extracted_content"`
	CreatedAt        time.Time `json:"created_at"`
}

type RenameInfo struct {
	OrganizationUuid string `json:"organization_uuid"`
	ConversationUuid string `json:"conversation_uuid"`
	Title            string `json:"title"`
}
