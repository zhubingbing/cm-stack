package openstack

import (
	"cm-stack/src/config"
	"fmt"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/volumes"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/projects"
	"sort"
	"strconv"

	az "github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/availabilityzones"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/hypervisors"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/keypairs"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/volumeattach"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/flavors"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/images"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/services"
	"log"
	"time"
)




//ListServers

type Image struct {
	Name  string `json:"name"`
	ID    string `json:"id"`
	Links []struct {
		Href string `json:"href"`
		Rel  string `json:"rel"`
	} `json:"links"`
}


type Tenant struct {
	// ID is a unique identifier for this tenant.
	ID string `json:"id"`

	// Name is a friendlier user-facing name for this tenant.
	Name string `json:"name"`

	// Description is a human-readable explanation of this Tenant's purpose.
	Description string `json:"description"`

	// Enabled indicates whether or not a tenant is active.
	Enabled bool `json:"enabled"`
}


type GetInstance struct {
	Instance struct {
		Name      string   `json:"name"`
		Network   string   `json:"network"`
		Flavor    string   `json:"flavor"`
		Subnet    string   `json:"subnet"`
		Port      string   `json:"port"`
		Key       string   `json:"key"`
		Zone      string   `json:"zone"`
		Networkid []string `json:"networkid"`
		Image     string   `json:"image"`
		Count     string   `json:"count"`
	} `json:"instance"`
}


type CreateServerOpts struct {
	// Name is the name to assign to the newly launched server.
	Name string `json:"name" required:"true"`

	// ImageRef [optional; required if ImageName is not provided] is the ID or
	// full URL to the image that contains the server's OS and initial state.
	// Also optional if using the boot-from-volume extension.
	ImageRef string `json:"imageRef"`

	// ImageName [optional; required if ImageRef is not provided] is the name of
	// the image that contains the server's OS and initial state.
	// Also optional if using the boot-from-volume extension.
	ImageName string `json:"-"`

	// FlavorRef [optional; required if FlavorName is not provided] is the ID or
	// full URL to the flavor that describes the server's specs.
	FlavorRef string `json:"flavorRef"`

	// FlavorName [optional; required if FlavorRef is not provided] is the name of
	// the flavor that describes the server's specs.
	FlavorName string `json:"-"`

	// SecurityGroups lists the names of the security groups to which this server
	// should belong.
	SecurityGroups []string `json:"-"`

	// UserData contains configuration information or scripts to use upon launch.
	// Create will base64-encode it for you, if it isn't already.
	UserData []byte `json:"-"`

	// AvailabilityZone in which to launch the server.
	AvailabilityZone string `json:"availability_zone,omitempty"`

	// Networks dictates how this server will be attached to available networks.
	// By default, the server will be attached to all isolated networks for the
	// tenant.
	Networks []servers.Network `json:"-"`

	// Metadata contains key-value pairs (up to 255 bytes each) to attach to the
	// server.
	Metadata map[string]string `json:"metadata,omitempty"`

	// ConfigDrive enables metadata injection through a configuration drive.
	ConfigDrive *bool `json:"config_drive,omitempty"`

	// AdminPass sets the root user password. If not set, a randomly-generated
	// password will be created and returned in the response.
	AdminPass string `json:"adminPass,omitempty"`

	// AccessIPv4 specifies an IPv4 address for the instance.
	AccessIPv4 string `json:"accessIPv4,omitempty"`

	// AccessIPv6 specifies an IPv6 address for the instance.
	AccessIPv6 string `json:"accessIPv6,omitempty"`

	// Min specifies Minimum number of servers to launch.
	Min int `json:"min_count,omitempty"`

	// Max specifies Maximum number of servers to launch.
	Max int `json:"max_count,omitempty"`

	// ServiceClient will allow calls to be made to retrieve an image or
	// flavor ID by name.
	ServiceClient *gophercloud.ServiceClient `json:"-"`

	// Tags allows a server to be tagged with single-word metadata.
	// Requires microversion 2.52 or later.
	Tags []string `json:"tags,omitempty"`
}


