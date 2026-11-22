package auth

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestExtractSessionHeaderWins(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("X-Session", "from-header")
	r.AddCookie(&http.Cookie{Name: "sid", Value: "from-cookie"})

	s, err := ExtractSession(r, "X-Session", "sid")
	if err != nil {
		t.Fatalf("ExtractSession() error = %v", err)
	}
	if s.Secret != "from-header" || s.Source != SourceHeader {
		t.Fatalf("Session = %+v, want header source", s)
	}
}

func TestExtractSessionFallsBackToCookie(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.AddCookie(&http.Cookie{Name: "sid", Value: "from-cookie"})

	s, err := ExtractSession(r, "X-Session", "sid")
	if err != nil {
		t.Fatalf("ExtractSession() error = %v", err)
	}
	if s.Secret != "from-cookie" || s.Source != SourceCookie {
		t.Fatalf("Session = %+v, want cookie source", s)
	}
}

func TestExtractSessionMissing(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	if _, err := ExtractSession(r, "X-Session", "sid"); !errors.Is(err, ErrNoSession) {
		t.Fatalf("error = %v, want ErrNoSession", err)
	}
	if _, err := ExtractSession(nil, "", "sid"); !errors.Is(err, ErrNoSession) {
		t.Fatalf("nil request error = %v, want ErrNoSession", err)
	}
}

type fakeVerifier struct {
	user *User
	err  error
}

func (f fakeVerifier) Verify(string) (*User, error) { return f.user, f.err }

func TestMiddlewareRejectsWithoutSession(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := NewWithVerifier(fakeVerifier{user: &User{ID: "u"}}, "sid")

	r := gin.New()
	r.GET("/me", svc.Middleware(""), func(c *gin.Context) {
		t.Error("handler reached without a session")
	})

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/me", nil))

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", w.Code)
	}
}

func TestMiddlewareRejectsInvalidSession(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := NewWithVerifier(fakeVerifier{err: ErrSessionInvalid}, "sid")

	r := gin.New()
	r.GET("/me", svc.Middleware(""), func(c *gin.Context) {
		t.Error("handler reached with an invalid session")
	})

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req.AddCookie(&http.Cookie{Name: "sid", Value: "bad"})

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", w.Code)
	}
}

func TestMiddlewareAttachesUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	want := &User{ID: "u1", Email: "a@b.c"}
	svc := NewWithVerifier(fakeVerifier{user: want}, "sid")

	r := gin.New()
	r.GET("/me", svc.Middleware("X-Session"), func(c *gin.Context) {
		got, ok := UserFrom(c)
		if !ok || got != want {
			t.Errorf("UserFrom() = %v, %v", got, ok)
		}
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req.Header.Set("X-Session", "secret")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
}

func TestMiddlewareFailsClosedWithoutVerifier(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := &Service{cookieName: "sid"}

	r := gin.New()
	r.GET("/me", svc.Middleware(""), func(c *gin.Context) {
		t.Error("handler reached without a verifier")
	})

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req.AddCookie(&http.Cookie{Name: "sid", Value: "x"})

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", w.Code)
	}
}
