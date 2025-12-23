package neuronet

import (
	"app/domain"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type Repository interface {
	Translate(ctx context.Context, wordString string) (domain.ResponseTranslate, error)
}

type GeminiRepository struct {
	client *genai.Client
	model  *genai.GenerativeModel
}

func NewGeminiRepository(apiKey string) (Repository, error) {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}

	model := client.GenerativeModel("gemini-2.5-flash")
	return &GeminiRepository{client: client, model: model}, nil
}

func (g *GeminiRepository) Translate(ctx context.Context, wordString string) (domain.ResponseTranslate, error) {
	prompt := fmt.Sprintf(
		`Переведи английское слово "%s" на русский и определи его уровень сложности (A1, A2, B1, B2, C1, C2).
Ответь СТРОГО в формате JSON, без markdown, без пояснений, без текста вне JSON.
Формат:
{"ru_word":"перевод","level":"A1"}`,
		wordString,
	)

	resp, err := g.model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return domain.ResponseTranslate{}, err
	}

	if len(resp.Candidates) == 0 ||
		resp.Candidates[0].Content == nil ||
		len(resp.Candidates[0].Content.Parts) == 0 {
		return domain.ResponseTranslate{}, fmt.Errorf("empty response from model")
	}

	part := resp.Candidates[0].Content.Parts[0]

	textPart, ok := part.(genai.Text)
	if !ok {
		return domain.ResponseTranslate{}, fmt.Errorf("unexpected response type")
	}

	raw := string(textPart)

	fmt.Printf(raw)

	start := strings.Index(raw, "{")
	end := strings.LastIndex(raw, "}")
	if start == -1 || end == -1 || end <= start {
		return domain.ResponseTranslate{}, fmt.Errorf("invalid JSON from model: %s", raw)
	}

	clean := raw[start : end+1]

	var result domain.ResponseTranslate
	if err := json.Unmarshal([]byte(clean), &result); err != nil {
		return domain.ResponseTranslate{}, fmt.Errorf("failed to parse JSON: %w", err)
	}

	if result.Level == "" || result.WordRu == "" {
		return domain.ResponseTranslate{}, errors.New("Word || level is empty")
	}

	return result, nil
}

func (r *GeminiRepository) Close() error {
	return r.client.Close()
}
