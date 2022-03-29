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

	address := []string{"10.244.0.235:11161"}

	c, err := gclient.New(ctx, client.Destination{
		Addrs:       address,
		Target:      "10.244.0.235",
		Timeout:     time.Second * 5,
		Credentials: nil,
		TLS:         nil,
	})

	if err != nil {
		// fmt.Errorf("could not create a gNMI client: %v", err)
		fmt.Print("Could not create a gNMI client: ")
		fmt.Println(err)
	}

	getRequest := pb.GetRequest{}

	response, err := c.(*gclient.Client).Get(ctx, &getRequest)
	if err != nil {
		fmt.Print("Target returned RPC error for Testing: ")
		fmt.Println(err)
	}

	var schema Schema
	json.Unmarshal(response.Notification[0].Update[0].Val.GetBytesVal(), &schema)

	// fmt.Println(schema)
	schemaTree := getTreeStructure(schema)

	// fmt.Println("#######################")
	fmt.Println(schemaTree.Name)
	fmt.Println("--------")
	for _, child := range schemaTree.Children {
		fmt.Print(" - ")
		fmt.Print(child.Name)
		fmt.Print(", ")
		fmt.Println(child.Namespace)
	}

	// fmt.Println(schemaTree)

	// fmt.Println("Client connected successfully")
}

// TODO: add pointer that traverse the tree based on tags, use that pointer to
// get correct parents.
func getTreeStructure(schema Schema) *SchemaTree {
	var newTree *SchemaTree
	tree := &SchemaTree{}
	lastNode := ""
	for _, entry := range schema.Entries {
		// fmt.Println("-------------------")
		// if index == 0 {
		// newTree = &SchemaTree{Parent: tree}
		// newTree.Name = entry.Name
		// newTree.Namespace = entry.Namespace
		// fmt.Println(tree.Name)
		// tree = &SchemaTree{Parent: tree}
		// continue
		// }
		if entry.Value == "" { // Directory
			if entry.Tag == "end" {
				if entry.Name != "data" {
					if lastNode != "leaf" {
						// fmt.Println(tree.Name)
						tree = tree.Parent
					}
					lastNode = ""
					// continue
				}
			} else {

				newTree = &SchemaTree{Parent: tree}

				newTree.Name = entry.Name
				newTree.Namespace = entry.Namespace
				newTree.Parent.Children = append(newTree.Parent.Children, newTree)

				tree = newTree

				// fmt.Print(tree.Name)
				// fmt.Print(", ")
				// fmt.Println(tree.Parent.Name)

				// tree = &SchemaTree{Parent: tree}
			}
		} else { // Leaf
			newTree = &SchemaTree{Parent: tree}

			newTree.Name = entry.Name
			newTree.Value = entry.Value
			newTree.Parent.Children = append(newTree.Parent.Children, newTree)

			// fmt.Print(newTree.Name)
			// fmt.Print(", ")
			// fmt.Println(newTree.Parent.Name)
			// fmt.Println(newTree.Value)

			lastNode = "leaf"
		}
		// fmt.Println("-------------------")
		// fmt.Print("name: ")
		// fmt.Print(tree.Name)
		// if tree.Name != "data" {
		// 	fmt.Print(", parent: ")
		// 	fmt.Println(tree.Parent.Name)

		// 	// fmt.Println("#######")
		// 	// for i, child := range tree.Parent.Children {
		// 	// 	if i < 10 {
		// 	// 		fmt.Print(child.Name)
		// 	// 		fmt.Print(", ")
		// 	// 	}
		// 	// }
		// 	// fmt.Println("\n******")
		// 	// for j, child := range tree.Children {
		// 	// 	if j < 10 {
		// 	// 		fmt.Print(child.Name)
		// 	// 		fmt.Print(": ")
		// 	// 		fmt.Print(child.Value)
		// 	// 		fmt.Print(", ")
		// 	// 	}
		// 	// }
		// } else {
		// 	// fmt.Println("")
		// 	// for _, child := range tree.Children {
		// 	// 	fmt.Print(child.Name)
		// 	// 	fmt.Print(" | ")
		// 	// }
		// }
		// fmt.Println("")
		// fmt.Println(entry)
		// fmt.Println(tree.Namespace)
		// fmt.Println("###################")
	}
	return tree
}

// func getTreeStructure(schema Schema) *SchemaTree {
// 	tree := &SchemaTree{}
// 	lastNode := ""
// 	for index, entry := range schema.Entries {
// 		// fmt.Println("-------------------")
// 		if index == 0 {
// 			tree.Name = entry.Name
// 			tree.Namespace = entry.Namespace
// 			// fmt.Println(tree.Name)
// 			tree = &SchemaTree{Parent: tree}
// 			continue
// 		}
// 		if entry.Value == "" { // Directory
// 			if entry.Tag == "end" {
// 				if lastNode != "leaf" {
// 					tree = tree.Parent
// 				}
// 				lastNode = ""
// 				continue
// 			}
// 			tree.Name = entry.Name
// 			tree.Namespace = entry.Namespace
// 			tree.Parent.Children = append(tree.Parent.Children, tree)

// 			// fmt.Print(tree.Name)
// 			// fmt.Print(", ")
// 			// fmt.Println(tree.Parent.Name)

// 			tree = &SchemaTree{Parent: tree}
// 		} else { // Leaf
// 			// NEED TO CREATE A NEW OBJECT FOR THE LEAVES???
// 			if lastNode == "leaf" {
// 				fmt.Println("************")
// 				tree = &SchemaTree{Parent: tree.Parent}
// 				for i, child := range tree.Parent.Children {
// 					if i < 5 {
// 						fmt.Print(child.Name)
// 						fmt.Print(" | ")
// 					}
// 				}
// 				fmt.Println("************")
// 			}

// 			tree.Name = entry.Name
// 			tree.Value = entry.Value
// 			tree.Parent.Children = append(tree.Parent.Children, tree)

// 			// fmt.Print(tree.Name)
// 			// fmt.Print(", ")
// 			// fmt.Println(tree.Parent.Name)
// 			// fmt.Println(tree.Value)

// 			lastNode = "leaf"
// 		}
// 		fmt.Println("-------------------")
// 		fmt.Print("name: ")
// 		fmt.Print(tree.Name)
// 		fmt.Print(", parent: ")
// 		fmt.Println(tree.Parent.Name)
// 		fmt.Println("#######")
// 		for i, child := range tree.Parent.Children {
// 			if i < 10 {
// 				fmt.Print(child.Name)
// 				fmt.Print(", ")
// 			}
// 		}
// 		fmt.Println("")
// 		// fmt.Println("###################")
// 	}
// 	return tree
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
