local hashKey = KEYS[1]
local setKey = KEYS[2]
local field = ARGV[1]

local hDelResult = redis.call('HDEL', hashKey, field)
if hDelResult == 0 then
    return 0
end

local sRemResult = redis.call('SREM', setKey, field)
if sRemResult == 0 then
    return 0
end

return 1