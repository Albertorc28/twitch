package stream

import (
	"time"
	"twitch/server/data/model"

	"github.com/gorilla/websocket"
)

var salas = make(map[string]*model.SalaStream)

//NuevaSala inicia una nueva sala (retransmisión) para un Usuario anfitrión
func NuevaSala(stream *model.RStream, usuario model.Usuario) {
	salas[usuario.Nombre] = &model.SalaStream{
		ID:        stream.ID,
		Anfitrion: usuario,
		Inicio:    stream.Inicio,
		Viewers:   0,
		Chat:      make(map[*websocket.Conn]model.Usuario),
		Broadcast: make(chan *model.Mensaje),
		KeepAlive: make(chan bool),
		Miniatura: model.DefaultStreamThumbnail,
	}
	go chatBroadcasting(salas[usuario.Nombre])
	go controlStream(salas[usuario.Nombre])
}

//CancelarSala cierra la sesión de streaming de una sala y elimina el objeto
func CancelarSala(anfitrion string) {
	sala := GetSala(anfitrion)
	if sala != nil {
		close(sala.Broadcast)
		close(sala.KeepAlive)
		for ws := range sala.Chat {
			ws.Close()
		}
		delete(salas, anfitrion)
	}
}

func chatBroadcasting(sala *model.SalaStream) {
	if sala != nil {
		for msg := range sala.Broadcast {
			if msg != nil {
				for ws := range sala.Chat {
					ws.WriteJSON(msg)
				}
			}
		}
	}
}

func controlStream(sala *model.SalaStream) {
	for {
		select {
		case vivo, abierto := <-sala.KeepAlive:
			if abierto {
				sala.Online = vivo
			} else {
				return
			}
		case <-time.After(2 * time.Second):
			sala.Online = false
		}
	}
}

//GetSala devuelve una sala por el nombre de anfitrión
func GetSala(anfitrion string) *model.SalaStream {
	var sala = salas[anfitrion]
	return sala
}

//GetStreams devuelve objetos RStream representando las salas
func GetStreams() []model.RStream {
	var streams = make([]model.RStream, 0)
	for _, v := range salas {
		s := model.RStream{
			Usuario:   v.Anfitrion.Nombre,
			Viewers:   len(v.Chat),
			ID:        v.ID,
			Inicio:    v.Inicio,
			Miniatura: v.Miniatura,
		}

		streams = append(streams, s)
	}
	return streams
}
