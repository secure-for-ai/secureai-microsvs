package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

type RedisConf struct {
	Addrs []string `json:"Addrs"`
	PW    string   `json:"PW"`
}

// RedisClusterClient struct
type RedisClient struct {
	rdb redis.UniversalClient
}

func NewRedisClient(conf RedisConf) (client *RedisClient, err error) {
	client = &RedisClient{}
	rdb := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:    conf.Addrs,
		Password: conf.PW,
		DB:       0, // use default DB
	})
	rdb.Context()
	pong, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	log.Println(pong, "Ping Redis Success!")

	client.rdb = rdb

	return client, err
}

func (c *RedisClient) Close() error {
	return c.rdb.Close()
}

func (c *RedisClient) GetClient() redis.UniversalClient {
	return c.rdb
}

// Get the value of key. If the key does not exist the special value
// nil is returned. An error is returned if the value stored at key
// is not a string, because GET only handles string values.
//
// Return value
//
// Bulk string reply: the value of key, or nil when key does not exist.
// See https://redis.io/commands/get
func (c *RedisClient) Get(ctx context.Context, key string) (string, error) {
	return c.rdb.Get(ctx, key).Result()
}

// Get the value of key. If the key does not exist the special value
// nil is returned. An error is returned if the value stored at key
// is not a string, because GET only handles string values.
//
// Return value
//
// Byte Array reply: the value of key, or nil when key does not exist.
// See https://redis.io/commands/get
func (c *RedisClient) GetBytes(ctx context.Context, key string) ([]byte, error) {
	return c.rdb.Get(ctx, key).Bytes()
}

// Get the value of key. If the key does not exist the special value
// nil is returned. An error is returned if the value stored at key
// is not a string, because GET only handles string values.
//
// Return value
//
// Int reply: the value of key, or nil when key does not exist.
// See https://redis.io/commands/get
func (c *RedisClient) GetInt(ctx context.Context, key string) (int, error) {
	return c.rdb.Get(ctx, key).Int()
}

// Get the value of key. If the key does not exist the special value
// nil is returned. An error is returned if the value stored at key
// is not a string, because GET only handles string values.
//
// Return value
//
// Int64 reply: the value of key, or nil when key does not exist.
// See https://redis.io/commands/get
func (c *RedisClient) GetInt64(ctx context.Context, key string) (int64, error) {
	return c.rdb.Get(ctx, key).Int64()
}

// Get the value of key. If the key does not exist the special value
// nil is returned. An error is returned if the value stored at key
// is not a string, because GET only handles string values.
//
// Return value
//
// Uint64 reply: the value of key, or nil when key does not exist.
// See https://redis.io/commands/get
func (c *RedisClient) GetUint64(ctx context.Context, key string) (uint64, error) {
	return c.rdb.Get(ctx, key).Uint64()
}

// Get the value of key. If the key does not exist the special value
// nil is returned. An error is returned if the value stored at key
// is not a string, because GET only handles string values.
//
// Return value
//
// Float32 reply: the value of key, or nil when key does not exist.
// See https://redis.io/commands/get
func (c *RedisClient) GetFloat32(ctx context.Context, key string) (float32, error) {
	return c.rdb.Get(ctx, key).Float32()
}

// Get the value of key. If the key does not exist the special value
// nil is returned. An error is returned if the value stored at key
// is not a string, because GET only handles string values.
//
// Return value
//
// Float64 reply: the value of key, or nil when key does not exist.
// See https://redis.io/commands/get
func (c *RedisClient) GetFloat64(ctx context.Context, key string) (float64, error) {
	return c.rdb.Get(ctx, key).Float64()
}

// Get the value of key. If the key does not exist the special value
// nil is returned. An error is returned if the value stored at key
// is not a string, because GET only handles string values.
//
// Return value
//
// Time reply: the value of key, or nil when key does not exist.
//
// See https://redis.io/commands/get
func (c *RedisClient) GetTime(ctx context.Context, key string) (time.Time, error) {
	return c.rdb.Get(ctx, key).Time()
}

// Set key to hold the string value. If key already holds a value,
// it is overwritten, regardless of its type. Any previous time to
// live associated with the key is discarded on successful SET operation.
//
//   - call classic SET command if expiration is 0
//   - call SETEX or SETPX command depending on expiration precision
//   - call SETEX if expiration is presented in seconds
//   - call SETPX if expiration is presented in million seconds
//
// Return value
//
// Simple string reply: OK if SET was executed correctly.
// Null reply: a Null Bulk Reply is returned if the SET
// operation was not performed because the user specified
// the NX or XX option but the condition was not met.
//
// See https://redis.io/commands/set
func (c *RedisClient) Set(ctx context.Context, key string,
	value interface{}, expiration time.Duration) (string, error) {
	return c.rdb.Set(ctx, key, value, expiration).Result()
}

