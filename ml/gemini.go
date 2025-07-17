package ml

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"google.golang.org/genai"
)

var titleRegex = regexp.MustCompile(`(.+?).\((\d\d\d\d)\);?`)

func SendGeminiRequest(collection *Collection) ([]*MovieInfo, error) {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		//TODO remove key
		APIKey:  "AIzaSyAa3yBFAwIOOsZb31ohe3tBO50w1RKoyeE",
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	fullPrompt := `Ты кинокритик. Дай подборку фильмов. Ответ должен быть только в виде списка названий, разделенных ТОЛЬКО точкой с запятой. 
Правила:
1. Используй только оригинальные названия
2. Указывай год выпуска
3. Используй только реально существующие фильмы
4. Формат: Title1 (Year1); Title2 (Year2); Title3 (Year3); ...
5. Не добавляй номера, точки, кавычки или другие символы
6. Не более 30 фильмов
7. Короткометражные фильмы и сериалы не включать в список
8. Строго запрещено повторять названия фильмов в ответе
9. Кино произведенное в СССР или России должны быть названия на русском
` + collection.Prompt

	contents := []*genai.Content{
		genai.NewContentFromText(fullPrompt, genai.RoleUser),
	}

	resp, err := client.Models.GenerateContent(ctx, "gemma-3-27b-it", contents, nil)
	if err != nil {
		return nil, fmt.Errorf("Gemini API error: %w", err)
	}

	if resp == nil || len(resp.Candidates) == 0 || resp.Candidates[0].Content == nil {
		return nil, fmt.Errorf("Gemini returned empty response")
	}

	var generatedTextBuilder strings.Builder
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				generatedTextBuilder.WriteString(string(part.Text))
			}
		}
	}

	fmt.Println("Resp:", generatedTextBuilder.String())
	return parseTitles(generatedTextBuilder.String()), nil
}

func parseTitles(response string) []*MovieInfo {
	matches := titleRegex.FindAllStringSubmatch(response, -1)
	var movies []*MovieInfo
	for _, match := range matches {
		if len(match) > 2 {
			title := cleanTitle(match[1])
			year := match[2]
			movies = append(movies, &MovieInfo{Title: title, Year: year})
		}
	}
	return movies
}

func cleanTitle(title string) string {
	title = strings.TrimSpace(title)
	return strings.Trim(title, `"'`)
}
