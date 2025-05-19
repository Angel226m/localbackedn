#!/bin/bash
# Script para configurar claves SSH para GitHub Actions
# Fecha: 2025-05-19
# Usuario: Angel226m

# Colores para mensajes
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # Sin Color

echo -e "${YELLOW}Configurando clave SSH para GitHub Actions${NC}"
echo "Fecha: 2025-05-19 04:28:04"
echo "Usuario: Angel226m"

# Directorio para las claves
SSH_DIR="/home/deployer/.ssh"
KEY_NAME="github_actions_deploy_key"

# Crear directorio si no existe
if [ ! -d "$SSH_DIR" ]; then
    mkdir -p "$SSH_DIR"
    chmod 700 "$SSH_DIR"
fi

# Generar par de claves SSH
echo -e "${YELLOW}Generando par de claves SSH...${NC}"
ssh-keygen -t ed25519 -C "github-actions-deploy@sistema-tours.com" -f "$SSH_DIR/$KEY_NAME" -N ""

# Configurar authorized_keys
echo -e "${YELLOW}Configurando clave en authorized_keys...${NC}"
cat "$SSH_DIR/$KEY_NAME.pub" >> "$SSH_DIR/authorized_keys"
chmod 600 "$SSH_DIR/authorized_keys"

# Ajustar permisos
chown -R deployer:deployer "$SSH_DIR"

# Instrucciones
echo -e "${GREEN}¡Clave SSH generada exitosamente!${NC}"
echo ""
echo -e "${YELLOW}Instrucciones para configurar GitHub Actions:${NC}"
echo "1. Añade la siguiente clave privada como secreto en tu repositorio GitHub:"
echo "   Nombre del secreto: HETZNER_SSH_KEY"
echo ""
echo "   Contenido del secreto (copia todo):"
echo "   ----------------------------------------"
cat "$SSH_DIR/$KEY_NAME"
echo "   ----------------------------------------"
echo ""
echo "2. También añade estos otros secretos:"
echo "   - HETZNER_HOST: La dirección IP de tu VPS"
echo "   - HETZNER_USERNAME: deployer"
echo ""
echo -e "${RED}¡IMPORTANTE! Guarda esta información de manera segura y luego borra este archivo${NC}"
echo "La clave privada ha sido guardada en: $SSH_DIR/$KEY_NAME"
echo "La clave pública ha sido guardada en: $SSH_DIR/$KEY_NAME.pub"