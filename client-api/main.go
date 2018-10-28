package main

import (
	"github.com/hashicorp/consul/api"
	"log"
	"net/http"
	"strconv"
	"io/ioutil"
)

func main() {
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		log.Println(err)
	}

	http.HandleFunc("/service/", func(w http.ResponseWriter, r *http.Request) {
		svcName := r.URL.Path[len("/service/"):]
		svcs, _, _ := client.Catalog().Service(svcName, "", nil)

		if len(svcs) < 1 {
			w.Write([]byte("Service not found"))
			return
		}
		rsp, err := http.Get("http://" + svcs[0].ServiceAddress + ":" + strconv.Itoa(svcs[0].ServicePort))
		if err != nil {
			log.Println(err.Error())
		}
		rspData, _ := ioutil.ReadAll(rsp.Body)
		w.Write(rspData)
	})

	kv, _, _ := client.KV().Get("ports/client-api", nil)

	log.Println("Client API waiting for requests...")
	http.ListenAndServe(":" + string(kv.Value), nil)
}
