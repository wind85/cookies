// package cookies provide a secure implementation based on gorilla securecookie
// to manage cookies, provides 3 functions SetCookieVal GetCookieVal and DelCookie
// The Cookie data type is an implementation the the CookieMng ( cookie manager )
// With this simple type is possible to exchange cookies securely over a non secure
// connection values are encrypted and decripted back, I don't suggest to pass
// sensitive data unless with an other encryption layer on top.
// It should be noted that for testing porpoises the HttpOnly and Secure are set
// to false, you must change that when using it on a production environment,to
// true.

package cookies

import (
	"net/http"
	"time"

	"github.com/gorilla/securecookie"
)

// CookieIface is the interface to be implemented in case other backends implementation
// are added in future right one is only one.
type CookieIface interface {
	Set(w http.ResponseWriter, r *http.Request, val map[string]interface{})
	Get(w http.ResponseWriter, r *http.Request) map[string]interface{}
	Del(w http.ResponseWriter, r *http.Request)
}

type Cookies struct {
	secure *securecookie.SecureCookie
	conf   *Conf
	name   string
}

// Conf is the struct the manages the settings for the cookie creation, must be passed
// to the constructor.
type Conf struct {
	HttpOnly bool
	Secure   bool
	MaxAge   int
}

// New accept a cookie name and a configuration and returns a valid cookie manager.
func New(name string, conf *Conf) *Cookies {
	return &Cookies{
		securecookie.New(
			securecookie.GenerateRandomKey(64),
			securecookie.GenerateRandomKey(32),
		),
		conf,
		name,
	}
}

// Set set the cookie with map of values map[string]string
func (c *Cookies) Set(w http.ResponseWriter, r *http.Request, val map[string]string) {
	c.setCookie(w, r, val, false, false, 0, time.Now().Add(168*time.Hour))
}

// Get gets the map from the request
func (c *Cookies) Get(w http.ResponseWriter, r *http.Request) map[string]string {
	value := make(map[string]string)

	if cookie, err := r.Cookie(c.name); err == nil {
		if err = c.secure.Decode(c.name, cookie.Value, &value); err != nil {
			http.Error(w, "Ops internal server error\n", http.StatusInternalServerError)
			return nil
		}
	}

	return value
}

// Del clears the cookie from the client
func (c *Cookies) Del(w http.ResponseWriter, r *http.Request) {
	c.setCookie(w, r, nil, false, false, -1, time.Now().Add(168*time.Hour))
}

func (c *Cookies) setCookie(w http.ResponseWriter, r *http.Request,
	val map[string]string, secure, httponly bool, age int, expiration time.Time) {

	if encoded, err := c.secure.Encode(c.name, val); err == nil {
		cookie := &http.Cookie{
			Name:     c.name,
			Value:    encoded,
			Path:     "/",
			HttpOnly: httponly,
			Secure:   secure,
			MaxAge:   age,
			Expires:  expiration,
		}

		http.SetCookie(w, cookie)
	}
}
