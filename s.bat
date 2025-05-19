@echo off
echo =====================================================================
echo CONFIGURACION COMPLETA DEVOPS PARA HETZNER VPS - SISTEMA-TOURS
echo =====================================================================
echo Fecha: 2025-05-19
echo Usuario: Angel226m
echo Repositorio: https://github.com/Angel226m/subidadevops.git
echo.
echo Creando estructura de archivos y configuraciones...
echo.

REM Crear directorios principales
mkdir devops
mkdir devops\hetzner-vps
mkdir devops\scripts
mkdir devops\terraform
mkdir devops\jenkins
mkdir devops\monitoring
mkdir devops\kubernetes
mkdir devops\database
mkdir .github\workflows
mkdir backend\migrations

REM ==================================================
REM CONFIGURACIÓN INICIAL DEL SERVIDOR HETZNER
REM ==================================================

echo #!/bin/bash > devops\hetzner-vps\setup-server.sh
echo # Script de configuración inicial para Hetzner VPS >> devops\hetzner-vps\setup-server.sh
echo # Fecha: 2025-05-19 >> devops\hetzner-vps\setup-server.sh
echo # Autor: Sistema-Tours >> devops\hetzner-vps\setup-server.sh
echo. >> devops\hetzner-vps\setup-server.sh
echo echo "Iniciando configuración del servidor..." >> devops\hetzner-vps\setup-server.sh
echo. >> devops\hetzner-vps\setup-server.sh
echo # Actualizar el sistema >> devops\hetzner-vps\setup-server.sh
echo apt update ^&^& apt upgrade -y >> devops\hetzner-vps\setup-server.sh
echo. >> devops\hetzner-vps\setup-server.sh
echo # Instalar paquetes esenciales >> devops\hetzner-vps\setup-server.sh
echo apt install -y apt-transport-https ca-certificates curl software-properties-common gnupg lsb-release git fail2ban unzip jq build-essential python3 python3-pip >> devops\hetzner-vps\setup-server.sh
echo. >> devops\hetzner-vps\setup-server.sh
echo # Configurar zona horaria >> devops\hetzner-vps\setup-server.sh
echo timedatectl set-timezone America/Mexico_City >> devops\hetzner-vps\setup-server.sh
echo. >> devops\hetzner-vps\setup-server.sh
echo # Configurar Fail2Ban para protección de SSH >> devops\hetzner-vps\setup-server.sh
echo cat ^> /etc/fail2ban/jail.local ^<^< EOL >> devops\hetzner-vps\setup-server.sh
echo [sshd] >> devops\hetzner-vps\setup-server.sh
echo enabled = true >> devops\hetzner-vps\setup-server.sh
echo port = ssh >> devops\hetzner-vps\setup-server.sh
echo filter = sshd >> devops\hetzner-vps\setup-server.sh
echo logpath = /var/log/auth.log >> devops\hetzner-vps\setup-server.sh
echo maxretry = 5 >> devops\hetzner-vps\setup-server.sh
echo bantime = 3600 >> devops\hetzner-vps\setup-server.sh
echo EOL >> devops\hetzner-vps\setup-server.sh
echo. >> devops\hetzner-vps\setup-server.sh
echo # Reiniciar Fail2Ban >> devops\hetzner-vps\setup-server.sh
echo systemctl restart fail2ban >> devops\hetzner-vps\setup-server.sh
echo. >> devops\hetzner-vps\setup-server.sh
echo # Configurar firewall >> devops\hetzner-vps\setup-server.sh
echo apt install -y ufw >> devops\hetzner-vps\setup-server.sh
echo ufw default deny incoming >> devops\hetzner-vps\setup-server.sh
echo ufw default allow outgoing >> devops\hetzner-vps\setup-server.sh
echo ufw allow ssh >> devops\hetzner-vps\setup-server.sh
echo ufw allow http >> devops\hetzner-vps\setup-server.sh
echo ufw allow https >> devops\hetzner-vps\setup-server.sh
echo ufw allow 9090/tcp >> devops\hetzner-vps\setup-server.sh
echo ufw allow 3000/tcp >> devops\hetzner-vps\setup-server.sh
echo ufw allow 8080/tcp >> devops\hetzner-vps\setup-server.sh
echo ufw allow 8000/tcp >> devops\hetzner-vps\setup-server.sh
echo echo "y" | ufw enable >> devops\hetzner-vps\setup-server.sh
echo. >> devops\hetzner-vps\setup-server.sh
echo # Crear usuario de despliegue >> devops\hetzner-vps\setup-server.sh
echo useradd -m -s /bin/bash deployer >> devops\hetzner-vps\setup-server.sh
echo echo "deployer:StrongPassword2025" | chpasswd >> devops\hetzner-vps\setup-server.sh
echo mkdir -p /home/deployer/.ssh >> devops\hetzner-vps\setup-server.sh
echo touch /home/deployer/.ssh/authorized_keys >> devops\hetzner-vps\setup-server.sh
echo chmod 700 /home/deployer/.ssh >> devops\hetzner-vps\setup-server.sh
echo chmod 600 /home/deployer/.ssh/authorized_keys >> devops\hetzner-vps\setup-server.sh
echo chown -R deployer:deployer /home/deployer/.ssh >> devops\hetzner-vps\setup-server.sh
echo. >> devops\hetzner-vps\setup-server.sh
echo # Añadir usuario al grupo sudo >> devops\hetzner-vps\setup-server.sh
echo usermod -aG sudo deployer >> devops\hetzner-vps\setup-server.sh
echo. >> devops\hetzner-vps\setup-server.sh
echo # Desactivar login con password para root (solo SSH key) >> devops\hetzner-vps\setup-server.sh
echo sed -i 's/PermitRootLogin yes/PermitRootLogin prohibit-password/g' /etc/ssh/sshd_config >> devops\hetzner-vps\setup-server.sh
echo systemctl restart ssh >> devops\hetzner-vps\setup-server.sh
echo. >> devops\hetzner-vps\setup-server.sh
echo echo "Configuración básica del servidor completada!" >> devops\hetzner-vps\setup-server.sh

REM ==================================================
REM INSTALACIÓN DE DOCKER Y DOCKER-COMPOSE
REM ==================================================

echo #!/bin/bash > devops\hetzner-vps\install-docker.sh
echo # Instalación de Docker y Docker Compose >> devops\hetzner-vps\install-docker.sh
echo. >> devops\hetzner-vps\install-docker.sh
echo # Desinstalar versiones antiguas si existen >> devops\hetzner-vps\install-docker.sh
echo apt remove -y docker docker-engine docker.io containerd runc >> devops\hetzner-vps\install-docker.sh
echo. >> devops\hetzner-vps\install-docker.sh
echo # Añadir repositorio oficial de Docker >> devops\hetzner-vps\install-docker.sh
echo curl -fsSL https://download.docker.com/linux/ubuntu/gpg | gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg >> devops\hetzner-vps\install-docker.sh
echo echo "deb [arch=amd64 signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null >> devops\hetzner-vps\install-docker.sh
echo. >> devops\hetzner-vps\install-docker.sh
echo # Instalar Docker Engine >> devops\hetzner-vps\install-docker.sh
echo apt update >> devops\hetzner-vps\install-docker.sh
echo apt install -y docker-ce docker-ce-cli containerd.io >> devops\hetzner-vps\install-docker.sh
echo. >> devops\hetzner-vps\install-docker.sh
echo # Instalar Docker Compose >> devops\hetzner-vps\install-docker.sh
echo curl -L "https://github.com/docker/compose/releases/download/v2.20.3/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose >> devops\hetzner-vps\install-docker.sh
echo chmod +x /usr/local/bin/docker-compose >> devops\hetzner-vps\install-docker.sh
echo. >> devops\hetzner-vps\install-docker.sh
echo # Añadir usuario al grupo docker >> devops\hetzner-vps\install-docker.sh
echo usermod -aG docker deployer >> devops\hetzner-vps\install-docker.sh
echo. >> devops\hetzner-vps\install-docker.sh
echo # Iniciar y habilitar Docker >> devops\hetzner-vps\install-docker.sh
echo systemctl start docker >> devops\hetzner-vps\install-docker.sh
echo systemctl enable docker >> devops\hetzner-vps\install-docker.sh
echo. >> devops\hetzner-vps\install-docker.sh
echo # Verificar instalación >> devops\hetzner-vps\install-docker.sh
echo docker --version >> devops\hetzner-vps\install-docker.sh
echo docker-compose --version >> devops\hetzner-vps\install-docker.sh
echo. >> devops\hetzner-vps\install-docker.sh
echo echo "Docker y Docker Compose instalados correctamente!" >> devops\hetzner-vps\install-docker.sh

