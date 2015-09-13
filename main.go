package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

func main() {
	//vm start name
	//vm stop name
	//vm img name size
	//vm start-all
	//vm stop-all
	//vm help

	if len(os.Args) < 2 {
		help()
		return
	}

	switch command := os.Args[1]; command {
	case "start":
		if len(os.Args) < 3 {
			help()
			return
		}
		start(os.Args[2])
	case "img":
		createImage(os.Args[2], os.Args[3])
	default:
		fmt.Println("Not Implement")
	}
}

func start(name string) {
	//usr/local/bin/qemu-kvm
	//  -m 8G -smp 2
	//  -drive file=bluemir-windows.img,if=virtio,media=disk
	//  -drive file=aaad,media=cdrom,index=2
	//  -daemonize
	//  -vga std
	//  -net nic,macaddr=DA:ED:DE:EF:0F:05
	//  -net tap,ifname=vnet5,script=/etc/kvm-ifup,downscript=/etc/kvm-ifdown
	kvm, lookErr := exec.LookPath("qemu-kvm")
	//kvm, lookErr := exec.LookPath("qemu-kvm")
	if lookErr != nil {
		fmt.Println("command not found : qemu-kvm")
		return
	}

	fmt.Println("use kvm in", kvm)

	config := readConfig("./config.yaml")

	vmconfig, ok := config.VM[name]
	if !ok {
		fmt.Println("config not exist!")
		os.Exit(-1)
	}

	cmd := exec.Command(kvm, vmconfig.makeArgs()...)

	var out, errout bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errout

	err := cmd.Run()
	if err != nil {
		fmt.Println("run fail", err)
		fmt.Println(errout.String())
		return
	}
	fmt.Println(out.String())

}
func createImage(name, size string) {
	//qemu-img create -f qcow2 -o size=200G bluemir-windows.img
	imgtool, lookErr := exec.LookPath("qemu-img")
	if lookErr != nil {
		fmt.Println("command not found : qemu-img")
		return
	}

	fmt.Println("use qemu-img in", imgtool)

	cmd := exec.Command(imgtool, "create", "-f", "qcow2", "-o", "size="+size, name)

	var out, errout bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errout

	err := cmd.Run()
	if err != nil {
		fmt.Println("create fail", err)
		fmt.Println(errout.String())
		return
	}
	fmt.Println(out.String())
}
func help() {
	fmt.Println("Usage:")
	fmt.Println("\tvm start <name>")
	fmt.Println("\tvm stop <name>")
	fmt.Println("\tvm img <image-name> <size>")
}
