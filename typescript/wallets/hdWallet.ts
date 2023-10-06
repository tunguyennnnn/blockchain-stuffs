import * as bip39 from "bip39";
import * as ecc from "tiny-secp256k1";
import BIP32Factory, { BIP32API } from "bip32";

import * as bitcoin from "bitcoinjs-lib";
import { ECPairFactory, ECPairAPI } from "ecpair";

export type Sizes = 12 | 24;

export const generateSeedPhrase = async (size: Sizes = 24): Promise<string> => {
  if (process.env.SEED_PHRASE) {
    return process.env.SEED_PHRASE;
  }
  const seedPhrase = await bip39.generateMnemonic(size === 12 ? 128 : 256);

  return seedPhrase;
};

type BitcoinWallet = {
  privateKey: string;
  publicKey: string;
  p2pkhAddress: string;
  p2wpkhAddress: string;
};

export class HdWallet {
  seedPhraseSize: Sizes;
  seedPhrase: string;
  masterPrivateKey: string;
  bip32: BIP32API;

  constructor(seedPhraseSize: Sizes = 24) {
    this.seedPhraseSize = seedPhraseSize;
    this.bip32 = BIP32Factory(ecc);
  }

  async generateSeedPhraseAndSecret(): Promise<this> {
    console.log("process.env.SEED_PHRASE", process.env.SEED_PHRASE);
    if (
      process.env.SEED_PHRASE &&
      process.env.SEED_PHRASE.split(" ").length === this.seedPhraseSize
    ) {
      this.seedPhrase = process.env.SEED_PHRASE;
    } else {
      this.seedPhrase = await bip39.generateMnemonic(
        this.seedPhraseSize === 12 ? 128 : 256
      );
    }

    const seed = await bip39.mnemonicToSeed(this.seedPhrase);

    const secret = await this.bip32.fromSeed(seed);

    this.masterPrivateKey = secret.toBase58();
    return this;
  }

  async generateBitcoinChildWallets(): Promise<BitcoinWallet> {
    const bitcoinPath = "m/44'/0'/0'/0/";
    const node = this.bip32.fromBase58(this.masterPrivateKey);

    const ECPair: ECPairAPI = ECPairFactory(ecc);

    const child = node.derivePath(bitcoinPath + "0");

    const publicKey = child.publicKey?.toString("hex") || "";
    const privateKey = child.privateKey?.toString("hex") || "";

    if (!publicKey || !privateKey) {
      throw new Error("No public or private key");
    }

    const ecPair = ECPair.fromPublicKey(child.publicKey);

    const { address: p2pkhAddress } = bitcoin.payments.p2pkh({
      pubkey: ecPair.publicKey,
    });

    const { address: p2wpkhAddress } = bitcoin.payments.p2wpkh({
      pubkey: ecPair.publicKey,
    });

    if (!p2pkhAddress || !p2wpkhAddress) {
      throw new Error("Unable to generate addresses");
    }

    return {
      privateKey,
      publicKey,
      p2pkhAddress,
      p2wpkhAddress,
    };
  }
}
