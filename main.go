package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type Agent struct {
	Name           string
	APIKey         string
	InitialMessage string
	ConversationID string
}

func NewAgent(name string, apiKey string, initialMessage string, conversationID string) *Agent {
	return &Agent{
		Name:           name,
		APIKey:         apiKey,
		InitialMessage: initialMessage,
		ConversationID: conversationID,
	}
}

func (a *Agent) GetChatText(prompt string) (string, error) {
	url := "https://api.openai.com/v1/engines/gpt-3.5-turbo/completions"

	payload := map[string]interface{}{
		"prompt":          prompt,
		"conversation_id": a.ConversationID,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.APIKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	}

	if choices, ok := result["choices"].([]interface{}); ok && len(choices) > 0 {
		if choice, ok := choices[0].(map[string]interface{}); ok {
			if text, ok := choice["text"].(string); ok {
				return text, nil
			}
		}
	}

	return "", fmt.Errorf("failed to get chat text")
}

func haveConversation(agents []*Agent) {
	speakerIndex := 0
	listenerIndex := 1

	for {
		speaker := agents[speakerIndex]
		listener := agents[listenerIndex]

		fmt.Printf("%s is speaking\n", speaker.Name)

		message := fmt.Sprintf("Hello, %s!", listener.Name)

		if speaker.InitialMessage != "" {
			message = speaker.InitialMessage
			speaker.InitialMessage = ""
		}

		if message == fmt.Sprintf("I want to speak to %s", listener.Name) {
			for i, agent := range agents {
				if agent.Name == listener.Name {
					listenerIndex = i
					break
				}
			}
			continue
		}

		response, err := speaker.GetChatText(message)
		if err != nil {
			fmt.Printf("get chat text failed: %v\n", err)
			break
		}

		fmt.Printf("%s: %s\n", speaker.Name, response)

		tempIndex := speakerIndex
		speakerIndex = listenerIndex
		listenerIndex = tempIndex
	}
}

func main() {
	apiKey := os.Getenv("OPENAI_API_KEY")
	initialMessages := strings.Split(os.Getenv("INITIAL_MESSAGES"), ",")
	conversationIDs := strings.Split(os.Getenv("CONVERSATION_IDS"), ",")

	var agents []*Agent
	for i, name := range strings.Split(os.Getenv("AGENT_NAMES"), ",") {
		initialMessage := ""
		if i < len(initialMessages) {
			initialMessage = initialMessages[i]
		}

		conversationID := ""
		if i < len(conversationIDs) {
			conversationID = conversationIDs[i]
		}

		agents = append(agents, NewAgent(name, apiKey, initialMessage, conversationID))
	}

	haveConversation(agents)
}
