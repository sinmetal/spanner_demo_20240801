package main

import (
	"context"

	demo "github.com/sinmetal/spanner_demo_20240801"
)

func main() {
	ctx := context.Background()

	demo.Ignition(ctx)

	//if len(os.Args) < 2 {
	//	log.Fatalln("Not enough args")
	//}
	//
	//dbName := fmt.Sprintf("projects/%s/instances/%s/databases/%s", "gcpug-public-spanner", "merpay-sponsored-instance", "sinmetal")
	//dspc := spanner.DefaultSessionPoolConfig
	//dspc.MinOpened = 3
	//spa, err := spanner.NewClientWithConfig(ctx, dbName,
	//	spanner.ClientConfig{
	//		SessionPoolConfig: dspc,
	//	})
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//
	//s, err := demo.NewService(ctx, spa)
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//
	//err = s.Insert(ctx, &demo.SampleMessage{
	//	Message: os.Args[1],
	//})
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//fmt.Println("DONE")
}
