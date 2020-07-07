package session

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"net/http"
	"template2/lib/cache"
	"template2/lib/util"
	"time"
)

// Amount of time for cookies/redis keys to expire.
var sessionExpire = 86400 * 30

var ErrNil = StoreError("session: not found")
var ErrInvalidCookie = StoreError("session: invalid session ID")
var ErrStoreFail = StoreError("session: storage fail")

type StoreError string

func (e StoreError) Error() string { return string(e) }

// DataSerializer provides an interface hook for alternative serializers
type DataSerializer interface {
	Deserialize(data []byte, sess *sessions.Session) error
	Serialize(sess *sessions.Session) ([]byte, error)
}

// JSONSerializer encode the session map to JSON.
type JSONSerializer struct{}

// Serialize to JSON. Will err if there are unmarshalable key values
func (s JSONSerializer) Serialize(sess *sessions.Session) ([]byte, error) {
	m := make(map[string]interface{}, len(sess.Values))
	for k, v := range sess.Values {
		ks, ok := k.(string)
		if !ok {
			err := fmt.Errorf("Non-string key value, cannot serialize session to JSON: %v", k)
			fmt.Printf("session.JSONSerializer.serialize() Error: %v", err)
			return nil, err
		}
		m[ks] = v
	}
	return json.Marshal(m)
}

// Deserialize back to map[string]interface{}
func (s JSONSerializer) Deserialize(data []byte, sess *sessions.Session) error {
	m := make(map[string]interface{})
	err := json.Unmarshal(data, &m)
	if err != nil {
		fmt.Printf("session.JSONSerializer.deserialize() Error: %v", err)
		return err
	}
	for k, v := range m {
		sess.Values[k] = v
	}
	return nil
}

// GobSerializer uses gob package to encode the session map
type GobSerializer struct{}

// Serialize using gob
func (s GobSerializer) Serialize(sess *sessions.Session) ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(sess.Values)
	if err == nil {
		return buf.Bytes(), nil
	}
	return nil, err
}

// Deserialize back to map[interface{}]interface{}
func (s GobSerializer) Deserialize(data []byte, sess *sessions.Session) error {
	dec := gob.NewDecoder(bytes.NewBuffer(data))
	return dec.Decode(&sess.Values)
}

// IdGenerator provides an interface hook for alternative ID generator
type IdGenerator interface {
	Generate(length int) (string, error)
}

// Generate Random ID and encoded into base64
type Base64ID struct{}

// Generate Random ID and encoded into base64
func (id Base64ID) Generate(length int) (string, error) {
	b, err := util.GenerateRandomKey(length)
	if err != nil {
		return "", err
	}

	return util.Base64Encode(b), nil
}

// CookieHandler provides an interface hook for alternative encode/decode cookie value
type CookieHandler interface {
	Encode(sess *sessions.Session, store *HybridStore) (string, error)
	Decode(c *http.Cookie, sess *sessions.Session, store *HybridStore) error
}

// Secure cookie handler
type SecureCookieHandler struct {
	Codecs []securecookie.Codec
}

func (h SecureCookieHandler) Encode(sess *sessions.Session, store *HybridStore) (string, error) {
	return securecookie.EncodeMulti(sess.Name(), sess.ID, h.Codecs...)
}

func (h SecureCookieHandler) Decode(c *http.Cookie, sess *sessions.Session, store *HybridStore) error {
	return securecookie.DecodeMulti(c.Name, c.Value, &sess.ID, h.Codecs...)
}

// Standard cookie handler
type StdCookieHandler struct{}

func (h StdCookieHandler) Encode(sess *sessions.Session, store *HybridStore) (string, error) {
	return sess.ID, nil
}

func (h StdCookieHandler) Decode(c *http.Cookie, sess *sessions.Session, store *HybridStore) error {
	sess.ID = c.Value
	return nil
}

// StoreEngine provides an interface hook for alternative session value storage such as
// redis, postgres, mysql, etc
type StoreEngine interface {
	load(ctx context.Context, key string) ([]byte, error)
	save(ctx context.Context, key string, value []byte, duration time.Duration) error
	delete(ctx context.Context, key string) error
}

type RedisStoreEngine struct {
	RedisClient *cache.RedisClient
}

func (r *RedisStoreEngine) load(ctx context.Context, key string) ([]byte, error) {
	data, err := r.RedisClient.GetBytes(ctx, key)
	if err != nil {
		switch err {
		case redis.Nil:
			return nil, ErrNil
		default:
			return nil, ErrStoreFail
		}
	}
	return data, nil
}

func (r *RedisStoreEngine) save(ctx context.Context, key string, value []byte, duration time.Duration) error {
	_, err := r.RedisClient.Set(ctx, key, value, duration)
	return err
}

func (r *RedisStoreEngine) delete(ctx context.Context, key string) error {
	_, err := r.RedisClient.Del(ctx, key)
	return err
}

