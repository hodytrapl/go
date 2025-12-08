package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func work(id int) {
	for i := 0; i < 3; i++ {
		fmt.Println("worker ", id, " iteration", i)
		time.Sleep(200 * time.Millisecond)
	}
}

func main() {
	/*http.HandleFunc("/", infoHandler)
	PORT := ":8080"
	err:=http.ListenAndServe(PORT, nil)
	fmt.Printf("сервак запущен на порту%s", PORT)
	if(err!=nil){
		log.Fatal()
	}*/

	/*go work(1)
	go work(2)
	time.Sleep(1 * time.Second)
	fmt.Println("done")*/

	/*for i:=0;i<1000;i++{
		func(id int){
			go fmt.Println("hello:" ,id)
		}(i)
	}
	time.Sleep(500 * time.Millisecond)*/
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("приветик"))
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	response := struct {
		Status      string `json:"status`
		ProjectName string `json:project_name`
	}{
		Status:      "Available",
		ProjectName: "GOLANG",
	}

	result, err := json.Marshal(response)

	if err != nil {
		http.Error(w, "server error", 500)
	}
	w.Header().Set("Content-Type", "json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func infoHandler(w http.ResponseWriter, r *http.Request) {
	response := struct {
		Name    string `json:"name"`
		Version string `json:"version"`
		Time    string `json:"server_time"`
	}{
		Name:    "g",
		Version: "1.2.4",
		Time:    time.Now().Format(time.RFC3339),
	}

	result, err := json.MarshalIndent(response, "", "   ")

	if err != nil {
		http.Error(w, "server error", 500)
	}
	w.Header().Set("Content-Type", "json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}
