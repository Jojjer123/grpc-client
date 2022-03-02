package adminServer

import (
	"golang.org/x/net/context"
)

type Server struct {
}

// Rename this file to adminServer.go

func (s *Server) MonitorDevice(ctx context.Context, message *MonitorMessage) (*MonitorResponse, error) {
	return &MonitorResponse{Respone: "Successfully created monitor"}, nil
}
