import { HdWallet } from "./wallets";
require("dotenv").config({
  path: require("path").resolve(__dirname, "../.env"),
});

async function main() {
  const wallet = new HdWallet();
  await wallet.generateSeedPhraseAndSecret();

  console.log("wallet", wallet);

  const bitcoinWallet = await wallet.generateBitcoinChildWallets();
  console.log("bitcoinWallet", bitcoinWallet);
}

main()
  .then(() => {
    console.log("done");
    process.exitCode = 0;
  })
  .catch((err) => {
    console.error(err);
    process.exitCode = 1;
  });
