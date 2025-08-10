-- Set a fixed seed for reproducibility (change the seed value as needed)
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

-- Function to sample from a Zipf distribution
function sampleZipf(N, alpha)
    local normalization = 0
    for i = 1, N do
        normalization = normalization + 1 / (i^alpha)
    end

    local u = math.random() -- Generate a random number between 0 and 1
    local cumulativeProb = 0

    for i = 1, N do
        local prob = 1 / ((i^alpha) * normalization)
        cumulativeProb = cumulativeProb + prob

        if u <= cumulativeProb then
            return i -- Return the sampled value
        end
    end

    return N -- Fallback (unlikely to reach this point)
end

-- Zipf distribution parameter (adjust as needed but alpha >= 0). At alpha=0, Zipf collapses to sampling from a uniform distribution, as alpha keeps getting larger the tail keeps shrinking
local alpha = 1.5


-- 1-to-1 hash function for picking a random ID
function hash(i, N)
    -- Ensure that N is a positive integer
    assert(type(N) == "number" and N > 0 and math.floor(N) == N, "N must be a positive integer")

    -- Ensure that i is in the range [1, N]
    assert(type(i) == "number" and i >= 1 and i <= N, "i must be in the range [1, N]")

    -- Choose a prime number for the multiplicative constant
    local prime = 17

    -- Calculate the hashed value using the multiplicative inverse
    local hashed_i = ((i - 1) * prime) % N + 1

    return hashed_i
end

detailCacheCapacity = 10 -- make sure this is the same as the variable in `main.go`
reviewCacheCapacity = 10 -- make sure this is the same as the variable in `main.go`
reservCacheCapacity = 10 -- make sure this is the same as the variable in `main.go`

-- dataset sizes, also the number of possible outcomes for the zipf dist
datasetMultiplier   = 10 -- ensures cache contains at most 10% of the total dataset at any given time
detailDatasetSize   = detailCacheCapacity * datasetMultiplier
reviewDatasetSize   = reviewCacheCapacity * datasetMultiplier
reservDatasetSize   = reservCacheCapacity * datasetMultiplier

local function post_detail()
    local method = "GET"

    -- Choose a random "sample" from detail dataset using zipf distribution
    local rand_id = hash(sampleZipf(detailDatasetSize, alpha), detailDatasetSize)

    local restaurant_name = urlencode("restaurant" .. tostring(rand_id))
    local location = urlencode("location" .. tostring(rand_id))
    local style = urlencode("style" .. tostring(rand_id))
    local capacity = tostring(math.random(40, 250))
    -- print(restaurant_name, location, style, capacity)

    -- Construct the path using the selected sample
    local path = url .. "/post-detail?restaurant_name=" .. restaurant_name .. "&location=" .. location .. "&style=" .. style .. "&capacity=" .. capacity
    -- print(path) -- (optional) uncomment me to print the URL query!
    local headers = {}
    return wrk.format(method, path, headers, nil)
end

local function get_detail()
    local method = "GET"

    -- Choose a random "sample" from detail dataset using zipf distribution
    local rand_id = hash(sampleZipf(detailDatasetSize, alpha), detailDatasetSize)
    local restaurant_name = urlencode("restaurant" .. tostring(rand_id))

    -- Construct the path using the selected sample
    local path = url .. "/get-detail?restaurant_name=" .. restaurant_name
    -- print(path) -- (optional) uncomment me to print the URL query!
    local headers = {}
    return wrk.format(method, path, headers, nil)
end

local function post_review()
    local method = "GET"

    -- Choose a random "sample" from detail dataset using zipf distribution
    local rand_id = hash(sampleZipf(reviewDatasetSize, alpha), reviewDatasetSize)

    local restaurant_name = urlencode("restaurant" .. tostring(rand_id))
    local user_name = urlencode("user" .. tostring(rand_id))
    local review = urlencode("review" .. tostring(rand_id))
    local rating = urlencode(tostring(math.random(1, 5)))
    -- print(restaurant_name, location, style, capacity)

    -- Construct the path using the selected sample
    local path = url .. "/post-review?restaurant_name=" .. restaurant_name .. "&user_name=" .. user_name .. "&review=" .. review .. "&rating=" .. rating
    -- print(path) -- (optional) uncomment me to print the URL query!
    local headers = {}
    return wrk.format(method, path, headers, nil)
end

local function get_review()
    local method = "GET"

    -- Choose a random "sample" from detail dataset using zipf distribution
    local rand_id = hash(sampleZipf(reviewDatasetSize, alpha), reviewDatasetSize)

    local restaurant_name = urlencode("restaurant" .. tostring(rand_id))
    local user_name = urlencode("user" .. tostring(rand_id))

    -- Construct the path using the selected sample
    local path = url .. "/get-review?user_name=" .. user_name .. "&restaurant_name=" .. restaurant_name
    -- print(path) -- (optional) uncomment me to print the URL query!
    local headers = {}
    return wrk.format(method, path, headers, nil)
end

local function search_reviews()
    local method = "GET"

    -- Choose a random "sample" from detail dataset using zipf distribution
    local rand_id = hash(sampleZipf(reviewDatasetSize, alpha), reviewDatasetSize)
    local restaurant_name = urlencode("restaurant" .. tostring(rand_id))

    -- Construct the path using the selected sample
    local path = url .. "/search-reviews?restaurant_name=" .. restaurant_name
    -- print(path) -- (optional) uncomment me to print the URL query!
    local headers = {}
    return wrk.format(method, path, headers, nil)
end

request = function(requestType)
    cur_time = math.floor(socket.gettime())
    local coin = math.random()
    local coin2 = math.random()

    local detail_ratio = 0.5

    if coin < detail_ratio then
        local post_detail_ratio = 0.2
        if coin2 < post_detail_ratio then
            return post_detail(url)
        else
            return get_detail(url)
        end
    else
        local post_review_ratio = 0.3
        local search_reviews_ratio = 0.4

        if coin2 < post_review_ratio then
            return post_review(url)
        elseif coin2 < post_review_ratio + search_reviews_ratio then
            return search_reviews(url)
        else
            return get_review(url)
        end
    end
end
