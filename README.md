# ParseWatchDog: The **SysAdmin's** best friend, always on guard.

ParseWatchDog is a log monitoring tool specifically designed to detect and alert on significant events in **Asterisk** systems, such as mass disconnection alerts. Currently, it focuses on monitoring SIP and PJSIP logs, specifically checking for "Unreachable" events, which indicate the disconnection of devices or users in the telephony network.

When ParseWatchDog detects a mass disconnection event, it generates alerts that can be sent through multiple notification channels, including email, Telegram, an API endpoint, and RabbitMQ. This multi-channel alerting capability ensures that responsible teams are immediately notified through the most convenient means.

Although its initial purpose is to monitor disconnection events in **Asterisk**, ParseWatchDog is designed with flexibility, allowing it to expand its capabilities in the future to monitor other events or services, adapting to the evolving needs of the telecommunications environment and ensuring comprehensive system supervision.

## Features

- Monitors specified log file for disconnection patterns
- Supports email, Telegram, API, RabbitMQ and Slack notifications
- Customizable configuration file
- Adjustable debug levels for granular logging

## Installation

1. Clone this repository:
    ```bash
    git clone https://github.com/lordbasex/parsewatchdog.git
    ```

2. Navigate into the project directory:
    ```bash
    cd parsewatchdog
    ```

3. Install the necessary dependencies:
    ```bash
    go mod tidy
    ```

4. Ensure RabbitMQ and other services used for notifications are properly configured.

## Configuration

The configuration file (`/etc/parsewatchdog.conf`) controls ParseWatchdog's behavior and notification settings.

### Example Configuration

```ini
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
enabled = false
webhook_url = https://hooks.slack.com/services/XXXXXXXXXXX/XXXXXXXXXXXX/XXXXXXXXXXXXXXXXXXXXXXXX

[debug]
debug_level=1 # Levels: 0 = no logs, 1 = only critical logs, 2 = all logs
```

## Usage

To run ParseWatchdog:

```bash
make build
make run
```

For production, you may want to compile the binary:

```bash
./dist/parsewatchdog-x86_64
```	

### Create Service (systemctl)

```bash
cp -fra ./dist/parsewatchdog-x86_64 /usr/local/bin/parsewatchdog-x86_64
chmod 777 /usr/local/bin/parsewatchdog-x86_64
```

```bash
cat > /etc/systemd/system/parsewatchdog.service <<ENDLINE
[Unit]
Description=ParseWatchDog
Requires=network.target
After=network.target network-online.target

[Service]
Type=simple
User=root
Group=root
WorkingDirectory=/usr/local/bin/
ExecStart=/usr/local/bin/parsewatchdog-x86_64
Restart=on-failure

# Redirige logs de salida estándar a un archivo
StandardOutput=append:/var/log/parsewatchdog.log
StandardError=append:/var/log/parsewatchdog.log

[Install]
WantedBy=multi-user.target
ENDLINE
```

```bash
systemctl enable parsewatchdog.service 
systemctl start parsewatchdog.service
systemctl status parsewatchdog.service 
```

 
## Debug Levels

* 0: No logs displayed
* 1: Only disconnection event logs
* 2: All logs for full debugging information

## Test

To generate simulated disconnection events in the Asterisk log file, run the following script:

```
bash test.sh
````


## Notification Channels

ParseWatchdog can send notifications via the following channels:

## Email
Sends an HTML email to specified recipients when a mass disconnection event is detected.

## Telegram
Sends a message to a specified Telegram chat. The Telegram bot token and chat ID are configured in the configuration file.

## API
Sends a JSON payload to a specified API endpoint.

## RabbitMQ
Publishes a JSON message to a RabbitMQ queue. The URI for RabbitMQ is configured in the configuration file.

License
This project is licensed under the MIT License. ""

Saving the content to a README.md file
