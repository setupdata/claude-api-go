package claudeapi

import (
	"errors"
	"fmt"
	"time"
)

type Base struct {
	Uuid      string    `json:"uuid"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Settings struct {
	ClaudeConsolePrivacy string `json:"claude_console_privacy"`
}

type Organization struct {
	Base         `json:",inline"`
	Name         string   `json:"name"`
	Settings     Settings `json:"settings"`
	Capabilities []string `json:"capabilities"`
	ActiveFlags  []string `json:"active_flags"`
}

type Organizations []*Organization

type ChatConversation struct {
	Base    `json:",inline"`
	Name    string `json:"name"`
	Summary string `json:"summary"`
}
type ChatConversations []*ChatConversation

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
