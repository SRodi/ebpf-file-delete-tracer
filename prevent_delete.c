// prevent_delete.c
#include "vmlinux.h"
#include <bpf/bpf_helpers.h>

SEC("tracepoint/syscalls/sys_enter_unlinkat")
int prevent_unlinkat(struct trace_event_raw_sys_enter* ctx) {
    // Return a negative value to indicate failure (EACCES: Permission denied)
    return -1; 
}

char _license[] SEC("license") = "GPL";

