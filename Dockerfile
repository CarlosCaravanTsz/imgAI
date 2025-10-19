FROM golang:1.24-alpine

# Crear directorio de trabajo
WORKDIR /app

# Copiar go.mod y go.sum y descargar dependencias
COPY go.mod go.sum ./
RUN go mod download

# Copiar el código
COPY . .

# Construir la app
RUN go build -o main ./cmd/server/main.go

# Exponer puerto
EXPOSE 8080

# Ejecutar app
CMD ["./main"]