package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"go-crud-zabbix/zabbix"
)

var zbx *zabbix.Client

func main() {
	var err error
	zbx, err = zabbix.Login("tro_admin", "gacor!") // change to your credentials
	if err != nil {
		log.Fatal("Login failed:", err)
	}

	http.HandleFunc("/hosts", getHosts)         // GET
	http.HandleFunc("/host/create", createHost) // POST
	http.HandleFunc("/host/delete", deleteHost) // DELETE

	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", nil)
}

// GET hosts
func getHosts(w http.ResponseWriter, r *http.Request) {
	result, err := zbx.Call("host.get", map[string]interface{}{
		"output": []string{"hostid", "host"},
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
}

// CREATE host
func createHost(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Host string `json:"host"`
		IP   string `json:"ip"`
	}
	json.NewDecoder(r.Body).Decode(&payload)

	params := map[string]interface{}{
		"host": payload.Host,
		"interfaces": []map[string]interface{}{
			{
				"type": 1, "main": 1, "useip": 1,
				"ip": payload.IP, "dns": "", "port": "10050",
			},
		},
		"groups": []map[string]string{
			{"groupid": "2"}, // Linux servers (adjust for your env)
		},
	}

	result, err := zbx.Call("host.create", params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
}

// DELETE host
func deleteHost(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		HostIDs []string `json:"hostids"`
	}
	json.NewDecoder(r.Body).Decode(&payload)

	result, err := zbx.Call("host.delete", payload.HostIDs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
}