REM ==================================================
REM INSTALACIÓN Y CONFIGURACIÓN DE POSTGRESQL
REM ==================================================

echo #!/bin/bash > devops\database\setup-postgres.sh
echo # Instalación y configuración de PostgreSQL >> devops\database\setup-postgres.sh
echo. >> devops\database\setup-postgres.sh
echo # Crear directorio para datos persistentes >> devops\database\setup-postgres.sh
echo mkdir -p /opt/postgres-data >> devops\database\setup-postgres.sh
echo. >> devops\database\setup-postgres.sh
echo # Crear archivo docker-compose para PostgreSQL >> devops\database\setup-postgres.sh
echo cat ^> /opt/postgres-docker-compose.yml ^<^< EOL >> devops\database\setup-postgres.sh
echo version: '3.8' >> devops\database\setup-postgres.sh
echo. >> devops\database\setup-postgres.sh
echo services: >> devops\database\setup-postgres.sh
echo   postgres: >> devops\database\setup-postgres.sh
echo     image: postgres:14 >> devops\database\setup-postgres.sh
echo     container_name: postgres-db >> devops\database\setup-postgres.sh
echo     restart: always >> devops\database\setup-postgres.sh
echo     environment: >> devops\database\setup-postgres.sh
echo       POSTGRES_PASSWORD: StrongDBPassword2025 >> devops\database\setup-postgres.sh
echo       POSTGRES_USER: sistema_tours >> devops\database\setup-postgres.sh
echo       POSTGRES_DB: sistema_tours_db >> devops\database\setup-postgres.sh
echo     ports: >> devops\database\setup-postgres.sh
echo       - "5432:5432" >> devops\database\setup-postgres.sh
echo     volumes: >> devops\database\setup-postgres.sh
echo       - /opt/postgres-data:/var/lib/postgresql/data >> devops\database\setup-postgres.sh
echo     networks: >> devops\database\setup-postgres.sh
echo       - backend-network >> devops\database\setup-postgres.sh
echo. >> devops\database\setup-postgres.sh
echo   pgadmin: >> devops\database\setup-postgres.sh
echo     image: dpage/pgadmin4 >> devops\database\setup-postgres.sh
echo     container_name: pgadmin >> devops\database\setup-postgres.sh
echo     restart: always >> devops\database\setup-postgres.sh
echo     environment: >> devops\database\setup-postgres.sh
echo       PGADMIN_DEFAULT_EMAIL: admin@sistema-tours.com >> devops\database\setup-postgres.sh
echo       PGADMIN_DEFAULT_PASSWORD: PgAdminPassword2025 >> devops\database\setup-postgres.sh
echo     ports: >> devops\database\setup-postgres.sh
echo       - "5050:80" >> devops\database\setup-postgres.sh
echo     depends_on: >> devops\database\setup-postgres.sh
echo       - postgres >> devops\database\setup-postgres.sh
echo     networks: >> devops\database\setup-postgres.sh
echo       - backend-network >> devops\database\setup-postgres.sh
echo. >> devops\database\setup-postgres.sh
echo networks: >> devops\database\setup-postgres.sh
echo   backend-network: >> devops\database\setup-postgres.sh
echo     driver: bridge >> devops\database\setup-postgres.sh
echo EOL >> devops\database\setup-postgres.sh
echo. >> devops\database\setup-postgres.sh
echo # Iniciar los contenedores de la base de datos >> devops\database\setup-postgres.sh
echo cd /opt && docker-compose -f postgres-docker-compose.yml up -d >> devops\database\setup-postgres.sh
echo. >> devops\database\setup-postgres.sh
echo # Configuración de respaldo automático (backup) >> devops\database\setup-postgres.sh
echo mkdir -p /opt/db-backups >> devops\database\setup-postgres.sh
echo. >> devops\database\setup-postgres.sh
echo # Script de respaldo diario >> devops\database\setup-postgres.sh
echo cat ^> /opt/db-backup.sh ^<^< EOL >> devops\database\setup-postgres.sh
echo #!/bin/bash >> devops\database\setup-postgres.sh
echo. >> devops\database\setup-postgres.sh
echo DATE=\$(date +%%Y-%%m-%%d-%%H%%M%%S) >> devops\database\setup-postgres.sh
echo BACKUP_DIR="/opt/db-backups" >> devops\database\setup-postgres.sh
echo. >> devops\database\setup-postgres.sh
echo # Crear respaldo >> devops\database\setup-postgres.sh
echo docker exec postgres-db pg_dump -U sistema_tours sistema_tours_db | gzip > \$BACKUP_DIR/sistema-tours-db-\$DATE.sql.gz >> devops\database\setup-postgres.sh
echo. >> devops\database\setup-postgres.sh
echo # Eliminar respaldos antiguos (más de 7 días) >> devops\database\setup-postgres.sh
echo find \$BACKUP_DIR -type f -name "*.sql.gz" -mtime +7 -delete >> devops\database\setup-postgres.sh
echo EOL >> devops\database\setup-postgres.sh
echo. >> devops\database\setup-postgres.sh
echo chmod +x /opt/db-backup.sh >> devops\database\setup-postgres.sh
echo. >> devops\database\setup-postgres.sh
echo # Programar respaldo diario a las 3 AM >> devops\database\setup-postgres.sh
echo (crontab -l 2>/dev/null; echo "0 3 * * * /opt/db-backup.sh") | crontab - >> devops\database\setup-postgres.sh
echo. >> devops\database\setup-postgres.sh
echo echo "PostgreSQL y pgAdmin instalados correctamente!" >> devops\database\setup-postgres.sh

REM ==================================================
REM CONFIGURACIÓN DE TERRAFORM
REM ==================================================

