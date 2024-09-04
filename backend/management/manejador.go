package management

import (
	"backend/beans"
	"backend/functions"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

/*-------------------------menu de comando--------------------*/
/*-------------------- Variables Globales --------------------*/
var lista_montajes = beans.New_lista()
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
	} else if data == "fdisk" {
		strRespuesta = strRespuesta + "\n" + fdisk(commandArray)
	} else if data == "mount" {
		strRespuesta = strRespuesta + "\n" + mount(commandArray)
	} else if data == "mkfs" {
		strRespuesta = strRespuesta + "\n" + mkfs(commandArray)
	} else if data == "rep" {
		strRespuesta = strRespuesta + "\n" + rep(commandArray)
	} else {
		strRespuesta = strRespuesta + "\n" + "[ERROR] El comando no fue reconocido..."
	}
}

/*-------------------------/menu de comando--------------------*/

/*--------------------------Comandos principales--------------------------*/

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
				band_unit = true
			} else {
				return "[ERROR-FDISK] El Valor del parametro -unit no es valido."
			}
		/* PARAMETRO OBLIGATORIO -> PATH */
		case strings.Contains(data, "path="):
			if band_path {
				return "[ERROR-FDISK] El parametro -path ya fue ingresado."
			}
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
				band_fit = true
				val_fit = "b"
			} else if val_fit == "ff" {
				band_fit = true
				val_fit = "f"
			} else if val_fit == "wf" {
				band_fit = true
				val_fit = "w"
			} else {
				return "[ERROR-FDISK] El Valor del parametro -fit no es valido."
			}
		/* PARAMETRO OBLIGATORIO -> NAME */
		case strings.Contains(data, "name="):
			// Valido si el parametro ya fue ingresado
			if band_name {
				return "[ERROR-FDISK] El parametro -name ya fue ingresado."
			}
			band_name = true
			// Reemplaza comillas
			val_name = strings.Replace(val_data, "\"", "", 2)
		/* PARAMETRO NO VALIDO */
		default:
			return "[ERROR-FDISK] Parametro no valido."
		}
	}

	// Verifico si no hay errores

	if band_size {
		if band_path {
			if band_name {
				if band_type {
					if val_type == "p" {
						// Primaria
						return functions.Crear_particion_primaria(val_path, val_name, val_size, val_fit, val_unit, "FDISK")
					} else if val_type == "e" {
						// Extendida
						return functions.Crear_particion_extendia(val_path, val_name, val_size, val_fit, val_unit, "FDISK")
					} else {
						// Logica
						return functions.Crear_particion_logica(val_path, val_name, val_size, val_fit, val_unit, "FDISK")
					}
				} else {
					// Si no lo indica se tomara como Primaria
					return functions.Crear_particion_primaria(val_path, val_name, val_size, val_fit, val_unit, "FDISK")
				}
			} else {
				return "[ERROR-FDISK] El parametro -name no fue ingresado"
			}
		} else {
			return "[ERROR-FDISK] El parametro -path no fue ingresado"
		}
	} else {
		return "[ERROR-FDISK] El parametro -size no fue ingresado"
	}

	//return "[MENSAJE-FDISK] finaliza comando, ejecutado exitosamente"
}

