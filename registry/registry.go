package registry

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
)

const (
	defaultPutPath = "/_geerpc_/registry/put"
	defaultGetPath = "/_geerpc_/registry/get"
)

type Registry struct {
	sync.Mutex
	servers map[string][]ServiceItem
}

type ServiceItem struct {
	ServiceName string `json:"service_name"`
	Ip          string `json:"ip"`
	Port        string `json:"port"`
}

func NewRegistry() *Registry {
	return &Registry{servers: make(map[string][]ServiceItem)}
}

var DefaultRegistry = NewRegistry()

func (r *Registry) PutServer(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		io.WriteString(w, "invalid http method\n")
		return
	}

	req.ParseForm()
	server := req.Form.Get("server")
	serverInfo := strings.Split(server, ":")
	log.Println("[PutServer] registry receive register:", server)

	if len(serverInfo) != 3 {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "invalid server param\n")
		return
	}

	serviceItem := &ServiceItem{
		ServiceName: serverInfo[0],
		Ip:          serverInfo[1],
		Port:        serverInfo[2],
	}
	r.Lock()
	if _, ok := r.servers[serviceItem.ServiceName]; !ok {
		r.servers[serviceItem.ServiceName] = []ServiceItem{}
	}
	r.servers[serviceItem.ServiceName] = append(r.servers[serviceItem.ServiceName], *serviceItem)
	r.Unlock()

	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "Success\n")
	return
}

func (r *Registry) GetServers(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		io.WriteString(w, "invalid http method\n")
		return
	}

	req.ParseForm()
	service := req.Form.Get("service")
	log.Println("[PutServer] registry get service:", service)

	r.Lock()
	if _, ok := r.servers[service]; !ok {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "no available address for service: %s\n", service)
		return
	}
	availableServers := r.servers[service]
	r.Unlock()

	fmt.Println(availableServers)

	resp, _ := json.Marshal(availableServers)
	fmt.Println(string(resp))

	w.WriteHeader(http.StatusOK)
	w.Write(resp)

	return
}

func (r *Registry) HandleHttp() {
	http.HandleFunc(defaultPutPath, r.PutServer)
	http.HandleFunc(defaultGetPath, r.GetServers)
}

func HandleHttp() {
	DefaultRegistry.HandleHttp()
	http.ListenAndServe(":8080", nil)
}
