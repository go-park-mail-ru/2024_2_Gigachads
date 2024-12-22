# Домашнее задание 3
## Заполнение базы  
1. Для заполнения базы мы написали вот такой скрипт post.lua
````
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
    if current > 100000 then
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
    Req/Sec    27.25      1.49k   12.53k    88.72%
  49584 requests in 10.00m, 49.58MB read
  Socket errors: connect 0, read 0, write 0, timeout 55
  Non-2xx or 3xx responses: 424
Requests/sec:    83.09
Transfer/sec:    84.63KB
````

3. Документация результатов и анализ значений:

Latency и Req/Sec показывают параметры функции Гаусса

Requests/sec "RPC" (запросы в секунду): показатель составляет 83.09, что является приемлемымм значением.

Transfer/sec (скорость передачи данных): показатель составляет 84.63 КБ/с, что является довольно низким значением.

## Чтение из базы
1. Для чтения из базы используем скрипт get.lua
````
local wrk = require "wrk"

wrk.method = "GET"
wrk.headers["Content-Type"] = "application/json"
wrk.headers["Cookie"] = "email=a1133744a67e4ef3f080b144e52ef6c1; csrf=8f4e081af9774aba27489a5e68a1351f"

request = function()
    return wrk.format(nil, "/api/email/30017", nil, nil)
end
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
    Req/Sec    84.62      2.22k   27.85k    96.62%
  69972 requests in 10.00m, 50.19MB read
  Socket errors: connect 0, read 0, write 0, timeout 76
  Non-2xx or 3xx responses: 973
Requests/sec:   116.62
Transfer/sec:    85.67KB

````

3. Документация результатов и анализ значений:

Latency и Req/Sec показывают параметры функции Гаусса

Requests/sec "RPC" (запросы в секунду): показатель составляет 116.62, что является приемлемым значением.

Transfer/sec (скорость передачи данных): показатель составляет 85.67 КБ/с, что является довольно низким значением.

## Выводы

Для оптимизации можно принять следующие действия:
- добавить пагинацию в страницу папки и грузить письма постранично, а не сразу все.
- для прорисовки содержимого папки отправлять на фронт только начало текста каждого письма, а весь текст уже при открытии конкретного письма