echo #!/bin/bash > devops\terraform\install-terraform.sh
echo # Instalación de Terraform >> devops\terraform\install-terraform.sh
echo. >> devops\terraform\install-terraform.sh
echo # Añadir repositorio de HashiCorp >> devops\terraform\install-terraform.sh
echo curl -fsSL https://apt.releases.hashicorp.com/gpg | apt-key add - >> devops\terraform\install-terraform.sh
echo apt-add-repository "deb [arch=amd64] https://apt.releases.hashicorp.com $(lsb_release -cs) main" >> devops\terraform\install-terraform.sh
echo. >> devops\terraform\install-terraform.sh
echo # Instalar Terraform >> devops\terraform\install-terraform.sh
echo apt update >> devops\terraform\install-terraform.sh
echo apt install -y terraform >> devops\terraform\install-terraform.sh
echo. >> devops\terraform\install-terraform.sh
echo # Verificar instalación >> devops\terraform\install-terraform.sh
echo terraform --version >> devops\terraform\install-terraform.sh
echo. >> devops\terraform\install-terraform.sh
echo # Crear directorios para configuración de Terraform >> devops\terraform\install-terraform.sh
echo mkdir -p /opt/terraform/sistema-tours >> devops\terraform\install-terraform.sh
echo. >> devops\terraform\install-terraform.sh
echo echo "Terraform instalado correctamente!" >> devops\terraform\install-terraform.sh

REM Crear archivo main.tf de Terraform
echo # Configuración de Terraform para Sistema-Tours > devops\terraform\main.tf
echo # Archivo principal de configuración >> devops\terraform\main.tf
echo. >> devops\terraform\main.tf
echo terraform { >> devops\terraform\main.tf
echo   required_providers { >> devops\terraform\main.tf
echo     hcloud = { >> devops\terraform\main.tf
echo       source = "hetznercloud/hcloud" >> devops\terraform\main.tf
echo       version = "~> 1.38.2" >> devops\terraform\main.tf
echo     } >> devops\terraform\main.tf
echo   } >> devops\terraform\main.tf
echo } >> devops\terraform\main.tf
echo. >> devops\terraform\main.tf
echo provider "hcloud" { >> devops\terraform\main.tf
echo   token = var.hcloud_token >> devops\terraform\main.tf
echo } >> devops\terraform\main.tf
echo. >> devops\terraform\main.tf
echo # Servidor VPS para Sistema-Tours >> devops\terraform\main.tf
echo resource "hcloud_server" "sistema_tours_server" { >> devops\terraform\main.tf
echo   name        = "sistema-tours" >> devops\terraform\main.tf
echo   image       = "ubuntu-22.04" >> devops\terraform\main.tf
echo   server_type = "cx31" >> devops\terraform\main.tf
echo   location    = "nbg1" >> devops\terraform\main.tf
echo   ssh_keys    = [hcloud_ssh_key.default.id] >> devops\terraform\main.tf
echo. >> devops\terraform\main.tf
echo   public_net { >> devops\terraform\main.tf
echo     ipv4_enabled = true >> devops\terraform\main.tf
echo     ipv6_enabled = true >> devops\terraform\main.tf
echo   } >> devops\terraform\main.tf
echo. >> devops\terraform\main.tf
echo   labels = { >> devops\terraform\main.tf
echo     environment = "production" >> devops\terraform\main.tf
echo     app = "sistema-tours" >> devops\terraform\main.tf
echo   } >> devops\terraform\main.tf
echo. >> devops\terraform\main.tf
echo   # Script de inicialización básica >> devops\terraform\main.tf
echo   user_data = <<-EOF >> devops\terraform\main.tf
echo               #!/bin/bash >> devops\terraform\main.tf
echo               apt update && apt upgrade -y >> devops\terraform\main.tf
echo               apt install -y git curl >> devops\terraform\main.tf
echo               EOF >> devops\terraform\main.tf
echo } >> devops\terraform\main.tf
echo. >> devops\terraform\main.tf
echo # Llave SSH para acceso al servidor >> devops\terraform\main.tf
echo resource "hcloud_ssh_key" "default" { >> devops\terraform\main.tf
echo   name       = "sistema-tours-key" >> devops\terraform\main.tf
echo   public_key = file(var.ssh_public_key_path) >> devops\terraform\main.tf
echo } >> devops\terraform\main.tf
echo. >> devops\terraform\main.tf
echo # Firewall básico >> devops\terraform\main.tf
echo resource "hcloud_firewall" "sistema_tours_firewall" { >> devops\terraform\main.tf
echo   name = "sistema-tours-firewall" >> devops\terraform\main.tf
echo. >> devops\terraform\main.tf
echo   # Regla para SSH >> devops\terraform\main.tf
echo   rule { >> devops\terraform\main.tf
echo     direction = "in" >> devops\terraform\main.tf
echo     protocol  = "tcp" >> devops\terraform\main.tf
echo     port      = "22" >> devops\terraform\main.tf
echo     source_ips = ["0.0.0.0/0", "::/0"] >> devops\terraform\main.tf
echo   } >> devops\terraform\main.tf
echo. >> devops\terraform\main.tf
echo   # Regla para HTTP >> devops\terraform\main.tf
echo   rule { >> devops\terraform\main.tf
echo     direction = "in" >> devops\terraform\main.tf
echo     protocol  = "tcp" >> devops\terraform\main.tf
echo     port      = "80" >> devops\terraform\main.tf
echo     source_ips = ["0.0.0.0/0", "::/0"] >> devops\terraform\main.tf
echo   } >> devops\terraform\main.tf
echo. >> devops\terraform\main.tf
echo   # Regla para HTTPS >> devops\terraform\main.tf
echo   rule { >> devops\terraform\main.tf
echo     direction = "in" >> devops\terraform\main.tf
echo     protocol  = "tcp" >> devops\terraform\main.tf
echo     port      = "443" >> devops\terraform\main.tf
echo     source_ips = ["0.0.0.0/0", "::/0"] >> devops\terraform\main.tf
echo   } >> devops\terraform\main.tf
echo } >> devops\terraform\main.tf
echo. >> devops\terraform\main.tf
echo # Asignar firewall al servidor >> devops\terraform\main.tf
echo resource "hcloud_firewall_attachment" "sistema_tours_firewall_attachment" { >> devops\terraform\main.tf
echo   firewall_id = hcloud_firewall.sistema_tours_firewall.id >> devops\terraform\main.tf
echo   server_ids  = [hcloud_server.sistema_tours_server.id] >> devops\terraform\main.tf
echo } >> devops\terraform\main.tf

REM Crear archivo variables.tf de Terraform
echo # Variables para la configuración de Terraform > devops\terraform\variables.tf
echo. >> devops\terraform\variables.tf
echo variable "hcloud_token" { >> devops\terraform\variables.tf
echo   description = "Token de API de Hetzner Cloud" >> devops\terraform\variables.tf
echo   type        = string >> devops\terraform\variables.tf
echo   sensitive   = true >> devops\terraform\variables.tf
echo } >> devops\terraform\variables.tf
echo. >> devops\terraform\variables.tf
echo variable "ssh_public_key_path" { >> devops\terraform\variables.tf
echo   description = "Ruta al archivo de clave pública SSH" >> devops\terraform\variables.tf
echo   type        = string >> devops\terraform\variables.tf
echo   default     = "~/.ssh/id_rsa.pub" >> devops\terraform\variables.tf
echo } >> devops\terraform\variables.tf

REM Crear archivo outputs.tf de Terraform
echo # Salidas de la configuración de Terraform > devops\terraform\outputs.tf
echo. >> devops\terraform\outputs.tf
echo output "server_ip" { >> devops\terraform\outputs.tf
echo   description = "Dirección IP pública del servidor" >> devops\terraform\outputs.tf
echo   value       = hcloud_server.sistema_tours_server.ipv4_address >> devops\terraform\outputs.tf
echo } >> devops\terraform\outputs.tf
echo. >> devops\terraform\outputs.tf
echo output "server_status" { >> devops\terraform\outputs.tf
echo   description = "Estado del servidor" >> devops\terraform\outputs.tf
echo   value       = hcloud_server.sistema_tours_server.status >> devops\terraform\outputs.tf
echo } >> devops\terraform\outputs.tf

