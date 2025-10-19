package ai

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	l "github.com/CarlosCaravanTsz/imgAI/internal/logger"
	"github.com/sirupsen/logrus"
)

type ChatRequest struct {
	Model    string      `json:"model"`
	Messages []ChatEntry `json:"messages"`
}

type ChatEntry struct {
	Role    string     `json:"role"`
	Content []ChatPart `json:"content"`
}

type ChatPart struct {
	Type     string      `json:"type"`
	Text     string      `json:"text,omitempty"`
	ImageURL *ImageParam `json:"image_url,omitempty"`
}

type ImageParam struct {
	URL string `json:"url"`
}

// Response parsing
type ChatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// Struct for our expected JSON output
type ImageAnalysis struct {
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}

func EncodeImageToBase64Auto(source string) (string, error) {
	var data []byte
	var mimeType string
	var err error

	if strings.HasPrefix(source, "http://") || strings.HasPrefix(source, "https://") {
		resp, err := http.Get(source)
		if err != nil {
			return "", fmt.Errorf("failed to download image: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return "", fmt.Errorf("failed to fetch image, status: %s", resp.Status)
		}

		data, err = io.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("failed to read image bytes: %w", err)
		}

		// Try to detect valid image type
		mimeType = resp.Header.Get("Content-Type")
		if mimeType == "" || !strings.HasPrefix(mimeType, "image/") {
			// Try to guess from URL extension
			ext := filepath.Ext(source)
			mimeType = mime.TypeByExtension(ext)
			if mimeType == "" || !strings.HasPrefix(mimeType, "image/") {
				// Fallback to JPEG if nothing else works
				mimeType = "image/jpeg"
			}
		}
	} else {
		// Handle local files
		data, err = os.ReadFile(source)
		if err != nil {
			return "", fmt.Errorf("failed to read local file: %w", err)
		}

		ext := filepath.Ext(source)
		mimeType = mime.TypeByExtension(ext)
		if mimeType == "" {
			mimeType = "image/jpeg"
		}
	}

	b64 := base64.StdEncoding.EncodeToString(data)
	return fmt.Sprintf("data:%s;base64,%s", mimeType, b64), nil
}

func ObtainDescription(url string) (*ImageAnalysis, error) {

	if os.Getenv("OPENAI_API_KEY") == "" {
		l.LogError("OPENAI_API_KEY not set", logrus.Fields{})
	}

	apiKey := os.Getenv("OPENAI_API_KEY")

	encoded, err := EncodeImageToBase64Auto(url)
	if err != nil {
		fmt.Errorf("Error parsing ")
	}
	// You can also load a local image, encode it to base64, and set ImageURL = "data:image/png;base64,...."

	reqBody := ChatRequest{
		Model: "gpt-4o-mini",
		Messages: []ChatEntry{
			{
				Role: "user",
				Content: []ChatPart{
					{Type: "text", Text: `
Analyze this image and return a JSON object with the following format:
{
  "description": "<short human-readable description>",
  "tags": ["tag1", "tag2", "tag3"]
}
`},
					{Type: "image_url", ImageURL: &ImageParam{URL: encoded}},
				},
			},
		},
	}

	data, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		fmt.Println("Error:", string(body))
		return nil, fmt.Errorf("Error ocurred while calling OpenAI API")
	}

	var chatResp ChatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return nil, err
	}

	content := chatResp.Choices[0].Message.Content

	// Parse the JSON returned by the model
	var analysis ImageAnalysis
	if err := json.Unmarshal([]byte(content), &analysis); err != nil {
		return nil, fmt.Errorf("failed to parse model JSON: %v\nRaw content: %s", err, content)
	}

	fmt.Print(analysis.Description, analysis.Tags)

	return &analysis, nil
}
