package services

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	"cse190-welp/proto/detail"
	"cse190-welp/proto/mycache"
	"cse190-welp/proto/mydatabase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

// Detail implements the detail service
type Detail struct {
	name string
	port int
	detail.DetailServiceServer
	lock                 sync.Mutex
	detailCacheClient    mycache.CacheServiceClient
	detailDatabaseClient mydatabase.DatabaseServiceClient
}

// NewDetail returns a new server for the detail service.
func NewDetail(name string, detailPort int, detailCacheAddr string, detailDatabaseAddr string) *Detail {
	return &Detail{
		name:                 name,
		port:                 detailPort,
		detailCacheClient:    mycache.NewCacheServiceClient(dial(detailCacheAddr)),
		detailDatabaseClient: mydatabase.NewDatabaseServiceClient(dial(detailDatabaseAddr)),
	}
}

// Run starts the Detail gRPC server and listens for incoming requests.
// It returns an error if the server fails to start or encounters an error.
func (s *Detail) Run() error {
	// Create a new gRPC server instance.
	srv := grpc.NewServer()

	// Register the Detail server implementation with the gRPC server.
	detail.RegisterDetailServiceServer(srv, s)

	// Create a TCP listener that listens for incoming requests on the specified port.
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// (Optional) Log a message indicating that the server is running and listening on the specified port.
	log.Printf("detail server running at port: %d", s.port)

	// Start serving incoming requests using the registered implementation.
	return srv.Serve(lis)
}

// GetDetail retrieves the details of a restaurant.
// It first checks if the data is cached in mycache.
// If not, it retrieves the data from mydb and stores it in mycache for future use.
// It returns an error if the requested restaurant does not exist.
func (s *Detail) GetDetail(ctx context.Context, req *detail.GetDetailRequest) (*detail.GetDetailResponse, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	// Get the name of the requested restaurant
	restaurantName := req.GetRestaurantName()

	// Check if the data is cached in mycache
	cacheRequest := &mycache.GetItemRequest{Key: restaurantName}
	cacheReply, err := s.detailCacheClient.GetItem(ctx, cacheRequest)
	replyStatus, _ := status.FromError(err)

	item := cacheReply.GetItem()
	detailResponse := &detail.GetDetailResponse{}

	switch replyStatus.Code() {
	case codes.OK:
		err = proto.Unmarshal(item.GetValue(), detailResponse)
		if err != nil {
			log.Fatal(err)
		}
		err = status.Error(codes.OK, "Cache hit while reading from service: mycache-detail")
	case codes.NotFound:
		// Cache miss, go to database
		err = nil // reset cache err to nil
		databaseRequest := &mydatabase.GetRecordRequest{Key: restaurantName}
		databaseReply, err := s.detailDatabaseClient.GetRecord(ctx, databaseRequest)
		databaseReplyStatus, _ := status.FromError(err)
		if databaseReplyStatus.Code() != codes.OK {
			return detailResponse, status.Error(codes.NotFound, "Item does not exist in cache or database")
		}

		// Unmarshal data from database record into response
		record := databaseReply.GetRecord()
		err = proto.Unmarshal(record.GetValue(), detailResponse)
		if err != nil { // err if bytes don't unmarshal
			log.Fatal(err)
		}

		// Populate cache with item
		item := &mycache.CacheItem{
			Key:   record.GetKey(),
			Value: record.GetValue(),
		}
		err = cacheSetHelper(s.detailCacheClient, ctx, item, s.name)
		if err != nil {
			log.Println("Failed to populate cache!") // don't exit if this occurs
		}
	case codes.Canceled:
		err = status.Errorf(codes.Canceled, "Error! GetDetail context canceled with message: %s", replyStatus.Message())
	default:
		log.Fatalf("Unexpected error getting item: %v", err)
	}

	// Return the response object and any error.
	return detailResponse, err
}

// PostDetail adds or updates the details of a restaurant.
func (s *Detail) PostDetail(ctx context.Context, req *detail.PostDetailRequest) (*detail.PostDetailResponse, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	// Extract fields from the request.
	restaurantName := req.GetRestaurantName()
	location := req.GetLocation()
	capacity := req.GetCapacity()
	style := req.GetStyle()

	msg := &detail.GetDetailResponse{
		RestaurantName: restaurantName,
		Location:       location,
		Style:          style,
		Capacity:       capacity,
	}

	data, err := proto.Marshal(msg)
	if err != nil {
		log.Fatal(err)
	}

	// Cache the data in mymcached
	item := &mycache.CacheItem{
		Key:   restaurantName,
		Value: data,
	}

	record := &mydatabase.DatabaseRecord{
		Key:   restaurantName,
		Value: data,
	}

	// Create a protobuf response indicating whether the detail was successfully posted
	detailResponse := &detail.PostDetailResponse{Status: true}

	err = cacheSetHelper(s.detailCacheClient, ctx, item, s.name)
	if err != nil {
		detailResponse.Status = false
	}

	err = storageSetHelper(s.detailDatabaseClient, ctx, record, s.name)
	if err != nil {
		detailResponse.Status = false
	}

	// Return the response object.
	return detailResponse, err
}
