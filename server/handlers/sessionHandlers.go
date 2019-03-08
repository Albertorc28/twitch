package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
	"twitch/server/core/sesion"
	"twitch/server/core/stream"

	m "twitch/server/data/model"
)

//LOGIN a
const LOGIN string = "/login"

//LOGOUT a
const LOGOUT string = "/logout"

func init() {
	Manejadores[LOGIN] = LoginHandler
	Manejadores[LOGOUT] = LogoutHandler
}

//LoginHandler maneja el inicio de sesión
func LoginHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}

	defer r.Body.Close()
	bytes, e := ioutil.ReadAll(r.Body)
	if e != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Println("Texto: " + string(bytes))
	var req m.LoginRequest

	e = json.Unmarshal(bytes, &req)
	if e != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var response m.GenericResponse
	//if req.User == "agu" && req.Password == "agu" {
	usuario := &m.Usuario{
		ID:     time.Now().Unix(),
		Nombre: req.User,
		Email:  req.User + "@metaphase07.es",
	}

	//Codificamos y encriptamos la cookie
	if cookie, err := sesion.GenerarCookie(usuario); err == nil {
		http.SetCookie(w, cookie)
		response = m.GenericResponse{
			Ok:   true,
			Data: usuario,
		}
	} else {
		response = m.GenericResponse{
			Ok:    false,
			Error: err,
		}
	}
	//} else {
	//	response = m.GenericResponse{
	//		Ok:    false,
	//		Error: errors.New("Credenciales no válidas"),
	//	}
	//}

	w.WriteHeader(http.StatusOK)
	respuesta, _ := json.Marshal(response)
	fmt.Fprint(w, string(respuesta))

}

//LogoutHandler maneja la finalización de la sesión
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}

	user := sesion.GetCookie(r)
	if user != nil {
		sesion.BorrarCookie(w)
		sala := stream.GetSala(user.Nombre)
		if sala != nil {
			stream.CancelarSala(user.Nombre)
		}
	}
	w.WriteHeader(http.StatusOK)
}
