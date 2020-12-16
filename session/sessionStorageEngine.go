package session

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/sessions"
	"github.com/secure-for-ai/secureai-microsvs/cache"
	"github.com/secure-for-ai/secureai-microsvs/db"
	"github.com/secure-for-ai/secureai-microsvs/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net"
	"strconv"
	"time"
)

//type Base64IDGenerator struct{
//	Length int
//}
//
//func (r Base64IDGenerator) newSessID(sess *sessions.Session) string {
//	b, _ := util.GenerateRandomKey(r.Length)
//	id := util.Base64Encode(b)
//	sess.Values["sid"] = util.Base64Encode(b)
//	return id
//}
//
//func (r Base64IDGenerator) encodeSessID(sess *sessions.Session) []byte {
//	return []byte(sess.Values["sid"].(string))
//}
//
//func (r Base64IDGenerator) decodeSessID(sess *sessions.Session, data []byte) error {
//	id := string(data)
//	sess.Values["sid"] = id
//	sess.ID = id
//	return nil
//}

type SessIDGenerator interface {
	newNonce(sess *sessions.Session)
	newSessID(sess *sessions.Session) string
	encodeSessID(sess *sessions.Session) []byte
	decodeSessID(sess *sessions.Session, data []byte) error
}

//type SUIDStrGenerator struct{}
//
//func (r SUIDStrGenerator) newNonce(sess *sessions.Session) {
//	nonce, _ := util.GenerateRandomKey(16)
//	sess.Values["nonce"] = nonce
//}
//
//func (r SUIDStrGenerator) newSessID(sess *sessions.Session) string {
//	r.newNonce(sess)
//	sess.Values["sid"] = string(util.GetNowTimestamp())
//	id := sess.Values["uid"].(string) + "_" + sess.Values["sid"].(string)
//	return id
//}
//
//func (r SUIDStrGenerator) encodeSessID(sess *sessions.Session) []byte {
//	return []byte(sess.Values["sid"].(string) + "|" + sess.Values["uid"].(string) +
//		"|" + util.Base64Encode(sess.Values["nonce"].([]byte)))
//}
//
//func (r SUIDStrGenerator) decodeSessID(sess *sessions.Session, data []byte) error {
//	strs := strings.Split(string(data), "|")
//	if len(strs) != 3 {
//		return ErrInvalidCookie
//	}
//	sess.Values["sid"] = strs[0]
//	sess.Values["uid"] = strs[1]
//	nonce, err := util.Base64Decode(strs[2])
//	if err != nil {
//		return err
//	}
//	sess.Values["nonce"] = nonce
//	sess.ID = strs[1] + "_" + strs[0]
//	return nil
//}

type SUIDInt64Generator struct{}

func (r SUIDInt64Generator) newNonce(sess *sessions.Session) {
	nonce, _ := util.GenerateRandomKey(16)
	sess.Values["nonce"] = nonce
}

func (r SUIDInt64Generator) newSessID(sess *sessions.Session) string {
	r.newNonce(sess)
	sess.Values["sid"] = time.Now().Unix()
	id := strconv.FormatInt(sess.Values["uid"].(int64), 10) +
		"_" + strconv.FormatInt(sess.Values["sid"].(int64), 10)
	return id
}

func (r SUIDInt64Generator) encodeSessID(sess *sessions.Session) []byte {
	data := make([]byte, 32)
	binary.LittleEndian.PutUint64(data, uint64(sess.Values["sid"].(int64)))
	binary.LittleEndian.PutUint64(data[8:], uint64(sess.Values["uid"].(int64)))
	copy(data[16:], sess.Values["nonce"].([]byte))
	return data
}

func (r SUIDInt64Generator) decodeSessID(sess *sessions.Session, data []byte) error {
	if len(data) != 32 {
		return ErrInvalidCookie
	}
	sid := int64(binary.LittleEndian.Uint64(data[0:8]))
	uid := int64(binary.LittleEndian.Uint64(data[8:]))
	sess.Values["sid"] = sid
	sess.Values["uid"] = uid
	nonce := make([]byte, 16)
	copy(nonce, data[16:])
	sess.Values["nonce"] = nonce
	sess.ID = strconv.FormatInt(sess.Values["uid"].(int64), 10) +
		"_" + strconv.FormatInt(sess.Values["sid"].(int64), 10)
	return nil
}

