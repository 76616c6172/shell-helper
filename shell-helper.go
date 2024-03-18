package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/sashabaranov/go-openai"
)

const apiKeyFile = ".config/shell-helper/key"

func getChatGPTResponse(client *openai.Client, prompt string) (string, error) {
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "You write modern, excellent concise bash oneliners. All codeblocks must start with the appropriate language identifiers.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)
	if err != nil {
		return "", err
	}
	return resp.Choices[0].Message.Content, nil
}

func extractCodeBlock(text string) string {
	pattern := regexp.MustCompile("```bash\n(.*?)\n```")
	match := pattern.FindStringSubmatch(text)
	if len(match) > 1 {
		return strings.TrimSpace(match[1])
	}
	return ""
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ? <prompt>")
		os.Exit(1)
	}

	userInput := strings.Join(os.Args[1:], " ")
	prompt := fmt.Sprintf("I have the linux terminal open, I want to %s. Output the correct bash one-liner to do this in a single codeblock. No japping!", userInput)

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting user home directory:", err)
		os.Exit(1)
	}

	keyPath := filepath.Join(homeDir, apiKeyFile)
	apiKey, err := os.ReadFile(keyPath)
	if err != nil {
		fmt.Println("Error reading API key file:", err)
		os.Exit(1)
	}

	client := openai.NewClient(strings.TrimSpace(string(apiKey)))
	response, err := getChatGPTResponse(client, prompt)
	if err != nil {
		fmt.Println("Error getting ChatGPT response:", err)
		os.Exit(1)
	}

	fmt.Println(response)

	shellCmd := extractCodeBlock(response)
	if shellCmd != "" {
		fmt.Print("Shell command detected. Execute? ")
		reader := bufio.NewReader(os.Stdin)
		userInput, _ := reader.ReadString('\n')
		if strings.TrimSpace(userInput) == "" {
			cmd := exec.Command("bash", "-c", shellCmd)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err := cmd.Run()
			if err != nil {
				fmt.Println("Error executing shell command:", err)
			}
		}
	}
}
