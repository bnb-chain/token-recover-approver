package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/bnb-chain/airdrop-service/internal/module/approval"
	"github.com/julienschmidt/httprouter"
)

func (server *HttpServer) Ping(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	server.Response(w, Success, "pong", nil)
}

func (server *HttpServer) GetClaimApproval(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		server.Response(w, InvalidRequest, nil, err)
		return
	}
	defer r.Body.Close()
	req := &approval.GetClaimApprovalRequest{}
	err = json.Unmarshal(body, req)
	if err != nil {
		server.Response(w, InvalidRequest, nil, err)
		return
	}

	resp, err := server.approvalService.GetClaimApproval(req)
	if err != nil {
		server.Response(w, InvalidRequest, nil, err)
		return
	}

	server.Response(w, Success, resp, nil)
}

func (server *HttpServer) Response(w http.ResponseWriter, code ResponseCode, data interface{}, err error) {
	resp := Response{
		Code:  code,
		Data:  data,
		Error: err,
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, resp.Marshal())
}
