package loggers

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

var (
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
)

// InitLogger inicializa los loggers para un componente específico (por ejemplo: "server" o "worker").
// Crea un archivo de log con la ruta: logs/{component}/{component}_YYYY-MM-DD.log
func InitLogger(component string) error {
	// Crea subdirectorio de logs por componente
	logDir := filepath.Join("logs", component)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("error creando directorio de logs: %w", err)
	}

	// Nombre del archivo con fecha
	currentTime := time.Now()
	logFileName := fmt.Sprintf("%s_%s.log", component, currentTime.Format("2006-01-02"))
	logPath := filepath.Join(logDir, logFileName)

	// Abre archivo
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("error abriendo archivo de log: %w", err)
	}

	// Escribe en archivo + consola
	multiWriter := io.MultiWriter(os.Stdout, logFile)

	InfoLogger = log.New(multiWriter, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(multiWriter, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	return nil
}

// RotateLogs elimina archivos de log antiguos (más de 7 días) de todas las carpetas dentro de "logs/"
func RotateLogs() error {
	components := []string{"server", "worker"} // Agrega aquí otros si es necesario

	now := time.Now()

	for _, component := range components {
		pattern := filepath.Join("logs", component, fmt.Sprintf("%s_*.log", component))
		files, err := filepath.Glob(pattern)
		if err != nil {
			return fmt.Errorf("error listando logs para %s: %w", component, err)
		}

		for _, file := range files {
			dateStr := file[len(filepath.Join("logs", component))+1 : len(file)-len(".log")]
			dateStr = dateStr[len(component)+1:] // Remover el prefijo del componente_
			fileDate, err := time.Parse("2006-01-02", dateStr)
			if err != nil {
				continue
			}

			if now.Sub(fileDate) > 7*24*time.Hour {
				if err := os.Remove(file); err != nil {
					ErrorLogger.Printf("No se pudo eliminar %s: %v", file, err)
				} else {
					InfoLogger.Printf("Log antiguo eliminado: %s", file)
				}
			}
		}
	}

	return nil
}
