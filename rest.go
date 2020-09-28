package chat_application

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

func listUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}

func listRooms(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	roomJsonMap := map[string]map[string]string{}
	roomJsonMap["room"] = map[string]string{}

	roomMtx.Lock()
	for room, _ := range roomMap {
		roomJsonMap["room"][room] = ""
	}
	roomMtx.Unlock()

	json.NewEncoder(w).Encode(roomJsonMap)
}

func RestAPIServer(wg sync.WaitGroup) {
	defer wg.Done()

	router := mux.NewRouter()

	router.HandleFunc("/api/user", listUsers)
	router.HandleFunc("/api/room", listRooms)

	srv := &http.Server{
		Handler: router,
		Addr:    fmt.Sprintf("%s:%d", chatConfig.Host, chatConfig.RestAPIPort),
	}

	fmt.Printf("Starting router on port %d\n", chatConfig.RestAPIPort)
	log.Fatal(srv.ListenAndServe())
}
