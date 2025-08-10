package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"cse190-welp/proto/detail"
	"cse190-welp/proto/reservation"
	"cse190-welp/proto/review"
)

type RequestType int

const (
	DETAIL      RequestType = 0
	RESERVATION RequestType = 1
	REVIEW      RequestType = 2
)

// func Schedule(t RequestType) {
// 	log.Printf("Scheduler is scheduling")
// 	// fmt.Printf("Scheduler is scheduling\n")
// 	switch t {
// 	case DETAIL:
// 		cpuMask := 0
// 		ret := C.sched_setaffinity(C.getpid(), C.get_nprocs_conf(), cpuMask)
// 		// ret := C.setAffinity(C.getpid(), C.get_nprocs_conf(), cpuMask)
// 		if ret != 0 {
// 			// error message
// 		}
// 		// C.sched_setaffinity(0, 0, 0)
// 	case RESERVATION:

// 	case REVIEW:

// 	default:
// 		// error message
// 	}
// }

// Frontend implements a service that acts as an interface to interact with different microservices.
type Frontend struct {
	port              int
	detailClient      detail.DetailServiceClient
	reviewClient      review.ReviewServiceClient
	reservationClient reservation.ReservationServiceClient
	User              string
}

// NewFrontend creates a new Frontend instance with the specified configuration.
func NewFrontend(port int, detailaddr string, reviewaddr string, reservationaddr string) *Frontend {
	f := &Frontend{
		port:              port,
		detailClient:      detail.NewDetailServiceClient(dial(detailaddr)),
		reviewClient:      review.NewReviewServiceClient(dial(reviewaddr)),
		reservationClient: reservation.NewReservationServiceClient(dial(reservationaddr)),
		User:              "None",
	}
	return f
}

// Run starts the Frontend server and listens for incoming requests on the specified port.
func (s *Frontend) Run() error {
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/get-detail", s.getDetailHandler)
	http.HandleFunc("/post-detail", s.postDetailHandler)
	http.HandleFunc("/get-review", s.getReviewHandler)
	http.HandleFunc("/post-review", s.postReviewHandler)
	http.HandleFunc("/search-reviews", s.searchReviewsHandler)
	http.HandleFunc("/get-reservation", s.getReservationHandler)
	http.HandleFunc("/make-reservation", s.makeReservationHandler)
	http.HandleFunc("/most-popular", s.mostPopularHandler)

	log.Printf("frontend server running at port dfsajlfsadj: %d", s.port)
	return http.ListenAndServe(fmt.Sprintf(":%d", s.port), nil)
}

// getDetailHandler handles requests for retrieving restaurant details.
func (s *Frontend) getDetailHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	restaurant_name := r.URL.Query().Get("restaurant_name")

	if restaurant_name == "" {
		http.Error(w, "Malformed request to `/get-detail` endpoint!", http.StatusBadRequest)
		return
	}

	req := &detail.GetDetailRequest{RestaurantName: restaurant_name}
	// log.Printf("In get Detail Handler")
	// Schedule(DETAIL)
	reply, err := s.detailClient.GetDetail(ctx, req)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(reply)
}

// postDetailHandler handles requests for posting restaurant details.
func (s *Frontend) postDetailHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	restaurant_name := r.URL.Query().Get("restaurant_name")
	location := r.URL.Query().Get("location")
	style := r.URL.Query().Get("style")
	capacity, err := strconv.Atoi(r.URL.Query().Get("capacity"))

	if restaurant_name == "" || location == "" || style == "" || err != nil {
		http.Error(w, "Malformed request to `/post-detail` endpoint!", http.StatusBadRequest)
		return
	}

	req := &detail.PostDetailRequest{
		RestaurantName: restaurant_name,
		Location:       location,
		Style:          style,
		Capacity:       int32(capacity),
	}

	// syscall.pthread_create()
	// make goroutine
	// core = Schedule(DETAIL)
	// pthread_join(above_thread)
	// log.Printf("In post Detail Handler")
	// Schedule(DETAIL)
	reply, err := s.detailClient.PostDetail(ctx, req)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(reply)
}

