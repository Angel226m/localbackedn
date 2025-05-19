#!/bin/bash
# Script para probar la integración con GitHub
# Fecha: 2025-05-19
# Usuario: Angel226m

# Colores para mensajes
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # Sin Color

echo -e "${YELLOW}Probando integración con GitHub${NC}"
echo "Fecha: 2025-05-19 04:28:04"
echo "Usuario: Angel226m"
echo "Repositorio: https://github.com/Angel226m/subidadevops.git"

# Verificar si git está instalado
if ! command -v git &> /dev/null; then
    echo -e "${RED}Error: Git no está instalado${NC}"
    exit 1
fi

# Verificar conexión a GitHub
echo -e "${YELLOW}Verificando conexión a GitHub...${NC}"
if curl -s https://api.github.com > /dev/null; then
    echo -e "${GREEN}Conexión a GitHub OK${NC}"
else
    echo -e "${RED}No se puede conectar a GitHub${NC}"
    exit 1
fi

# Verificar si el repositorio es accesible
echo -e "${YELLOW}Verificando acceso al repositorio...${NC}"
if git ls-remote https://github.com/Angel226m/subidadevops.git > /dev/null 2>&1; then
    echo -e "${GREEN}Repositorio accesible${NC}"
else
    echo -e "${RED}No se puede acceder al repositorio${NC}"
    echo -e "${YELLOW}Comprueba que el repositorio existe y es público, o que tienes acceso a él.${NC}"
    exit 1
fi

# Crear directorio temporal para pruebas
TMP_DIR=$(mktemp -d)
echo -e "${YELLOW}Creando clon temporal en $TMP_DIR${NC}"

# Clonar el repositorio
if git clone https://github.com/Angel226m/subidadevops.git "$TMP_DIR/repo"; then
    echo -e "${GREEN}Repositorio clonado exitosamente${NC}"
    
    # Crear un archivo de prueba
    echo "Test de integración con GitHub - $(date)" > "$TMP_DIR/repo/test-github-integration.txt"
    
    # Configurar git
    cd "$TMP_DIR/repo"
    git config user.name "Test Integration"
    git config user.email "test@sistema-tours.com"
    
    # Intentar hacer commit y push
    echo -e "${YELLOW}Intentando push de prueba (esto puede fallar si no tienes permisos)...${NC}"
    git add test-github-integration.txt
    git commit -m "Test de integración desde Hetzner VPS"
    
    # Esta parte probablemente fallará si no hay credenciales configuradas, lo cual es esperado
    if git push origin main; then
        echo -e "${GREEN}¡Push exitoso! La integración está completamente funcional.${NC}"
    else
        echo -e "${YELLOW}El push falló, lo cual es normal si no tienes credenciales configuradas.${NC}"
        echo -e "${YELLOW}Para configurar credenciales, ejecuta: devops/scripts/github-sync.sh${NC}"
    fi
else
    echo -e "${RED}Error al clonar el repositorio${NC}"
fi

# Limpiar
echo -e "${YELLOW}Limpiando directorio temporal...${NC}"
rm -rf "$TMP_DIR"

echo -e "${GREEN}Prueba de integración con GitHub completada${NC}"
echo "Para configurar completamente la integración automática, sigue estos pasos:"
echo "1. Ejecuta devops/scripts/github-sync.sh para configurar la sincronización"
echo "2. Ejecuta devops/scripts/setup-github-deploy-key.sh para configurar claves SSH"
echo "3. Configura los secretos en el repositorio GitHub según las instrucciones anteriores"