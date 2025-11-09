# Use Go 1.25 ou superior
FROM golang:1.25-alpine

WORKDIR /app

# Copia arquivos de dependências e baixa módulos
COPY go.mod go.sum ./
RUN go mod download

# Copia o restante do código
COPY . .

# Compila o binário
RUN go build -o server .

# Expõe a porta do app
EXPOSE 8080

# Comando para rodar o app
CMD ["./server"]
