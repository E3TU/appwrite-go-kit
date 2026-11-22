package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Middleware returns a Gin middleware that requires an authenticated user.
//
// It fails closed: a missing or unverifiable session aborts with 401 and never
// falls through to the handler. On success the *User is stored under UserKey
// and also available through UserFrom.
//
// header is the optional custom session header; empty means cookie only.
func (s *Service) Middleware(header string) gin.HandlerFunc {
	return func(c *gin.Context) {
		sess, err := ExtractSession(c.Request, header, s.cookieName)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		if s.verify == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		user, err := s.verify.Verify(sess.Secret)
		if err != nil || user == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		c.Set(UserKey, user)
		c.Next()
	}
}

// UserFrom returns the authenticated user attached by Middleware.
func UserFrom(c *gin.Context) (*User, bool) {
	return userFromContext(c)
}
