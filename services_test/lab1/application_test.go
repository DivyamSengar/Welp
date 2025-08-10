package services_test

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"cse190-welp/proto/detail"
	"cse190-welp/proto/reservation"
	"cse190-welp/proto/review"
)

// Programmatically issue an HTTP GET request to the specified URL string
func httpRequest(getURL string) ([]byte, error) {
	respGet, err := http.Get(getURL)
	if err != nil {
		return nil, fmt.Errorf("Failed to make GET request: %v", err)
	}
	defer respGet.Body.Close()

	// Read the GET response body
	body, err := ioutil.ReadAll(respGet.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read response body: %v", err)
	}
	return body, nil
}

func loadCSV(csvFilePath string) ([][]string, error) {
	file, err := os.Open(csvFilePath)
	if err != nil {
		return nil, fmt.Errorf("Failed to open CSV file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("Failed to read CSV records: %v", err)
	}
	return records, nil
}

func testDetailHelper(t *testing.T, expectedData *url.Values) {
	clusterFrontendIP := "10.96.88.88"
	port := 8080

	expectedRestaurantName := expectedData.Get("restaurant_name")
	expectedLocation := expectedData.Get("location")
	expectedCapacity, err := strconv.Atoi(expectedData.Get("capacity"))
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	expectedStyle := expectedData.Get("style")

	// Make a POST request to the URL
	postURL := fmt.Sprintf("http://%s:%v/post-detail?%s", clusterFrontendIP, port, expectedData.Encode())
	rawStatus, err := httpRequest(postURL)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	var postResponse detail.PostDetailResponse
	log.Printf("Raw Response Body: %v", string(rawStatus))
	if err := json.Unmarshal(rawStatus, &postResponse); err != nil {
		t.Fatalf("Failed to decode JSON response: %v", err)
	}
	// Read the POST response body (usually you would parse this response if needed)
	if !postResponse.Status {
		t.Errorf("Expected status name: true \nActual restaurant name: %v", &postResponse.Status)
	}

	// Make a GET request to the URL
	getData := url.Values{}
	getData.Add("restaurant_name", expectedRestaurantName)
	getURL := fmt.Sprintf("http://%s:%v/get-detail?%s", clusterFrontendIP, port, getData.Encode())

	body, err := httpRequest(getURL)
	log.Printf("Raw Response Body: %v", string(body))
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	// Decode the JSON response body into GetDetailRequest struct
	var actualResp detail.GetDetailResponse
	if err := json.Unmarshal(body, &actualResp); err != nil {
		t.Fatalf("Failed to decode JSON response: %v", err)
	}

	// Compare the decoded struct with the expected values
	if actualResp.RestaurantName != expectedRestaurantName {
		t.Errorf("Expected restaurant name: %s\nActual restaurant name: %s", expectedRestaurantName, actualResp.RestaurantName)
	}

	if actualResp.Location != expectedLocation {
		t.Errorf("Expected restaurant name: %s\nActual restaurant name: %s", expectedRestaurantName, actualResp.RestaurantName)
	}

	if actualResp.Style != expectedStyle {
		t.Errorf("Expected restaurant name: %s\nActual restaurant name: %s", expectedRestaurantName, actualResp.RestaurantName)
	}

	if actualResp.Capacity != int32(expectedCapacity) {
		t.Errorf("Expected restaurant name: %s\nActual restaurant name: %s", expectedRestaurantName, actualResp.RestaurantName)
	}
}

func testReservationHelper(t *testing.T, expectedData *url.Values) {
	clusterFrontendIP := "10.96.88.88"
	port := 8080

	// t.Logf("Post data: %v", expectedData)
	expectedUserName := expectedData.Get("user_name")
	expectedRestaurantName := expectedData.Get("restaurant_name")
	expectedYear, err := strconv.Atoi(expectedData.Get("year"))
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	expectedMonth, err := strconv.Atoi(expectedData.Get("month"))
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	expectedDay, err := strconv.Atoi(expectedData.Get("day"))
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	// Make a POST request to the URL
	postURL := fmt.Sprintf("http://%s:%v/make-reservation?%s", clusterFrontendIP, port, expectedData.Encode())
	rawStatus, err := httpRequest(postURL)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	var postResponse reservation.MakeReservationResponse
	log.Printf("Raw Response Body: %v", string(rawStatus))
	if err := json.Unmarshal(rawStatus, &postResponse); err != nil {
		t.Fatalf("Failed to decode JSON response: %v", err)
	}
	// Read the POST response body (usually you would parse this response if needed)
	if !postResponse.Status {
		t.Errorf("Expected status name: true \nActual restaurant name: %v", &postResponse.Status)
	}

	// Make a GET request to the URL
	getData := url.Values{}
	getData.Add("user_name", expectedUserName)
	getData.Add("restaurant_name", expectedRestaurantName)
	getURL := fmt.Sprintf("http://%s:%v/get-reservation?%s", clusterFrontendIP, port, getData.Encode())

	body, err := httpRequest(getURL)
	log.Printf("Raw Response Body: %v", string(body))
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	// Decode the JSON response body into GetDetailRequest struct
	var actualResp reservation.GetReservationResponse
	if err := json.Unmarshal(body, &actualResp); err != nil {
		t.Fatalf("Failed to decode JSON response: %v", err)
	}

	// Compare the decoded struct with the expected values
	if actualResp.RestaurantName != expectedRestaurantName {
		t.Errorf("Expected restaurant name: %s\nActual restaurant name: %s", expectedRestaurantName, actualResp.RestaurantName)
	}

	if actualResp.UserName != expectedUserName {
		t.Errorf("Expected restaurant name: %s\nActual restaurant name: %s", expectedRestaurantName, actualResp.RestaurantName)
	}

	if actualResp.Time.Year != int32(expectedYear) {
		t.Errorf("Expected restaurant name: %s\nActual restaurant name: %s", expectedRestaurantName, actualResp.RestaurantName)
	}

	if actualResp.Time.Month != int32(expectedMonth) {
		t.Errorf("Expected restaurant name: %s\nActual restaurant name: %s", expectedRestaurantName, actualResp.RestaurantName)
	}

	if actualResp.Time.Day != int32(expectedDay) {
		t.Errorf("Expected restaurant name: %s\nActual restaurant name: %s", expectedRestaurantName, actualResp.RestaurantName)
	}
}

func testPopularHelper(t *testing.T, topk int, expectedSlice []string) {
	clusterFrontendIP := "10.96.88.88"
	port := 8080

	data := url.Values{}
	data.Add("topk", fmt.Sprint(topk))
	url := fmt.Sprintf("http://%s:%v/most-popular?%s", clusterFrontendIP, port, data.Encode())
	body, err := httpRequest(url)
	log.Printf("Raw Response Body: %v", string(body))
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	var actualResp reservation.MostPopularResponse
	if err := json.Unmarshal(body, &actualResp); err != nil {
		t.Fatalf("Failed to decode JSON response: %v", err)
	}
	for i := 0; i < topk; i++ {
		expected := expectedSlice[i]
		actual := actualResp.TopKRestaurants[i]
		if actual != expected {
			t.Errorf("Expected %v, but got %v", expected, actual)
		}
	}
}

func testReviewEqual(t *testing.T, actualResp, expected *review.GetReviewResponse) {
	// Compare the decoded struct with the expected values
	if actualResp.RestaurantName != expected.RestaurantName {
		t.Errorf("Expected restaurant name: %s\nActual restaurant name: %s", expected.RestaurantName, actualResp.RestaurantName)
	}

	if actualResp.UserName != expected.UserName {
		t.Errorf("Expected user name: %s\nActual user name: %s", expected.UserName, actualResp.UserName)
	}

	if actualResp.Rating != int32(expected.Rating) {
		t.Errorf("Expected rating: %d\nActual rating: %d", expected.Rating, actualResp.Rating)
	}

	if actualResp.Review != expected.Review {
		t.Errorf("Expected review: %s\nActual review: %s", expected.Review, actualResp.Review)
	}
}

func testReviewHelper(t *testing.T, expectedData *url.Values) *review.GetReviewResponse {
	clusterFrontendIP := "10.96.88.88"
	port := 8080

	expectedUserName := expectedData.Get("user_name")
	expectedRestaurantName := expectedData.Get("restaurant_name")
	expectedReview := expectedData.Get("review")
	expectedRating, err := strconv.Atoi(expectedData.Get("rating"))
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	// Make a POST request to the URL
	postURL := fmt.Sprintf("http://%s:%v/post-review?%s", clusterFrontendIP, port, expectedData.Encode())
	rawStatus, err := httpRequest(postURL)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	var postResponse review.PostReviewResponse
	log.Printf("Raw Response Body: %v", string(rawStatus))
	if err := json.Unmarshal(rawStatus, &postResponse); err != nil {
		t.Fatalf("Failed to decode JSON response: %v", err)
	}
	// Read the POST response body (usually you would parse this response if needed)
	if !postResponse.Status {
		t.Errorf("Expected status name: true \nActual restaurant name: %v", &postResponse.Status)
	}

	// Make a GET request to the URL
	getData := url.Values{}
	getData.Add("user_name", expectedUserName)
	getData.Add("restaurant_name", expectedRestaurantName)
	getURL := fmt.Sprintf("http://%s:%v/get-review?%s", clusterFrontendIP, port, getData.Encode())

	body, err := httpRequest(getURL)
	log.Printf("Raw Response Body: %v", string(body))
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	// Decode the JSON response body into GetDetailRequest struct
	var actualResp review.GetReviewResponse
	if err := json.Unmarshal(body, &actualResp); err != nil {
		t.Fatalf("Failed to decode JSON response: %v", err)
	}

	// Compare the decoded struct with the expected values
	if actualResp.RestaurantName != expectedRestaurantName {
		t.Errorf("Expected restaurant name: %s\nActual restaurant name: %s", expectedRestaurantName, actualResp.RestaurantName)
	}

	if actualResp.UserName != expectedUserName {
		t.Errorf("Expected user name: %s\nActual user name: %s", expectedUserName, actualResp.UserName)
	}

	if actualResp.Rating != int32(expectedRating) {
		t.Errorf("Expected rating: %d\nActual rating: %d", expectedRating, actualResp.Rating)
	}

	if actualResp.Review != expectedReview {
		t.Errorf("Expected review: %s\nActual review: %s", expectedReview, actualResp.Review)
	}
	return &actualResp
}

func testMultiReviewHelper(t *testing.T, expectedSearchData map[string]map[string]*review.GetReviewResponse) {
	clusterFrontendIP := "10.96.88.88"
	port := 8080

	for restName := range expectedSearchData {
		expectedReviews := expectedSearchData[restName]
		urlData := url.Values{}
		urlData.Add("restaurant_name", restName)

		// Make a POST request to the URL
		searchURL := fmt.Sprintf("http://%s:%v/search-reviews?%s", clusterFrontendIP, port, urlData.Encode())
		body, err := httpRequest(searchURL)
		log.Printf("Raw Search Response Body: %v", string(body))
		if err != nil {
			t.Fatalf("Error: %v", err)
		}
		var actualSearchResp review.SearchReviewsResponse
		if err := json.Unmarshal(body, &actualSearchResp); err != nil {
			t.Fatalf("Failed to decode JSON response: %v", err)
		}
		for userName := range expectedReviews {
			if actualReview, ok := actualSearchResp.GetReviewsMap()[userName]; ok {
				expectedReview := expectedReviews[userName]
				testReviewEqual(t, actualReview, expectedReview)
			} else {
				// review doesn't exist, fail now
				t.Errorf("Expected user: %v review for restaurant: %v, instead found nothing.", userName, restName)
			}
		}
	}

}
func TestDetail(t *testing.T) {
	csvFilePath := filepath.Join("..", "..", "samples", "detail_samples.csv")
	t.Log("CSV File Path:", csvFilePath)

	records, err := loadCSV(csvFilePath)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	for row, record := range records {
		if row == 0 {
			continue
		}
		if len(record) >= 4 {
			postData := url.Values{}
			postData.Add("restaurant_name", record[0])
			postData.Add("location", record[1])
			postData.Add("style", record[2])
			postData.Add("capacity", record[3])
			testDetailHelper(t, &postData)
		}
	}
}

func TestReview(t *testing.T) {
	csvFilePath := filepath.Join("..", "..", "samples", "review_samples.csv")
	t.Log("CSV File Path:", csvFilePath)

	records, err := loadCSV(csvFilePath)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	searchKeys := map[string]bool{
		"Chick-fil-A":     false, // dummy boolean value
		"In-N-Out Burger": false, // dummy boolean value
		"Chipotle":        false, // dummy boolean value
	}
	// map of restaurant name to map of user name to GetReview
	expectedSearchData := make(map[string]map[string]*review.GetReviewResponse)
	for key := range searchKeys {
		expectedSearchData[key] = make(map[string]*review.GetReviewResponse)
	}

	// Process post and get
	for row, record := range records {
		if row == 0 {
			continue
		}
		if len(record) >= 4 {
			userName := record[0]
			restName := record[1]

			postData := url.Values{}
			postData.Add("user_name", userName)
			postData.Add("restaurant_name", restName)
			postData.Add("review", record[2])
			postData.Add("rating", record[3])

			// part of search-reviews query, add response to map
			if _, isSearchKey := searchKeys[restName]; isSearchKey {
				resp := testReviewHelper(t, &postData)
				expectedSearchData[restName][userName] = resp
			} else {
				_ = testReviewHelper(t, &postData)
			}
		}
	}

	// Run search-reviews
	testMultiReviewHelper(t, expectedSearchData)

}

func TestReservation(t *testing.T) {
	csvFilePath := filepath.Join("..", "..", "samples", "reservation_samples.csv")
	t.Log("CSV File Path:", csvFilePath)

	records, err := loadCSV(csvFilePath)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	// Process post-get
	for row, record := range records {
		if row == 0 {
			continue
		}
		if len(record) >= 5 {
			postData := url.Values{}
			postData.Add("user_name", record[0])
			postData.Add("restaurant_name", record[1])
			postData.Add("year", record[2])
			postData.Add("month", record[3])
			postData.Add("day", record[4])
			testReservationHelper(t, &postData)
		}
	}

	// Run most-popular
	// ordering of equal numbered elements doesn't matter
	expectedSlice := []string{"In-N-Out Burger", "Chick-fil-A", "Chipotle", "Starbucks"}
	topk := len(expectedSlice)
	testPopularHelper(t, topk, expectedSlice)
}
