local wrk = require "wrk"
local json = require("cjson")

------------генерация рандомной строки
math.randomseed(os.time())  -- Инициализация генератора случайных чисел

local function randomString(length)
    local chars = 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789'
    local result = ''
    for i = 1, length do
        local index = math.random(1, #chars)
        result = result .. chars:sub(index, index)
    end
    return result
end
-------------------------------------

wrk.method = "POST"
wrk.headers["Content-Type"] = "application/json"
wrk.headers["Cookie"] = "email=b7698ba355803ae8ac44993d9283e47;csrf=b12ca9de73d413ec524cfa80fbd9432c"

wrk.requests = 10

wrk.thread = function()
    local titleRand = randomString(10)
    local descriptionRand = randomString(200)
    local body = json.encode({
       parentID = 0,
       recipient = "sonya@giga-mail.ru",
       title = titleRand,
       description = descriptionRand
    })
    wrk.connections = 10 -- Количество одновременных подключений на поток
    return wrk.format(nil, "/api/email", nil, body)
end