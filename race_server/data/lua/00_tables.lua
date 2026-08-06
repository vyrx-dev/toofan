-- Topic: Tables and Functions
local config = {
    host = "localhost",
    port = 8080,
    debug = true,
}

local function getAddr(cfg)
    return cfg.host .. ":" .. tostring(cfg.port)
end

print(getAddr(config))

local fruits = { "apple", "banana", "cherry" }

for i, fruit in ipairs(fruits) do
    print(i, fruit)
end

for key, val in pairs(config) do
    print(key .. " = " .. tostring(val))
end
