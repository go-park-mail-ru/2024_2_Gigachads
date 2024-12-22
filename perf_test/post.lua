local wrk = require "wrk"

local counter = (function()
    local count = 0
    return function()
        count = count + 1
        return count
    end
end)()

local function randomString(length)
    local chars = 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789'
    local result = ''
    for i = 1, length do
        local index = math.random(1, #chars)
        result = result .. chars:sub(index, index)
    end
    return result
end

wrk.method = "POST"
wrk.headers["Content-Type"] = "application/json"
wrk.headers["Cookie"] = "email=76c79a5058762eb6b8a4abd068ae2965e; csrf=e5996a98a36a54c852af49abb0989cf3"
wrk.headers["X-CSRF-Token"] = "e5996a98a36a54c852af49abb0989cf3"
wrk.headers["Origin"] = "https://giga-mail.ru"
wrk.headers["Referer"] = "https://giga-mail.ru/mail"
wrk.headers["Accept"] = "application/json"
wrk.headers["User-Agent"] = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36"
wrk.headers["Connection"] = "keep-alive"

request = function()
    local current = counter()
    if current > 100000 then
        return nil
    end
    
    local titleRand = randomString(10)
    local descriptionRand = randomString(200)
    
    return wrk.format("POST", wrk.path, wrk.headers, 
        string.format('{"type":"email","parentID":0,"recipient":"sonya@giga-mail.ru","title":"%s","description":"%s"}',
            titleRand, descriptionRand))
end

wrk.thread = function()
    math.randomseed(os.time() + wrk.thread:get("id"))
    wrk.delay = function()
        return 10
    end
end