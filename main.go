package main

import (
	"flag"
	"fmt"

	"teissem.fr/bind9_rest_api/webserver"
)

func main() {
	// Command line params
	var port int
	var namedConfLocal string
	var commandReloadBind string
	flag.IntVar(&port, "p", 8000, "Provide a port number")
	flag.StringVar(&namedConfLocal, "f", "/etc/bind/named.conf.local", "Provide the location of named.conf.local")
	flag.StringVar(&commandReloadBind, "c", "systemctl reload bind", "Provide the command to reload bind9")
	flag.Parse()
	// Start the web server
	ws := webserver.InitWebServer(namedConfLocal, commandReloadBind)
	ws.GinEngine.Run(fmt.Sprintf(":%d", port))
}
