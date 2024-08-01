package spanner_demo_20240801

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"time"

	"cloud.google.com/go/spanner"
	"google.golang.org/api/iterator"
)

const SampleMessagesTable = "SampleMessages"

type SampleMessage struct {
	SampleMessageID string    `json:"sampleMessageID"`
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

func (s *Service) SearchMessage(ctx context.Context, text string) ([]*SampleMessage, error) {
	const sql = `SELECT SampleMessageID, Message, CreatedAt
 FROM SampleMessages
 WHERE SEARCH(SampleMessages_Message_Tokens, @text)
`
	sts := spanner.NewStatement(sql)
	sts.Params = map[string]interface{}{
		"text": text,
	}

	var rets []*SampleMessage
	iter := s.spa.Single().Query(ctx, sts)
	for {
		row, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return nil, err
		}
		var v SampleMessage
		if err := row.ToStruct(&v); err != nil {
			return nil, err
		}
		rets = append(rets, &v)
	}
	return rets, nil
}
