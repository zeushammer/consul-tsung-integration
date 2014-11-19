package main

/*
Take in parameter where file is -f config.json
Inject a config file that has N paramters

inner join servers from client
curl http://172.20.4.60:8500/v1/catalog/nodes

assume all settings are:
{
    "hostname": "loadnode01",
    "cpu": "4",
    "maxusers": "60000",
    "ipscan": "eth0",
    "largeportrange": true
}

config file lookts like:
{
    "hostname": "tsungmaster",
    "consulclientip": "172.20.4.60:8500",
}

create array of struct(s) that hold the infomation needed to execute a test
*/

import "encoding/json"
import "fmt"
import "flag"
import "io/ioutil"
import "net/http"

//import "reflect"
import "encoding/xml"
import "os"
import "bytes"

type (
	Config struct {
		Hostname       string
		ConsulClientIP string
	}

	Clients struct {
		XMLName xml.Name    `xml:"clients"`
		Clients []XMLclient `xml:"client"`
	}

	XMLclient struct {
		XMLName  xml.Name `xml:"client"`
		Host     string   `xml:"host,attr"`
		Weight   int      `xml:"weight,attr"`
		CPU      int      `xml:"cpu,attr"`
		Maxusers int      `xml:"maxusers,attr"`
		Ip       []XMLip  `xml:"ip"`
	}

	XMLip struct {
		XMLName xml.Name `xml:"ip"`
		Scan    string   `xml:"scan,attr"`
		Value   string   `xml:"value,attr"`
	}

	AvailableNodes []struct {
		Node    string `json:"Node"`
		Address string `json:"Address"`
	}
)

var (
	config  Config
	servers Clients
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func NewXMLclient(host string) *XMLclient {
	return &XMLclient{Host: host, Weight: 1, CPU: 4, Maxusers: 60000, Ip: []XMLip{
		XMLip{Scan: "yes", Value: "eth0"}}}
} //NewXMLclient

func main() {

	// load config file
	// get nodes from catalog
	// return arrary of xml structs built in a new Client struct

	configFilePtr := flag.String("f", "config.json", "Fully qualified path to config file [config.json]")
	flag.Parse()
	fmt.Println("config file:", *configFilePtr)

	fd, err := os.Open(*configFilePtr)
	check(err)
	err = json.NewDecoder(fd).Decode(&config)
	check(err)

	getAvailableNodesFromCatalog()

} // main

func getAvailableNodesFromCatalog() {

	fmt.Printf("http://%s/v1/catalog/nodes\n", config.ConsulClientIP)
	url := fmt.Sprintf("http://%s/v1/catalog/nodes", config.ConsulClientIP)
	response, err := http.Get(url)
	check(err)

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	check(err)

	// fmt.Printf("%s\n", string(contents))

	var catalog AvailableNodes
	err = json.NewDecoder(bytes.NewReader(contents)).Decode(&catalog)
	check(err)

	for _, node := range catalog {
		fmt.Println(node.Node)
		fmt.Println(node.Address)
		servers.Clients = append(servers.Clients, *NewXMLclient(node.Node))
	}

	fmt.Printf("%+v\n", servers)
}
