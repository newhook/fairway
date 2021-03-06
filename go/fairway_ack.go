package fairway

func FairwayAck() string {
	return `
local namespace = KEYS[1];

local k = function (queue, subkey)
  return namespace .. queue .. ':' .. subkey;
end

local queue   = ARGV[1];
local facet   = ARGV[2];
local message = ARGV[3];

local active_facets  = k(queue, 'active_facets');
local round_robin    = k(queue, 'facet_queue');
local inflight       = k(queue, 'inflight');
local messages       = k(queue, facet);
local inflight_total = k(queue, facet .. ':inflight');
local inflight_limit = k(queue, 'limit');

local removed = tonumber(redis.call('zrem', inflight, message))

if removed > 0 then
  local current = tonumber(redis.call('decr', inflight_total));
  local max     = tonumber(redis.call('get', inflight_limit)) or 0;
  local length  = redis.call('llen', messages);

  -- If we decremented current to one less than the max,
  -- then we were at the limit, so re-add the facet to
  -- the round robin queue to allow additional messages
  -- to be dequeued.
  if max > 0 and length > 0 and current + 1 == max then
    redis.call('lpush', round_robin, facet);
  end

  if current == 0 then
    redis.call('del', inflight_total);

    -- If we have no currently inflight messages
    -- and there are no more messages in the
    -- facet, remove facet from list of active
    -- facets.
    if length == 0 then
      redis.call('srem', active_facets, facet);
    end
  end
end

return removed
`
}
