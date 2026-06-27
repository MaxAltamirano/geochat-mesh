package main

import (
	"encoding/json"
	"fmt"
	"geochat-mesh/telemetria"
	"geochat-mesh/transport"
	"github.com/gorilla/websocket"
	"log"
	"math"
	"net/http"
	"os"
	"time"
)

var node *GeoNode

// Configuración de WebSocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("[ERROR] Fallo al establecer WebSocket:", err)
		return
	}
	defer conn.Close()

	log.Println("[MESH] Frontend conectado. Canal bidireccional abierto.")

	// Ubicación base del nodo (Avellaneda)
	const miLat, miLon = -34.6611, -58.3644

	// BUCLE PRINCIPAL: Recepción de alertas desde el frontend (Usuario)
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("[MESH] Cliente desconectado:", err)
			break
		}

		// 1. Decodificar el mensaje del usuario
		var newAlert Alert
		if err := json.Unmarshal(msg, &newAlert); err == nil {
			
			// 2. Aplicar Tacto Espacial (Filtro)
			// Si el mensaje no trae lat/lon propias, usamos las del evento o valores por defecto
			if newAlert.Lat == 0 { newAlert.Lat = -34.6650 }
			if newAlert.Lon == 0 { newAlert.Lon = -58.3700 }
			
			distancia := calcularDistancia(miLat, miLon, newAlert.Lat, newAlert.Lon)

			// 3. Validación de relevancia
			if distancia <= 2.0 {
				newAlert.Timestamp = time.Now().Format("15:04:05")
				newAlert.DistanciaKm = distancia
				newAlert.Location = "Avellaneda-Mesh"

				// 4. Persistimos en disco mediante la estructura del nodo
				node.AddAlert(newAlert)
				node.SaveState("./local_data")

				log.Printf("[MESH] ✅ Alerta relevante guardada: %s (a %.2f km)", newAlert.Type, distancia)
				
				// Opcional: Enviar confirmación al frontend
				conn.WriteMessage(websocket.TextMessage, []byte("ALERTA_GUARDADA_EXITOSAMENTE"))
			} else {
				log.Printf("[MESH] 🔇 Ignorando alerta lejana: %s (a %.2f km)", newAlert.Type, distancia)
				conn.WriteMessage(websocket.TextMessage, []byte("ALERTA_IGNORADA_FUERA_DE_RANGO"))
			}

		} else {
			log.Println("[ERROR] Mensaje mal formado recibido:", err)
		}
	}
}

// Ejemplo de lógica de Tacto Espacial
const RadioAccionKm = 5.0 // Solo nos importa lo que pase a 5km a la redonda

func EsRelevante(lat, lon float64, miLat, miLon float64) bool {
	// Aquí iría el cálculo de distancia Haversine
	distancia := calcularDistancia(lat, lon, miLat, miLon)
	return distancia <= RadioAccionKm
}

// calcularDistancia usa la fórmula de Haversine para medir distancia entre dos coordenadas en KM
func calcularDistancia(lat1, lon1, lat2, lon2 float64) float64 {
	const radioTierra = 6371 // Radio en KM

	dLat := (lat2 - lat1) * (math.Pi / 180)
	dLon := (lon2 - lon1) * (math.Pi / 180)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*(math.Pi/180))*math.Cos(lat2*(math.Pi/180))*
			math.Sin(dLon/2)*math.Sin(dLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return radioTierra * c
}

func main() {
	fmt.Println("🚀 Iniciando Nodo Mesh - Modo: Soberano")

	// 1. Inicialización de Capa de Transporte (Interfaz Agnóstica)
	var network transport.Transport
	if os.Getenv("ENV") == "PROD_LORA" {
		network = &transport.LoRaDriver{Puerto: "/dev/ttyUSB0"}
	} else {
		network = &transport.MockDriver{} // Simulación activa
	}

	// 2. Inicialización del Nodo (Persistencia y Carga de Estado)
	// IMPORTANTE: Asignamos a la variable global 'node' definida arriba en tu código
	node = NewGeoNode("LLAVERO_001")
	dataDir := "./local_data"

	err := os.MkdirAll(dataDir, 0755)
	if err != nil {
		log.Fatalf("[ERROR] No se pudo crear el directorio de datos: %v", err)
	}

	err = node.SaveState(dataDir)
	if err != nil {
		log.Printf("[ERROR] No se pudo persistir el estado del nodo: %v", err)
	}

	// Visualización del estado inicial del Nodo
	nodeData, _ := json.MarshalIndent(node, "", "  ")
	fmt.Printf("📦 Estado del Nodo Llavero:\n%s\n", string(nodeData))

	// 3. Telemetría a Google (Reporte de salud)
	telemetria.ReportarEstado("nodo_status", 1.0)

	// 4. Servidor de API / Status
	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(node)
	})

	// 5. Canal de recepción de alertas (El "oído" del nodo)
	go func() {
		mensajes, err := network.Recibir()
		if err != nil {
			log.Printf("[ERROR] Fallo en la capa de transporte: %v", err)
			return
		}

		// Posición actual del nodo (Avellaneda)
		const miLat = -34.6611
		const miLon = -58.3644

		for msg := range mensajes {
			// Simulamos coordenadas del evento (luego vendrán dentro de msg)
			latEvento, lonEvento := -34.6650, -58.3700
			distancia := calcularDistancia(miLat, miLon, latEvento, lonEvento)

			// Tacto Espacial: Filtro de proximidad (2km)
			if distancia <= 2.0 {
				fmt.Printf("⚠️ ALERTA RELEVANTE en %s (a %.2f km)\n", msg.Tipo, distancia)

				// Creamos la alerta con coordenadas
				nuevaAlerta := Alert{
					Type:      msg.Tipo,
					Location:  "Avellaneda-Mesh",
					Timestamp: time.Now().Format("15:04:05"),
					Lat:       latEvento,
					Lon:       lonEvento,
				}

				// A. Procesar y guardar en memoria
				node.AddAlert(nuevaAlerta)

				// B. Persistir el estado completo en disco (Cromosoma de Memoria)
				err := node.SaveState("./local_data")
				if err != nil {
					log.Printf("[ERROR] Fallo al actualizar Cromosoma: %v", err)
				} else {
					log.Println("💾 [CROMOSOMA] Alerta persistida exitosamente.")
				}
			} else {
				fmt.Printf("🔇 Ignorando alerta lejana: %s (a %.2f km)\n", msg.Tipo, distancia)
			}
		}
	}()

	// 6. Lanzamiento del Hub (Servicio de red)
	http.HandleFunc("/ws", handleWebSocket)
	fmt.Println("GeoChat Mesh Hub operando en puerto 8081...")

	// El servidor se queda bloqueado aquí escuchando
	log.Fatal(http.ListenAndServe(":8081", nil))
}
