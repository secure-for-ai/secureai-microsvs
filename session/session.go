package session

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/secure-for-ai/secureai-microsvs/util"
	"net"
	"net/http"
	"time"
)

// Amount of time for cookies/redis keys to expire.
var sessionExpire = 86400 * 30

var ErrNil = StoreError("session: not found")
var ErrInvalidCookie = StoreError("session: invalid session ID")
var ErrStoreFail = StoreError("session: storage fail")
var ErrNonce = StoreError("session: storage fail")

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
	m := make(map[string]any, len(sess.Values))
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

// Deserialize back to map[string]any
func (s JSONSerializer) Deserialize(data []byte, sess *sessions.Session) error {
	m := make(map[string]any)
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

// Deserialize back to map[any]any
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

	return util.Base64EncodeToString(b), nil
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
	return securecookie.EncodeMulti(sess.Name(), store.EncodeSessionId(sess), h.Codecs...)
}

func (h SecureCookieHandler) Decode(c *http.Cookie, sess *sessions.Session, store *HybridStore) error {
	var data []byte
	err := securecookie.DecodeMulti(c.Name, c.Value, &data, h.Codecs...)
	if err != nil {
		return err
	}
	return store.DecodeSessionId(sess, data)
}

// Base64 cookie handler
type Base64CookieHandler struct{}

func (h Base64CookieHandler) Encode(sess *sessions.Session, store *HybridStore) (string, error) {
	return base64.RawURLEncoding.EncodeToString(store.EncodeSessionId(sess)), nil
}

func (h Base64CookieHandler) Decode(c *http.Cookie, sess *sessions.Session, store *HybridStore) error {
	data, err := base64.RawURLEncoding.DecodeString(c.Value)
	if err != nil {
		return err
	}
	return store.DecodeSessionId(sess, data)
}

// StoreEngine provides an interface hook for alternative session value storage such as
// redis, postgres, mysql, etc
type StoreEngine interface {
	init() error
	load(ctx context.Context, sess *sessions.Session) (bool, error)
	save(ctx context.Context, sess *sessions.Session, duration time.Duration) error
	delete(ctx context.Context, sess *sessions.Session) error
	newSessID(sess *sessions.Session) string
	encodeSessID(sess *sessions.Session) []byte
	decodeSessID(sess *sessions.Session, data []byte) error
}

type HybridStoreConf struct {
	Options       *sessions.Options `json:"Options"`
	CookieHandler string            `json:"CookieHandler"`
	KeyPairs      []string          `json:"KeyPairs"`
}

type HybridStore struct {
	Storage       StoreEngine
	Options       *sessions.Options // default configuration
	cookieHandler CookieHandler
}

func NewSessionStore(storage StoreEngine, conf *HybridStoreConf) *HybridStore {
	var (
		err           error
		cookieHandler CookieHandler = Base64CookieHandler{}
	)

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
	case "base64":
		cookieHandler = Base64CookieHandler{}
	case "secure":
		// initial keyPairs
		keyPairs := make([][]byte, len(conf.KeyPairs))
		for i, v := range conf.KeyPairs {
			if keyPairs[i], err = util.Base64DecodeString(v); err != nil {
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
		encKey, err := util.Base64DecodeString(conf.KeyPairs[0])
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
		cookieHandler: cookieHandler,
	}

	return store
}

func (s *HybridStore) Get(r *http.Request, name string) (*sessions.Session, error) {
	return sessions.GetRegistry(r).Get(s, name)
}

type SessionEx struct {
	*sessions.Session
	IP net.IP
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
	fmt.Println(err)

	return session, err
}

// Save the session to backend storage
//   - If MaxAge < 0, then use Set-Coolie Header to delete the cookie in client side.
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

		session.Values["ip"] = util.GetIP(r)
		session.Values["userAgent"] = r.Header.Get("User-Agent")
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

// return the
func (s *HybridStore) EncodeSessionId(sess *sessions.Session) []byte {
	return s.Storage.encodeSessID(sess)
}

func (s *HybridStore) DecodeSessionId(sess *sessions.Session, data []byte) error {
	return s.Storage.decodeSessID(sess, data)
}

// save stores the session in redis.
func (s *HybridStore) save(ctx context.Context, session *sessions.Session) error {
	age := session.Options.MaxAge
	if age == 0 {
		age = s.Options.MaxAge
	}

	err := s.Storage.save(ctx, session, time.Duration(age))
	fmt.Println("save ID", session.ID)
	fmt.Println("save err", err)
	return err
}

// load reads the session from redis.
// returns true if there is a sessoin data in DB
func (s *HybridStore) load(ctx context.Context, session *sessions.Session) (bool, error) {
	return s.Storage.load(ctx, session)
}

// delete removes keys from redis if MaxAge<0
func (s *HybridStore) delete(ctx context.Context, session *sessions.Session) error {
	return s.Storage.delete(ctx, session)
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

func (s *Collection) UpdateValue(name string, key, value any) error {
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
