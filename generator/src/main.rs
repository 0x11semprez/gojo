mod all_computer_infos;
mod bitcoin_generate_keys;
mod monero_generate_keys;

use std::env;
use std::process::ExitCode;

use all_computer_infos::{cpu_status, total_cores, total_memory};
use bitcoin_generate_keys::generate_private_key as generate_bitcoin_keys;
use machine_info::Machine;
use monero_generate_keys::generate_private_key as generate_monero_keys;
use serde::Serialize;

#[derive(Serialize)]
#[serde(tag = "network", rename_all = "lowercase")]
enum KeyMaterial {
    Bitcoin {
        private_key: String,
        public_key: String,
        address: String,
    },
    Monero {
        private_spend_key: String,
        public_spend_key: String,
        private_view_key: String,
        public_view_key: String,
    },
}

fn main() -> ExitCode {
    let network = match env::args().nth(1) {
        Some(arg) => arg,
        None => {
            eprintln!("usage: generator <bitcoin|monero>");
            return ExitCode::from(2);
        }
    };

    let mut machine = Machine::new();

    // Entropy diagnostics go to stderr, never stdout: stdout is
    // reserved for the JSON key material the Go caller parses.
    eprintln!("cpu status: {:?}", cpu_status(&mut machine));
    eprintln!("total memory: {} bytes", total_memory(&mut machine));
    eprintln!("total cores: {}", total_cores(&mut machine));

    let output = match network.as_str() {
        "bitcoin" => {
            let keys = generate_bitcoin_keys(&mut machine);
            KeyMaterial::Bitcoin {
                private_key: hex::encode(keys.private_key.inner.secret_bytes()),
                public_key: hex::encode(keys.public_key.to_bytes()),
                address: keys.address.to_string(),
            }
        }
        "monero" => {
            let keys = generate_monero_keys(&mut machine);
            KeyMaterial::Monero {
                private_spend_key: hex::encode(keys.private_spend_key.to_bytes()),
                public_spend_key: hex::encode(keys.public_spend_key.to_bytes()),
                private_view_key: hex::encode(keys.private_view_key.to_bytes()),
                public_view_key: hex::encode(keys.public_view_key.to_bytes()),
            }
        }
        other => {
            eprintln!("unsupported network: {other:?} (expected \"bitcoin\" or \"monero\")");
            return ExitCode::from(2);
        }
    };

    // Only the JSON key material goes to stdout: the Go caller reads
    // stdout as the contract. Anything else (logs, diagnostics) must
    // go to stderr instead.
    match serde_json::to_string(&output) {
        Ok(json) => {
            println!("{json}");
            ExitCode::SUCCESS
        }
        Err(err) => {
            eprintln!("failed to serialize key material: {err}");
            ExitCode::FAILURE
        }
    }
}
