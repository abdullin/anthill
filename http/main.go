package main


import (
    "log"
    "net/http"
	  "encoding/json"

	"github.com/gorilla/mux"
)
type VersionResponse struct {
	Name string `json:"name"`
	Env string `json:"env"`
	ID string `json:"id"`
	BuildID string `json:"buildId"`
	BuildStamp string `json:"buildStamp"`
	UptimeHours int `json:"uptimeHours"`
}


func main() {

	router := mux.NewRouter().StrictSlash(true)
    router.HandleFunc("/status/version", Index)
	

    log.Fatal(http.ListenAndServe(":8080", router))

}

func Index(w http.ResponseWriter, r *http.Request) {

	resp := VersionResponse {
		Name: "APIv2",
		Env: "production",
		ID: "SkuVault.ApiRole_IN_0",
		BuildID: "v2production_s84d_r86b",
		BuildStamp: "20161129-1519",
		UptimeHours: 0,
		
		
	}
	json.NewEncoder(w).Encode(resp)
}