type HybridStoreConf struct {
	IdLength      int               `json:"IdLength"`
	KeyPrefix     string            `json:"KeyPrefix"`
	Options       *sessions.Options `json:"Options"`
	IdGenerator   string            `json:"IdGenerator"`
	Serializer    string            `json:"Serializer"`
	CookieHandler string            `json:"CookieHandler"`
	KeyPairs      []string          `json:"KeyPairs"`
}

type HybridStore struct {
	Storage       StoreEngine
	Options       *sessions.Options // default configuration
	idLength      int
	keyPrefix     string
	idGenerator   IdGenerator
	serializer    DataSerializer
	cookieHandler CookieHandler
}

func NewSessionStore(storage StoreEngine, conf *HybridStoreConf) *HybridStore {
	var (
		err           error
		idGenerator   IdGenerator    = Base64ID{}
		serializer    DataSerializer = GobSerializer{}
		cookieHandler CookieHandler  = StdCookieHandler{}
	)

	// config session id generator
	switch conf.IdGenerator {
	case "base64":
		idGenerator = Base64ID{}
	}

	// config session serializer
	switch conf.Serializer {
	case "gob":
		serializer = GobSerializer{}
	case "json":
		serializer = JSONSerializer{}
	}

	// initial cookie Options
	Options := sessions.Options{
		Path:   "/",
		MaxAge: sessionExpire,
	}
	if conf.Options != nil {
		Options = *conf.Options
	}

	// config cookie encode/decode handler
	switch conf.CookieHandler {
	case "standard":
		cookieHandler = StdCookieHandler{}
	case "std":
		cookieHandler = StdCookieHandler{}
	case "secure":
		// initial keyPairs
		keyPairs := make([][]byte, len(conf.KeyPairs))
		for i, v := range conf.KeyPairs {
			if keyPairs[i], err = util.Base64Decode(v); err != nil {
				panic("Error: loading Invalid key for secure cookie handler.")
			}
		}
		codecs := securecookie.CodecsFromPairs(keyPairs...)

		for _, s := range codecs {
			if cookie, ok := s.(*securecookie.SecureCookie); ok {
				cookie.MaxAge(Options.MaxAge)
				//cookie.SetSerializer(securecookie.JSONEncoder{})
				//cookie.HashFunc(sha512.New512_256)
			}
		}
		cookieHandler = SecureCookieHandler{
			Codecs: codecs,
		}
	case "aes_gcm":
		encKey, err := util.Base64Decode(conf.KeyPairs[0])
		if err != nil {
			panic("Error: loading Invalid key for aes_gcm cookie handler.")
		}
		cookieHandler, err = NewAesGcmCookieHandler(encKey)
		if err != nil {
			panic("Create AES_GCM cookie handler fail: " + err.Error())
		}
	}

	// initial session store
	store := &HybridStore{
		Storage:       storage,
		Options:       &Options,
		idLength:      conf.IdLength,
		keyPrefix:     conf.KeyPrefix,
		idGenerator:   idGenerator,
		serializer:    serializer,
		cookieHandler: cookieHandler,
	}

	return store
}

func (s *HybridStore) Get(r *http.Request, name string) (*sessions.Session, error) {
	return sessions.GetRegistry(r).Get(s, name)
}

// Create a new session by searching cookie name
//
// Return value:
//   - a new session if given cookie name is not found, no error. Session.ID is an empty string
//   - load an exist session if given cookie name is found, and the key value corresponding to
//     the session is found in our database
//   - return nil, ErrNil if given cookie name is found but the key value corresponding to
//     the session does not exist in our database or the format of session ID is wrong
//   - return nil, ErrStoreFail if there is something wrong with the backend storage
func (s *HybridStore) New(r *http.Request, name string) (*sessions.Session, error) {
	var (
		err error
		ok  bool
	)
	session := sessions.NewSession(s, name)
	// make a copy
	options := *s.Options
	session.Options = &options
	session.IsNew = true

	// it will not look up the storage if cookie not find
	if c, errCookie := r.Cookie(name); errCookie == nil {
		fmt.Println("find cookie", name)
		err = s.cookieHandler.Decode(c, session, s)
		if err == nil {
			ok, err = s.load(r.Context(), session)
			fmt.Println("New session value", session.Values)
			session.IsNew = !(err == nil && ok) // not new if no error and data available
		} else {
			err = ErrInvalidCookie
		}

		// there is error for either decoding or loading session, reset session ID to empty
		if err != nil {
			session.ID = ""
		}
	}

	return session, err
}

