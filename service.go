package spanner_demo_20240801

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"time"

	"cloud.google.com/go/spanner"
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

func (s *Service) Insert(ctx context.Context, value *SampleMessage) error {
	value.SampleMessageID = s.CreateSampleMessageID(value.Message)
	value.CreatedAt = spanner.CommitTimestamp
	_, err := s.spa.ReadWriteTransaction(ctx, func(ctx context.Context, tx *spanner.ReadWriteTransaction) error {
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
		return err
	}
	return nil
}
