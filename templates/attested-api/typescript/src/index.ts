import dotenv from 'dotenv';
import Fastify from 'fastify';
import { randomBytes } from 'node:crypto';
import { hashMessage } from 'viem';
import { mnemonicToAccount } from 'viem/accounts';

dotenv.config();

async function main() {
  const mnemonic = process.env.MNEMONIC;

  if (!mnemonic) {
    console.error('MNEMONIC environment variable is not set');
    process.exit(1);
  }

  let account;
  try {
    account = mnemonicToAccount(mnemonic);
  } catch (error) {
    console.error('Error deriving signing account:', error);
    process.exit(1);
  }

  const server = Fastify({ logger: true });

  server.get('/random', async () => {
    const entropy = randomBytes(32);
    const randomNumber = `0x${entropy.toString('hex')}`;
    const randomNumberDecimal = BigInt(randomNumber).toString();
    const timestamp = new Date().toISOString();
    const message = `RandomnessBeacon|${randomNumber}|${timestamp}`;
    const messageHash = hashMessage(message);
    const signature = await account.signMessage({ message });

    return {
      randomNumber,
      randomNumberDecimal,
      timestamp,
      message,
      messageHash,
      signature,
      signer: account.address,
    };
  });

  const port = Number(process.env.PORT ?? 8080);
  try {
    await server.listen({ port, host: '0.0.0.0' });
  } catch (error) {
    server.log.error(error);
    process.exit(1);
  }
}

main().catch((error) => {
  console.error('Fatal error starting server:', error);
  process.exit(1);
});