//type AbsStoreEngine struct {
//	Prefix string
//	IDGenerator SessIDGenerator
//}
//
//func (r AbsStoreEngine) newSessID(sess *sessions.Session) string {
//	return r.Prefix + r.IDGenerator.newSessID(sess)
//}
//
//func (r AbsStoreEngine) encodeSessID(sess *sessions.Session) []byte {
//	return r.IDGenerator.encodeSessID(sess)
//}
//
//func (r AbsStoreEngine) decodeSessID(sess *sessions.Session, data []byte) error {
//	err := r.IDGenerator.decodeSessID(sess, data)
//	if err != nil {
//		return err
//	}
//	sess.ID = r.Prefix + sess.ID
//	return nil
//}

/* RedisStoreEngine only use redis as storage, session data will lost after reboot */
type RedisStoreEngine struct {
	RedisClient *cache.RedisClient
	Serializer  DataSerializer
	Prefix      string
	IDGenerator SessIDGenerator
}

func (r *RedisStoreEngine) init() error { return nil }

func (r *RedisStoreEngine) load(ctx context.Context, sess *sessions.Session) (bool, error) {
	data, err := r.RedisClient.GetBytes(ctx, sess.ID)
	//fmt.Println("load session id", sess.ID)
	//fmt.Println("load", data)
	//fmt.Println("err", err)
	if err != nil {
		switch err {
		case redis.Nil:
			return false, ErrNil
		default:
			return false, ErrStoreFail
		}
	}
	nonceClient := sess.Values["nonce"].([]byte)
	err = r.Serializer.Deserialize(data, sess)
	if err != nil {
		return false, err
	}
	nonceServer := sess.Values["nonce"].([]byte)
	res := bytes.Compare(nonceClient, nonceServer)
	if res == 0 {
		return true, nil
	}
	return false, ErrNonce
}

func (r *RedisStoreEngine) save(ctx context.Context, sess *sessions.Session, duration time.Duration) error {
	if sess.ID == "" {
		// we use timestamp as session id
		// developer is responsible to assign sess.Values["uid"]
		sess.ID = r.newSessID(sess)
	} else {
		r.IDGenerator.newNonce(sess)
	}

	data, err := r.Serializer.Serialize(sess)
	fmt.Println("save session value", sess.Values, 1000)
	fmt.Println("save session serial", data)
	_, err = r.RedisClient.Set(ctx, sess.ID, data, duration)
	return err
}

func (r *RedisStoreEngine) delete(ctx context.Context, sess *sessions.Session) error {
	_, err := r.RedisClient.Del(ctx, sess.ID)
	return err
}

func (r RedisStoreEngine) newSessID(sess *sessions.Session) string {
	return r.Prefix + r.IDGenerator.newSessID(sess)
}

func (r *RedisStoreEngine) encodeSessID(sess *sessions.Session) []byte {
	return r.IDGenerator.encodeSessID(sess)
}

func (r *RedisStoreEngine) decodeSessID(sess *sessions.Session, data []byte) error {
	err := r.IDGenerator.decodeSessID(sess, data)
	if err != nil {
		return err
	}
	sess.ID = r.Prefix + sess.ID
	return nil
}

/* RedisMongoStoreEngine, use redis as cache and mongo as persistent storage */
type RedisMongoStoreEngine struct {
	RedisClient   *cache.RedisClient
	MongoDBClient *db.MongoDBClient
	Serializer    DataSerializer
	Prefix        string
	IDGenerator   SessIDGenerator
	Table         string
	CacheAge      time.Duration
}

type SessValue struct {
	Sid        int64                  `json:"sid"        bson:"sid"         pg:"sid"`
	Uid        int64                  `json:"uid"        bson:"uid"         pg:"uid"`
	Nonce      []byte                 `json:"nonce"      bson:"nonce"       pg:"nonce"`
	Data       map[string]interface{} `json:"data"       bson:"data"        pg:"data"`
	IP         net.IP                 `json:"ip"         bson:"ip"          pg:"ip"`
	UserAgent  string                 `json:"userAgent"  bson:"user_agent"  pg:"user_agent"`
	CreateTime int64                  `json:"createTime" bson:"create_time" pg:"create_time"`
	UpdateTime int64                  `json:"updateTime" bson:"update_time" pg:"update_time"`
	ExpireTime int64                  `json:"expireTime" bson:"expire_time" pg:"expire_time"`
}

