package claudeapi

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/imroc/req/v3"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"time"
)

const (
	idleConnsPerHost = 25
	defaultTimeout   = 30 * time.Second
)

var defaultHeaders = map[string]string{
	"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36",
	"Referer":    "https://claude.ai/chats",
	"Origin":     "https://claude.ai",
}

type RequestHandler func(req *req.Request) error
type ClaudeOption func(claude *Claude) error

type Claude struct {
	organizationId string
	baseReqHandler RequestHandler
	client         *req.Client
}

func WithBaseHttpHandler(reqHandler RequestHandler) ClaudeOption {
	return func(claude *Claude) error {
		claude.baseReqHandler = reqHandler
		return nil
	}
}

func (c *Claude) GetOrganizationId(organizationIndex int) (string, error) {
	claudeUrl := "https://claude.ai/api/organizations"

	var organizations Organizations
	var resError ErrorRes

	request := c.client.R().SetSuccessResult(&organizations).SetErrorResult(&resError)

	request.SetHeaders(defaultHeaders)
	request.SetHeader("Content-Type", "application/json")

	if c.baseReqHandler != nil {
		err := c.baseReqHandler(request)
		if err != nil {
			return "", err
		}
	}

	resp, err := request.Get(claudeUrl)
	if err != nil {
		return "", err
	}

	if resp.IsErrorState() {
		return "", resError.Err()
	}

	if len(organizations) == 0 {
		return "", errors.New("organizations is empty")
	}
	if organizationIndex > len(organizations)-1 {
		return "", errors.New("organizationIndex exceeds length")
	}
	return organizations[organizationIndex].Uuid, nil
}

func (c *Claude) ListConversations() (chatConversations ListConversations, err error) {
	claudeUrl := fmt.Sprintf("https://claude.ai/api/organizations/%s/chat_conversations", c.organizationId)

	var resError ErrorRes

	request := c.client.R().SetSuccessResult(&chatConversations).SetErrorResult(&resError)

	request.SetHeaders(defaultHeaders)
	request.SetHeader("Content-Type", "application/json")

	if c.baseReqHandler != nil {
		err = c.baseReqHandler(request)
		if err != nil {
			return nil, err
		}
	}

	resp, err := request.Get(claudeUrl)
	if err != nil {
		return nil, err
	}

	if resp.IsErrorState() {
		err = resError.Err()
		return nil, err
	}

	return chatConversations, nil
}

func (c *Claude) CreateConversation(name string, summary string) (*ConversationInfo, error) {
	claudeUrl := fmt.Sprintf("https://claude.ai/api/organizations/%s/chat_conversations", c.organizationId)

	var err error
	conversationInfo := &ConversationInfo{}
	resError := &ErrorRes{}

	body := &ConversationInfo{
		Uuid:    uuid.New().String(),
		Name:    name,
		Summary: summary,
	}
	bytesData, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	request := c.client.R().SetSuccessResult(conversationInfo).SetErrorResult(resError).SetBody(bytesData)

	request.SetHeaders(defaultHeaders)
	request.SetHeader("Content-Type", "application/json")

	if c.baseReqHandler != nil {
		err = c.baseReqHandler(request)
		if err != nil {
			return nil, err
		}
	}

	resp, err := request.Post(claudeUrl)
	if err != nil {
		return nil, err
	}

	if resp.IsErrorState() {
		err = resError.Err()
		return nil, err
	}

	return conversationInfo, nil
}

func (c *Claude) DeleteConversation(conversationId string) (bool, error) {
	claudeUrl := fmt.Sprintf("https://claude.ai/api/organizations/%s/chat_conversations/%s", c.organizationId, conversationId)

	var err error
	resError := &ErrorRes{}

	body := conversationId

	request := c.client.R().SetErrorResult(resError).SetBody(body)

	request.SetHeaders(defaultHeaders)
	request.SetHeader("Content-Type", "application/json")

	if c.baseReqHandler != nil {
		err = c.baseReqHandler(request)
		if err != nil {
			return false, err
		}
	}

	resp, err := request.Delete(claudeUrl)
	if err != nil {
		return false, err
	}

	if resp.IsErrorState() {
		err = resError.Err()
		return false, err
	}

	return true, nil
}

