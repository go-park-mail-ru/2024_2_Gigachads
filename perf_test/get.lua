local wrk = require "wrk"

wrk.method = "GET"
wrk.headers["Content-Type"] = "application/json"
wrk.headers["Cookie"] = "email=a1133744a67e4ef3f080b144e52ef6c1; csrf=8f4e081af9774aba27489a5e68a1351f"

request = function()
    return wrk.format(nil, "/api/email/30017", nil, nil)
end
