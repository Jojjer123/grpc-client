package main

import (
	"context"
	"fmt"

	adminServer "github.com/onosproject/grpc-client/cmd/generated"
	"google.golang.org/grpc"
)

const numberOfComponents = 7

func main() {
	var conn *grpc.ClientConn

	conn, err := grpc.Dial("10.244.0.19:4040", grpc.WithInsecure())

	if err != nil {
		fmt.Println("Did not connect: ", err)
	}

	defer conn.Close()

	c := adminServer.NewMonitorAdminInterfaceClient(conn)

	response, err := c.MonitorDevice(context.Background(), &adminServer.MonitorMessage{Action: "supa-action", Target: "supa-server"})
	if err != nil {
		fmt.Println("Error when calling MonitorDevice: ", err)
	}

	fmt.Println("Response from server: ", response.Response)
}
