package main

import (
	"github.com/hashicorp/consul/api"
	"log"
	"net/http"
	"strconv"
	"fmt"
)

func main() {
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		log.Println(err)
	}

	kv, _, _ := client.KV().Get("ports/user-service", nil)
	port, _ := strconv.Atoi(string(kv.Value))

	serviceDef := &api.AgentServiceRegistration{
		ID:"user-service",
		Name:"user-service",
		Address:"127.0.0.1",
		Port:port,
		Tags:[]string{"product"},
		Check:&api.AgentServiceCheck{
			HTTP:fmt.Sprintf("http://127.0.0.1:%s/health", string(kv.Value)),
			Interval:"5s",
			Timeout:"3s",
		},
	}

	client.Agent().ServiceRegister(serviceDef)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("I'm User Service"))
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	log.Println("User Service srated...")
	http.ListenAndServe(":" + string(kv.Value), nil)
}
