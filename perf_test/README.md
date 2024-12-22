# Домашнее задание 3
## Заполнение базы  
1. Для заполнения базы мы написали вот такой скрипт post.lua
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

2. Вывод программы:

````
Running 10m test @ https://giga-mail.ru
  5 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   110.65ms  332.35ms   2.00s    90.30%
    Req/Sec   527.25      1.49k   12.53k    88.72%
  499999 requests in 10.00m, 225.43MB read
  Socket errors: connect 0, read 0, write 325353, timeout 12655
  Non-2xx or 3xx responses: 481424
Requests/sec:    833.09
Transfer/sec:    384.63KB
````

## Чтение из базы
1. Для чтения из базы используем скрипт get.lua
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

2. Вывод программы:

````
Running 10m test @ https://giga-mail.ru
  5 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   176.19ms  404.66ms   2.00s    86.35%
    Req/Sec   384.62      2.22k   27.85k    96.62%
  610068 requests in 10.00m, 284.62MB read
  Socket errors: connect 0, read 0, write 0, timeout 14376
  Non-2xx or 3xx responses: 584973
Requests/sec:   1016.62
Transfer/sec:    485.67KB

````
