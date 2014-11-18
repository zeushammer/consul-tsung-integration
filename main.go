package main

/*
Take in parameter where file is -f config.json
Inject a config file that has N paramters
If kddey = host
    then
If key = ConsulClient
    then
If key = CPU
If key = MaxUsers
If key = ipscan
If key = portMax

inner join servers from client
curl http://172.20.4.60:8500/v1/catalog/nodes
[{"Node":"loadnode01","Address":"172.20.4.50"},{"Node":"loadnode02","Address":"172.20.4.53"},{"Node":"loadnode03","Address":"172.20.4.68"},{"Node":"loadnode04","Address":"172.20.4.71"},{"Node":"tsungmaster","Address":"172.20.4.60"}]

and servers in key/value store

to remove any that are potentially "down"

create array of struct(s) that hold the infomation needed to execute a test
*/

import "encoding/json"
import "fmt"
import "flag"
import "io/ioutil"
import b64 "encoding/base64"
import "net/http"
import "reflect"

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {

	// load config file
	// get nodes from catalog
	// get nodes from kv
	// inner join on node
	// return xml structs built

	configFilePtr := flag.String("f", "config.json", "Fully qualified path to config file [config.json]")

	flag.Parse()

	fmt.Println("config file:", *configFilePtr)

	configBytes, err := ioutil.ReadFile(*configFilePtr)
	check(err)
	fmt.Print(string(configBytes))

	var configFile interface{}
	err = json.Unmarshal(configBytes, &configFile)
	check(err)

	m := configFile.(map[string]interface{})

	fmt.Println(reflect.TypeOf(m).String())

	getAvailableNodesFromCatalog(m)
	//getRegisteredNodesFromKeyValueStore()

} // main

func getAvailableNodesFromCatalog(m map[string]interface{}) {

	fmt.Printf("http://%s/v1/catalog/nodes\n", m["consulclientip"])
	url := fmt.Sprintf("http://%s/v1/catalog/nodes", m["consulclientip"])
	response, err := http.Get(url)
	check(err)

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	check(err)

	fmt.Printf("%s\n", string(contents))

	var nodeCatalog []interface{}
	err = json.Unmarshal(contents, &nodeCatalog)
	check(err)

	fmt.Println(reflect.TypeOf(nodeCatalog).String())

	json := nodeCatalog[0].(map[string]interface{})
	_ = json
	//fmt.Printf("http://%s/v1/kv/nodes?recurse\n", m["consulclientip"])

}

func getRegisteredNodesFromKeyValueStore() {
	/*
	   [
	     {
	       "CreateIndex": 1270,
	       "ModifyIndex": 1270,
	       "LockIndex": 0,
	       "Key": "nodes\/loadnode01",
	       "Flags": 0,
	       "Value": "eyJob3N0bmFtZSI6ICJsb2Fkbm9kZTAxIiwiY3B1IjogIjQiLCJtYXh1c2VycyI6ICI2MDAwMCIsImlwc2NhbiI6ICJldGgwIn0="
	     },
	     {
	       "CreateIndex": 1272,
	       "ModifyIndex": 1272,
	       "LockIndex": 0,
	       "Key": "nodes\/loadnode02",
	       "Flags": 0,
	       "Value": "eyJob3N0bmFtZSI6ICJsb2Fkbm9kZTAyIiwiY3B1IjogIjQiLCJtYXh1c2VycyI6ICIzMDAwMCIsImlwc2NhbiI6ICJldGgwIn0="
	     },
	     {
	       "CreateIndex": 1278,
	       "ModifyIndex": 1278,
	       "LockIndex": 0,
	       "Key": "nodes\/loadnode03",
	       "Flags": 0,
	       "Value": "eyJob3N0bmFtZSI6ICJsb2Fkbm9kZTAzIiwiY3B1IjogIjQiLCJtYXh1c2VycyI6ICIzMDAwMCIsImlwc2NhbiI6ICJldGgwIn0="
	     },
	     {
	       "CreateIndex": 1269,
	       "ModifyIndex": 1318,
	       "LockIndex": 0,
	       "Key": "nodes\/",
	       "Flags": 0,
	       "Value": "eyJob3N0bmFtZSI6ICIkaG9zdCIsImNwdSI6ICIkY3B1IiwibWF4dXNlcnMiOiAiMzAwMDAiLCJpcHNjYW4iOiAiZXRoMCJ9"
	     }
	   ]
	*/

	// for each dict in array
	// iterate over k/v pairs
	// if trimmed key from kv store is in nodes from catalog
	// make new slice of dicts
	// add meta data to new dict
	// use for later translation to structs for delivery to xml builder
	sEnc := "=="
	sDec, _ := b64.StdEncoding.DecodeString(sEnc)
	fmt.Println(string(sDec))
	fmt.Println()
}
