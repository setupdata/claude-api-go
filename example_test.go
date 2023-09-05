package claudeapi

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
)

var claude *Claude

func getClaude() (*Claude, error) {
	if claude == nil {
		var err error
		claude, err = NewClaude(
			&Config{
				Proxy: func(request *http.Request) (*url.URL, error) {
					return url.Parse("http://localhost:7890")
				},
				Cookies: CreateCookies(map[string]string{
					"sessionKey": "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",
				}),
			},
		)
		if err != nil {
			// Handle error.
			log.Println(err.Error())
			return nil, err
		}
	}
	return claude, nil
}

func ExampleNewClaude() {
	var err error
	claude, err = getClaude()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	chatConversations, err := claude.ListConversations()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for _, chatConversation := range chatConversations {
		log.Printf("%#v\n", chatConversation)
	}
	fmt.Printf("true")
	return
	// Output: true
}

func ExampleClaude_ConvertDocument() {
	var err error
	claude, err = getClaude()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	//fileName := "https://filesamples.com/samples/document/docx/sample3.docx"
	document, err := claude.ConvertDocument("fileName", nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	log.Printf("%#v\n", document)
	fmt.Printf("true")
	return
	// Output: true
}

func ExampleClaude_AppendMessage() {
	var err error
	claude, err = getClaude()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	conversations, err := claude.ListConversations()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	completion, err := claude.AppendMessage("who are you?", conversations[0].Uuid, nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	log.Println(completion)

	fmt.Printf("true")
	return
	// Output: true
}
