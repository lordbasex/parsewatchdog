package notification

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/lordbasex/parsewatchdog/config"
)

type APINotifier struct {
	config *config.Config
}

func NewAPINotifier(cfg *config.Config) *APINotifier {
	return &APINotifier{config: cfg}
}

func (n *APINotifier) Send(subject, message string) error {
	payload := map[string]string{"subject": subject, "message": message}
	jsonData, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", n.config.API.Endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("API-Key", n.config.API.APIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send API notification: %w", err)
	}
	return nil
}
