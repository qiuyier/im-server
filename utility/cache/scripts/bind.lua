local hashKey = KEYS[1]
local setKey = KEYS[2]
local field = ARGV[1]
local value = ARGV[2]

local hSetResult = redis.call('HSET', hashKey, field, value)
if hSetResult == 0 then
    return 0
end

local sAddResult = redis.call('SADD', setKey, field)
if sAddResult == 0 then
    return 0
end

return 1