// VolumeAttachment contains attachment information between a volume
type VolumeAttachment struct {
	// ID is a unique id of the attachment.
	ID string `json:"id"`

	//Name is VolumeAttachment device name
	Name string `json:"name"`

	// Device is what device the volume is attached as.
	Device string `json:"device"`

	// VolumeID is the ID of the attached volume.
	VolumeID string `json:"volumeId"`

	// ServerID is the ID of the instance that has the volume attached.
	ServerID string `json:"serverId"`

	// Current status of the volume.
	Status string `json:"status"`


	// Instances onto which the volume is attached.
	Attachments []map[string]interface{} `json:"attachments"`
	// This parameter is no longer used.
	AvailabilityZone string `json:"availability_zone"`


	// Indicates whether this is a bootable volume.
	Bootable string `json:"bootable"`

	// The date when this volume was created.
	CreatedAt time.Time `json:"-"`

	// Human-readable description for the volume.
	Description string `json:"display_description"`

	// The type of volume to create, either SATA or SSD.
	VolumeType string `json:"volume_type"`

	// The ID of the snapshot from which the volume was created
	SnapshotID string `json:"snapshot_id"`

	// The ID of another block storage volume from which the current volume was created
	SourceVolID string `json:"source_volid"`

	// Arbitrary key-value pairs defined by the user.
	Metadata map[string]string `json:"metadata"`

	// Size of the volume in GB.
	Size int `json:"size"`
}

// Flavor records represent (virtual) hardware configurations for server resources in a region.
type Flavor struct {
	ID         string   `json:"id"`
	Disk       int      `json:"disk"`
	RAM        int      `json:"ram"`
	Name       string   `json:"name"`
	RxTxFactor float64  `json:"rxtx_factor"`
	Swap       int      `json:"-"`
	VCPUs      int      `json:"vcpus"`
	IsPublic   bool     `json:"os-flavor-access:is_public"`
	Ephemeral  int      `json:"OS-FLV-EXT-DATA:ephemeral"`
}

type Flavors []Flavor

// Len()方法和Swap()方法不用变化
// 获取此 slice 的长度
func (p Flavors) Len() int { return len(p) }

// 交换数据
func (p Flavors) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

type SortByCPU struct{ Flavors }

type SortByRAM struct{ Flavors }

// ascending order by vcpus
func (p SortByCPU) Less(i, j int) bool {
	return p.Flavors[i].VCPUs < p.Flavors[j].VCPUs
}

// ascending order by ram
func (p SortByRAM) Less(i, j int) bool {
	return p.Flavors[i].RAM < p.Flavors[j].RAM
}


type Server struct {
	// ID uniquely identifies this server amongst all other servers, including those not accessible to the current tenant.
	ID string

	// TenantID identifies the tenant owning this server resource.
	TenantID string `mapstructure:"tenant_id"`

	// UserID uniquely identifies the user account owning the tenant.
	UserID string `mapstructure:"user_id"`

	// Name contains the human-readable name for the server.
	Name string

	// Updated and Created contain ISO-8601 timestamps of when the state of the server last changed, and when it was created.
	Updated time.Time
	Created time.Time

	HostID string

	// Status contains the current operational status of the server, such as IN_PROGRESS or ACTIVE.
	Status string

	// Progress ranges from 0..100.
	// A request made against the server completes only once Progress reaches 100.
	Progress int

	// AvailabilityZone in which to launch the server.
	AvailabilityZone string `json:"availability_zone,omitempty"`

	// AccessIPv4 and AccessIPv6 contain the IP addresses of the server, suitable for remote access for administration.
	AccessIPv4, AccessIPv6 string

	// Image refers to a JSON object, which itself indicates the OS image used to deploy the server.
	Images Image

	// Project refers to a JSON object, which itself indicates the Tenant to deploy the server.

	Project projects.Project

	//VolumeAttachment refers to a JSON object
	Volume []VolumeAttachment

	// Flavor refers to a JSON object, which itself indicates the hardware configuration of the deployed server.
	Flavor map[string]interface{}

	// Addresses includes a list of all IP addresses assigned to the server, keyed by pool.
	Addresses map[string]interface{}

	// Metadata includes a list of all user-specified key-value pairs attached to the server.
	Metadata map[string]interface{}

	// Links includes HTTP references to the itself, useful for passing along to other APIs that might want a server reference.
	Links []interface{}

	// KeyName indicates which public key was injected into the server on launch.
	KeyName string `json:"key_name" mapstructure:"key_name"`

	// AdminPass will generally be empty ("").  However, it will contain the administrative password chosen when provisioning a new server without a set AdminPass setting in the first place.
	// Note that this is the ONLY time this field will be valid.
	AdminPass string `json:"adminPass" mapstructure:"adminPass"`

	// SecurityGroups includes the security groups that this instance has applied to it
	SecurityGroups []map[string]interface{} `json:"security_groups" mapstructure:"security_groups"`
}


