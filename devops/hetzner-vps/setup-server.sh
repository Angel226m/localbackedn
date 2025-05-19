#!/bin/bash 
# Script de configuración inicial para Hetzner VPS 
# Fecha: 2025-05-19 
# Autor: Sistema-Tours 
 
echo "Iniciando configuración del servidor..." 
 
# Actualizar el sistema 
apt update && apt upgrade -y 
 
# Instalar paquetes esenciales 
apt install -y apt-transport-https ca-certificates curl software-properties-common gnupg lsb-release git fail2ban unzip jq build-essential python3 python3-pip 
 
# Configurar zona horaria 
timedatectl set-timezone America/Mexico_City 
 
# Configurar Fail2Ban para protección de SSH 
cat > /etc/fail2ban/jail.local << EOL 
[sshd] 
enabled = true 
port = ssh 
filter = sshd 
logpath = /var/log/auth.log 
maxretry = 5 
bantime = 3600 
EOL 
 
# Reiniciar Fail2Ban 
systemctl restart fail2ban 
 
# Configurar firewall 
apt install -y ufw 
ufw default deny incoming 
ufw default allow outgoing 
ufw allow ssh 
ufw allow http 
ufw allow https 
ufw allow 9090/tcp 
ufw allow 3000/tcp 
ufw allow 8080/tcp 
ufw allow 8000/tcp 
