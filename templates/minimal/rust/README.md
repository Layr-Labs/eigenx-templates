# eigenx-tee-rust-app

## Development

### Setup & Local Testing
```bash
cargo build
cp .env.example .env
cargo run
```

### Docker Testing
```bash
docker build -t my-app .
docker run --rm --env-file .env my-app
```