func GetIDFromeInterface(a map[string]interface{}) string {

	var id string
	for _, v := range a {

		switch vv := v.(type) {
		case string:
			id = vv

		case nil:
				///fmt.Println(k, "is nil", "null")

		default:
				//fmt.Println(k, "is of a type I don't know how to handle ", fmt.Sprintf("%T", v))
		}
	}

	return id
}

func (instance *Server)GetToken() (*gophercloud.ProviderClient, error) {

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
		log.Fatalln(err)
	}

    return provider, nil

}


func (instance *Server)ListServers() (allInstance []Server) {

	//var allInstance []Server
	fmt.Println(config.Conf.OpenStackConfig.IdentityEndpoint)


	opts := gophercloud.AuthOptions{
		IdentityEndpoint: config.Conf.OpenStackConfig.IdentityEndpoint,
		Username:         config.Conf.OpenStackConfig.User,
		Password:         config.Conf.OpenStackConfig.PassWord,
		DomainName:       config.Conf.OpenStackConfig.DomainName,
		TenantName:       config.Conf.OpenStackConfig.TenantName,

	}
	fmt.Println(opts)

	provider, err := openstack.AuthenticatedClient(opts)
	if err != nil {
		fmt.Printf("AuthenticatedClient : %v", err)
		log.Fatalln(err)
	}

	client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{
		Region: config.Conf.OpenStackConfig.Region,
	})
	if err != nil {
		fmt.Printf("AuthenticatedClient : %v", err)
		log.Fatalln(err)
	}


	test, err := services.List(client).AllPages()
	if err != nil {
		panic(err)
	}

	allServices, err := services.ExtractServices(test)
	if err != nil {
		panic(err)
	}

	for _, service := range allServices {
		fmt.Printf("%+v\n", service)
	}



	identityClient, err := openstack.NewIdentityV3(provider,gophercloud.EndpointOpts{})
	if err != nil {
		fmt.Printf("NewIdentityV3 : %v", err)
		log.Fatalln(err)
	}


	blockstorageClient , err := openstack.NewBlockStorageV3(provider, gophercloud.EndpointOpts{})
	if err != nil {
		fmt.Printf("NewBlockStorageV3 : %v", err)
		log.Fatalln(err)
	}


	listOps := servers.ListOpts{}
	listOps.AllTenants = true

	allPages, err := servers.List(client, listOps).AllPages()
	if err != nil {
		log.Fatalln(err)
	}
	//解析返回值
	allServers, err := servers.ExtractServers(allPages)
	if err != nil {
		log.Fatalln(err)
	}

	for _, s := range allServers {
		var b Server

		b.ID = s.ID
		b.Name = s.Name
		b.UserID = s.UserID
		b.Created = s.Created
		b.Updated = s.Updated
		b.AccessIPv4 = s.AccessIPv4
		b.Addresses = s.Addresses
		b.Status = s.Status
		b.Flavor = s.Flavor
		b.KeyName = s.KeyName
		b.AdminPass = s.AdminPass
		b.HostID = s.HostID
		b.SecurityGroups = s.SecurityGroups
		b.Images.ID = GetIDFromeInterface(s.Image)
		b.TenantID = s.TenantID

		client, _ := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{
			Region: "RegionOne",
		})

		b.Images = instance.GetImage(client, b.Images.ID)

		b.Volume = instance.GetVolumeattach(client,blockstorageClient, b.ID)

		b.Project = instance.GetProject(identityClient, b.TenantID)


		allInstance = append(allInstance, b)

	}

	return allInstance
	
}


func (instance *Server)GetImage(client *gophercloud.ServiceClient, id string) (image Image) {
	im ,err:= images.Get(client, id).Extract()
	if err != nil {
		fmt.Printf("images.Get : %v", err)
		log.Fatalln(err)
	}

	image.Name = im.Name
	image.ID = im.ID

	return image
}

