package main

import (
	"cm-stack/src/server"
	"fmt"
	"golang.org/x/sync/errgroup"
	"gopkg.in/alecthomas/kingpin.v2"

	"cm-stack/src/config"
)

var (
	// 版本
	//Version = "v0.0.8"

	// 配置参数
	//Conf config.ConfYaml

	err error
)


func main() {
	opts := config.ConfYaml{}

	var (
		configFile           = kingpin.Flag("config", "config file").Default("./config.yaml").String()
		address              = kingpin.Flag("address", "listen address; default any").Default("").String()
		port                 = kingpin.Flag("port", "listen port; default 9098").Default("9098").String()

	)

	kingpin.HelpFlag.Short('h')
	kingpin.Parse()


	opts.Core.Address = *address
	opts.Core.Port = *port
	config.Conf, err = config.Init(*configFile)
	if err != nil {
		fmt.Printf("Load yaml config file error: '%v'", err)
		return
	}

	var g errgroup.Group

	g.Go(func() error {
		// 启动服务
		return server.RunHTTPServer()
	})

	if err = g.Wait(); err != nil {
		fmt.Println("接口服务出错了：", err)
	}

}
