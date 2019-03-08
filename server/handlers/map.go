package handlers

import "net/http"

//ManejadorHTTP encapsula como tipo la función de manejo de peticiones HTTP, para que sea posible almacenar sus referencias en un diccionario
type ManejadorHTTP = func(w http.ResponseWriter, r *http.Request)

//Manejadores es el diccionario general de las peticiones que son manejadas por nuestro servidor
var Manejadores = make(map[string]ManejadorHTTP)
