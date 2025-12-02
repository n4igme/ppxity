package perplexity

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestChatClient_Start(t *testing.T) {
	// Skip test if no API key is provided
	apiKey := os.Getenv("PPLX_API_KEY")
	if apiKey == "" {
		t.Skip("Skipping test: PPLX_API_KEY not set")
	}

	ppxity := NewChatClient(true, false)
	ppxity.SetAPIKey(apiKey)
	err := ppxity.Connect()
	if err != nil {
		t.Fatalf("Error: %s", err)
	}
	defer ppxity.Close()
	err = ppxity.SendMessage("Hello, World!", CLAUDE)
	if err != nil {
		t.Fatalf("Error: %s", err)
	}

	// Get response directly instead of using ReceiveMessage
	response := ppxity.GetLastResponse()
	fmt.Println("Response:", response)
	err = ppxity.SendMessage("What message did i send?", CLAUDE)
	if err != nil {
		t.Fatalf("Error: %s", err)
	}

	// Get second response
	response = ppxity.GetLastResponse()
	fmt.Println("Response:", response)
	time.Sleep(time.Second * 2) // Reduced sleep time since API responses are immediate
}
