package spanner_demo_20240801

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type MessageHandler struct {
	s *Service
}

type PostMessageReq struct {
	Text string `json:"text"`
}

type PostMessageResp struct {
	Result *SampleMessage `json:"result"`
}

func (h *MessageHandler) PostMessageHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var body PostMessageReq
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	msg, err := h.s.Insert(ctx, &SampleMessage{
		Message: body.Text,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/javascript;charset=utf-8")
	w.WriteHeader(http.StatusOK)
	resp := PostMessageResp{
		Result: msg,
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		fmt.Println(err)
	}
}

type SearchReq struct {
	Text string `json:"text"`
}

type SearchResp struct {
	Results []*SampleMessage `json:"results"`
}

func (h *MessageHandler) SearchHandler(w http.ResponseWriter, r *http.Request) {
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
