package main

import (
	"context"
	"encoding/json"
	"errors"

	// "encoding/xml"
	"fmt"

	// "os"
	// "strconv"
	"time"

	"github.com/Juniper/go-netconf/netconf"
	"github.com/atomix/atomix-go-client/pkg/atomix"
	"golang.org/x/crypto/ssh"

	// _map "github.com/atomix/atomix-go-client/pkg/atomix/map"
	// configmodel "github.com/onosproject/onos-config-model/pkg/model"

	// "github.com/golang/protobuf/proto"
	"github.com/openconfig/gnmi/client"
	gclient "github.com/openconfig/gnmi/client/gnmi"
	pb "github.com/openconfig/gnmi/proto/gnmi"

	// "github.com/openconfig/gnmi/proto/gnmi_ext"

	// types "github.com/onosproject/grpc-client/Types"

	// adapterResp "github.com/onosproject/grpc-client/adapterResponse"

	// "github.com/openconfig/goyang/pkg/yang"
	// "github.com/openconfig/ygot/ygot"

	model "github.com/onosproject/grpc-client/pkg/modelPlugin"
	// "github.com/onosproject/onos-lib-go/pkg/errors"

	_map "github.com/atomix/atomix-go-client/pkg/atomix/map"
)

func main() {
	fmt.Println("Start")

	getConfig("192.168.0.3")

	// getFullConfigFromSwitch("192.168.0.3")

	// resp := getFullConfig("192.168.0.3")
	// fmt.Printf("Config: %v", resp)

	// var adapterResponse adapterResp.AdapterResponse

	// if err := proto.Unmarshal(resp.Notification[0].Update[0].Val.GetProtoBytes(), &adapterResponse); err != nil {
	// 	fmt.Printf("Failed to unmarshal ProtoBytes: %v", err)
	// }

	// schemaTree := getNewTreeStructure(adapterResponse.Entries)

	// printTree(schemaTree, 0)

	// testApplyingConfig()

	// model, err := getModel("tsn-model:1.0.2")
	// if err != nil {
	// 	fmt.Printf("Failed loading model: %v\n", err)
	// }

	// setReq("Start", "192.168.0.2", "2")

	// time.Sleep(55 * time.Second)

	// setReq("Stop", "192.168.0.2")

	// testNetworkChangeRequest()

	// fmt.Println("End")

	for {
		time.Sleep(10 * time.Second)
	}
}

func getConfig(ip string) {
	reply, err := sendRPCRequest(netconf.MethodGetConfig("running"), ip)
	if err != nil {
		fmt.Printf("Failed sending RPC request: %v\n", err)
	}

	fmt.Printf("Config: %v\n", reply.Data)
}

