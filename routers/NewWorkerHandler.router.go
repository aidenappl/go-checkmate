package routers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/aidenappl/go-checkmate/poller"
	"github.com/aidenappl/go-checkmate/responder"
)

type NewWorkerHandlerRequest struct {
	Endpoint string `json:"endpoint"`
	Name     string `json:"name"`
	Service  string `json:"service"`
	Interval int    `json:"interval"`
	Timeout  int    `json:"timeout"`
}

type WorkerAPI struct {
	Poller *poller.Manager
}

func (h *WorkerAPI) NewWorkerHandler(w http.ResponseWriter, r *http.Request) {
	var req NewWorkerHandlerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responder.ParsingError(w, err)
		return
	}

	// Check if the endpoint is provided
	if req.Endpoint == "" {
		responder.ParamError(w, "endpoint")
		return
	}
	// Check if the name is provided
	if req.Name == "" {
		responder.ParamError(w, "name")
		return
	}
	// Check if the service is provided
	if req.Service == "" {
		responder.ParamError(w, "service")
		return
	}
	// Check if the interval is provided
	if req.Interval == 0 {
		responder.ParamError(w, "interval")
		return
	}
	// Check if the timeout is provided
	if req.Timeout == 0 {
		responder.ParamError(w, "timeout")
		return
	}

	id := req.Service + ":" + req.Name

	err := h.Poller.AddOrUpdate(poller.Spec{
		ID:       id,
		Name:     req.Name,
		Service:  req.Service,
		Endpoint: req.Endpoint,
		Interval: time.Duration(req.Interval) * time.Second,
		Timeout:  5 * time.Second, // override per spec if you like
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

}
