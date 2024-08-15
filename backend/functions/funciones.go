package functions

import (
	"MIA_2S_P1_201701078/management"
	"fmt"
	"os"
)

func Menu() {
	management.CrearArchivo()
	var comando string
	fmt.Println("")
	fmt.Println("     ##################################### HT1 - Byron Par - Carnet: 201701078")
	fmt.Println("     1. Registro de Profesor")
	fmt.Println("     2. Registro de Estudiante")
	fmt.Println("     3. Ver Registros")
	fmt.Println("     4. Salir")
	fmt.Println("     #####################################")
	fmt.Print("     Ingrese una opción: ")
	fmt.Println("")

	fmt.Print("Ingrese una opción: ")
	fmt.Scanln(&comando)

	if comando == "1" {
		management.RegistroProfesor()
	} else if comando == "2" {
		management.RegistroEstudiante()
	} else if comando == "3" {
		management.VerRegistros()
	} else if comando == "4" {
		os.Exit(0)
	} else {
		fmt.Println("Opción no válida")
		fmt.Println("")

	}
	Menu()

}
