package agent

import (
	"encoding/json"
	"strings"
	"unicode/utf8"

	"github.com/jumpserver/wisp/pkg/agent/provider"
)

const (
	truncatedPromptMarker       = "[earlier content truncated]\n"
	middleTruncatedPromptMarker = "\n[middle content truncated]\n"
)

func withResponseLanguage(system, responseLanguage string) string {
	if responseLanguage == "" {
		return system
	}
	return system + "\nThe trusted interface language is " +
		responseLanguage + ". Write every user-visible natural-language field " +
		"in this language, including answers, status text, explanations and " +
		"summaries. The interface language takes precedence over the language " +
		"of the request and evidence. Do not translate SQL, identifiers or " +
		"quoted content."
}

func normalizeResponseLanguage(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	value = strings.ReplaceAll(value, "_", "-")
	switch {
	case value == "":
		return ""
	case value == "zh-hant" || strings.HasPrefix(value, "zh-hant-") ||
		value == "zh-tw" || strings.HasPrefix(value, "zh-tw-") ||
		value == "zh-hk" || strings.HasPrefix(value, "zh-hk-") ||
		value == "zh-mo" || strings.HasPrefix(value, "zh-mo-"):
		return "Traditional Chinese (繁體中文)"
	case value == "zh" || value == "zh-hans" ||
		strings.HasPrefix(value, "zh-hans-") ||
		value == "zh-cn" || strings.HasPrefix(value, "zh-cn-") ||
		value == "zh-sg" || strings.HasPrefix(value, "zh-sg-"):
		return "Simplified Chinese (简体中文)"
	case value == "ja" || strings.HasPrefix(value, "ja-"):
		return "Japanese"
	case value == "ko" || strings.HasPrefix(value, "ko-"):
		return "Korean"
	case value == "es" || strings.HasPrefix(value, "es-"):
		return "Spanish"
	case value == "pt" || strings.HasPrefix(value, "pt-"):
		return "Portuguese"
	case value == "ru" || strings.HasPrefix(value, "ru-"):
		return "Russian"
	case value == "en" || strings.HasPrefix(value, "en-"):
		return "English"
	default:
		return "English"
	}
}

func decodeModelJSON(content string, output any) error {
	content = strings.TrimSpace(content)
	content = strings.TrimPrefix(content, "```json")
	content = strings.TrimPrefix(content, "```")
	content = strings.TrimSuffix(content, "```")
	if err := json.Unmarshal([]byte(strings.TrimSpace(content)), output); err != nil {
		return provider.NewOutputError(
			provider.ErrorInvalidOutput, "decode model JSON: %v", err,
		)
	}
	return nil
}

func mustJSON(value any) string {
	result, _ := json.Marshal(value)
	return string(result)
}

func promptTail(value string, limit int) string {
	value = strings.ToValidUTF8(value, "\uFFFD")
	if len(value) <= limit {
		return value
	}
	available := limit - len(truncatedPromptMarker)
	if available <= 0 {
		return truncatedPromptMarker[:limit]
	}
	start := len(value) - available
	for start < len(value) && !utf8.RuneStart(value[start]) {
		start++
	}
	return truncatedPromptMarker + value[start:]
}

func headTailPrompt(value string, limit int) string {
	value = strings.ToValidUTF8(value, "\uFFFD")
	if limit <= 0 {
		return ""
	}
	if len(value) <= limit {
		return value
	}
	available := limit - len(middleTruncatedPromptMarker)
	if available <= 0 {
		return promptTail(value, limit)
	}
	headBytes := available / 2
	tailBytes := available - headBytes
	headEnd := headBytes
	for headEnd > 0 && headEnd < len(value) && !utf8.RuneStart(value[headEnd]) {
		headEnd--
	}
	tailStart := len(value) - tailBytes
	for tailStart < len(value) && !utf8.RuneStart(value[tailStart]) {
		tailStart++
	}
	return value[:headEnd] + middleTruncatedPromptMarker + value[tailStart:]
}
