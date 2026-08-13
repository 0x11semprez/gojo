use bitcoin::hashes::{sha256, Hash};
use bitcoin::secp256k1::{Secp256k1, SecretKey};
use bitcoin::{Address, CompressedPublicKey, Network, NetworkKind, PrivateKey};
use machine_info::Machine;
use rand::rngs::OsRng;
use rand::RngCore;

use crate::all_computer_infos::cpu_info;

#[derive(Debug)]
pub struct BitcoinKeys {
    pub private_key: PrivateKey,
    pub public_key: CompressedPublicKey,
    pub address: Address,
}

pub fn generate_private_key(machine: &mut Machine) -> BitcoinKeys {
    let processor = cpu_info(machine);

    let mut seed_material = Vec::new();
    seed_material.extend_from_slice(processor.vendor.as_bytes());
    seed_material.extend_from_slice(processor.brand.as_bytes());
    seed_material.extend_from_slice(&processor.frequency.to_le_bytes());

    let mut random_bytes = [0u8; 32];
    OsRng.fill_bytes(&mut random_bytes);
    seed_material.extend_from_slice(&random_bytes);

    // Unlike Monero's scalar field, secp256k1 rejects 0 and values >= the curve
    // order. That is astronomically unlikely (~2^-128) with a SHA-256 output, but
    // we bump a nonce and rehash instead of unwrapping, just in case.
    let mut nonce: u32 = 0;
    let secret_key = loop {
        let mut data = seed_material.clone();
        data.extend_from_slice(&nonce.to_le_bytes());
        let hash = sha256::Hash::hash(&data);
        match SecretKey::from_slice(hash.as_byte_array()) {
            Ok(sk) => break sk,
            Err(_) => nonce += 1,
        }
    };

    let private_key = PrivateKey::new(secret_key, NetworkKind::Main);

    let secp = Secp256k1::new();
    let public_key = CompressedPublicKey::from_private_key(&secp, &private_key)
        .expect("private key is always constructed as compressed");

    let address = Address::p2wpkh(&public_key, Network::Bitcoin);

    BitcoinKeys {
        private_key,
        public_key,
        address,
    }
}


#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_generate_private_key() {
        let mut machine = Machine::new();
        let test1 = generate_private_key(&mut machine);
        let private_key = test1.private_key;
        let secp = Secp256k1::new();
        let public_key = CompressedPublicKey::from_private_key(&secp, &private_key);

        assert_eq!(public_key, Ok(test1.public_key));
    }

}
