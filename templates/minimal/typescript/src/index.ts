import { mnemonicToAccount } from 'viem/accounts';
import dotenv from 'dotenv';

dotenv.config();

async function main() {
  const mnemonic = process.env.MNEMONIC;
  
  if (!mnemonic) {
    console.error('MNEMONIC environment variable is not set');
    process.exit(1);
  }
  
  try {
    const account = mnemonicToAccount(mnemonic);
    console.log('First wallet address:', account.address);
  } catch (error) {
    console.error('Error generating wallet:', error);
    process.exit(1);
  }
}

main().catch(console.error);
