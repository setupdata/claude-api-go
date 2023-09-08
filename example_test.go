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

func ExampleClaude_CreateConversation() {
	var err error
	claude, err = getClaude()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	conversationInfo, err := claude.CreateConversation("test_test", "test-test-test")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	log.Println(conversationInfo)
	fmt.Printf("true")
	return
	// Output: true
}

func ExampleClaude_DeleteConversation() {
	var err error
	claude, err = getClaude()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	conversationInfo, err := claude.CreateConversation("test_test", "test-test-test")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	log.Println(conversationInfo)

	success, err := claude.DeleteConversation(conversationInfo.Uuid)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	log.Println(success)
	fmt.Printf("true")
	return
	// Output: true
}

func ExampleClaude_RenameConversation() {
	var err error
	claude, err = getClaude()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	conversationInfo, err := claude.CreateConversation("test_test", "test-test-test")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	log.Println(conversationInfo)

	success, err := claude.DeleteConversation(conversationInfo.Uuid)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	log.Println(success)
	fmt.Printf("true")
	return
	// Output: true
}
