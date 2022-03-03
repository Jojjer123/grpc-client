package main

import (
	"context"
	"fmt"
	"sync"

	adminServer "github.com/onosproject/grpc-client/cmd/generated"
	"google.golang.org/grpc"
)

const numberOfComponents = 7

func main() {
	var waitGroup sync.WaitGroup
	waitGroup.Add(2)

	var conn *grpc.ClientConn

	conn, err := grpc.Dial("10.244.0.16:4040", grpc.WithInsecure())

	if err != nil {
		fmt.Println("Did not connect: ", err)
	}

	defer conn.Close()

	c := adminServer.NewMonitorAdminInterfaceClient(conn)
	go sendRequest(c, "supa action", "supa target", &waitGroup)

	c2 := adminServer.NewMonitorAdminInterfaceClient(conn)
	go sendRequest(c2, "action in uganda", "target acquired", &waitGroup)

	waitGroup.Wait()
}

func sendRequest(c adminServer.MonitorAdminInterfaceClient, action string, target string, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()
	response, err := c.MonitorDevice(context.Background(), &adminServer.MonitorMessage{Action: action, Target: target})
	if err != nil {
		fmt.Println("Error when calling MonitorDevice: ", err)
	}

	fmt.Println("Response from server: ", response.Response)
}
