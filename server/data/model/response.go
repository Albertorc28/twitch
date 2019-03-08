package model

//GenericResponse a
type GenericResponse struct {
	Ok    bool
	Data  interface{}
	Error error
}
