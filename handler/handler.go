package handler

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/liamsorsby/tokeniser/encryption"

	"github.com/urfave/negroni"
)

// Body to hold struct for json body
type Body struct {
	Body string
}

// Server object to hold global services for the server
type Server struct {
	Encryption encryption.Service
}

// New to create a new handler
func (s *Server) New() *negroni.Negroni {
	r := http.NewServeMux()
	r.HandleFunc("/v1/tokenise", s.tokenise)
	r.HandleFunc("/health", s.healthCheck)
	n := negroni.Classic()
	n.UseHandler(r)

	return n
}

func (s *Server) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if _, err := io.WriteString(w, `{"alive": true}`); err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (s *Server) tokenise(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 200))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	result := &Body{}

	if err := json.Unmarshal(body, result); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	if string(result.Body) == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token := s.Encryption.Encrypt([]byte(result.Body), []byte("tokenise"))

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(token); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
