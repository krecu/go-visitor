local charset = {}

function randInteger(a, b)
    if a == nil and b == nil then
        return math.random(0, 100)
    end
    if b == nil then
        return math.random(a)
    end
    return math.random(a, b)
end

function randIpv4()
    local str = ''
    for i=1, 4 do
        str = str .. randInteger(0, 255)
        if i ~= 4 then str = str .. '.' end
    end
    return str
end

function randUa()
    lines = {}
    for line in io.lines("/Users/kretsu/Work/Go/src/github.com/krecu/go-visitor/tests/fixtures/user-agent.txt") do
        table.insert(lines, line)
    end

    return lines[math.random(#lines)]
end

request = function()

    ip = randIpv4()
    ua = randUa()

    wrk.method = "POST"
    wrk.body   = "{\"id\":\"" .. randInteger(10000, 100000000) .. "\",\"ip\":\"" .. ip .. "\",\"ua\":\"" .. ua .. "\"}"
    wrk.headers["Content-Type"] = "application/json"

    path = "/api/visitor"

    return wrk.format(nil, path)
end