REM ==================================================
REM CONFIGURACIÓN DE JENKINS
REM ==================================================

echo #!/bin/bash > devops\jenkins\install-jenkins.sh
echo # Instalación de Jenkins >> devops\jenkins\install-jenkins.sh
echo. >> devops\jenkins\install-jenkins.sh
echo # Crear directorio para datos persistentes >> devops\jenkins\install-jenkins.sh
echo mkdir -p /opt/jenkins-data >> devops\jenkins\install-jenkins.sh
echo. >> devops\jenkins\install-jenkins.sh
echo # Crear archivo docker-compose para Jenkins >> devops\jenkins\install-jenkins.sh
echo cat ^> /opt/jenkins-docker-compose.yml ^<^< EOL >> devops\jenkins\install-jenkins.sh
echo version: '3.8' >> devops\jenkins\install-jenkins.sh
echo. >> devops\jenkins\install-jenkins.sh
echo services: >> devops\jenkins\install-jenkins.sh
echo   jenkins: >> devops\jenkins\install-jenkins.sh
echo     image: jenkins/jenkins:lts >> devops\jenkins\install-jenkins.sh
echo     container_name: jenkins >> devops\jenkins\install-jenkins.sh
echo     restart: always >> devops\jenkins\install-jenkins.sh
echo     privileged: true >> devops\jenkins\install-jenkins.sh
echo     user: root >> devops\jenkins\install-jenkins.sh
echo     ports: >> devops\jenkins\install-jenkins.sh
echo       - "8080:8080" >> devops\jenkins\install-jenkins.sh
echo       - "50000:50000" >> devops\jenkins\install-jenkins.sh
echo     volumes: >> devops\jenkins\install-jenkins.sh
echo       - /opt/jenkins-data:/var/jenkins_home >> devops\jenkins\install-jenkins.sh
echo       - /var/run/docker.sock:/var/run/docker.sock >> devops\jenkins\install-jenkins.sh
echo     environment: >> devops\jenkins\install-jenkins.sh
echo       - JAVA_OPTS=-Djenkins.install.runSetupWizard=false >> devops\jenkins\install-jenkins.sh
echo     networks: >> devops\jenkins\install-jenkins.sh
echo       - jenkins-network >> devops\jenkins\install-jenkins.sh
echo. >> devops\jenkins\install-jenkins.sh
echo networks: >> devops\jenkins\install-jenkins.sh
echo   jenkins-network: >> devops\jenkins\install-jenkins.sh
echo     driver: bridge >> devops\jenkins\install-jenkins.sh
echo EOL >> devops\jenkins\install-jenkins.sh
echo. >> devops\jenkins\install-jenkins.sh
echo # Iniciar el contenedor de Jenkins >> devops\jenkins\install-jenkins.sh
echo cd /opt && docker-compose -f jenkins-docker-compose.yml up -d >> devops\jenkins\install-jenkins.sh
echo. >> devops\jenkins\install-jenkins.sh
echo # Esperar a que Jenkins esté listo >> devops\jenkins\install-jenkins.sh
echo echo "Esperando a que Jenkins esté en funcionamiento..." >> devops\jenkins\install-jenkins.sh
echo sleep 30 >> devops\jenkins\install-jenkins.sh
echo. >> devops\jenkins\install-jenkins.sh
echo # Mostrar la contraseña inicial de administrador >> devops\jenkins\install-jenkins.sh
echo echo "Contraseña inicial de administrador de Jenkins:" >> devops\jenkins\install-jenkins.sh
echo docker exec jenkins cat /var/jenkins_home/secrets/initialAdminPassword >> devops\jenkins\install-jenkins.sh
echo. >> devops\jenkins\install-jenkins.sh
echo echo "Jenkins instalado correctamente! Accede desde http://tu-ip:8080" >> devops\jenkins\install-jenkins.sh

REM ==================================================
REM CONFIGURACIÓN DE PROMETHEUS Y GRAFANA
REM ==================================================

