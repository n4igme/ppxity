package main

import (
	"bufio"
	"fmt"
	"github.com/0xInception/ppxity/perplexity"
	"github.com/0xInception/ppxity/prompt"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strings"
	"time"
)

var directories []string
var files []string
var extensions []string
var model string
var promptText string
var debug bool
var showInitialPrompt bool
var apiKey string

// var conversationMode bool
var timeout int

var rootCmd = &cobra.Command{
	Use:     "ppxity",
	Example: "ppxity -d C:\\Users\\User\\GolandProjects\\exampleProject -p \"Explain what this project is about.\"",
	Run: func(cmd *cobra.Command, args []string) {
		p := prompt.NewPrompt(promptText, extensions)
		for _, dir := range directories {
			err := p.AddDirectory(dir)
			if err != nil {
				log.Fatalf("Error: %s", err)
			}
		}

		for _, file := range files {
			err := p.AddFile(file)
			if err != nil {
				log.Fatalf("Error: %s", err)
			}
		}

		comp, err := p.Compile()
		if err != nil {
			log.Fatalf("Error: %s", err)
			return
		}

		ppxity := perplexity.NewChatClient(debug, false)

		// Set API key from flag or environment variable
		if apiKey == "" {
			apiKey = os.Getenv("PPLX_API_KEY")
		}
		if apiKey == "" {
			log.Fatal("Error: API key not provided. Set --api-key flag or PPLX_API_KEY environment variable")
		}
		ppxity.SetAPIKey(apiKey)

		err = ppxity.Connect()
		if err != nil {
			log.Fatalf("Error: %s", err)
			return
		}
		defer ppxity.Close()

		// Send the initial message
		err = ppxity.SendMessage(comp, model)
		if err != nil {
			log.Fatalf("Error: %s", err)
			return
		}

		// Get the response directly
		resp := ppxity.GetLastResponse()

		if !showInitialPrompt {
			fmt.Println("User: Initial prompt sent")
		} else {
			log.Println("User: " + comp)
		}

		// Handle continuation if needed (for longer responses)
		// Note: This approach is different from the original socket.io version
		// The new API returns complete responses in one call
		fmt.Println("Assistant:\r\n" + resp + "\r\n\r\n")

		handleUserInput(ppxity, model, time.Second*time.Duration(timeout))
	},
}

func init() {
	rootCmd.Flags().StringSliceVarP(&directories, "directories", "d", []string{}, "Directories to use for the initial prompt")
	rootCmd.Flags().StringSliceVarP(&files, "files", "f", []string{}, "Files to use for the initial prompt")
	rootCmd.Flags().StringSliceVarP(&extensions, "extensions", "e", []string{"go", "txt", "mod", "cs", "c", "rs", "js", "ts"}, "Allowed file extensions to use for the initial prompt")
	rootCmd.Flags().StringVarP(&model, "model", "m", perplexity.CLAUDE, "Perplexity model to use: e.g. '"+perplexity.CLAUDE+"' ("+strings.Join(perplexity.ALL_MODELS, ", ")+")")
	rootCmd.Flags().StringVarP(&promptText, "prompt", "p", "", "Initial prompt for the conversation: e.g. 'Hello, World!'")
	rootCmd.Flags().StringVarP(&apiKey, "api-key", "k", "", "Perplexity API key (or set PPLX_API_KEY environment variable)")
	rootCmd.Flags().BoolVarP(&debug, "debug", "D", false, "Enable debug mode")
	rootCmd.Flags().BoolVarP(&showInitialPrompt, "show-initial-prompt", "s", false, "Show the initial prompt")
	//rootCmd.Flags().BoolVarP(&conversationMode, "conversation", "C", false, "Enable conversation mode")
	rootCmd.Flags().IntVarP(&timeout, "timeout", "t", 50, "Timeout in seconds for receiving messages")

	_ = rootCmd.MarkFlagRequired("prompt")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func handleUserInput(ppxity *perplexity.ChatClient, model string, timeout time.Duration) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		x := scanner.Text()
		if x == "exit" {
			return
		} else if x == "backtrack" {
			err := ppxity.Backtrack()
			if err != nil {
				log.Printf("Error backtracking: %v", err)
			}
			continue
		}
		err := ppxity.SendMessage(x, model)
		if err != nil {
			log.Printf("Error sending message: %v", err)
			continue
		}

		// Get the response directly
		resp := ppxity.GetLastResponse()
		fmt.Println(resp)
	}
}
