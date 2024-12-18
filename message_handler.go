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
	Title   string   `json:"title"`
	Message string   `json:"message"`
	Tags    []string `json:"tags"`
}

type PostMessageResp struct {
	Result *SampleMessage `json:"result"`
}

func (h *MessageHandler) PostMessageHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var body PostMessageReq
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Printf("failed request body json decode. %s\n", err)
		return
	}

	msg, err := h.s.Insert(ctx, &SampleMessage{
		Title:   body.Title,
		Message: body.Message,
		Tags:    body.Tags,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("failed SampleMessage.Insert. %s\n", err)
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
	Results []*SearchMessageResult `json:"results"`
}

func (h *MessageHandler) SearchHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var body SearchReq
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Printf("failed request body json decode. %s\n", err)
		return
	}

	msgs, err := h.s.SearchMessage(ctx, body.Text)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("failed SampleMessage.SearchMessage. %s\n", err)
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

type SearchCompositeReq struct {
	Tag     string `json:"tag"`
	Title   string `json:"title"`
	Message string `json:"message"`
}

type SearchCompositeResp struct {
	Results []*SearchSampleMessagesResp `json:"results"`
}

func (h *MessageHandler) SearchCompositeHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var body SearchCompositeReq
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Printf("failed request body json decode. %s\n", err)
		return
	}

	msgs, err := h.s.SearchSampleMessages(ctx, &SearchSampleMessagesReq{
		Tag:     body.Tag,
		Title:   body.Title,
		Message: body.Message,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("failed SampleMessage.SearchMessage. %s\n", err)
		return
	}

	w.Header().Set("Content-Type", "text/javascript;charset=utf-8")
	w.WriteHeader(http.StatusOK)
	resp := SearchCompositeResp{
		Results: msgs,
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		fmt.Println(err)
	}
}
