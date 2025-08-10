package services

import (
	"context"
	"fmt"
	"log"
	"net"

	apps "cse190-welp/applications"
	"cse190-welp/proto/mydatabase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// MyDatabase represents a gRPC service for interacting with a database.
type MyDatabase struct {
	name string
	port int
	mydatabase.DatabaseServiceServer
	app *apps.EmulatedStorageApp
}

// NewMyDatabase creates a new instance of MyDatabase.
// serverName: The name of the database server.
// databasePort: The port on which the server should listen.
// deviceType: The type of storage device to use. (ssd, disk, or cloud)
func NewMyDatabase(serverName string, databasePort int, deviceType string) *MyDatabase {
	// Initialize and return a new MyDatabase instance.
	app, err := apps.NewEmulatedStorageApp(deviceType)
	if err != nil {
		log.Fatalf("failed to initialize application: %v", err)
	}
	return &MyDatabase{
		name: serverName,
		port: databasePort,
		app:  app,
	}
}

// Run starts the MyDatabase gRPC server and listens for incoming requests.
// It returns an error if the server fails to start or encounters an error.
func (s *MyDatabase) Run() error {
	// Create a new gRPC server instance.
	srv := grpc.NewServer()

	// Register the Database server implementation with the gRPC server.
	mydatabase.RegisterDatabaseServiceServer(srv, s)

	// Create a TCP listener that listens for incoming requests on the specified port.
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// (Optional) Log a message indicating that the server is running and listening on the specified port.
	log.Printf("storage server <%s> running at port: %d", s.name, s.port)
	return srv.Serve(lis)
}

// GetRecord retrieves a record from the database.
func (s *MyDatabase) GetRecord(ctx context.Context, req *mydatabase.GetRecordRequest) (*mydatabase.GetRecordResponse, error) {
	// Get the name of the requested item
	key := req.GetKey()

	// Retrieve record from the database application
	record, ok := s.app.Get(key)
	msg := &mydatabase.GetRecordResponse{
		Record: record, // will be nil if an error occurs
	}
	var err error
	if !ok {
		err = status.Errorf(codes.NotFound, "Record not found in storage!")
	} else {
		err = status.Error(codes.OK, "Record found in storage!")
	}
	return msg, err
}

// SetRecord sets a record in the database.
func (s *MyDatabase) SetRecord(ctx context.Context, req *mydatabase.SetRecordRequest) (*mydatabase.SetRecordResponse, error) {
	record := req.GetRecord()

	msg := &mydatabase.SetRecordResponse{
		Success: true,
	}
	s.app.Set(record)
	return msg, status.Error(codes.OK, "Record placed in storage!")
}

// DeleteRecord deletes a record from the database.
func (s *MyDatabase) DeleteRecord(ctx context.Context, req *mydatabase.DeleteRecordRequest) (*mydatabase.DeleteRecordResponse, error) {
	key := req.GetKey()
	log.Printf("DeleteKey: %s", key)
	msg := &mydatabase.DeleteRecordResponse{
		Success: true,
	}

	s.app.Delete(key)
	return msg, status.Error(codes.OK, "Record deleted from database!")
}
