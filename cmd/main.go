package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	// "github.com/gogo/protobuf/proto"
	"github.com/golang/protobuf/proto"
	"github.com/openconfig/gnmi/client"
	gclient "github.com/openconfig/gnmi/client/gnmi"
	pb "github.com/openconfig/gnmi/proto/gnmi"

	types "github.com/onosproject/grpc-client/Types"
)

func main() {
	// setCreate("Create", "192.168.1.34", "0")
	// time.Sleep(10 * time.Second)
	// setCreate("Create", "192.168.1.82", "1")

	// type Test struct {
	// 	Name string
	// 	Interval int
	// 	Path string
	// }

	// COMMENTED FOR TESTING OF XML CONVERSION
	// var config0 []types.DeviceCounters //struct

	// c := types.DeviceCounters{
	// 	Name:     "second",
	// 	Interval: 123,
	// 	Path:     "elem: <name: 'test'>",
	// }

	// config0 = append(config0, c)
	// // config0[0].Name = "second"
	// // config0[0].Interval = 123
	// // config0[0].Path = "elem: <name: 'test'>"

	// configs := []types.Conf{
	// 	types.Conf{
	// 		Counter: config0,
	// 	},
	// }

	// conf := types.ConfigRequest{
	// 	DeviceIP: "192.168.1.82",
	// 	Configs:  configs,
	// }

	// setUpdate(conf)

	testing()

	for {
		time.Sleep(10 * time.Second)
	}
}

func testing() {
	ctx := context.Background()

	address := []string{"storage-service:11161"}

	c, err := gclient.New(ctx, client.Destination{
		Addrs:       address,
		Target:      "storage-service",
		Timeout:     time.Second * 5,
		Credentials: nil,
		TLS:         nil,
	})

	if err != nil {
		// fmt.Errorf("could not create a gNMI client: %v", err)
		fmt.Print("Could not create a gNMI client: ")
		fmt.Println(err)
	}

	// var updateList []*pb.Update

	// data := pb.TypedValue_StringVal{
	// 	StringVal: "Create", // Set command
	// }

	// actionMap := make(map[string]string)
	// actionMap["Action"] = action

	// pathElements := []*pb.PathElem{}

	// pathElements = append(pathElements, &pb.PathElem{
	// 	Name: "Action",
	// 	// Key:  actionMap,
	// })

	// configMap := make(map[string]string)
	// configMap["ConfigIndex"] = confIndex

	// pathElements = append(pathElements, &pb.PathElem{
	// 	Name: "ConfigIndex",
	// 	// Key:  configMap,
	// })

	// // TODO: REWORK, could place it all in Path
	// update := pb.Update{
	// 	Path: &pb.Path{
	// 		Elem:   pathElements,
	// 	},
	// 	// Val: &pb.TypedValue{
	// 	// 	Value: &data,
	// 	// },
	// }

	// updateList = append(updateList, &update)

	getRequest := pb.GetRequest{}

	response, err := c.(*gclient.Client).Get(ctx, &getRequest)
	if err != nil {
		fmt.Print("Target returned RPC error for Testing: ")
		fmt.Println(err)
	}

	// fmt.Print("Response from storage-service is: ")
	// fmt.Println(response)

	// schema := Schema{}
	// dec := gob.NewDecoder(bytes.NewReader(response.Notification[0].Update[0].Val.GetBytesVal()))
	// err = dec.Decode(&schema)
	// if err != nil {
	// 	fmt.Println("Failed to decode byte slice!")
	// 	fmt.Println(err)
	// }

	var schema Schema

	json.Unmarshal(response.Notification[0].Update[0].Val.GetBytesVal(), &schema)

	// fmt.Println(schema)
	getTreeStructure(schema)
	// schemaTree := getTreeStructure(schema)

	// fmt.Println(schemaTree)

	// fmt.Println("Client connected successfully")
}

