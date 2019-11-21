package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/hashicorp/consul/api"
)

// Op main object to make operation calls
type Op struct{}

var cfgFiles = []string{"haproxy-template.cfg", "haproxy-final.cfg"}

// Clear removes previously created config files
func (o *Op) Clear() {
	for _, file := range cfgFiles {
		filePath := "./resources/" + file
		e := os.Remove(filePath)
		if e != nil {
			log.Print("WARN " + e.Error())
		}
	}

}

// Show previously created config files
func (o *Op) Show() {
	for _, file := range cfgFiles {
		filePath := "./resources/" + file
		b, e := ioutil.ReadFile(filePath)
		if e != nil {
			log.Print("WARN File " + filePath + " not found.")
		} else {
			log.Print("INFO Showing contents of " + filePath)
			out := color.CyanString(string(b))
			fmt.Println(out)
		}
	}

}

// GenerateConfig creates an basic HAProxy configuration
func (o *Op) GenerateConfig() {
	b, e := ioutil.ReadFile("./samples/haproxy-sample.cfg")
	if e != nil {
		log.Print("WARN Failed to read HAProxy sample file.")
		os.Exit(0)
	}
	e = ioutil.WriteFile("./resources/haproxy-template.cfg", b, 0644)
	if e != nil {
		log.Print("WARN Failed to write HAProxy template file.")
		os.Exit(0)
	}
	log.Print("INFO HAProxy template created sucessfully!")
}

// UpdateBackends queries consul cluster to get available backend servers
func (o *Op) UpdateBackends() {

	// Reads current config
	b, e := ioutil.ReadFile("./resources/haproxy-template.cfg")
	if e != nil {
		log.Print("WARN Failed to load pre-defined HAProxy template! Use generate-config first.")
		os.Exit(0)
	}

	cfg := string(b)

	// Get a new client
	clientConfig := &api.Config{
		Address:    "localhost:8500",
		Datacenter: "dc-local",
		Scheme:     "http",
	}
	client, err := api.NewClient(clientConfig)
	if err != nil {
		panic(err)
	}

	//spew.Dump(client)

	catalog := client.Catalog()

	svc, _, _ := catalog.Service("cluster_rest_api", "", nil)

	serverList := make([]string, 0)
	for _, s := range svc {
		serverList = append(serverList, fmt.Sprintf("server %s %s:%d check inter 5s", s.Node, s.Address, s.ServicePort))
	}

	cfg = strings.Replace(cfg, "#CLUSTER_REST_API_BACKEND_SERVERS#", strings.Join(serverList, "\n    "), 1)

	svc, _, _ = catalog.Service("cluster_soap_api", "", nil)

	serverList = serverList[:0]
	for _, s := range svc {
		serverList = append(serverList, fmt.Sprintf("server %s %s:%d check inter 5s", s.Node, s.Address, s.ServicePort))
	}

	cfg = strings.Replace(cfg, "#CLUSTER_SOAP_API_BACKEND_SERVERS#", strings.Join(serverList, "\n    "), 1)

	if flagCommit {
		e := ioutil.WriteFile("./resources/haproxy-final.cfg", []byte(cfg), 0644)
		if e != nil {
			log.Print("WARN Failed to write final HAProxy configuration: " + e.Error())
		}
		log.Print("INFO Sucessfully written HAProxy final configuration!")
	} else {
		fmt.Println(cfg)
	}
}
