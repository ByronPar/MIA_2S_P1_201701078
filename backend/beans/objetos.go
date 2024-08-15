package beans

// Definición de la estructura de datos
type Profesor struct {
	Tipo        int32
	Id_profesor int32
	CUI         [13]byte
	Nombre      [25]byte
	Curso       [65]byte
}

type Estudiante struct {
	Tipo          int32
	Id_estudiante int32
	CUI           [13]byte
	Nombre        [25]byte
	Carnet        [65]byte
}
