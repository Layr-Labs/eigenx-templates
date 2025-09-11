use alloy_signer_local::{coins_bip39::English, MnemonicBuilder};
use std::env;

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let _ = dotenvy::dotenv();

    let mnemonic = env::var("MNEMONIC")?;
    
    let wallet = MnemonicBuilder::<English>::default()
        .phrase(mnemonic)
        .index(0)?
        .build()?;
    
    println!("First wallet address: {}", wallet.address());
    
    Ok(())
}
