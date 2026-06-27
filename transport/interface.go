package transport

// MensajeGeoChat es el tipo público que debe ser visible para main.go
type MensajeGeoChat struct {
    ID      string  `json:"id"`
    Tipo    string  `json:"tipo"`
    Lat     float64 `json:"lat"`   // Coordenada latitud
    Lon     float64 `json:"lon"`   // Coordenada longitud
    Payload []byte  `json:"payload"`
    Firma   string  `json:"firma"`
}

// Transport es la interfaz pública
type Transport interface {
    Enviar(msg MensajeGeoChat) error
    Recibir() (<-chan MensajeGeoChat, error)
    EstadoConexion() string
}