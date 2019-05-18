package openstack

import (
	"fmt"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"cm-stack/src/config"
)



type NetWorks struct {
	// UUID for the network
	ID string `json:"id"`

	// Human-readable name for the network. Might not be unique.
	Name string `json:"name"`

	// The administrative state of network. If false (down), the network does not
	// forward packets.
	AdminStateUp bool `json:"admin_state_up"`

	// Indicates whether network is currently operational. Possible values include
	// `ACTIVE', `DOWN', `BUILD', or `ERROR'. Plug-ins might define additional
	// values.
	Status string `json:"status"`

	// Subnets associated with this network.
	Subnets []string `json:"subnets"`

	// TenantID is the project owner of the network.
	TenantID string `json:"tenant_id"`

	// ProjectID is the project owner of the network.
	ProjectID string `json:"project_id"`

	// Specifies whether the network resource can be accessed by any tenant.
	Shared bool `json:"shared"`

	// Availability zone hints groups network nodes that run services like DHCP, L3, FW, and others.
	// Used to make network resources highly available.
	AvailabilityZoneHints []string `json:"availability_zone_hints"`

	// Tags optionally set via extensions/attributestags
	Tags []string `json:"tags"`
}



func (network *NetWorks)List() (allNetwork []NetWorks) {
	//var allInstance []Server
	fmt.Println(config.Conf.OpenStackConfig.User)

	opts := gophercloud.AuthOptions{
		IdentityEndpoint: config.Conf.OpenStackConfig.IdentityEndpoint,
		Username:         config.Conf.OpenStackConfig.User,
		Password:         config.Conf.OpenStackConfig.PassWord,
		DomainName:       config.Conf.OpenStackConfig.DomainName,
		TenantName:       config.Conf.OpenStackConfig.TenantName,
	}

	provider, err := openstack.AuthenticatedClient(opts)
	if err != nil {
		fmt.Printf("AuthenticatedClient : %v", err)
	}

	client, err := openstack.NewNetworkV2(provider, gophercloud.EndpointOpts{
		Region: config.Conf.OpenStackConfig.Region,
		Name:   "neutron",
	})

	allPages, err := networks.List(client, networks.ListOpts{}).AllPages()

	ns, _ := networks.ExtractNetworks(allPages)

	for _, a := range ns {
		var nets NetWorks

		nets.Name = a.Name
		nets.Status = a.Status
		nets.AdminStateUp = a.AdminStateUp
		nets.Subnets = a.Subnets
		nets.ID = a.ID
		nets.Shared = a.Shared

		allNetwork = append(allNetwork, nets)
	}

    return allNetwork

}