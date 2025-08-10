local socket = require("socket")
math.randomseed(socket.gettime()*1000)
math.random(); math.random(); math.random()

char_to_hex = function(c)
    return string.format("%%%02X", string.byte(c))
end

function urlencode(url)
    if url == nil then
        return
    end
    url = url:gsub("\n", "\r\n")
    url = url:gsub("([^%w ])", char_to_hex)
    url = url:gsub(" ", "+")
    return url
end

hex_to_char = function(x)
    return string.char(tonumber(x, 16))
end

urldecode = function(url)
    if url == nil then
        return
    end
    url = url:gsub("+", " ")
    url = url:gsub("%%(%x%x)", hex_to_char)
    return url
end

-- Load the data from the CSV files into tables
local detail_samples = {}
local reservation_samples = {}
local review_samples = {}

local function load_review_samples(file_path, t)
    local file = io.open(file_path, "r")

    -- Skip the first line containing headers
    file:read("*line")

    for line in file:lines() do
        local user_name, restaurant_name, review, rating = line:match("\"(.-)\",\"(.-)\",\"(.-)\",(%d+)")
        if user_name and restaurant_name and review and rating then
            table.insert(t, {user_name = user_name, restaurant_name = restaurant_name, review = review, rating = tonumber(rating)})
        end
    end
    file:close()
end

