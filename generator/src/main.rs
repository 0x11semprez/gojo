mod all_computer_infos;
mod bitcoin_generate_keys;
mod monero_generate_keys;

use all_computer_infos::{cpu_info, cpu_status, total_cores, total_memory};
use bitcoin_generate_keys::generate_private_key as generate_bitcoin_keys;
use machine_info::Machine;
use monero_generate_keys::generate_private_key as generate_monero_keys;

fn main() {
    let mut machine = Machine::new();

    let info = cpu_info(&mut machine);
    println!("{:?}", info);

    let status = cpu_status(&mut machine);
    println!("{:?}", status);

    let memory = total_memory(&mut machine);
    println!("Total memory: {} bytes", memory);

    let cores = total_cores(&mut machine);
    println!("Total cores: {}", cores);

    let monero_keys = generate_monero_keys(&mut machine);
    println!("{:?}", monero_keys);

    let bitcoin_keys = generate_bitcoin_keys(&mut machine);
    println!("{:?}", bitcoin_keys);
}
