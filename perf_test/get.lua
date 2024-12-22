local wrk = require "wrk"

wrk.method = "GET"
wrk.headers["Content-Type"] = "application/json"
wrk.headers["Cookie"] = "email=b7698ba355803ae8ac44993d9283e47;csrf=b12ca9de73d413ec524cfa80fbd9432c"

wrk.requests = 100000

wrk.thread = function()
    wrk.path = "/api/messages?email=sonya@giga-mail.ru"
    wrk.connections = 10
end