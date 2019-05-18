package server

import (
	"cm-stack/src/config"
	"cm-stack/src/router"

	"crypto/tls"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/acme/autocert"
	"net/http"
)



// New returns a app instance
func New() *gin.Engine {
	// Set gin mode.
	gin.SetMode(config.Conf.Core.Mode)

	// Create the Gin engine.
	g := gin.New()

	// Routes
	router.Load(g)

	return g
}

func autoTLSServer() *http.Server {
	m := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(config.Conf.Core.AutoTLS.Host),
		Cache:      autocert.DirCache(config.Conf.Core.AutoTLS.Folder),
	}
	return &http.Server{
		Addr:      	":https",
		TLSConfig: 	&tls.Config{GetCertificate: m.GetCertificate},
		Handler:  	New(),
	}
}

func defaultTLSServer() *http.Server {
	return &http.Server{
		Addr: 			"0.0.0.0:" + config.Conf.Core.TLS.Port,
		Handler:	  New(),
	}
}

func defaultServer() *http.Server {
	return &http.Server{
		Addr: 			"0.0.0.0:" + config.Conf.Core.Port,
		Handler:	  New(),
	}
}

// RunHTTPServer provide run http or https protocol.
func RunHTTPServer() (err error) {
	if !config.Conf.Core.Enabled {
		config.Log.Debug("httpd server is disabled.")
		return nil
	} else {
		s := defaultServer()
		//handleSignal(s)
		config.Log.Infof("3. Start to listening the incoming requests on http address: %s", config.Conf.Core.Port)
		err = s.ListenAndServe()
	}

	return
}

