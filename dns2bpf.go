package main

import (
	"fmt"
	"os/exec"
)

func main() {
	// Prompt for input of DNS name
	var dnsName string
	fmt.Print("Enter DNS name: ")
	fmt.Scanln(&dnsName)

	// Generate BPF bytecode
	bpfBytecode, err := exec.Command("bpf_asm", "-d", fmt.Sprintf("bpf_dns_match_name %s", dnsName)).Output()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Load BPF bytecode into kernel
	err = exec.Command("bpf", "-w", "-d", string(bpfBytecode)).Run()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Create IPtables rule
	err = exec.Command("iptables", "-A", "INPUT", "-p", "udp", "--dport", "53", "-m", "bpf", "--bpf", string(bpfBytecode), "-j", "DROP").Run()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("BPF bytecode and IPtables rule created successfully")
}