func mount(commandArray []string) string {
	// Variables para los valores de los parametros
	val_path := ""
	val_name := ""

	// Banderas para verificar los parametros y ver si se repiten
	band_path := false
	band_name := false

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
				return "[ERROR-MOUNT] El parametro -path ya fue ingresado."
			}
			band_path = true
			// Reemplaza comillas
			val_path = strings.Replace(val_data, "\"", "", 2)
		/* PARAMETRO OBLIGATORIO -> NAME */
		case strings.Contains(data, "name="):
			// Valido si el parametro ya fue ingresado
			if band_name {
				return "[ERROR-MOUNT] El parametro -name ya fue ingresado."
			}

			// Activo la bandera del parametro
			band_name = true

			// Reemplaza comillas
			val_name = strings.Replace(val_data, "\"", "", 2)
		/* PARAMETRO NO VALIDO */
		default:
			return "[ERROR-MOUNT] Parametro no valido."
		}
	}

	if band_path {
		if band_name {
			index_p := functions.Buscar_particion_p_e(val_path, val_name)
			// Si existe
			if index_p != -1 {
				// Apertura del archivo
				f, err := os.OpenFile(val_path, os.O_RDWR, 0660)

				if err == nil {
					mbr_empty := beans.Mbr{}

					// Calculo del tamaño de struct en bytes
					mbr2 := functions.Struct_a_bytes(mbr_empty)
					sstruct := len(mbr2)

					// Lectrura del archivo binario desde el inicio
					lectura := make([]byte, sstruct)
					f.Seek(0, io.SeekStart)
					f.Read(lectura)

					// Conversion de bytes a struct
					master_boot_record := functions.Bytes_a_struct_mbr(lectura)

					// Colocamos la particion ocupada
					copy(master_boot_record.Mbr_partition[index_p].Part_status[:], "2")

					// Conversion de struct a bytes
					mbr_byte := functions.Struct_a_bytes(master_boot_record)

					// Se posiciona al inicio del archivo para guardar la informacion del disco
					f.Seek(0, io.SeekStart)
					f.Write(mbr_byte)
					f.Close()

					// Verifico si la particion ya esta montada
					if beans.Buscar_particion(val_path, val_name, lista_montajes) {
						fmt.Println("[ERROR] La particion ya esta montada...")
					} else {
						// Numero de particion
						num := beans.Buscar_numero(val_path, lista_montajes)
						// Letra de disco
						letra := beans.Buscar_letra(val_path, lista_montajes)
						// Terminacion de su Carnet (los ultimos dos digitos)
						id := "30" + strconv.Itoa(num) + letra

						var n = beans.New_nodo(id, val_path, val_name, letra, num)
						beans.Insertar(n, lista_montajes)
						return "[MENSAJE-MOUNT] Particion montada con exito!"
						//beans.Imprimir_contenido(lista_montajes)
					}
				} else {
					return "[ERROR-MOUNT] No se encuentra el disco."
				}
			} else {
				//Posiblemente logica
				index_p := functions.Buscar_particion_l(val_path, val_name)
				if index_p != -1 {
					// Apertura del archivo
					f, err := os.OpenFile(val_path, os.O_RDWR, 0660)

					if err == nil {
						ebr_empty := beans.Ebr{}

						// Calculo del tamaño de struct en bytes
						ebr2 := functions.Struct_a_bytes(ebr_empty)
						sstruct := len(ebr2)

						// Lectrura del archivo binario desde el inicio
						lectura := make([]byte, sstruct)
						f.Seek(int64(index_p), io.SeekStart)
						f.Read(lectura)

						// Conversion de bytes a struct
						extended_boot_record := functions.Bytes_a_struct_ebr(lectura)

						// Colocamos la particion ocupada
						copy(extended_boot_record.Part_status[:], "2")

						// Conversion de struct a bytes
						mbr_byte := functions.Struct_a_bytes(extended_boot_record)

						// Se posiciona al inicio del archivo para guardar la informacion del disco
						f.Seek(int64(index_p), io.SeekStart)
						f.Write(mbr_byte)
						f.Close()

						// Verifico si la particion ya esta montada
						if beans.Buscar_particion(val_path, val_name, lista_montajes) {
							fmt.Println("[ERROR] La particion ya esta montada...")
						} else {
							// Generacion de id
							// Numero de particion
							num := beans.Buscar_numero(val_path, lista_montajes)
							// Letra de disco
							letra := beans.Buscar_letra(val_path, lista_montajes)
							// Terminacion de su Carnet (los ultimos dos digitos)
							id := "30" + strconv.Itoa(num) + letra

							var n = beans.New_nodo(id, val_path, val_name, letra, num)
							beans.Insertar(n, lista_montajes)
							return "[MENSAJE-MOUNT] Particion montada con exito!"
							//beans.Imprimir_contenido(lista_montajes)
						}
					} else {
						return "[ERROR-MOUNT] No se encuentra el disco."
					}

				} else {
					return "[ERROR-MOUNT] No se encuentra la particion a montar."
				}
			}
		} else {
			return "[ERROR-MOUNT] Parametro -name no definido."
		}
	} else {
		return "[ERROR-MOUNT] Parametro -path no definido..."
	}

	return "[MENSAJE-MOUNT] El comando MOUNT aqui finaliza"
}

