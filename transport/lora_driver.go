package transport

// La 'L' y 'D' en mayúsculas hacen que sea público (Exportado)
type LoRaDriver struct {
	Puerto string
}

// Estos métodos también deben empezar con mayúscula
func (l *LoRaDriver) Enviar(msg MensajeGeoChat) error { return nil }
func (l *LoRaDriver) Recibir() (<-chan MensajeGeoChat, error) { return nil, nil }
func (l *LoRaDriver) EstadoConexion() string { return "LORA_OK" }