package claudeapi

import "time"

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
