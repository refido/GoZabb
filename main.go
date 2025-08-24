package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"go-crud-zabbix/zabbix"
	"github.com/joho/godotenv"
)

var zbx *zabbix.Client

func mustGetenv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("environment variable %s is not set", key)
	}
	return v
}

func main() {
	// Load .env if present (non-fatal if missing; useful in prod where real env is provided)
	_ = godotenv.Load()

	api := mustGetenv("ZABBIX_API")
	user := mustGetenv("ZABBIX_USER")
	pass := mustGetenv("ZABBIX_PASSWORD")

	zbx = zabbix.New(api)
	if err := zbx.Login(user, pass); err != nil {
		log.Fatal("Login failed:", err)
	}

	http.HandleFunc("/version", getVersion)
	http.HandleFunc("/hosts", getHosts)
	http.HandleFunc("/host/create", createHost)
	http.HandleFunc("/host/delete", deleteHost)
	http.HandleFunc("/test", getTest)

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getHosts(w http.ResponseWriter, r *http.Request) {
	result, err := zbx.Call("host.get", map[string]interface{}{
		"output":           []string{"hostid", "host", "name"},
		"selectInterfaces": []string{"interfaceid", "ip"},
		"selectGroups":     []string{"groupid", "name"},
		"sortfield":        "host",
		"sortorder":        "ASC",
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Zabbix response:", string(result))

	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
}

func createHost(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Host string `json:"host"`
		IP   string `json:"ip"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	params := map[string]interface{}{
		"host": payload.Host,
		"interfaces": []map[string]interface{}{
			{"type": 1, "main": 1, "useip": 1, "ip": payload.IP, "dns": "", "port": "10050"},
		},
		"groups": []map[string]string{
			{"groupid": "2"},
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

func deleteHost(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		HostIDs []string `json:"hostids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	result, err := zbx.Call("host.delete", payload.HostIDs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
}

func getTest(w http.ResponseWriter, r *http.Request) {
	result, err := zbx.Call("host.get", map[string]interface{}{
		"output": "extend",
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 👇 raw response from Zabbix
	log.Println("Zabbix raw response:", string(result))

	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
}

func getVersion(w http.ResponseWriter, r *http.Request) {
	version, err := zbx.Version()
	if err != nil {
		http.Error(w, "Failed to get version: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"version": version})
}
