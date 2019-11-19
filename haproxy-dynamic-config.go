package main

import (
	"flag"
	"fmt"

	"github.com/hashicorp/consul/api"
)

var (
	flagOp string // Pode ser generate-config ou update-backends
)

func main() {

	flag.Parse()

	//	Operações:
	//	* generate-config
	//	* update-config

	// Get a new client
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		panic(err)
	}

	catalog := client.Catalog()

	svc, _, _ := catalog.Service("cluster_rest_api", "", nil)

	//spew.Dump(svc)

	for _, s := range svc {
		fmt.Println(s.Node)
	}

	// O arquivo haproxy.cfg tem backends definidos da seguinte forma:
	// backend cluster_rest_api
	//		option balance leastconn
	//		timeout server 60s
	//		option forceclose
	//		http-request set-header client-ip %[src]
	//		http-response add-header X-Backend %b
	//		http-response add-header X-Server %s
	//		balance leastconn
	//		#{BACKEND_SERVERS}
}

func init() {
	flag.StringVar(&flagOp, "op", "", "Informe uma operação válida")
}
