package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"text/template"
	"time"
	"twitch/server/core/sesion"
	"twitch/server/core/stream"
	"twitch/server/data/model"

	"github.com/gorilla/mux"
)

//PAGEINDEX es la ruta de inicio
const PAGEINDEX string = "/"

//PAGEAPP es la ruta de inicio
const PAGEAPP string = "/app"

//PAGELOGIN a
const PAGELOGIN string = "/sesion"

//PAGESTREAM a
const PAGESTREAM string = "/stream/{streamer}"

func init() {
	Manejadores[PAGEINDEX] = PageIndexHandler
	Manejadores[PAGEAPP] = PageAppIndexHandler
	Manejadores[PAGELOGIN] = PageLoginHandler
	Manejadores[PAGESTREAM] = PageStreamHandler
}

//PageIndexHandler a
func PageIndexHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/app", 302)
}

//PageLoginHandler a
func PageLoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}
	http.ServeFile(w, r, "html/login.html")
}

//PageAppIndexHandler devuelve la página de inicio
func PageAppIndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}

	user := sesion.GetCookie(r)

	pagina := model.PaginaIndex{}
	if user != nil {
		pagina.Logado = user != nil
		pagina.Usuario = user.Nombre
		pagina.TengoSalaActiva = stream.GetSala(pagina.Usuario) != nil
	}

	var streams = stream.GetStreams()
	//var streams = sqlclient.Streams()
	fmt.Printf("Encontrados %d streams\n", len(streams))
	pagina.HayStreams = len(streams) > 0
	pagina.Streams = make([]interface{}, 0)
	for i := range streams {
		s := &streams[i]
		//s.Miniatura = model.DefaultStreamThumbnail
		pagina.Streams = append(pagina.Streams, s)
	}

	t, err := template.ParseFiles("html/index.html")
	if err != nil {
		panic(err)
	}
	t.Execute(w, pagina)
}

//PageStreamHandler devuelve la página de inicio
func PageStreamHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}

	var user = sesion.GetCookie(r)
	parametros := mux.Vars(r)
	sala := stream.GetSala(parametros["streamer"])

	t, err := template.ParseFiles("html/stream.html")
	pagina := model.PaginaStream{}
	if sala != nil {
		if err != nil {
			panic(err)
		}
		pagina.Streamer = sala.Anfitrion.Nombre
		if user == nil {
			pagina.Usuario = "anonimo" + strconv.FormatInt(time.Now().Unix(), 10)
		} else {
			pagina.Logado = true
			pagina.Usuario = user.Nombre
		}

		if pagina.Usuario == pagina.Streamer {
			pagina.EsElStreamer = true
		}
	} else {
		pagina.Streamer = ""
		if user == nil {
			pagina.Usuario = "anonimo" + strconv.FormatInt(time.Now().Unix(), 10)
		} else {
			pagina.Logado = true
			pagina.Usuario = user.Nombre
		}
	}
	t.Execute(w, pagina)
}
