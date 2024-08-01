package spanner_demo_20240801

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
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
	var body SearchReq
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Println(body.Text)

	w.Header().Set("Content-Type", "text/javascript;charset=utf-8")
	w.WriteHeader(http.StatusOK)
	resp := SearchResp{
		[]*SampleMessage{
			{
				SampleMessageID: "sampleID-1",
				Message:         "hello",
				CreatedAt:       time.Now(),
			},
			{
				SampleMessageID: "sampleID-2",
				Message:         fmt.Sprintf("echo %s", body.Text),
				CreatedAt:       time.Now(),
			},
		},
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		fmt.Println(err)
	}
}
