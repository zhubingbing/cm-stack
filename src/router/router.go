package router

import (
	"cm-stack/src/service/openstack"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/api"
	"net/http"

	"github.com/prometheus/client_golang/api/prometheus/v1"
)



func Load(g *gin.Engine, mw ...gin.HandlerFunc) *gin.Engine {
	// Middlewares.
	g.Use(gin.Logger())
	g.Use(gin.Recovery())
	//g.Use(middleware.NoCache)
	//g.Use(middleware.Options)
	//g.Use(middleware.Secure)
	//g.Use(mw...)

	//g.GET("/version", versionHandler)
	//g.GET("/", rootHandler)

	// 404 Handler.
	g.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "不存在的接口地址.")
	})

	// User API
	OpenStack_Hypervisors := g.Group("/openstack/os-hypervisors")
	{
		// u.POST("", user.Create)
		//OpenStack.POST("", user.Create)
		OpenStack_Hypervisors.GET("/", ListHypervisors)
		//OpenStack_Compute.GET("/:id", .Get)
	}


	OpenStack_Compute := g.Group("/openstack/servers")
	{
		// u.POST("", user.Create)
		OpenStack_Compute.POST("/", CreateServer)

		OpenStack_Compute.GET("/", ListAll)
		OpenStack_Compute.GET("/flavor", ListFlavor)
		OpenStack_Compute.GET("/key", ListKey)
		OpenStack_Compute.GET("/zone", ListZone)
		OpenStack_Compute.GET("/image", ListImage)
		OpenStack_Compute.GET("/server/", GetServer)
		//OpenStack_Compute.GET("/:id", .Get)
	}

	OpenStack_Network := g.Group("/openstack/networks")
	{
		// u.POST("", user.Create)
		//OpenStack.POST("", user.Create)
		OpenStack_Network.GET("/", ListNetwork)
		//OpenStack_Compute.GET("/:id", .Get)
	}


	OpenStack_Image:= g.Group("/openstack/images")
	{
		// u.POST("", user.Create)
		//OpenStack.POST("", user.Create)
		OpenStack_Image.GET("/", ListAll)
		//OpenStack_Compute.GET("/:id", .Get)
	}

	Promethues_Api := g.Group("/promethues")
	{
		Promethues_Api.GET("/targets", DoTargets)
		Promethues_Api.GET("/alerts", DoAlerts)
		Promethues_Api.GET("/alert", GetAlerts)

	}


	a := g.Group("/api")
	{
		// u.POST("", user.Create)
		//OpenStack.POST("", user.Create)
		a.POST("/login", Login)

		a.OPTIONS("/login", func(c *gin.Context) {
			c.JSON(200,"OPTIONS")
		})

		a.GET("/info", GetUser)
		//OpenStack_Compute.GET("/", ListAll)
		//OpenStack_Compute.GET("/:id", .Get)
	}

	return g
}


func ListHypervisors(c *gin.Context)  {
	p := openstack.Server{}

	a := p.List_Hypervisors()

	c.JSON(http.StatusOK, gin.H{
		"data": a,
		"code": 20000,

	})
}


func ListAll(c *gin.Context)  {
	p := openstack.Server{}

	instances := p.ListServers()

	//fmt.Println(instances)

	c.JSON(http.StatusOK, gin.H{
		"data": instances,
		"code": 20000,

	})
}


func GetServer(c *gin.Context)  {
	p := openstack.Server{}

	instanceUuid := c.Query("uuid")

	instances := p.GetServer(instanceUuid)

	//fmt.Println(instances)

	c.JSON(http.StatusOK, gin.H{
		"data": instances,
		"code": 20000,

	})
}



func CreateServer(c *gin.Context)  {
	p := openstack.Server{}

	opts := &openstack.GetInstance{}
	c.Bind(opts)

	p.CreateServers(opts)

	//fmt.Println(opts.Instance.Count)


	c.JSON(http.StatusOK, gin.H{
		"data": opts,
		"code": 20000,

	})
}

func ListNetwork(c *gin.Context)  {
	p := openstack.NetWorks{}

	network := p.List()

	c.JSON(http.StatusOK, gin.H{
		"data": network,
		"code": 20000,

	})
}

func ListFlavor(c *gin.Context)  {
	p := openstack.Server{}

	flavor := p.ListFlavor()

	c.JSON(http.StatusOK, gin.H{
		"data": flavor,
		"code": 20000,

	})
}

func ListImage(c *gin.Context)  {
	p := openstack.Server{}

	image := p.ListImage()

	c.JSON(http.StatusOK, gin.H{
		"data": image,
		"code": 20000,

	})
}

func ListKey(c *gin.Context)  {
	p := openstack.Server{}

	key := p.ListKey()

	c.JSON(http.StatusOK, gin.H{
		"data": key,
		"code": 20000,

	})
}

func ListZone(c *gin.Context)  {
	p := openstack.Server{}

	zone := p.ListZone()

	c.JSON(http.StatusOK, gin.H{
		"data": zone,
		"code": 20000,

	})
}

func Login(c*gin.Context)  {


	token:= map[string]string{
		"token": "admin",
	}

	c.JSON(http.StatusOK, gin.H{
		"data": token,
		"code": 20000,

	})

}


type Role struct {
	Name string        `json:"name"`
	Roles []string     `json:"roles"`

	Avatar string      `json:"avatar"`

}


func GetUser(c*gin.Context)  {

	var role Role
	s := []string{"admin"}

	role.Name = "admin"
	role.Avatar = "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif"
	role.Roles = s


	c.JSON(http.StatusOK, gin.H{
		"code": 20000,
		"data": role,

	})

}

func DoTargets(c *gin.Context)  {
	//client := &prometheus.ApiTest{}

	cfg := api.Config{}

	cfg.Address = "http://172.21.21.200:9091"
    cfg.RoundTripper = api.DefaultRoundTripper

	client ,_ := api.NewClient(cfg)


	apiClient := v1.NewAPI(client)

	a , _:= apiClient.Targets(context.Background())

	fmt.Println(a.Active)

	c.JSON(http.StatusOK, gin.H{
		"data": a,
		"code": 20000,

	})
}

func DoAlerts(c *gin.Context)  {
	//client := &prometheus.ApiTest{}

	cfg := api.Config{}

	cfg.Address = "http://172.21.21.200:9091"
	cfg.RoundTripper = api.DefaultRoundTripper

	client ,_ := api.NewClient(cfg)

	apiClient := v1.NewAPI(client)

	a , _:= apiClient.Alerts(context.Background())

	fmt.Println(a.Alerts)

	c.JSON(http.StatusOK, gin.H{
		"data": a.Alerts,
		"code": 20000,

	})
}

func GetAlerts(c *gin.Context)  {
	//client := &prometheus.ApiTest{}

	fmt.Println(c.Request.Method)

	fmt.Println(c)

	c.JSON(http.StatusOK, gin.H{
		"data": "test",
		"code": 20000,

	})
}
