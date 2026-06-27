package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type GeoNode struct {
	NodeID              string                 `json:"node_id"`
	Status              string                 `json:"status"`
	BlockchainHandshake bool                   `json:"blockchain_handshake"`
	Metadata            map[string]interface{} `json:"metadata"`
	Alerts              []Alert                `json:"alerts"` // <--- Lista persistente
}

type Alert struct {
	Type      string `json:"type"`      // policía, accidente, etc.
	Location  string `json:"location"`  // coordenadas o referencia
	Timestamp string `json:"timestamp"`
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
	}
}

// SaveState guarda el estado actual del nodo en un archivo JSON en el directorio especificado.
func (n *GeoNode) SaveState(dataDir string) error {
	// Asegurar que el directorio exista
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return err
	}

	path := filepath.Join(dataDir, "state_"+n.NodeID+".json")
	file, err := json.MarshalIndent(n, "", "  ")
	if err != nil {
		return err
	}
	
	return os.WriteFile(path, file, 0644)
}

// Authenticate marca el handshake como exitoso y actualiza el estado.
func (n *GeoNode) Authenticate() {
	n.BlockchainHandshake = true
}

// Agrega esto en node.go
func (n *GeoNode) AddAlert(a Alert) {
    n.Alerts = append(n.Alerts, a)
    n.SaveState("./local_data") // Persistimos en disco cada vez que agregamos algo
}
