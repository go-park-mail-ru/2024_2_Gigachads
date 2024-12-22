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
wrk.headers["Cookie"] = "email=4a070bdf28b1b7a38ff1d7283164fa18; csrf=e0e43a3881528791c0c4f546d438c952"
wrk.headers["Accept"] = "application/json"

request = function()
    local current = counter()
    if current > 10 then
        return nil
    end
    
    local titleRand = randomString(10)
    local descriptionRand = randomString(200)
    return wrk.format(nil, "/api/email", nil,
        string.format('{"type":"email","parentID":0,"recipient":"sonya2@giga-mail.ru","title":"%s","description":"%s"}',
            titleRand, descriptionRand))
end

wrk.thread = function()
    math.randomseed(os.time() + wrk.thread:get("id"))
    wrk.delay = function()
        return 10
    end
end