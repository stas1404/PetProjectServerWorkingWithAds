package cookie

import (
	"crypto/sha1"
	"fmt"
	"server/internal/user"
	"time"
)

const CookieName = "ads_server_cookie"

func GenerateCookieValue(u user.User) string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(u.Email+u.Password+u.Nickname+time.Now().String())))
}
