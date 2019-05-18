package main

import (
	"fmt"
	"github.com/levigross/grequests"
	"github.com/libvirt/libvirt-go"
	"github.com/libvirt/libvirt-go-xml"
	"github.com/prometheus/client_golang/api"
	//"github.com/prometheus/client_golang/prometheus/promhttp"
	//promv1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"net/http"
	"strings"
)



func main()  {
	conn, err := libvirt.NewConnect("qemu+tcp://172.21.21.121:16509/system")
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()


	cfg := api.Config{
		Address:      "http://172.21.21.200:9091/api/v1/label/job/values",
		RoundTripper: &http.Transport{},
	}


	resp, _:= grequests.Get(cfg.Address, nil)


	fmt.Println(resp)







	hostname,_ := conn.GetHostname()

	fmt.Println(hostname)
	//libvirt.CONNECT_LIST_DOMAINS_ACTIVE

	doms, err := conn.ListAllDomains(libvirt.CONNECT_LIST_DOMAINS_RUNNING)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%d running domains:\n", len(doms))
	for _, dom := range doms {

		id, err := dom.GetID()



		name, err := dom.GetName()


		a , _:= dom.GetInfo()
		b ,_:= dom.MemoryStats( uint32(id), 0)

		fmt.Println(b[0].Tag)

		hostname ,_ := dom.GetHostname(0)



		fmt.Println(hostname)

		fmt.Println(a.Memory, a.State, a.MaxMem, a.NrVirtCpu, a.CpuTime)
		if err == nil {
			fmt.Printf("sss %s , %s\n",name ,id)
		}

		xml , _:= dom.GetXMLDesc(0)


		//fmt.Println(xml)
		domcfg := &libvirtxml.Domain{}

		err = domcfg.Unmarshal(xml)


		if err != nil {
			fmt.Printf("%s ",err)
		}



		if err != nil {
			fmt.Printf("%s ",err)
		}


		fmt.Println(domcfg.Devices.Interfaces[0].Target.Dev)


		fmt.Println(domcfg.Devices.Disks[0].Target.Bus)

		fmt.Println("debug----")
		fmt.Println(domcfg.Devices.Disks[0].Source.Network.Name)

		s := strings.Split(domcfg.Devices.Disks[0].Source.Network.Name, "/")
		fmt.Println(s, len(s))
		//fmt.Println(s[0])
		//fmt.Println(s[1])


		if strings.HasPrefix(s[1], "volume-") {
			fmt.Println(strings.TrimPrefix(s[1], "volume-"))
		} else if strings.HasSuffix(s[1], "_disk") {
			fmt.Println(strings.TrimSuffix(s[1], "_disk"))
		} else if strings.HasSuffix(s[1], "_disk.config"){
			fmt.Println(strings.TrimSuffix(s[1], "_disk.config"))
		}


		fmt.Println(domcfg.Devices.Interfaces[0].Target.Dev)


		//fmt.Println(xml)

		ss ,_:= dom.BlockStats("sda")

		fmt.Println(ss.WrBytes, ss.RdBytes)

		cc,_ := dom.GetBlockInfo("sda", 0)

		fmt.Println(cc.Allocation, cc.Capacity, cc.Physical)

		network ,_ := dom.InterfaceStats(domcfg.Devices.Interfaces[0].Target.Dev)

		fmt.Println(network.RxBytes, network.RxBytesSet, network.RxErrs)

		//fmt.Println(xml)
		//dom.ManagedSaveGetXMLDesc()
		dom.Free()
	}





}