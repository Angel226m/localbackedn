#!/bin/bash
# Script para configurar y sincronizar automáticamente con GitHub
# Fecha: 2025-05-19
# Usuario: Angel226m

# Directorio principal
APP_DIR="/opt/sistema-tours"

# Colores para mensajes
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # Sin Color

echo -e "${YELLOW}Configurando sincronización con GitHub para Sistema-Tours${NC}"
echo "Repositorio: https://github.com/Angel226m/subidadevops.git"
echo "Fecha: 2025-05-19 04:28:04"

# Verificar si git está instalado
if ! command -v git &> /dev/null; then
    echo -e "${RED}Git no está instalado. Instalando...${NC}"
    apt update && apt install -y git
fi

# Verificar si el directorio existe, si no, crearlo
if [ ! -d "$APP_DIR" ]; then
    echo -e "${YELLOW}Creando directorio principal...${NC}"
    mkdir -p "$APP_DIR"
fi

# Configurar git si no existe el repositorio
if [ ! -d "$APP_DIR/.git" ]; then
    echo -e "${YELLOW}Inicializando repositorio Git...${NC}"
    cd "$APP_DIR"
    git init
    git remote add origin https://github.com/Angel226m/subidadevops.git
else
    echo -e "${GREEN}Repositorio Git ya está configurado${NC}"
    cd "$APP_DIR"
fi

# Configurar usuario git global
git config --global user.name "Angel226m"
git config --global user.email "angel226m@sistema-tours.com"

# Configurar GitHub CLI si está disponible
if command -v gh &> /dev/null; then
    echo -e "${YELLOW}Configurando GitHub CLI...${NC}"
    # Si tienes un token de GitHub, puedes configurarlo así:
    # echo "$GITHUB_TOKEN" | gh auth login --with-token
else
    echo -e "${YELLOW}GitHub CLI no está instalado. Usando Git directamente.${NC}"
fi

# Configurar webhook para recibir actualizaciones automáticas
WEBHOOK_DIR="/opt/github-webhook"
if [ ! -d "$WEBHOOK_DIR" ]; then
    echo -e "${YELLOW}Configurando webhook de GitHub...${NC}"
    mkdir -p "$WEBHOOK_DIR"
    
    # Instalar dependencias para el webhook (nodejs)
    apt update && apt install -y nodejs npm
    
    # Crear un pequeño servidor webhook
    cat > "$WEBHOOK_DIR/webhook.js" << 'EOL'
const http = require('http');
const crypto = require('crypto');
const { exec } = require('child_process');

const SECRET = 'CAMBIAR_ESTO_POR_UN_SECRETO_SEGURO';
const PORT = 9000;
const REPO_PATH = '/opt/sistema-tours';

http.createServer((req, res) => {
    if (req.method === 'POST' && req.url === '/webhook') {
        let body = '';
        req.on('data', chunk => {
            body += chunk.toString();
        });
        
        req.on('end', () => {
            const signature = req.headers['x-hub-signature'];
            
            if (!signature) {
                console.log('No se recibió la firma');
                return res.end('No signature');
            }
            
            const hmac = crypto.createHmac('sha1', SECRET);
            const digest = 'sha1=' + hmac.update(body).digest('hex');
            
            if (signature !== digest) {
                console.log('Firma no válida');
                return res.end('Invalid signature');
            }
            
            const payload = JSON.parse(body);
            if (payload.ref === 'refs/heads/main') {
                console.log('Recibido push a rama main, actualizando...');
                exec(`cd ${REPO_PATH} && git pull && docker-compose down && docker-compose up -d`, 
                    (error, stdout, stderr) => {
                        if (error) {
                            console.error(`Error en el deploy: ${error}`);
                            return;
                        }
                        console.log(`Salida: ${stdout}`);
                    }
                );
                res.end('OK');
            } else {
                res.end('No action needed');
            }
        });
    } else {
        res.end('Not a POST request or wrong endpoint');
    }
}).listen(PORT, () => {
    console.log(`Webhook escuchando en el puerto ${PORT}`);
});
EOL
    
    # Crear servicio systemd para el webhook
    cat > "/etc/systemd/system/github-webhook.service" << EOL
[Unit]
Description=GitHub Webhook para Sistema-Tours
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=${WEBHOOK_DIR}
ExecStart=/usr/bin/node ${WEBHOOK_DIR}/webhook.js
Restart=always

[Install]
WantedBy=multi-user.target
EOL

    # Activar y arrancar el servicio
    systemctl daemon-reload
    systemctl enable github-webhook
    systemctl start github-webhook
    
    echo -e "${GREEN}Webhook configurado y en funcionamiento en el puerto 9000${NC}"
    echo -e "${YELLOW}NOTA: Agrega un webhook en GitHub con la URL http://tu-ip:9000/webhook${NC}"
    echo -e "${YELLOW}      y configura el secreto en el archivo ${WEBHOOK_DIR}/webhook.js${NC}"
fi

# Sincronizar con el repositorio
echo -e "${YELLOW}Sincronizando con GitHub...${NC}"
git pull origin main

echo -e "${GREEN}¡Sincronización con GitHub completada!${NC}"
echo "Repositorio configurado en: $APP_DIR"
echo "Webhook escuchando en: http://tu-ip:9000/webhook"
echo ""
echo -e "${YELLOW}Recuerda añadir los secrets necesarios en tu repositorio de GitHub:${NC}"
echo "- DOCKER_HUB_USERNAME: Tu usuario de Docker Hub"
echo "- DOCKER_HUB_TOKEN: Tu token de acceso a Docker Hub"
echo "- HETZNER_HOST: La dirección IP de tu VPS"
echo "- HETZNER_USERNAME: Usuario SSH para acceder al VPS"
echo "- HETZNER_SSH_KEY: Clave SSH privada para acceder al VPS"