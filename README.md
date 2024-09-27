# eBPF file-delete tracer

## Overview
This project demonstrates the use of eBPF (Extended Berkeley Packet Filter) to trace file deletion events on a Linux system. The eBPF program is loaded from an ELF file and attached to the appropriate kernel hooks to monitor file deletions.

## Prerequisites
- Go (version 1.16 or higher)
- Linux system with eBPF support
- `clang` and `llvm` for compiling eBPF programs
- `libbpf` library
- `bpftool` command-line tool to inspect and manage BPF objects

## Installation
1. **Install Go**: Follow the instructions on the [official Go website](https://golang.org/doc/install).
2. **Install clang and llvm**:
    ```sh
    sudo apt-get install clang llvm
    ```
3. **Install libbpf**:
    ```sh
    sudo apt-get install libbpf-dev
    ```
4. **Install bpftool**
    ```sh
    # Package bpftool is a virtual package provided by linux-tools-common
    sudo apt-get install linux-tools-common
    ```
4. **Clone the repository**:
    ```sh
    git clone https://github.com/srodi/ebpf-file-delete-tracer.git
    cd ebpf-file-delete-tracer
    ```

## Usage
1. **Compile the eBPF program**:
    ```sh
    clang -O2 -g -target bpf -c trace_file_delete.c -o trace_file_delete.o
    ```
2. **Run the Go program**:
    ```sh
    go run main.go
    ```

## Code Explanation
- **main.go**: The main Go file that sets up the environment, loads the eBPF program, and attaches it to the kernel hooks.
    - Sets the `RLIMIT_MEMLOCK` to allow locking memory for eBPF.
    - Loads the eBPF program from the specified ELF file.
    - Attaches the eBPF program to the appropriate kernel hooks to trace file deletions.

## Contributing
We welcome contributions! Please follow these steps:
1. Fork the repository.
2. Create a new branch (`git checkout -b feature-branch`).
3. Make your changes.
4. Commit your changes (`git commit -am 'Add new feature'`).
5. Push to the branch (`git push origin feature-branch`).
6. Create a new Pull Request.

## License
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---


## Environment set-up on WSL2 - Ubuntu-24.04

Update the package lists for upgrades and new package installations, and upgrade all installed packages to their latest versions without prompting for confirmation.
```sh
sudo apt update && sudo apt upgrade --yes
```

Install the libbpf development library and searche for the file named `"bpf_helpers.h"` within the `/usr/include` directory.

```sh
sudo apt install libbpf-dev
find /usr/include -name "bpf_helpers.h"
```

Install and use `bpftool`. If you are using `WSL2` you will have to install manually

```sh
git clone --recurse-submodules https://github.com/libbpf/bpftool.git
cd bpftool/src
sudo make install
```

Generate `vmlinux.h` using `bpftools`

```sh
bpftool btf dump file /sys/kernel/btf/vmlinux format c > vmlinux.h
```

Ensure that the BPF filesystem is mounted 

```sh
sudo mount -t bpf bpf /sys/fs/bpf/
```

## Compile and load eBPF program

Compile the `C` code 

```sh
clang -O2 -g -target bpf -c trace_file_delete.c -o trace_file_delete.o -I /usr/include
```

Load the BPF Program using `bpftool`

```sh
sudo bpftool prog load trace_file_delete.o /sys/fs/bpf/trace_file_delete autoattach
sudo bpftool --json --pretty prog show |jq -r '.[] | select(.name=="trace_unlinkat")'
```

Check for kernel logs errors

```sh
dmesg | tail
```

## Run the application

Run `go` application with `sudo` privileges

```sh
sudo go run main.go
```

## Development iterations

To detach the program and remove the pinned file (prior to re-compiling `C` code)

```sh
sudo rm /sys/fs/bpf/trace_unlinkat
```