// Save the session to backend storage
//   - If MaxAge <= 0, then use Set-Coolie Header to delete the cookie in client side.
//   - If session ID does not exist, create a new session ID before saving.
func (s *HybridStore) Save(r *http.Request, w http.ResponseWriter, session *sessions.Session) error {
	// Marked for deletion.
	if session.Options.MaxAge < 0 {
		// delete the cookie in the database only if session.ID is set.
		if session.ID != "" {
			if err := s.delete(r.Context(), session); err != nil {
				return err
			}
		}
		http.SetCookie(w, sessions.NewCookie(session.Name(), "", session.Options))
	} else {
		var err error

		if session.ID == "" {
			if session.ID, err = s.idGenerator.Generate(s.idLength); err != nil {
				return err
			}
		}
		if err := s.save(r.Context(), session); err != nil {
			return err
		}
		encoded, err := s.cookieHandler.Encode(session, s)

		if err != nil {
			return err
		}

		fmt.Println("set session cookie")
		http.SetCookie(w, sessions.NewCookie(session.Name(), encoded, session.Options))
	}

	return nil
}

// save stores the session in redis.
func (s *HybridStore) save(ctx context.Context, session *sessions.Session) error {
	data, err := s.serializer.Serialize(session)
	fmt.Println("save session value", session.Values)
	fmt.Println("save session serial", data)

	if err != nil {
		return err
	}

	//if s.maxLength != 0 && len(b) > s.maxLength {
	//	return errors.New("HybridStore: the value to store is too big")
	//}

	age := session.Options.MaxAge
	if age == 0 {
		age = s.Options.MaxAge
	}
	fmt.Println("save ID", s.keyPrefix+session.ID)
	err = s.Storage.save(ctx, s.keyPrefix+session.ID, data, time.Duration(age))
	fmt.Println("save err", err)
	return err
}

// load reads the session from redis.
// returns true if there is a sessoin data in DB
func (s *HybridStore) load(ctx context.Context, session *sessions.Session) (bool, error) {
	data, err := s.Storage.load(ctx, s.keyPrefix+session.ID)
	fmt.Println("load session id", s.keyPrefix+session.ID)
	fmt.Println("load", data)
	fmt.Println("err", err)
	if err != nil {
		return false, err
	}

	if data == nil {
		return false, nil // no data was associated with this key
	}

	return true, s.serializer.Deserialize(data, session)
}

// delete removes keys from redis if MaxAge<0
func (s *HybridStore) delete(ctx context.Context, session *sessions.Session) error {
	return s.Storage.delete(ctx, s.keyPrefix+session.ID)
}

// Get the idGenerator
func (s *HybridStore) IdGenerator() IdGenerator {
	return s.idGenerator
}

// Get the idLength of session store
func (s *HybridStore) IdLength() int {
	return s.idLength
}

// contextKey is the type used to store the registry in the context.
type contextKey int

// collectionKey is the key used to store the registry in the context.
const collectionKey contextKey = 0

// Collection stores sessions used during a request.
type Collection struct {
	request  *http.Request
	writer   http.ResponseWriter
	sessions map[string]*sessions.Session
}

func errorSessionNotExist(name string) error {
	return fmt.Errorf("sessions: error session %s not exist", name)
}

func (s *Collection) Get(name string) (session *sessions.Session, err error) {
	session, ok := s.sessions[name]
	if !ok {
		return nil, errorSessionNotExist(name)
	}
	return session, nil
}

func (s *Collection) UpdateValue(name string, key, value interface{}) error {
	session, ok := s.sessions[name]
	if !ok {
		return errorSessionNotExist(name)
	}
	session.Values[key] = value
	session.IsNew = true
	return nil
}

func (s *Collection) MaxAge(name string, age int) error {
	session, ok := s.sessions[name]
	if !ok {
		return errorSessionNotExist(name)
	}

	session.Options.MaxAge = age
	session.IsNew = true
	return nil
}

func (s *Collection) Save(name string) error {
	session, ok := s.sessions[name]
	if !ok {
		return fmt.Errorf("sessions: error session %s not exist", name)
	}
	return session.Save(s.request, s.writer)
}

func (s *Collection) SaveAll(name string) error {
	var errMulti util.MultiError
	for name, session := range s.sessions {
		if err := session.Save(s.request, s.writer); err != nil {
			errMulti = append(errMulti, fmt.Errorf(
				"sessions: error saving session %q -- %v", name, err))
		}
	}

	if errMulti != nil {
		return errMulti
	}
	return nil
}

func NewCollection(r *http.Request, w http.ResponseWriter, sessions map[string]*sessions.Session) *Collection {
	return &Collection{
		request:  r,
		writer:   w,
		sessions: sessions,
	}
}

// NewContext returns a new Context that carries value u.
func NewCollectionContext(ctx context.Context, s *Collection) context.Context {
	return context.WithValue(ctx, collectionKey, s)
}

// FromContext returns the User value stored in ctx, if any.
func FromCollectionContext(ctx context.Context) (*Collection, bool) {
	s, ok := ctx.Value(collectionKey).(*Collection)
	return s, ok
}
