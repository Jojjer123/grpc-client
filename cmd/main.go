package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/Juniper/go-netconf/netconf"
	"github.com/golang/protobuf/proto"
	"github.com/openconfig/gnmi/client"
	gclient "github.com/openconfig/gnmi/client/gnmi"
	pb "github.com/openconfig/gnmi/proto/gnmi"

	types "github.com/onosproject/grpc-client/Types"
	"golang.org/x/crypto/ssh"
	// "github.com/vrgakos/go-netconf-client/netconf"
	// "github.com/vrgakos/go-netconf-client/netconf/message"
)

func main() {
	fmt.Println("Starting to test client")

	// setReq("Start", "192.168.0.1", "0")
	// time.Sleep(30 * time.Second)

	// // setReq("Update", "192.168.0.1", "1")

	// setReq("Stop", "192.168.0.1")

	sshConfig := &ssh.ClientConfig{
		User:            "root",
		Auth:            []ssh.AuthMethod{ssh.Password("")},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	s, err := netconf.DialSSH("192.168.0.1", sshConfig)

	if err != nil {
		log.Fatal(err)
	}

	defer s.Close()

	// fmt.Println(s.ServerCapabilities)
	// fmt.Println(s.SessionID)

	capabilities := netconf.DefaultCapabilities
	s.Transport.SendHello(&netconf.HelloMessage{Capabilities: capabilities})

	start := time.Now().UnixNano()
	reply, err := s.Exec(netconf.RawMethod("<get><filter type='subtree'><interfaces xmlns='urn:ietf:params:xml:ns:yang:ietf-interfaces'><interface><name>sw0p1</name><ethernet xmlns='urn:ieee:std:802.3:yang:ieee802-ethernet-interface'><statistics><frame><in-total-frames></in-total-frames></frame></statistics></ethernet></interface></interfaces></filter></get>"))
	end := time.Now().UnixNano()

	fmt.Printf("Delay: %v", end-start)

	if err != nil {
		panic(err)
	}
	fmt.Printf("Reply: %+v", reply)

	fmt.Println("Done!")

	for {
		time.Sleep(10 * time.Second)
	}
}

// func testNetconfClient() {
// 	fmt.Println("Creating session...")
// 	session := createSession()
// 	fmt.Println("Session created!")
// 	defer session.Close()
// 	execRPC(session)
// }

// func execRPC(session *netconf.Session) {
// 	// gt := message.NewGet("", "")
// 	// gt := message.NewGet("", "<interfaces xmlns=\"urn:ietf:params:xml:ns:yang:ietf-interfaces\"><interface><name>sw0p1</name><ethernet xmlns=\"urn:ieee:std:802.3:yang:ieee802-ethernet-interface\"><statistics><frame><in-total-frames/></frame></statistics></ethernet></interface></interfaces>")
// 	gt := message.NewRPC(`<get><filter type="subtree"><interfaces xmlns="urn:ietf:params:xml:ns:yang:ietf-interfaces"><interface><name>sw0p1</name><ethernet xmlns="urn:ieee:std:802.3:yang:ieee802-ethernet-interface"><statistics><frame><in-total-frames></in-total-frames></frame></statistics></ethernet></interface></interfaces></filter></get>`)
// 	fmt.Println("Message created!")
// 	start := time.Now().UnixNano()
// 	session.AsyncRPC(gt, defaultLogRpcReplyCallback(gt.MessageID, start))
// 	time.Sleep(100 * time.Millisecond)

// 	fmt.Printf("MessageID: %v\n", gt.MessageID)
// 	// fmt.Printf("delay: %v\n", time.Now().UnixNano()-start)
// 	// if err != nil {
// 	// 	fmt.Printf("Failed RPC: %v\n", err)
// 	// } else {
// 	// 	fmt.Println(reply.RawReply)
// 	// }

// 	// d2 := message.NewCloseSession()
// 	// start2 := time.Now().UnixNano()
// 	// session.AsyncRPC(d2, defaultLogRpcReplyCallback(d2.MessageID, start2))

// 	session.Listener.WaitForMessages()
// }

// func createSession() *netconf.Session {
// 	sshConfig := &ssh.ClientConfig{
// 		User:            "root",
// 		Auth:            []ssh.AuthMethod{ssh.Password("")},
// 		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
// 	}
// 	s, err := netconf.DialSSH("192.168.0.1:830", sshConfig)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	capabilities := netconf.DefaultCapabilities
// 	err = s.SendHello(&message.Hello{Capabilities: capabilities})
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	return s
// }

// func defaultLogRpcReplyCallback(eventId string, start int64) netconf.Callback {
// 	return func(event netconf.Event) {
// 		reply := event.RPCReply()
// 		fmt.Printf("delay for event %v: %v\n", eventId, time.Now().UnixNano()-start)
// 		if reply == nil {
// 			println("Failed to execute RPC")
// 		}
// 		if event.EventID() == eventId {
// 			println("Successfully executed RPC")
// 			println(reply.RawReply)
// 		}
// 	}
// }

func testSequences() {
	// fmt.Println("Start batch monitoring on switch_one, switch_two, and switch_three")
	// setReq("Start", "192.168.0.1", "0")
	// setReq("Start", "192.168.0.2", "0")
	// setReq("Start", "192.168.0.3", "0")
	// time.Sleep(2 * time.Minute)
	// setReq("Stop", "192.168.0.1")
	// setReq("Stop", "192.168.0.2")
	// setReq("Stop", "192.168.0.3")

	// time.Sleep(1 * time.Minute)

	// fmt.Println("Start batch monitoring on switch_one, switch_two, and switch_three")
	// setReq("Start", "192.168.0.1", "0")
	// setReq("Start", "192.168.0.2", "0")
	// setReq("Start", "192.168.0.3", "0")
	// time.Sleep(1 * time.Minute)
	// setReq("Stop", "192.168.0.1")
	// setReq("Stop", "192.168.0.2")
	// setReq("Stop", "192.168.0.3")

	// time.Sleep(1 * time.Minute)

	// fmt.Println("Start update test")
	// setReq("Start", "192.168.0.1", "0")
	// setReq("Start", "192.168.0.2", "0")
	// setReq("Start", "192.168.0.3", "0")

	// time.Sleep(30 * time.Second)

	// fmt.Println("Starting to update now")

	// setReq("Update", "192.168.0.1", "1")
	// setReq("Update", "192.168.0.2", "1")
	// setReq("Update", "192.168.0.3", "1")

	// time.Sleep(30 * time.Second)

	// fmt.Println("Updating second time")

	// setReq("Update", "192.168.0.1", "0")
	// setReq("Update", "192.168.0.2", "0")
	// setReq("Update", "192.168.0.3", "0")

	// time.Sleep(60 * time.Second)

	// fmt.Println("Updating third time")

	// setReq("Update", "192.168.0.1", "1")
	// setReq("Update", "192.168.0.2", "1")
	// setReq("Update", "192.168.0.3", "1")

	// time.Sleep(30 * time.Second)

	// fmt.Println("Updating fourth time")

	// setReq("Update", "192.168.0.1", "0")
	// setReq("Update", "192.168.0.2", "0")
	// setReq("Update", "192.168.0.3", "0")

	// time.Sleep(60 * time.Second)

	// fmt.Println("Updating fifth and final time")

	// setReq("Update", "192.168.0.1", "1")
	// setReq("Update", "192.168.0.2", "1")
	// setReq("Update", "192.168.0.3", "1")

	// fmt.Println("Start non-batch monitoring on switch_one, switch_two, and switch_three")
	// setReq("Start", "192.168.0.1", "1")
	// setReq("Start", "192.168.0.2", "1")
	// setReq("Start", "192.168.0.3", "1")
	// time.Sleep(2 * time.Minute)
	// setReq("Stop", "192.168.0.1")
	// setReq("Stop", "192.168.0.2")
	// setReq("Stop", "192.168.0.3")

	setReq("Start", "192.168.0.1", "0")
	time.Sleep(30 * time.Second)
	setReq("Stop", "192.168.0.1")
}

func getFullConfig() {
	ctx := context.Background()

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

	getRequest := pb.GetRequest{
		Path: []*pb.Path{
			{
				Target: "192.168.0.2",
			},
		},
		Type: pb.GetRequest_STATE,
	}

	response, err := c.(*gclient.Client).Get(ctx, &getRequest)
	if err != nil {
		fmt.Print("Target returned RPC error for Testing: ")
		fmt.Println(err)
	}

	c.Close()

	fmt.Println(response)
}

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
			{
				// Elem: []*pb.PathElem{
				// 	{
				// 		Name: "interfaces",
				// 		Key:  interfaceKeyMap,
				// 	},
				// 	{
				// 		Name: "interface",
				// 		Key:  loopbackKeyMap,
				// 	},
				// },
				Target: "192.168.0.2",
			},
		},
		Type: pb.GetRequest_STATE,
	}

	response, err := c.(*gclient.Client).Get(ctx, &getRequest)
	if err != nil {
		fmt.Print("Target returned RPC error for Testing: ")
		fmt.Println(err)
	}

	c.Close()

	fmt.Println(response)

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
							"name": "sw0p1",
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

	fmt.Println(response.Notification[0].Update[0])

	// Send get request to adapter

	address = []string{"gnmi-netconf-adapter:11161"}

	c, err = gclient.New(ctx, client.Destination{
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

	getRequest = pb.GetRequest{
		Path: []*pb.Path{
			response.Notification[0].Update[0].Path,
		},
		Type: pb.GetRequest_STATE,
	}

	response, err = c.(*gclient.Client).Get(ctx, &getRequest)
	if err != nil {
		fmt.Print("Target returned RPC error for Testing: ")
		fmt.Println(err)
	}

	c.Close()

	// fmt.Println(response)

	var schema Schema
	var schemaTree *SchemaTree
	if len(response.Notification) > 0 {
		json.Unmarshal(response.Notification[0].Update[0].Val.GetBytesVal(), &schema)

		// fmt.Println(schema)
		schemaTree = getTreeStructure(schema)
	}

	printSchemaTree(schemaTree)
}

func printSchemaTree(schemaTree *SchemaTree) {
	fmt.Printf("%s - %s - %v - %s\n", schemaTree.Parent.Name, schemaTree.Name, schemaTree.Namespace, schemaTree.Value)
	for _, child := range schemaTree.Children {
		printSchemaTree(child)
	}
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

// func setDelete(action string, target string) {
// 	ctx := context.Background()
//
// 	address := []string{"monitor-service:11161"}
//
// 	c, err := gclient.New(ctx, client.Destination{
// 		Addrs:       address,
// 		Target:      "monitor-service",
// 		Timeout:     time.Second * 5,
// 		Credentials: nil,
// 		TLS:         nil,
// 	})
//
// 	if err != nil {
// 		// fmt.Errorf("could not create a gNMI client: %v", err)
// 		fmt.Print("Could not create a gNMI client: ")
// 		fmt.Println(err)
// 	}
//
// 	actionMap := make(map[string]string)
// 	actionMap["Action"] = action
//
// 	setRequest := pb.SetRequest{
// 		Update: []*pb.Update{
// 			{
// 				Path: &pb.Path{
// 					Target: target,
// 					Elem: []*pb.PathElem{
// 						{
// 							Name: "Action",
// 							Key:  actionMap,
// 						},
// 					},
// 				},
// 			},
// 		},
// 	}
//
// 	response, err := c.(*gclient.Client).Set(ctx, &setRequest)
//
// 	fmt.Print("Response from device-monitor is: ")
// 	fmt.Println(response)
//
// 	if err != nil {
// 		fmt.Print("Target returned RPC error for Set: ")
// 		fmt.Println(err)
// 	}
//
// 	fmt.Println("Client connected successfully")
// }

func setReq(action string, target string, confIndex ...string) {
	ctx := context.Background()

	address := []string{"monitor-service:11161"}

	c, err := gclient.New(ctx, client.Destination{
		Addrs:       address,
		Target:      "monitor-service",
		Timeout:     time.Second * 5,
		Credentials: nil,
		TLS:         nil,
	})

	if err != nil {
		// fmt.Errorf("could not create a gNMI client: %v", err)
		fmt.Print("Could not create a gNMI client: ")
		fmt.Println(err)
	}

	setRequest := pb.SetRequest{
		Update: []*pb.Update{
			{
				Path: &pb.Path{
					Target: target,
					Elem: []*pb.PathElem{
						{
							Name: "Action",
							Key: map[string]string{
								"Action": action,
							},
						},
					},
				},
			},
		},
	}

	if confIndex != nil {
		setRequest.Update[0].Path.Elem = append(setRequest.Update[0].Path.Elem, &pb.PathElem{
			Name: "ConfigIndex",
			Key: map[string]string{
				"ConfigIndex": confIndex[0],
			},
		})
	}

	response, err := c.(*gclient.Client).Set(ctx, &setRequest)

	fmt.Print("Response from device-monitor is: ")
	fmt.Println(response)

	if err != nil {
		fmt.Print("Target returned RPC error for Set: ")
		fmt.Println(err)
	}

	// fmt.Println("Client connected successfully")
}

func setUpdate(config types.ConfigRequest) {
	ctx := context.Background()

	address := []string{"monitor-service:11161"}

	c, err := gclient.New(ctx, client.Destination{
		Addrs:       address,
		Target:      "monitor-service",
		Timeout:     time.Second * 5,
		Credentials: nil,
		TLS:         nil,
	})

	if err != nil {
		// fmt.Errorf("could not create a gNMI client: %v", err)
		fmt.Print("Could not create a gNMI client: ")
		fmt.Println(err)
	}

	pathElements := []*pb.PathElem{
		{
			Name: "Action",
			Key: map[string]string{
				"Action": "Change config",
			},
		},
		{
			Name: "Info",
			Key: map[string]string{
				"DeviceIP":   config.DeviceIP,
				"DeviceName": config.DeviceName,
				"Protocol":   config.Protocol,
			},
		},
	}

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

	setRequest := pb.SetRequest{
		Update: []*pb.Update{
			{
				Path: &pb.Path{
					Target: config.DeviceIP,
					Elem:   pathElements,
				},
			},
		},
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

	address := []string{"monitor-service:11161"}

	c, err := gclient.New(ctx, client.Destination{
		Addrs:       address,
		Target:      "monitor-service",
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

	fmt.Println(query)

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

// Updates will be sent here,
func protoCallback(msg proto.Message) error {
	fmt.Print("protoCallback msg: ")
	fmt.Println(msg)
	return nil
}