func mkfs(commandArray []string) string {

	// Variables para los valores de los parametros
	val_id := ""
	//val_type := ""

	// Banderas para verificar los parametros y ver si se repiten
	band_id := false
	band_type := false

	// Obtengo solo los parametros validos
	for i := 1; i < len(commandArray); i++ {
		aux_data := strings.SplitAfter(commandArray[i], "=")
		data := strings.ToLower(aux_data[0])
		val_data := aux_data[1]

		// Identifica los parametos
		switch {
		/* PARAMETRO OBLIGATORIO -> ID */
		case strings.Contains(data, "id="):
			if band_id {
				return "[ERROR-MKFS] El parametro -path ya fue ingresado."
			}

			// Activo la bandera del parametro
			band_id = true
			val_id = val_data
		/* PARAMETRO OBLIGATORIO -> TYPE */
		case strings.Contains(data, "type="):
			// Valido si el parametro ya fue ingresado
			if band_type {
				return "[ERROR-MKFS] El parametro -name ya fue ingresado."
			}

			// Activo la bandera del parametro
			band_type = true
			//val_type = val_data
		/* PARAMETRO NO VALIDO */
		default:
			return "[ERROR-MKFS] Parametro no valido."
		}
	}

	if band_id {
		var aux = beans.Obtener_nodo(val_id, lista_montajes)
		if aux != nil {
			index := functions.Buscar_particion_p_e(aux.Direccion, aux.Nombre)

			// Si existe la particion
			if index != -1 {
				// Apertura del archivo
				f, err := os.OpenFile(aux.Direccion, os.O_RDWR, 0660)

				if err == nil {
					mbr_empty := beans.Mbr{}

					// Calculo del tamaño de struct en bytes
					mbr2 := functions.Struct_a_bytes(mbr_empty)
					sstruct := len(mbr2)

					// Lectrura del archivo binario desde el inicio
					lectura := make([]byte, sstruct)
					f.Seek(0, io.SeekStart)
					f.Read(lectura)

					// Conversion de bytes a struct
					master_boot_record := functions.Bytes_a_struct_mbr(lectura)

					// Obtengo el inicio
					s_part_start := string(master_boot_record.Mbr_partition[index].Part_start[:])
					// Le quito los caracteres null
					s_part_start = strings.Trim(s_part_start, "\x00")
					inicio, _ := strconv.Atoi(s_part_start)

					// Obtengo el espacio utilizado
					s_part_size := string(master_boot_record.Mbr_partition[index].Part_size[:])
					// Le quito los caracteres null
					s_part_size = strings.Trim(s_part_size, "\x00")
					tamano, _ := strconv.Atoi(s_part_size)

					//return "[MENSAJE] Formateando " + val_type + "\\n"
					f.Close()
					return functions.Formatear_ext2(inicio, tamano, aux.Direccion, "MKFS")

				} else {
					return "[ERROR] No se puede abrir el archivo."
				}

			} else {
				index = functions.Buscar_particion_l(aux.Direccion, aux.Nombre)
				//return "[MENSAJE] Index de la logica" + strconv.Itoa(index) + "\\n"
			}
		} else {
			return "[ERROR-MKFS] No se encuentra ninguna particion montada con ese id."
		}
	} else {
		return "[ERROR-MKFS] El Parametro -id no fue ingresado."
	}

	return "[MENSAJE-MKFS] El comando MKFS finaliza"
}

func rep(commandArray []string) string {
	// Variables para los valores de los parametros
	val_name := ""
	val_path := ""
	val_id := ""

	// Banderas para verificar los parametros y ver si se repiten
	band_name := false
	band_path := false
	band_id := false
	band_ruta := false

	// Obtengo solo los parametros validos
	for i := 1; i < len(commandArray); i++ {
		aux_data := strings.SplitAfter(commandArray[i], "=")
		data := strings.ToLower(aux_data[0])
		val_data := aux_data[1]

		// Identifica los parametos
		switch {
		/* PARAMETRO OBLIGATORIO -> NAME */
		case strings.Contains(data, "name="):
			// Valido si el parametro ya fue ingresado
			if band_name {
				return "[ERROR-REP] El parametro -name ya fue ingresado."
			}

			// Activo la bandera del parametro
			band_name = true

			// Reemplaza comillas
			val_name = strings.Replace(val_data, "\"", "", 2)
		/* PARAMETRO OBLIGATORIO -> PATH */
		case strings.Contains(data, "path="):
			if band_path {
				return "[ERROR-REP] El parametro -path ya fue ingresado."
			}

			// Activo la bandera del parametro
			band_path = true

			// Reemplaza comillas
			val_path = strings.Replace(val_data, "\"", "", 2)
		/* PARAMETRO OBLIGATORIO -> ID */
		case strings.Contains(data, "id="):
			// Valido si el parametro ya fue ingresado
			if band_id {
				return "[ERROR-REP] El parametro -id ya fue ingresado."
			}

			// Activo la bandera del parametro
			band_id = true

			// Reemplaza comillas
			val_id = val_data
		/* PARAMETRO OBLIGATORIO -> RUTA */
		case strings.Contains(data, "ruta="):
			if band_ruta {
				return "[ERROR-REP] El parametro -ruta ya fue ingresado."
			}

			// Activo la bandera del parametro
			band_ruta = true

		/* PARAMETRO NO VALIDO */
		default:
			return "[ERROR-REP] Parametro no valido."
		}
	}

	if band_path {
		if band_name {
			if band_id {
				var aux = beans.Obtener_nodo(val_id, lista_montajes)

				if aux != nil {
					// Reportes validos
					if val_name == "disk" {
						fmt.Println(functions.Graficar_disk(aux.Direccion, val_path, "jpg"))
						return "[MENSAJE-REP] Reporte generado exitosamente."
					} else {
						return "[ERROR-REP] Reporte no valido."
					}
				} else {
					return "[ERROR-REP] No encuentra la particion."
				}
			} else {
				return "[ERROR-REP] El parametro -id no fue ingresado."
			}
		} else {
			return "[ERROR-REP] El parametro -name no fue ingresado."
		}
	} else {
		return "[ERROR-REP] El parametro -path no fue ingresado."
	}

}

/*--------------------------/Comandos--------------------------*/

/*--------------------------Comandos carpetas y archivos--------------------------*/
/*--------------------------/Comandos carpetas y archivos--------------------------*/
