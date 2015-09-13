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
    example:
        cpu :
            core : 1
        memory : 512M
        disks :
            - path: example-hda.img
              interface: virtio
            - path: example-hdb.img
        network :
            ifname : vnet1
            mac: DE:AD:BE:EF:00:01
        logfile: test.log
        cdroms :
            - example-install.iso
        vnc : :10
        others :
            - -usb
            - -name
            - example
    example-non-network:
        cpu :
            type : host
            core : 2
            virtualcore : 1
        memory : 2G
```

To make disk image file, excute below line
```
vm-manager img example.img 50G
```

Then, you can check image file maked, and you can run your vm :
```
vm-manager start example
```

