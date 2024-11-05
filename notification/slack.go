package notification

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/lordbasex/parsewatchdog/config"
)

// SlackNotifier manages Slack notifications
type SlackNotifier struct {
	config *config.Config
}

// NewSlackNotifier initializes a SlackNotifier
func NewSlackNotifier(cfg *config.Config) *SlackNotifier {
	return &SlackNotifier{config: cfg}
}

// Send formats and sends a structured message to Slack
func (n *SlackNotifier) Send(subject, message string) error {
	// Parse message to extract timestamp, totalExtensions, and extensions
	parsedData, err := parseSlackMessage(message)
	if err != nil {
		return fmt.Errorf("error parsing message: %v", err)
	}

	// Formatted message for Slack
	formattedMessage := fmt.Sprintf("🚨 *Mass Disconnection Alert* 🚨\n\n📅 *Time:* %s\n🔢 *Total Extensions Disconnected:* %s\n📋 *Extensions List:*\n%s",
		parsedData["Timestamp"],
		parsedData["TotalExtensions"],
		parsedData["Extensions"])

	// Prepare data for Slack webhook
	data := map[string]string{
		"text": formattedMessage,
	}

	// Convert data to JSON
	jsonData, _ := json.Marshal(data)

	// Send request to Slack webhook
	req, err := http.NewRequest("POST", n.config.Slack.WebhookURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error sending message to Slack: status code %d", resp.StatusCode)
	}
	return nil
}

func parseSlackMessage(message string) (map[string]string, error) {
	lines := strings.Split(message, "\n")
	if len(lines) < 3 {
		return nil, fmt.Errorf("message format is invalid")
	}

	// Extraer campos asumiendo un formato específico
	timestamp := strings.TrimSpace(strings.TrimPrefix(lines[0], "Mass disconnection detected at"))
	totalExtensionsStr := strings.TrimSpace(strings.TrimPrefix(lines[1], "Total: "))
	totalExtensions := strings.Split(totalExtensionsStr, " ")[0]
	extensions := strings.TrimSpace(strings.TrimPrefix(lines[2], "Extensions:"))

	// Formatear la lista de extensiones para Slack
	formattedExtensions := strings.ReplaceAll(extensions, ", ", "\n• ")

	return map[string]string{
		"Timestamp":       timestamp,
		"TotalExtensions": totalExtensions,
		"Extensions":      formattedExtensions,
	}, nil
}
