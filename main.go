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
import "bytes"
import "encoding/xml"
import "os"
import "strings"

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

	XMLoptions struct {
		XMLName xml.Name    `xml:"options"`
		Options []XMLoption `xml:"option"`
	}

	XMLoption struct {
		XMLName xml.Name `xml:"option"`
		Name    string   `xml:"name,attr"`
		Min     int      `xml:"min,attr"`
		Max     int      `xml:"max,attr"`
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

	XMLload struct {
		XMLName       xml.Name          `xml:"load"`
		Duration      int               `xml:"duration,attr"`
		Unit          string            `xml:"unit,attr"`
		ArrivalPhases []XMLarrivalphase `xml:"arrivalphase"`
	}

	XMLarrivalphase struct {
		XMLName  xml.Name `xml:"arrivalphase"`
		Phase    int      `xml:"phase,attr"`
		Duration int      `xml:"duration,attr"`
		Unit     string   `xml:"unit,attr"`
		Users    XMLusers `xml:"users"`
	}

	XMLusers struct {
		XMLName     xml.Name `xml:"users"`
		Maxnumber   int      `xml:"maxnumber,attr"`
		Arrivalrate int      `xml:"arrivalrate,attr"`
		Unit        string   `xml:"unit,attr"`
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
		XMLip{Scan: "true", Value: "eth0"}}}
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
		// do not use the hostname as a loadserver
		if node.Node != config.Hostname {
			fmt.Println(node.Node)
			fmt.Println(node.Address)
			servers.Clients = append(servers.Clients, *NewXMLclient(node.Node))
		}
	}

	fmt.Printf("%+v\n\n", servers)

	load := XMLload{Duration: 1, Unit: "minute",
		ArrivalPhases: []XMLarrivalphase{
			XMLarrivalphase{
				Phase:    1,
				Duration: 1,
				Unit:     "minute",
				Users: XMLusers{
					Maxnumber:   60000,
					Arrivalrate: 100,
					Unit:        "minute"}}}}

	options := XMLoptions{Options: []XMLoption{
		XMLoption{Name: "ports_range", Min: 1024, Max: 65535}}}

	fd, err := os.Open("/home/andrewgerhold/json-and-go/boris_go_loadtest.xml")
	check(err)

	bufferForFile := new(bytes.Buffer)
	bufferForFile.ReadFrom(fd)
	fileAsString := bufferForFile.String()

	modifyXmlValues(fileAsString, load, servers, options)
}

func testOutput() {
	cloud := Clients{}
	cloud = Clients{
		Clients: []XMLclient{
			XMLclient{Host: "cmbload01", Weight: 1, CPU: 4, Maxusers: 40000},
			XMLclient{Host: "cmbload02", Weight: 1, CPU: 4, Maxusers: 40000},
			XMLclient{Host: "cmbload03", Weight: 1, CPU: 4, Maxusers: 40000}}}
	fmt.Printf("%+v\n", cloud)
}

func modifyXmlValues(content string, load XMLload, cloud Clients, options XMLoptions) {
	var buffer bytes.Buffer
	inputReader := strings.NewReader(content)
	decoder := xml.NewDecoder(inputReader)
	encoder := xml.NewEncoder(&buffer)
	encoder.Indent("", " ")
	buffer.WriteString(xml.Header)
	buffer.WriteString("<!DOCTYPE tsung SYSTEM '/usr/share/tsung/tsung-1.0.dtd'>\n")
	for {
		t, _ := decoder.Token()
		if t == nil {
			break
		}
		switch token := t.(type) {
        //Damn this feels soooooo hacky
		case xml.StartElement:
			fmt.Println(t)
			switch token.Name.Local {
			case "load":
				encoder.Encode(load)
			case "clients":
				encoder.Encode(cloud)
			case "options":
				encoder.Encode(options)
			case "client", "arrivalphase", "users", "option":
				// allows me to ignore the inner element of this
				// I probably should explain why this is necessary
				// I can't remember myself as of right now
			default:
				err := encoder.EncodeToken(t)
				check(err)
			}
		case xml.EndElement:
			switch token.Name.Local {
			case "clients", "load", "client", "arrivalphase", "users", "options", "option":
				// allows me to ignore end tag errors like
				// xml: end tag </client> does not match start tag <tsung>
			default:
				err := encoder.EncodeToken(t)
				check(err)
			}
		}
	}
	encoder.Flush()
	fmt.Println(buffer.String())
}