func (instance *Server)GetProject(identityClient *gophercloud.ServiceClient, id string) (project projects.Project) {

	pro , err := projects.Get(identityClient, id).Extract()

	if err != nil {
		fmt.Printf("domains.Get : %v", err)
		log.Fatalln(err)
	}

	project.Name = pro.Name
	project.Enabled = pro.Enabled
	project.Description = pro.Description
	project.ID = pro.ID
	project.DomainID = pro.DomainID
	project.IsDomain = pro.IsDomain
	project.ParentID = pro.ParentID

	return project

}




func (instance *Server)GetVolumeattach(client *gophercloud.ServiceClient, blockstorageClient *gophercloud.ServiceClient, serverId string) (volume []VolumeAttachment){

	allPages, _ := volumeattach.List(client, serverId).AllPages()

	allVolumes, _ := volumeattach.ExtractVolumeAttachments(allPages)


	for _, v := range allVolumes {
		var vo VolumeAttachment

		vo.ID = v.ID
		vo.Device = v.Device
		vo.VolumeID = v.VolumeID

		volumeInfo, err := volumes.Get(blockstorageClient,v.ID).Extract()
		if err != nil {
			fmt.Printf("volumes.Get : %v", err)
			log.Fatalln(err)
		}

		vo.Name = volumeInfo.Name
		vo.Status = volumeInfo.Status
		vo.Size = volumeInfo.Size
		vo.AvailabilityZone = volumeInfo.AvailabilityZone
		vo.CreatedAt = volumeInfo.CreatedAt
		vo.Description = volumeInfo.Description
		vo.VolumeType = volumeInfo.VolumeType

		volume = append(volume, vo)
	}



	return volume
}


func (instance *Server)ListFlavor() (allFlavor []Flavor) {

	provider,err := instance.GetToken()
	if err != nil {
		fmt.Printf("AuthenticatedClient : %v", err)
		log.Fatalln(err)
	}

	client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{
		Region: config.Conf.OpenStackConfig.Region,
	})
	if err != nil {
		fmt.Printf("AuthenticatedClient : %v", err)
	    log.Fatalln(err)
    }

	allPages, _ := flavors.ListDetail(client, nil).AllPages()
	allFlavors, _ := flavors.ExtractFlavors(allPages)

	for _, flavor := range allFlavors {
		var f Flavor

		f.Name = flavor.Name
		f.ID = flavor.ID
		f.RAM = flavor.RAM
		f.Disk = flavor.Disk
		f.IsPublic = flavor.IsPublic
		f.VCPUs = flavor.VCPUs
		f.Ephemeral = flavor.Ephemeral
		f.Swap = flavor.Swap

		allFlavor = append(allFlavor, f)
	}

	sort.Sort(SortByRAM{allFlavor})
	sort.Sort(SortByCPU{allFlavor})

	return allFlavor
}

func (instance *Server)ListKey() (allKey []keypairs.KeyPair) {

	provider,err := instance.GetToken()
	if err != nil {
		fmt.Printf("AuthenticatedClient : %v", err)
		log.Fatalln(err)
	}

	client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{
		Region: config.Conf.OpenStackConfig.Region,
	})
	if err != nil {
		fmt.Printf("AuthenticatedClient : %v", err)
		log.Fatalln(err)
	}

	allPages, _ := keypairs.List(client).AllPages()
	allKeys, _ := keypairs.ExtractKeyPairs(allPages)

	for _, keys := range allKeys {
		var f keypairs.KeyPair

		f.Name = keys.Name
		f.UserID = keys.UserID
		f.PublicKey = keys.PublicKey
		f.PrivateKey = keys.PrivateKey

		allKey = append(allKey, f)
	}


	return allKey
}

func (instance *Server)ListZone() (allzone []az.AvailabilityZone) {

	provider,err := instance.GetToken()
	if err != nil {
		fmt.Printf("AuthenticatedClient : %v", err)
		log.Fatalln(err)
	}

	client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{
		Region: config.Conf.OpenStackConfig.Region,
	})
	if err != nil {
		fmt.Printf("AuthenticatedClient : %v", err)
		log.Fatalln(err)
	}

	allPages, _ := az.List(client).AllPages()
	allZone, _ := az.ExtractAvailabilityZones(allPages)

	for _, zone := range allZone {
		var f az.AvailabilityZone

		f.ZoneName = zone.ZoneName
		f.Hosts = zone.Hosts
		f.ZoneState = zone.ZoneState

		allzone = append(allzone, f)
	}

	return allzone
}


