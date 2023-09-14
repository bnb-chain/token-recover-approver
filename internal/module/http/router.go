package http

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (server *HttpServer) Ping(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Pong!")
}

func (server *HttpServer) GetClaimApproval(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "GetClaimApproval!")
}

func (server *HttpServer) GetRegisterTokenApproval(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "GetRegisterTokenApproval!")
}