// Takes in an RPCMethod function and executes it, then returns the reply from the network device
func sendRPCRequest(fn netconf.RPCMethod, switchAddr string) (*netconf.RPCReply, error) {
	//  Define config for connection to network device
	sshConfig := &ssh.ClientConfig{
		User:            "root",
		Auth:            []ssh.AuthMethod{ssh.Password("")},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// Start connection to network device
	s, err := netconf.DialSSH(switchAddr, sshConfig)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// Close connetion to network device when this function is done executing
	defer s.Close()

	reply, err := s.Exec(fn)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return reply, nil
}

// Get the model from Atomix (k/v store)
func getModel(name string) (*model.ModelPlugin, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	store, err := atomix.GetMap(ctx, "device-models")
	if err != nil {
		fmt.Printf("Failed gettings store: %v\n", err)
		return nil, err
	}

	modelFields := []string{"info", "roPaths", "rwPaths"}

	modelPlugin := &model.ModelPlugin{}

	for _, field := range modelFields {
		data, err := store.Get(ctx, name+":"+field)
		if err != nil {
			return nil, err
		}

		err = deserialize(modelPlugin, field, data)
		if err != nil {
			return nil, err
		}
	}

	return modelPlugin, nil
}

func deserialize(modelPlugin *model.ModelPlugin, field string, data *_map.Entry) error {
	switch field {
	case "info":
		err := json.Unmarshal(data.Value, &modelPlugin.Info)
		if err != nil {
			return err
		}
	case "roPaths":
		err := json.Unmarshal(data.Value, &modelPlugin.ReadOnlyPaths)
		if err != nil {
			return err
		}

	case "rwPaths":
		err := json.Unmarshal(data.Value, &modelPlugin.ReadWritePaths)
		if err != nil {
			return err
		}

	default:
		return errors.New("field not recognized")
	}

	return nil
}

// func getModel(name string) (*model.ModelPlugin, error) {
// 	ctx := context.Background()

// 	store, err := atomix.GetMap(ctx, "device-models")
// 	if err != nil {
// 		fmt.Printf("Failed gettings store: %v\n", err)
// 		return nil, err
// 	}

// 	modelFields := []string{"info", "roPaths", "rwPaths"}

// 	modelPlugin := &model.ModelPlugin{}

// 	for _, field := range modelFields {
// 		data, err := store.Get(ctx, name+":"+field)
// 		if err != nil {
// 			fmt.Printf("Failed getting resource: %v\n", err)
// 			return nil, err
// 		}

// 		switch field {
// 		case "info":
// 			err := json.Unmarshal(data.Value, &modelPlugin.Info)
// 			if err != nil {
// 				return nil, err
// 			}
// 		case "roPaths":
// 			err := json.Unmarshal(data.Value, &modelPlugin.ReadOnlyPaths)
// 			if err != nil {
// 				return nil, err
// 			}

// 		case "rwPaths":
// 			err := json.Unmarshal(data.Value, &modelPlugin.ReadWritePaths)
// 			if err != nil {
// 				return nil, err
// 			}

// 		default:
// 			return nil, err
// 		}
// 	}

// 	return modelPlugin, nil
// }

// func readChan(ch chan _map.Entry, num int) {
// 	msg := <-ch

// 	fmt.Println(msg.Key)

// 	if num < 1 {
// 		go readChan(ch, num+1)
// 	}
// }

// // ModelPlugin is a config model
// type ModelPlugin struct {
// 	Info           configmodel.ModelInfo
// 	Model          configmodel.ConfigModel
// 	ReadOnlyPaths  ReadOnlyPathMap  `json:"readOnlyPathMap"`
// 	ReadWritePaths ReadWritePathMap `json:"readWritePathMap"`
// }

// // Name is a config model name
// type Name string

// // Version is a config model version
// type Version string

// // Revision is a config module revision
// type Revision string

// // GetStateMode defines the Getstate handling
// type GetStateMode string

// const (
// 	// GetStateNone - device type does not support Operational State at all
// 	GetStateNone GetStateMode = "GetStateNone"
// 	// GetStateOpState - device returns all its op state attributes by querying
// 	// GetRequest_STATE and GetRequest_OPERATIONAL
// 	GetStateOpState GetStateMode = "GetStateOpState"
// 	// GetStateExplicitRoPaths - device returns all its op state attributes by querying
// 	// exactly what the ReadOnly paths from YANG - wildcards are handled by device
// 	GetStateExplicitRoPaths GetStateMode = "GetStateExplicitRoPaths"
// 	// GetStateExplicitRoPathsExpandWildcards - where there are wildcards in the
// 	// ReadOnly paths 2 calls have to be made - 1) to expand the wildcards in to
// 	// real paths (since the device doesn't do it) and 2) to query those expanded
// 	// wildcard paths - this is the Stratum 1.0.0 method
// 	GetStateExplicitRoPathsExpandWildcards GetStateMode = "GetStateExplicitRoPathsExpandWildcards"
// )

// func (m ModelInfo) String() string {
// 	return fmt.Sprintf("%s@%s", m.Name, m.Version)
// }

// // ModuleInfo is a config module info
// type ModuleInfo struct {
// 	Name         Name     `json:"name"`
// 	File         string   `json:"file"`
// 	Organization string   `json:"organization"`
// 	Revision     Revision `json:"revision"`
// }

// // FileInfo is a config file info
// type FileInfo struct {
// 	Path string `json:"path"`
// 	Data []byte `json:"data"`
// }

// // PluginInfo is config model plugin info
// type PluginInfo struct {
// 	Name    Name    `json:"name"`
// 	Version Version `json:"version"`
// }

// // ConfigModel is a configuration model data
// type ConfigModel interface {
// 	// Info returns the config model info
// 	Info() ModelInfo

// 	// Data returns the config model data
// 	Data() []*pb.ModelData

// 	// Schema returns the config model schema
// 	Schema() (map[string]*yang.Entry, error)

// 	// GetStateMode returns the get state mode
// 	GetStateMode() GetStateMode

// 	// Unmarshaler returns the config model unmarshaler function
// 	Unmarshaler() Unmarshaler

// 	// Validator returns the config model validator function
// 	Validator() Validator
// }

// // Unmarshaler is a config model unmarshaler function
// type Unmarshaler func([]byte) (*ygot.ValidatedGoStruct, error)

// // Validator is a config model validator function
// type Validator func(model *ygot.ValidatedGoStruct, opts ...ygot.ValidationOption) error

// // ModelInfo is config model info
// type ModelInfo struct {
// 	Name         Name         `json:"name"`
// 	Version      Version      `json:"version"`
// 	GetStateMode GetStateMode `json:"getStateMode"`
// 	Files        []FileInfo   `json:"files"`
// 	Modules      []ModuleInfo `json:"modules"`
// 	Plugin       PluginInfo   `json:"plugin"`
// }

// // ValueType is the type for a value
// type ValueType int32

// type ReadOnlyAttrib struct {
// 	ValueType   ValueType      `json:"valueType"`
// 	TypeOpts    []uint8        `json:"typeOpts"`
// 	Description string         `json:"description"`
// 	Units       string         `json:"units"`
// 	Enum        map[int]string `json:"enum"`
// 	IsAKey      bool           `json:"isAKey"`
// 	AttrName    string         `json:"attrName"`
// }

// // ReadOnlySubPathMap abstracts the read only subpath
// type ReadOnlySubPathMap map[string]ReadOnlyAttrib

// // ReadOnlyPathMap abstracts the read only path
// type ReadOnlyPathMap map[string]ReadOnlySubPathMap

// // ReadWritePathElem holds data about a leaf or container
// type ReadWritePathElem struct {
// 	ReadOnlyAttrib
// 	Mandatory bool
// 	Default   string
// 	Range     []string
// 	Length    []string
// }

// // ReadWritePathMap is a map of ReadWrite paths a their metadata
// type ReadWritePathMap map[string]ReadWritePathElem
// --------------------------------------------------------------------------------------------

func testApplyingConfig() {
	ctx := context.Background()

	address := []string{"gnmi-netconf-adapter:11161"}

	c, err := gclient.New(ctx, client.Destination{
		Addrs:       address,
		Timeout:     time.Second * 5,
		Credentials: nil,
		TLS:         nil,
	})

	if err != nil {
		fmt.Print("Could not create a gNMI client: ")
		fmt.Println(err)
	}

	setRequest := pb.SetRequest{
		Replace: []*pb.Update{
			{
				Path: &pb.Path{
					Elem: []*pb.PathElem{
						{
							Name: "interfaces",
							Key:  map[string]string{"namespace": "urn:ietf:params:xml:ns:yang:ietf-interfaces"},
						},
						{
							Name: "interface",
							Key:  map[string]string{"name": "sw0p1"},
						},
						{
							Name: "gate-parameters",
							Key:  map[string]string{"namespace": "urn:ieee:std:802.1Q:yang:ieee802-dot1q-sched"},
						},
						{
							Name: "gate-enabled",
							Key:  map[string]string{},
						},
					},
					Target: "192.168.0.3",
				},
				Val: &pb.TypedValue{
					Value: &pb.TypedValue_BoolVal{
						BoolVal: true,
					},
				},
			},
			{
				Path: &pb.Path{
					Elem: []*pb.PathElem{
						{
							Name: "interfaces",
							Key:  map[string]string{"namespace": "urn:ietf:params:xml:ns:yang:ietf-interfaces"},
						},
						{
							Name: "interface",
							Key:  map[string]string{"name": "sw0p1"},
						},
						{
							Name: "gate-parameters",
							Key:  map[string]string{"namespace": "urn:ieee:std:802.1Q:yang:ieee802-dot1q-sched"},
						},
						{
							Name: "admin-gate-states",
							Key:  map[string]string{},
						},
					},
					Target: "192.168.0.3",
				},
				Val: &pb.TypedValue{
					Value: &pb.TypedValue_UintVal{
						UintVal: 255,
					},
				},
			},
			{
				Path: &pb.Path{
					Elem: []*pb.PathElem{
						{
							Name: "interfaces",
							Key:  map[string]string{"namespace": "urn:ietf:params:xml:ns:yang:ietf-interfaces"},
						},
						{
							Name: "interface",
							Key:  map[string]string{"name": "sw0p1"},
						},
						{
							Name: "gate-parameters",
							Key:  map[string]string{"namespace": "urn:ieee:std:802.1Q:yang:ieee802-dot1q-sched"},
						},
						{
							Name: "admin-control-list-length",
							Key:  map[string]string{},
						},
					},
					Target: "192.168.0.3",
				},
				Val: &pb.TypedValue{
					Value: &pb.TypedValue_UintVal{
						UintVal: uint64(4),
					},
				},
			},
			{
				Path: &pb.Path{
					Elem: []*pb.PathElem{
						{
							Name: "interfaces",
							Key:  map[string]string{"namespace": "urn:ietf:params:xml:ns:yang:ietf-interfaces"},
						},
						{
							Name: "interface",
							Key:  map[string]string{"name": "sw0p1"},
						},
						{
							Name: "gate-parameters",
							Key:  map[string]string{"namespace": "urn:ieee:std:802.1Q:yang:ieee802-dot1q-sched"},
						},
						{
							Name: "admin-control-list",
							Key:  map[string]string{"index": fmt.Sprint(0)},
						},
						{
							Name: "operation-name",
							Key:  map[string]string{},
						},
					},
					Target: "192.168.0.3",
				},
				Val: &pb.TypedValue{
					Value: &pb.TypedValue_StringVal{
						StringVal: "set-gate-states",
					},
				},
			},
			{
				Path: &pb.Path{
					Elem: []*pb.PathElem{
						{
							Name: "interfaces",
							Key:  map[string]string{"namespace": "urn:ietf:params:xml:ns:yang:ietf-interfaces"},
						},
						{
							Name: "interface",
							Key:  map[string]string{"name": "sw0p1"},
						},
						{
							Name: "gate-parameters",
							Key:  map[string]string{"namespace": "urn:ieee:std:802.1Q:yang:ieee802-dot1q-sched"},
						},
						{
							Name: "admin-control-list",
							Key:  map[string]string{"index": fmt.Sprint(0)},
						},
						{
							Name: "sgs-params",
							Key:  map[string]string{},
						},
						{
							Name: "gate-states-value",
							Key:  map[string]string{},
						},
					},
					Target: "192.168.0.3",
				},
				Val: &pb.TypedValue{
					Value: &pb.TypedValue_UintVal{
						UintVal: 19,
					},
				},
			},
			{
				Path: &pb.Path{
					Elem: []*pb.PathElem{
						{
							Name: "interfaces",
							Key:  map[string]string{"namespace": "urn:ietf:params:xml:ns:yang:ietf-interfaces"},
						},
						{
							Name: "interface",
							Key:  map[string]string{"name": "sw0p1"},
						},
						{
							Name: "gate-parameters",
							Key:  map[string]string{"namespace": "urn:ieee:std:802.1Q:yang:ieee802-dot1q-sched"},
						},
						{
							Name: "admin-control-list",
							Key:  map[string]string{"index": fmt.Sprint(0)},
						},
						{
							Name: "sgs-params",
							Key:  map[string]string{},
						},
						{
							Name: "time-interval-value",
							Key:  map[string]string{},
						},
					},
					Target: "192.168.0.3",
				},
				Val: &pb.TypedValue{
					Value: &pb.TypedValue_UintVal{
						UintVal: 250000,
					},
				},
			},
			{
				Path: &pb.Path{
					Elem: []*pb.PathElem{
						{
							Name: "interfaces",
							Key:  map[string]string{"namespace": "urn:ietf:params:xml:ns:yang:ietf-interfaces"},
						},
						{
							Name: "interface",
							Key:  map[string]string{"name": "sw0p1"},
						},
						{
							Name: "gate-parameters",
							Key:  map[string]string{"namespace": "urn:ieee:std:802.1Q:yang:ieee802-dot1q-sched"},
						},
						{
							Name: "admin-control-list",
							Key:  map[string]string{"index": fmt.Sprint(1)},
						},
						{
							Name: "operation-name",
							Key:  map[string]string{},
						},
					},
					Target: "192.168.0.3",
				},
				Val: &pb.TypedValue{
					Value: &pb.TypedValue_StringVal{
						StringVal: "set-gate-states",
					},
				},
			},
			{
				Path: &pb.Path{
					Elem: []*pb.PathElem{
						{
							Name: "interfaces",
							Key:  map[string]string{"namespace": "urn:ietf:params:xml:ns:yang:ietf-interfaces"},
						},
						{
							Name: "interface",
							Key:  map[string]string{"name": "sw0p1"},
						},
						{
							Name: "gate-parameters",
							Key:  map[string]string{"namespace": "urn:ieee:std:802.1Q:yang:ieee802-dot1q-sched"},
						},
						{
							Name: "admin-control-list",
							Key:  map[string]string{"index": fmt.Sprint(1)},
						},
						{
							Name: "sgs-params",
							Key:  map[string]string{},
						},
						{
							Name: "gate-states-value",
							Key:  map[string]string{},
						},
					},
					Target: "192.168.0.3",
				},
				Val: &pb.TypedValue{
					Value: &pb.TypedValue_UintVal{
						UintVal: 35,
					},
				},
			},
			{
				Path: &pb.Path{
					Elem: []*pb.PathElem{
						{
							Name: "interfaces",
							Key:  map[string]string{"namespace": "urn:ietf:params:xml:ns:yang:ietf-interfaces"},
						},
						{
							Name: "interface",
							Key:  map[string]string{"name": "sw0p1"},
						},
						{
							Name: "gate-parameters",
							Key:  map[string]string{"namespace": "urn:ieee:std:802.1Q:yang:ieee802-dot1q-sched"},
						},
						{
							Name: "admin-control-list",
							Key:  map[string]string{"index": fmt.Sprint(1)},
						},
						{
							Name: "sgs-params",
							Key:  map[string]string{},
						},
						{
							Name: "time-interval-value",
							Key:  map[string]string{},
						},
					},
					Target: "192.168.0.3",
				},
				Val: &pb.TypedValue{
					Value: &pb.TypedValue_UintVal{
						UintVal: 250000,
					},
				},
			},
			{
				Path: &pb.Path{
					Elem: []*pb.PathElem{
						{
							Name: "interfaces",
							Key:  map[string]string{"namespace": "urn:ietf:params:xml:ns:yang:ietf-interfaces"},
						},
						{
							Name: "interface",
							Key:  map[string]string{"name": "sw0p1"},
						},
						{
							Name: "gate-parameters",
							Key:  map[string]string{"namespace": "urn:ieee:std:802.1Q:yang:ieee802-dot1q-sched"},
						},
						{
							Name: "admin-control-list",
							Key:  map[string]string{"index": fmt.Sprint(2)},
						},
						{
							Name: "operation-name",
							Key:  map[string]string{},
						},
					},
					Target: "192.168.0.3",
				},
				Val: &pb.TypedValue{
					Value: &pb.TypedValue_StringVal{
						StringVal: "set-gate-states",
					},
				},
			},
			{
				Path: &pb.Path{
					Elem: []*pb.PathElem{
						{
							Name: "interfaces",
							Key:  map[string]string{"namespace": "urn:ietf:params:xml:ns:yang:ietf-interfaces"},
						},
						{
							Name: "interface",
							Key:  map[string]string{"name": "sw0p1"},
						},
						{
							Name: "gate-parameters",
							Key:  map[string]string{"namespace": "urn:ieee:std:802.1Q:yang:ieee802-dot1q-sched"},
						},
						{
							Name: "admin-control-list",
							Key:  map[string]string{"index": fmt.Sprint(2)},
						},
						{
							Name: "sgs-params",
							Key:  map[string]string{},
						},
						{
							Name: "gate-states-value",
							Key:  map[string]string{},
						},
					},
					Target: "192.168.0.3",
				},
				Val: &pb.TypedValue{
					Value: &pb.TypedValue_UintVal{
						UintVal: 67,
					},
				},
			},
			{
				Path: &pb.Path{
					Elem: []*pb.PathElem{
						{
							Name: "interfaces",
							Key:  map[string]string{"namespace": "urn:ietf:params:xml:ns:yang:ietf-interfaces"},
						},
						{
							Name: "interface",
							Key:  map[string]string{"name": "sw0p1"},
						},
						{
							Name: "gate-parameters",
							Key:  map[string]string{"namespace": "urn:ieee:std:802.1Q:yang:ieee802-dot1q-sched"},
						},
						{
							Name: "admin-control-list",
							Key:  map[string]string{"index": fmt.Sprint(2)},
						},
						{
							Name: "sgs-params",
							Key:  map[string]string{},
						},
						{
							Name: "time-interval-value",
							Key:  map[string]string{},
						},
					},
					Target: "192.168.0.3",
				},
				Val: &pb.TypedValue{
					Value: &pb.TypedValue_UintVal{
						UintVal: 250000,
					},
				},
			},
			{
				Path: &pb.Path{
					Elem: []*pb.PathElem{
						{
							Name: "interfaces",
							Key:  map[string]string{"namespace": "urn:ietf:params:xml:ns:yang:ietf-interfaces"},
						},
						{
							Name: "interface",
							Key:  map[string]string{"name": "sw0p1"},
						},
						{
							Name: "gate-parameters",
							Key:  map[string]string{"namespace": "urn:ieee:std:802.1Q:yang:ieee802-dot1q-sched"},
						},
						{
							Name: "admin-control-list",
							Key:  map[string]string{"index": fmt.Sprint(3)},
						},
						{
							Name: "operation-name",
							Key:  map[string]string{},
						},
					},
					Target: "192.168.0.3",
				},
				Val: &pb.TypedValue{
					Value: &pb.TypedValue_StringVal{
						StringVal: "set-gate-states",
					},
				},
			},
			{
				Path: &pb.Path{
					Elem: []*pb.PathElem{
						{
							Name: "interfaces",
							Key:  map[string]string{"namespace": "urn:ietf:params:xml:ns:yang:ietf-interfaces"},
						},
						{
							Name: "interface",
							Key:  map[string]string{"name": "sw0p1"},
						},
						{
							Name: "gate-parameters",
							Key:  map[string]string{"namespace": "urn:ieee:std:802.1Q:yang:ieee802-dot1q-sched"},
						},
						{
							Name: "admin-control-list",
							Key:  map[string]string{"index": fmt.Sprint(3)},
						},
						{
							Name: "sgs-params",
							Key:  map[string]string{},
						},
						{
							Name: "gate-states-value",
							Key:  map[string]string{},
						},
					},
					Target: "192.168.0.3",
				},
				Val: &pb.TypedValue{
					Value: &pb.TypedValue_UintVal{
						UintVal: 131,
					},
				},
			},
			{
				Path: &pb.Path{
					Elem: []*pb.PathElem{
						{
							Name: "interfaces",
							Key:  map[string]string{"namespace": "urn:ietf:params:xml:ns:yang:ietf-interfaces"},
						},
						{
							Name: "interface",
							Key:  map[string]string{"name": "sw0p1"},
						},
						{
							Name: "gate-parameters",
							Key:  map[string]string{"namespace": "urn:ieee:std:802.1Q:yang:ieee802-dot1q-sched"},
						},
						{
							Name: "admin-control-list",
							Key:  map[string]string{"index": fmt.Sprint(3)},
						},
						{
							Name: "sgs-params",
							Key:  map[string]string{},
						},
						{
							Name: "time-interval-value",
							Key:  map[string]string{},
						},
					},
					Target: "192.168.0.3",
				},
				Val: &pb.TypedValue{
					Value: &pb.TypedValue_UintVal{
						UintVal: 240000,
					},
				},
			},
			{
				Path: &pb.Path{
					Elem: []*pb.PathElem{
						{
							Name: "interfaces",
							Key:  map[string]string{"namespace": "urn:ietf:params:xml:ns:yang:ietf-interfaces"},
						},
						{
							Name: "interface",
							Key:  map[string]string{"name": "sw0p1"},
						},
						{
							Name: "gate-parameters",
							Key:  map[string]string{"namespace": "urn:ieee:std:802.1Q:yang:ieee802-dot1q-sched"},
						},
						{
							Name: "admin-cycle-time",
							Key:  map[string]string{},
						},
						{
							Name: "numerator",
							Key:  map[string]string{},
						},
					},
					Target: "192.168.0.3",
				},
				Val: &pb.TypedValue{
					Value: &pb.TypedValue_IntVal{
						IntVal: int64(1),
					},
				},
			},
			{
				Path: &pb.Path{
					Elem: []*pb.PathElem{
						{
							Name: "interfaces",
							Key:  map[string]string{"namespace": "urn:ietf:params:xml:ns:yang:ietf-interfaces"},
						},
						{
							Name: "interface",
							Key:  map[string]string{"name": "sw0p1"},
						},
						{
							Name: "gate-parameters",
							Key:  map[string]string{"namespace": "urn:ieee:std:802.1Q:yang:ieee802-dot1q-sched"},
						},
						{
							Name: "admin-cycle-time",
							Key:  map[string]string{},
						},
						{
							Name: "denominator",
							Key:  map[string]string{},
						},
					},
					Target: "192.168.0.3",
				},
				Val: &pb.TypedValue{
					Value: &pb.TypedValue_IntVal{
						IntVal: 1000,
					},
				},
			},
			{
				Path: &pb.Path{
					Elem: []*pb.PathElem{
						{
							Name: "interfaces",
							Key:  map[string]string{"namespace": "urn:ietf:params:xml:ns:yang:ietf-interfaces"},
						},
						{
							Name: "interface",
							Key:  map[string]string{"name": "sw0p1"},
						},
						{
							Name: "gate-parameters",
							Key:  map[string]string{"namespace": "urn:ieee:std:802.1Q:yang:ieee802-dot1q-sched"},
						},
						{
							Name: "admin-cycle-time-extension",
							Key:  map[string]string{},
						},
					},
					Target: "192.168.0.3",
				},
				Val: &pb.TypedValue{
					Value: &pb.TypedValue_UintVal{
						UintVal: 0,
					},
				},
			},
			{
				Path: &pb.Path{
					Elem: []*pb.PathElem{
						{
							Name: "interfaces",
							Key:  map[string]string{"namespace": "urn:ietf:params:xml:ns:yang:ietf-interfaces"},
						},
						{
							Name: "interface",
							Key:  map[string]string{"name": "sw0p1"},
						},
						{
							Name: "gate-parameters",
							Key:  map[string]string{"namespace": "urn:ieee:std:802.1Q:yang:ieee802-dot1q-sched"},
						},
						{
							Name: "admin-base-time",
							Key:  map[string]string{},
						},
						{
							Name: "seconds",
							Key:  map[string]string{},
						},
					},
					Target: "192.168.0.3",
				},
				Val: &pb.TypedValue{
					Value: &pb.TypedValue_StringVal{
						StringVal: "0",
					},
				},
			},
			{
				Path: &pb.Path{
					Elem: []*pb.PathElem{
						{
							Name: "interfaces",
							Key:  map[string]string{"namespace": "urn:ietf:params:xml:ns:yang:ietf-interfaces"},
						},
						{
							Name: "interface",
							Key:  map[string]string{"name": "sw0p1"},
						},
						{
							Name: "gate-parameters",
							Key:  map[string]string{"namespace": "urn:ieee:std:802.1Q:yang:ieee802-dot1q-sched"},
						},
						{
							Name: "admin-base-time",
							Key:  map[string]string{},
						},
						{
							Name: "fractional-seconds",
							Key:  map[string]string{},
						},
					},
					Target: "192.168.0.3",
				},
				Val: &pb.TypedValue{
					Value: &pb.TypedValue_StringVal{
						StringVal: "0",
					},
				},
			},
			{
				Path: &pb.Path{
					Elem: []*pb.PathElem{
						{
							Name: "interfaces",
							Key:  map[string]string{"namespace": "urn:ietf:params:xml:ns:yang:ietf-interfaces"},
						},
						{
							Name: "interface",
							Key:  map[string]string{"name": "sw0p1"},
						},
						{
							Name: "gate-parameters",
							Key:  map[string]string{"namespace": "urn:ieee:std:802.1Q:yang:ieee802-dot1q-sched"},
						},
						{
							Name: "config-change",
							Key:  map[string]string{},
						},
					},
					Target: "192.168.0.3",
				},
				Val: &pb.TypedValue{
					Value: &pb.TypedValue_BoolVal{
						BoolVal: true,
					},
				},
			},
		},
	}

	response, err := c.(*gclient.Client).Set(ctx, &setRequest)
	if err != nil {
		fmt.Print("Target returned RPC error for Set: ")
		fmt.Println(err)
		return
	}

	fmt.Println(response)
}

// func getNewTreeStructure(schemaEntries []*adapterResp.SchemaEntry) *SchemaTree {
// 	var newTree *SchemaTree
// 	tree := &SchemaTree{}
// 	lastNode := ""
// 	for _, entry := range schemaEntries {
// 		if entry.Value == "" {
// 			// In a directory
// 			if entry.Tag == "end" {
// 				if entry.Name != "data" {
// 					if lastNode != "leaf" {
// 						tree = tree.Parent
// 					}
// 					lastNode = ""
// 				}
// 			} else {
// 				newTree = &SchemaTree{Parent: tree}

// 				newTree.Name = entry.Name
// 				newTree.Namespace = entry.Namespace
// 				newTree.Parent.Children = append(newTree.Parent.Children, newTree)

// 				tree = newTree
// 			}
// 		} else {
// 			// In a leaf
// 			newTree = &SchemaTree{Parent: tree}

// 			newTree.Name = entry.Name
// 			newTree.Value = entry.Value
// 			newTree.Parent.Children = append(newTree.Parent.Children, newTree)

// 			lastNode = "leaf"
// 		}
// 	}
// 	return tree
// }

// func printTree(tree *SchemaTree, tabLevels int) {
// 	tabs := ""
// 	for i := 0; i < tabLevels; i++ {
// 		tabs += "  "
// 	}

// 	fmt.Println(tabs + tree.Name + "---" + tree.Namespace + "---" + tree.Value)
// 	for _, child := range tree.Children {
// 		printTree(child, tabLevels+1)
// 	}
// }

// func getFullConfigFromSwitch(addr string) {
// 	sshConfig := &ssh.ClientConfig{
// 		User:            "root",
// 		Auth:            []ssh.AuthMethod{ssh.Password("")},
// 		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
// 	}

// 	s, err := netconf.DialSSH(addr, sshConfig)
// 	if err != nil {
// 		fmt.Printf("Failed creating connection: %v\n", err)
// 	}

// 	defer s.Close()

// 	reply, err := s.Exec(netconf.MethodGetConfig("running"))
// 	if err != nil {
// 		fmt.Printf("Failed getting config: %v\n", err)
// 		return
// 	}

// 	fmt.Println(reply.Data)
// }

// func testNetworkChangeRequest(switchAddr string) {
// 	ctx := context.Background()

// 	address := []string{"gnmi-netconf-adapter:11161"}

// 	c, err := gclient.New(ctx, client.Destination{
// 		Addrs:       address,
// 		Timeout:     time.Second * 5,
// 		Credentials: nil,
// 		TLS:         nil,
// 	})

// 	if err != nil {
// 		fmt.Print("Could not create a gNMI client: ")
// 		fmt.Println(err)
// 	}

// 	// myMap := map[string]string{}
// 	// myMap["Hello"] = "World!"

// 	setRequest := pb.SetRequest{
// 		Update: []*pb.Update{
// 			{
// 				Path: &pb.Path{
// 					Target: switchAddr,
// 					Elem: []*pb.PathElem{
// 						{
// 							Name: "interfaces",
// 							Key:  map[string]string{"namespace": "urn:ietf:params:xml:ns:yang:ietf-interfaces"},
// 						},
// 						{
// 							Name: "interface",
// 							Key:  map[string]string{"name": "sw0p1"},
// 						},
// 						{
// 							Name: "max-sdu-table",
// 							Key:  map[string]string{"namespace": "urn:ieee:std:802.1Q:yang:ieee802-dot1q-sched", "traffic-class": "0"},
// 						},
// 						{
// 							Name: "queue-max-sdu",
// 						},
// 					}, // Path to an element that should be updated
// 				},
// 				Val: &pb.TypedValue{
// 					Value: &pb.TypedValue_StringVal{
// 						StringVal: "1503",
// 					},
// 				},
// 			},
// 			// {
// 			// 	Path: &pb.Path{
// 			// 		Target: "192.168.0.1",
// 			// 		Elem: []*pb.PathElem{
// 			// 			{
// 			// 				Name: "interfaces",
// 			// 				Key:  map[string]string{"namespace": "urn:ietf:params:xml:ns:yang:ietf-interfaces"},
// 			// 			},
// 			// 			{
// 			// 				Name: "interface",
// 			// 				Key:  map[string]string{"name": "sw0p1"},
// 			// 			},
// 			// 			{
// 			// 				Name: "max-sdu-table",
// 			// 				Key:  map[string]string{"namespace": "urn:ieee:std:802.1Q:yang:ieee802-dot1q-sched", "traffic-class": "1"},
// 			// 			},
// 			// 			{
// 			// 				Name: "queue-max-sdu",
// 			// 			},
// 			// 		}, // Path to an element that should be updated
// 			// 	},
// 			// 	Val: &pb.TypedValue{
// 			// 		Value: &pb.TypedValue_StringVal{
// 			// 			StringVal: "1505",
// 			// 		},
// 			// 	},
// 			// },
// 			// {
// 			// 	Path: &pb.Path{
// 			// 		Target: "192.168.0.2",
// 			// 	},
// 			// },
// 			// {
// 			// 	Path: &pb.Path{
// 			// 		Target: "192.168.0.2",
// 			// 	},
// 			// },
// 		},
// 		Extension: []*gnmi_ext.Extension{
// 			{
// 				Ext: &gnmi_ext.Extension_RegisteredExt{
// 					RegisteredExt: &gnmi_ext.RegisteredExtension{
// 						Id:  gnmi_ext.ExtensionID(100),
// 						Msg: []byte("my_network_change"),
// 					},
// 				},
// 			},
// 			{
// 				Ext: &gnmi_ext.Extension_RegisteredExt{
// 					RegisteredExt: &gnmi_ext.RegisteredExtension{
// 						Id:  gnmi_ext.ExtensionID(101),
// 						Msg: []byte("1.0.2"),
// 					},
// 				},
// 			},
// 			{
// 				Ext: &gnmi_ext.Extension_RegisteredExt{
// 					RegisteredExt: &gnmi_ext.RegisteredExtension{
// 						Id:  gnmi_ext.ExtensionID(102),
// 						Msg: []byte("tsn-model"),
// 					},
// 				},
// 			},
// 		},
// 	}

// 	response, err := c.(*gclient.Client).Set(ctx, &setRequest)
// 	if err != nil {
// 		fmt.Print("Target returned RPC error for Set: ")
// 		fmt.Println(err)
// 		return
// 	}

// 	fmt.Print("Response from gnmi-netconf-adapter is: ")
// 	fmt.Println(response)
// }

// func testAtomixStore() {
// 	ctx := context.Background()

// 	fmt.Println("Creating client")
// 	atomixClient := atomix.NewClient(atomix.WithClientID(os.Getenv("POD_NAME")))

// 	fmt.Println("Getting map")

// 	myMap, err := atomixClient.GetMap(ctx, "monitor-config")

// 	if err != nil {
// 		fmt.Printf("Error from atomixClient.GetMap:%+v\n", err)
// 		return
// 	}

// 	fmt.Println("Pushing value to myMap")

// 	newVal, err := myMap.Put(ctx, "Test", []byte("This works now"))
// 	if err != nil {
// 		fmt.Printf("Error pushing new entry to myMap: %v\n", err)
// 		return
// 	}

// 	fmt.Printf("myMap now contains entry: %v\n", newVal)

// 	myMap.Close(ctx)
// }

// func testSwitchDelay() {
// 	sshConfig := &ssh.ClientConfig{
// 		User:            "root",
// 		Auth:            []ssh.AuthMethod{ssh.Password("")},
// 		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
// 	}

// 	// request := netconf.RawMethod("<get xmlns='urn:ietf:params:xml:ns:netconf:base:1.0'><filter type='subtree'><interfaces xmlns='urn:ietf:params:xml:ns:yang:ietf-interfaces'><interface><name>sw0p1</name><ethernet xmlns='urn:ieee:std:802.3:yang:ieee802-ethernet-interface'><statistics><frame><in-total-frames></in-total-frames></frame></statistics></ethernet></interface></interfaces></filter></get>")

// 	s, err := netconf.DialSSH("192.168.0.1", sshConfig)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	defer s.Close()

// 	// fmt.Println(s.ServerCapabilities)
// 	// fmt.Println(s.SessionID)

// 	hello := &netconf.HelloMessage{Capabilities: netconf.DefaultCapabilities}

// 	val, err := xml.Marshal(hello)
// 	if err != nil {
// 		fmt.Printf("Failed to marshal hello message: %v\n", err)
// 	}

// 	header := []byte(xml.Header)
// 	val = append(header, val...)

// 	fmt.Println(val)

// 	var req string
// 	err = xml.Unmarshal(val, &req)
// 	if err != nil {
// 		fmt.Printf("Failed to unmarshal value: %v\n", err)
// 	}

// 	fmt.Println(req)

// 	// start := time.Now().UnixNano()
// 	// _, err = s.Exec(request)
// 	// end := time.Now().UnixNano()

// 	// fmt.Printf("Delay: %v\n", end-start)

// 	// if err != nil {
// 	// 	panic(err)
// 	// }

// 	fmt.Println("Done!")
// }

// // func testNetconfClient() {
// // 	fmt.Println("Creating session...")
// // 	session := createSession()
// // 	fmt.Println("Session created!")
// // 	defer session.Close()
// // 	execRPC(session)
// // }

// // func execRPC(session *netconf.Session) {
// // 	// gt := message.NewGet("", "")
// // 	// gt := message.NewGet("", "<interfaces xmlns=\"urn:ietf:params:xml:ns:yang:ietf-interfaces\"><interface><name>sw0p1</name><ethernet xmlns=\"urn:ieee:std:802.3:yang:ieee802-ethernet-interface\"><statistics><frame><in-total-frames/></frame></statistics></ethernet></interface></interfaces>")
// // 	gt := message.NewRPC(`<get><filter type="subtree"><interfaces xmlns="urn:ietf:params:xml:ns:yang:ietf-interfaces"><interface><name>sw0p1</name><ethernet xmlns="urn:ieee:std:802.3:yang:ieee802-ethernet-interface"><statistics><frame><in-total-frames></in-total-frames></frame></statistics></ethernet></interface></interfaces></filter></get>`)
// // 	fmt.Println("Message created!")
// // 	start := time.Now().UnixNano()
// // 	session.AsyncRPC(gt, defaultLogRpcReplyCallback(gt.MessageID, start))
// // 	time.Sleep(100 * time.Millisecond)
// //
// // 	fmt.Printf("MessageID: %v\n", gt.MessageID)
// // 	// fmt.Printf("delay: %v\n", time.Now().UnixNano()-start)
// // 	// if err != nil {
// // 	// 	fmt.Printf("Failed RPC: %v\n", err)
// // 	// } else {
// // 	// 	fmt.Println(reply.RawReply)
// // 	// }
// //
// // 	// d2 := message.NewCloseSession()
// // 	// start2 := time.Now().UnixNano()
// // 	// session.AsyncRPC(d2, defaultLogRpcReplyCallback(d2.MessageID, start2))
// //
// // 	session.Listener.WaitForMessages()
// // }

// // func createSession() *netconf.Session {
// // 	sshConfig := &ssh.ClientConfig{
// // 		User:            "root",
// // 		Auth:            []ssh.AuthMethod{ssh.Password("")},
// // 		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
// // 	}
// // 	s, err := netconf.DialSSH("192.168.0.1:830", sshConfig)
// // 	if err != nil {
// // 		log.Fatal(err)
// // 	}
// //
// // 	capabilities := netconf.DefaultCapabilities
// // 	err = s.SendHello(&message.Hello{Capabilities: capabilities})
// // 	if err != nil {
// // 		log.Fatal(err)
// // 	}
// //
// // 	return s
// // }

// // func defaultLogRpcReplyCallback(eventId string, start int64) netconf.Callback {
// // 	return func(event netconf.Event) {
// // 		reply := event.RPCReply()
// // 		fmt.Printf("delay for event %v: %v\n", eventId, time.Now().UnixNano()-start)
// // 		if reply == nil {
// // 			println("Failed to execute RPC")
// // 		}
// // 		if event.EventID() == eventId {
// // 			println("Successfully executed RPC")
// // 			println(reply.RawReply)
// // 		}
// // 	}
// // }

// func testSequences() {
// 	// fmt.Println("Start batch monitoring on switch_one, switch_two, and switch_three")
// 	// setReq("Start", "192.168.0.1", "0")
// 	// setReq("Start", "192.168.0.2", "0")
// 	// setReq("Start", "192.168.0.3", "0")
// 	// time.Sleep(2 * time.Minute)
// 	// setReq("Stop", "192.168.0.1")
// 	// setReq("Stop", "192.168.0.2")
// 	// setReq("Stop", "192.168.0.3")

// 	// time.Sleep(1 * time.Minute)

// 	// fmt.Println("Start batch monitoring on switch_one, switch_two, and switch_three")
// 	// setReq("Start", "192.168.0.1", "0")
// 	// setReq("Start", "192.168.0.2", "0")
// 	// setReq("Start", "192.168.0.3", "0")
// 	// time.Sleep(1 * time.Minute)
// 	// setReq("Stop", "192.168.0.1")
// 	// setReq("Stop", "192.168.0.2")
// 	// setReq("Stop", "192.168.0.3")

// 	// time.Sleep(1 * time.Minute)

// 	// fmt.Println("Start update test")
// 	// setReq("Start", "192.168.0.1", "0")
// 	// setReq("Start", "192.168.0.2", "0")
// 	// setReq("Start", "192.168.0.3", "0")

// 	// time.Sleep(30 * time.Second)

// 	// fmt.Println("Starting to update now")

// 	// setReq("Update", "192.168.0.1", "1")
// 	// setReq("Update", "192.168.0.2", "1")
// 	// setReq("Update", "192.168.0.3", "1")

// 	// time.Sleep(30 * time.Second)

// 	// fmt.Println("Updating second time")

// 	// setReq("Update", "192.168.0.1", "0")
// 	// setReq("Update", "192.168.0.2", "0")
// 	// setReq("Update", "192.168.0.3", "0")

// 	// time.Sleep(60 * time.Second)

// 	// fmt.Println("Updating third time")

// 	// setReq("Update", "192.168.0.1", "1")
// 	// setReq("Update", "192.168.0.2", "1")
// 	// setReq("Update", "192.168.0.3", "1")

// 	// time.Sleep(30 * time.Second)

// 	// fmt.Println("Updating fourth time")

// 	// setReq("Update", "192.168.0.1", "0")
// 	// setReq("Update", "192.168.0.2", "0")
// 	// setReq("Update", "192.168.0.3", "0")

// 	// time.Sleep(60 * time.Second)

// 	// fmt.Println("Updating fifth and final time")

// 	// setReq("Update", "192.168.0.1", "1")
// 	// setReq("Update", "192.168.0.2", "1")
// 	// setReq("Update", "192.168.0.3", "1")

// 	// fmt.Println("Start non-batch monitoring on switch_one, switch_two, and switch_three")
// 	// setReq("Start", "192.168.0.1", "1")
// 	// setReq("Start", "192.168.0.2", "1")
// 	// setReq("Start", "192.168.0.3", "1")
// 	// time.Sleep(2 * time.Minute)
// 	// setReq("Stop", "192.168.0.1")
// 	// setReq("Stop", "192.168.0.2")
// 	// setReq("Stop", "192.168.0.3")

// 	setReq("Start", "192.168.0.1", "0")
// 	time.Sleep(30 * time.Second)
// 	setReq("Stop", "192.168.0.1")
// }

// func getFullConfig(switchAddr string) *pb.GetResponse {
// 	ctx := context.Background()

// 	address := []string{"gnmi-netconf-adapter:11161"}

// 	c, err := gclient.New(ctx, client.Destination{
// 		Addrs:       address,
// 		Target:      "gnmi-netconf-adapter",
// 		Timeout:     time.Second * 5,
// 		Credentials: nil,
// 		TLS:         nil,
// 	})

// 	if err != nil {
// 		// fmt.Errorf("could not create a gNMI client: %v", err)
// 		fmt.Print("Could not create a gNMI client: ")
// 		fmt.Println(err)
// 	}

// 	getRequest := pb.GetRequest{
// 		Path: []*pb.Path{
// 			{
// 				Target: switchAddr,
// 			},
// 		},
// 		Type: pb.GetRequest_STATE,
// 	}

// 	response, err := c.(*gclient.Client).Get(ctx, &getRequest)
// 	if err != nil {
// 		fmt.Print("Target returned RPC error for Testing: ")
// 		fmt.Println(err)
// 	}

// 	c.Close()

// 	// fmt.Println(response)

// 	return response
// }

// func testing() {
// 	ctx := context.Background()

// 	// Send get request to adapter

// 	address := []string{"gnmi-netconf-adapter:11161"}

// 	c, err := gclient.New(ctx, client.Destination{
// 		Addrs:       address,
// 		Target:      "gnmi-netconf-adapter",
// 		Timeout:     time.Second * 5,
// 		Credentials: nil,
// 		TLS:         nil,
// 	})

// 	if err != nil {
// 		// fmt.Errorf("could not create a gNMI client: %v", err)
// 		fmt.Print("Could not create a gNMI client: ")
// 		fmt.Println(err)
// 	}

// 	// interfaceKeyMap := map[string]string{}
// 	// interfaceKeyMap["namespace"] = "urn:ietf:params:xml:ns:yang:ietf-interfaces"

// 	// loopbackKeyMap := map[string]string{}
// 	// loopbackKeyMap["name"] = "lo"

// 	getRequest := pb.GetRequest{
// 		Path: []*pb.Path{
// 			{
// 				// Elem: []*pb.PathElem{
// 				// 	{
// 				// 		Name: "interfaces",
// 				// 		Key:  interfaceKeyMap,
// 				// 	},
// 				// 	{
// 				// 		Name: "interface",
// 				// 		Key:  loopbackKeyMap,
// 				// 	},
// 				// },
// 				Target: "192.168.0.2",
// 			},
// 		},
// 		Type: pb.GetRequest_STATE,
// 	}

// 	response, err := c.(*gclient.Client).Get(ctx, &getRequest)
// 	if err != nil {
// 		fmt.Print("Target returned RPC error for Testing: ")
// 		fmt.Println(err)
// 	}

// 	c.Close()

// 	fmt.Println(response)

// 	// Send set request to storage

// 	address = []string{"storage-service:11161"}

// 	c, err = gclient.New(ctx, client.Destination{
// 		Addrs:       address,
// 		Target:      "storage-service",
// 		Timeout:     time.Second * 5,
// 		Credentials: nil,
// 		TLS:         nil,
// 	})

// 	if err != nil {
// 		// fmt.Errorf("could not create a gNMI client: %v", err)
// 		fmt.Print("Could not create a gNMI client: ")
// 		fmt.Println(err)
// 	}

// 	//

// 	if len(response.Notification) <= 0 {
// 		return
// 	}

// 	actionMap := map[string]string{}
// 	actionMap["Action"] = "Store namespaces"

// 	setRequest := pb.SetRequest{
// 		Update: []*pb.Update{
// 			{
// 				Path: &pb.Path{
// 					Elem: []*pb.PathElem{
// 						{
// 							Name: "Action",
// 							Key:  actionMap,
// 						},
// 					},
// 				},
// 				Val: response.Notification[0].Update[0].Val,
// 			},
// 		},
// 	}

// 	setResponse, err := c.(*gclient.Client).Set(ctx, &setRequest)
// 	if err != nil {
// 		fmt.Print("Target returned RPC error for Testing: ")
// 		fmt.Println(err)
// 	}

// 	c.Close()

// 	fmt.Println(setResponse)

// 	if err != nil {
// 		// fmt.Errorf("could not create a gNMI client: %v", err)
// 		fmt.Print("Could not create a gNMI client: ")
// 		fmt.Println(err)
// 	}

// 	// Send get request to storage

// 	address = []string{"storage-service:11161"}

// 	c, err = gclient.New(ctx, client.Destination{
// 		Addrs:       address,
// 		Target:      "storage-service",
// 		Timeout:     time.Second * 5,
// 		Credentials: nil,
// 		TLS:         nil,
// 	})

// 	if err != nil {
// 		// fmt.Errorf("could not create a gNMI client: %v", err)
// 		fmt.Print("Could not create a gNMI client: ")
// 		fmt.Println(err)
// 	}

// 	getRequest = pb.GetRequest{
// 		Path: []*pb.Path{
// 			{
// 				Elem: []*pb.PathElem{
// 					{
// 						Name: "interfaces",
// 					},
// 					{
// 						Name: "interface",
// 						Key: map[string]string{
// 							"name": "sw0p1",
// 						},
// 					},
// 					{
// 						Name: "bridge-port",
// 					},
// 					{
// 						Name: "traffic-class",
// 					},
// 					{
// 						Name: "priority0",
// 					},
// 				},
// 			},
// 		},
// 		Type: 4,
// 	}

// 	response, err = c.(*gclient.Client).Get(ctx, &getRequest)
// 	if err != nil {
// 		fmt.Print("Target returned RPC error for Testing: ")
// 		fmt.Println(err)
// 	}

// 	c.Close()

// 	fmt.Println(response.Notification[0].Update[0])

// 	// Send get request to adapter

// 	address = []string{"gnmi-netconf-adapter:11161"}

// 	c, err = gclient.New(ctx, client.Destination{
// 		Addrs:       address,
// 		Target:      "gnmi-netconf-adapter",
// 		Timeout:     time.Second * 5,
// 		Credentials: nil,
// 		TLS:         nil,
// 	})

// 	if err != nil {
// 		// fmt.Errorf("could not create a gNMI client: %v", err)
// 		fmt.Print("Could not create a gNMI client: ")
// 		fmt.Println(err)
// 	}

// 	getRequest = pb.GetRequest{
// 		Path: []*pb.Path{
// 			response.Notification[0].Update[0].Path,
// 		},
// 		Type: pb.GetRequest_STATE,
// 	}

// 	response, err = c.(*gclient.Client).Get(ctx, &getRequest)
// 	if err != nil {
// 		fmt.Print("Target returned RPC error for Testing: ")
// 		fmt.Println(err)
// 	}

// 	c.Close()

// 	// fmt.Println(response)

// 	var schema Schema
// 	var schemaTree *SchemaTree
// 	if len(response.Notification) > 0 {
// 		json.Unmarshal(response.Notification[0].Update[0].Val.GetBytesVal(), &schema)

// 		// fmt.Println(schema)
// 		schemaTree = getTreeStructure(schema)
// 	}

// 	printSchemaTree(schemaTree)
// }

// func printSchemaTree(schemaTree *SchemaTree) {
// 	fmt.Printf("%s - %s - %v - %s\n", schemaTree.Parent.Name, schemaTree.Name, schemaTree.Namespace, schemaTree.Value)
// 	for _, child := range schemaTree.Children {
// 		printSchemaTree(child)
// 	}
// }

// // TODO: add pointer that traverse the tree based on tags, use that pointer to
// // get correct parents.
// func getTreeStructure(schema Schema) *SchemaTree {
// 	var newTree *SchemaTree
// 	tree := &SchemaTree{}
// 	lastNode := ""
// 	for _, entry := range schema.Entries {
// 		// fmt.Println("-------------------")
// 		// if index == 0 {
// 		// newTree = &SchemaTree{Parent: tree}
// 		// newTree.Name = entry.Name
// 		// newTree.Namespace = entry.Namespace
// 		// fmt.Println(tree.Name)
// 		// tree = &SchemaTree{Parent: tree}
// 		// continue
// 		// }
// 		if entry.Value == "" { // Directory
// 			if entry.Tag == "end" {
// 				if entry.Name != "data" {
// 					if lastNode != "leaf" {
// 						// fmt.Println(tree.Name)
// 						tree = tree.Parent
// 					}
// 					lastNode = ""
// 					// continue
// 				}
// 			} else {

// 				newTree = &SchemaTree{Parent: tree}

// 				newTree.Name = entry.Name
// 				newTree.Namespace = entry.Namespace
// 				newTree.Parent.Children = append(newTree.Parent.Children, newTree)

// 				tree = newTree

// 				// fmt.Print(tree.Name)
// 				// fmt.Print(", ")
// 				// fmt.Println(tree.Parent.Name)

// 				// tree = &SchemaTree{Parent: tree}
// 			}
// 		} else { // Leaf
// 			newTree = &SchemaTree{Parent: tree}

// 			newTree.Name = entry.Name
// 			newTree.Value = entry.Value
// 			newTree.Parent.Children = append(newTree.Parent.Children, newTree)

// 			// fmt.Print(newTree.Name)
// 			// fmt.Print(", ")
// 			// fmt.Println(newTree.Parent.Name)
// 			// fmt.Println(newTree.Value)

// 			lastNode = "leaf"
// 		}
// 		// fmt.Println("-------------------")
// 		// fmt.Print("name: ")
// 		// fmt.Print(tree.Name)
// 		// if tree.Name != "data" {
// 		// 	fmt.Print(", parent: ")
// 		// 	fmt.Println(tree.Parent.Name)

// 		// 	// fmt.Println("#######")
// 		// 	// for i, child := range tree.Parent.Children {
// 		// 	// 	if i < 10 {
// 		// 	// 		fmt.Print(child.Name)
// 		// 	// 		fmt.Print(", ")
// 		// 	// 	}
// 		// 	// }
// 		// 	// fmt.Println("\n******")
// 		// 	// for j, child := range tree.Children {
// 		// 	// 	if j < 10 {
// 		// 	// 		fmt.Print(child.Name)
// 		// 	// 		fmt.Print(": ")
// 		// 	// 		fmt.Print(child.Value)
// 		// 	// 		fmt.Print(", ")
// 		// 	// 	}
// 		// 	// }
// 		// } else {
// 		// 	// fmt.Println("")
// 		// 	// for _, child := range tree.Children {
// 		// 	// 	fmt.Print(child.Name)
// 		// 	// 	fmt.Print(" | ")
// 		// 	// }
// 		// }
// 		// fmt.Println("")
// 		// fmt.Println(entry)
// 		// fmt.Println(tree.Namespace)
// 		// fmt.Println("###################")
// 	}
// 	return tree
// }

// type SchemaTree struct {
// 	Name      string
// 	Namespace string
// 	Children  []*SchemaTree
// 	Parent    *SchemaTree
// 	Value     string
// }

// type Schema struct {
// 	Entries []SchemaEntry
// }

// type SchemaEntry struct {
// 	Name      string
// 	Tag       string
// 	Namespace string
// 	Value     string
// }

// // func setDelete(action string, target string) {
// // 	ctx := context.Background()
// //
// // 	address := []string{"monitor-service:11161"}
// //
// // 	c, err := gclient.New(ctx, client.Destination{
// // 		Addrs:       address,
// // 		Target:      "monitor-service",
// // 		Timeout:     time.Second * 5,
// // 		Credentials: nil,
// // 		TLS:         nil,
// // 	})
// //
// // 	if err != nil {
// // 		// fmt.Errorf("could not create a gNMI client: %v", err)
// // 		fmt.Print("Could not create a gNMI client: ")
// // 		fmt.Println(err)
// // 	}
// //
// // 	actionMap := make(map[string]string)
// // 	actionMap["Action"] = action
// //
// // 	setRequest := pb.SetRequest{
// // 		Update: []*pb.Update{
// // 			{
// // 				Path: &pb.Path{
// // 					Target: target,
// // 					Elem: []*pb.PathElem{
// // 						{
// // 							Name: "Action",
// // 							Key:  actionMap,
// // 						},
// // 					},
// // 				},
// // 			},
// // 		},
// // 	}
// //
// // 	response, err := c.(*gclient.Client).Set(ctx, &setRequest)
// //
// // 	fmt.Print("Response from device-monitor is: ")
// // 	fmt.Println(response)
// //
// // 	if err != nil {
// // 		fmt.Print("Target returned RPC error for Set: ")
// // 		fmt.Println(err)
// // 	}
// //
// // 	fmt.Println("Client connected successfully")
// // }

// func setReq(action string, target string, confIndex ...string) {
// 	ctx := context.Background()

// 	address := []string{"monitor-service:11161"}

// 	c, err := gclient.New(ctx, client.Destination{
// 		Addrs:       address,
// 		Target:      "monitor-service",
// 		Timeout:     time.Second * 5,
// 		Credentials: nil,
// 		TLS:         nil,
// 	})

// 	if err != nil {
// 		// fmt.Errorf("could not create a gNMI client: %v", err)
// 		fmt.Print("Could not create a gNMI client: ")
// 		fmt.Println(err)
// 	}

// 	setRequest := pb.SetRequest{
// 		Update: []*pb.Update{
// 			{
// 				Path: &pb.Path{
// 					Target: target,
// 					Elem: []*pb.PathElem{
// 						{
// 							Name: "Action",
// 							Key: map[string]string{
// 								"Action": action,
// 							},
// 						},
// 					},
// 				},
// 			},
// 		},
// 	}

// 	if confIndex != nil {
// 		setRequest.Update[0].Path.Elem = append(setRequest.Update[0].Path.Elem, &pb.PathElem{
// 			Name: "ConfigIndex",
// 			Key: map[string]string{
// 				"ConfigIndex": confIndex[0],
// 			},
// 		})
// 	}

// 	response, err := c.(*gclient.Client).Set(ctx, &setRequest)

// 	fmt.Print("Response from device-monitor is: ")
// 	fmt.Println(response)

// 	if err != nil {
// 		fmt.Print("Target returned RPC error for Set: ")
// 		fmt.Println(err)
// 	}
// }

// func setUpdate(config types.ConfigRequest) {
// 	ctx := context.Background()

// 	address := []string{"monitor-service:11161"}

// 	c, err := gclient.New(ctx, client.Destination{
// 		Addrs:       address,
// 		Target:      "monitor-service",
// 		Timeout:     time.Second * 5,
// 		Credentials: nil,
// 		TLS:         nil,
// 	})

// 	if err != nil {
// 		// fmt.Errorf("could not create a gNMI client: %v", err)
// 		fmt.Print("Could not create a gNMI client: ")
// 		fmt.Println(err)
// 	}

// 	pathElements := []*pb.PathElem{
// 		{
// 			Name: "Action",
// 			Key: map[string]string{
// 				"Action": "Change config",
// 			},
// 		},
// 		{
// 			Name: "Info",
// 			Key: map[string]string{
// 				"DeviceIP":   config.DeviceIP,
// 				"DeviceName": config.DeviceName,
// 				"Protocol":   config.Protocol,
// 			},
// 		},
// 	}

// 	for confIndex, conf := range config.Configs {
// 		counterMap := make(map[string]string)
// 		for counterIndex, counter := range conf.Counter {
// 			counterMap["Name"+strconv.Itoa(counterIndex)] = counter.Name
// 			counterMap["Interval"+strconv.Itoa(counterIndex)] = strconv.Itoa(counter.Interval)
// 			counterMap["Path"+strconv.Itoa(counterIndex)] = counter.Path
// 		}

// 		pathElements = append(pathElements, &pb.PathElem{
// 			Name: "Config" + strconv.Itoa(confIndex),
// 			Key:  counterMap,
// 		})
// 	}

// 	setRequest := pb.SetRequest{
// 		Update: []*pb.Update{
// 			{
// 				Path: &pb.Path{
// 					Target: config.DeviceIP,
// 					Elem:   pathElements,
// 				},
// 			},
// 		},
// 	}

// 	response, err := c.(*gclient.Client).Set(ctx, &setRequest)

// 	fmt.Print("Response from device-monitor is: ")
// 	fmt.Println(response)

// 	if err != nil {
// 		fmt.Print("Target returned RPC error for Set: ")
// 		fmt.Println(err)
// 	}
// }

// func sub() {
// 	ctx := context.Background()

// 	address := []string{"monitor-service:11161"}

// 	c, err := gclient.New(ctx, client.Destination{
// 		Addrs:       address,
// 		Target:      "monitor-service",
// 		Timeout:     time.Second * 5,
// 		Credentials: nil,
// 		TLS:         nil,
// 	})

// 	if err != nil {
// 		// fmt.Errorf("could not create a gNMI client: %v", err)
// 		fmt.Print("Could not create a gNMI client: ")
// 		fmt.Println(err)
// 	}

// 	var path []client.Path

// 	path = append(path, []string{"test/testing/lol"})

// 	query := client.Query{
// 		Addrs:  address,
// 		Target: "supa-target",
// 		// Replica: ,
// 		UpdatesOnly: true,
// 		Queries:     path,
// 		Type:        3, // 1 - Once, 2 - Poll, 3 - Stream
// 		// Timeout: ,
// 		// NotificationHandler: callback,
// 		ProtoHandler: protoCallback,
// 		// ...
// 	}

// 	fmt.Println(query)

// 	err = c.(*gclient.Client).Subscribe(ctx, query)

// 	if err != nil {
// 		fmt.Print("Target returned RPC error for Subscribe: ")
// 		fmt.Println(err)
// 	}

// 	fmt.Println("Client connected successfully")

// 	for {
// 		fmt.Println(c.(*gclient.Client).Recv())
// 		time.Sleep(10 * time.Second)
// 	}
// }

// // Updates will be sent here,
// func protoCallback(msg proto.Message) error {
// 	fmt.Print("protoCallback msg: ")
// 	fmt.Println(msg)
// 	return nil
// }
