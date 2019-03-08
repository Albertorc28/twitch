package handlers

import (
	"net/http"
	"strings"
)

//JS es la ruta inicial de todo JS
const JS string = "/js/{file}"

//JSLIBS es la ruta inicial de toda lib externa JS
const JSLIBS string = "/js/libs/{file}"

//CSS es la ruta inicial de todo CSS
const CSS string = "/css/{file}"

//IMAGES es la ruta inicial de toda imagen
const IMAGES string = "/static/img/{file}"

//ICONS es la ruta inicial de todo icono
const ICONS string = "/static/ico/{file}"

//THUMBNAILS es la ruta inicial de todo thumbnail
const THUMBNAILS string = "/static/thumbnails/streams/{streamer}/{streamId}/{fileName}"

func init() {
	Manejadores[JS] = JSHandler
	Manejadores[JSLIBS] = JSHandler
	Manejadores[CSS] = CSSHandler
	Manejadores[IMAGES] = ImagesHandler
	Manejadores[ICONS] = IconsHandler
	Manejadores[THUMBNAILS] = ThumbnailsHandler
}

//JSHandler devuelve los archivos JS
func JSHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}
	if !strings.HasSuffix(r.URL.Path, ".js") {
		http.NotFound(w, r)
		return
	}
	var file = strings.TrimLeft(r.URL.Path, "/")
	http.ServeFile(w, r, file)
}

//CSSHandler devuelve los archivos CSS
func CSSHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}

	if !strings.HasSuffix(r.URL.Path, ".css") {
		http.NotFound(w, r)
		return
	}
	var file = strings.TrimLeft(r.URL.Path, "/")
	http.ServeFile(w, r, file)
}

//ImagesHandler devuelve los archivos de imagen
func ImagesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}

	if !strings.HasSuffix(r.URL.Path, ".png") &&
		!strings.HasSuffix(r.URL.Path, ".jpg") &&
		!strings.HasSuffix(r.URL.Path, ".jpeg") {
		http.NotFound(w, r)
		return
	}
	var file = strings.TrimLeft(r.URL.Path, "/")
	http.ServeFile(w, r, file)
}

//IconsHandler devuelve los archivos de imagen
func IconsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}

	if !strings.HasSuffix(r.URL.Path, ".ico") {
		http.NotFound(w, r)
		return
	}
	var file = strings.TrimLeft(r.URL.Path, "/")
	http.ServeFile(w, r, file)
}

//ThumbnailsHandler devuelve los archivos de imagen
func ThumbnailsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}

	if !strings.HasSuffix(r.URL.Path, ".jpg") &&
		!strings.HasSuffix(r.URL.Path, ".jpeg") {
		http.NotFound(w, r)
		return
	}
	var file = strings.TrimLeft(r.URL.Path, "/")
	http.ServeFile(w, r, file)
}
