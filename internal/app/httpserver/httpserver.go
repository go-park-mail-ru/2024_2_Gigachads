package httpserver

import (
	"github.com/gorilla/mux"
	"log"
	config "mail/config"
	"net/http"
)

type HTTPServer struct {
	server *http.Server
}

func (s *HTTPServer) Start(config *config.Config) {
	s.server = new(http.Server)
	log.Println(config.HTTPServer.IP + ":" + config.HTTPServer.Port)
	s.server.Addr = config.HTTPServer.IP + ":" + config.HTTPServer.Port
	s.configureRouter()
	log.Println("Server is running on port 8080")
	if err := s.server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
func (s *HTTPServer) configureRouter() {
	router := mux.NewRouter()
	router.HandleFunc("/hello", HelloHandler).Methods("GET")
	router.HandleFunc("/mail/list", getAllMails).Methods("GET")
	s.server.Handler = router
}
