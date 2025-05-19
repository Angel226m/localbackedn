FROM golang:1.24-alpine

WORKDIR /app

# Instalar dependencias del sistema
RUN apk add --no-cache gcc musl-dev

# Copiar archivos de módulos para aprovechar el caché
COPY go.mod go.sum ./
RUN go mod download

# Copiar el resto del código de la aplicación
COPY . .

# Compilar la aplicación (se asume que el punto de entrada está en cmd/main.go)
RUN go build -o main ./cmd/main.go

# Exponer el puerto que utiliza la aplicación
EXPOSE 8080

# Comando para ejecutar la aplicación
CMD ["./main"]