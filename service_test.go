package spanner_demo_20240801

import (
	"context"
	"fmt"
	"log"
	"testing"

	"cloud.google.com/go/spanner"
)

func TestService_Insert(t *testing.T) {
	ctx := context.Background()

	const msg = `これはGo言語の本です。
`

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

	err = s.Insert(ctx, &SampleMessage{
		Message: msg,
	})
	if err != nil {
		log.Fatalln(err)
	}
}