echo #!/bin/bash > devops\monitoring\install-monitoring.sh
echo # Instalación de Prometheus y Grafana >> devops\monitoring\install-monitoring.sh
echo. >> devops\monitoring\install-monitoring.sh
echo # Crear directorios para datos persistentes >> devops\monitoring\install-monitoring.sh
echo mkdir -p /opt/prometheus-data /opt/grafana-data >> devops\monitoring\install-monitoring.sh
echo. >> devops\monitoring\install-monitoring.sh
echo # Crear directorio para configuración de Prometheus >> devops\monitoring\install-monitoring.sh
echo mkdir -p /opt/prometheus-config >> devops\monitoring\install-monitoring.sh
echo. >> devops\monitoring\install-monitoring.sh
echo # Configurar Prometheus >> devops\monitoring\install-monitoring.sh
echo cat ^> /opt/prometheus-config/prometheus.yml ^<^< EOL >> devops\monitoring\install-monitoring.sh
echo global: >> devops\monitoring\install-monitoring.sh
echo   scrape_interval: 15s >> devops\monitoring\install-monitoring.sh
echo   evaluation_interval: 15s >> devops\monitoring\install-monitoring.sh
echo. >> devops\monitoring\install-monitoring.sh
echo scrape_configs: >> devops\monitoring\install-monitoring.sh
echo   - job_name: 'prometheus' >> devops\monitoring\install-monitoring.sh
echo     static_configs: >> devops\monitoring\install-monitoring.sh
echo       - targets: ['localhost:9090'] >> devops\monitoring\install-monitoring.sh
echo. >> devops\monitoring\install-monitoring.sh
echo   - job_name: 'node_exporter' >> devops\monitoring\install-monitoring.sh
echo     static_configs: >> devops\monitoring\install-monitoring.sh
echo       - targets: ['node-exporter:9100'] >> devops\monitoring\install-monitoring.sh
echo. >> devops\monitoring\install-monitoring.sh
echo   - job_name: 'cadvisor' >> devops\monitoring\install-monitoring.sh
echo     static_configs: >> devops\monitoring\install-monitoring.sh
echo       - targets: ['cadvisor:8080'] >> devops\monitoring\install-monitoring.sh
echo EOL >> devops\monitoring\install-monitoring.sh
echo. >> devops\monitoring\install-monitoring.sh
echo # Crear archivo docker-compose para monitoreo >> devops\monitoring\install-monitoring.sh
echo cat ^> /opt/monitoring-docker-compose.yml ^<^< EOL >> devops\monitoring\install-monitoring.sh
echo version: '3.8' >> devops\monitoring\install-monitoring.sh
echo. >> devops\monitoring\install-monitoring.sh
echo services: >> devops\monitoring\install-monitoring.sh
echo   prometheus: >> devops\monitoring\install-monitoring.sh
echo     image: prom/prometheus >> devops\monitoring\install-monitoring.sh
echo     container_name: prometheus >> devops\monitoring\install-monitoring.sh
echo     restart: always >> devops\monitoring\install-monitoring.sh
echo     volumes: >> devops\monitoring\install-monitoring.sh
echo       - /opt/prometheus-config:/etc/prometheus >> devops\monitoring\install-monitoring.sh
echo       - /opt/prometheus-data:/prometheus >> devops\monitoring\install-monitoring.sh
echo     command: >> devops\monitoring\install-monitoring.sh
echo       - '--config.file=/etc/prometheus/prometheus.yml' >> devops\monitoring\install-monitoring.sh
echo       - '--storage.tsdb.path=/prometheus' >> devops\monitoring\install-monitoring.sh
echo       - '--web.console.libraries=/usr/share/prometheus/console_libraries' >> devops\monitoring\install-monitoring.sh
echo       - '--web.console.templates=/usr/share/prometheus/consoles' >> devops\monitoring\install-monitoring.sh
echo     ports: >> devops\monitoring\install-monitoring.sh
echo       - "9090:9090" >> devops\monitoring\install-monitoring.sh
echo     networks: >> devops\monitoring\install-monitoring.sh
echo       - monitoring-network >> devops\monitoring\install-monitoring.sh
echo. >> devops\monitoring\install-monitoring.sh
echo   grafana: >> devops\monitoring\install-monitoring.sh
echo     image: grafana/grafana >> devops\monitoring\install-monitoring.sh
echo     container_name: grafana >> devops\monitoring\install-monitoring.sh
echo     restart: always >> devops\monitoring\install-monitoring.sh
echo     volumes: >> devops\monitoring\install-monitoring.sh
echo       - /opt/grafana-data:/var/lib/grafana >> devops\monitoring\install-monitoring.sh
echo     environment: >> devops\monitoring\install-monitoring.sh
echo       - GF_SECURITY_ADMIN_PASSWORD=StrongGrafanaPassword2025 >> devops\monitoring\install-monitoring.sh
echo       - GF_USERS_ALLOW_SIGN_UP=false >> devops\monitoring\install-monitoring.sh
echo     ports: >> devops\monitoring\install-monitoring.sh
echo       - "3000:3000" >> devops\monitoring\install-monitoring.sh
echo     networks: >> devops\monitoring\install-monitoring.sh
echo       - monitoring-network >> devops\monitoring\install-monitoring.sh
echo. >> devops\monitoring\install-monitoring.sh
echo   node-exporter: >> devops\monitoring\install-monitoring.sh
echo     image: prom/node-exporter >> devops\monitoring\install-monitoring.sh
echo     container_name: node-exporter >> devops\monitoring\install-monitoring.sh
echo     restart: always >> devops\monitoring\install-monitoring.sh
echo     volumes: >> devops\monitoring\install-monitoring.sh
echo       - /proc:/host/proc:ro >> devops\monitoring\install-monitoring.sh
echo       - /sys:/host/sys:ro >> devops\monitoring\install-monitoring.sh
echo       - /:/rootfs:ro >> devops\monitoring\install-monitoring.sh
echo     command: >> devops\monitoring\install-monitoring.sh
echo       - '--path.procfs=/host/proc' >> devops\monitoring\install-monitoring.sh
echo       - '--path.sysfs=/host/sys' >> devops\monitoring\install-monitoring.sh
echo       - '--collector.filesystem.ignored-mount-points=^/(sys|proc|dev|host|etc)(.*)' >> devops\monitoring\install-monitoring.sh
echo     networks: >> devops\monitoring\install-monitoring.sh
echo       - monitoring-network >> devops\monitoring\install-monitoring.sh
echo. >> devops\monitoring\install-monitoring.sh
echo   cadvisor: >> devops\monitoring\install-monitoring.sh
echo     image: gcr.io/cadvisor/cadvisor >> devops\monitoring\install-monitoring.sh
echo     container_name: cadvisor >> devops\monitoring\install-monitoring.sh
echo     restart: always >> devops\monitoring\install-monitoring.sh
echo     volumes: >> devops\monitoring\install-monitoring.sh
echo       - /:/rootfs:ro >> devops\monitoring\install-monitoring.sh
echo       - /var/run:/var/run:ro >> devops\monitoring\install-monitoring.sh
echo       - /sys:/sys:ro >> devops\monitoring\install-monitoring.sh
echo       - /var/lib/docker/:/var/lib/docker:ro >> devops\monitoring\install-monitoring.sh
echo     networks: >> devops\monitoring\install-monitoring.sh
echo       - monitoring-network >> devops\monitoring\install-monitoring.sh
echo. >> devops\monitoring\install-monitoring.sh
echo networks: >> devops\monitoring\install-monitoring.sh
echo   monitoring-network: >> devops\monitoring\install-monitoring.sh
echo     driver: bridge >> devops\monitoring\install-monitoring.sh
echo EOL >> devops\monitoring\install-monitoring.sh
echo. >> devops\monitoring\install-monitoring.sh
echo # Iniciar los contenedores de monitoreo >> devops\monitoring\install-monitoring.sh
echo cd /opt && docker-compose -f monitoring-docker-compose.yml up -d >> devops\monitoring\install-monitoring.sh
echo. >> devops\monitoring\install-monitoring.sh
echo echo "Prometheus y Grafana instalados correctamente!" >> devops\monitoring\install-monitoring.sh
echo echo "Accede a Prometheus: http://tu-ip:9090" >> devops\monitoring\install-monitoring.sh
echo echo "Accede a Grafana: http://tu-ip:3000 (usuario: admin, contraseña: StrongGrafanaPassword2025)" >> devops\monitoring\install-monitoring.sh

REM ==================================================
REM CONFIGURACIÓN DE LA APLICACIÓN SISTEMA-TOURS
REM ==================================================

