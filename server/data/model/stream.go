package model

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

//RStream define una fila de la tabla Stream
type RStream struct {
	ID        int64
	Usuario   string
	Inicio    time.Time
	Expira    time.Time
	Viewers   int
	Miniatura string
}

//SalaStream es la estructura de control de los streams
type SalaStream struct {
	ID              int64
	Online          bool
	Anfitrion       Usuario
	Frame           []byte
	UltimaMiniatura time.Time
	Chat            map[*websocket.Conn]Usuario
	Broadcast       chan *Mensaje
	KeepAlive       chan bool
	Inicio          time.Time
	Viewers         int
	Miniatura       string
}

//GuardarFrame guarda el frame y setea la miniatura y la fecha de último frame
func (sala *SalaStream) GuardarFrame(frame []byte) {
	sala.Frame = frame
	if time.Since(sala.UltimaMiniatura) > time.Duration(5*time.Second) {
		go sala.crearMiniatura(frame)
	}
}

func (sala *SalaStream) crearMiniatura(frame []byte) {
	carpeta := "static/thumbnails/streams/" + sala.Anfitrion.Nombre + "/" + strconv.FormatInt(sala.ID, 10) + "/"
	err := os.MkdirAll(carpeta, os.ModeDir)
	if err == nil {
		archivo := carpeta + "thumb.jpg"
		err = ioutil.WriteFile(archivo, frame, 0644)
		if err == nil {
			sala.Miniatura = "/" + archivo
			sala.UltimaMiniatura = time.Now()
		} else {
			log.Println("Crear thumbnail: " + err.Error())
		}
	} else {
		log.Println("Crear carpetas thumbnail: " + err.Error())
	}
}

//Mensaje es la estructura que define un mensaje corriente
type Mensaje struct {
	Anfitrion string
	Texto     string
	Autor     string
}

//ToJSON devuelve la estructura en JSON
func (msg *Mensaje) ToJSON() string {
	bytes, _ := json.Marshal(msg)
	return string(bytes)
}

//FromJSON compone la estructura desde su representación JSON
func (msg *Mensaje) FromJSON(js string) {
	json.Unmarshal([]byte(js), msg)
}
