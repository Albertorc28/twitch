package handlers

import (
	"io/ioutil"
	"net/http"
	"time"
	"twitch/server/core/sesion"
	"twitch/server/core/stream"
	"twitch/server/data/model"

	"github.com/gorilla/mux"
)

//NEWSTREAM Nuevo stream
const NEWSTREAM string = "/newstream"

//EXITSTREAM Cerrar un stream (el propio)
const EXITSTREAM string = "/stream/{streamer}/exit"

//GETFRAME Nuevo stream
const GETFRAME string = "/stream/{streamer}/frame"

//SETFRAME Nuevo stream
const SETFRAME string = "/stream/{streamer}/addframe" //gg wp
//gl para que no te cuelen un "frame" en tu propio stream

func init() {
	Manejadores[NEWSTREAM] = NewStreamHandler
	Manejadores[EXITSTREAM] = ExitStreamHandler
	Manejadores[GETFRAME] = GetFrameHandler
	Manejadores[SETFRAME] = SetFrameHandler
}

//NewStreamHandler inicia un nuevo stream
func NewStreamHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}

	anfitrion := sesion.GetCookie(r)

	if anfitrion != nil {

		nuevoStream := model.RStream{
			Usuario: anfitrion.Nombre,
			Inicio:  time.Now().UTC(),
			Expira:  time.Now().UTC().Add(time.Minute * 1),
			ID:      time.Now().Unix(),
		}

		//sqlclient.CrearStream(&nuevoStream)

		stream.NuevaSala(&nuevoStream, *anfitrion)

		http.Redirect(w, r, "/stream/"+anfitrion.Nombre, 301)
	} else {
		http.Redirect(w, r, "/sesion", 301)
	}
}

//GetFrameHandler descarga el frame actual del stream
func GetFrameHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}

	params := mux.Vars(r)

	sala := stream.GetSala(params["streamer"])
	w.WriteHeader(http.StatusOK)
	if sala != nil {
		if sala.Online && sala.Frame != nil {
			w.Write(sala.Frame)
		} else {
			bytes, _ := ioutil.ReadFile("static/img/stream_offline.jpg")
			w.Write(bytes)
		}
	} else {
		bytes, _ := ioutil.ReadFile("static/img/stream_inexistente.jpg")
		w.Write(bytes)
	}
}

//SetFrameHandler proporciona un nuevo frame al stream
func SetFrameHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}

	params := mux.Vars(r)

	sala := stream.GetSala(params["streamer"])
	if sala != nil {
		defer r.Body.Close()

		mpf, _, err := r.FormFile("frame")

		frame, err := ioutil.ReadAll(mpf)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			frame, _ = ioutil.ReadFile("static/img/stream_offline.jpg")
			sala.KeepAlive <- false
		} else {
			w.WriteHeader(http.StatusOK)
			sala.KeepAlive <- true
		}
		sala.GuardarFrame(frame)
	} else {
		http.Error(w, "El usuario no tiene una sala activa", http.StatusServiceUnavailable)
	}
}

//ExitStreamHandler cierra un stream y destruye la sala si el llamador es el propietario
func ExitStreamHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}

	user := sesion.GetCookie(r)

	params := mux.Vars(r)

	if user != nil {
		if user.Nombre == params["streamer"] {
			if sala := stream.GetSala(user.Nombre); sala != nil {
				stream.CancelarSala(user.Nombre)
			} else {
				http.Error(w, "Tu sala no existe", http.StatusNotFound)
			}
		} else {
			http.Error(w, "No eres el propietario de este stream", http.StatusUnauthorized)
		}
	} else {
		http.Error(w, "No has iniciado sesión", http.StatusUnauthorized)
	}
}