func (c *Claude) RenameConversation(conversationId string, title string) (bool, error) {
	claudeUrl := "https://claude.ai/api/rename_chat"

	var err error
	resError := &ErrorRes{}

	body := &RenameInfo{
		OrganizationUuid: c.organizationId,
		ConversationUuid: conversationId,
		Title:            title,
	}
	bytesData, err := json.Marshal(body)
	if err != nil {
		return false, err
	}

	request := c.client.R().SetErrorResult(resError).SetBody(bytesData)

	request.SetHeaders(defaultHeaders)
	request.SetHeader("Content-Type", "application/json")

	if c.baseReqHandler != nil {
		err = c.baseReqHandler(request)
		if err != nil {
			return false, err
		}
	}

	resp, err := request.Post(claudeUrl)
	if err != nil {
		return false, err
	}

	if resp.IsErrorState() {
		err = resError.Err()
		return false, err
	}

	return true, nil
}

func (c *Claude) isBinaryFile(data []byte) bool {
	if len(data) > 512 {
		data = data[:512]
	}
	for _, v := range data {
		if v == 0 {
			return true
		}
	}
	return false
}

func (c *Claude) ConvertDocument(fileName string, content []byte) (document *Document, err error) {
	if len(content) > 10*1024*1024 {
		return nil, errors.New("the file is larger than 10MB")
	}

	contentType := http.DetectContentType(content)
	if !c.isBinaryFile(content) {
		document.FileName = fileName
		document.FileSize = len(content)
		document.FileType = contentType
		document.ExtractedContent = string(content)
		return document, nil
	}

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filepath.Base(fileName))
	if err != nil {
		return nil, err
	}

	_, err = part.Write(content)
	if err != nil {
		return nil, err
	}

	err = writer.WriteField("orgUuid", c.organizationId)
	if err != nil {
		return nil, err
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	claudeUrl := "https://claude.ai/api/convert_document"
	var resError ErrorRes
	request := c.client.R().SetSuccessResult(&document).SetErrorResult(&resError)

	request.SetHeaders(defaultHeaders)
	request.SetHeader("Content-Type", writer.FormDataContentType())

	if c.baseReqHandler != nil {
		err = c.baseReqHandler(request)
		if err != nil {
			return nil, err
		}
	}

	request.SetBody(body)

	resp, err := request.Post(claudeUrl)
	if err != nil {
		return nil, err
	}

	if resp.IsErrorState() {
		err = resError.Err()
		return nil, err
	}

	return document, nil
}

func (c *Claude) AppendMessage(prompt string, conversationId string, attachments []*Document) (completion string, err error) {
	if attachments == nil {
		attachments = make([]*Document, 0)
	}

	claudeUrl := "https://claude.ai/api/append_message"
	var resError ErrorRes

	location := time.Local

	body := &Message{
		Completion: Completion{
			Prompt:   prompt,
			Timezone: location.String(),
			Model:    "claude-2",
		},
		OrganizationUuid: c.organizationId,
		ConversationUuid: conversationId,
		Text:             prompt,
		Attachments:      attachments,
	}
	bytesData, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	request := c.client.R().SetErrorResult(&resError).SetBody(bytes.NewReader(bytesData))

	request.SetHeaders(defaultHeaders)
	request.SetHeader("Content-Type", "application/json")

	if c.baseReqHandler != nil {
		err = c.baseReqHandler(request)
		if err != nil {
			return "", err
		}
	}

	resp, err := request.Post(claudeUrl)
	if err != nil {
		return "", err
	}

	if resp.IsErrorState() {
		err = resError.Err()
		return "", err
	}

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		resStr := scanner.Bytes()
		if len(resStr) > 0 {
			resStr = resStr[6:]
			messageRes := &MessageRes{}
			err = json.Unmarshal(resStr, messageRes)
			if err != nil {
				return "", err
			}
			if messageRes.StopReason == "stop_sequence" {
				break
			}
			completion += messageRes.Completion
		}
	}

	return completion, nil
}

func NewClaude(config *Config, options ...ClaudeOption) (*Claude, error) {
	var client *req.Client
	if config.Debug == true {
		client = req.DevMode()
	} else {
		client = req.DefaultClient()
	}
	client.ImpersonateChrome()

	proxy := http.ProxyFromEnvironment
	if config.Proxy != nil {
		proxy = config.Proxy
	}
	client.SetProxy(proxy)

	timeout := defaultTimeout
	if config.Timeout != 0 {
		timeout = config.Timeout
	}
	client.SetTimeout(timeout)

	if len(config.Cookies) == 0 {
		return nil, errors.New("cookies is nil")
	}
	client.SetCommonCookies(config.Cookies...)

	claude := &Claude{
		client: client,
	}

	for _, option := range options {
		if err := option(claude); err != nil {
			return nil, err
		}
	}

	organizationId, err := claude.GetOrganizationId(config.OrganizationIndex)
	if err != nil {
		return nil, err
	}

	claude.organizationId = organizationId

	return claude, nil
}
