use machine_info::Machine;
use monero::{Hash, PrivateKey, PublicKey};
use rand::rngs::OsRng;
use rand::RngCore;

use crate::all_computer_infos::cpu_info;

#[derive(Debug)]
pub struct MoneroKeys {
    pub private_spend_key: PrivateKey,
    pub public_spend_key: PublicKey,
    pub private_view_key: PrivateKey,
    pub public_view_key: PublicKey,
}

pub fn generate_private_key(machine: &mut Machine) -> MoneroKeys {
    let processor = cpu_info(machine);

    let mut seed_material = Vec::new();
    seed_material.extend_from_slice(processor.vendor.as_bytes());
    seed_material.extend_from_slice(processor.brand.as_bytes());
    seed_material.extend_from_slice(&processor.frequency.to_le_bytes());

    let mut random_bytes = [0u8; 32];
    OsRng.fill_bytes(&mut random_bytes);
    seed_material.extend_from_slice(&random_bytes);

    let private_spend_key = Hash::hash_to_scalar(&seed_material);
    let public_spend_key = PublicKey::from_private_key(&private_spend_key);

    let private_view_key = Hash::hash_to_scalar(private_spend_key.as_bytes());
    let public_view_key = PublicKey::from_private_key(&private_view_key);

    MoneroKeys {
        private_spend_key,
        public_spend_key,
        private_view_key,
        public_view_key,
    }
}

#[cfg(test)]

mod tests {
    use super::*;

    fn test_generate_private_key() {
        let mut machine = Machine::new();
        let test2 = generate_private_key(&mut machine);

        let private_key = test2.private_spend_key;
        let public_key = PublicKey::from_private_key(&private_key);

        assert_eq!(public_key, test2.public_spend_key );
    }
}
