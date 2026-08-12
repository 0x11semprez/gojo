use machine_info::{Machine, Processor, SystemStatus};

pub fn cpu_status(machine: &mut Machine) -> SystemStatus {
    machine.system_status().expect("failed to read /proc/stat or /proc/meminfo")
}

pub fn cpu_info(machine: &mut Machine) -> Processor {
    machine.system_info().processor
}

pub fn total_memory(machine: &mut Machine) -> u64 {
    machine.system_info().memory
}

pub fn total_cores(machine: &mut Machine) -> usize {
    machine.system_info().total_processors
}