// Set key to hold string value if key does not exist. In that case,
// it is equal to SET. When key already holds a value, no operation
// is performed. SETNX is short for "SET if Not eXists".
//
// Return value
//
// Integer reply, specifically:
//   - 1 if the key was set
//   - 0 if the key was not set
//
// See https://redis.io/commands/setnx
func (c *RedisClient) SetNX(ctx context.Context, key string,
	value interface{}, expiration time.Duration) (bool, error) {
	return c.rdb.SetNX(ctx, key, value, expiration).Result()
}

// Increments the number stored at key by one. If the key does not exist,
// it is set to 0 before performing the operation. An error is returned if
// the key contains a value of the wrong type or contains a string that
// can not be represented as integer. This operation is limited to 64 bit
// signed integers.
//
// Note: this is a string operation because Redis does not have a dedicated
// integer type. The string stored at the key is interpreted as a base-10
// 64 bit signed integer to execute the operation.
//
// Redis stores integers in their integer representation, so for string
// values that actually hold an integer, there is no overhead for storing
// the string representation of the integer.
//
// Return value
//
// Integer reply: the value of key after the increment
//
// See https://redis.io/commands/incr
func (c *RedisClient) Incr(ctx context.Context, key string) (int64, error) {
	return c.rdb.Incr(ctx, key).Result()
}

// Increments the number stored at key by increment. If the key does not exist,
// it is set to 0 before performing the operation. An error is returned if the
// key contains a value of the wrong type or contains a string that can not be
// represented as integer. This operation is limited to 64 bit signed integers.
//
// See INCR for extra information on increment/decrement operations.
//
// Return value
//
// Integer reply: the value of key after the increment
//
// See https://redis.io/commands/incrby
func (c *RedisClient) IncrBy(ctx context.Context, key string, value int64) (int64, error) {
	return c.rdb.IncrBy(ctx, key, value).Result()
}

// Increment the string representing a floating point number stored at key by
// the specified increment. By using a negative increment value, the result is
// that the value stored at the key is decremented (by the obvious properties
// of addition). If the key does not exist, it is set to 0 before performing the
// operation. An error is returned if one of the following conditions occur:
//
// Return value
//
// Bulk string reply: the value of key after the increment.
//
// See https://redis.io/commands/incrbyfloat
func (c *RedisClient) IncrByFloat(ctx context.Context, key string, value float64) (float64, error) {
	return c.rdb.IncrByFloat(ctx, key, value).Result()
}

// Returns the values of all specified keys. For every key that does not hold
// a string value or does not exist, the special value nil is returned. Because
// of this, the operation never fails.
//
// Return value
//
// Array reply: list of values at the specified keys.
//
// See https://redis.io/commands/mget
func (c *RedisClient) MGet(ctx context.Context, key ...string) ([]interface{}, error) {
	return c.rdb.MGet(ctx, key...).Result()
}

// Sets the given keys to their respective values. MSET replaces existing values
// with new values, just as regular SET. See MSETNX if you don't want to overwrite
// existing values.
//
// MSET is atomic, so all given keys are set at once. It is not possible for clients
// to see that some of the keys were updated while others are unchanged.
//
// Return value
//
// Simple string reply: always OK since MSET can't fail.
//
// See https://redis.io/commands/mset
//
// Example: MSet is like Set but accepts multiple values:
//   - MSet("key1", "value1", "key2", "value2")
//   - MSet([]string{"key1", "value1", "key2", "value2"})
//   - MSet(map[string]interface{}{"key1": "value1", "key2": "value2"})
func (c *RedisClient) MSet(ctx context.Context, values ...interface{}) (string, error) {
	return c.rdb.MSet(ctx, values...).Result()
}

// Sets the given keys to their respective values. MSETNX will not perform any
// operation at all even if just a single key already exists.
//
// Return value
// Integer reply, specifically:
//
//   - 1 if the all the keys were set.
//   - 0 if no key was set (at least one key already existed).
//
// See https://redis.io/commands/msetnx
func (c *RedisClient) MSETNX(ctx context.Context, values ...interface{}) (bool, error) {
	return c.rdb.MSetNX(ctx, values...).Result()
}

// Returns all keys matching pattern.
//
// Return value
// Array reply: list of keys matching pattern.
//
// See https://redis.io/commands/keys
func (c *RedisClient) KEYS(ctx context.Context, pattern string) ([]string, error) {
	return c.rdb.Keys(ctx, pattern).Result()
}

