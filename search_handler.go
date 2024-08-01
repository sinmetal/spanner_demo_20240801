package spanner_demo_20240801

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type SearchReq struct {
	Text string `json:"text"`
}

type SearchResp struct {
	Results []*SampleMessage `json:"results"`
}

type SearchHandler struct {
	s *Service
}

func (h *SearchHandler) Handler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var body SearchReq
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	msgs, err := h.s.SearchMessage(ctx, body.Text)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/javascript;charset=utf-8")
	w.WriteHeader(http.StatusOK)
	resp := SearchResp{
		Results: msgs,
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		fmt.Println(err)
	}
}
