package main

import (
	"context"
	"fmt"
	"time"

	// "github.com/gogo/protobuf/proto"
	"github.com/golang/protobuf/proto"
	"github.com/openconfig/gnmi/client"
	gclient "github.com/openconfig/gnmi/client/gnmi"
	pb "github.com/openconfig/gnmi/proto/gnmi"
)

func main() {
	set()
}

func set() {
	ctx := context.Background()

	address := []string{"device-monitor:11161"}

	c, err := gclient.New(ctx, client.Destination{
		Addrs:       address,
		Target:      "device-monitor",
		Timeout:     time.Second * 5,
		Credentials: nil,
		TLS:         nil,
	})

	if err != nil {
		// fmt.Errorf("could not create a gNMI client: %v", err)
		fmt.Print("Could not create a gNMI client: ")
		fmt.Println(err)
	}

	var updateList []*pb.Update

	data := pb.TypedValue_StringVal{
		StringVal: "set request data here",
	}

	// TODO: REWORK, could place it all in Path
	update := pb.Update{
		Path: &pb.Path{
			Target: "switch 123",
		},
		Val: &pb.TypedValue{
			Value: &data,
		},
	}

	updateList = append(updateList, &update)

	setRequest := pb.SetRequest{
		Update: updateList,
	}

	response, err := c.(*gclient.Client).Set(ctx, &setRequest)

	fmt.Print("Response from device-monitor is: ")
	fmt.Println(response)

	if err != nil {
		fmt.Print("Target returned RPC error for Set: ")
		fmt.Println(err)
	}

	fmt.Println("Client connected successfully")

	for {
		// fmt.Println(c.(*gclient.Client).Recv())
		time.Sleep(10 * time.Second)
	}
}

func sub() {
	ctx := context.Background()

	address := []string{"device-monitor:11161"}

	c, err := gclient.New(ctx, client.Destination{
		Addrs:       address,
		Target:      "device-monitor",
		Timeout:     time.Second * 5,
		Credentials: nil,
		TLS:         nil,
	})

	if err != nil {
		// fmt.Errorf("could not create a gNMI client: %v", err)
		fmt.Print("Could not create a gNMI client: ")
		fmt.Println(err)
	}

	var path []client.Path

	path = append(path, []string{"test/testing/lol"})

	query := client.Query{
		Addrs:  address,
		Target: "supa-target",
		// Replica: ,
		UpdatesOnly: true,
		Queries:     path,
		Type:        3, // 1 - Once, 2 - Poll, 3 - Stream
		// Timeout: ,
		// NotificationHandler: callback,
		ProtoHandler: protoCallback,
		// ...
	}

	err = c.(*gclient.Client).Subscribe(ctx, query)

	if err != nil {
		fmt.Print("Target returned RPC error for Subscribe: ")
		fmt.Println(err)
	}

	fmt.Println("Client connected successfully")

	for {
		fmt.Println(c.(*gclient.Client).Recv())
		time.Sleep(10 * time.Second)
	}
}

// func callback(msg client.Notification) error {
// 	fmt.Print("callback msg: ")
// 	fmt.Println(msg)
// 	return nil
// }

// Updates will be sent here,
func protoCallback(msg proto.Message) error {
	fmt.Print("protoCallback msg: ")
	fmt.Println(msg)
	return nil
}
