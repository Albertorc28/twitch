package model

//PaginaIndex representa la estructura de la plantilla index.html
type PaginaIndex struct {
	Logado          bool
	TengoSalaActiva bool
	Usuario         string
	HayStreams      bool
	Streams         []interface{}
}

//PaginaStream representa la estructura de la plantilla index.html
type PaginaStream struct {
	Logado       bool
	EsElStreamer bool
	Streamer     string
	Usuario      string
}
