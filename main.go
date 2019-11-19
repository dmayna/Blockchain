package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"uri"
	"uri/handlers"
)

func main() {
	router := uri.NewRouter()
	ifaces, _ := net.Interfaces()
	// handle err
	var ip net.IP
	for _, i := range ifaces {
		addrs, _ := i.Addrs()
		// handle err
		for _, addr := range addrs {

			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
				//fmt.Println(ip)
			case *net.IPAddr:
				ip = v.IP
				//fmt.Println(ip)
			}
			// process IP address
			//fmt.Println(ip)
		}
	}
	var port string
	if len(os.Args) > 1 {
		port = os.Args[1]
	} else {
		port = "8080"
	}

	handlers.Miner.Id = ip.String() + port
	handlers.Miner.Port = port
	log.Println("Listening on port: ", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
