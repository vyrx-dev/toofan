-- Topic: String and File I/O
local function readLines(path)
    local lines = {}
    local file = io.open(path, "r")
    if not file then
        return nil, "cannot open " .. path
    end
    for line in file:lines() do
        lines[#lines + 1] = line
    end
    file:close()
    return lines
end

local data, err = readLines("input.txt")
if not data then
    print("Error: " .. err)
    os.exit(1)
end

for _, line in ipairs(data) do
    local word = string.match(line, "^(%S+)")
    if word then
        print(string.format("-> %s", word))
    end
end
