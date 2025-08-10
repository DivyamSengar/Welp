package main

import (
	"flag"
	"log"
	"os"
	"runtime"

	services "cse190-welp/services"
)

type server interface {
	Run() error
}

func main() {
	// Define the flags to specify port numbers and addresses
	var (
		frontendPort     = flag.Int("frontend", 8080, "frontend server port")
		detailPort       = flag.Int("detailport", 8081, "detail service port")
		reviewPort       = flag.Int("reviewport", 8082, "review service port")
		reservationsPort = flag.Int("reservationport", 8083, "reservation service port")

		detailAddr      = flag.String("detailaddr", "detail:8081", "detail service address")
		reviewAddr      = flag.String("reviewaddr", "review:8082", "review service addr")
		reservationAddr = flag.String("reservationaddr", "reservation:8083", "reservation service addr")

		cachePort            = flag.Int("cacheport", 11211, "port used by all caches")
		detailCacheAddr      = flag.String("detail_mycache_addr", "mycache-detail:11211", "details mycache address")
		reviewCacheAddr      = flag.String("review_mycache_addr", "mycache-review:11211", "review mycache address")
		reservationCacheAddr = flag.String("reservation_mycache_addr", "mycache-reservation:11211", "reservation mycache address")

		detailCacheCapacity      = flag.Int("detail_mycache_capacity", 10, "maximum number of K-V entries allowed in the detail cache service")
		reviewCacheCapacity      = flag.Int("review_mycache_capacity", 10, "maximum number of K-V entries allowed in the review cache service")
		reservationCacheCapacity = flag.Int("reservation_mycache_capacity", 10, "maximum number of K-V entries allowed in the reservation cache service")

		databasePort            = flag.Int("databaseport", 27017, "port used by all databases")
		storageDeviceType       = flag.String("storage_device_type", "cloud", "specifies emulated storage device type, e.g. option `ssd`, `disk`, or `cloud`")
		detailDatabaseAddr      = flag.String("detail_mydatabase_addr", "mydatabase-detail:27017", "details mydatabase address")
		reviewDatabaseAddr      = flag.String("review_mydatabase_addr", "mydatabase-review:27017", "review mydatabase address")
		reservationDatabaseAddr = flag.String("reservation_mydatabase_addr", "mydatabase-reservation:27017", "reservation mydatabase address")
	)

	// Limit to 1 thread
	runtime.GOMAXPROCS(1)

	// Parse the flags
	flag.Parse()

	var srv server
	var cmd = os.Args[1]

	// Switch statement to create the correct service based on the command
	switch cmd {
	case "frontend":
		// Create a new frontend service with the specified ports and addresses
		srv = services.NewFrontend(
			*frontendPort,
			*detailAddr,
			*reviewAddr,
			*reservationAddr,
		)
	case "detail":
		switch {
		case len(os.Args) < 3:
			// Create a new detail service with the specified port
			srv = services.NewDetail(
				"detail",
				*detailPort,
				*detailCacheAddr,
				*detailDatabaseAddr,
			)
		case os.Args[2] == "cache":
			srv = services.NewMyCache(
				"detail-cache",
				*cachePort,
				*detailCacheCapacity,
			)
		case os.Args[2] == "database":
			srv = services.NewMyDatabase(
				"detail-database",
				*databasePort,
				*storageDeviceType,
			)
		default:
			log.Fatalf("unknown subcmd for detail service: %s", os.Args[2])
		}
	case "reservation":
		switch {
		case len(os.Args) < 3:
			// Create a new reservation service with the specified port
			srv = services.NewReservation(
				"reservation",
				*reservationsPort,
				*reservationCacheAddr,
				*reservationDatabaseAddr,
			)
		case os.Args[2] == "cache":
			srv = services.NewMyCache(
				"reservation-cache",
				*cachePort,
				*reservationCacheCapacity,
			)
		case os.Args[2] == "database":
			srv = services.NewMyDatabase(
				"reservation-database",
				*databasePort,
				*storageDeviceType,
			)
		default:
			log.Fatalf("unknown subcmd for reservation service: %s", os.Args[2])
		}
	case "review":
		switch {
		case len(os.Args) < 3:
			// Create a new review service with the specified port
			srv = services.NewReview(
				"review",
				*reviewPort,
				*reviewCacheAddr,
				*reviewDatabaseAddr,
			)
		case os.Args[2] == "cache":
			srv = services.NewMyCache(
				"review-cache",
				*cachePort,
				*reviewCacheCapacity,
			)
		case os.Args[2] == "database":
			srv = services.NewMyDatabase(
				"review-database",
				*databasePort,
				*storageDeviceType,
			)
		default:
			log.Fatalf("unknown subcmd for review service: %s", os.Args[2])
		}
	default:
		// If an unknown command is provided, log an error and exit
		log.Fatalf("unknown cmd: %s", cmd)
	}

	// Start the server and log any errors that occur
	if err := srv.Run(); err != nil {
		log.Fatalf("run %s error: %v", cmd, err)
	}
}
