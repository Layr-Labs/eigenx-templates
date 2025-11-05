# eigenx-tee-python-app

## Development

### Setup & Local Testing
```bash
pip install -r requirements.txt
cp .env.example .env
python src/main.py
```

### Docker Testing
```bash
docker build -t my-app .
docker run --rm --env-file .env my-app
```

### Environment
- `MNEMONIC`: 12/24-word BIP39 phrase used to derive the signer (`m/44'/60'/0'/0/0`).
- `PORT` (optional): server port, defaults to `8080`.

### API
- `GET /random` → `{ randomNumber, timestamp, message, messageHash, signature, signer }`