// Removes the specified keys. A key is ignored if it does not exist.
//
// Return value
// Integer reply: The number of keys that were removed.
//
// See https://redis.io/commands/del
func (c *RedisClient) Del(ctx context.Context, keys ...string) (int64, error) {
	return c.rdb.Del(ctx, keys...).Result()
}

// Removes the specified fields from the hash stored at key. Specified fields
// that do not exist within this hash are ignored. If key does not exist, it
// is treated as an empty hash and this command returns 0.
//
// Return value
//
// Integer reply: the number of fields that were removed from the hash, not
// including specified but non existing fields.
//
// See https://redis.io/commands/hdel
func (c *RedisClient) HDel(ctx context.Context, key string, fields ...string) (int64, error) {
	return c.rdb.HDel(ctx, key, fields...).Result()
}

// Returns if field is an existing field in the hash stored at key.
//
// Return value
//
// Integer reply, specifically:
//  - 1 if the hash contains field.
//  - 0 if the hash does not contain field, or key does not exist.
//
// See https://redis.io/commands/hexists
func (c *RedisClient) HExists(ctx context.Context, key, field string) (bool, error) {
	return c.rdb.HExists(ctx, key, field).Result()
}

// Returns the value associated with field in the hash stored at key.
//
// Return value
//
// Bulk string reply: the value associated with field, or nil when field is
// not present in the hash or key does not exist.
//
// See https://redis.io/commands/hget
func (c *RedisClient) HGet(ctx context.Context, key, field string) (string, error) {
	return c.rdb.HGet(ctx, key, field).Result()
}

// Returns the value associated with field in the hash stored at key.
//
// Return value
//
// Byte Array reply: the value associated with field, or nil when field is
// not present in the hash or key does not exist.
//
// See https://redis.io/commands/hget
func (c *RedisClient) HGetBytes(ctx context.Context, key, field string) ([]byte, error) {
	return c.rdb.HGet(ctx, key, field).Bytes()
}

// Returns the value associated with field in the hash stored at key.
//
// Return value
//
// Int reply: the value associated with field, or nil when field is
// not present in the hash or key does not exist.
//
// See https://redis.io/commands/hget
func (c *RedisClient) HGetInt(ctx context.Context, key, field string) (int, error) {
	return c.rdb.HGet(ctx, key, field).Int()
}

// Returns the value associated with field in the hash stored at key.
//
// Return value
//
// Int64 reply: the value associated with field, or nil when field is
// not present in the hash or key does not exist.
//
// See https://redis.io/commands/hget
func (c *RedisClient) HGetIn64(ctx context.Context, key, field string) (int64, error) {
	return c.rdb.HGet(ctx, key, field).Int64()
}

// Returns the value associated with field in the hash stored at key.
//
// Return value
//
// Uint64 reply: the value associated with field, or nil when field is
// not present in the hash or key does not exist.
//
// See https://redis.io/commands/hget
func (c *RedisClient) HGetUin64(ctx context.Context, key, field string) (uint64, error) {
	return c.rdb.HGet(ctx, key, field).Uint64()
}

// Returns the value associated with field in the hash stored at key.
//
// Return value
//
// Float32 reply: the value associated with field, or nil when field is
// not present in the hash or key does not exist.
//
// See https://redis.io/commands/hget
func (c *RedisClient) HGetFloat32(ctx context.Context, key, field string) (float32, error) {
	return c.rdb.HGet(ctx, key, field).Float32()
}

// Returns the value associated with field in the hash stored at key.
//
// Return value
//
// Float64 reply: the value associated with field, or nil when field is
// not present in the hash or key does not exist.
//
// See https://redis.io/commands/hget
func (c *RedisClient) HGetFloat64(ctx context.Context, key, field string) (float64, error) {
	return c.rdb.HGet(ctx, key, field).Float64()
}

// Returns the value associated with field in the hash stored at key.
//
// Return value
//
// Time reply: the value associated with field, or nil when field is
// not present in the hash or key does not exist.
//
// See https://redis.io/commands/hget
func (c *RedisClient) HGetTime(ctx context.Context, key, field string) (time.Time, error) {
	return c.rdb.HGet(ctx, key, field).Time()
}

// Returns all fields and values of the hash stored at key. In the
// returned value, every field name is followed by its value, so the
// length of the reply is twice the size of the hash.
//
// Return value
//
// Array reply: list of fields and their values stored in the hash,
// or an empty list when key does not exist.
//
// See https://redis.io/commands/hgetall
func (c *RedisClient) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return c.rdb.HGetAll(ctx, key).Result()
}