func (r *RedisMongoStoreEngine) init() error { return nil }

func (r *RedisMongoStoreEngine) load(ctx context.Context, sess *sessions.Session) (bool, error) {
	nonceClient := sess.Values["nonce"].([]byte)
	data, err := r.RedisClient.GetBytes(ctx, sess.ID)
	//fmt.Println("load session id", sess.ID)
	//fmt.Println("load", data)
	//fmt.Println("err", err)

	if data == nil {
		fmt.Println("========== cache not hit ===========", err)
		// cache not hit
		filter := bson.M{
			"sid": sess.Values["sid"],
			"uid": sess.Values["uid"],
		}
		var sessValue = &SessValue{}
		errDB := r.MongoDBClient.FindOne(ctx, r.Table, &filter, &sessValue)
		if errDB != nil {
			return false, errDB
		}
		//fmt.Println(errDB)
		//fmt.Println(sessValue)
		//fmt.Printf("%T\n", sessValue.Data)
		if res := bytes.Compare(nonceClient, sessValue.Nonce); res == 0 {
			sess.Values["data"] = sessValue.Data
			sess.Values["ip"] = sessValue.IP
			sess.Values["userAgent"] = sessValue.UserAgent
			sess.Values["createTime"] = sessValue.CreateTime
			sess.Values["updateTime"] = sessValue.UpdateTime
			sess.Values["expireTime"] = sessValue.ExpireTime
			// save the cache in redis
			data, err := r.Serializer.Serialize(sess)
			if err == nil {
				_, err = r.RedisClient.Set(ctx, sess.ID, data, r.CacheAge)
				fmt.Println("*****", err)
			}
			//fmt.Println("save session value", sess.Values, 1000)
			//fmt.Println("save session serial", data)
			return true, nil
		}
		return false, ErrNonce
	} else {
		fmt.Println("========== cache hit ===========", err)
		// cache hit
		err = r.Serializer.Deserialize(data, sess)
		if err != nil {
			return false, err
		}
		nonceServer := sess.Values["nonce"].([]byte)
		res := bytes.Compare(nonceClient, nonceServer)
		if res == 0 {
			return true, nil
		}
		return false, ErrNonce
	}
}

func (r *RedisMongoStoreEngine) save(ctx context.Context, sess *sessions.Session, duration time.Duration) error {
	currentTime := util.GetNowTimestamp()
	expireTime := currentTime + int64(duration)

	if sess.ID == "" {
		// we use timestamp as session id
		// developer is responsible to assign sess.Values["uid"]
		sess.ID = r.newSessID(sess)
		//fmt.Println(sess.ID)
		sess.Values["createTime"] = currentTime
	} else {
		r.IDGenerator.newNonce(sess)
	}

	sess.Values["updateTime"] = currentTime
	sess.Values["expireTime"] = expireTime

	// save the data in mongoDB
	sessValue := SessValue{
		Sid:        sess.Values["sid"].(int64),
		Uid:        sess.Values["uid"].(int64),
		Nonce:      sess.Values["nonce"].([]byte),
		Data:       sess.Values["data"].(map[string]interface{}),
		IP:         sess.Values["ip"].(net.IP),
		UserAgent:  sess.Values["userAgent"].(string),
		CreateTime: sess.Values["createTime"].(int64),
		UpdateTime: currentTime,
		ExpireTime: expireTime,
	}

	opts := options.Update().SetUpsert(true)
	filter := bson.M{
		"sid": sess.Values["sid"],
		"uid": sess.Values["uid"],
	}
	update := bson.M{
		"$set": sessValue,
	}
	result, err := r.MongoDBClient.UpdateOne(ctx, r.Table, &filter, &update, opts)

	if err != nil {
		log.Fatal(err)
		return err
	}

	if result.MatchedCount != 0 {
		fmt.Println("matched and replaced an existing document")
	}

	if result.UpsertedCount != 0 {
		fmt.Printf("inserted a new document with ID %v\n", result.UpsertedID)
	}

	// delete cache
	//_, _ = r.RedisClient.Del(ctx, sess.ID)
	// save the cache in redis
	data, err := r.Serializer.Serialize(sess)
	//fmt.Println("save session value", sess.Values, 1000)
	//fmt.Println("save session serial", data)
	_, err = r.RedisClient.Set(ctx, sess.ID, data, r.CacheAge)
	return nil
}

