package model

import "strconv"

//LoginRequest a
type LoginRequest struct {
	User     string
	Password string
}

//Usuario a
type Usuario struct {
	ID     int64
	Email  string
	Nombre string
}

//ToCookieMap convierte la estructura en un mapa de strings
func (user *Usuario) ToCookieMap() map[string]string {
	value := map[string]string{
		"id":     strconv.FormatInt(user.ID, 10),
		"nombre": user.Nombre,
		"email":  user.Email,
	}
	return value
}

//FromCookie obtiene la estructura Usuario a partir de una cookie
func (user *Usuario) FromCookie(mapa map[string]string) {
	user.ID, _ = strconv.ParseInt(mapa["id"], 10, 64)
	user.Nombre = mapa["nombre"]
	user.Email = mapa["email"]
}
