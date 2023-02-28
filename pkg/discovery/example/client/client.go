package main

import (
	"context"
	"fmt"
	"github.com/treeforest/zut.evidence/pkg/discovery"
	"github.com/treeforest/zut.evidence/pkg/discovery/example/pb"
	"time"
)

func main() {
	dis, err := discovery.New(discovery.BaseUrl)
	if err != nil {
		panic(err)
	}
	defer dis.Close()

	conn, err := dis.Dial("Greeter")
	if err != nil {
		panic(err)
	}
	greeterClient := pb.NewGreeterClient(conn)

	for {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
		reply, err := greeterClient.Hello(ctx, &pb.HelleReq{Name: "tony"})
		if err != nil {
			cancel()
			panic(err)
		}
		cancel()
		fmt.Println(reply.Say)
		time.Sleep(time.Millisecond * 500)
	}
}
