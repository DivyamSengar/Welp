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

-- Load the data from the CSV file into a table
local samples = {}
local file = io.open("./samples/review_samples.csv", "r")

-- Skip the first line containing headers
file:read("*line")

for line in file:lines() do
    local user_name, restaurant_name, review, rating = line:match("\"(.-)\",\"(.-)\",\"(.-)\",(%d+)")
    if user_name and restaurant_name and review and rating then
        table.insert(samples, {user_name = user_name, restaurant_name = restaurant_name, review = review, rating = tonumber(rating)})
    end
end
file:close()

local function post_review()
    local method = "GET"

    -- Choose a random sample to post from the loaded detail data
    local random_index = math.random(1, #samples)
    local sample = samples[random_index]

    -- Construct the path using the selected sample
    local path = url .. "/post-review?restaurant_name=" .. urlencode(sample.restaurant_name) .. "&user_name=" .. urlencode(sample.user_name) .. "&review=" .. urlencode(sample.review) .. "&rating=" .. sample.rating
    -- print(path) -- (optional) uncomment me to print the URL query!
    local headers = {}
    return wrk.format(method, path, headers, nil)
end

local function get_review()
    local method = "GET"

    -- Choose a random sample to get from the loaded review data
    local random_index = math.random(1, #samples)
    local sample = samples[random_index]

    -- Construct the path using the selected sample
    local path = url .. "/get-review?user_name=" .. urlencode(sample.user_name) .. "&restaurant_name=" .. urlencode(sample.restaurant_name)
    -- print(path) -- (optional) uncomment me to print the URL query!
    local headers = {}
    return wrk.format(method, path, headers, nil)
end

local function search_reviews()
    local method = "GET"

    -- Choose a random sample to get from the loaded review data
    local random_index = math.random(1, #samples)
    local sample = samples[random_index]

    -- Construct the path using the selected sample
    local path = url .. "/search-reviews?restaurant_name=" .. urlencode(sample.restaurant_name)
    -- print(path) -- (optional) uncomment me to print the URL query!
    local headers = {}
    return wrk.format(method, path, headers, nil)
end

request = function()
    cur_time = math.floor(socket.gettime())
    local post_review_ratio = 0.6
    local get_review_ratio = 0.15
    local search_reviews_ratio = 0.25

    local coin = math.random()
    if coin < post_review_ratio then
        return post_review(url)
    elseif coin < post_review_ratio + get_review_ratio then
        return get_review(url)
    else
        return search_reviews(url)
    end
end