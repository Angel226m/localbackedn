#!/bin/bash
# Script para automatizar los push al repositorio GitHub
# Fecha: 2025-05-19
# Usuario: Angel226m

# Directorio principal
APP_DIR="/opt/sistema-tours"

# Configuración para el push automático
GIT_USER="Angel226m"
GIT_EMAIL="angel226m@sistema-tours.com"
REPO_URL="https://github.com/Angel226m/subidadevops.git"

# Colores para mensajes
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # Sin Color

echo -e "${YELLOW}Configurando auto-push a GitHub${NC}"
echo "Repositorio: $REPO_URL"
echo "Usuario: $GIT_USER"
echo "Fecha: 2025-05-19 04:28:04"

# Verificar si estamos en el directorio correcto
if [ ! -d "$APP_DIR/.git" ]; then
    echo -e "${RED}Error: No se encuentra el repositorio Git en $APP_DIR${NC}"
    exit 1
fi

# Ir al directorio de la aplicación
cd "$APP_DIR"

# Configurar credenciales de Git si es necesario
git config user.name "$GIT_USER"
git config user.email "$GIT_EMAIL"

# Verificar si hay cambios que hacer commit
if [[ -z $(git status -s) ]]; then
    echo -e "${GREEN}No hay cambios que subir a GitHub${NC}"
    exit 0
fi

# Preguntar por un mensaje de commit
echo -e "${YELLOW}Se han detectado cambios. Por favor, proporciona un mensaje para el commit:${NC}"
read -p "Mensaje (o presiona Enter para mensaje por defecto): " COMMIT_MSG

# Si no se proporciona mensaje, usar uno por defecto
if [ -z "$COMMIT_MSG" ]; then
    COMMIT_MSG="Actualización automática - $(date '+%Y-%m-%d %H:%M:%S')"
fi

# Añadir todos los cambios
git add .

# Realizar commit
git commit -m "$COMMIT_MSG"

# Intentar hacer push
echo -e "${YELLOW}Haciendo push a GitHub...${NC}"
if git push origin main; then
    echo -e "${GREEN}¡Push a GitHub exitoso!${NC}"
else
    echo -e "${RED}Error al hacer push. Puede que necesites proporcionar credenciales.${NC}"
    # Opción para usar token personal
    echo -e "${YELLOW}¿Quieres intentar usar un token personal? (s/n)${NC}"
    read -p "Respuesta: " USE_TOKEN
    
    if [[ "$USE_TOKEN" == "s" || "$USE_TOKEN" == "S" ]]; then
        read -p "Ingresa tu token personal de GitHub: " GH_TOKEN
        # Usar el token para hacer push
        git remote set-url origin https://$GIT_USER:$GH_TOKEN@github.com/Angel226m/subidadevops.git
        if git push origin main; then
            echo -e "${GREEN}¡Push con token exitoso!${NC}"
        else
            echo -e "${RED}Falló el push incluso con token. Revisa tus credenciales.${NC}"
        fi
    fi
fi