func (instance *Server)ListImage() (allimage []images.Image) {

	provider,err := instance.GetToken()
	if err != nil {
		fmt.Printf("AuthenticatedClient : %v", err)
		log.Fatalln(err)
	}

	client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{
		Region: config.Conf.OpenStackConfig.Region,
	})
	if err != nil {
		fmt.Printf("AuthenticatedClient : %v", err)
		log.Fatalln(err)
	}

	listOpts := images.ListOpts{
		Status: "active",
	}

	allPages, _ := images.ListDetail(client, listOpts).AllPages()
	allimages, _ := images.ExtractImages(allPages)

	for _, image := range allimages {
		var f images.Image

		f.Name = image.Name
		f.ID = image.ID
		f.Status = image.Status
		f.Created = image.Created
		f.Updated = image.Updated
		f.MinDisk = image.MinDisk
		f.Metadata = image.Metadata
		f.Progress = image.Progress
		f.MinRAM = image.MinRAM

		allimage = append(allimage, f)
	}

	return allimage
}

func (instance *Server)CreateServers(opts *GetInstance) {

	provider,err := instance.GetToken()
	if err != nil {
		fmt.Printf("AuthenticatedClient : %v", err)
		log.Fatalln(err)
	}

	client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{
		Region: config.Conf.OpenStackConfig.Region,
	})
	if err != nil {
		fmt.Printf("AuthenticatedClient : %v", err)
		log.Fatalln(err)
	}


	c,_:=strconv.Atoi(opts.Instance.Count)


	CreateOpts := servers.CreateOpts{
		Name:         opts.Instance.Name,
		FlavorRef:    opts.Instance.Flavor,
		ImageRef:     opts.Instance.Image,
		Networks:[]servers.Network{
			{UUID:opts.Instance.Networkid[0]},
		},
		AvailabilityZone: opts.Instance.Zone,
		Min: c,
		Max: c,

	}


	ps := keypairs.CreateOptsExt{
		CreateOptsBuilder: CreateOpts,
		KeyName:             opts.Instance.Key,
	}
	allPages, err := servers.Create(client, ps).Extract()
	if err != nil {
		log.Fatalln(err)
	}

	//解析返回值
	fmt.Println(allPages)

	//return  allPages

}

func (instance *Server)List_Hypervisors() (allhypervisors []hypervisors.Hypervisor) {
	provider,err := instance.GetToken()
	if err != nil {
		fmt.Printf("AuthenticatedClient : %v", err)
		log.Fatalln(err)
	}

	client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{
		Region: config.Conf.OpenStackConfig.Region,
	})
	if err != nil {
		fmt.Printf("AuthenticatedClient : %v", err)
		log.Fatalln(err)
	}

	allPages, err := hypervisors.List(client).AllPages()
	if err != nil {
		log.Fatalln(err)
	}

	allHypervisors, err := hypervisors.ExtractHypervisors(allPages)
	if err != nil {
		log.Fatalln(err)
	}
	for _, hypervisor := range allHypervisors {
		var a hypervisors.Hypervisor
		a = hypervisor

		allhypervisors = append(allhypervisors, a)

	}

	return allhypervisors
}

func (instance *Server)GetServer(id string) (r *servers.Server) {

	provider,err := instance.GetToken()
	if err != nil {
		config.Log.Errorf("openstack Compute get token error")
	}

	client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{
		Region: config.Conf.OpenStackConfig.Region,
	})

	fmt.Println(config.Conf.OpenStackConfig)
	if err != nil {
		config.Log.Errorf("openstack Compute AuthenticatedClient error")
	}

	actual, err := servers.Get(client, id).Extract()
	if err != nil {
		config.Log.Errorf("Unexpected Get instance error")
	}

	return actual
}

func (instance *Server)GetServerCpuInfo(id string) (r *servers.Server) {

	provider,err := instance.GetToken()
	if err != nil {
		fmt.Printf("AuthenticatedClient : %v", err)
		log.Fatalln(err)
	}

	client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{
		Region: config.Conf.OpenStackConfig.Region,
	})
	if err != nil {
		fmt.Printf("AuthenticatedClient : %v", err)
		log.Fatalln(err)
	}

	actual, err := servers.Get(client, id).Extract()
	if err != nil {
		log.Fatalln("Unexpected Get error: %v", err)
	}

	return actual
}