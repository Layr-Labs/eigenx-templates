# eigenx-tee-typescript-app

## Development

### Setup & Local Testing
```bash
npm install
cp .env.example .env
npm run dev
```

### Docker Testing
```bash
docker build -t my-app .
docker run --rm --env-file .env my-app
```
