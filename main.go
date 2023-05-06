package main

import (
	"fmt"
	"net/http"
	"time"

	chatgpt "github.com/chatgp/chatgpt-go"
)

const (
	PromptEngineer = iota
	ChatGPTAgent
	DomainExpert
)

func main() {
	token := `copy-from-cookies`
	cfValue := "copy-from-cookies"
	puid := "copy-from-cookies"
	cookies := []*http.Cookie{
		{Name: "__Secure-next-auth.session-token", Value: token},
		{Name: "cf_clearance", Value: cfValue},
		{Name: "_puid", Value: puid},
	}

	// create three separate instances of the ChatGPT client
	promptEngineer := chatgpt.NewClient(
		chatgpt.WithDebug(true),
		chatgpt.WithTimeout(60*time.Second),
		chatgpt.WithCookies(cookies),
	)

	chatGPTAgent := chatgpt.NewClient(
		chatgpt.WithDebug(true),
		chatgpt.WithTimeout(60*time.Second),
		chatgpt.WithCookies(cookies),
	)

	domainExpert := chatgpt.NewClient(
		chatgpt.WithDebug(true),
		chatgpt.WithTimeout(60*time.Second),
		chatgpt.WithCookies(cookies),
	)

	// define a simple state machine to keep track of the conversation
	speaker := PromptEngineer
	listener := ChatGPTAgent

	// define a set of rules or conditions that determine when and how the agents communicate with each other
	for {
		var message string
		var err error

		switch speaker {
		case PromptEngineer:
			fmt.Println("Prompt engineer is speaking")

			// decide whether to switch to communicating with the domain expert
			if someConditionIsTrue() {
				fmt.Println("Prompt engineer is switching to communicating with the domain expert")
				listener = DomainExpert
				message = "switch to domain expert"
			} else {
				fmt.Println("Prompt engineer is communicating with the ChatGPT agent")
				listener = ChatGPTAgent
				message = "Hello, ChatGPT agent!"
			}

			_, err = promptEngineer.GetChatText(message)
		case ChatGPTAgent:
			fmt.Println("ChatGPT agent is speaking")
			message = "Hello, prompt engineer!"
			_, err = chatGPTAgent.GetChatText(message)
		case DomainExpert:
			fmt.Println("Domain expert is speaking")
			message = "Hello, prompt engineer!"
			_, err = domainExpert.GetChatText(message)
		}

		if err != nil {
			fmt.Printf("get chat text failed: %v\n", err)
			break
		}

		// update the state of the conversation
		speaker = (speaker + 1) % 3
	}
}
