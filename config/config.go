package config

import (
	"gopkg.in/ini.v1"
)

// Version information (will be populated at build time using ldflags)
var (
	Version            string
	DaemonGitBuild     string
	DaemonGitBuildDate string
)

type SMTPConfig struct {
	Enabled    bool
	Host       string
	Port       int
	User       string
	Pass       string
	Recipients []string
}

type TelegramConfig struct {
	Enabled bool
	Token   string
	ChatID  int64
}

type APIConfig struct {
	Enabled  bool
	Endpoint string
	APIKey   string
}

type DebugConfig struct {
	DebugLevel int
}

type RabbitMQConfig struct {
	Enabled  bool
	Type     string
	User     string
	Password string
	IP       string
	Port     int
	Queue    string
}

type SlackConfig struct {
	Enabled    bool
	WebhookURL string
}

type Config struct {
	SMTP     SMTPConfig
	Telegram TelegramConfig
	API      APIConfig
	Debug    DebugConfig
	RabbitMQ RabbitMQConfig
	Slack    SlackConfig
}

// LoadConfig carga el archivo parsewatchdog.conf
func LoadConfig(filePath string) (*Config, error) {
	cfg, err := ini.Load(filePath)
	if err != nil {
		return nil, err
	}

	config := &Config{}

	// Leer configuración de SMTP
	smtpSection := cfg.Section("smtp")
	config.SMTP.Enabled = smtpSection.Key("enabled").MustBool(false)
	config.SMTP.Host = smtpSection.Key("host").String()
	config.SMTP.Port = smtpSection.Key("port").MustInt(587)
	config.SMTP.User = smtpSection.Key("user").String()
	config.SMTP.Pass = smtpSection.Key("pass").String()
	config.SMTP.Recipients = smtpSection.Key("recipients").Strings(",")

	// Leer configuración de Telegram
	telegramSection := cfg.Section("telegram")
	config.Telegram.Enabled = telegramSection.Key("enabled").MustBool(false)
	config.Telegram.Token = telegramSection.Key("token").String()
	config.Telegram.ChatID = telegramSection.Key("chat_id").MustInt64(0)

	// Leer configuración de API
	apiSection := cfg.Section("api")
	config.API.Enabled = apiSection.Key("enabled").MustBool(false)
	config.API.Endpoint = apiSection.Key("endpoint").String()
	config.API.APIKey = apiSection.Key("api_key").String()

	// Leer configuración de Debug
	debugSection := cfg.Section("debug")
	config.Debug.DebugLevel = debugSection.Key("debug_level").MustInt(0)

	// Leer configuración de RabbitMQ
	rabbitSection := cfg.Section("rabbitmq")
	config.RabbitMQ.Enabled = rabbitSection.Key("enabled").MustBool(false)
	config.RabbitMQ.Type = rabbitSection.Key("type").String()
	config.RabbitMQ.User = rabbitSection.Key("user").String()
	config.RabbitMQ.Password = rabbitSection.Key("password").String()
	config.RabbitMQ.IP = rabbitSection.Key("ip").String()
	config.RabbitMQ.Port = rabbitSection.Key("port").MustInt(5672)
	config.RabbitMQ.Queue = rabbitSection.Key("queue").String()

	// Leer configuración de Slack
	slackSection := cfg.Section("slack")
	config.Slack.Enabled = slackSection.Key("enabled").MustBool(false)
	config.Slack.WebhookURL = slackSection.Key("webhook_url").String()

	return config, nil
}
