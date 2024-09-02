package functions

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

/*--------------------------Metodos o Funciones--------------------------*/

// Muestra el mensaje de error
func Msg_error(err error, comando string) string {
	return fmt.Sprintf("[ERROR-"+comando+"] %s", err.Error())
}

func Struct_a_bytes(p interface{}) []byte {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(p)
	// ERROR
	if err != nil && err != io.EOF {
		Msg_error(err, "fncGeneral")
	}
	return buf.Bytes()
}

// Crea la Particion Primaria
func crear_particion_primaria(direccion string, nombre string, size int, fit string, unit string) {
	//aux_fit := ""
	aux_unit := ""
	aux_path := direccion
	size_bytes := 1024
	//buffer := "1"

	mbr_empty := mbr{}
	var empty [100]byte

	/* Pendiente */
	// Verifico si tiene Ajuste
	if fit != "" {
		//aux_fit = fit
	} else {
		// Por default es Peor ajuste
		//aux_fit = "w"
	}

	// Verifico si tiene Unidad
	if unit != "" {
		aux_unit = unit

		// *Bytes
		if aux_unit == "b" {
			size_bytes = size
		} else if aux_unit == "k" {
			// *Kilobytes
			size_bytes = size * 1024
		} else {
			// *Megabytes
			size_bytes = size * 1024 * 1024
		}
	} else {
		// Por default Kilobytes
		size_bytes = size * 1024
	}

	// Abro el archivo para lectura con opcion a modificar
	// * OpenFile(name string, flag int, perm FileMode)
	// * O_RDWR -> Lectura y Escritura
	// * 0660 -> Permisos de lectura y escritura
	f, err := os.OpenFile(aux_path, os.O_RDWR, 0660)

	// ERROR
	if err != nil {
		msg_error(err)
	} else {
		// Procede a leer el archivo
		band_particion := false
		num_particion := 0

		// Calculo del tamano de struct en bytes
		mbr2 := struct_a_bytes(mbr_empty)
		sstruct := len(mbr2)

		// Lectrura del archivo binario desde el inicio
		// make -> Crea un slice de bytes con el tamaño indicado (sstruct)
		// ReadAt -> Lee el archivo binario desde la posicion indicada (0) y lo guarda en el slice de bytes
		// Slice de byte es un arreglo de bytes que se puede modificar y con ReadAt se llena con los bytes del archivo
		lectura := make([]byte, sstruct)
		_, err = f.ReadAt(lectura, 0)

		// ERROR
		if err != nil && err != io.EOF {
			msg_error(err)
		}

		// Conversion de bytes a struct
		master_boot_record := bytes_a_struct_mbr(lectura)

		// ERROR
		if err != nil {
			msg_error(err)
		}

		// Si el disco esta creado
		if master_boot_record.Mbr_tamano != empty {
			s_part_start := ""

			// Recorro las 4 particiones
			for i := 0; i < 4; i++ {
				// Antes de comparar limpio la cadena
				s_part_start = string(master_boot_record.Mbr_partition[i].Part_start[:])
				// Le quito los caracteres null
				s_part_start = strings.Trim(s_part_start, "\x00")

				// Verifico si en las particiones hay espacio
				if s_part_start == "-1" && band_particion == false {
					band_particion = true
					num_particion = i
				}
			}

			// Verifico si hay espacio
			if band_particion {
				espacio_usado := 0

				// Recorro las 4 particiones
				for i := 0; i < 4; i++ {
					// Obtengo el espacio utilizado
					s_size := string(master_boot_record.Mbr_partition[i].Part_size[:])
					// Le quito los caracteres null
					s_size = strings.Trim(s_size, "\x00")
					i_size, err := strconv.Atoi(s_size)

					// ERROR
					if err != nil {
						msg_error(err)
					}

					// Le sumo el valor al espacio
					espacio_usado += i_size
				}

				/* Tamaño del disco */

				// Obtengo el tamaño del disco
				s_tamaño_disco := string(master_boot_record.Mbr_tamano[:])
				// Le quito los caracteres null
				s_tamaño_disco = strings.Trim(s_tamaño_disco, "\x00")
				i_tamaño_disco, err2 := strconv.Atoi(s_tamaño_disco)

				// ERROR
				if err2 != nil {
					msg_error(err)
				}

				espacio_disponible := i_tamaño_disco - espacio_usado

				fmt.Println("[ESPACIO DISPONIBLE] ", espacio_disponible, " Bytes")
				fmt.Println("[ESPACIO NECESARIO] ", size_bytes, " Bytes")
				fmt.Println(num_particion)

				// Verifico que haya espacio suficiente
				if espacio_disponible >= size_bytes {
					fmt.Println("Si cumple " + nombre + " !")
				}
			}
		}
		f.Close()
	}
}

/*--------------------------/Metodos o Funciones--------------------------*/
