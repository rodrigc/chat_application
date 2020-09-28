package chat_application

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

var router *mux.Router

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

func listEndpoints(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	htmlOutput := `
<!DOCTYPE html>
<html>
<head>
  <title>Chat server REST API endpoints</title>
</head>
<body>
  <h1>Chat server REST API endpoints</h1>
  <ul>
`
	err := router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, err := route.GetPathTemplate()
		if err == nil {
			htmlOutput += fmt.Sprintf("<li><a href=\"%s\">%s</a></li>", pathTemplate, pathTemplate)
		}
		return err
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	htmlOutput += `
</ul>
</body>
</html>
`
	io.WriteString(w, htmlOutput)

}

func RestAPIServer(wg *sync.WaitGroup) {
	defer wg.Done()

	router = mux.NewRouter()

	router.HandleFunc("/", listEndpoints)
	router.HandleFunc("/api/user", listUsers)
	router.HandleFunc("/api/room", listRooms)

	srv := &http.Server{
		Handler: router,
		Addr:    fmt.Sprintf("%s:%d", chatConfig.Host, chatConfig.RestAPIPort),
	}

	fmt.Printf("Starting REST API endpoint on port %d\n", chatConfig.RestAPIPort)
	log.Fatal(srv.ListenAndServe())
}