// TODO: add pointer that traverse the tree based on tags, use that pointer to
// get correct parents.
func getTreeStructure(schema Schema) *SchemaTree {
	tree := &SchemaTree{}
	for index, entry := range schema.Entries {
		fmt.Println("-------------------")
		if index == 0 {
			// fmt.Println("in data")
			tree.Name = entry.Name
			tree.Namespace = entry.Namespace
			fmt.Println(tree.Name)
			tree = &SchemaTree{Parent: tree}
			continue
		}
		if entry.Value == "" { //&& schema.Entries[index + 1].ParentName  { // Directory
			if entry.Tag == "end" {
				// tree = tree.Parent
				continue
			}
			// fmt.Println("in dir")
			tree.Name = entry.Name
			tree.Namespace = entry.Namespace
			tree.Parent.Children = append(tree.Parent.Children, tree)

			fmt.Println(tree.Name)
			// fmt.Println(tree.Parent.Name)
			// fmt.Println(tree)

			tree = &SchemaTree{Parent: tree}
		} else { // Leaf
			// fmt.Println("in leaf")
			tree.Name = entry.Name
			tree.Value = entry.Value
			tree.Parent.Children = append(tree.Parent.Children, tree)
			fmt.Println(tree.Name)
			fmt.Println(tree.Value)
			// fmt.Println(tree.Parent.Name)
			// tree = tree.Parent
		}
		// fmt.Println("-------------------")
		// fmt.Println(entry)
		// fmt.Println("###################")
	}
	return tree
}

// func (s *Schema) UnmarshalBinary(data []byte) (err error) {
// 	dec := gob.NewDecoder(bytes.NewReader(data))
// 	if err = dec.Decode(&s.Entry); err != nil {
// 		return
// 	}
// 	if err = dec.Decode(&s.Children); err != nil {
// 		return
// 	}
// 	// var isCyclic bool
// 	// if err = dec.Decode(&isCyclic); err != nil {
// 	//     return
// 	// }
// 	// err = dec.Decode(&p.Q)
// 	// if isCyclic {
// 	//     p.Q.P = p
// 	// }
// 	return
// }

// type Schema struct {
// 	Entry    *SchemaEntry
// 	Children map[string]interface{}
// 	Parent   *Schema
// }

type SchemaTree struct {
	Name      string
	Namespace string
	Children  []*SchemaTree
	Parent    *SchemaTree
	Value     string
}

type Schema struct {
	Entries []SchemaEntry
}

type SchemaEntry struct {
	Name      string
	Tag       string
	Namespace string
	Value     string
}

func setCreate(action string, target string, confIndex string) {
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

	// data := pb.TypedValue_StringVal{
	// 	StringVal: "Create", // Set command
	// }

	actionMap := make(map[string]string)
	actionMap["Action"] = action

	pathElements := []*pb.PathElem{}

	pathElements = append(pathElements, &pb.PathElem{
		Name: "Action",
		Key:  actionMap,
	})

	configMap := make(map[string]string)
	configMap["ConfigIndex"] = confIndex

	pathElements = append(pathElements, &pb.PathElem{
		Name: "ConfigIndex",
		Key:  configMap,
	})

	// TODO: REWORK, could place it all in Path
	update := pb.Update{
		Path: &pb.Path{
			Target: target,
			Elem:   pathElements,
		},
		// Val: &pb.TypedValue{
		// 	Value: &data,
		// },
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
}
func setUpdate(config types.ConfigRequest) {
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

	actionMap := make(map[string]string)
	actionMap["Action"] = "Change config"

	pathElements := []*pb.PathElem{}

	pathElements = append(pathElements, &pb.PathElem{
		Name: "Action",
		Key:  actionMap,
	})

	configMap := make(map[string]string)
	configMap["DeviceIP"] = config.DeviceIP
	configMap["DeviceName"] = config.DeviceName
	configMap["Protocol"] = config.Protocol

	// pathElements := []*pb.PathElem{}

	pathElements = append(pathElements, &pb.PathElem{
		Name: "Info",
		Key:  configMap,
	})

	for confIndex, conf := range config.Configs {
		counterMap := make(map[string]string)
		for counterIndex, counter := range conf.Counter {
			counterMap["Name"+strconv.Itoa(counterIndex)] = counter.Name
			counterMap["Interval"+strconv.Itoa(counterIndex)] = strconv.Itoa(counter.Interval)
			counterMap["Path"+strconv.Itoa(counterIndex)] = counter.Path
		}

		pathElements = append(pathElements, &pb.PathElem{
			Name: "Config" + strconv.Itoa(confIndex),
			Key:  counterMap,
		})
	}

	update := pb.Update{
		Path: &pb.Path{
			Target: config.DeviceIP,
			Elem:   pathElements,
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