// Increments the number stored at field in the hash stored at key by increment.
// If key does not exist, a new key holding a hash is created. If field does not
// exist the value is set to 0 before the operation is performed.
//
// The range of values supported by HINCRBY is limited to 64 bit signed integers.
//
// Return value
//
// Integer reply: the value at field after the increment operation.
//
// See https://redis.io/commands/hincrby
func (c *RedisClient) HIncrBy(ctx context.Context, key, field string, incr int64) (int64, error) {
	return c.rdb.HIncrBy(ctx, key, field, incr).Result()
}

// Increment the specified field of a hash stored at key, and representing a floating
// point number, by the specified increment. If the increment value is negative,
// the result is to have the hash field value decremented instead of incremented.
//
// If the field does not exist, it is set to 0 before performing the operation.
// An error is returned if one of the following conditions occur:
//
//   - The field contains a value of the wrong type (not a string).
//   - The current field content or the specified increment are not parsable as a double
//     precision floating point number.
//
// Return value
//
// Bulk string reply: the value of field after the increment.
//
// See https://redis.io/commands/hincrbyfloat
func (c *RedisClient) HIncrByFloat(ctx context.Context, key, field string, incr float64) (float64, error) {
	return c.rdb.HIncrByFloat(ctx, key, field, incr).Result()
}

// Returns all field names in the hash stored at key.
//
// Return value
//
// Array reply: list of fields in the hash, or an empty list when key does not exist.
//
// See https://redis.io/commands/hkeys
func (c *RedisClient) HKeys(ctx context.Context, key string) ([]string, error) {
	return c.rdb.HKeys(ctx, key).Result()
}

// Returns the number of fields contained in the hash stored at key.
//
// Return value
//
// Integer reply: number of fields in the hash, or 0 when key does not exist.
//
// See https://redis.io/commands/hlen
func (c *RedisClient) HLen(ctx context.Context, key string) (int64, error) {
	return c.rdb.HLen(ctx, key).Result()
}

// Returns the values associated with the specified fields in the hash stored at key.
//
// For every field that does not exist in the hash, a nil value is returned. Because
// non-existing keys are treated as empty hashes, running HMGET against a non-existing
// key will return a list of nil values.
//
// Return value
//
// Array reply: list of values associated with the given fields, in the same order as they are requested.
// See https://redis.io/commands/hmget
func (c *RedisClient) HMGet(ctx context.Context, key string, fields ...string) ([]interface{}, error) {
	return c.rdb.HMGet(ctx, key, fields...).Result()
}

// Sets field in the hash stored at key to value. If key does not exist, a new key
// holding a hash is created. If field already exists in the hash, it is overwritten.
//
// As of Redis 4.0.0, HSET is variadic and allows for multiple field/value pairs.
//
// Return value
//
// Integer reply: The number of fields that were added.
//
// Example:
// HSet accepts values in following formats:
//   - HSet("myhash", "key1", "value1", "key2", "value2")
//   - HSet("myhash", []string{"key1", "value1", "key2", "value2"})
//   - HSet("myhash", map[string]interface{}{"key1": "value1", "key2": "value2"})
//
// Note that it requires Redis v4 for multiple field/value pairs support.
//
// See https://redis.io/commands/hset
func (c *RedisClient) HSet(ctx context.Context, key string, fields ...interface{}) (int64, error) {
	return c.rdb.HSet(ctx, key, fields...).Result()
}

// [Deprecated] As per Redis 4.0.0, HMSET is considered deprecated. Please use HSET in new code.
func (c *RedisClient) HMSet(ctx context.Context, key string, fields ...interface{}) (bool, error) {
	return c.rdb.HMSet(ctx, key, fields...).Result()
}

// Sets field in the hash stored at key to value, only if field does not yet exist.
// If key does not exist, a new key holding a hash is created. If field already
// exists, this operation has no effect.
//
// Return value
//
// Integer reply, specifically:
//   - 1 if field is a new field in the hash and value was set.
//   - 0 if field already exists in the hash and no operation was performed.
//
// See https://redis.io/commands/hsetnx
func (c *RedisClient) HSetNX(ctx context.Context, key, field string, value interface{}) (bool, error) {
	return c.rdb.HSetNX(ctx, key, field, value).Result()
}

// Returns all values in the hash stored at key.
//
// Return value
//
// Array reply: list of values in the hash, or an empty list when key does not exist.
//
// See https://redis.io/commands/hvals
func (c *RedisClient) HVals(ctx context.Context, key string) ([]string, error) {
	return c.rdb.HVals(ctx, key).Result()
}
