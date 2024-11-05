#!/bin/bash

# ParseWatchDog Installer Script

# Base URL for downloading binaries
BASE_URL="https://github.com/lordbasex/parsewatchdog/raw/refs/heads/main/dist"

# Installation directories and files
INSTALL_DIR="/usr/local/bin"
SERVICE_FILE="/etc/systemd/system/parsewatchdog.service"
CONFIG_FILE="/etc/parsewatchdog.conf"
LOG_FILE="/var/log/parsewatchdog.log"

# URLs for each architecture
DOWNLOAD_URL_AMD64="$BASE_URL/parsewatchdog-x86_64"
DOWNLOAD_URL_ARM64="$BASE_URL/parsewatchdog-arm64"
DOWNLOAD_URL_I386="$BASE_URL/parsewatchdog-i386"

# Stop Daemon
stop_daemon() {
    systemctl start parsewatchdog.service
}

# Function to check OS and architecture and set the appropriate download URL
check_os_and_arch() {
    # Check if the OS is Linux
    if [ "$(uname)" != "Linux" ]; then
        echo "Unsupported operating system: $(uname)"
        exit 1
    fi

    # Load additional OS information
    if [ -f "/etc/os-release" ]; then
        . /etc/os-release
        echo "Detected OS: $NAME ($ID), Version: $VERSION_ID"
    else
        echo "Warning: Unable to detect OS information (missing /etc/os-release)"
    fi

    # Detect architecture and set the download URL
    local os_arch
    os_arch=$(uname -m)
    case "$os_arch" in
        "x86_64")
            echo "Architecture detected: x86_64 (64-bit AMD/Intel)"
            DOWNLOAD_URL="$DOWNLOAD_URL_AMD64"
            ;;
        "aarch64")
            echo "Architecture detected: ARM64 (64-bit)"
            DOWNLOAD_URL="$DOWNLOAD_URL_ARM64"
            ;;
        "i386" | "i686")
            echo "Architecture detected: x86 (32-bit Intel)"
            DOWNLOAD_URL="$DOWNLOAD_URL_I386"
            ;;
        *)
            echo "Unsupported architecture: $os_arch"
            exit 1
            ;;
    esac
}

# Function to download the binary
download_binary() {
    echo "Downloading binary from $DOWNLOAD_URL..."
    curl -L -o "$INSTALL_DIR/parsewatchdog" "$DOWNLOAD_URL"
    if [ $? -ne 0 ]; then
        echo "Error downloading the binary."
        exit 1
    fi
    chmod +x "$INSTALL_DIR/parsewatchdog"
    echo "Binary downloaded and installed to $INSTALL_DIR/parsewatchdog"
}

# Function to create a default configuration file
create_config() {
    if [ ! -f "$CONFIG_FILE" ]; then
        echo "Creating default configuration at $CONFIG_FILE..."
        cat <<EOF > "$CONFIG_FILE"
[smtp]
enabled=true
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
EOF
        echo "Default configuration created."
    else
        echo "Configuration file already exists at $CONFIG_FILE, skipping creation."
    fi
}

# Function to create systemd service file
create_service() {
    echo "Creating systemd service file..."
    cat <<EOF > "$SERVICE_FILE"
[Unit]
Description=ParseWatchDog Service
Requires=network.target
After=network.target network-online.target

[Service]
Type=simple
ExecStart=$INSTALL_DIR/parsewatchdog
Restart=on-failure
StandardOutput=append:$LOG_FILE
StandardError=append:$LOG_FILE

[Install]
WantedBy=multi-user.target
EOF

    # Reload systemd, enable and start the service
    echo "Enabling and starting the ParseWatchDog service..."
    systemctl daemon-reload
    systemctl enable parsewatchdog.service
    systemctl start parsewatchdog.service

    # Check the status of the service
    systemctl status parsewatchdog.service --no-pager
    echo "ParseWatchDog service installed and started successfully."
}

# Main installation function
main() {
    # Ensure script is run as root
    if [ "$EUID" -ne 0 ]; then
        echo "Please run as root."
        exit 1
    fi

    stop_daemon
    check_os_and_arch
    download_binary
    create_config
    create_service
}

# Execute main function
main
