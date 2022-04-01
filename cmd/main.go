package main

import (
	"context"
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

	// netconfConv(xmlString)

	testing()

	for {
		time.Sleep(10 * time.Second)
	}
}

// func netconfConv(xmlString string) *types.Schema {
// 	decoder := xml.NewDecoder(strings.NewReader(xmlString))
// 	schema := &types.Schema{}

// 	var newEntry *types.SchemaEntry
// 	// var lastNamespace string

// 	var nsParser *types.NamespaceParser

// 	index := 0
// 	for {
// 		tok, _ := decoder.Token()

// 		if tok == nil {
// 			fmt.Println("")
// 			return schema
// 		}

// 		switch tokType := tok.(type) {
// 		case xml.StartElement:
// 			// fmt.Println(tokType.Name.Local)
// 			newEntry = &types.SchemaEntry{}
// 			newEntry.Name = tokType.Name.Local

// 			// nsParser = &types.NamespaceParser{}

// 			if index > 0 {
// 				newNsParser := &types.NamespaceParser{
// 					Parent:              nsParser,
// 					LastParentNamespace: nsParser.LastParentNamespace,
// 				}

// 				// fmt.Print(tokType.Name.Local)
// 				// fmt.Printf(" - %v", tokType.Attr)

// 				if nsParser.LastParentNamespace != tokType.Name.Space {
// 					newNsParser.LastParentNamespace = tokType.Name.Space
// 					newEntry.Namespace = tokType.Name.Space
// 					// fmt.Printf(" - %s\n", newNsParser.LastParentNamespace)
// 				} else {
// 					// fmt.Print("\n")
// 				}

// 				nsParser.Children = append(nsParser.Children, newNsParser)

// 				nsParser = newNsParser

// 				// TODO: Fix namespaces, it currently won't add modules-state which is a module
// 				// with exactly the same namespace as the previous module yang-library for state-data
// 				// if tokType.Name.Space != lastNamespace {
// 				// 	lastNamespace = tokType.Name.Space
// 				// 	newEntry.Namespace = lastNamespace
// 				// }
// 				newEntry.Tag = "start"
// 			} else {
// 				nsParser = &types.NamespaceParser{
// 					LastParentNamespace: tokType.Name.Space,
// 					Parent:              nil,
// 				}

// 				// fmt.Printf("%s - %s\n", newEntry.Name, nsParser.LastParentNamespace)

// 				// lastNamespace = tokType.Name.Space
// 				newEntry.Namespace = tokType.Name.Space
// 			}

// 			schema.Entries = append(schema.Entries, *newEntry)
// 			index++

// 		case xml.EndElement:
// 			// fmt.Printf("Exiting %s now\n", tokType.Name.Local)
// 			nsParser = nsParser.Parent

// 			newEntry = &types.SchemaEntry{}
// 			newEntry.Name = tokType.Name.Local
// 			newEntry.Tag = "end"
// 			schema.Entries = append(schema.Entries, *newEntry)
// 			index++

// 		case xml.CharData:
// 			bytes := xml.CharData(tokType)
// 			schema.Entries[index-1].Value = string([]byte(bytes))

// 		default:
// 			fmt.Printf("Token type was not recognized with type: %v", tokType)
// 		}
// 	}
// }

func testing() {
	ctx := context.Background()

	// Send get request to adapter

	address := []string{"gnmi-netconf-adapter:11161"}

	c, err := gclient.New(ctx, client.Destination{
		Addrs:       address,
		Target:      "gnmi-netconf-adapter",
		Timeout:     time.Second * 5,
		Credentials: nil,
		TLS:         nil,
	})

	if err != nil {
		// fmt.Errorf("could not create a gNMI client: %v", err)
		fmt.Print("Could not create a gNMI client: ")
		fmt.Println(err)
	}

	// interfaceKeyMap := map[string]string{}
	// interfaceKeyMap["namespace"] = "urn:ietf:params:xml:ns:yang:ietf-interfaces"

	// loopbackKeyMap := map[string]string{}
	// loopbackKeyMap["name"] = "lo"

	getRequest := pb.GetRequest{
		Path: []*pb.Path{
			// {
			// 	Elem: []*pb.PathElem{
			// 		{
			// 			Name: "interfaces",
			// 			Key:  interfaceKeyMap,
			// 		},
			// 		{
			// 			Name: "interface",
			// 			Key:  loopbackKeyMap,
			// 		},
			// 	},
			// },
		},
		Type: pb.GetRequest_STATE,
	}

	response, err := c.(*gclient.Client).Get(ctx, &getRequest)
	if err != nil {
		fmt.Print("Target returned RPC error for Testing: ")
		fmt.Println(err)
	}

	c.Close()

	// Send set request to storage

	address = []string{"storage-service:11161"}

	c, err = gclient.New(ctx, client.Destination{
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

	//

	if len(response.Notification) <= 0 {
		return
	}

	actionMap := map[string]string{}
	actionMap["Action"] = "Store namespaces"

	setRequest := pb.SetRequest{
		Update: []*pb.Update{
			{
				Path: &pb.Path{
					Elem: []*pb.PathElem{
						{
							Name: "Action",
							Key:  actionMap,
						},
					},
				},
				Val: response.Notification[0].Update[0].Val,
			},
		},
	}

	setResponse, err := c.(*gclient.Client).Set(ctx, &setRequest)
	if err != nil {
		fmt.Print("Target returned RPC error for Testing: ")
		fmt.Println(err)
	}

	c.Close()

	fmt.Println(setResponse)

	if err != nil {
		// fmt.Errorf("could not create a gNMI client: %v", err)
		fmt.Print("Could not create a gNMI client: ")
		fmt.Println(err)
	}

	// Send get request to storage

	address = []string{"storage-service:11161"}

	c, err = gclient.New(ctx, client.Destination{
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

	getRequest = pb.GetRequest{
		Path: []*pb.Path{
			{
				Elem: []*pb.PathElem{
					{
						Name: "interfaces",
					},
					{
						Name: "interface",
						Key: map[string]string{
							"Name": "sw0p1",
						},
					},
					{
						Name: "bridge-port",
					},
					{
						Name: "traffic-class",
					},
					{
						Name: "priority0",
					},
				},
			},
		},
		Type: 4,
	}

	response, err = c.(*gclient.Client).Get(ctx, &getRequest)
	if err != nil {
		fmt.Print("Target returned RPC error for Testing: ")
		fmt.Println(err)
	}

	c.Close()

	fmt.Println(response.Notification[0].Update)

	// var schema Schema
	// if len(response.Notification) > 0 {
	// 	json.Unmarshal(response.Notification[0].Update[0].Val.GetBytesVal(), &schema)

	// 	// fmt.Println(schema)
	// 	schemaTree := getTreeStructure(schema)

	// 	// fmt.Println("#######################")
	// 	fmt.Println(schemaTree.Name)
	// 	fmt.Println("--------")
	// 	for _, child := range schemaTree.Children {
	// 		fmt.Print(" - ")
	// 		fmt.Print(child.Name)
	// 		fmt.Print(", ")
	// 		fmt.Println(child.Namespace)
	// 	}
	// }
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
