package session

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"net/http"
	"template2/lib/cache"
	"template2/lib/util"
	"time"
)

// Amount of time for cookies/redis keys to expire.
var sessionExpire = 86400 * 30

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
type SecureCookieHandler struct{}

func (h SecureCookieHandler) Encode(sess *sessions.Session, store *HybridStore) (string, error) {
	return securecookie.EncodeMulti(sess.Name(), sess.ID, store.Codecs...)
}

func (h SecureCookieHandler) Decode(c *http.Cookie, sess *sessions.Session, store *HybridStore) error {
	return securecookie.DecodeMulti(c.Name, c.Value, &sess.ID, store.Codecs...)
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
	RedisClient   *cache.RedisClient
	Codecs        []securecookie.Codec
	Options       *sessions.Options // default configuration
	idLength      int
	keyPrefix     string
	idGenerator   IdGenerator
	serializer    DataSerializer
	cookieHandler CookieHandler
}

func NewSessionStore(client *cache.RedisClient, conf *HybridStoreConf) *HybridStore {
	var (
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

	// config cookie encode/decode handler
	switch conf.CookieHandler {
	case "standard":
		cookieHandler = StdCookieHandler{}
	case "std":
		cookieHandler = StdCookieHandler{}
	case "secure":
		cookieHandler = SecureCookieHandler{}
	}

	// initial cookie Options
	Options := sessions.Options{
		Path:   "/",
		MaxAge: sessionExpire,
	}
	if conf.Options != nil {
		Options = *conf.Options
	}

	// initial keyPairs
	keyPairs := make([][]byte, len(conf.KeyPairs))
	for i, v := range conf.KeyPairs {
		keyPairs[i], _ = util.Base64Decode(v)
	}
	codecs := securecookie.CodecsFromPairs(keyPairs...)

	for _, s := range codecs {
		if cookie, ok := s.(*securecookie.SecureCookie); ok {
			cookie.MaxAge(Options.MaxAge)
			//cookie.SetSerializer(securecookie.JSONEncoder{})
			//cookie.HashFunc(sha512.New512_256)
		}
	}
	// initial session store
	store := &HybridStore{
		RedisClient:   client,
		Codecs:        codecs,
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

	if c, errCookie := r.Cookie(name); errCookie == nil {
		err = s.cookieHandler.Decode(c, session, s)
		if err == nil {
			ok, err = s.load(session)
			fmt.Println("New session value", session.Values)
			session.IsNew = !(err == nil && ok) // not new if no error and data available
		}
	}
	//else {
	//	err = errCookie
	//}

	return session, err
}

func (s *HybridStore) Save(r *http.Request, w http.ResponseWriter, session *sessions.Session) error {
	// Marked for deletion.
	if session.Options.MaxAge <= 0 {
		if err := s.delete(session); err != nil {
			return err
		}
		http.SetCookie(w, sessions.NewCookie(session.Name(), "", session.Options))
	} else {
		var err error

		if session.ID == "" {
			if session.ID, err = s.idGenerator.Generate(s.idLength); err != nil {
				return err
			}
		}
		if err := s.save(session); err != nil {
			return err
		}
		encoded, err := s.cookieHandler.Encode(session, s)

		if err != nil {
			return err
		}

		http.SetCookie(w, sessions.NewCookie(session.Name(), encoded, session.Options))
	}

	return nil
}

// save stores the session in redis.
func (s *HybridStore) save(session *sessions.Session) error {
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
	count, err := s.RedisClient.Set(context.Background(), s.keyPrefix+session.ID, data, time.Duration(age)*time.Second)
	fmt.Println("save count", count)
	return err
}

// load reads the session from redis.
// returns true if there is a sessoin data in DB
func (s *HybridStore) load(session *sessions.Session) (bool, error) {
	data, err := s.RedisClient.GetBytes(context.Background(), s.keyPrefix+session.ID)
	fmt.Println("load session id", s.keyPrefix+session.ID)
	fmt.Println("load", data)
	if err != nil {
		return false, err
	}

	if data == nil {
		return false, nil // no data was associated with this key
	}

	return true, s.serializer.Deserialize(data, session)
}

// delete removes keys from redis if MaxAge<0
func (s *HybridStore) delete(session *sessions.Session) error {
	_, err := s.RedisClient.Del(context.Background(), s.keyPrefix+session.ID)
	return err
}
