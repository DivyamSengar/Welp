import argparse
import csv
import subprocess
import requests
import sys
import os

detailCacheCapacity = 10
reviewCacheCapacity = 10

datasetMultiplier   = 10
detailDatasetSize   = detailCacheCapacity * datasetMultiplier
reviewDatasetSize   = reviewCacheCapacity * datasetMultiplier

def execute_post_detail(restaurant_name, location, style, capacity):
    url = f"http://10.96.88.88:8080/post-detail?restaurant_name={restaurant_name}&location={location}&style={style}&capacity={capacity}"
#    print(url)
    response = requests.post(url)
    if response.status_code != 200:
        print("Error! PostDetail failed with status code:", response.status_code)
        print(url)
        sys.exit(1)

def execute_post_review(user_name, restaurant_name, review, rating):
    url = f"http://10.96.88.88:8080/post-review?user_name={user_name}&restaurant_name={restaurant_name}&review={review}&rating={rating}"
#    print(url)
    response = requests.post(url)
    if response.status_code != 200:
        print("Error! PostReview failed with status code:", response.status_code)
        sys.exit(1)

if __name__ == '__main__':
    for i in range(1, detailDatasetSize + 1):
        print("adding detail %d" % i)
        execute_post_detail("restaurant" + str(i), "location" + str(i), "style" + str(i), str(i))

    for i in range(1, reviewDatasetSize + 1):
        print("adding review %d" % i)
        execute_post_review("user" + str(i), "restaurant" + str(i), "review" + str(i), str(i))
