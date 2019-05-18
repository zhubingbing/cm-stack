package main

import (
	"flag"
	"fmt"
	"github.com/digitalocean/go-qemu/qemu"
	"log"
	"net"
	"time"

	hypervisor "github.com/digitalocean/go-qemu/hypervisor"
)

var (
	network = flag.String("network", "unix", `Named network used to connect on. For Unix sockets -network=unix, for TCP connections: -network=tcp`)
	address = flag.String("address", "/var/run/libvirt/libvirt-sock", `Address of the hypervisor. This could be in the form of Unix or TCP sockets. For TCP connections: -address="host:16509"`)
	timeout = flag.Duration("timeout", 2*time.Second, "Connection timeout. Another valid value could be -timeout=500ms")
)

func main() {
	flag.Parse()

	fmt.Printf("\nConnecting to %s://%s\n", *network, *address)
	newConn := func() (net.Conn, error) {
		return net.DialTimeout(*network, *address, *timeout)
	}

	driver := hypervisor.NewRPCDriver(newConn)

	hv := hypervisor.New(driver)

	fmt.Printf("\n**********Domains**********\n")
	domains, err := hv.Domains()
	if err != nil {
		log.Fatalf("Unable to get domains from hypervisor: %v", err)
	}
	for _, dom := range domains {
		fmt.Printf("%s\n", dom.Name)

		displayPCIDevices(dom)
		displayBlockDevices(dom)
	}
	fmt.Printf("\n***************************\n")

}


func displayPCIDevices(domain *qemu.Domain) {

	fmt.Println(domain.Events())
	pciDevices, err := domain.PCIDevices()
	if err != nil {
		log.Fatalf("Error getting PCIDevices: %v\n", pciDevices)
	}
	fmt.Printf("\n[ PCIDevices ]\n")
	fmt.Printf("======================================\n")
	fmt.Printf("%10s %20s\n", "[ID]", "[Description]")
	fmt.Printf("======================================\n")
	for _, pciDevice := range pciDevices {
		fmt.Printf("[%10s] [%20s]\n", pciDevice.QdevID, pciDevice.ClassInfo.Desc)
	}
}

func displayBlockDevices(domain *qemu.Domain) {
	blockDevices, err := domain.BlockDevices()
	if err != nil {
		log.Fatalf("Error getting blockDevices: %v\n", blockDevices)
	}
	fmt.Printf("\n[ BlockDevices ]\n")
	fmt.Printf("========================================================================\n")
	fmt.Printf("%20s %8s %30s\n", "Device", "Driver", "File")
	fmt.Printf("========================================================================\n")
	for _, blockDevice := range blockDevices {
		fmt.Printf("%20s %8s %30s\n",
			blockDevice.Device, blockDevice.Inserted.Driver, blockDevice.Inserted.File)
	}
}