echo #!/bin/bash > devops\scripts\deploy-app.sh
echo # Script para desplegar la aplicación Sistema-Tours >> devops\scripts\deploy-app.sh
echo. >> devops\scripts\deploy-app.sh
echo # Crear directorio para la aplicación >> devops\scripts\deploy-app.sh
echo mkdir -p /opt/sistema-tours >> devops\scripts\deploy-app.sh
echo. >> devops\scripts\deploy-app.sh
echo # Clonar repositorio >> devops\scripts\deploy-app.sh
echo git clone https://github.com/Angel226m/subidadevops.git /opt/sistema-tours >> devops\scripts\deploy-app.sh
echo. >> devops\scripts\deploy-app.sh
echo # Crear archivo .env para variables de entorno >> devops\scripts\deploy-app.sh
echo cat ^> /opt/sistema-tours/.env ^<^< EOL >> devops\scripts\deploy-app.sh
echo # Variables de entorno para Sistema-Tours >> devops\scripts\deploy-app.sh
echo DB_HOST=postgres-db >> devops\scripts\deploy-app.sh
echo DB_PORT=5432 >> devops\scripts\deploy-app.sh
echo DB_USER=sistema_tours >> devops\scripts\deploy-app.sh
echo DB_PASSWORD=StrongDBPassword2025 >> devops\scripts\deploy-app.sh
echo DB_NAME=sistema_tours_db >> devops\scripts\deploy-app.sh
echo JWT_SECRET=SistemaToursJwtSecretKey2025 >> devops\scripts\deploy-app.sh
echo EOL >> devops\scripts\deploy-app.sh
echo. >> devops\scripts\deploy-app.sh
echo # Crear docker-compose para la aplicación >> devops\scripts\deploy-app.sh
echo cat ^> /opt/sistema-tours/docker-compose.yml ^<^< EOL >> devops\scripts\deploy-app.sh
echo version: '3.8' >> devops\scripts\deploy-app.sh
echo. >> devops\scripts\deploy-app.sh
echo services: >> devops\scripts\deploy-app.sh
echo   backend: >> devops\scripts\deploy-app.sh
echo     build: >> devops\scripts\deploy-app.sh
echo       context: ./backend >> devops\scripts\deploy-app.sh
echo       dockerfile: Dockerfile >> devops\scripts\deploy-app.sh
echo     container_name: sistema-tours-backend >> devops\scripts\deploy-app.sh
echo     restart: always >> devops\scripts\deploy-app.sh
echo     environment: >> devops\scripts\deploy-app.sh
echo       DB_HOST: \${DB_HOST} >> devops\scripts\deploy-app.sh
echo       DB_PORT: \${DB_PORT} >> devops\scripts\deploy-app.sh
echo       DB_USER: \${DB_USER} >> devops\scripts\deploy-app.sh
echo       DB_PASSWORD: \${DB_PASSWORD} >> devops\scripts\deploy-app.sh
echo       DB_NAME: \${DB_NAME} >> devops\scripts\deploy-app.sh
echo       JWT_SECRET: \${JWT_SECRET} >> devops\scripts\deploy-app.sh
echo     ports: >> devops\scripts\deploy-app.sh
echo       - "8080:8080" >> devops\scripts\deploy-app.sh
echo     depends_on: >> devops\scripts\deploy-app.sh
echo       - migrations >> devops\scripts\deploy-app.sh
echo     networks: >> devops\scripts\deploy-app.sh
echo       - app-network >> devops\scripts\deploy-app.sh
echo       - backend-network >> devops\scripts\deploy-app.sh
echo. >> devops\scripts\deploy-app.sh
echo   migrations: >> devops\scripts\deploy-app.sh
echo     build: >> devops\scripts\deploy-app.sh
echo       context: ./backend >> devops\scripts\deploy-app.sh
echo       dockerfile: Dockerfile.migrations >> devops\scripts\deploy-app.sh
echo     container_name: sistema-tours-migrations >> devops\scripts\deploy-app.sh
echo     environment: >> devops\scripts\deploy-app.sh
echo       DB_HOST: \${DB_HOST} >> devops\scripts\deploy-app.sh
echo       DB_PORT: \${DB_PORT} >> devops\scripts\deploy-app.sh
echo       DB_USER: \${DB_USER} >> devops\scripts\deploy-app.sh
echo       DB_PASSWORD: \${DB_PASSWORD} >> devops\scripts\deploy-app.sh
echo       DB_NAME: \${DB_NAME} >> devops\scripts\deploy-app.sh
echo     networks: >> devops\scripts\deploy-app.sh
echo       - backend-network >> devops\scripts\deploy-app.sh
echo. >> devops\scripts\deploy-app.sh
echo   nginx: >> devops\scripts\deploy-app.sh
echo     image: nginx:alpine >> devops\scripts\deploy-app.sh
echo     container_name: sistema-tours-nginx >> devops\scripts\deploy-app.sh
echo     restart: always >> devops\scripts\deploy-app.sh
echo     ports: >> devops\scripts\deploy-app.sh
echo       - "80:80" >> devops\scripts\deploy-app.sh
echo       - "443:443" >> devops\scripts\deploy-app.sh
echo     volumes: >> devops\scripts\deploy-app.sh
echo       - ./nginx/nginx.conf:/etc/nginx/nginx.conf >> devops\scripts\deploy-app.sh
echo       - ./nginx/conf.d:/etc/nginx/conf.d >> devops\scripts\deploy-app.sh
echo       - ./ssl:/etc/nginx/ssl >> devops\scripts\deploy-app.sh
echo       - ./certbot/conf:/etc/letsencrypt >> devops\scripts\deploy-app.sh
echo       - ./certbot/www:/var/www/certbot >> devops\scripts\deploy-app.sh
echo     depends_on: >> devops\scripts\deploy-app.sh
echo       - backend >> devops\scripts\deploy-app.sh
echo     networks: >> devops\scripts\deploy-app.sh
echo       - app-network >> devops\scripts\deploy-app.sh
echo. >> devops\scripts\deploy-app.sh
echo networks: >> devops\scripts\deploy-app.sh
echo   app-network: >> devops\scripts\deploy-app.sh
echo     driver: bridge >> devops\scripts\deploy-app.sh
echo   backend-network: >> devops\scripts\deploy-app.sh
echo     external: true >> devops\scripts\deploy-app.sh
echo EOL >> devops\scripts\deploy-app.sh
echo. >> devops\scripts\deploy-app.sh
echo # Crear directorios necesarios >> devops\scripts\deploy-app.sh
echo mkdir -p /opt/sistema-tours/nginx/conf.d /opt/sistema-tours/ssl /opt/sistema-tours/certbot/conf /opt/sistema-tours/certbot/www >> devops\scripts\deploy-app.sh
echo. >> devops\scripts\deploy-app.sh
echo # Configurar Nginx >> devops\scripts\deploy-app.sh
echo cat ^> /opt/sistema-tours/nginx/nginx.conf ^<^< EOL >> devops\scripts\deploy-app.sh
echo user nginx; >> devops\scripts\deploy-app.sh
echo worker_processes auto; >> devops\scripts\deploy-app.sh
echo error_log /var/log/nginx/error.log; >> devops\scripts\deploy-app.sh
echo pid /run/nginx.pid; >> devops\scripts\deploy-app.sh
echo. >> devops\scripts\deploy-app.sh
echo events { >> devops\scripts\deploy-app.sh
echo     worker_connections 1024; >> devops\scripts\deploy-app.sh
echo } >> devops\scripts\deploy-app.sh
echo. >> devops\scripts\deploy-app.sh
echo http { >> devops\scripts\deploy-app.sh
echo     include       /etc/nginx/mime.types; >> devops\scripts\deploy-app.sh
echo     default_type  application/octet-stream; >> devops\scripts\deploy-app.sh
echo. >> devops\scripts\deploy-app.sh
echo     log_format  main  '\$remote_addr - \$remote_user [\$time_local] "\$request" ' >> devops\scripts\deploy-app.sh
echo                       '\$status \$body_bytes_sent "\$http_referer" ' >> devops\scripts\deploy-app.sh
echo                       '"\$http_user_agent" "\$http_x_forwarded_for"'; >> devops\scripts\deploy-app.sh
echo. >> devops\scripts\deploy-app.sh
echo     access_log  /var/log/nginx/access.log  main; >> devops\scripts\deploy-app.sh
echo. >> devops\scripts\deploy-app.sh
echo     sendfile        on; >> devops\scripts\deploy-app.sh
echo     keepalive_timeout  65; >> devops\scripts\deploy-app.sh
echo     gzip  on; >> devops\scripts\deploy-app.sh
echo. >> devops\scripts\deploy-app.sh
echo     include /etc/nginx/conf.d/*.conf; >> devops\scripts\deploy-app.sh
echo } >> devops\scripts\deploy-app.sh
echo EOL >> devops\scripts\deploy-app.sh
echo. >> devops\scripts\deploy-app.sh
echo # Configurar sitio en Nginx >> devops\scripts\deploy-app.sh
echo cat ^> /opt/sistema-tours/nginx/conf.d/sistema-tours.conf ^<^< EOL >> devops\scripts\deploy-app.sh
echo server { >> devops\scripts\deploy-app.sh
echo     listen 80; >> devops\scripts\deploy-app.sh
echo     server_name sistema-tours.com www.sistema-tours.com; >> devops\scripts\deploy-app.sh
echo. >> devops\scripts\deploy-app.sh
echo     location /.well-known/acme-challenge/ { >> devops\scripts\deploy-app.sh
echo         root /var/www/certbot; >> devops\scripts\deploy-app.sh
echo     } >> devops\scripts\deploy-app.sh
echo. >> devops\scripts\deploy-app.sh
echo     location / { >> devops\scripts\deploy-app.sh
echo         return 301 https://\$host\$request_uri; >> devops\scripts\deploy-app.sh
echo     } >> devops\scripts\deploy-app.sh
echo } >> devops\scripts\deploy-app.sh
echo. >> devops\scripts\deploy-app.sh
echo server { >> devops\scripts\deploy-app.sh
echo     listen 443 ssl; >> devops\scripts\deploy-app.sh
echo     server_name sistema-tours.com www.sistema-tours.com; >> devops\scripts\deploy-app.sh
echo. >> devops\scripts\deploy-app.sh
echo     ssl_certificate /etc/letsencrypt/live/sistema-tours.com/fullchain.pem; >> devops\scripts\deploy-app.sh
echo     ssl_certificate_key /etc/letsencrypt/live/sistema-tours.com/privkey.pem; >> devops\scripts\deploy-app.sh
echo. >> devops\scripts\deploy-app.sh
echo     location / { >> devops\scripts\deploy-app.sh
echo         proxy_pass http://backend:8080; >> devops\scripts\deploy-app.sh
echo         proxy_set_header Host \$host; >> devops\scripts\deploy-app.sh
echo         proxy_set_header X-Real-IP \$remote_addr; >> devops\scripts\deploy-app.sh
echo         proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for; >> devops\scripts\deploy-app.sh
echo         proxy_set_header X-Forwarded-Proto \$scheme; >> devops\scripts\deploy-app.sh
echo     } >> devops\scripts\deploy-app.sh
echo } >> devops\scripts\deploy-app.sh
echo EOL >> devops\scripts\deploy-app.sh
echo. >> devops\scripts\deploy-app.sh
echo # Crear Dockerfile para el backend >> devops\scripts\deploy-app.sh
echo cat ^> /opt/sistema-tours/backend/Dockerfile ^<^< EOL >> devops\scripts\deploy-app.sh
echo FROM golang:1.20-alpine as builder >> devops\scripts\deploy-app.sh
echo. >> devops\scripts\deploy-app.sh
echo WORKDIR /app >> devops\scripts\deploy-app.sh
echo. >> devops\scripts\deploy-app.sh
echo COPY go.mod go.sum ./ >> devops\scripts\deploy-app.sh
echo RUN go mod download >> devops\scripts\deploy-app.sh
echo. >> devops\scripts\deploy-app.sh
echo COPY . . >> devops\scripts\deploy-app.sh
echo RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main.go >> devops\scripts\deploy-app.sh
echo. >> devops\scripts\deploy-app.sh
echo FROM alpine:latest >> devops\scripts\deploy-app.sh
echo RUN apk --no-cache add ca-certificates tzdata >> devops\scripts\deploy-app.sh
echo. >> devops\scripts\deploy-app.sh
echo WORKDIR /app >> devops\scripts\deploy-app.sh
echo COPY --from=builder /app/main . >> devops\scripts\deploy-app.sh
echo. >> devops\scripts\deploy-app.sh
echo EXPOSE 8080 >> devops\scripts\deploy-app.sh
echo CMD ["./main"] >> devops\scripts\deploy-app.sh
echo EOL >> devops\scripts\deploy-app.sh
echo. >> devops\scripts\deploy-app.sh
echo # Crear Dockerfile para migraciones >> devops\scripts\deploy-app.sh
echo cat ^> /opt/sistema-tours/backend/Dockerfile.migrations ^<^< EOL >> devops\scripts\deploy-app.sh
echo FROM golang:1.20-alpine >> devops\scripts\deploy-app.sh
echo. >> devops\scripts\deploy-app.sh
echo # Instalar herramientas necesarias para migraciones >> devops\scripts\deploy-app.sh
echo RUN apk add --no-cache postgresql-client >> devops\scripts\deploy-app.sh
echo. >> devops\scripts\deploy-app.sh
echo WORKDIR /app >> devops\scripts\deploy-app.sh
echo COPY migrations/ ./migrations/ >> devops\scripts\deploy-app.sh
echo. >> devops\scripts\deploy-app.sh
echo # Script para ejecutar migraciones >> devops\scripts\deploy-app.sh
echo COPY scripts/run-migrations.sh . >> devops\scripts\deploy-app.sh
echo RUN chmod +x run-migrations.sh >> devops\scripts\deploy-app.sh
echo. >> devops\scripts\deploy-app.sh
echo CMD ["./run-migrations.sh"] >> devops\scripts\deploy-app.sh
echo EOL >> devops\scripts\deploy-app.sh
echo. >> devops\scripts\deploy-app.sh
echo # Crear script para ejecutar migraciones >> devops\scripts\deploy-app.sh
echo mkdir -p /opt/sistema-tours/backend/scripts >> devops\scripts\deploy-app.sh
echo cat ^> /opt/sistema-tours/backend/scripts/run-migrations.sh ^<^< EOL >> devops\scripts\deploy-app.sh
echo #!/bin/sh >> devops\scripts\deploy-app.sh
echo. >> devops\scripts\deploy-app.sh
echo echo "Esperando a que PostgreSQL esté disponible..." >> devops\scripts\deploy-app.sh
echo. >> devops\scripts\deploy-app.sh
echo # Esperar a que PostgreSQL esté listo >> devops\scripts\deploy-app.sh
echo until PGPASSWORD=\$DB_PASSWORD psql -h \$DB_HOST -U \$DB_USER -d \$DB_NAME -c '\q'; do >> devops\scripts\deploy-app.sh
echo   echo "PostgreSQL no está disponible todavía - esperando..." >> devops\scripts\deploy-app.sh
echo   sleep 2 >> devops\scripts\deploy-app.sh
echo done >> devops\scripts\deploy-app.sh
echo. >> devops\scripts\deploy-app.sh
echo echo "PostgreSQL está listo! Ejecutando migraciones..." >> devops\scripts\deploy-app.sh
echo. >> devops\scripts\deploy-app.sh
echo # Ejecutar el archivo SQL de migraciones >> devops\scripts\deploy-app.sh
echo PGPASSWORD=\$DB_PASSWORD psql -h \$DB_HOST -U \$DB_USER -d \$DB_NAME -f /app/migrations/crear_tablas.sql >> devops\scripts\deploy-app.sh
echo. >> devops\scripts\deploy-app.sh
echo echo "Migraciones completadas con éxito!" >> devops\scripts\deploy-app.sh
echo EOL >> devops\scripts\deploy-app.sh
echo chmod +x /opt/sistema-tours/backend/scripts/run-migrations.sh >> devops\scripts\deploy-app.sh
echo. >> devops\scripts\deploy-app.sh
echo # Obtener certificados SSL con Let's Encrypt (Certbot) >> devops\scripts\deploy-app.sh
echo echo "Para obtener certificados SSL, ejecuta el siguiente comando (ajusta el dominio):" >> devops\scripts\deploy-app.sh
echo echo "certbot --nginx -d sistema-tours.com -d www.sistema-tours.com" >> devops\scripts\deploy-app.sh
echo. >> devops\scripts\deploy-app.sh
echo # Iniciar la aplicación >> devops\scripts\deploy-app.sh
echo cd /opt/sistema-tours && docker-compose up -d >> devops\scripts\deploy-app.sh
echo. >> devops\scripts\deploy-app.sh
echo echo "Aplicación Sistema-Tours desplegada correctamente!" >> devops\scripts\deploy-app.sh

REM ==================================================
REM SCRIPT GENERAL DE INSTALACIÓN
REM ==================================================

echo #!/bin/bash > devops\scripts\master-setup.sh
echo # Script maestro para configurar todo el entorno >> devops\scripts\master-setup.sh
echo. >> devops\scripts\master-setup.sh
echo echo "=====================================================" >> devops\scripts\master-setup.sh
echo echo "CONFIGURACIÓN COMPLETA DE DEVOPS PARA SISTEMA-TOURS" >> devops\scripts\master-setup.sh
echo echo "=====================================================" >> devops\scripts\master-setup.sh
echo echo >> devops\scripts\master-setup.sh
echo. >> devops\scripts\master-setup.sh
echo # 1. Configuración inicial del servidor >> devops\scripts\master-setup.sh
echo echo "1. Configurando el servidor..." >> devops\scripts\master-setup.sh
echo bash /opt/sistema-tours/devops/hetzner-vps/setup-server.sh >> devops\scripts\master-setup.sh
echo echo "✓ Servidor configurado correctamente" >> devops\scripts\master-setup.sh
echo echo >> devops\scripts\master-setup.sh
echo. >> devops\scripts\master-setup.sh
echo # 2. Instalar Docker y Docker Compose >> devops\scripts\master-setup.sh
echo echo "2. Instalando Docker y Docker Compose..." >> devops\scripts\master-setup.sh
echo bash /opt/sistema-tours/devops/hetzner-vps/install-docker.sh >> devops\scripts\master-setup.sh
echo echo "✓ Docker instalado correctamente" >> devops\scripts\master-setup.sh
echo echo >> devops\scripts\master-setup.sh
echo. >> devops\scripts\master-setup.sh
echo # 3. Configurar PostgreSQL >> devops\scripts\master-setup.sh
echo echo "3. Configurando PostgreSQL..." >> devops\scripts\master-setup.sh
echo bash /opt/sistema-tours/devops/database/setup-postgres.sh >> devops\scripts\master-setup.sh
echo echo "✓ PostgreSQL configurado correctamente" >> devops\scripts\master-setup.sh
echo echo >> devops\scripts\master-setup.sh
echo. >> devops\scripts\master-setup.sh
echo # 4. Instalar Terraform >> devops\scripts\master-setup.sh
echo echo "4. Instalando Terraform..." >> devops\scripts\master-setup.sh
echo bash /opt/sistema-tours/devops/terraform/install-terraform.sh >> devops\scripts\master-setup.sh
echo echo "✓ Terraform instalado correctamente" >> devops\scripts\master-setup.sh
echo echo >> devops\scripts\master-setup.sh
echo. >> devops\scripts\master-setup.sh
echo # 5. Instalar Jenkins >> devops\scripts\master-setup.sh
echo echo "5. Instalando Jenkins..." >> devops\scripts\master-setup.sh
echo bash /opt/sistema-tours/devops/jenkins/install-jenkins.sh >> devops\scripts\master-setup.sh
echo echo "✓ Jenkins instalado correctamente" >> devops\scripts\master-setup.sh
echo echo >> devops\scripts\master-setup.sh
echo. >> devops\scripts\master-setup.sh
echo # 6. Configurar Prometheus y Grafana >> devops\scripts\master-setup.sh
echo echo "6. Configurando Prometheus y Grafana..." >> devops\scripts\master-setup.sh
echo bash /opt/sistema-tours/devops/monitoring/install-monitoring.sh >> devops\scripts\master-setup.sh
echo echo "✓ Monitoreo configurado correctamente" >> devops\scripts\master-setup.sh
echo echo >> devops\scripts\master-setup.sh
echo. >> devops\scripts\master-setup.sh
echo # 7. Desplegar la aplicación >> devops\scripts\master-setup.sh
echo echo "7. Desplegando la aplicación Sistema-Tours..." >> devops\scripts\master-setup.sh
echo bash /opt/sistema-tours/devops/scripts/deploy-app.sh >> devops\scripts\master-setup.sh
echo echo "✓ Aplicación desplegada correctamente" >> devops\scripts\master-setup.sh
echo echo >> devops\scripts\master-setup.sh
echo. >> devops\scripts\master-setup.sh
echo echo "=====================================================" >> devops\scripts\master-setup.sh
echo echo "¡CONFIGURACIÓN COMPLETA FINALIZADA!" >> devops\scripts\master-setup.sh
echo echo "=====================================================" >> devops\scripts\master-setup.sh
echo echo >> devops\scripts\master-setup.sh
echo echo "URLs de acceso:" >> devops\scripts\master-setup.sh
echo echo "- Aplicación: https://sistema-tours.com" >> devops\scripts\master-setup.sh
echo echo "- Jenkins: http://[IP-SERVIDOR]:8080" >> devops\scripts\master-setup.sh
echo echo "- Grafana: http://[IP-SERVIDOR]:3000" >> devops\scripts\master-setup.sh
echo echo "- Prometheus: http://[IP-SERVIDOR]:9090" >> devops\scripts\master-setup.sh
echo echo "- PostgreSQL: [IP-SERVIDOR]:5432" >> devops\scripts\master-setup.sh
echo echo "- pgAdmin: http://[IP-SERVIDOR]:5050" >> devops\scripts\master-setup.sh
echo echo >> devops\scripts\master-setup.sh
echo echo "Recuerda cambiar las contraseñas predeterminadas y actualizar la configuración según tus necesidades." >> devops\scripts\master-setup.sh

REM ==================================================
REM AUTOMATIZACIÓN CON GITHUB ACTIONS
REM ==================================================

echo name: CI/CD Pipeline > .github\workflows\ci-cd.yml
echo. >> .github\workflows\ci-cd.yml
echo on: >> .github\workflows\ci-cd.yml
echo   push: >> .github\workflows\ci-cd.yml
echo     branches: [ main ] >> .github\workflows\ci-cd.yml
echo   pull_request: >> .github\workflows\ci-cd.yml
echo     branches: [ main ] >> .github\workflows\ci-cd.yml
echo   workflow_dispatch: >> .github\workflows\ci-cd.yml
echo. >> .github\workflows\ci-cd.yml
echo jobs: >> .github\workflows\ci-cd.yml
echo   test: >> .github\workflows\ci-cd.yml
echo     name: Test >> .github\workflows\ci-cd.yml
echo     runs-on: ubuntu-latest >> .github\workflows\ci-cd.yml
echo     steps: >> .github\workflows\ci-cd.yml
echo       - uses: actions/checkout@v3 >> .github\workflows\ci-cd.yml
echo. >> .github\workflows\ci-cd.yml
echo       - name: Set up Go >> .github\workflows\ci-cd.yml
echo         uses: actions/setup-go@v4 >> .github\workflows\ci-cd.yml
echo         with: >> .github\workflows\ci-cd.yml
echo           go-version: '1.20' >> .github\workflows\ci-cd.yml
echo. >> .github\workflows\ci-cd.yml
echo       - name: Test >> .github\workflows\ci-cd.yml
echo         run: | >> .github\workflows\ci-cd.yml
echo           cd backend >> .github\workflows\ci-cd.yml
echo           go test -v ./... >> .github\workflows\ci-cd.yml
echo. >> .github\workflows\ci-cd.yml
echo   build: >> .github\workflows\ci-cd.yml
echo     name: Build >> .github\workflows\ci-cd.yml
echo     needs: test >> .github\workflows\