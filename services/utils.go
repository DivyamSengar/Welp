package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"cse190-welp/proto/mycache"
	"cse190-welp/proto/mydatabase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func UnaryClientInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (err error) {
	start := time.Now()
	defer func() {
		in, _ := json.Marshal(req)
		out, _ := json.Marshal(reply)
		inStr, outStr := string(in), string(out)
		duration := int64(time.Since(start).Microseconds())

		delimiter := ";"
		errStr := fmt.Sprintf("%v", err)
		if err == nil {
			errStr = "<nil>"
		}
		logMessage := fmt.Sprintf("grpc%s%s%s%s%s%s%s%s%s%d", delimiter, method, delimiter, inStr, delimiter, outStr, delimiter, errStr, delimiter, duration)
		log.Println(logMessage)

	}()

	return invoker(ctx, method, req, reply, cc, opts...)
}

// dial creates a new gRPC client connection to the specified address and returns a client connection object.
func dial(addr string) *grpc.ClientConn {
	// Define gRPC dial options for the client connection.
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(UnaryClientInterceptor),
	}

	// Create a new gRPC client connection to the specified address using the dial options.
	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		// If there was an error creating the client connection, panic with an error message.
		panic(fmt.Sprintf("ERROR: dial error: %v", err))
	}

	// Return the created client connection object.
	return conn
}

// Retrieves UUID associated with a particular restaurant and user pair. Useful for mapping reviews and/or reservations to data
func GetQueryUUID(restaurantName, userName string) (string, error) {
	combinedString := restaurantName + userName
	uuidNamespace := uuid.NewSHA1(uuid.Nil, []byte(combinedString))
	uuidBytes := uuidNamespace[:]
	uuidObj, err := uuid.FromBytes(uuidBytes)
	if err != nil {
		return "", err
	}
	return uuidObj.String(), nil
}

// private helper function
func cacheSetHelper(client mycache.CacheServiceClient, ctx context.Context, item *mycache.CacheItem, serverName string) error {
	cacheRequest := &mycache.SetItemRequest{Item: item}
	_, err := client.SetItem(ctx, cacheRequest)
	cacheReplyStatus, _ := status.FromError(err)

	switch cacheReplyStatus.Code() {
	case codes.OK:
		err = status.Errorf(codes.OK, "Successfully cached for service: %s", serverName)
	case codes.Canceled:
		err = status.Errorf(codes.Canceled, "Error! Service %s context canceled with message: %s", serverName, cacheReplyStatus.Message())
	default:
		log.Fatal(err)
	}
	return err
}

// private helper function
func storageSetHelper(client mydatabase.DatabaseServiceClient, ctx context.Context, record *mydatabase.DatabaseRecord, serverName string) error {
	databaseRequest := &mydatabase.SetRecordRequest{Record: record}
	_, err := client.SetRecord(ctx, databaseRequest)
	databaseReplyStatus, _ := status.FromError(err)

	switch databaseReplyStatus.Code() {
	case codes.OK:
		err = status.Errorf(codes.OK, "Successfully placed in database: %s", serverName)
	case codes.Canceled:
		err = status.Errorf(codes.Canceled, "Error! Service %s context canceled with message: %s", serverName, databaseReplyStatus.Message())
	default:
		log.Fatal(err)
	}
	return err
}
