// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"storj.io/drpc/drpcconn"

	"storj.io/drpc/examples/drpc/pb"
)

func main() {
	err := Main(context.Background())
	if err != nil {
		panic(err)
	}
}

func Main(ctx context.Context) error {
	// dial the drpc server
	rawconn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		return err
	}
	// N.B.: If you want TLS, you need to wrap the net.Conn with TLS before
	// making a DRPC conn.

	// convert the net.Conn to a drpc.Conn
	conn := drpcconn.New(rawconn)
	defer conn.Close()

	// make a drpc proto-specific client
	client := pb.NewDRPCCookieMonsterClient(conn)

	for _, d := range []time.Duration{1, 3} {
		select {
		case <-conn.Closed():
			log.Print("closed")
		default:
			log.Print("not closed")
		}
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
			fmt.Println(crumbs.Cookie.Type.String())
		}()
	}

	return nil
}
