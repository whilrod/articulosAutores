# Etapa 1: Build
FROM golang:1.25.5-alpine AS builder

WORKDIR /app

# Copiar archivos de dependencias
COPY go.mod go.sum ./
RUN go mod download

# Copiar el código fuente
COPY . .

# Compilar la aplicación
RUN go build -o main cmd/api/main.go

# Etapa 2: Imagen final
FROM alpine:latest

WORKDIR /app

# Copiar el binario compilado
COPY --from=builder /app/main .

# Copiar migraciones
COPY --from=builder /app/migrations ./migrations

# Exponer puerto
EXPOSE 8080

# Ejecutar la aplicación
CMD ["./main"]