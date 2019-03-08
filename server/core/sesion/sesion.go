package sesion

import (
	"net/http"
	"time"

	"twitch/server/data/model"
	m "twitch/server/data/model"

	"github.com/gorilla/securecookie"
)

//SESSIONCOOKIE define el nombre de la cookie de sesión
const SESSIONCOOKIE = "TwitchAuth"

//hashKey y blockKey para codificar y encriptar la cookie.
//Al ser aleatorios, cada vez que se reinicie el servidor estos valores cambiarán.
//Por esto el usuario deberá de volver a logearse al reiniciarse el servidor.
var hashKey = []byte(securecookie.GenerateRandomKey(32))
var blockKey = []byte(securecookie.GenerateRandomKey(32))

//SecureCookie
var sc = securecookie.New(hashKey, blockKey)

func GetCookie(r *http.Request) *m.Usuario {
	var user *m.Usuario
	if cookie, err := r.Cookie(SESSIONCOOKIE); err == nil {
		value := make(map[string]string)
		if err = sc.Decode(SESSIONCOOKIE, cookie.Value, &value); err == nil {
			user = &m.Usuario{}
			user.FromCookie(value)
		}
	}

	return user

}

func GenerarCookie(usuario *model.Usuario) (*http.Cookie, error) {
	var err error
	if encoded, err := sc.Encode(SESSIONCOOKIE, usuario.ToCookieMap()); err == nil {
		expire := time.Now().UTC().AddDate(0, 0, 1)
		cookie := &http.Cookie{
			Name:     SESSIONCOOKIE,
			Value:    encoded,
			Path:     "/",
			Expires:  expire,
			HttpOnly: true,
		}
		return cookie, nil
	}
	return nil, err
}

func BorrarCookie(w http.ResponseWriter) {
	c := &http.Cookie{
		Name:    SESSIONCOOKIE,
		Path:    "/",
		Expires: time.Unix(0, 0),
		MaxAge:  -1,
	}
	http.SetCookie(w, c)
}
