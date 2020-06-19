package main

import (
	"fmt"
	"github.com/gorilla/sessions"
	"net/http"
	"template2/lib/cache"
	"template2/lib/session"
)

var (
	redisConf = cache.RedisConf{
		Addrs: []string{"localhost:6379"},
		PW:    "password",
	}
	sessRedisClient, _ = cache.NewRedisClient(redisConf)

	storeConf = session.HybridStoreConf{
		IdLength:  32,
		KeyPrefix: "session_",
		Options: &sessions.Options{
			Path:   "/",
			MaxAge: 86400 * 30,
		},
		IdGenerator:   "base64",
		Serializer:    "gob",
		CookieHandler: "standard", // "secure",
		// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256), encoded in base64RawUrl
		KeyPairs: []string{
			"pR6kDdHYqNMRO74kUxFiGgpv3A6qKFeCY6IDHxDH8NY",
			"4M3xrW-JjbMYRYDqUVBPNAJKR0LW8ehvm_jcwD6iyT0",
		},
	}
	storage = session.RedisStoreEngine{
		RedisClient: sessRedisClient,
	}
	//key   = []byte("super-secret-key")
	//store = sessions.NewCookieStore(key)
	store = session.NewSessionStore(&storage, &storeConf)
)

func secret(w http.ResponseWriter, r *http.Request) {
	sess, _ := store.Get(r, "cookie-name")

	fmt.Println(sess.Values)
	// Check if user is authenticated
	if auth, ok := sess.Values["authenticated"].(bool); !ok || !auth {
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
	sess.Values["authenticated"] = true
	sess.Save(r, w)
}

func logout(w http.ResponseWriter, r *http.Request) {
	sess, _ := store.Get(r, "cookie-name")

	// Revoke users authentication
	sess.Values["authenticated"] = false
	sess.Options.MaxAge = -1
	sess.Save(r, w)
}

func main() {
	http.HandleFunc("/secret", secret)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)

	http.ListenAndServe(":8080", nil)
}
