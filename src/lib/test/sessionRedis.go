package main

import (
	"fmt"
	"github.com/gorilla/sessions"
	"net/http"
	"template2/lib/cache"
	"template2/lib/db"
	"template2/lib/session"
)

var (
	redisConf = cache.RedisConf{
		Addrs: []string{"localhost:6379"},
		PW:    "password",
	}
	sessRedisClient, _ = cache.NewRedisClient(redisConf)

	storeConf = session.HybridStoreConf{
		//IdLength:  32,
		//KeyPrefix: "session_",
		Options: &sessions.Options{
			Path:   "/",
			MaxAge: 86400 * 30,
		},
		//IdGenerator:   "base64",
		//Serializer:    "gob",
		CookieHandler: "aes_gcm", // "base64", // "aes_gcm", // "secure",
		// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256), encoded in base64RawUrl
		KeyPairs: []string{
			"pR6kDdHYqNMRO74kUxFiGgpv3A6qKFeCY6IDHxDH8NY",
			"4M3xrW-JjbMYRYDqUVBPNAJKR0LW8ehvm_jcwD6iyT0",
		},
	}

	redisStore = session.RedisStoreEngine{
		RedisClient: sessRedisClient,
		Serializer:  session.GobSerializer{},
		Prefix:      "session_",
		IDGenerator: session.SUIDInt64Generator{},
	}

	mongoConf = db.MongoDBConf{
		Host:        "localhost",
		Port:        "27017",
		DBName:      "gtest",
		User:        "test",
		PW:          "password",
		AdminDBName: "admin",
	}

	mongoClient, _  = db.NewMongoDB(mongoConf)
	redisMongoStore = session.RedisMongoStoreEngine{
		RedisClient:   sessRedisClient,
		MongoDBClient: mongoClient,
		Serializer:    session.GobSerializer{},
		Prefix:        "session_",
		IDGenerator:   session.SUIDInt64Generator{},
		Table:         "test_session",
		CacheAge:      10,
	}

	pgConf = db.PGPoolConf{
		Host:   "localhost",
		Port:   "7000",
		DBName: "test",
		User:   "test",
		PW:     "password",
	}

	pgClient, _  = db.NewPGClient(pgConf)
	redisPGStore = session.RedisPGStoreEngine{
		RedisClient: sessRedisClient,
		PGClient:    pgClient,
		Serializer:  session.GobSerializer{},
		Prefix:      "session_",
		IDGenerator: session.SUIDInt64Generator{},
		Table:       "test_session",
		CacheAge:    10,
	}

	//key   = []byte("super-secret-key")
	//store = sessions.NewCookieStore(key)
	//store = session.NewSessionStore(&redisStore, &storeConf)
	store = session.NewSessionStore(&redisMongoStore, &storeConf)
	//store = session.NewSessionStore(&redisPGStore, &storeConf)
)

func secret(w http.ResponseWriter, r *http.Request) {
	sess, _ := store.Get(r, "cookie-name")

	fmt.Println(sess.Values)
	// Check if user is authenticated
	if data, ok := sess.Values["data"].(map[string]interface{}); !ok || !data["authenticated"].(bool) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Print secret message
	fmt.Fprintln(w, "The cake is a lie!")
}

func login(w http.ResponseWriter, r *http.Request) {
	sess, _ := store.Get(r, "cookie-name")

	// Authentication goes here
	// ...

	// Set user as authenticated
	sess.Values["uid"] = int64(0)
	sess.Values["data"] = map[string]interface{}{"authenticated": true}
	sess.Save(r, w)
}

func logout(w http.ResponseWriter, r *http.Request) {
	sess, _ := store.Get(r, "cookie-name")

	// Revoke users authentication
	sess.Values["data"] = map[string]interface{}{"authenticated": false}
	sess.Options.MaxAge = -1
	sess.Save(r, w)
}

func main() {
	http.HandleFunc("/secret", secret)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)

	http.ListenAndServe(":8080", nil)
}
