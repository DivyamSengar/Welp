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
local file = io.open("./samples/reservation_samples.csv", "r")

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
        table.insert(samples, {user_name = user_name, restaurant_name = restaurant_name, year = year, month = month, day = day})
    end
end

local function make_reservation()
    local method = "GET"

    -- Choose a random sample to post from the loaded detail data
    local random_index = math.random(1, #samples)
    local sample = samples[random_index]

    -- Construct the path using the selected sample
    local path = url .. "/make-reservation?restaurant_name=" .. urlencode(sample.restaurant_name) .. "&user_name=" .. urlencode(sample.user_name) .. "&year=" .. sample.year .. "&month=" .. sample.month .. "&day=" .. sample.day
    -- print(path) -- (optional) uncomment me to print the URL query!
    local headers = {}
    return wrk.format(method, path, headers, nil)
end

local function get_reservation()
    local method = "GET"

    -- Choose a random sample to get from the loaded detail data
    local random_index = math.random(1, #samples)
    local sample = samples[random_index]

    -- Construct the path using the selected sample
    local path = url .. "/get-reservation?user_name=" .. urlencode(sample.user_name)
    -- print(path) -- (optional) uncomment me to print the URL query!
    local headers = {}
    return wrk.format(method, path, headers, nil)
end

local function most_popular()
    local method = "GET"

    -- Construct the path using the selected sample
    local topk = tostring(math.random(1, 10))
    local path = url .. "/most-popular?topk=" .. topk
    -- print(path) -- (optional) uncomment me to print the URL query!
    local headers = {}
    return wrk.format(method, path, headers, nil)
end



request = function()
    cur_time = math.floor(socket.gettime())
    local make_reservation_ratio  = 0.6
    local get_reservation_ratio   = 0.1
    local most_popular_ratio = 0.3
  
    local coin = math.random()
    if coin < make_reservation_ratio then
        return make_reservation(url)
    elseif coin < make_reservation_ratio + get_reservation_ratio then
        return get_reservation(url)
    else 
        return most_popular(url)
    end
end
