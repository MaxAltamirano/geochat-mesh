package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
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

	// 1. GOROUTINE: Envío de alertas desde el nodo al frontend
	// Esto corre en paralelo para no bloquear la recepción
	go func() {
		for {
			// Aquí podrías enviar el estado actual de las alertas guardadas
			// Si el nodo detecta algo, lo manda por aquí
			time.Sleep(10 * time.Second) 
		}
	}()

	// 2. BUCLE PRINCIPAL: Recepción de alertas desde el frontend (Usuario)
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("[MESH] Cliente desconectado:", err)
			break
		}

		// Procesamos la alerta enviada por el usuario
		var newAlert Alert
		if err := json.Unmarshal(msg, &newAlert); err == nil {
			// Asignamos tiempo si no viene definido
			if newAlert.Timestamp == "" {
				newAlert.Timestamp = time.Now().Format("15:04:05")
			}
			
			// Persistimos en disco mediante la estructura del nodo
			// Nota: Asegúrate de tener acceso a la variable 'node' aquí
			node.AddAlert(newAlert)
			
			log.Printf("[MESH] Nueva alerta guardada por el usuario: %s en %s", newAlert.Type, newAlert.Location)
		} else {
			log.Println("[ERROR] Mensaje mal formado recibido:", err)
		}
	}
}

func main() {

	// 1. Inicialización: Obtenemos el nodo desde node.go
	node := NewGeoNode("LLAVERO_001")
	dataDir := "./local_data"
    
	// 2. Persistencia: Aseguramos entorno y guardamos estado
	err := os.MkdirAll(dataDir, 0755)
	if err != nil {
		log.Fatalf("[ERROR] No se pudo crear el directorio de datos: %v", err)
	}

	err = node.SaveState(dataDir)
	if err != nil {
		log.Printf("[ERROR] No se pudo persistir el estado del nodo: %v", err)
	}

	// Visualización de arranque
	nodeData, _ := json.MarshalIndent(node, "", "  ")
	fmt.Printf("GeoChat Mesh Node activo con estado persistente:\n%s\n", string(nodeData))

	// 3. Servidor: Levantamos el hub en el puerto 8080
	http.HandleFunc("/ws", handleWebSocket)
	fmt.Println("GeoChat Mesh Hub operando en puerto 8081...")

	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(node) // Devuelve el JSON completo del nodo
	})
	
	log.Fatal(http.ListenAndServe(":8081", nil))
}