package services

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	"cse190-welp/proto/mycache"
	"cse190-welp/proto/mydatabase"
	"cse190-welp/proto/review"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

// Review implements the review service
type Review struct {
	name string
	port int
	review.ReviewServiceServer
	reviewCacheClient    mycache.CacheServiceClient
	reviewDatabaseClient mydatabase.DatabaseServiceClient
	idLookupTable        map[string]map[string]struct{} // Map of restaurant names to sets of review IDs
	lock                 sync.Mutex
}

// NewReview returns a new server
func NewReview(name string, reviewPort int, reviewCacheAddr string, reviewDatabaseAddr string) *Review {
	return &Review{
		name:                 name,
		port:                 reviewPort,
		reviewCacheClient:    mycache.NewCacheServiceClient(dial(reviewCacheAddr)),
		reviewDatabaseClient: mydatabase.NewDatabaseServiceClient(dial(reviewDatabaseAddr)),
		idLookupTable:        make(map[string]map[string]struct{}),
	}
}

func (s *Review) updateLookupTable(restaurantName string, reviewID string) error {
	set, exists := s.idLookupTable[restaurantName]
	if !exists {
		set = make(map[string]struct{})
		s.idLookupTable[restaurantName] = set
	}
	set[reviewID] = struct{}{}
	return nil
}

func (s *Review) getFromLookupTable(restaurantName string) ([]string, error) {
	set, exists := s.idLookupTable[restaurantName]
	if !exists {
		// no IDs associated with restaurantName
		return nil, nil
	}
	result := make([]string, 0, len(set))
	for value := range set {
		result = append(result, value)
	}
	return result, nil
}

func (s *Review) getResponseHelper(ctx context.Context, reviewID string) (*review.GetReviewResponse, error) {
	// Check if the data is cached in mycache
	cacheRequest := &mycache.GetItemRequest{Key: reviewID}
	cacheReply, err := s.reviewCacheClient.GetItem(ctx, cacheRequest)
	cacheReplyStatus, _ := status.FromError(err)

	item := cacheReply.GetItem()
	reviewResponse := &review.GetReviewResponse{}

	switch cacheReplyStatus.Code() {
	case codes.OK: // inside cache
		err = proto.Unmarshal(item.GetValue(), reviewResponse)
		if err != nil {
			log.Fatal(err)
		}
		err = status.Error(codes.OK, "Cache hit while reading from service: mycache-review")
	case codes.NotFound:
		err = nil
		// Cache miss, go to database
		databaseRequest := &mydatabase.GetRecordRequest{Key: reviewID}
		databaseReply, err := s.reviewDatabaseClient.GetRecord(ctx, databaseRequest)
		databaseReplyStatus, _ := status.FromError(err)
		if databaseReplyStatus.Code() != codes.OK {
			return reviewResponse, status.Error(codes.NotFound, "Item does not exist in cache or database")
		}

		// Unmarshal data from database record into response
		record := databaseReply.GetRecord()
		err = proto.Unmarshal(record.GetValue(), reviewResponse)
		if err != nil { // err if bytes don't unmarshal
			log.Fatal(err)
		}

		// Populate cache with item
		item := &mycache.CacheItem{
			Key:   record.GetKey(),
			Value: record.GetValue(),
		}
		err = cacheSetHelper(s.reviewCacheClient, ctx, item, s.name)
		if err != nil {
			log.Println("failed to populate cache!") // don't fail if this occurs
		}
	case codes.Canceled:
		err = status.Errorf(codes.Canceled, "Error! GetReview context canceled with message: %s", cacheReplyStatus.Message())
	default:
		// This should NOT happen, and we should restart the container
		log.Fatalf("Unexpected error getting item: %v", err)
	}

	return reviewResponse, err
}

// Run starts the Review gRPC server and listens for incoming requests.
// It returns an error if the server fails to start or encounters an error.
func (s *Review) Run() error {
	// Create a new gRPC server instance.
	srv := grpc.NewServer()

	// Register the Review server implementation with the gRPC server.
	review.RegisterReviewServiceServer(srv, s)

	// Create a TCP listener that listens for incoming requests on the specified port.
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// (Optional) Log a message indicating that the server is running and listening on the specified port.
	log.Printf("review server running at port: %d", s.port)

	// Start serving incoming requests using the registered implementation.
	return srv.Serve(lis)
}

// GetReview returns the review of a restaurant
func (s *Review) GetReview(ctx context.Context, req *review.GetReviewRequest) (*review.GetReviewResponse, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	// Get the restaurant and user names
	restaurantName := req.GetRestaurantName()
	userName := req.GetUserName()

	reviewID, _ := GetQueryUUID(restaurantName, userName)
	return s.getResponseHelper(ctx, reviewID)
}

func (s *Review) SearchReviews(ctx context.Context, req *review.SearchReviewsRequest) (*review.SearchReviewsResponse, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	restaurantName := req.GetRestaurantName()

	// maps usernames to review responses
	userReviews := make(map[string]*review.GetReviewResponse)

	var reviewIDs []string
	reviewIDs, _ = s.getFromLookupTable(restaurantName)

	for _, reviewID := range reviewIDs {
		r, err := s.getResponseHelper(ctx, reviewID)
		if err != nil {
			return &review.SearchReviewsResponse{}, err
		}
		userReviews[r.UserName] = r
	}
	return &review.SearchReviewsResponse{ReviewsMap: userReviews}, nil
}

// PostReview posts a review of a restaurant
func (s *Review) PostReview(ctx context.Context, req *review.PostReviewRequest) (*review.PostReviewResponse, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	// Get the name, review, rating from the request
	restaurantName := req.GetRestaurantName()
	userName := req.GetUserName()
	restaurantReview := req.GetReview()
	restaurantRating := req.GetRating()

	// Create a protobuf message containing the review data
	msg := &review.GetReviewResponse{
		UserName:       userName,
		RestaurantName: restaurantName,
		Review:         restaurantReview,
		Rating:         restaurantRating,
	}

	data, err := proto.Marshal(msg)
	if err != nil {
		log.Fatal(err)
	}

	reviewID, _ := GetQueryUUID(restaurantName, userName)
	s.updateLookupTable(restaurantName, reviewID)

	// Cache the data in mycache
	item := &mycache.CacheItem{
		Key:   reviewID,
		Value: data,
	}

	record := &mydatabase.DatabaseRecord{
		Key:   reviewID,
		Value: data,
	}

	// Create a protobuf response indicating whether the review was successfully posted
	reviewResponse := &review.PostReviewResponse{
		Status: true,
	}

	err = cacheSetHelper(s.reviewCacheClient, ctx, item, s.name)
	if err != nil {
		reviewResponse.Status = false
	}

	err = storageSetHelper(s.reviewDatabaseClient, ctx, record, s.name)
	if err != nil {
		reviewResponse.Status = false
	}

	return reviewResponse, err
}
