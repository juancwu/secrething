package cookies

import (
	"net/http"
	"time"
)

// Cookie represents the settings for an HTTP cookie
type Cookie struct {
	Name     string
	Value    string
	Path     string
	Domain   string
	MaxAge   int
	Secure   bool
	HttpOnly bool
	SameSite http.SameSite
}

// New creates a new cookie with default values
func New(name, value string) *Cookie {
	return &Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
}

// WithDomain sets the domain for the cookie
func (c *Cookie) WithDomain(domain string) *Cookie {
	c.Domain = domain
	return c
}

// WithPath sets the path for the cookie
func (c *Cookie) WithPath(path string) *Cookie {
	c.Path = path
	return c
}

// WithMaxAge sets the max age for the cookie in seconds
func (c *Cookie) WithMaxAge(seconds int) *Cookie {
	c.MaxAge = seconds
	return c
}

// WithExpiry sets the expiry time for the cookie
func (c *Cookie) WithExpiry(expiry time.Time) *Cookie {
	c.MaxAge = int(time.Until(expiry).Seconds())
	return c
}

// WithSecure sets whether the cookie should only be sent over HTTPS
func (c *Cookie) WithSecure(secure bool) *Cookie {
	c.Secure = secure
	return c
}

// WithHttpOnly sets whether the cookie should be accessible via JavaScript
func (c *Cookie) WithHttpOnly(httpOnly bool) *Cookie {
	c.HttpOnly = httpOnly
	return c
}

// WithSameSite sets the SameSite attribute for the cookie
func (c *Cookie) WithSameSite(sameSite http.SameSite) *Cookie {
	c.SameSite = sameSite
	return c
}

// ToHTTPCookie converts the Cookie to an http.Cookie
func (c *Cookie) ToHTTPCookie() *http.Cookie {
	return &http.Cookie{
		Name:     c.Name,
		Value:    c.Value,
		Path:     c.Path,
		Domain:   c.Domain,
		MaxAge:   c.MaxAge,
		Secure:   c.Secure,
		HttpOnly: c.HttpOnly,
		SameSite: c.SameSite,
	}
}

// SetCookie adds a cookie to the response
func SetCookie(w http.ResponseWriter, cookie *Cookie) {
	http.SetCookie(w, cookie.ToHTTPCookie())
}

// GetCookie retrieves a cookie from the request
func GetCookie(r *http.Request, name string) (*http.Cookie, error) {
	return r.Cookie(name)
}

// ClearCookie clears a cookie by setting its MaxAge to -1
func ClearCookie(w http.ResponseWriter, name string) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})
}