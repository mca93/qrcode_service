# QR Code Service API

API em Go para gestão de QR Codes, Templates e ClientApps com encriptação e geração visual de QRCodes.

## 🔧 Tecnologias
- Golang + Gin
- PostgreSQL + GORM
- Docker + Docker Compose
- Geração de QR Codes com `github.com/skip2/go-qrcode`
- Encriptação AES (em repouso e em trânsito)
- Swagger UI
- Frontend com Next.js + TailwindCSS + shadcn/ui

## 🚀 Como correr localmente

```bash
docker-compose up --build
```

A app estará disponível em `http://localhost:8080`.

Swagger estará disponível em:
```
http://localhost:8080/swagger/index.html
```

## 🔐 Autenticação

Todas as rotas da API estão protegidas por API Key.
Adiciona o header `X-API-Key` com o valor da tua chave.

```bash
curl -H "X-API-Key: minha-chave" http://localhost:8080/v1/clientapps
```

## 🔒 Encriptação de dados

A app suporta encriptação de dados sensíveis (como dados de QRCode e templates) usando AES.

Exemplo para encriptar/desencriptar um dado:

```go
import "qrcode_service/utils"

cipherText, err := utils.Encrypt("meu_dado_sensivel")
plainText, err := utils.Decrypt(cipherText)
```

A chave usada encontra-se em `utils/crypto.go` e deve ser armazenada de forma segura em produção (via secrets manager ou variável de ambiente).

## 🧪 Testes

```bash
go test ./...
```

## 🖥️ Frontend

Este projeto inclui um dashboard web moderno com Next.js + TailwindCSS + shadcn/ui.

### Para correr o frontend:

```bash
cd qrcode_service/frontend
npm install
npm run dev
```

Cria um ficheiro `.env.local` com:

```
NEXT_PUBLIC_API_URL=http://localhost:8080/v1
NEXT_PUBLIC_API_KEY=minha-chave-secreta
```

## 📦 Endpoints

Veja o Swagger ou o `swagger.yaml` para a lista completa de endpoints.

## ✍️ Autor
Feito com ❤️ por pavulla. Todos os direitos reservados
