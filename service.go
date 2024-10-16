package spanner_demo_20240801

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"text/template"
	"time"

	"cloud.google.com/go/spanner"
	"google.golang.org/api/iterator"
)

const SampleMessagesTable = "SampleMessages"

type SampleMessage struct {
	SampleMessageID string    `json:"sampleMessageID"`
	Tags            []string  `json:"tags"`
	Title           string    `json:"title"`
	Message         string    `json:"message"`
	CreatedAt       time.Time `json:"createdAt"`
}

type Service struct {
	spa *spanner.Client
}

func NewService(ctx context.Context, spa *spanner.Client) (*Service, error) {
	return &Service{
		spa: spa,
	}, nil
}

func (s *Service) CreateSampleMessageID(message string) string {
	v := sha1.Sum([]byte(message))
	return hex.EncodeToString(v[:])
}

func (s *Service) Insert(ctx context.Context, value *SampleMessage) (*SampleMessage, error) {
	value.SampleMessageID = s.CreateSampleMessageID(value.Message)
	value.CreatedAt = spanner.CommitTimestamp
	commitTimestamp, err := s.spa.ReadWriteTransaction(ctx, func(ctx context.Context, tx *spanner.ReadWriteTransaction) error {
		m, err := spanner.InsertOrUpdateStruct(SampleMessagesTable, value)
		if err != nil {
			return err
		}
		if err := tx.BufferWrite([]*spanner.Mutation{m}); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	value.CreatedAt = commitTimestamp
	return value, nil
}

type SearchMessageResult struct {
	SampleMessageID string    `json:"sampleMessageID"`
	Message         string    `json:"message"`
	Score           float64   `json:"score"`
	CreatedAt       time.Time `json:"createdAt"`
}

func (s *Service) SearchMessage(ctx context.Context, text string) ([]*SearchMessageResult, error) {
	const sql = `SELECT SampleMessageID, Message, SCORE(SampleMessages_Message_Tokens, @text) AS score, CreatedAt
 FROM SampleMessages
 WHERE SEARCH(SampleMessages_Message_Tokens, @text, enhance_query=>true, language_tag=>'ja')
 ORDER BY SCORE(SampleMessages_Message_Tokens, @text)
 LIMIT 50
`
	sts := spanner.NewStatement(sql)
	sts.Params = map[string]interface{}{
		"text": text,
	}

	var rets []*SearchMessageResult
	iter := s.spa.Single().Query(ctx, sts)
	for {
		row, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return nil, err
		}
		var v SearchMessageResult
		if err := row.ToStruct(&v); err != nil {
			return nil, err
		}
		rets = append(rets, &v)
	}
	return rets, nil
}

type SearchSampleMessagesReq struct {
	Tag     string `json:"tag"`
	Title   string `json:"title"`
	Message string `json:"message"`
}

type SearchSampleMessagesResp struct {
	SampleMessageID string    `json:"sampleMessageID"`
	Tags            []string  `json:"tags"`
	Title           string    `json:"title"`
	Message         string    `json:"message"`
	CreatedAt       time.Time `json:"createdAt"`
}

func (s *Service) SearchSampleMessages(ctx context.Context, req *SearchSampleMessagesReq) ([]*SearchSampleMessagesResp, error) {
	const sql = `SELECT SampleMessageID, Tags, Title, Message, CreatedAt
 FROM SampleMessages
 WHERE {{.WHERE}}
 LIMIT 50
`
	params := make(map[string]interface{})
	var whereClause []string
	if req.Tag != "" {
		whereClause = append(whereClause, "SEARCH(SampleMessages_Tags_Tokens, @tag)")
		params["tag"] = req.Tag
	}
	if req.Title != "" {
		whereClause = append(whereClause, "SEARCH(SampleMessages_Title_Tokens, @title)")
		params["title"] = req.Title
	}
	if req.Message != "" {
		whereClause = append(whereClause, "SEARCH(SampleMessages_Message_Tokens, @message)")
		params["message"] = req.Message
	}
	tmpl := template.New("mytemplate")
	tmpl, err := tmpl.Parse(sql)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	data := map[string]string{"WHERE": strings.Join(whereClause, "AND")}
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return nil, err
	}

	finishSQL := buf.String()
	fmt.Println(finishSQL)
	sts := spanner.NewStatement(finishSQL)
	sts.Params = params
	var rets []*SearchSampleMessagesResp
	iter := s.spa.Single().Query(ctx, sts)
	for {
		row, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return nil, err
		}
		v := &SearchSampleMessagesResp{}
		var sampleMessageID string
		if err := row.ColumnByName("SampleMessageID", &sampleMessageID); err != nil {
			return nil, fmt.Errorf("failed get sampleMessageID column : %w", err)
		}
		v.SampleMessageID = sampleMessageID

		var tags []string
		if err := row.ColumnByName("Tags", &tags); err != nil {
			return nil, fmt.Errorf("failed get Tags column : %w", err)
		}
		v.Tags = tags

		var title *string
		if err := row.ColumnByName("Title", &title); err != nil {
			return nil, fmt.Errorf("failed get Title column : %w", err)
		}
		if title != nil {
			v.Title = *title
		}

		var message string
		if err := row.ColumnByName("Message", &message); err != nil {
			return nil, fmt.Errorf("failed get message column : %w", err)
		}
		v.Message = message

		var createdAt time.Time
		if err := row.ColumnByName("CreatedAt", &createdAt); err != nil {
			return nil, fmt.Errorf("failed get CreatedAt column : %w", err)
		}
		v.CreatedAt = createdAt

		rets = append(rets, v)
	}
	return rets, nil
}
