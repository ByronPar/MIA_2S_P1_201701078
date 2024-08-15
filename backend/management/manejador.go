package management

import (
	"MIA_2S_P1_201701078/beans"
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

// Función RegistroProfesor
func RegistroProfesor() {
	var id int32
	var cui string
	var nombre string
	var curso string

	// abrir archivo en modo escritura
	arch, err := os.OpenFile("HT_1.dat", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)

	// Verificar si hubo un error
	if err != nil {
		fmt.Println(err)
		return
	}
	// Cerrar el archivo al terminar la función
	defer arch.Close()

	// Mover el cursor al final del archivo
	arch.Seek(0, io.SeekEnd)

	// crear un profesor
	var profesorNuevo beans.Profesor
	profesorNuevo.Tipo = int32(1)
	fmt.Println("")
	fmt.Println("")
	fmt.Println("***************************************************************************************")

	// Solicitar ID
	fmt.Print("Ingrese el ID del profesor: ")
	fmt.Scanln(&id)
	profesorNuevo.Id_profesor = id

	// Solicitar CUI
	fmt.Print("Ingrese el CUI del profesor: ")
	fmt.Scanln(&cui)
	copy(profesorNuevo.CUI[:], cui)

	// Solicitar Nombre
	fmt.Print("Ingrese el nombre del profesor: ")
	fmt.Scanln(&nombre)
	copy(profesorNuevo.Nombre[:], nombre)

	// Solicitar Curso
	fmt.Print("Ingrese el curso del profesor: ")
	fmt.Scanln(&curso)
	copy(profesorNuevo.Curso[:], curso)

	// Escribir el profesor en el archivo
	binary.Write(arch, binary.LittleEndian, &profesorNuevo)
	arch.Close()
	fmt.Println("")
	fmt.Println("Profesor registrado con éxito")

	fmt.Println("***************************************************************************************")
	fmt.Println("")
	fmt.Println("")
}

// Función RegistroEstudiante
func RegistroEstudiante() {

	var id int32
	var cui string
	var nombre string
	var carnet string

	// abrir archivo en modo escritura
	arch, err := os.OpenFile("HT_1.dat", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)

	// Verificar si hubo un error
	if err != nil {
		fmt.Println(err)
		return
	}
	// Cerrar el archivo al terminar la función
	defer arch.Close()

	// Mover el cursor al final del archivo
	arch.Seek(0, io.SeekEnd)

	// crear un profesor
	var estudianteNuevo beans.Estudiante
	estudianteNuevo.Tipo = int32(2)
	fmt.Println("")
	fmt.Println("")
	fmt.Println("***************************************************************************************")

	// Solicitar ID
	fmt.Print("Ingrese el ID del estudiante: ")
	fmt.Scanln(&id)
	estudianteNuevo.Id_estudiante = id

	// Solicitar CUI
	fmt.Print("Ingrese el CUI del estudiante: ")
	fmt.Scanln(&cui)
	copy(estudianteNuevo.CUI[:], cui)

	// Solicitar Nombre
	fmt.Print("Ingrese el nombre del estudiante: ")
	fmt.Scanln(&nombre)
	copy(estudianteNuevo.Nombre[:], nombre)

	// Solicitar Carnet
	fmt.Print("Ingrese el carnet del estudiante: ")
	fmt.Scanln(&carnet)
	copy(estudianteNuevo.Carnet[:], carnet)

	// Escribir el estudiante en el archivo
	binary.Write(arch, binary.LittleEndian, &estudianteNuevo)
	arch.Close()
	fmt.Println("")
	fmt.Println("Estudiante registrado con éxito")

	fmt.Println("***************************************************************************************")
	fmt.Println("")
	fmt.Println("")

}

// Función VerRegistros
func VerRegistros() {
	// abrir archivo en modo lectura
	arch, err := os.OpenFile("HT_1.dat", os.O_RDONLY, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer arch.Close()

	// leer el archivo con bucle para leer todos los registros
	for {
		var profesor beans.Profesor
		err = binary.Read(arch, binary.LittleEndian, &profesor)
		if err == io.EOF {
			break
		}
		if err != nil {
			break
		}

		// Imprimir información del registro
		fmt.Println("*****************************************************************")
		fmt.Println("Tipo: ", profesor.Tipo)
		if profesor.Tipo == 1 {
			fmt.Println("Rol: Profesor")
			fmt.Println("ID: ", profesor.Id_profesor)
			fmt.Println("CUI: ", string(profesor.CUI[:]))
			fmt.Println("Nombre: ", string(profesor.Nombre[:]))
			fmt.Println("Curso: ", string(profesor.Curso[:]))
			fmt.Println("*****************************************************************")
			fmt.Println("")
		} else if profesor.Tipo == 2 {
			fmt.Println("Rol: Estudiante")
			fmt.Println("ID: ", profesor.Id_profesor) // Debería ser Id_estudiante
			fmt.Println("CUI: ", string(profesor.CUI[:]))
			fmt.Println("Nombre: ", string(profesor.Nombre[:]))
			fmt.Println("Carnet: ", string(profesor.Curso[:]))
			fmt.Println("*****************************************************************")
			fmt.Println("")
		}
	}
}

// Función CrearArchivo
func CrearArchivo() {
	if _, err := os.Stat("HT_1.dat"); os.IsNotExist(err) {
		// Crear archivo si no existe
		arch, err := os.Create("HT_1.dat")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer arch.Close()
	}
}