func (r *RedisMongoStoreEngine) delete(ctx context.Context, sess *sessions.Session) error {
	// Delete the cache first, return error if not success
	_, errCache := r.RedisClient.Del(ctx, sess.ID)
	if errCache != nil {
		return errCache
	}
	query := bson.M{
		"sid": sess.Values["sid"],
		"uid": sess.Values["uid"],
	}
	// delete in the mongodb, return error if not success
	_, errDB := r.MongoDBClient.DeleteOne(ctx, r.Table, &query)
	return errDB
}

func (r RedisMongoStoreEngine) newSessID(sess *sessions.Session) string {
	return r.Prefix + r.IDGenerator.newSessID(sess)
}

func (r *RedisMongoStoreEngine) encodeSessID(sess *sessions.Session) []byte {
	return r.IDGenerator.encodeSessID(sess)
}

func (r *RedisMongoStoreEngine) decodeSessID(sess *sessions.Session, data []byte) error {
	err := r.IDGenerator.decodeSessID(sess, data)
	if err != nil {
		return err
	}
	sess.ID = r.Prefix + sess.ID
	return nil
}

/* RedisPGEngine, use redis as cache and postgres as persistent storage */
type RedisPGStoreEngine struct {
	RedisClient *cache.RedisClient
	PGClient    *db.PGClient
	Serializer  DataSerializer
	Prefix      string
	IDGenerator SessIDGenerator
	Table       string
	CacheAge    time.Duration
}

func (r *RedisPGStoreEngine) init() error { return nil }

func (r *RedisPGStoreEngine) load(ctx context.Context, sess *sessions.Session) (bool, error) {
	nonceClient := sess.Values["nonce"].([]byte)
	data, err := r.RedisClient.GetBytes(ctx, sess.ID)

	if data == nil {
		fmt.Println("========== cache not hit ===========", err)
		// cache not hit
		var sessValue = &SessValue{}
		conn, errDB := r.PGClient.GetConn(ctx)
		if errDB != nil {
			return false, errDB
		}
		defer conn.Release()
		errDB = conn.FindOne(
			ctx, "SELECT * FROM "+r.Table+" WHERE sid=$1 AND uid=$2;", sessValue,
			sess.Values["sid"], sess.Values["uid"])
		if errDB != nil {
			return false, errDB
		}
		fmt.Println(errDB)
		fmt.Println("PG:", sessValue)
		//fmt.Printf("%T\n", sessValue.Data)
		if res := bytes.Compare(nonceClient, sessValue.Nonce); res == 0 {
			sess.Values["data"] = sessValue.Data
			sess.Values["ip"] = sessValue.IP
			sess.Values["userAgent"] = sessValue.UserAgent
			sess.Values["createTime"] = sessValue.CreateTime
			sess.Values["updateTime"] = sessValue.UpdateTime
			sess.Values["expireTime"] = sessValue.ExpireTime
			// save the cache in redis
			data, err := r.Serializer.Serialize(sess)
			if err == nil {
				_, err = r.RedisClient.Set(ctx, sess.ID, data, r.CacheAge)
				fmt.Println("*****", err)
			}
			//fmt.Println("save session value", sess.Values, 1000)
			//fmt.Println("save session serial", data)
			return true, nil
		}
		return false, ErrNonce
	} else {
		fmt.Println("========== cache hit ===========", err)
		// cache hit
		err = r.Serializer.Deserialize(data, sess)
		if err != nil {
			return false, err
		}
		nonceServer := sess.Values["nonce"].([]byte)
		res := bytes.Compare(nonceClient, nonceServer)
		if res == 0 {
			return true, nil
		}
		return false, ErrNonce
	}
	//if err != nil {
	//	switch err {
	//	case redis.Nil:
	//		return nil, ErrNil
	//	default:
	//		return nil, ErrStoreFail
	//	}
	//}
	//return data, nil
}

