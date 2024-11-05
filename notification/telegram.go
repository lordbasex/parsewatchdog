package notification

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/lordbasex/parsewatchdog/config"
)

// TelegramNotifier manages Telegram notifications
type TelegramNotifier struct {
	config *config.Config
}

// NewTelegramNotifier initializes a TelegramNotifier
func NewTelegramNotifier(cfg *config.Config) *TelegramNotifier {
	return &TelegramNotifier{config: cfg}
}

// Send formats and sends a structured message to Telegram with emojis
func (n *TelegramNotifier) Send(subject, message string) error {
	// Parse message to extract timestamp, totalExtensions, and extensions
	parsedData, err := parseTelegramMessage(message)
	if err != nil {
		return fmt.Errorf("error parsing message: %v", err)
	}

	// Formatted message with emojis for Telegram
	formattedMessage := fmt.Sprintf("🚨 *Mass Disconnection Alert* 🚨\n\n📅 *Time:* %s\n🔢 *Total Extensions Disconnected:* %s\n📋 *Extensions List:*\n%s",
		parsedData["Timestamp"],
		parsedData["TotalExtensions"],
		parsedData["Extensions"])

	// Telegram API URL
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", n.config.Telegram.Token)

	// Prepare data for Telegram API
	data := map[string]string{
		"chat_id":    fmt.Sprintf("%d", n.config.Telegram.ChatID),
		"text":       formattedMessage,
		"parse_mode": "Markdown", // Use Markdown for bold and emojis
	}

	// Convert data to JSON and send request
	jsonData, _ := json.Marshal(data)
	_, err = http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	return err
}

// parseTelegramMessage extracts timestamp, totalExtensions, and extensions for Telegram format
func parseTelegramMessage(message string) (map[string]string, error) {
	lines := strings.Split(message, "\n")
	if len(lines) < 3 {
		return nil, fmt.Errorf("message format is invalid")
	}

	// Extract fields assuming specific format
	timestamp := strings.TrimSpace(strings.TrimPrefix(lines[0], "Mass disconnection detected at"))
	totalExtensionsStr := strings.TrimSpace(strings.TrimPrefix(lines[1], "Total: "))
	totalExtensions := strings.Split(totalExtensionsStr, " ")[0] // extract number only
	extensions := strings.TrimSpace(strings.TrimPrefix(lines[2], "Extensions:"))

	// Prepare formatted extensions list with line breaks for readability
	formattedExtensions := strings.ReplaceAll(extensions, ", ", "\n- ")

	// Prepare the data map for the template
	return map[string]string{
		"Timestamp":       timestamp,
		"TotalExtensions": totalExtensions,
		"Extensions":      formattedExtensions,
	}, nil
}
