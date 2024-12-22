local wrk = require "wrk"

local counter = (function()
    local count = 0
    return function()
        count = count + 1
        return count
    end
end)()

wrk.method = "GET"
wrk.headers["Content-Type"] = "application/json"
wrk.headers["Cookie"] = "email=6c79a5058762eb6b8a4abd068ae2965e;csrf=e5996a98a36a54c852af49abb0989cf3"
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
    
    return wrk.format("GET", "/api/messages?email=sonya@giga-mail.ru")
end

wrk.thread = function()
    wrk.delay = function()
        return 10
    end
end