package management

import (
	"backend/beans"
	"backend/functions"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

/*-------------------------menu de comando--------------------*/
var strRespuesta = ""

func Analizar(strInput string) string {

	resultado := strings.Split(strInput, "\n")

	//  ejecuta cada linea ingresada
	for _, comando := range resultado {
		split_comando(comando)
	}
	return strRespuesta
}

// Separa los diferentes comando con sus parametros si tienen
func split_comando(comando string) {
	var commandArray []string
	// Elimina los saltos de linea y retornos de carro
	comando = strings.Replace(comando, "\n", "", 1)
	comando = strings.Replace(comando, "\r", "", 1)
	// Banderas para verficar comentarios
	band_comentario := false
	if strings.Contains(comando, "#") {
		// Comentario
		band_comentario = true
		fmt.Println(comando)
	} else {
		// Comando con Parametros
		commandArray = strings.Split(comando, " -")
	}
	// Ejecuta el comando leido si no es un comentario
	if !band_comentario {
		ejecutar_comando(commandArray)
	}
}

// Identifica y ejecuta el comando encontrado
func ejecutar_comando(commandArray []string) {
	// Convierte el comando a minusculas
	data := strings.ToLower(commandArray[0])

	// Identifica el comando a ejecutar
	if data == "mkdisk" {
		strRespuesta = strRespuesta + "\n" + mkdisk(commandArray)
	} else if data == "rmdisk" {
		strRespuesta = strRespuesta + "\n" + rmdisk(commandArray)
		/*} else if data == "fdisk" {
			fdisk(commandArray)
		} else if data == "rep" {
			rep()
		*/
	} else {
		strRespuesta = strRespuesta + "\n" + "[ERROR] El comando no fue reconocido..."
	}
}

/*-------------------------/menu de comando--------------------*/

/*--------------------------Comandos--------------------------*/

func mkdisk(commandArray []string) string {
	// Variables para los valores de los parametros
	val_size := 0
	val_fit := ""
	val_unit := ""
	val_path := ""

	// Banderas para verificar los parametros y ver si se repiten
	band_size := false
	band_fit := false
	band_unit := false
	band_path := false
	band_error := false

	// Obtengo solo los parametros validos
	for i := 1; i < len(commandArray); i++ {
		aux_data := strings.SplitAfter(commandArray[i], "=")
		data := strings.ToLower(aux_data[0])
		val_data := aux_data[1]

		// Identifica los parametos
		switch {
		/* PARAMETRO OBLIGATORIO -> SIZE */
		case strings.Contains(data, "size="):
			// Valido si el parametro ya fue ingresado
			if band_size {
				return "[ERROR-MKDISK] El parametro -size ya fue ingresado."
			}

			// Activo la bandera del parametro
			band_size = true

			// Conversion a entero
			aux_size, err := strconv.Atoi(val_data)
			val_size = aux_size

			// ERROR de conversion
			if err != nil {
				band_error = true
				return functions.Msg_error(err, "MKDISK")
			}

			// Valido que el tamaño sea positivo
			if val_size < 0 {
				return "[ERROR-MKDISK] El parametro -size es negativo."
			}
		/* PARAMETRO OPCIONAL -> FIT */
		case strings.Contains(data, "fit="):
			// Valido si el parametro ya fue ingresado
			if band_fit {
				return "[ERROR-MKDISK] El parametro -fit ya fue ingresado."
			}

			// Le quito las comillas y lo paso a minusculas
			val_fit = strings.Replace(val_data, "\"", "", 2)
			val_fit = strings.ToLower(val_fit)

			if val_fit == "bf" {
				// Activo la bandera del parametro y obtengo el caracter que me interesa
				band_fit = true
				val_fit = "b"
			} else if val_fit == "ff" {
				// Activo la bandera del parametro y obtengo el caracter que me interesa
				band_fit = true
				val_fit = "f"
			} else if val_fit == "wf" {
				// Activo la bandera del parametro y obtengo el caracter que me interesa
				band_fit = true
				val_fit = "w"
			} else {
				return "[ERROR-MKDISK] El Valor del parametro -fit no es valido..."
			}
		/* PARAMETRO OPCIONAL -> UNIT */
		case strings.Contains(data, "unit="):
			// Valido si el parametro ya fue ingresado
			if band_unit {
				return "[ERROR-MKDISK] El parametro -unit ya fue ingresado..."
			}

			// Reemplaza comillas y lo paso a minusculas
			val_unit = strings.Replace(val_data, "\"", "", 2)
			val_unit = strings.ToLower(val_unit)

			if val_unit == "k" || val_unit == "m" {
				// Activo la bandera del parametro
				band_unit = true
			} else {
				// Parametro no valido
				return "[ERROR-MKDISK] El Valor del parametro -unit no es valido..."
			}
		/* PARAMETRO OBLIGATORIO -> PATH */
		case strings.Contains(data, "path="):
			if band_path {
				return "[ERROR-MKDISK] El parametro -path ya fue ingresado..."
			}
			// Activo la bandera del parametro
			band_path = true
			// Reemplaza comillas
			val_path = strings.Replace(val_data, "\"", "", 2)
		/* PARAMETRO NO VALIDO */
		default:
			return "[ERROR-MKDISK] Parametro no valido..."
		}
	}

	// Verifico si no hay errores
	// Verifico que el parametro "Path" (Obligatorio) este ingresado
	// Verifico que el parametro "Size" (Obligatorio) este ingresado
	if !band_error && band_path && band_size {
		total_size := 1024
		master_boot_record := beans.Mbr{}
		// Disco -> Archivo Binario
		// crear el archivo que simula el disco
		aux, err := filepath.Abs(val_path)

		// ERROR
		if err != nil {
			return functions.Msg_error(err, "MKDISK")
		}

		// Verifica si el directorio ya existe
		if _, err := os.Stat(aux); !errors.Is(err, os.ErrNotExist) {
			if err == nil {
				// El directorio ya existe
				return "[ERROR-MKDISK] El disco ya existe."
			}
		}
		// Crea el directiorio de forma recursiva
		cmd1 := exec.Command("/bin/sh", "-c", "sudo mkdir -p '"+filepath.Dir(aux)+"'")
		cmd1.Dir = "/"
		_, err1 := cmd1.Output()
		// ERROR
		if err1 != nil {
			return functions.Msg_error(err, "MKDISK")
		}
		// Da los permisos al directorio
		cmd2 := exec.Command("/bin/sh", "-c", "sudo chmod -R 777 '"+filepath.Dir(aux)+"'")
		cmd2.Dir = "/"
		_, err2 := cmd2.Output()

		// ERROR
		if err2 != nil {
			return functions.Msg_error(err, "MKDISK")
		}
		// Verifica si existe la ruta para el archivo
		if _, err := os.Stat(filepath.Dir(aux)); errors.Is(err, os.ErrNotExist) {
			if err != nil {
				return "[ERROR-MKDISK] No se pudo crear el disco."
			}
		}
		// finaliza creación del disco
		// Fecha
		fecha := time.Now()
		str_fecha := fecha.Format("02/01/2006 15:04:05")
		// Copio valor al Struct
		copy(master_boot_record.Mbr_fecha_creacion[:], str_fecha)
		// Numero aleatorio
		rng := rand.New(rand.NewSource(time.Now().UnixNano()))
		min := 0
		max := 100
		num_random := rng.Intn(max-min+1) + min
		// Copio valor al Struct
		copy(master_boot_record.Mbr_dsk_signature[:], strconv.Itoa(num_random))
		// Verifico si existe el parametro "Fit" (Opcional)
		if band_fit {
			// Copio valor al Struct
			copy(master_boot_record.Dsk_fit[:], val_fit)
		} else {
			// Si no especifica -> "Primer Ajuste"
			copy(master_boot_record.Dsk_fit[:], "f")
		}
		// Verifico si existe el parametro "Unit" (Opcional)
		if band_unit {
			// Megabytes
			if val_unit == "m" {
				copy(master_boot_record.Mbr_tamano[:], strconv.Itoa(val_size*1024*1024))
				total_size = val_size * 1024
			} else {
				// Kilobytes
				copy(master_boot_record.Mbr_tamano[:], strconv.Itoa(val_size*1024))
				total_size = val_size
			}
		} else {
			// Si no especifica -> Megabytes
			copy(master_boot_record.Mbr_tamano[:], strconv.Itoa(val_size*1024*1024))
			total_size = val_size * 1024
		}

		// Inicializar Parcticiones
		for i := 0; i < 4; i++ {
			copy(master_boot_record.Mbr_partition[i].Part_status[:], "0")
			copy(master_boot_record.Mbr_partition[i].Part_type[:], "0")
			copy(master_boot_record.Mbr_partition[i].Part_fit[:], "0")
			copy(master_boot_record.Mbr_partition[i].Part_start[:], "-1")
			copy(master_boot_record.Mbr_partition[i].Part_size[:], "0")
			copy(master_boot_record.Mbr_partition[i].Part_name[:], "")
			copy(master_boot_record.Mbr_partition[i].Part_correlative[:], "")
			copy(master_boot_record.Mbr_partition[i].Part_id[:], "")
		}

		// Convierto de entero a string
		str_total_size := strconv.Itoa(total_size)
		// Comando para definir el tamaño (Kilobytes) y llenarlo de ceros
		cmd := exec.Command("/bin/sh", "-c", "dd if=/dev/zero of=\""+val_path+"\" bs=1024 count="+str_total_size)
		cmd.Dir = "/"
		_, err = cmd.Output()

		// ERROR
		if err != nil {
			return functions.Msg_error(err, "MKDISK")
		}

		// Se escriben los datos en disco
		// Apertura del archivo
		disco, err := os.OpenFile(val_path, os.O_RDWR, 0660)

		// ERROR
		if err != nil {
			return functions.Msg_error(err, "MKDISK")
		}

		// Conversion de struct a bytes
		mbr_byte := functions.Struct_a_bytes(master_boot_record)

		// Se posiciona al inicio del archivo para guardar la informacion del disco
		newpos, err := disco.Seek(0, 0)

		// ERROR
		if err != nil {
			return functions.Msg_error(err, "MKDISK")
		}

		// Escritura de struct en archivo binario
		_, err = disco.WriteAt(mbr_byte, newpos)

		// ERROR
		if err != nil {
			return functions.Msg_error(err, "MKDISK")
		}
		disco.Close()
	}
	return "[MENSAJE-MKDISK] finaliza comando, ejecutado exitosamente"
}

func rmdisk(commandArray []string) string {
	// Variables para los valores de los parametros
	val_path := ""
	// Banderas para verificar los parametros y ver si se repiten
	band_path := false
	band_error := false
	// Obtengo solo los parametros validos
	for i := 1; i < len(commandArray); i++ {
		aux_data := strings.SplitAfter(commandArray[i], "=")
		data := strings.ToLower(aux_data[0])
		val_data := aux_data[1]
		// Identifica los parametos
		switch {
		/* PARAMETRO OBLIGATORIO -> PATH */
		case strings.Contains(data, "path="):
			if band_path {
				fmt.Println("[ERROR-RMDISK] El parametro -path ya fue ingresado...")
				band_error = true
				break
			}
			// Activo la bandera del parametro
			band_path = true
			// Reemplaza comillas
			val_path = strings.Replace(val_data, "\"", "", 2)
		/* PARAMETRO NO VALIDO */
		default:
			return "[ERROR-RMDISK] Parametro no valido."
		}
	}
	// Verifico si no hay errores
	if !band_error {
		// Verifico que el parametro "Path" (Obligatorio) este ingresado
		// Si existe el archivo binario
		_, e := os.Stat(val_path)
		if e != nil {
			// Si no existe
			if os.IsNotExist(e) {
				return "[ERROR-RMDISK]  No existe el disco que desea eliminar."
			}
		} else {
			// Si existe
			cmd := exec.Command("/bin/sh", "-c", "rm \""+val_path+"\"")
			cmd.Dir = "/"
			_, err := cmd.Output()
			// ERROR
			if err != nil {
				return functions.Msg_error(err, "RMDISK")
			} else {
				return "[MENSAJE-RMDISK] El Disco fue eliminado!"
			}
		}
	}
	return "[MENSAJE-RMDISK] finaliza comando, ejecutado exitosamente"
}

func fdisk(commandArray []string) string {
	// Variables para los valores de los parametros
	val_size := 0
	val_unit := ""
	val_path := ""
	val_type := ""
	val_fit := ""
	val_name := ""
	// Banderas para verificar los parametros y ver si se repiten
	band_size := false
	band_unit := false
	band_path := false
	band_type := false
	band_fit := false
	band_name := false
	// Obtengo solo los parametros validos
	for i := 1; i < len(commandArray); i++ {
		aux_data := strings.SplitAfter(commandArray[i], "=")
		data := strings.ToLower(aux_data[0])
		val_data := aux_data[1]

		// Identifica los parametos
		switch {
		/* PARAMETRO OBLIGATORIO -> SIZE */
		case strings.Contains(data, "size="):
			// Valido si el parametro ya fue ingresado
			if band_size {
				return "[ERROR-FDISK] El parametro -size ya fue ingresado."
			}
			// Activo la bandera del parametro
			band_size = true
			// Conversion a entero
			aux_size, err := strconv.Atoi(val_data)
			val_size = aux_size
			// ERROR de conversion
			if err != nil {
				return functions.Msg_error(err, "FDISK")
			}
			// Valido que el tamaño sea positivo
			if val_size < 0 {
				return "[ERROR-FDISK] El parametro -size es negativo."
			}
		/* PARAMETRO OPCIONAL -> UNIT */
		case strings.Contains(data, "unit="):
			// Valido si el parametro ya fue ingresado
			if band_unit {
				return "[ERROR-FDISK] El parametro -unit ya fue ingresado."
			}
			// Reemplaza comillas y lo paso a minusculas
			val_unit = strings.Replace(val_data, "\"", "", 2)
			val_unit = strings.ToLower(val_unit)
			fmt.Println("Unit: ", val_unit)
			if val_unit == "b" || val_unit == "k" || val_unit == "m" {
				// Activo la bandera del parametro
				band_unit = true
			} else {
				// Parametro no valido
				return "[ERROR-FDISK] El Valor del parametro -unit no es valido."
			}
		/* PARAMETRO OBLIGATORIO -> PATH */
		case strings.Contains(data, "path="):
			if band_path {
				return "[ERROR-FDISK] El parametro -path ya fue ingresado."
			}
			// Activo la bandera del parametro
			band_path = true
			// Reemplaza comillas
			val_path = strings.Replace(val_data, "\"", "", 2)
		/* PARAMETRO OPCIONAL -> TYPE */
		case strings.Contains(data, "type="):
			if band_type {
				return "[ERROR-FDISK] El parametro -type ya fue ingresado."
			}
			// Reemplaza comillas y lo paso a minusculas
			val_type = strings.Replace(val_data, "\"", "", 2)
			val_type = strings.ToLower(val_type)
			fmt.Println("Type: ", val_type)
			if val_type == "p" || val_type == "e" || val_type == "l" {
				// Activo la bandera del parametro
				band_type = true
			} else {
				// Parametro no valido
				return "[ERROR-FDISK] El Valor del parametro -type no es valido."
			}
		/* PARAMETRO OPCIONAL -> FIT */
		case strings.Contains(data, "fit="):
			// Valido si el parametro ya fue ingresado
			if band_fit {
				return "[ERROR-FDISK] El parametro -fit ya fue ingresado."
			}
			// Le quito las comillas y lo paso a minusculas
			val_fit = strings.Replace(val_data, "\"", "", 2)
			val_fit = strings.ToLower(val_fit)
			if val_fit == "bf" {
				// Activo la bandera del parametro y obtengo el caracter que me interesa
				band_fit = true
				val_fit = "b"
			} else if val_fit == "ff" {
				// Activo la bandera del parametro y obtengo el caracter que me interesa
				band_fit = true
				val_fit = "f"
			} else if val_fit == "wf" {
				// Activo la bandera del parametro y obtengo el caracter que me interesa
				band_fit = true
				val_fit = "w"
			} else {
				return "[ERROR-FDISK] El Valor del parametro -fit no es valido."
			}
			fmt.Println("fit: ", val_fit)
		/* PARAMETRO OBLIGATORIO -> NAME */
		case strings.Contains(data, "name="):
			// Valido si el parametro ya fue ingresado
			if band_name {
				return "[ERROR-FDISK] El parametro -name ya fue ingresado."
			}
			// Activo la bandera del parametro
			band_name = true
			// Reemplaza comillas
			val_name = strings.Replace(val_data, "\"", "", 2)
		/* PARAMETRO NO VALIDO */
		default:
			return "[ERROR-FDISK] Parametro no valido."
		}
	}
	// Verifico si no hay errores
	if band_type {
		if val_type == "p" {
			// Primaria
			crear_particion_primaria(val_path, val_name, val_size, val_fit, val_unit)
		} else if val_type == "e" {
			// Extendida
		} else {
			// Logica
		}
	} else {
		// Si no lo indica se tomara como Primaria
		crear_particion_primaria(val_path, val_name, val_size, val_fit, val_unit)
	}
	return "[MENSAJE-FDISK] finaliza comando, ejecutado exitosamente"
}

/*--------------------------/Comandos--------------------------*/
