use std::{env, net::SocketAddr, sync::Arc};

use axum::{extract::State, http::StatusCode, response::IntoResponse, routing::get, Json, Router};
use chrono::{SecondsFormat, SubsecRound, Utc};
use ethers_signers::{coins_bip39::English, LocalWallet, MnemonicBuilder, Signer};
use hex::encode as hex_encode;
use rand::RngCore;
use serde::Serialize;
use sha3::{Digest, Keccak256};
use serde_json::json;

#[derive(Clone)]
struct AppState {
    wallet: Arc<LocalWallet>,
    signer: String,
}

#[allow(non_snake_case)]
#[derive(Serialize)]
struct BeaconResponse {
    randomNumber: String,
    randomNumberDecimal: String,
    timestamp: String,
    message: String,
    messageHash: String,
    signature: String,
    signer: String,
}

async fn random_beacon(State(state): State<AppState>) -> impl IntoResponse {
    let mut entropy = [0u8; 32];
    rand::rngs::OsRng.fill_bytes(&mut entropy);

    let random_number = format!("0x{}", hex_encode(entropy));
    let random_number_decimal = entropy
        .iter()
        .fold(num_bigint::BigUint::from(0u32), |acc, &byte| {
            (acc << 8) + num_bigint::BigUint::from(byte)
        })
        .to_string();
    let timestamp = Utc::now()
        .round_subsecs(3)
        .to_rfc3339_opts(SecondsFormat::Millis, true);
    let message = format!("RandomnessBeacon|{}|{}", random_number, timestamp);

    let signature = match state.wallet.sign_message(message.clone()).await {
        Ok(sig) => sig,
        Err(err) => {
            return (
                StatusCode::INTERNAL_SERVER_ERROR,
                Json(json!({ "error": format!("failed to sign message: {err}") })),
            )
                .into_response();
        }
    };

    let hash = ethereum_message_hash(&message);

    Json(BeaconResponse {
        randomNumber: random_number.clone(),
        randomNumberDecimal: random_number_decimal,
        timestamp,
        message,
        messageHash: format!("0x{}", hex_encode(hash)),
        signature: signature.to_string(),
        signer: state.signer.clone(),
    })
    .into_response()
}

fn ethereum_message_hash(message: &str) -> [u8; 32] {
    let prefix = format!("\x19Ethereum Signed Message:\n{}", message.len());
    let mut hasher = Keccak256::new();
    hasher.update(prefix.as_bytes());
    hasher.update(message.as_bytes());
    let result = hasher.finalize();
    let mut output = [0u8; 32];
    output.copy_from_slice(&result);
    output
}

#[tokio::main]
async fn main() {
    dotenvy::dotenv().ok();

    let mnemonic = env::var("MNEMONIC").unwrap_or_else(|_| {
        eprintln!("MNEMONIC environment variable is not set");
        std::process::exit(1);
    });

    let wallet = MnemonicBuilder::<English>::default()
        .phrase(mnemonic.trim())
        .derivation_path("m/44'/60'/0'/0/0")
        .expect("invalid derivation path")
        .build()
        .unwrap_or_else(|err| {
            eprintln!("Failed to derive wallet: {err}");
            std::process::exit(1);
        });

    let signer_address = format!("{:#x}", wallet.address());

    let state = AppState {
        wallet: Arc::new(wallet),
        signer: signer_address.clone(),
    };

    let app = Router::new().route("/random", get(random_beacon)).with_state(state);

    let port: u16 = env::var("PORT")
        .ok()
        .and_then(|value| value.parse().ok())
        .unwrap_or(8080);
    let addr = SocketAddr::from(([0, 0, 0, 0], port));

    println!("Randomness beacon serving on {addr} with signer {signer_address}");

    let listener = tokio::net::TcpListener::bind(addr)
        .await
        .unwrap_or_else(|err| {
            eprintln!("Failed to bind to {addr}: {err}");
            std::process::exit(1);
        });

    if let Err(err) = axum::serve(listener, app).await {
        eprintln!("Server error: {err}");
        std::process::exit(1);
    }
}
