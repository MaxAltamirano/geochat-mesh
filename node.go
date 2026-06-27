package main

import (
	"encoding/json"
	"os"
	//"path/filepath"
)

type Alert struct {
	Type      string `json:"type"`
	Location  string `json:"location"`
	Timestamp string `json:"timestamp"`
	Lat       float64 `json:"lat"` // Nueva coordenada
	Lon       float64 `json:"lon"` // Nueva coordenada
	DistanciaKm   float64 `json:"distancia_km"` // NUEVO CAMPO
}

type GeoNode struct {
	NodeID              string                 `json:"node_id"`
	Status              string                 `json:"status"`
	BlockchainHandshake bool                   `json:"blockchain_handshake"` // <--- ESTE ES EL CAMPO QUE FALTABA
	Metadata            map[string]interface{} `json:"metadata"`             // <--- ESTE ES EL CAMPO QUE FALTABA
	Alerts              []Alert                `json:"alerts"`
}

// NewGeoNode inicializa un nuevo nodo con valores por defecto.
func NewGeoNode(nodeID string) *GeoNode {
	return &GeoNode{
		NodeID:              nodeID,
		Status:              "ACTIVE",
		BlockchainHandshake: false,
		Metadata: map[string]interface{}{
			"version": "1.0.0",
		},
		Alerts: []Alert{}, // Inicializamos como lista vacía
	}
}


func (n *GeoNode) AddAlert(a Alert) {
	n.Alerts = append(n.Alerts, a)
}

func (n *GeoNode) SaveState(dir string) error {
	file, _ := json.MarshalIndent(n, "", "  ")
	return os.WriteFile(dir + "/state_" + n.NodeID + ".json", file, 0644)
}

// Authenticate marca el handshake como exitoso y actualiza el estado.
func (n *GeoNode) Authenticate() {
	n.BlockchainHandshake = true
}