func (r *RedisPGStoreEngine) save(ctx context.Context, sess *sessions.Session, duration time.Duration) error {
	currentTime := util.GetNowTimestamp()
	expireTime := currentTime + int64(duration)

	if sess.ID == "" {
		// we use timestamp as session id
		// developer is responsible to assign sess.Values["uid"]
		sess.ID = r.newSessID(sess)
		//fmt.Println(sess.ID)
		sess.Values["createTime"] = currentTime
	} else {
		r.IDGenerator.newNonce(sess)
	}

	sess.Values["updateTime"] = currentTime
	sess.Values["expireTime"] = expireTime

	// save the data in mongoDB
	sessValue := SessValue{
		Sid:        sess.Values["sid"].(int64),
		Uid:        sess.Values["uid"].(int64),
		Nonce:      sess.Values["nonce"].([]byte),
		Data:       sess.Values["data"].(map[string]interface{}),
		IP:         sess.Values["ip"].(net.IP),
		UserAgent:  sess.Values["userAgent"].(string),
		CreateTime: sess.Values["createTime"].(int64),
		UpdateTime: currentTime,
		ExpireTime: expireTime,
	}

	// delete in the mongodb, return error if not success
	conn, err := r.PGClient.GetConn(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()
	result, err := conn.Update(
		ctx,
		"INSERT INTO "+r.Table+" (sid, uid, nonce, data, ip, "+
			"user_agent, create_time, update_time, expire_time) VALUES "+
			"($1, $2, $3, $4, $5, $6, $7, $8, $9) "+
			"ON CONFLICT (sid, uid) "+
			"Do UPDATE SET nonce=EXCLUDED.nonce, data=EXCLUDED.data, "+
			"ip=EXCLUDED.ip, user_agent=EXCLUDED.user_agent, "+
			"create_time=EXCLUDED.create_time, update_time=EXCLUDED.update_time,"+
			"expire_time=EXCLUDED.expire_time;",
		sessValue.Sid,
		sessValue.Uid, sessValue.Nonce, sessValue.Data, sessValue.IP,
		sessValue.UserAgent, sessValue.CreateTime, sessValue.UpdateTime,
		sessValue.ExpireTime,
	)
	// there is an error on upsert
	if err != nil {
		log.Fatal(err)
		return err
	}

	fmt.Println("Upsert Session Affected rows: ", result)

	// delete cache
	//_, _ = r.RedisClient.Del(ctx, sess.ID)
	// save the cache in redis
	data, err := r.Serializer.Serialize(sess)
	//fmt.Println("save session value", sess.Values, 1000)
	//fmt.Println("save session serial", data)
	_, err = r.RedisClient.Set(ctx, sess.ID, data, r.CacheAge)
	return nil
}

func (r *RedisPGStoreEngine) delete(ctx context.Context, sess *sessions.Session) error {
	// Delete the cache first, return error if not success
	_, errCache := r.RedisClient.Del(ctx, sess.ID)
	if errCache != nil {
		return errCache
	}

	// delete in the mongodb, return error if not success
	conn, errDB := r.PGClient.GetConn(ctx)
	if errDB != nil {
		return errDB
	}
	defer conn.Release()
	_, errDB = conn.Delete(
		ctx, "DELETE FROM "+r.Table+" WHERE sid=$1 AND uid=$2",
		sess.Values["sid"], sess.Values["uid"])

	return errDB
}

func (r *RedisPGStoreEngine) newSessID(sess *sessions.Session) string {
	return r.Prefix + r.IDGenerator.newSessID(sess)
}

func (r *RedisPGStoreEngine) encodeSessID(sess *sessions.Session) []byte {
	return r.IDGenerator.encodeSessID(sess)
}

func (r *RedisPGStoreEngine) decodeSessID(sess *sessions.Session, data []byte) error {
	err := r.IDGenerator.decodeSessID(sess, data)
	if err != nil {
		return err
	}
	sess.ID = r.Prefix + sess.ID
	return nil
}
