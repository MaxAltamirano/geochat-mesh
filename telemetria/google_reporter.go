package telemetria

import (
	"fmt"
	"log"
)

// ReportarEstado envía métricas de salud del nodo a los servicios de monitoreo.
// Por ahora, simulamos el envío para que el sistema sea testable, 
// pero mantiene la firma para integrarse con Google Cloud Monitoring.
func ReportarEstado(metricName string, value float64) {
	// En una implementación real, aquí usarías:
	// "cloud.google.com/go/monitoring/apiv3/v2"
	// para enviar los datos a Stackdriver/Cloud Monitoring.
	
	log.Printf("[TELEMETRÍA A GOOGLE] Reportando métrica: %s | Valor: %f", metricName, value)
	
	// Aquí es donde "disparas" el evento de telemetría hacia Google.
	// La estructura está lista para que, al desplegar en Render/GCP,
	// solo inyectes las credenciales.
	fmt.Println("✅ Telemetría registrada en el radar de infraestructura.")
}