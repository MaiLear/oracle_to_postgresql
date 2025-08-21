package constants

// Tipos de programa (PRF_TIPO_PROGRAMA)
const (
	TypeProgramComplementary = "C" // Complementaria
	TypeProgramTitled        = "T" // Titulada
)

// Mapeo de códigos de tipo de programa a descripciones
var ProgramTypeMap = map[string]string{
	TypeProgramComplementary: "COMPLEMENTARIA",
	TypeProgramTitled:        "TITULADA",
}

// GetProgramTypeDescription retorna la descripción del tipo de programa
func GetProgramTypeDescription(programType string) string {
	if description, exists := ProgramTypeMap[programType]; exists {
		return description
	}
	return "Tipo Desconocido"
}
