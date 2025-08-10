import argparse
import csv
import subprocess
import requests
import sys
import os

def execute_post_detail(restaurant_name, location, style, capacity):
    url = f"http://10.96.88.88:8080/post-detail?restaurant_name={restaurant_name}&location={location}&style={style}&capacity={capacity}"
    print(url)
    response = requests.post(url)
    if response.status_code != 200:
        print("Error! PostReview failed with status code:", response.status_code)
        print(url)
        sys.exit(1)

def execute_get_detail(restaurant_name):
    url = f"http://10.96.88.88:8080/get-detail?restaurant_name={restaurant_name}"
    response = requests.post(url)
    if response.status_code != 200:
        print(url)
        print("Error! GetReview failed with status code:", response.status_code)
        print("Message content:", response.content)
        sys.exit(1)
    print(response.content)

# Post reviews using sample data
def execute_post_review(user_name, restaurant_name, review, rating):
    url = f"http://10.96.88.88:8080/post-review?user_name={user_name}&restaurant_name={restaurant_name}&review={review}&rating={rating}"
    print(url)
    response = requests.post(url)
    if response.status_code != 200:
        print("Error! PostReview failed with status code:", response.status_code)
        sys.exit(1)

def execute_get_review(restaurant_name):
    url = f"http://10.96.88.88:8080/get-review?restaurant_name={restaurant_name}"
    response = requests.post(url)
    if response.status_code != 200:
        print("Error! GetReview failed with status code:", response.status_code)
        print(url)
        sys.exit(1)
    print(response.content)

# Make reservations using sample data
def execute_make_reservation(user_name, restaurant_name, year, month, day):
    url = f"http://10.96.88.88:8080/make-reservation?user_name={user_name}&restaurant_name={restaurant_name}&year={year}&month={month}&day={day}"
    print(url)
    response = requests.post(url)
    if response.status_code != 200:
        print("Error! MakeReservation failed with status code:", response.status_code)
        sys.exit(1)

def execute_get_reservation(user_name):
    url = f"http://10.96.88.88:8080/get-reservation?user_name={user_name}"
    response = requests.post(url)
    if response.status_code != 200:
        print("Error! GetReservation failed with status code:", response.status_code)
        sys.exit(1)
    print(response.content)

# Read and populate review samples
def post_reviews(filename, op='GET'):
    with open(filename, mode="r") as file:
        reader = csv.DictReader(file, fieldnames=["UserName", "RestaurantName", "Review", "Rating"])
        # skip header column
        next(reader)
        for row in reader:
            user_name = row["UserName"]
            restaurant_name = row["RestaurantName"]
            review = row["Review"]
            rating = row["Rating"]
            # print(user_name, "|", restaurant_name, "|", review, "|", rating)
            if op == 'PUT':
                execute_post_review(user_name, restaurant_name, review, rating)
            elif op == 'GET':
                execute_get_review(restaurant_name)

def make_reservations(filename, op='GET'):
    with open(filename, mode="r") as file:
        reader = csv.DictReader(file, fieldnames=["UserName", "RestaurantName", "Year", "Month", "Day"])
        # skip header column
        next(reader)
        for row in reader:
            user_name = row["UserName"]
            restaurant_name = row["RestaurantName"]
            year = row["Year"]
            month = row["Month"]
            day = row["Day"]
            # print(user_name, "|", restaurant_name, "|", year, "|", month,"|", day)
            if op == 'PUT':
                execute_make_reservation(user_name, restaurant_name, year, month, day)
            elif op == 'GET':
                execute_get_reservation(user_name)
            
def post_details(filename, op='GET'):
    with open(filename, mode="r") as file:
        reader = csv.DictReader(file, fieldnames=["RestaurantName","Location","Style","Capacity"])
        # skip header column
        next(reader)
        for row in reader:
            restaurant_name = row["RestaurantName"]
            location = row["Location"]
            style = row["Style"]
            capacity = row["Capacity"]
            # print(user_name, "|", restaurant_name, "|", review, "|", rating)
            if op == 'PUT':
                execute_post_detail(restaurant_name, location, style, capacity)
            elif op == 'GET':
                execute_get_detail(restaurant_name)
        

if __name__ == '__main__':
    parser = argparse.ArgumentParser(description='Process CSV files in a specified directory')
    parser.add_argument('service', help="Detail, Review, or Reservation...")
    parser.add_argument('directory', help='Path to the directory containing CSV files for desired service')
    parser.add_argument('service_op', help="PUT or GET")

    # Parse the command-line arguments
    args = parser.parse_args()

    # Retrieve the directory path from the parsed arguments
    directory = args.directory
    service = args.service
    service_op = args.service_op
    print("Operation:", service_op)

    review_samples = os.path.join(directory, "review_samples.csv")
    reservation_samples = os.path.join(directory, "reservation_samples.csv")
    detail_samples = os.path.join(directory, "detail_samples.csv")

    if service == 'review':
        post_reviews(review_samples, service_op)
    elif service == 'reservation':
        make_reservations(reservation_samples, service_op)
    elif service == 'detail':
        post_details(detail_samples, service_op)
    elif service == 'all':
        post_reviews(review_samples, service_op)
        make_reservations(reservation_samples, service_op)
        post_details(detail_samples, service_op)
    else:
        print(service)
        print("usage: python generate-samples.py <service_name> <samples_dir> <PUT|GET>")


