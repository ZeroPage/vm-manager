#VM Manager

This tool based on [Qemu & KVM guide #1](http://www.slideshare.net/baramdori/qemu-kvm-guide-1-intro-basic)

##How to install & use

First, install go

https://golang.org/dl

Then, set $GOPATH and $PATH

```
# append below lines to .bashrc or .zshrc or other shell profile
export $GOPATH=your-go-code-path
export $PATH=$GOPATH/bin:$PATH
```
(and you may re-login shell)

Then, install vm-manager

```
go get github.com/zeropage/vm-manager
```

And make config file like below.
```
#This is example file
#config.yaml

vm:
	- example-vm
		cpu : 1
		memory : 2G
		disk :
			path : example.img
		network :
			mac : DE:AD:BE:EF:00:01
			ifname : vnet1
		cdrom : your-os-install-image-file.iso
```

To make disk image file, excute below line
```
vm-manager img example.img 50G
```

Then, you can check image file maked, and you can run your vm :
```
vm-manager start example-vm
```

