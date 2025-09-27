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
