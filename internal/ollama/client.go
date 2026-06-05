package ollama

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type Client struct {
	Host  string
	Model string
}

func New(host, model string) *Client {
	return &Client{Host: host, Model: model}
}

type message struct {
	Role    string   `json:"role"`
	Content string   `json:"content"`
	Images  []string `json:"images,omitempty"`
}

type chatRequest struct {
	Model    string    `json:"model"`
	Messages []message `json:"messages"`
	Options  any       `json:"options"`
	Stream   bool      `json:"stream"`
}

type chatResponse struct {
	Message struct {
		Content string `json:"content"`
	} `json:"message"`
	Done bool `json:"done"`
}

// Ask sends a prompt with optional image paths to the model and returns the response.
func (c *Client) Ask(prompt string, imagePaths []string) (string, error) {
	images := make([]string, 0, len(imagePaths))
	for _, p := range imagePaths {
		b64, err := encodeImage(p)
		if err != nil {
			return "", fmt.Errorf("encoding %s: %w", p, err)
		}
		images = append(images, b64)
	}

	body := chatRequest{
		Model: c.Model,
		Messages: []message{
			{Role: "user", Content: prompt, Images: images},
		},
		Options: map[string]any{"temperature": 0.3},
		Stream:  true,
	}

	data, _ := json.Marshal(body)
	req, err := http.NewRequest("POST", c.Host+"/api/chat", bytes.NewReader(data))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	hc := &http.Client{Timeout: 120 * time.Second}
	resp, err := hc.Do(req)
	if err != nil {
		return "", fmt.Errorf("ollama request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var buf bytes.Buffer
		buf.ReadFrom(resp.Body)
		return "", fmt.Errorf("ollama error %d: %s", resp.StatusCode, buf.String())
	}

	var result bytes.Buffer
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		var chunk chatResponse
		if err := json.Unmarshal(scanner.Bytes(), &chunk); err != nil {
			continue
		}
		result.WriteString(chunk.Message.Content)
		if chunk.Done {
			break
		}
	}

	return result.String(), scanner.Err()
}

func encodeImage(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}
