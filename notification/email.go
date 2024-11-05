package notification

import (
	"bytes"
	"fmt"
	"net/smtp"

	"github.com/lordbasex/parsewatchdog/config"

	"strings"
	"text/template"
)

// EmailNotifier manages email notifications
type EmailNotifier struct {
	config *config.Config
}

// NewEmailNotifier initializes an EmailNotifier
func NewEmailNotifier(cfg *config.Config) *EmailNotifier {
	return &EmailNotifier{config: cfg}
}

// HTML template for the email content
const emailTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Mass Disconnection Alert</title>
</head>
<body style="font-family: Arial, sans-serif; color: #333;">
    <h2>Mass Disconnection Alert</h2>
    <p><strong>Timestamp:</strong> {{.Timestamp}}</p>
    <p><strong>Total Extensions Disconnected:</strong> {{.TotalExtensions}}</p>
    <p><strong>Extensions:</strong> {{.Extensions}}</p>
    <hr>
    <p>Please review this issue as soon as possible.</p>
</body>
</html>
`

// Send sends an email using the HTML template and provided arguments
func (n *EmailNotifier) Send(subject, message string) error {
	from := n.config.SMTP.User
	pass := n.config.SMTP.Pass
	host := n.config.SMTP.Host
	port := n.config.SMTP.Port
	to := n.config.SMTP.Recipients

	// Concatenate host and port
	address := fmt.Sprintf("%s:%d", host, port)
	auth := smtp.PlainAuth("", from, pass, host)

	// Parse `message` to extract timestamp, totalExtensions, and extensions
	parsedData, err := parseMessage(message)
	if err != nil {
		return fmt.Errorf("error parsing message: %v", err)
	}

	// Parse and execute the email template with extracted data
	tmpl, err := template.New("emailTemplate").Parse(emailTemplate)
	if err != nil {
		return fmt.Errorf("error parsing template: %v", err)
	}

	var body bytes.Buffer
	body.WriteString("To: " + strings.Join(to, ",") + "\r\n")
	body.WriteString("Subject: " + subject + "\r\n")
	body.WriteString("MIME-version: 1.0;\r\n")
	body.WriteString("Content-Type: text/html; charset=\"UTF-8\";\r\n")
	body.WriteString("\r\n")

	// Execute template with the parsed data
	err = tmpl.Execute(&body, parsedData)
	if err != nil {
		return fmt.Errorf("error executing template: %v", err)
	}

	// Send the email
	err = smtp.SendMail(address, auth, from, to, body.Bytes())
	if err != nil {
		return fmt.Errorf("error sending email: %v", err)
	}

	return nil
}

// parseMessage extracts timestamp, totalExtensions, and extensions from the message string
func parseMessage(message string) (map[string]interface{}, error) {
	lines := strings.Split(message, "\n")
	if len(lines) < 3 {
		return nil, fmt.Errorf("message format is invalid")
	}

	// Extract fields from lines assuming specific format
	timestamp := strings.TrimSpace(strings.TrimPrefix(lines[0], "Mass disconnection detected at"))
	totalExtensionsStr := strings.TrimSpace(strings.TrimPrefix(lines[1], "Total: "))
	totalExtensions := strings.Split(totalExtensionsStr, " ")[0] // extract number only
	extensions := strings.TrimSpace(strings.TrimPrefix(lines[2], "Extensions:"))

	// Prepare the data map for the template
	return map[string]interface{}{
		"Timestamp":       timestamp,
		"TotalExtensions": totalExtensions,
		"Extensions":      extensions,
	}, nil
}
