package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/lordbasex/parsewatchdog/config"
	"github.com/lordbasex/parsewatchdog/notification"
)

var lastAlertTimestamps = make(map[string]struct{})

const defaultConfigPath = "/etc/parsewatchdog.conf"
const defaultConfigContent = `
[smtp]
enabled=false
host=smtp.gmail.com
port=587
user=your_email@gmail.com
pass=your_email_password
recipients=recipient1@example.com,recipient2@example.com

[telegram]
enabled=false
token=your_telegram_bot_token
chat_id=123456789

[api]
enabled=false
endpoint=https://example.com/notify
api_key=your_api_key

[rabbitmq]
enabled=false
type=amqp
user=rabbitmq_user
password=rabbitmq_password
ip=192.168.0.10
port=5672
queue=parsewatchdog_notifications

[slack]
enabled=false
webhook_url=https://hooks.slack.com/services/XXXXXXXXXXX/XXXXXXXXXXXX/XXXXXXXXXXXXXXXXXXXXXXXX

[debug]
debug_level=1
`

func main() {
	//log.SetFlags(log.LstdFlags | log.Lshortfile)

	fmt.Println("Starting ParseWatchdog...")

	// Ensure configuration file exists, creating it with default content if necessary
	err := createDefaultConfig(defaultConfigPath, defaultConfigContent)
	if err != nil {
		log.Fatalf("Error creating default config: %v", err)
	}

	// Load configuration
	cfg, err := config.LoadConfig(defaultConfigPath)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	fmt.Printf("ParseWatchdog Pid(%d)\n", os.Getpid())
	fmt.Printf("\n [*] Version: %s (%s)", config.Version, config.DaemonGitBuild)
	fmt.Printf("\n [*] Build Date: %s \n\n", config.DaemonGitBuildDate)

	filePath := "/var/log/asterisk/full"
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	defer file.Close()

	// Start from the end of the file
	file.Seek(0, os.SEEK_END)

	// Monitor at regular intervals
	for {
		checkLogForUnreachable(file, cfg)
		time.Sleep(1 * time.Second)
	}
}

func createDefaultConfig(path string, content string) error {
	// Check if the config file already exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// Create the config file with default content
		file, err := os.Create(path)
		if err != nil {
			return fmt.Errorf("failed to create config file: %w", err)
		}
		defer file.Close()

		_, err = file.WriteString(content)
		if err != nil {
			return fmt.Errorf("failed to write default config content: %w", err)
		}

		log.Printf("Default configuration created at %s", path)
	}
	return nil
}

func logMessage(cfg *config.Config, level int, message string) {
	if cfg.Debug.DebugLevel >= level {
		log.Println(message)
	}
}

func checkLogForUnreachable(file *os.File, cfg *config.Config) {
	scanner := bufio.NewScanner(file)
	unreachableEvents := make(map[string][]string)

	// Updated regex to support both chan_sip (Peer) and pjsip (Endpoint), considering case sensitivity for UNREACHABLE
	re := regexp.MustCompile(`^\[(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2})] .+(Endpoint|Peer) ['"]?(\d+)['"]? is now (Unreachable|UNREACHABLE)`)

	logMessage(cfg, 2, "Incremental log file reading...")

	for scanner.Scan() {
		line := scanner.Text()
		match := re.FindStringSubmatch(line)
		if match != nil {
			timestampStr := match[1]
			extension := match[3]

			logMessage(cfg, 2, fmt.Sprintf("Reading log line: %s", line))

			// Process only "now Unreachable/UNREACHABLE" events
			unreachableEvents[timestampStr] = append(unreachableEvents[timestampStr], extension)
		}
	}

	for timestamp, extensions := range unreachableEvents {
		if len(extensions) > 1 {
			logMessage(cfg, 1, fmt.Sprintf("Mass disconnection detected at %s: %d extensions", timestamp, len(extensions)))

			// Check if the timestamp has already been registered
			if _, alreadyAlerted := lastAlertTimestamps[timestamp]; alreadyAlerted {
				logMessage(cfg, 2, fmt.Sprintf("Alert already sent for timestamp %s, skipping...", timestamp))
				continue
			}

			// Register new timestamp in lastAlertTimestamps
			logMessage(cfg, 2, fmt.Sprintf("Registering alert for timestamp %s", timestamp))
			lastAlertTimestamps[timestamp] = struct{}{}

			// Generate and send alert
			subject := fmt.Sprintf("Mass Disconnection Alert: %d extensions disconnected at %s", len(extensions), timestamp)
			message := fmt.Sprintf("Mass disconnection detected at %s:\nTotal: %d extensions disconnected.\nExtensions: %s", timestamp, len(extensions), strings.Join(extensions, ", "))
			notification.NotifyAll(cfg, subject, message)
			logMessage(cfg, 1, fmt.Sprintf("Alert sent: %s", subject))
		}
	}

	// Check for scanner errors
	if err := scanner.Err(); err != nil {
		logMessage(cfg, 2, fmt.Sprintf("Error reading log lines: %v", err))
	}
}
