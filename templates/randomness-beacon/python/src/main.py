#!/usr/bin/env python3
import os
import secrets
import sys
from datetime import datetime, timezone

import uvicorn
from dotenv import load_dotenv
from eth_account import Account
from eth_account.messages import encode_defunct
from fastapi import FastAPI

load_dotenv()

mnemonic = os.environ.get("MNEMONIC")

if not mnemonic:
    print("MNEMONIC environment variable is not set", file=sys.stderr)
    sys.exit(1)

try:
    Account.enable_unaudited_hdwallet_features()
    ACCOUNT = Account.from_mnemonic(mnemonic, account_path="m/44'/60'/0'/0/0")
except Exception as exc:  # pragma: no cover - startup crash
    print(f"Error deriving signing account: {exc}", file=sys.stderr)
    sys.exit(1)

app = FastAPI(title="EigenX Randomness Beacon", version="1.0.0")


def _current_timestamp() -> str:
    return (
        datetime.now(timezone.utc)
        .isoformat(timespec="milliseconds")
        .replace("+00:00", "Z")
    )


@app.get("/random")
def get_random_entropy():
    random_bytes = secrets.token_bytes(32)
    random_number = "0x" + random_bytes.hex()
    random_number_decimal = str(int.from_bytes(random_bytes, byteorder="big"))
    timestamp = _current_timestamp()
    message = f"RandomnessBeacon|{random_number}|{timestamp}"
    signed_message = ACCOUNT.sign_message(encode_defunct(text=message))

    return {
        "randomNumber": random_number,
        "randomNumberDecimal": random_number_decimal,
        "timestamp": timestamp,
        "message": message,
        "messageHash": signed_message.message_hash.hex(),
        "signature": signed_message.signature.hex(),
        "signer": ACCOUNT.address,
    }


def main() -> None:
    port = int(os.environ.get("PORT", "8080"))
    uvicorn.run(app, host="0.0.0.0", port=port, log_level="info")


if __name__ == "__main__":
    main()
