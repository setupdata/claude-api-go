package claudeapi

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
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

func (c *Claude) GetOrganizationId(apiReqHandles ...RequestHandler) (string, error) {
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

	for _, apiReqHandle := range apiReqHandles {
		err := apiReqHandle(request)
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
	return organizations[0].Uuid, nil
}

func (c *Claude) ListConversations(apiReqHandles ...RequestHandler) (chatConversations ChatConversations, err error) {
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

	for _, apiReqHandle := range apiReqHandles {
		err = apiReqHandle(request)
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

// ConvertDocument Make sure the file exists and has read permission
func (c *Claude) ConvertDocument(fileName string, content []byte, apiReqHandles ...RequestHandler) (document *Document, err error) {
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

	for _, apiReqHandle := range apiReqHandles {
		err = apiReqHandle(request)
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

func (c *Claude) AppendMessage(prompt string, conversationId string, attachments []*Document, apiReqHandles ...RequestHandler) (completion string, err error) {
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

	request := c.client.R().SetErrorResult(&resError)

	request.SetHeaders(defaultHeaders)
	request.SetHeader("Content-Type", "application/json")

	if c.baseReqHandler != nil {
		err = c.baseReqHandler(request)
		if err != nil {
			return "", err
		}
	}

	for _, apiReqHandle := range apiReqHandles {
		err = apiReqHandle(request)
		if err != nil {
			return "", err
		}
	}

	request.SetBody(bytes.NewReader(bytesData))

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
			var messageRes *MessageRes
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
	client := req.DevMode()
	client.ImpersonateFirefox()

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

	if len(config.cookies) == 0 {
		return nil, errors.New("cookies is nil")
	}
	client.SetCommonCookies(config.cookies...)

	claude := &Claude{
		client: client,
	}

	for _, option := range options {
		if err := option(claude); err != nil {
			return nil, err
		}
	}

	organizationId, err := claude.GetOrganizationId()
	if err != nil {
		return nil, err
	}

	claude.organizationId = organizationId

	return claude, nil
}
