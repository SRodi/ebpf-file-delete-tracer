// main.go
package main

import (
    "log"
    "os"
    "os/signal"
    "syscall"

    "github.com/cilium/ebpf"
    "github.com/cilium/ebpf/link"
    "golang.org/x/sys/unix"
)

const bpfProgPath = "./prevent_delete.o"

func main() {
    // Set the RLIMIT_MEMLOCK resource limit
    var rLimit unix.Rlimit
    rLimit.Cur = unix.RLIM_INFINITY
    rLimit.Max = unix.RLIM_INFINITY
    if err := unix.Setrlimit(unix.RLIMIT_MEMLOCK, &rLimit); err != nil {
        log.Fatalf("Failed to set RLIMIT_MEMLOCK: %v", err)
    }
		
    // Load the compiled eBPF program from ELF
    spec, err := ebpf.LoadCollectionSpec(bpfProgPath)
    if err != nil {
        log.Fatalf("Failed to load eBPF program: %v", err)
    }

    // Create a new eBPF Collection
    coll, err := ebpf.NewCollection(spec)
    if err != nil {
        log.Fatalf("Failed to create eBPF collection: %v", err)
    }
    defer coll.Close()

    // Find the eBPF program by its section name
    prog := coll.Programs["prevent_unlinkat"]
    if prog == nil {
        log.Fatalf("Program 'prevent_unlinkat' not found")
    }

    // Attach the eBPF program to the unlinkat syscall tracepoint
    tp, err := link.Tracepoint("syscalls", "sys_enter_unlinkat", prog, nil)
    if err != nil {
        log.Fatalf("Failed to attach eBPF program: %v", err)
    }
    defer tp.Close()

    log.Println("eBPF program attached. Press Ctrl+C to exit.")
    
    // Wait for a signal to exit
    sig := make(chan os.Signal, 1)
    signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
    <-sig
    log.Println("Exiting program")
}

