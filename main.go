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

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func main () {
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

for k, v := range m {
    switch v := v.(type) {
    case string:
        fmt.Println(k, "is string", v)
    case int:
        fmt.Println(k, "is int", v)
    case []interface{}:
        fmt.Println(k, "is an array:")
        for i, u := range v {
            fmt.Println(i, u)
        }
    default:
        fmt.Println(k, "is of a type I don't know how to handle")
    }
} // for

fmt.Printf("http://%s/v1/catalog/nodes\n", m["hostname"])

} // main
