#!/bin/bash

# Archivo de log destino
LOG_FILE="/var/log/asterisk/full"  # Cambia si necesitas un archivo de prueba

# Limpia o crea el archivo de log para comenzar de cero
echo "Generando eventos de desconexión en $LOG_FILE..."
echo "" >> "$LOG_FILE"

# Número de dispositivos en el evento de desconexión masiva
NUM_DEVICES_MASSIVE=20

# Generar un evento de desconexión único con la fecha actual
current_time=$(date +"%Y-%m-%d %H:%M:%S")
unique_device=$((RANDOM % 9000 + 1000))
echo "[$current_time] VERBOSE[9014] res_pjsip/pjsip_configuration.c: Endpoint $unique_device is now Unreachable" >> "$LOG_FILE"

# Generar eventos de desconexión masiva con la fecha de hoy en un mismo segundo
massive_timestamp=$(date +"%Y-%m-%d %H:%M:27")  # Ejemplo: todos en el segundo 27
for ((i=1; i<=NUM_DEVICES_MASSIVE; i++)); do
    endpoint_id=$((1100 + i))
    echo "[$massive_timestamp] VERBOSE[18528] res_pjsip/pjsip_configuration.c: Endpoint $endpoint_id is now Unreachable" >> "$LOG_FILE"
done

# Generar eventos de desconexión masiva con una fecha pasada
past_timestamp="2024-11-03 13:05:06"
past_device=8864
echo "[$past_timestamp] VERBOSE[9014] res_pjsip/pjsip_configuration.c: Endpoint $past_device is now Unreachable" >> "$LOG_FILE"

echo "Eventos generados en $LOG_FILE."

