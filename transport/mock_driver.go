package transport

import (
	"fmt"
	"time"
	"math/rand"
)

type MockDriver struct{}

func (m *MockDriver) Enviar(msg MensajeGeoChat) error {
	// Simulamos el envío por aire
	fmt.Printf("[MOCK-LORA] Enviando paquete: %s | Contenido: %s\n", msg.ID, string(msg.Payload))
	return nil
}



func (m *MockDriver) EstadoConexion() string {
	return "SIMULADO: Conectado a Red LoRa Virtual"
}

func (m *MockDriver) Recibir() (<-chan MensajeGeoChat, error) {
	ch := make(chan MensajeGeoChat)
	
	go func() {
		for {
			// 1. Generamos coordenadas dinámicas alrededor de Avellaneda
			lat := -34.6611 + (rand.Float64()*0.1 - 0.05)
			lon := -58.3644 + (rand.Float64()*0.1 - 0.05)
			
			// 2. Enviamos el mensaje con las coordenadas incluidas
			ch <- MensajeGeoChat{
				Tipo: "POLICÍA",
				Lat:  lat,
				Lon:  lon,
			}
			
			// 3. Intervalo de simulación
			time.Sleep(5 * time.Second)
		}
	}()
	
	return ch, nil
}