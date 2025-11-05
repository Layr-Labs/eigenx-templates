# eigenx-tee-golang-app

## Development

### Setup & Local Testing
```bash
go mod download
cp .env.example .env
go run src/main.go
```

### Docker Testing
```bash
docker build -t my-app .
docker run --rm --env-file .env my-app
```

### Environment
- `MNEMONIC`: 12/24-word BIP39 phrase used to derive the signing key (`m/44'/60'/0'/0/0`).
- `PORT` (optional): server port, defaults to `8080`.

### API
- `GET /random` → `{ randomNumber, timestamp, message, messageHash, signature, signer }`
