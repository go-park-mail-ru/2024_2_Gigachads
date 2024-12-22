# Домашнее задание 3
## Заполнение базы  
Для заполнения базы мы написали вот такой скрипт post.lua
````
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
    
    wrk.requests = 100000
    
    wrk.thread = function()
        local titleRand = randomString(10)
        local descriptionRand = randomString(200)
        wrk.body = json.encode({
        parentID = 0,
        recipient = "sonya@giga-mail.ru",
        title = titleRand,
        description = descriptionRand
    })
    wrk.connections = 10 -- Количество одновременных подключений на поток
    end
```` 
В этом скрипте мы задаем тип запроса - POST, передаем в теле запроса соответствующий нашему API json со сгенерированными названием и описанием письма, а также устанавливаем ограничение в 10 одновременных подключений на один поток  
Команда для запуска:  
``
    wrk -t5 -c100 -d600s -s post.lua https://giga-mail.ru
``  

## Чтение из базы
Для чтения из базы используем скрипт get.lua
````
    local wrk = require "wrk"

    wrk.method = "GET"
    wrk.headers["Content-Type"] = "application/json"
    wrk.headers["Cookie"] = "email=b7698ba355803ae8ac44993d9283e47"

    wrk.requests = 100000
````
В этом скрипте мы задаем тип запроса - GET, передаем в заголовке запроса email, который мы используем для авторизации, и устанавливаем ограничение в 10 одновременных подключений на один поток  
Команда для запуска:  
``
    wrk -t5 -c100 -d600s -s get.lua https://giga-mail.ru
``