package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	api "github.com/tetran/proglog-example/api/v1"
)

func main() {
	addr := flag.String("addr", ":8400", "service address")
	flag.Parse()
	conn, err := grpc.Dial(*addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	client := api.NewLogClient(conn)
	ctx := context.Background()
	res, err := client.GetServers(ctx, &api.GetServersRequest{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("servers:")
	for _, server := range res.Servers {
		fmt.Printf("- %v\n", server)
	}
}
