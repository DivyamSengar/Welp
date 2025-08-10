package services

import (
	"context"
	"fmt"
	"log"
	"net"

	apps "cse190-welp/applications"
	"cse190-welp/proto/mycache"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// MyCache represents a gRPC service for interacting with a cache.
type MyCache struct {
	name string
	port int
	mycache.CacheServiceServer
	app *apps.LRUCacheApp
}

// NewMyCache creates a new instance of MyCache.
// serverName: The name of the cache server.
// cachePort: The port on which the server should listen.
// capacity: The maximum capacity of the cache.
func NewMyCache(serverName string, cachePort int, capacity int) *MyCache {
	return &MyCache{
		name: serverName,
		port: cachePort,
		app:  apps.NewLRUCacheApp(capacity),
	}
}

// Run starts the MyCache gRPC server and listens for incoming requests.
// It returns an error if the server fails to start or encounters an error.
func (s *MyCache) Run() error {
	// Create a new gRPC server instance.
	srv := grpc.NewServer()

	// Register the Cache server implementation with the gRPC server.
	mycache.RegisterCacheServiceServer(srv, s)

	// Create a TCP listener that listens for incoming requests on the specified port.
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// (Optional) Log a message indicating that the server is running and listening on the specified port.
	log.Printf("cache server <%s> running at port: %d", s.name, s.port)
	return srv.Serve(lis)
}

// GetItem retrieves an item from the cache.
func (s *MyCache) GetItem(ctx context.Context, req *mycache.GetItemRequest) (*mycache.GetItemResponse, error) {
	// TODO: implement GetItem function
	item, error := s.app.Get(req.Key)
	if error == nil {
		return &mycache.GetItemResponse{Item: item}, status.Errorf(codes.OK, "Item gotten successfully")
	} else {
		return &mycache.GetItemResponse{}, status.Errorf(codes.NotFound, "Key not found in cache")
	}
}

// SetItem sets an item in the cache.
func (s *MyCache) SetItem(ctx context.Context, req *mycache.SetItemRequest) (*mycache.SetItemResponse, error) {
	// TODO: implement SetItem function
	item := req.Item
	error := s.app.Set(item)
	if error != nil {
		return &mycache.SetItemResponse{}, status.Errorf(codes.Unknown, "Item could not be set in cache")
	} else {
		return &mycache.SetItemResponse{Success: true}, status.Errorf(codes.OK, "Item set successfully")
	}
}

// DeleteItem deletes an item from the cache.
func (s *MyCache) DeleteItem(ctx context.Context, req *mycache.DeleteItemRequest) (*mycache.DeleteItemResponse, error) {
	// TODO: implement DeleteItem function
	key := req.Key
	error := s.app.Delete(key)
	if error != nil {
		return &mycache.DeleteItemResponse{}, status.Errorf(codes.NotFound, "Item to be deleted not found in cache")
	} else {
		return &mycache.DeleteItemResponse{}, status.Errorf(codes.OK, "Item delete successfully")
	}
}