local function load_reservation_samples(file_path, t)
    local file = io.open(file_path, "r")

    -- Skip the first line containing headers
    file:read("*line")

    for line in file:lines() do
        -- Split the line using commas
        local fields = {}
        for field in line:gmatch("[^,]+") do
            fields[#fields + 1] = field
        end
    
        local user_name, restaurant_name, year, month, day = fields[1], fields[2], tonumber(fields[3]), tonumber(fields[4]), tonumber(fields[5])
        -- print(user_name, restaurant_name, year, month, day) -- (optional) uncomment to look at data
    
        -- Check if the fields are not nil
        if user_name and restaurant_name and year and month and day then
            table.insert(t, {user_name = user_name, restaurant_name = restaurant_name, year = year, month = month, day = day})
        end
    end
    file:close()
end

local function load_detail_samples(file_path, t)
    local file = io.open(file_path, "r")

    -- Skip the first line containing headers
    file:read("*line")

    for line in file:lines() do
        local restaurant_name, location, style, capacity = line:match("\"(.-)\",\"(.-)\",\"(.-)\",(%d+)")
        -- print(restaurant_name, location, style, capacity) -- (optional) toggle to look at data
        if restaurant_name and location and style and capacity then
            table.insert(t, {restaurant_name = restaurant_name, location = location, style = style, capacity = tonumber(capacity)})
        end
    end
    file:close()
end

load_detail_samples("./samples/detail_samples.csv", detail_samples)
load_reservation_samples("./samples/reservation_samples.csv", reservation_samples)
load_review_samples("./samples/review_samples.csv", review_samples)

local function post_detail()
    local method = "GET"

    -- Choose a random sample to post from the loaded detail data
    local random_index = math.random(1, #detail_samples)
    local sample = detail_samples[random_index]

    -- Construct the path using the selected sample
    local path = url .. "/post-detail?restaurant_name=" .. urlencode(sample.restaurant_name) .. "&location=" .. urlencode(sample.location) .. "&style=" .. urlencode(sample.style) .. "&capacity=" .. sample.capacity
    local headers = {}
    return wrk.format(method, path, headers, nil)
end

local function get_detail()
    local method = "GET"

    -- Choose a random sample to get from the loaded detail data
    local random_index = math.random(1, #detail_samples)
    local sample = detail_samples[random_index]

    -- Construct the path using the selected sample
    local path = url .. "/get-detail?restaurant_name=" .. urlencode(sample.restaurant_name)
    local headers = {}
    return wrk.format(method, path, headers, nil)
end

local function make_reservation()
    local method = "GET"

    -- Choose a random sample to post from the loaded reservation data
    local random_index = math.random(1, #reservation_samples)
    local sample = reservation_samples[random_index]

    -- Construct the path using the selected sample
    local path = url .. "/make-reservation?restaurant_name=" .. urlencode(sample.restaurant_name) .. "&user_name=" .. urlencode(sample.user_name) .. "&year=" .. sample.year .. "&month=" .. sample.month .. "&day=" .. sample.day
    local headers = {}
    return wrk.format(method, path, headers, nil)
end

local function get_reservation()
    local method = "GET"

    -- Choose a random sample to get from the loaded reservation data
    local random_index = math.random(1, #reservation_samples)
    local sample = reservation_samples[random_index]

    -- Construct the path using the selected sample
    local path = url .. "/get-reservation?user_name=" .. urlencode(sample.user_name)
    local headers = {}
    return wrk.format(method, path, headers, nil)
end

local function most_popular()
    local method = "GET"

    -- Construct the path using the selected sample
    local topk = tostring(math.random(1, 10))
    local path = url .. "/most-popular?topk=" .. topk
    local headers = {}
    return wrk.format(method, path, headers, nil)
end

local function post_review()
    local method = "GET"

    -- Choose a random sample to post from the loaded review data
    local random_index = math.random(1, #review_samples)
    local sample = review_samples[random_index]

    -- Construct the path using the selected sample
    local path = url .. "/post-review?restaurant_name=" .. urlencode(sample.restaurant_name) .. "&user_name=" .. urlencode(sample.user_name) .. "&review=" .. urlencode(sample.review) .. "&rating=" .. sample.rating
    local headers = {}
    return wrk.format(method, path, headers, nil)
end

local function get_review()
    local method = "GET"

    -- Choose a random sample to get from the loaded review data
    local random_index = math.random(1, #review_samples)
    local sample = review_samples[random_index]

    -- Construct the path using the selected sample
    local path = url .. "/get-review?user_name=" .. urlencode(sample.user_name) .. "&restaurant_name=" .. urlencode(sample.restaurant_name)
    local headers = {}
    return wrk.format(method, path, headers, nil)
end

local function search_reviews()
    local method = "GET"

    -- Choose a random sample to get from the loaded review data
    local random_index = math.random(1, #review_samples)
    local sample = review_samples[random_index]

    -- Construct the path using the selected sample
    local path = url .. "/search-reviews?restaurant_name=" .. urlencode(sample.restaurant_name)
    local headers = {}
    return wrk.format(method, path, headers, nil)
end

request = function()
    cur_time = math.floor(socket.gettime())
    -- Rebalancing the work load to detail:reservation:review = 2:3:0 to make it easier to identify a bottleneck.
    local post_detail_ratio = 2
    local get_detail_ratio = 2
    local make_reservation_ratio = 2
    local get_reservation_ratio = 1
    local most_popular_ratio = 3
    local post_review_ratio = 0
    local get_review_ratio = 0
    local search_reviews_ratio = 0

    local sum_ratio = post_detail_ratio + get_detail_ratio + make_reservation_ratio + get_reservation_ratio + most_popular_ratio + post_review_ratio + get_review_ratio + search_reviews_ratio

    local coin = math.random() * sum_ratio
    if coin < post_detail_ratio then
        return post_detail(url)
    elseif coin < post_detail_ratio + get_detail_ratio then
        return get_detail(url)
    elseif coin < post_detail_ratio + get_detail_ratio + make_reservation_ratio then
        return make_reservation(url)
    elseif coin < post_detail_ratio + get_detail_ratio + make_reservation_ratio + get_reservation_ratio then
        return get_reservation(url)
    elseif coin < post_detail_ratio + get_detail_ratio + make_reservation_ratio + get_reservation_ratio + most_popular_ratio then
        return most_popular(url)
    elseif coin < post_detail_ratio + get_detail_ratio + make_reservation_ratio + get_reservation_ratio + most_popular_ratio + post_review_ratio then
        return post_review(url)
    elseif coin < post_detail_ratio + get_detail_ratio + make_reservation_ratio + get_reservation_ratio + most_popular_ratio + post_review_ratio + get_review_ratio then
        return get_review(url)
    else
        return search_reviews(url)
    end
end
