// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"

	"storj.io/drpc/examples/grpc/pb"
)

func main() {
	err := Main(context.Background())
	if err != nil {
		panic(err)
	}
}

func Main(ctx context.Context) error {
	// dial the grpc server (without TLS)
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return err
	}
	defer conn.Close()

	// make a grpc proto-specific client
	client := pb.NewCookieMonsterClient(conn)

	for _, d := range []time.Duration{1, 3} {
		log.Println(conn.GetState())
		func() {
			// set a deadline for the operation
			ctx, cancel := context.WithTimeout(ctx, d*time.Second)
			defer cancel()

			// run the RPC
			crumbs, err := client.EatCookie(ctx, &pb.Cookie{
				Type: pb.Cookie_Oatmeal,
			})
			if err != nil {
				log.Println(err)
				return
			}

			// check the results
			_, err = fmt.Println(crumbs.Cookie.Type.String())
		}()
	}

	return nil
}