// getReviewHandler handles requests for retrieving reviews.
func (s *Frontend) getReviewHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	restaurant_name := r.URL.Query().Get("restaurant_name")
	user_name := r.URL.Query().Get("user_name")

	if restaurant_name == "" || user_name == "" {
		http.Error(w, "Malformed request to `/get-review` endpoint!", http.StatusBadRequest)
		return
	}

	req := &review.GetReviewRequest{
		RestaurantName: restaurant_name,
		UserName:       user_name,
	}
	reply, err := s.reviewClient.GetReview(ctx, req)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(reply)
}

// postReviewHandler handles requests for posting reviews.
func (s *Frontend) postReviewHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user_name := r.URL.Query().Get("user_name")
	restaurant_name := r.URL.Query().Get("restaurant_name")
	restaurant_review := r.URL.Query().Get("review")
	restaurant_rating, err := strconv.Atoi(r.URL.Query().Get("rating"))

	if restaurant_name == "" || user_name == "" || restaurant_review == "" || err != nil {
		http.Error(w, "Malformed request to `/post-review` endpoint!", http.StatusBadRequest)
		return
	}

	req := &review.PostReviewRequest{
		UserName:       user_name,
		RestaurantName: restaurant_name,
		Review:         restaurant_review,
		Rating:         int32(restaurant_rating),
	}
	reply, err := s.reviewClient.PostReview(ctx, req)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(reply)
}

// searchReviewHandler handles requests for searching reviews.
func (s *Frontend) searchReviewsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	restaurant_name := r.URL.Query().Get("restaurant_name")
	if restaurant_name == "" {
		http.Error(w, "Malformed request to `/search-reviews` endpoint!", http.StatusBadRequest)
		return
	}

	req := &review.SearchReviewsRequest{
		RestaurantName: restaurant_name,
	}
	reply, err := s.reviewClient.SearchReviews(ctx, req)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(reply)
}

// getReservationHandler handles requests for retrieving reservations.
func (s *Frontend) getReservationHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user_name := r.URL.Query().Get("user_name")
	restaurant_name := r.URL.Query().Get("restaurant_name")

	if user_name == "" || restaurant_name == "" {
		http.Error(w, "Malformed request to `/get-reservation` endpoint!", http.StatusBadRequest)
		return
	}

	req := &reservation.GetReservationRequest{
		UserName:       user_name,
		RestaurantName: restaurant_name,
	}
	reply, err := s.reservationClient.GetReservation(ctx, req)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(reply)
}

// makeReservationHandler handles requests for making reservations.
func (s *Frontend) makeReservationHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user_name := r.URL.Query().Get("user_name")
	restaurant_name := r.URL.Query().Get("restaurant_name")
	year, year_err := strconv.Atoi(r.URL.Query().Get("year"))
	month, month_err := strconv.Atoi(r.URL.Query().Get("month"))
	day, day_err := strconv.Atoi(r.URL.Query().Get("day"))

	if restaurant_name == "" || user_name == "" || year_err != nil || month_err != nil || day_err != nil {
		http.Error(w, "Malformed request to `/make-reservation` endpoint!", http.StatusBadRequest)
		return
	}

	req := &reservation.MakeReservationRequest{
		UserName:       user_name,
		RestaurantName: restaurant_name,
		Time:           &reservation.Date{Year: int32(year), Month: int32(month), Day: int32(day)},
	}
	reply, err := s.reservationClient.MakeReservation(ctx, req)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(reply)
}

// mostPopularHandler handles requests for retrieving most popular restaurants.
func (s *Frontend) mostPopularHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	topk, err := strconv.Atoi(r.URL.Query().Get("topk"))

	if err != nil {
		http.Error(w, "Malformed request to `/most-popular` endpoint!", http.StatusBadRequest)
		return
	}

	req := &reservation.MostPopularRequest{
		TopK: int32(topk),
	}
	reply, err := s.reservationClient.MostPopular(ctx, req)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(reply)
}
