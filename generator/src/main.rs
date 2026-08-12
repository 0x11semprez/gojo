mod all_computer_infos;

use all_computer_infos::{cpu_info, cpu_status, total_cores, total_memory};
use machine_info::Machine;

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
}
