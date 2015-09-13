package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
)

var flags = struct {
	verbose    bool
	configFile string
}{}

func main() {
	//vm start name
	//vm stop name
	//vm img name size
	//vm start-all
	//vm stop-all
	//vm help
	flagParser := flag.NewFlagSet("", flag.ExitOnError)

	flagParser.BoolVar(&flags.verbose, "v", false, "verbose option")
	flagParser.StringVar(&flags.configFile, "c", "config.yaml", "search `directory` for include files")
	flagParser.Int("test", 0, "set `test`configration")

	flagParser.Parse(os.Args[1:])

	switch command := flagParser.Arg(0); command {
	case "start":
		if flagParser.Arg(1) == "" {
			help()
			return
		}
		start(flagParser.Arg(1))
	case "img":
		createImage(flagParser.Arg(1), flagParser.Arg(2))
	case "":
		help()
	default:
		fmt.Println("Not Implement", command)
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

	//kvm, lookErr := exec.LookPath("echo")
	kvm, lookErr := exec.LookPath("qemu-kvm")
	if lookErr != nil {
		fmt.Println("command not found : qemu-kvm")
		return
	}

	fmt.Println("use kvm in", kvm)

	config := readConfig(flags.configFile)

	cmd := exec.Command(kvm, config.getConfigArgs(name)...)

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
	fmt.Println("  vm start <name>")
	fmt.Println("  vm stop <name>")
	fmt.Println("  vm img <image-name> <size>")
}
