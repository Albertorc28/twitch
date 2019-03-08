package handlers

import (
	"log"
	"net/http"
	"twitch/server/core/sesion"
	"twitch/server/core/stream"
	"twitch/server/data/model"

	"github.com/gorilla/mux"

	"github.com/gorilla/websocket"
)

//CHATSOCKET Nuevo stream
const CHATSOCKET string = "/stream/{streamer}/chat/{userName}"

func init() {
	Manejadores[CHATSOCKET] = ChatHandler
}

var upg = websocket.Upgrader{}

//ChatHandler maneja la conexión de un usuario con el chat de un anfitrión
func ChatHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	streamer := params["streamer"]
	sala := stream.GetSala(streamer)

	if sala == nil {
		http.Error(w, "El streamer no ha iniciado una sala", http.StatusNotFound)
		return
	}

	c, err := upg.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Error upgradeando la conexión", http.StatusInternalServerError)
		return
	}

	defer c.Close()
	user := sesion.GetCookie(r)
	if user == nil {
		user = &model.Usuario{
			ID:     0,
			Nombre: params["userName"],
			Email:  "noemail",
		}
	}
	sala.Chat[c] = *user

	bienvenida := &model.Mensaje{
		Autor:     "Bot",
		Texto:     user.Nombre + " ha entrado en la sala",
		Anfitrion: streamer,
	}
	sala.Broadcast <- bienvenida

	for {
		var mensaje = &model.Mensaje{}
		err = c.ReadJSON(mensaje)

		if err != nil {
			log.Println("ReadJSON: " + err.Error())
			mensaje.Autor = "Bot"
			mensaje.Texto = user.Nombre + " ha salido de la sala"
			mensaje.Anfitrion = streamer
			delete(sala.Chat, c)
			sala.Broadcast <- mensaje
			return
		}

		if sala == nil {
			return
		}

		sala.Broadcast <- mensaje
	}
}
