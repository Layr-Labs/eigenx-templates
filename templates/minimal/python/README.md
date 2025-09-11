# eigenx-tee-python-app

## Prerequisites

Before deploying, you'll need:

- **Docker** - To package and publish your application image
  - [Download Docker](https://www.docker.com/get-started/)
  - You'll also need to `docker login` to push images to your registry
- **Sepolia ETH** - To pay for deployment transactions
  - [Google Cloud Faucet](https://cloud.google.com/application/web3/faucet/ethereum/sepolia)
  - [Alchemy Faucet](https://sepoliafaucet.com/)

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

## Deployment

```bash
# Store your private key (generate new or use existing)
eigenx auth generate --store
# OR: eigenx auth login (if you have an existing key)

eigenx app deploy username/image-name
```

The CLI will automatically detect the `Dockerfile` and build your app before deploying.

## Management & Monitoring

### App Lifecycle
```bash
eigenx app list                    # List all apps
eigenx app info [app-name]         # Get app details
eigenx app logs [app-name]         # View logs
eigenx app start [app-name]        # Start stopped app
eigenx app stop [app-name]         # Stop running app
eigenx app terminate [app-name]    # Terminate app
eigenx app upgrade [app-name] [image] # Update deployment
```

### App Naming
```bash
eigenx app name [app-id] [new-name]  # Update friendly name
```

## Documentation

[EigenX CLI Documentation](https://github.com/Layr-Labs/eigenx-cli/blob/main/README.md)
