package spanner_demo_20240801

import (
	"context"
	"fmt"
	"log"
	"testing"

	"cloud.google.com/go/spanner"
	"github.com/k0kubun/pp/v3"
)

func TestService_Insert(t *testing.T) {
	ctx := context.Background()

	const msg = `これはGo言語の本です。この本を読み終わるとあなたはGo言語についての理解が深まり、沼に沈んでいったことを自覚します。`

	dbName := fmt.Sprintf("projects/%s/instances/%s/databases/%s", "gcpug-public-spanner", "merpay-sponsored-instance", "sinmetal")
	dspc := spanner.DefaultSessionPoolConfig
	dspc.MinOpened = 3
	spa, err := spanner.NewClientWithConfig(ctx, dbName,
		spanner.ClientConfig{
			SessionPoolConfig: dspc,
		})
	if err != nil {
		log.Fatalln(err)
	}

	s, err := NewService(ctx, spa)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = s.Insert(ctx, &SampleMessage{
		Title:   "Goを完全に理解することができる可能性がある本",
		Tags:    []string{"Go", "Software"},
		Message: msg,
	})
	if err != nil {
		log.Fatalln(err)
	}
}

func TestService_SearchMessage(t *testing.T) {
	ctx := context.Background()

	dbName := fmt.Sprintf("projects/%s/instances/%s/databases/%s", "gcpug-public-spanner", "merpay-sponsored-instance", "sinmetal")
	dspc := spanner.DefaultSessionPoolConfig
	dspc.MinOpened = 3
	spa, err := spanner.NewClientWithConfig(ctx, dbName,
		spanner.ClientConfig{
			SessionPoolConfig: dspc,
		})
	if err != nil {
		log.Fatalln(err)
	}

	s, err := NewService(ctx, spa)
	if err != nil {
		log.Fatalln(err)
	}

	got, err := s.SearchMessage(ctx, "横浜")
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range got {
		pp.Println(v)
	}
}

func TestService_SearchSampleMessages(t *testing.T) {
	ctx := context.Background()

	dbName := fmt.Sprintf("projects/%s/instances/%s/databases/%s", "gcpug-public-spanner", "merpay-sponsored-instance", "sinmetal")
	dspc := spanner.DefaultSessionPoolConfig
	dspc.MinOpened = 3
	spa, err := spanner.NewClientWithConfig(ctx, dbName,
		spanner.ClientConfig{
			SessionPoolConfig: dspc,
		})
	if err != nil {
		log.Fatalln(err)
	}

	s, err := NewService(ctx, spa)
	if err != nil {
		log.Fatalln(err)
	}

	got, err := s.SearchSampleMessages(ctx, &SearchSampleMessagesReq{
		Tag:     "",
		Title:   "",
		Message: "Go言語",
	})
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range got {
		pp.Println(v)
	}
}
