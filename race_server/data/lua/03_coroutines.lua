-- Topic: Coroutines
local function producer()
    return coroutine.wrap(function()
        for i = 1, 5 do
            coroutine.yield("item_" .. i)
        end
    end)
end

local next = producer()
for val in next do
    print(val)
end
