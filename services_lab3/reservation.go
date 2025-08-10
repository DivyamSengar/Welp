package services

import (
	"context"
	"fmt"
	"log"
	"net"
	"sort"
	"sync"

	"cse190-welp/proto/mycache"
	"cse190-welp/proto/mydatabase"
	"cse190-welp/proto/reservation"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

// Reservation implements the reservation service
type Reservation struct {
	name string
	port int
	reservation.ReservationServiceServer
	reservationCacheClient    mycache.CacheServiceClient
	reservationDatabaseClient mydatabase.DatabaseServiceClient
	popularityTable           map[string]int
	lock                      sync.Mutex // Mutex to synchronize access to popularityTable
}

// NewReservation returns a new server
func NewReservation(name string, reservationPort int, reservationCacheAddr string, reservationDatabaseAddr string) *Reservation {
	return &Reservation{
		name:                      name,
		port:                      reservationPort,
		reservationCacheClient:    mycache.NewCacheServiceClient(dial(reservationCacheAddr)),
		reservationDatabaseClient: mydatabase.NewDatabaseServiceClient(dial(reservationDatabaseAddr)),
		popularityTable:           make(map[string]int),
	}
}

// Run starts the Reservation gRPC server and listens for incoming requests.
// It returns an error if the server fails to start or encounters an error.
func (s *Reservation) Run() error {
	// Create a new gRPC server instance.
	srv := grpc.NewServer()

	// Register the Reservation server implementation with the gRPC server.
	reservation.RegisterReservationServiceServer(srv, s)

	// Create a TCP listener that listens for incoming requests on the specified port.
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// (Optional) Log a message indicating that the server is running and listening on the specified port.
	log.Printf("reservation server running at port: %d", s.port)

	// Start serving incoming requests using the registered implementation.
	return srv.Serve(lis)
}

func (s *Reservation) GetReservation(ctx context.Context, req *reservation.GetReservationRequest) (*reservation.GetReservationResponse, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	// Get the restaurant and user names
	restaurantName := req.GetRestaurantName()
	userName := req.GetUserName()

	// Check if the data is cached in mycache
	reservationID, _ := GetQueryUUID(restaurantName, userName)

	cacheRequest := &mycache.GetItemRequest{Key: reservationID}
	cacheReply, err := s.reservationCacheClient.GetItem(ctx, cacheRequest)
	replyStatus, _ := status.FromError(err)

	item := cacheReply.GetItem()
	reservationResponse := &reservation.GetReservationResponse{}

	switch replyStatus.Code() {
	case codes.OK:
		// Unmarshal the cached data into the reservationResponse struct
		err = proto.Unmarshal(item.Value, reservationResponse)
		if err != nil {
			log.Fatal(err)
		}
		err = status.Error(codes.OK, "Cache hit while reading from service: mycache-reservation")
	case codes.NotFound:
		err = nil
		// Cache miss, go to database
		databaseRequest := &mydatabase.GetRecordRequest{Key: reservationID}
		databaseReply, err := s.reservationDatabaseClient.GetRecord(ctx, databaseRequest)
		databaseReplyStatus, _ := status.FromError(err)
		if databaseReplyStatus.Code() != codes.OK {
			return reservationResponse, status.Error(codes.NotFound, "Item does not exist in cache or database")
		}

		// Unmarshal data from database record into response
		record := databaseReply.GetRecord()
		err = proto.Unmarshal(record.GetValue(), reservationResponse)
		if err != nil { // err if bytes don't unmarshal
			log.Fatal(err)
		}

		// Populate cache with item
		item := &mycache.CacheItem{
			Key:   record.GetKey(),
			Value: record.GetValue(),
		}
		err = cacheSetHelper(s.reservationCacheClient, ctx, item, s.name)
		if err != nil {
			log.Println("failed to populate cache!") // don't fail if this occurs
		}
	case codes.Canceled:
		err = status.Errorf(codes.Canceled, "Error! GetReservation context canceled with message: %s", replyStatus.Message())
	default:
		// This should NOT happen, and we should restart the container
		log.Fatalf("Unexpected error getting item: %v", err)
	}

	return reservationResponse, err
}

// This function takes a MakeReservationRequest message, saves the reservation to the database and caches it in mycache.
func (s *Reservation) MakeReservation(ctx context.Context, req *reservation.MakeReservationRequest) (*reservation.MakeReservationResponse, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	// Retrieve the user's name, restaurant name, and time from the request message
	userName := req.GetUserName()
	restaurantName := req.GetRestaurantName()
	time := req.GetTime()

	// Create a GetReservationResponse message containing the reservation information
	msg := &reservation.GetReservationResponse{
		UserName:       userName,
		RestaurantName: restaurantName,
		Time:           time,
	}

	// Safely update Restaurant popularity using a mutex
	s.popularityTable[restaurantName]++

	// Marshal the message to binary data for storage
	data, err := proto.Marshal(msg)
	if err != nil {
		log.Fatal(err)
	}

	// Cache the data in mycache
	reservationID, _ := GetQueryUUID(restaurantName, userName)

	item := &mycache.CacheItem{
		Key:   reservationID,
		Value: data,
	}
	record := &mydatabase.DatabaseRecord{
		Key:   reservationID,
		Value: data,
	}

	// Create a protobuf response indicating whether the reservation was successfully posted
	reservationResponse := &reservation.MakeReservationResponse{Status: true}

	err = cacheSetHelper(s.reservationCacheClient, ctx, item, s.name)
	if err != nil {
		reservationResponse.Status = false
	}

	err = storageSetHelper(s.reservationDatabaseClient, ctx, record, s.name)
	if err != nil {
		reservationResponse.Status = false
	}

	return reservationResponse, err
}

func (s *Reservation) MostPopular(ctx context.Context, req *reservation.MostPopularRequest) (*reservation.MostPopularResponse, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	topK := int(req.GetTopK())

	keys := make([]string, 0, len(s.popularityTable))

	for key := range s.popularityTable {
		keys = append(keys, key)
	}

	sort.SliceStable(keys, func(i, j int) bool {
		return s.popularityTable[keys[i]] > s.popularityTable[keys[j]]
	})

	topKeys := make([]string, 0, topK)
	for i := 0; i < topK && i < len(keys); i++ {
		topKeys = append(topKeys, keys[i])
	}

	resp := &reservation.MostPopularResponse{
		TopKRestaurants: topKeys,
	}
	return resp, nil
}
