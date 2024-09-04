package functions

import (
	"backend/beans"
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

/*--------------------------Variables globales--------------------------*/
var graphDot = ""

/*--------------------------Metodos o Funciones--------------------------*/

// Muestra el mensaje de error
func Msg_error(err error, comando string) string {
	return fmt.Sprintf("[ERROR-"+comando+"] %s", err.Error())
}

// Codifica de Struct a []Bytes
func Struct_a_bytes(p interface{}) []byte {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(p)

	// ERROR
	if err != nil && err != io.EOF {
		fmt.Println("[ERROR] ", err)
	}

	return buf.Bytes()
}

// Decodifica de [] Bytes a Struct
func Bytes_a_struct_mbr(s []byte) beans.Mbr {
	p := beans.Mbr{}
	dec := gob.NewDecoder(bytes.NewReader(s))
	err := dec.Decode(&p)

	// ERROR
	if err != nil && err != io.EOF {
		fmt.Println("[ERROR] ", err)
	}

	return p
}

// Decodifica de [] Bytes a Struct
func Bytes_a_struct_ebr(s []byte) beans.Ebr {
	p := beans.Ebr{}
	dec := gob.NewDecoder(bytes.NewReader(s))
	err := dec.Decode(&p)

	// ERROR
	if err != nil && err != io.EOF {
		fmt.Println("[ERROR] ", err)
	}
	return p
}

// Crea la Particion Primaria
func Crear_particion_primaria(direccion string, nombre string, size int, fit string, unit string, metodo string) string {
	aux_fit := ""
	aux_unit := ""
	size_bytes := 1024

	mbr_empty := beans.Mbr{}
	var empty [100]byte

	// Verifico si tiene Ajuste
	if fit != "" {
		aux_fit = fit
	} else {
		// Por default es Peor ajuste
		aux_fit = "w"
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
	f, err := os.OpenFile(direccion, os.O_RDWR, 0660)

	// ERROR
	if err != nil {
		return "[ERROR-" + metodo + "] No existe un disco duro con ese nombre."
	} else {
		// Bandera para ver si hay una particion disponible
		band_particion := false
		// Valor del numero de particion
		num_particion := 0

		// Calculo del tamaño de struct en bytes
		mbr2 := Struct_a_bytes(mbr_empty)
		sstruct := len(mbr2)

		// Lectrura del archivo binario desde el inicio
		lectura := make([]byte, sstruct)
		f.Seek(0, io.SeekStart)
		f.Read(lectura)

		// Conversion de bytes a struct
		master_boot_record := Bytes_a_struct_mbr(lectura)

		// Si el disco esta creado
		if master_boot_record.Mbr_tamano != empty {
			s_part_start := ""

			// Recorro las 4 particiones
			for i := 0; i < 4; i++ {
				// Antes de comparar limpio la cadena
				s_part_start = string(master_boot_record.Mbr_partition[i].Part_start[:])
				s_part_start = strings.Trim(s_part_start, "\x00")

				// Verifico si en las particiones hay espacio
				if s_part_start == "-1" {
					band_particion = true
					num_particion = i
					break
				}
			}

			// Si hay una particion disponible
			if band_particion {
				espacio_usado := 0
				s_part_size := ""
				i_part_size := 0
				s_part_status := ""

				// Recorro las 4 particiones
				for i := 0; i < 4; i++ {
					// Obtengo el espacio utilizado
					s_part_size = string(master_boot_record.Mbr_partition[i].Part_size[:])
					// Le quito los caracteres null
					s_part_size = strings.Trim(s_part_size, "\x00")
					i_part_size, _ = strconv.Atoi(s_part_size)

					// Obtengo el status de la particion
					s_part_status = string(master_boot_record.Mbr_partition[i].Part_status[:])
					// Le quito los caracteres null
					s_part_status = strings.Trim(s_part_status, "\x00")

					if s_part_status != "1" {
						// Le sumo el valor al espacio
						espacio_usado += i_part_size
					}
				}

				/* Tamaño del disco */

				// Obtengo el tamaño del disco
				s_tamaño_disco := string(master_boot_record.Mbr_tamano[:])
				// Le quito los caracteres null
				s_tamaño_disco = strings.Trim(s_tamaño_disco, "\x00")
				i_tamaño_disco, _ := strconv.Atoi(s_tamaño_disco)

				espacio_disponible := i_tamaño_disco - espacio_usado

				// Verifico que haya espacio suficiente
				if espacio_disponible >= size_bytes {
					// Verifico si no existe una particion con ese nombre
					if !existe_particion(direccion, nombre) {
						// Antes de comparar limpio la cadena
						s_dsk_fit := string(master_boot_record.Dsk_fit[:])
						s_dsk_fit = strings.Trim(s_dsk_fit, "\x00")

						/*  Primer Ajuste  */
						if s_dsk_fit == "f" {
							copy(master_boot_record.Mbr_partition[num_particion].Part_type[:], "p")
							copy(master_boot_record.Mbr_partition[num_particion].Part_fit[:], aux_fit)

							// Si esta iniciando
							if num_particion == 0 {
								// Guardo el inicio de la particion y dejo un espacio de separacion
								mbr_empty_byte := Struct_a_bytes(mbr_empty)
								copy(master_boot_record.Mbr_partition[num_particion].Part_start[:], strconv.Itoa(len(mbr_empty_byte)))
							} else {
								// Obtengo el inicio de la particion anterior
								s_part_start_ant := string(master_boot_record.Mbr_partition[num_particion-1].Part_start[:])
								// Le quito los caracteres null
								s_part_start_ant = strings.Trim(s_part_start_ant, "\x00")
								i_part_start_ant, _ := strconv.Atoi(s_part_start_ant)

								// Obtengo el tamaño de la particion anterior
								s_part_size_ant := string(master_boot_record.Mbr_partition[num_particion-1].Part_size[:])
								// Le quito los caracteres null
								s_part_size_ant = strings.Trim(s_part_size_ant, "\x00")
								i_part_size_ant, _ := strconv.Atoi(s_part_size_ant)

								copy(master_boot_record.Mbr_partition[num_particion].Part_start[:], strconv.Itoa(i_part_start_ant+i_part_size_ant))
							}

							copy(master_boot_record.Mbr_partition[num_particion].Part_size[:], strconv.Itoa(size_bytes))
							copy(master_boot_record.Mbr_partition[num_particion].Part_status[:], "0")
							copy(master_boot_record.Mbr_partition[num_particion].Part_name[:], nombre)

							// Se guarda de nuevo el MBR

							// Conversion de struct a bytes
							mbr_byte := Struct_a_bytes(master_boot_record)

							// Se posiciona al inicio del archivo para guardar la informacion del disco
							f.Seek(0, io.SeekStart)
							f.Write(mbr_byte)

							// Obtengo el inicio de la particion
							s_part_start = string(master_boot_record.Mbr_partition[num_particion].Part_start[:])
							// Le quito los caracteres null
							s_part_start = strings.Trim(s_part_start, "\x00")
							i_part_start, _ := strconv.Atoi(s_part_start)

							// Se posiciona en el inicio de la particion
							f.Seek(int64(i_part_start), io.SeekStart)

							// Lo llena de unos
							for i := 0; i < size_bytes; i++ {
								f.Write([]byte{1})
							}

							return "[METODO-" + metodo + "] La Particion primaria fue creada con exito!"
						} else if s_dsk_fit == "b" {
							/*  Mejor Ajuste  */
							best_index := num_particion

							// Variables para conversiones
							s_part_start_act := ""
							s_part_status_act := ""
							s_part_size_act := ""
							i_part_size_act := 0
							s_part_start_best := ""
							i_part_start_best := 0
							s_part_start_best_ant := ""
							i_part_start_best_ant := 0
							s_part_size_best := ""
							i_part_size_best := 0
							s_part_size_best_ant := ""
							i_part_size_best_ant := 0

							for i := 0; i < 4; i++ {
								// Obtengo el inicio de la particion actual
								s_part_start_act = string(master_boot_record.Mbr_partition[i].Part_start[:])
								// Le quito los caracteres null
								s_part_start_act = strings.Trim(s_part_start_act, "\x00")

								// Obtengo el size de la particion actual
								s_part_status_act = string(master_boot_record.Mbr_partition[i].Part_status[:])
								// Le quito los caracteres null
								s_part_status_act = strings.Trim(s_part_status_act, "\x00")

								// Obtengo la posicion de la particion actual
								s_part_size_act = string(master_boot_record.Mbr_partition[i].Part_size[:])
								// Le quito los caracteres null
								s_part_size_act = strings.Trim(s_part_size_act, "\x00")
								i_part_size_act, _ = strconv.Atoi(s_part_size_act)

								if s_part_start_act == "-1" || (s_part_status_act == "1" && i_part_size_act >= size_bytes) {
									if i != num_particion {
										// Obtengo el tamaño de la particion del mejor indice
										s_part_size_best = string(master_boot_record.Mbr_partition[best_index].Part_size[:])
										// Le quito los caracteres null
										s_part_size_best = strings.Trim(s_part_size_best, "\x00")
										i_part_size_best, _ = strconv.Atoi(s_part_size_best)

										// Obtengo la posicion de la particion actual
										s_part_size_act = string(master_boot_record.Mbr_partition[i].Part_size[:])
										// Le quito los caracteres null
										s_part_size_act = strings.Trim(s_part_size_act, "\x00")
										i_part_size_act, _ = strconv.Atoi(s_part_size_act)

										if i_part_size_best > i_part_size_act {
											best_index = i
											break
										}
									}
								}
							}

							// Primaria
							copy(master_boot_record.Mbr_partition[best_index].Part_type[:], "p")
							copy(master_boot_record.Mbr_partition[best_index].Part_fit[:], aux_fit)

							// Si esta iniciando
							if best_index == 0 {
								// Guardo el inicio de la particion y dejo un espacio de separacion
								mbr_empty_byte := Struct_a_bytes(mbr_empty)
								copy(master_boot_record.Mbr_partition[best_index].Part_start[:], strconv.Itoa(len(mbr_empty_byte)))
							} else {
								// Obtengo el inicio de la particion actual
								s_part_start_best_ant = string(master_boot_record.Mbr_partition[best_index-1].Part_start[:])
								// Le quito los caracteres null
								s_part_start_best_ant = strings.Trim(s_part_start_best_ant, "\x00")
								i_part_start_best_ant, _ = strconv.Atoi(s_part_start_best_ant)

								// Obtengo el inicio de la particion actual
								s_part_size_best_ant = string(master_boot_record.Mbr_partition[best_index-1].Part_size[:])
								// Le quito los caracteres null
								s_part_size_best_ant = strings.Trim(s_part_size_best_ant, "\x00")
								i_part_size_best_ant, _ = strconv.Atoi(s_part_size_best_ant)

								copy(master_boot_record.Mbr_partition[best_index].Part_start[:], strconv.Itoa(i_part_start_best_ant+i_part_size_best_ant))
							}

							copy(master_boot_record.Mbr_partition[best_index].Part_size[:], strconv.Itoa(size_bytes))
							copy(master_boot_record.Mbr_partition[best_index].Part_status[:], "0")
							copy(master_boot_record.Mbr_partition[best_index].Part_name[:], nombre)

							// Se guarda de nuevo el MBR

							// Conversion de struct a bytes
							mbr_byte := Struct_a_bytes(master_boot_record)

							// Se posiciona al inicio del archivo para guardar la informacion del disco
							f.Seek(0, io.SeekStart)
							f.Write(mbr_byte)

							// Obtengo el inicio de la particion best
							s_part_start_best = string(master_boot_record.Mbr_partition[best_index].Part_start[:])
							// Le quito los caracteres null
							s_part_start_best = strings.Trim(s_part_start_best, "\x00")
							i_part_start_best, _ = strconv.Atoi(s_part_start_best)

							// Conversion de struct a bytes

							// Se posiciona en el inicio de la particion
							f.Seek(int64(i_part_start_best), io.SeekStart)

							// Lo llena de unos
							for i := 1; i < size_bytes; i++ {
								f.Write([]byte{1})
							}

							return "[MENSAJE-" + metodo + "] La Particion primaria fue creada con exito!"
						} else {
							/*  Peor ajuste  */
							worst_index := num_particion

							// Variables para conversiones
							s_part_start_act := ""
							s_part_status_act := ""
							s_part_size_act := ""
							i_part_size_act := 0
							s_part_start_worst := ""
							i_part_start_worst := 0
							s_part_start_worst_ant := ""
							i_part_start_worst_ant := 0
							s_part_size_worst := ""
							i_part_size_worst := 0
							s_part_size_worst_ant := ""
							i_part_size_worst_ant := 0

							for i := 0; i < 4; i++ {
								// Obtengo el inicio de la particion actual
								s_part_start_act = string(master_boot_record.Mbr_partition[i].Part_start[:])
								// Le quito los caracteres null
								s_part_start_act = strings.Trim(s_part_start_act, "\x00")

								// Obtengo el size de la particion actual
								s_part_status_act = string(master_boot_record.Mbr_partition[i].Part_status[:])
								// Le quito los caracteres null
								s_part_status_act = strings.Trim(s_part_status_act, "\x00")

								// Obtengo la posicion de la particion actual
								s_part_size_act = string(master_boot_record.Mbr_partition[i].Part_size[:])
								// Le quito los caracteres null
								s_part_size_act = strings.Trim(s_part_size_act, "\x00")
								i_part_size_act, _ = strconv.Atoi(s_part_size_act)

								if s_part_start_act == "-1" || (s_part_status_act == "1" && i_part_size_act >= size_bytes) {
									if i != num_particion {
										// Obtengo el tamaño de la particion del mejor indice
										s_part_size_worst = string(master_boot_record.Mbr_partition[worst_index].Part_size[:])
										// Le quito los caracteres null
										s_part_size_worst = strings.Trim(s_part_size_worst, "\x00")
										i_part_size_worst, _ = strconv.Atoi(s_part_size_worst)

										// Obtengo la posicion de la particion actual
										s_part_size_act = string(master_boot_record.Mbr_partition[i].Part_size[:])
										// Le quito los caracteres null
										s_part_size_act = strings.Trim(s_part_size_act, "\x00")
										i_part_size_act, _ = strconv.Atoi(s_part_size_act)

										if i_part_size_worst < i_part_size_act {
											worst_index = i
											break
										}
									}
								}
							}

							// Particiones Primarias
							copy(master_boot_record.Mbr_partition[worst_index].Part_type[:], "p")
							copy(master_boot_record.Mbr_partition[worst_index].Part_fit[:], aux_fit)

							// Se esta iniciando
							if worst_index == 0 {
								// Guardo el inicio de la particion y dejo un espacio de separacion
								mbr_empty_byte := Struct_a_bytes(mbr_empty)
								copy(master_boot_record.Mbr_partition[worst_index].Part_start[:], strconv.Itoa(len(mbr_empty_byte)))
							} else {
								// Obtengo el inicio de la particion anterior
								s_part_start_worst_ant = string(master_boot_record.Mbr_partition[worst_index-1].Part_start[:])
								// Le quito los caracteres null
								s_part_start_worst_ant = strings.Trim(s_part_start_worst_ant, "\x00")
								i_part_start_worst_ant, _ = strconv.Atoi(s_part_start_worst_ant)

								// Obtengo el tamaño de la particion anterior
								s_part_size_worst_ant = string(master_boot_record.Mbr_partition[worst_index-1].Part_size[:])
								// Le quito los caracteres null
								s_part_size_worst_ant = strings.Trim(s_part_size_worst_ant, "\x00")
								i_part_size_worst_ant, _ = strconv.Atoi(s_part_size_worst_ant)

								copy(master_boot_record.Mbr_partition[worst_index].Part_start[:], strconv.Itoa(i_part_start_worst_ant+i_part_size_worst_ant))
							}

							copy(master_boot_record.Mbr_partition[worst_index].Part_size[:], strconv.Itoa(size_bytes))
							copy(master_boot_record.Mbr_partition[worst_index].Part_status[:], "0")
							copy(master_boot_record.Mbr_partition[worst_index].Part_name[:], nombre)

							// Se guarda de nuevo el MBR

							// Conversion de struct a bytes
							mbr_byte := Struct_a_bytes(master_boot_record)

							// Escribe desde el inicio del archivo
							f.Seek(0, io.SeekStart)
							f.Write(mbr_byte)

							// Obtengo el inicio de la particion best
							s_part_start_worst = string(master_boot_record.Mbr_partition[worst_index].Part_start[:])
							// Le quito los caracteres null
							s_part_start_worst = strings.Trim(s_part_start_worst, "\x00")
							i_part_start_worst, _ = strconv.Atoi(s_part_start_worst)

							// Se posiciona en el inicio de la particion
							f.Seek(int64(i_part_start_worst), io.SeekStart)

							// Lo llena de unos
							for i := 1; i < size_bytes; i++ {
								f.Write([]byte{1})
							}
							f.Close()
							return "[MENSAJE-" + metodo + "] La Particion primaria fue creada con exito!"
						}
					} else {
						f.Close()
						return "[ERROR" + metodo + "] Ya existe una particion creada con ese nombre."
					}
				} else {
					f.Close()
					return "[ERROR-" + metodo + "] La particion que desea crear excede el espacio disponible."
				}
			} else {
				f.Close()
				return "[ERROR-" + metodo + "] La suma de particiones primarias y extendidas no debe exceder de 4 particiones. \n " +
					"[MENSAJE-" + metodo + "] Se recomienda eliminar alguna particion para poder crear otra particion primaria o extendida."
			}
		} else {
			f.Close()
			return "[ERROR-" + metodo + "] el disco se encuentra vacio..."
		}

	}
}

// Crea la Particion Extendida
func Crear_particion_extendia(direccion string, nombre string, size int, fit string, unit string, metodo string) string {
	aux_fit := ""
	aux_unit := ""
	size_bytes := 1024

	mbr_empty := beans.Mbr{}
	var empty [100]byte

	// Verifico si tiene Ajuste
	if fit != "" {
		aux_fit = fit
	} else {
		// Por default es Peor ajuste
		aux_fit = "w"
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
	f, err := os.OpenFile(direccion, os.O_RDWR, 0660)

	// ERROR
	if err != nil {
		return Msg_error(err, metodo)
	} else {
		// Procede a leer el archivo
		band_particion := false
		band_extendida := false
		num_particion := 0

		// Calculo del tamaño de struct en bytes
		mbr2 := Struct_a_bytes(mbr_empty)
		sstruct := len(mbr2)

		// Lectrura del archivo binario desde el inicio
		lectura := make([]byte, sstruct)
		f.Seek(0, io.SeekStart)
		f.Read(lectura)

		// Conversion de bytes a struct
		master_boot_record := Bytes_a_struct_mbr(lectura)

		// Si el disco esta creado
		if master_boot_record.Mbr_tamano != empty {
			s_part_type := ""

			// Recorro las 4 particiones
			for i := 0; i < 4; i++ {
				// Antes de comparar limpio la cadena
				s_part_type = string(master_boot_record.Mbr_partition[i].Part_type[:])
				s_part_type = strings.Trim(s_part_type, "\x00")

				if s_part_type == "e" {
					band_extendida = true
					break
				}
			}

			// Si no es extendida
			if !band_extendida {
				s_part_start := ""
				s_part_status := ""
				s_part_size := ""
				i_part_size := 0

				// Recorro las 4 particiones
				for i := 0; i < 4; i++ {
					// Antes de comparar limpio la cadena
					s_part_start = string(master_boot_record.Mbr_partition[i].Part_start[:])
					s_part_start = strings.Trim(s_part_start, "\x00")

					s_part_status = string(master_boot_record.Mbr_partition[i].Part_status[:])
					s_part_status = strings.Trim(s_part_status, "\x00")

					s_part_size = string(master_boot_record.Mbr_partition[i].Part_size[:])
					s_part_size = strings.Trim(s_part_size, "\x00")
					i_part_size, _ = strconv.Atoi(s_part_size)

					// Verifica si existe una particion disponible
					if s_part_start == "-1" || (s_part_status == "1" && i_part_size >= size_bytes) {
						band_particion = true
						num_particion = i
						break
					}
				}

				// Si hay una particion
				if band_particion {
					espacio_usado := 0

					// Recorro las 4 particiones
					for i := 0; i < 4; i++ {
						s_part_status = string(master_boot_record.Mbr_partition[i].Part_status[:])
						s_part_status = strings.Trim(s_part_status, "\x00")

						if s_part_status != "1" {
							// Obtengo el espacio utilizado
							s_part_size = string(master_boot_record.Mbr_partition[i].Part_size[:])
							// Le quito los caracteres null
							s_part_size = strings.Trim(s_part_size, "\x00")
							i_part_size, _ = strconv.Atoi(s_part_size)

							// Le sumo el valor al espacio
							espacio_usado += i_part_size
						}
					}

					/* Tamaño del disco */

					// Obtengo el tamaño del disco
					s_tamaño_disco := string(master_boot_record.Mbr_tamano[:])
					// Le quito los caracteres null
					s_tamaño_disco = strings.Trim(s_tamaño_disco, "\x00")
					i_tamaño_disco, _ := strconv.Atoi(s_tamaño_disco)

					espacio_disponible := i_tamaño_disco - espacio_usado

					//fmt.Println("[ESPACIO DISPONIBLE] ", espacio_disponible, " Bytes")
					//fmt.Println("[ESPACIO NECESARIO] ", size_bytes, " Bytes")

					// Verifico que haya espacio suficiente
					if espacio_disponible >= size_bytes {
						// Verifico si no existe una particion con ese nombre
						if !existe_particion(direccion, nombre) {
							// Antes de comparar limpio la cadena
							s_dsk_fit := string(master_boot_record.Dsk_fit[:])
							s_dsk_fit = strings.Trim(s_dsk_fit, "\x00")

							/*  Primer Ajuste  */
							if s_dsk_fit == "f" {
								copy(master_boot_record.Mbr_partition[num_particion].Part_type[:], "e")
								copy(master_boot_record.Mbr_partition[num_particion].Part_fit[:], aux_fit)

								// Si esta iniciando
								if num_particion == 0 {
									// Guardo el inicio de la particion y dejo un espacio de separacion
									mbr_empty_byte := Struct_a_bytes(mbr_empty)
									copy(master_boot_record.Mbr_partition[num_particion].Part_start[:], strconv.Itoa(len(mbr_empty_byte)))
								} else {
									// Obtengo el inicio de la particion anterior
									s_part_start_ant := string(master_boot_record.Mbr_partition[num_particion-1].Part_start[:])
									// Le quito los caracteres null
									s_part_start_ant = strings.Trim(s_part_start_ant, "\x00")
									i_part_start_ant, _ := strconv.Atoi(s_part_start_ant)

									// Obtengo el tamaño de la particion anterior
									s_part_size_ant := string(master_boot_record.Mbr_partition[num_particion-1].Part_size[:])
									// Le quito los caracteres null
									s_part_size_ant = strings.Trim(s_part_size_ant, "\x00")
									i_part_size_ant, _ := strconv.Atoi(s_part_size_ant)

									copy(master_boot_record.Mbr_partition[num_particion].Part_start[:], strconv.Itoa(i_part_start_ant+i_part_size_ant))
								}

								copy(master_boot_record.Mbr_partition[num_particion].Part_size[:], strconv.Itoa(size_bytes))
								copy(master_boot_record.Mbr_partition[num_particion].Part_status[:], "0")
								copy(master_boot_record.Mbr_partition[num_particion].Part_name[:], nombre)

								// Se guarda de nuevo el MBR

								// Conversion de struct a bytes
								mbr_byte := Struct_a_bytes(master_boot_record)

								// Escribe en la posicion inicial del archivo
								f.Seek(0, io.SeekStart)
								f.Write(mbr_byte)

								// Obtengo el tamaño de la particion
								s_part_start = string(master_boot_record.Mbr_partition[num_particion].Part_start[:])
								// Le quito los caracteres null
								s_part_start = strings.Trim(s_part_start, "\x00")
								i_part_start, _ := strconv.Atoi(s_part_start)

								// Se posiciona en el inicio de la particion
								f.Seek(int64(i_part_start), io.SeekStart)

								extended_boot_record := beans.Ebr{}
								copy(extended_boot_record.Part_fit[:], aux_fit)
								copy(extended_boot_record.Part_status[:], "0")
								copy(extended_boot_record.Part_start[:], s_part_start)
								copy(extended_boot_record.Part_size[:], "0")
								copy(extended_boot_record.Part_next[:], "-1")
								copy(extended_boot_record.Part_name[:], "")
								ebr_byte := Struct_a_bytes(extended_boot_record)
								f.Write(ebr_byte)

								// Lo llena de unos
								for i := 0; i < (size_bytes - len(ebr_byte)); i++ {
									f.Write([]byte{1})
								}
								f.Close()
								return "[MENSAJE-" + metodo + "] La Particion extendida fue creada con exito!"
							} else if s_dsk_fit == "b" {
								/*  Mejor Ajuste  */
								best_index := num_particion

								// Variables para conversiones
								s_part_start_act := ""
								s_part_status_act := ""
								s_part_size_act := ""
								i_part_size_act := 0
								s_part_start_best := ""
								i_part_start_best := 0
								s_part_start_best_ant := ""
								i_part_start_best_ant := 0
								s_part_size_best := ""
								i_part_size_best := 0
								s_part_size_best_ant := ""
								i_part_size_best_ant := 0

								for i := 0; i < 4; i++ {
									// Obtengo el inicio de la particion actual
									s_part_start_act = string(master_boot_record.Mbr_partition[i].Part_start[:])
									// Le quito los caracteres null
									s_part_start_act = strings.Trim(s_part_start_act, "\x00")

									// Obtengo el size de la particion actual
									s_part_status_act = string(master_boot_record.Mbr_partition[i].Part_status[:])
									// Le quito los caracteres null
									s_part_status_act = strings.Trim(s_part_status_act, "\x00")

									// Obtengo la posicion de la particion actual
									s_part_size_act = string(master_boot_record.Mbr_partition[i].Part_size[:])
									// Le quito los caracteres null
									s_part_size_act = strings.Trim(s_part_size_act, "\x00")
									i_part_size_act, _ = strconv.Atoi(s_part_size_act)

									if s_part_start_act == "-1" || (s_part_status_act == "1" && i_part_size_act >= size_bytes) {
										if i != num_particion {
											// Obtengo el tamaño de la particion del mejor indice
											s_part_size_best = string(master_boot_record.Mbr_partition[best_index].Part_size[:])
											// Le quito los caracteres null
											s_part_size_best = strings.Trim(s_part_size_best, "\x00")
											i_part_size_best, _ = strconv.Atoi(s_part_size_best)

											// Obtengo la posicion de la particion actual
											s_part_size_act = string(master_boot_record.Mbr_partition[i].Part_size[:])
											// Le quito los caracteres null
											s_part_size_act = strings.Trim(s_part_size_act, "\x00")
											i_part_size_act, _ = strconv.Atoi(s_part_size_act)

											if i_part_size_best > i_part_size_act {
												best_index = i
												break
											}
										}
									}
								}

								// Extendida
								copy(master_boot_record.Mbr_partition[best_index].Part_type[:], "e")
								copy(master_boot_record.Mbr_partition[best_index].Part_fit[:], aux_fit)

								// Si esta iniciando
								if best_index == 0 {
									// Guardo el inicio de la particion y dejo un espacio de separacion
									mbr_empty_byte := Struct_a_bytes(mbr_empty)
									copy(master_boot_record.Mbr_partition[best_index].Part_start[:], strconv.Itoa(len(mbr_empty_byte)))
								} else {
									// Obtengo el inicio de la particion actual
									s_part_start_best_ant = string(master_boot_record.Mbr_partition[best_index-1].Part_start[:])
									// Le quito los caracteres null
									s_part_start_best_ant = strings.Trim(s_part_start_best_ant, "\x00")
									i_part_start_best_ant, _ = strconv.Atoi(s_part_start_best_ant)

									// Obtengo el inicio de la particion actual
									s_part_size_best_ant = string(master_boot_record.Mbr_partition[best_index-1].Part_size[:])
									// Le quito los caracteres null
									s_part_size_best_ant = strings.Trim(s_part_size_best_ant, "\x00")
									i_part_size_best_ant, _ = strconv.Atoi(s_part_size_best_ant)

									copy(master_boot_record.Mbr_partition[best_index].Part_start[:], strconv.Itoa(i_part_start_best_ant+i_part_size_best_ant))
								}

								copy(master_boot_record.Mbr_partition[best_index].Part_size[:], strconv.Itoa(size_bytes))
								copy(master_boot_record.Mbr_partition[best_index].Part_status[:], "0")
								copy(master_boot_record.Mbr_partition[best_index].Part_name[:], nombre)

								// Se guarda de nuevo el MBR

								// Conversion de struct a bytes
								mbr_byte := Struct_a_bytes(master_boot_record)

								// Se escribe al inicio del archivo
								f.Seek(0, io.SeekStart)
								f.Write(mbr_byte)

								// Obtengo el inicio de la particion best
								s_part_start_best = string(master_boot_record.Mbr_partition[best_index].Part_start[:])
								// Le quito los caracteres null
								s_part_start_best = strings.Trim(s_part_start_best, "\x00")
								i_part_start_best, _ = strconv.Atoi(s_part_start_best)

								// Se posiciona en el inicio de la particion
								f.Seek(int64(i_part_start_best), io.SeekStart)

								extended_boot_record := beans.Ebr{}
								copy(extended_boot_record.Part_fit[:], aux_fit)
								copy(extended_boot_record.Part_status[:], "0")
								copy(extended_boot_record.Part_start[:], s_part_start_best)
								copy(extended_boot_record.Part_size[:], "0")
								copy(extended_boot_record.Part_next[:], "-1")
								copy(extended_boot_record.Part_name[:], "")
								ebr_byte := Struct_a_bytes(extended_boot_record)
								f.Write(ebr_byte)

								// Lo llena de unos
								for i := 0; i < (size_bytes - len(ebr_byte)); i++ {
									f.Write([]byte{1})
								}
								f.Close()
								return "[MENSAJE-" + metodo + "] La Particion extendida fue creada con exito!"
							} else {
								/*  Peor ajuste  */
								worst_index := num_particion

								// Variables para conversiones
								s_part_start_act := ""
								s_part_status_act := ""
								s_part_size_act := ""
								i_part_size_act := 0
								s_part_start_worst := ""
								i_part_start_worst := 0
								s_part_start_worst_ant := ""
								i_part_start_worst_ant := 0
								s_part_size_worst := ""
								i_part_size_worst := 0
								s_part_size_worst_ant := ""
								i_part_size_worst_ant := 0

								for i := 0; i < 4; i++ {
									// Obtengo el inicio de la particion actual
									s_part_start_act = string(master_boot_record.Mbr_partition[i].Part_start[:])
									// Le quito los caracteres null
									s_part_start_act = strings.Trim(s_part_start_act, "\x00")

									// Obtengo el size de la particion actual
									s_part_status_act = string(master_boot_record.Mbr_partition[i].Part_status[:])
									// Le quito los caracteres null
									s_part_status_act = strings.Trim(s_part_status_act, "\x00")

									// Obtengo la posicion de la particion actual
									s_part_size_act = string(master_boot_record.Mbr_partition[i].Part_size[:])
									// Le quito los caracteres null
									s_part_size_act = strings.Trim(s_part_size_act, "\x00")
									i_part_size_act, _ = strconv.Atoi(s_part_size_act)

									if s_part_start_act == "-1" || (s_part_status_act == "1" && i_part_size_act >= size_bytes) {
										if i != num_particion {
											// Obtengo el tamaño de la particion del mejor indice
											s_part_size_worst = string(master_boot_record.Mbr_partition[worst_index].Part_size[:])
											// Le quito los caracteres null
											s_part_size_worst = strings.Trim(s_part_size_worst, "\x00")
											i_part_size_worst, _ = strconv.Atoi(s_part_size_worst)

											// Obtengo la posicion de la particion actual
											s_part_size_act = string(master_boot_record.Mbr_partition[i].Part_size[:])
											// Le quito los caracteres null
											s_part_size_act = strings.Trim(s_part_size_act, "\x00")
											i_part_size_act, _ = strconv.Atoi(s_part_size_act)

											if i_part_size_worst < i_part_size_act {
												worst_index = i
												break
											}
										}
									}
								}

								// Particiones Extendidas
								copy(master_boot_record.Mbr_partition[worst_index].Part_type[:], "e")
								copy(master_boot_record.Mbr_partition[worst_index].Part_fit[:], aux_fit)

								// Se esta iniciando
								if worst_index == 0 {
									// Guardo el inicio de la particion y dejo un espacio de separacion
									mbr_empty_byte := Struct_a_bytes(mbr_empty)
									copy(master_boot_record.Mbr_partition[worst_index].Part_start[:], strconv.Itoa(len(mbr_empty_byte)))
								} else {
									// Obtengo el inicio de la particion actual
									s_part_start_worst_ant = string(master_boot_record.Mbr_partition[worst_index-1].Part_start[:])
									// Le quito los caracteres null
									s_part_start_worst_ant = strings.Trim(s_part_start_worst_ant, "\x00")
									i_part_start_worst_ant, _ = strconv.Atoi(s_part_start_worst_ant)

									// Obtengo el inicio de la particion actual
									s_part_size_worst_ant = string(master_boot_record.Mbr_partition[worst_index-1].Part_size[:])
									// Le quito los caracteres null
									s_part_size_worst_ant = strings.Trim(s_part_size_worst_ant, "\x00")
									i_part_size_worst_ant, _ = strconv.Atoi(s_part_size_worst_ant)

									copy(master_boot_record.Mbr_partition[worst_index].Part_start[:], strconv.Itoa(i_part_start_worst_ant+i_part_size_worst_ant))
								}

								copy(master_boot_record.Mbr_partition[worst_index].Part_size[:], strconv.Itoa(size_bytes))
								copy(master_boot_record.Mbr_partition[worst_index].Part_status[:], "0")
								copy(master_boot_record.Mbr_partition[worst_index].Part_name[:], nombre)

								// Se guarda de nuevo el MBR

								// Conversion de struct a bytes
								mbr_byte := Struct_a_bytes(master_boot_record)

								// Se escribe desde el inicio del archivo
								f.Seek(0, io.SeekStart)
								f.Write(mbr_byte)

								// Obtengo el inicio de la particion best
								s_part_start_worst = string(master_boot_record.Mbr_partition[worst_index].Part_start[:])
								// Le quito los caracteres null
								s_part_start_worst = strings.Trim(s_part_start_worst, "\x00")
								i_part_start_worst, _ = strconv.Atoi(s_part_start_worst)

								// Se posiciona en el inicio de la particion
								f.Seek(int64(i_part_start_worst), io.SeekStart)

								extended_boot_record := beans.Ebr{}
								copy(extended_boot_record.Part_fit[:], aux_fit)
								copy(extended_boot_record.Part_status[:], "0")
								copy(extended_boot_record.Part_start[:], s_part_start_worst)
								copy(extended_boot_record.Part_size[:], "0")
								copy(extended_boot_record.Part_next[:], "-1")
								copy(extended_boot_record.Part_name[:], "")
								ebr_byte := Struct_a_bytes(extended_boot_record)
								f.Write(ebr_byte)

								// Lo llena de unos
								for i := 0; i < (size_bytes - len(ebr_byte)); i++ {
									f.Write([]byte{1})
								}
								f.Close()
								return "[MENSAJE-" + metodo + "] La Particion extendida fue creada con exito!"
							}
						} else {
							f.Close()
							return "[ERROR-" + metodo + "] Ya existe una particion creada con ese nombre."
						}
					} else {
						f.Close()
						return "[ERROR-" + metodo + "] La particion que desea crear excede el espacio disponible."
					}
				} else {
					f.Close()
					return "[ERROR-" + metodo + "] La suma de particiones primarias y extendidas no debe exceder de 4 particiones. \n" +
						"[ERROR-" + metodo + "] Se recomienda eliminar alguna particion para poder crear otra particion primaria o extendida"
				}
			} else {
				f.Close()
				return "[ERROR-" + metodo + "] Solo puede haber una particion extendida por disco."
			}
		} else {
			f.Close()
			return "[ERROR-" + metodo + "] el disco se encuentra vacio."
		}
	}
}

// Crea la Particion Logica
func Crear_particion_logica(direccion string, nombre string, size int, fit string, unit string, metodo string) string {
	aux_fit := ""
	aux_unit := ""
	size_bytes := 1024

	mbr_empty := beans.Mbr{}
	ebr_empty := beans.Ebr{}
	var empty [100]byte

	// Verifico si tiene Ajuste
	if fit != "" {
		aux_fit = fit
	} else {
		// Por default es Peor ajuste
		aux_fit = "w"
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
	f, err := os.OpenFile(direccion, os.O_RDWR, 0660)

	// ERROR
	if err != nil {
		return "[ERROR-" + metodo + "] No existe el disco duro con ese nombre."
	} else {
		// Calculo del tamaño de struct en bytes
		mbr2 := Struct_a_bytes(mbr_empty)
		sstruct := len(mbr2)

		// Lectrura del archivo binario desde el inicio
		lectura := make([]byte, sstruct)
		f.Seek(0, io.SeekStart)
		f.Read(lectura)

		// Conversion de bytes a struct
		master_boot_record := Bytes_a_struct_mbr(lectura)

		// Si el disco esta creado
		if master_boot_record.Mbr_tamano != empty {
			s_part_type := ""
			num_extendida := -1

			// Recorro las 4 particiones
			for i := 0; i < 4; i++ {
				// Antes de comparar limpio la cadena
				s_part_type = string(master_boot_record.Mbr_partition[i].Part_type[:])
				s_part_type = strings.Trim(s_part_type, "\x00")

				if s_part_type == "e" {
					num_extendida = i
					break
				}
			}

			if !existe_particion(direccion, nombre) {
				if num_extendida != -1 {
					s_part_start := string(master_boot_record.Mbr_partition[num_extendida].Part_start[:])
					s_part_start = strings.Trim(s_part_start, "\x00")
					i_part_start, _ := strconv.Atoi(s_part_start)

					cont := i_part_start

					// Se posiciona en el inicio de la particion
					f.Seek(int64(cont), io.SeekStart)

					// Calculo del tamaño de struct en bytes
					ebr2 := Struct_a_bytes(ebr_empty)
					sstruct := len(ebr2)

					// Lectrura del archivo binario desde el inicio
					lectura := make([]byte, sstruct)
					f.Read(lectura)

					// Conversion de bytes a struct
					extended_boot_record := Bytes_a_struct_ebr(lectura)

					// Obtencion de datos
					s_part_size_ext := string(extended_boot_record.Part_size[:])
					s_part_size_ext = strings.Trim(s_part_size_ext, "\x00")

					if s_part_size_ext == "0" {
						// Obtencion de datos
						s_part_size := string(master_boot_record.Mbr_partition[num_extendida].Part_size[:])
						s_part_size = strings.Trim(s_part_size, "\x00")
						i_part_size, _ := strconv.Atoi(s_part_size)

						//fmt.Println("[ESPACIO DISPONIBLE] ", i_part_size, " Bytes")
						//fmt.Println("[ESPACIO NECESARIO] ", size_bytes, " Bytes")

						// Si excede el tamaño de la extendida
						if i_part_size < size_bytes {
							f.Close()
							return "[ERROR-" + metodo + "] La particion logica a crear excede el espacio disponible de la particion extendida."
						} else {
							copy(extended_boot_record.Part_status[:], "0")
							copy(extended_boot_record.Part_fit[:], aux_fit)

							// Posicion actual en el archivo
							pos_actual, _ := f.Seek(0, io.SeekCurrent)
							ebr_empty_byte := Struct_a_bytes(ebr_empty)

							copy(extended_boot_record.Part_start[:], strconv.Itoa(int(pos_actual)-len(ebr_empty_byte)))
							copy(extended_boot_record.Part_size[:], strconv.Itoa(size_bytes))
							copy(extended_boot_record.Part_next[:], "-1")
							copy(extended_boot_record.Part_name[:], nombre)

							// Obtencion de datos
							s_part_start := string(master_boot_record.Mbr_partition[num_extendida].Part_start[:])
							s_part_start = strings.Trim(s_part_start, "\x00")
							i_part_start, _ := strconv.Atoi(s_part_start)

							// Se posiciona en el inicio de la particion
							ebr_byte := Struct_a_bytes(extended_boot_record)
							f.Seek(int64(i_part_start), io.SeekStart)
							f.Write(ebr_byte)
							f.Close()
							return "[MENSAJE-" + metodo + "] La Particion logica fue creada con exito!"
						}
					} else {
						// Obtencion de datos
						s_part_size := string(master_boot_record.Mbr_partition[num_extendida].Part_size[:])
						s_part_size = strings.Trim(s_part_size, "\x00")
						i_part_size, _ := strconv.Atoi(s_part_size)

						// Obtencion de datos
						s_part_start := string(master_boot_record.Mbr_partition[num_extendida].Part_start[:])
						s_part_start = strings.Trim(s_part_start, "\x00")
						i_part_start, _ := strconv.Atoi(s_part_start)

						//fmt.Println("[ESPACIO DISPONIBLE] ", i_part_size+i_part_start, " Bytes")
						//fmt.Println("[ESPACIO NECESARIO] ", size_bytes, " Bytes")

						// Obtencion de datos
						s_part_next := string(extended_boot_record.Part_next[:])
						s_part_next = strings.Trim(s_part_next, "\x00")
						i_part_next, _ := strconv.Atoi(s_part_next)

						pos_actual, _ := f.Seek(0, io.SeekCurrent)

						for (i_part_next != -1) && (int(pos_actual) < (i_part_size + i_part_start)) {
							// Se posiciona en el inicio de la particion
							f.Seek(int64(i_part_next), io.SeekStart)

							// Calculo del tamaño de struct en bytes
							ebr2 := Struct_a_bytes(ebr_empty)
							sstruct := len(ebr2)

							// Lectrura del archivo binario desde el inicio
							lectura := make([]byte, sstruct)
							f.Read(lectura)

							// Posicion actual
							pos_actual, _ = f.Seek(0, io.SeekCurrent)

							// Conversion de bytes a struct
							extended_boot_record = Bytes_a_struct_ebr(lectura)

							if extended_boot_record.Part_next == empty {
								break
							}

							// Obtencion de datos
							s_part_next = string(extended_boot_record.Part_next[:])
							s_part_next = strings.Trim(s_part_next, "\x00")
							i_part_next, _ = strconv.Atoi(s_part_next)
						}

						// Obtencion de datos
						s_part_start_ext := string(extended_boot_record.Part_start[:])
						s_part_start_ext = strings.Trim(s_part_start_ext, "\x00")
						i_part_start_ext, _ := strconv.Atoi(s_part_start_ext)

						// Obtencion de datos
						s_part_size_ext := string(extended_boot_record.Part_size[:])
						s_part_size_ext = strings.Trim(s_part_size_ext, "\x00")
						i_part_size_ext, _ := strconv.Atoi(s_part_size_ext)

						// Obtencion de datos
						s_part_size_mbr := string(master_boot_record.Mbr_partition[num_extendida].Part_size[:])
						s_part_size_mbr = strings.Trim(s_part_size_mbr, "\x00")
						i_part_size_mbr, _ := strconv.Atoi(s_part_size_mbr)

						// Obtencion de datos
						s_part_start_mbr := string(master_boot_record.Mbr_partition[num_extendida].Part_start[:])
						s_part_start_mbr = strings.Trim(s_part_start_mbr, "\x00")
						i_part_start_mbr, _ := strconv.Atoi(s_part_start_mbr)

						espacio_necesario := i_part_start_ext + i_part_size_ext + size_bytes

						if espacio_necesario <= (i_part_size_mbr + i_part_start_mbr) {
							copy(extended_boot_record.Part_next[:], strconv.Itoa(i_part_start_ext+i_part_size_ext))

							// Escribo el nedxto del ultimo ebr
							pos_actual, _ = f.Seek(0, io.SeekCurrent)
							ebr_byte := Struct_a_bytes(extended_boot_record)
							// Escribo el next del ultimo EBR
							f.Seek(int64(int(pos_actual)-len(ebr_byte)), io.SeekStart)
							f.Write(ebr_byte)

							// Escribo el nuevo EBR
							f.Seek(int64(i_part_start_ext+i_part_size_ext), io.SeekStart)
							copy(extended_boot_record.Part_status[:], "0")
							copy(extended_boot_record.Part_fit[:], aux_fit)
							// Posicion actual del archivo
							pos_actual, _ = f.Seek(0, io.SeekCurrent)
							copy(extended_boot_record.Part_start[:], strconv.Itoa(int(pos_actual)))
							copy(extended_boot_record.Part_size[:], strconv.Itoa(size_bytes))
							copy(extended_boot_record.Part_next[:], "-1")
							copy(extended_boot_record.Part_name[:], nombre)
							ebr_byte = Struct_a_bytes(extended_boot_record)
							f.Write(ebr_byte)
							f.Close()
							return "[MENSAJE-" + metodo + "] La Particion logica fue creada con exito!"
						} else {
							f.Close()
							return "[ERROR-" + metodo + "] La particion logica a crear excede el espacio disponible de la particion extendida."
						}
					}
				} else {
					f.Close()
					return "[ERROR-" + metodo + "] No se puede crear una particion logica si no hay una extendida."
				}
			} else {
				f.Close()
				return "[ERROR-" + metodo + "] Ya existe una particion con ese nombre."
			}
		} else {
			f.Close()
			return "[ERROR-" + metodo + "] el disco se encuentra vacio."
		}
	}
}

// Verifica si el nombre de la particion esta disponible
func existe_particion(direccion string, nombre string) bool {
	extendida := -1
	mbr_empty := beans.Mbr{}
	ebr_empty := beans.Ebr{}
	var empty [100]byte

	// Abro el archivo para lectura con opcion a modificar
	f, err := os.OpenFile(direccion, os.O_RDWR, 0660)

	if err == nil {
		// Procedo a leer el archivo

		// Calculo del tamaño de struct en bytes
		mbr2 := Struct_a_bytes(mbr_empty)
		sstruct := len(mbr2)

		// Lectrura del archivo binario desde el inicio
		lectura := make([]byte, sstruct)
		f.Seek(0, io.SeekStart)
		f.Read(lectura)

		// Conversion de bytes a struct
		master_boot_record := Bytes_a_struct_mbr(lectura)

		// Si el disco esta creado
		if master_boot_record.Mbr_tamano != empty {
			s_part_name := ""
			s_part_type := ""

			// Recorro las 4 particiones
			for i := 0; i < 4; i++ {
				// Antes de comparar limpio la cadena
				s_part_name = string(master_boot_record.Mbr_partition[i].Part_name[:])
				s_part_name = strings.Trim(s_part_name, "\x00")

				// Verifico si ya existe una particion con ese nombre
				if s_part_name == nombre {
					f.Close()
					return true
				}

				// Antes de comparar limpio la cadena
				s_part_type = string(master_boot_record.Mbr_partition[i].Part_type[:])
				s_part_type = strings.Trim(s_part_type, "\x00")

				// Verifico si de tipo extendida
				if s_part_type == "e" {
					extendida = i
				}
			}

			// Si es extendida
			if extendida != -1 {
				// Obtengo el inicio de la particion
				s_part_start := string(master_boot_record.Mbr_partition[extendida].Part_start[:])
				// Le quito los caracteres null
				s_part_start = strings.Trim(s_part_start, "\x00")
				i_part_start, _ := strconv.Atoi(s_part_start)

				// Obtengo el espacio de la partcion
				s_part_size := string(master_boot_record.Mbr_partition[extendida].Part_size[:])
				// Le quito los caracteres null
				s_part_size = strings.Trim(s_part_size, "\x00")
				i_part_size, _ := strconv.Atoi(s_part_size)

				// Calculo del tamaño de struct en bytes
				ebr2 := Struct_a_bytes(ebr_empty)
				sstruct := len(ebr2)

				// Lectrura de conjunto de bytes en archivo binario
				lectura := make([]byte, sstruct)
				// Lee a partir del inicio de la particion
				n_leidos, _ := f.Read(lectura)

				// Posicion actual en el archivo
				f.Seek(int64(i_part_start), io.SeekStart)

				// Posicion actual en el archivo
				pos_actual, _ := f.Seek(0, io.SeekCurrent)

				// Lectrura de conjunto de bytes desde el inicio de la particion
				for n_leidos != 0 && (pos_actual < int64(i_part_size+i_part_start)) {
					// Lectrura de conjunto de bytes en archivo binario
					lectura := make([]byte, sstruct)
					// Lee a partir del inicio de la particion
					n_leidos, _ = f.Read(lectura)

					// Posicion actual en el archivo
					pos_actual, _ = f.Seek(0, io.SeekCurrent)

					// Conversion de bytes a struct
					extended_boot_record := Bytes_a_struct_ebr(lectura)

					if extended_boot_record.Part_size == empty {
						break
					} else {
						// Antes de comparar limpio la cadena
						s_part_name = string(extended_boot_record.Part_name[:])
						s_part_name = strings.Trim(s_part_name, "\x00")

						// Verifico si ya existe una particion con ese nombre
						if s_part_name == nombre {
							f.Close()
							return true
						}

						// Obtengo el espacio utilizado
						s_part_next := string(extended_boot_record.Part_next[:])
						// Le quito los caracteres null
						s_part_next = strings.Trim(s_part_next, "\x00")

						// Si ya termino
						if s_part_next != "-1" {
							f.Close()
							return false
						}
					}
				}
			}
		} /*else {
			fmt.Println("[ERROR] el disco se encuentra vacio."
		}*/
	}
	//f.Close()
	return false
}

// Busca particiones Primarias o Extendidas
func Buscar_particion_p_e(direccion string, nombre string) int {
	// Apertura del archivo
	f, err := os.OpenFile(direccion, os.O_RDWR, 0660)

	if err == nil {
		mbr_empty := beans.Mbr{}

		// Calculo del tamaño de struct en bytes
		mbr2 := Struct_a_bytes(mbr_empty)
		sstruct := len(mbr2)

		// Lectrura del archivo binario desde el inicio
		lectura := make([]byte, sstruct)
		f.Seek(0, io.SeekStart)
		f.Read(lectura)

		// Conversion de bytes a struct
		master_boot_record := Bytes_a_struct_mbr(lectura)

		s_part_status := ""
		s_part_name := ""

		// Recorro las 4 particiones
		for i := 0; i < 4; i++ {
			// Antes de comparar limpio la cadena
			s_part_status = string(master_boot_record.Mbr_partition[i].Part_status[:])
			s_part_status = strings.Trim(s_part_status, "\x00")

			if s_part_status != "1" {
				// Antes de comparar limpio la cadena
				s_part_name = string(master_boot_record.Mbr_partition[i].Part_name[:])
				s_part_name = strings.Trim(s_part_name, "\x00")
				if s_part_name == nombre {
					return i
				}
			}

		}
	}

	return -1
}

// Busca particiones Logicas
func Buscar_particion_l(direccion string, nombre string) int {
	// Apertura del archivo
	f, err := os.OpenFile(direccion, os.O_RDWR, 0660)

	if err == nil {
		extendida := -1
		mbr_empty := beans.Mbr{}

		// Calculo del tamaño de struct en bytes
		mbr2 := Struct_a_bytes(mbr_empty)
		sstruct := len(mbr2)

		// Lectrura del archivo binario desde el inicio
		lectura := make([]byte, sstruct)
		f.Seek(0, io.SeekStart)
		f.Read(lectura)

		// Conversion de bytes a struct
		master_boot_record := Bytes_a_struct_mbr(lectura)

		s_part_type := ""

		// Recorro las 4 particiones
		for i := 0; i < 4; i++ {
			// Antes de comparar limpio la cadena
			s_part_type = string(master_boot_record.Mbr_partition[i].Part_type[:])
			s_part_type = strings.Trim(s_part_type, "\x00")

			if s_part_type != "e" {
				extendida = i
				break
			}
		}

		// Si hay extendida
		if extendida != -1 {
			ebr_empty := beans.Ebr{}

			ebr2 := Struct_a_bytes(ebr_empty)
			sstruct := len(ebr2)

			// Lectrura del archivo binario desde el inicio
			lectura := make([]byte, sstruct)

			s_part_start := string(master_boot_record.Mbr_partition[extendida].Part_start[:])
			s_part_start = strings.Trim(s_part_start, "\x00")
			i_part_start, _ := strconv.Atoi(s_part_start)

			f.Seek(int64(i_part_start), io.SeekStart)

			n_leidos, _ := f.Read(lectura)
			pos_actual, _ := f.Seek(0, io.SeekCurrent)

			s_part_size := string(master_boot_record.Mbr_partition[extendida].Part_start[:])
			s_part_size = strings.Trim(s_part_size, "\x00")
			i_part_size, _ := strconv.Atoi(s_part_size)

			for (n_leidos != 0) && (pos_actual < int64(i_part_start+i_part_size)) {
				n_leidos, _ = f.Read(lectura)
				pos_actual, _ = f.Seek(0, io.SeekCurrent)

				// Conversion de bytes a struct
				extended_boot_record := Bytes_a_struct_ebr(lectura)

				s_part_name_ext := string(extended_boot_record.Part_name[:])
				s_part_name_ext = strings.Trim(s_part_name_ext, "\x00")

				ebr_empty_byte := Struct_a_bytes(ebr_empty)

				if s_part_name_ext == nombre {
					return int(pos_actual) - len(ebr_empty_byte)
				}
			}
		}
		f.Close()
	}

	return -1
}

func Formatear_ext2(inicio int, tamano int, direccion string, metodo string) string {
	sb_empty := beans.Super_bloque{}
	sb := Struct_a_bytes(sb_empty)

	in_empty := beans.Inodo{}
	in := Struct_a_bytes(in_empty)

	ba_empty := beans.Bloque_archivo{}
	ba := Struct_a_bytes(ba_empty)

	// Despejo n de la formula tamaño particion
	n := (tamano - len(sb)) / (4 + len(in) + 3*len(ba))
	num_estructuras := n
	num_bloques := 3 * num_estructuras

	Super_bloque := beans.Super_bloque{}

	fecha := time.Now()
	str_fecha := fecha.Format("02/01/2006 15:04:05")

	copy(Super_bloque.S_filesystem_type[:], "2")
	copy(Super_bloque.S_inodes_count[:], strconv.Itoa(num_estructuras))
	copy(Super_bloque.S_blocks_count[:], strconv.Itoa(num_bloques))
	copy(Super_bloque.S_free_blocks_count[:], strconv.Itoa(num_bloques-2))
	copy(Super_bloque.S_free_inodes_count[:], strconv.Itoa(num_estructuras-2))
	copy(Super_bloque.S_mtime[:], str_fecha)
	copy(Super_bloque.S_mnt_count[:], "0")
	copy(Super_bloque.S_magic[:], "0xEF53")
	copy(Super_bloque.S_inode_size[:], strconv.Itoa(len(in)))
	copy(Super_bloque.S_block_size[:], strconv.Itoa(len(ba)))
	copy(Super_bloque.S_firts_ino[:], "2")
	copy(Super_bloque.S_first_blo[:], "2")
	copy(Super_bloque.S_bm_inode_start[:], strconv.Itoa(inicio+len(sb)))
	copy(Super_bloque.S_bm_block_start[:], strconv.Itoa(inicio+len(sb)+num_estructuras))
	copy(Super_bloque.S_inode_start[:], strconv.Itoa(inicio+len(sb)+num_estructuras+num_bloques))
	copy(Super_bloque.S_block_start[:], strconv.Itoa(inicio+len(sb)+num_estructuras+num_bloques+len(in)+num_estructuras))

	inodo := beans.Inodo{}
	bloque := beans.Bloque_carpeta{}

	// Libre
	buffer := "0"
	// Usado o archivo
	buffer2 := "1"
	// Carpeta
	buffer3 := "2"

	f, err := os.OpenFile(direccion, os.O_RDWR, 0660)

	if err == nil {
		// Super Bloque
		f.Seek(int64(inicio), io.SeekStart)
		super_bloque_byte := Struct_a_bytes(Super_bloque)
		f.Write(super_bloque_byte)

		s_bm_inode_start := string(Super_bloque.S_bm_inode_start[:])
		s_bm_inode_start = strings.Trim(s_bm_inode_start, "\x00")
		i_bm_inode_start, _ := strconv.Atoi(s_bm_inode_start)

		for i := 0; i < num_estructuras; i++ {
			f.Seek(int64(i_bm_inode_start+i), io.SeekStart)
			f.Write([]byte(buffer))
		}

		f.Seek(int64(i_bm_inode_start), io.SeekStart)
		f.Write([]byte(buffer2))
		f.Write([]byte(buffer2))

		s_bm_block_start := string(Super_bloque.S_bm_block_start[:])
		s_bm_block_start = strings.Trim(s_bm_block_start, "\x00")
		i_bm_block_start, _ := strconv.Atoi(s_bm_block_start)

		for i := 0; i < num_bloques; i++ {
			f.Seek(int64(i_bm_block_start+i), io.SeekStart)
			f.Write([]byte(buffer))
		}

		// Marcando el Bitmap de Inodos para la carpeta "/" y el archivo "users.txt"
		f.Seek(int64(i_bm_block_start), io.SeekStart)
		f.Write([]byte(buffer2))
		f.Write([]byte(buffer3))

		/* Inodo para la carpeta "/" */
		copy(inodo.I_uid[:], "1")
		copy(inodo.I_gid[:], "1")
		copy(inodo.I_size[:], "0")
		copy(inodo.I_atime[:], str_fecha)
		copy(inodo.I_ctime[:], str_fecha)
		copy(inodo.I_mtime[:], str_fecha)
		copy(inodo.I_block[0:1], "0")

		// Nodos libres
		for i := 1; i < 15; i++ {
			copy(inodo.I_block[i:i+1], "$")
		}

		// 1 = archivo o 0 = Carpeta
		copy(inodo.I_type[:], "0")
		copy(inodo.I_perm[:], "664")

		s_inode_start := string(Super_bloque.S_inode_start[:])
		s_inode_start = strings.Trim(s_inode_start, "\x00")
		i_inode_start, _ := strconv.Atoi(s_inode_start)

		f.Seek(int64(i_inode_start), io.SeekStart)
		inodo_byte := Struct_a_bytes(inodo)
		f.Write(inodo_byte)

		/* Bloque Para Carpeta "/" */
		copy(bloque.B_content[0].B_name[:], ".")
		copy(bloque.B_content[0].B_inodo[:], "0")

		// Directorio Padre
		copy(bloque.B_content[1].B_name[:], "..")
		copy(bloque.B_content[1].B_inodo[:], "0")

		// Nombre de la carpeta o archivo
		copy(bloque.B_content[2].B_name[:], "users.txt")
		copy(bloque.B_content[2].B_inodo[:], "1")
		copy(bloque.B_content[3].B_name[:], ".")
		copy(bloque.B_content[3].B_inodo[:], "-1")

		s_block_start := string(Super_bloque.S_block_start[:])
		s_block_start = strings.Trim(s_block_start, "\x00")
		i_block_start, _ := strconv.Atoi(s_block_start)

		f.Seek(int64(i_block_start), io.SeekStart)
		bloque_byte := Struct_a_bytes(bloque)
		f.Write(bloque_byte)

		/* Inodo Para "users.txt" */
		copy(inodo.I_uid[:], "1")
		copy(inodo.I_gid[:], "1")
		copy(inodo.I_size[:], "29")
		copy(inodo.I_atime[:], str_fecha)
		copy(inodo.I_ctime[:], str_fecha)
		copy(inodo.I_mtime[:], str_fecha)
		copy(inodo.I_block[0:1], "1")

		// Nodos libres
		for i := 1; i < 15; i++ {
			copy(inodo.I_block[i:i+1], "$")
		}

		copy(inodo.I_type[:], "1")
		copy(inodo.I_perm[:], "755")

		f.Seek(int64(i_inode_start+len(in)), io.SeekStart)
		inodo_byte = Struct_a_bytes(inodo)
		f.Write(inodo_byte)

		/* Bloque Para "users.txt" */
		archivo := beans.Bloque_archivo{}

		for i := 0; i < 100; i++ {
			copy(archivo.B_content[i:i+1], "0")
		}

		bc_empty := beans.Bloque_carpeta{}
		bc := Struct_a_bytes(bc_empty)

		// GID, TIPO, GRUPO
		// UID, TIPO, GRUPO, USUARIO, CONTRASEÑA
		copy(archivo.B_content[:], "1,G,root\n1,U,root,root,123\n")
		f.Seek(int64(i_block_start+len(bc)), io.SeekStart)
		archivo_byte := Struct_a_bytes(archivo)
		f.Write(archivo_byte)
		f.Close()
		return "[MENSAJE-" + metodo + "] El Disco se formateo en el sistema EXT2 con exito!"
	}

	return "[ERROR-" + metodo + "] al abrir el disco!"
}

// Reporte DISK
func Graficar_disk(direccion string, destino string, extension string) string {
	graphDot = ""
	mbr_empty := beans.Mbr{}
	var empty [100]byte

	// Abro el archivo para lectura con opcion a modificar
	f, err := os.OpenFile(direccion, os.O_RDWR, 0660)

	// Calculo del tamaño de struct en bytes
	mbr2 := Struct_a_bytes(mbr_empty)
	sstruct := len(mbr2)

	// Lectrura del archivo binario desde el inicio
	lectura := make([]byte, sstruct)
	f.Seek(0, io.SeekStart)
	f.Read(lectura)

	// Conversion de bytes a struct
	master_boot_record := Bytes_a_struct_mbr(lectura)

	if master_boot_record.Mbr_tamano != empty {
		if err == nil {
			graphDot += "digraph G{\\n\\n"
			graphDot += "  tbl [\\n    shape=box\\n    label=<\\n"
			graphDot += "     <table border='0' cellborder='2' width='600' height=\\\"150\\\" color='dodgerblue1'>\\n"
			graphDot += "     <tr>\\n"
			graphDot += "     <td height='150' width='110'> MBR </td>\\n"

			// Obtengo el espacio utilizado
			s_mbr_tamano := string(master_boot_record.Mbr_tamano[:])
			// Le quito los caracteres null
			s_mbr_tamano = strings.Trim(s_mbr_tamano, "\x00")
			i_mbr_tamano, _ := strconv.Atoi(s_mbr_tamano)
			total := i_mbr_tamano

			var espacioUsado float64
			espacioUsado = 0

			for i := 0; i < 4; i++ {
				// Obtengo el espacio utilizado
				s_part_s := string(master_boot_record.Mbr_partition[i].Part_size[:])
				// Le quito los caracteres null
				s_part_s = strings.Trim(s_part_s, "\x00")
				i_part_s, _ := strconv.Atoi(s_part_s)

				parcial := i_part_s

				// Obtengo el espacio utilizado
				s_part_start := string(master_boot_record.Mbr_partition[i].Part_start[:])
				// Le quito los caracteres null
				s_part_start = strings.Trim(s_part_start, "\x00")

				if s_part_start != "-1" {
					var porcentaje_real float64
					porcentaje_real = (float64(parcial) * 100) / float64(total)
					var porcentaje_aux float64
					porcentaje_aux = (porcentaje_real * 500) / 100

					espacioUsado += porcentaje_real

					// Obtengo el espacio utilizado
					s_part_status := string(master_boot_record.Mbr_partition[i].Part_status[:])
					// Le quito los caracteres null
					s_part_status = strings.Trim(s_part_status, "\x00")

					if s_part_status != "1" {
						// Obtengo el espacio utilizado
						s_part_type := string(master_boot_record.Mbr_partition[i].Part_type[:])
						// Le quito los caracteres null
						s_part_type = strings.Trim(s_part_type, "\x00")

						if s_part_type == "p" {
							graphDot += "     <td height='200' width='" + strconv.FormatFloat(porcentaje_aux, 'g', 3, 64) + "'>Primaria <br/> " + strconv.FormatFloat(porcentaje_real, 'g', 3, 64) + " por ciento del Disco </td>\\n"

							if i != 3 {
								// Obtengo el espacio utilizado
								s_part_s = string(master_boot_record.Mbr_partition[i].Part_size[:])
								// Le quito los caracteres null
								s_part_s = strings.Trim(s_part_s, "\x00")
								i_part_s, _ = strconv.Atoi(s_part_s)

								// Obtengo el espacio utilizado
								s_part_start := string(master_boot_record.Mbr_partition[i].Part_start[:])
								// Le quito los caracteres null
								s_part_start = strings.Trim(s_part_start, "\x00")
								i_part_start, _ := strconv.Atoi(s_part_start)

								p1 := i_part_start + i_part_s

								// Obtengo el espacio utilizado
								s_part_start = string(master_boot_record.Mbr_partition[i+1].Part_start[:])
								// Le quito los caracteres null
								s_part_start = strings.Trim(s_part_start, "\x00")
								i_part_start, _ = strconv.Atoi(s_part_start)

								p2 := i_part_start

								if s_part_start != "-1" {
									if (p2 - p1) != 0 {
										fragmentacion := p2 - p1
										porcentaje_real = float64(fragmentacion) * 100 / float64(total)
										porcentaje_aux = (porcentaje_real * 500) / 100

										graphDot += "     <td height='200' width='" + strconv.FormatFloat(porcentaje_aux, 'g', 3, 64) + "'>Libre<br/> " + strconv.FormatFloat(porcentaje_real, 'g', 3, 64) + " por ciento del Disco </td>\\n"
									}
								}
							} else {
								// Obtengo el espacio utilizado
								s_part_s = string(master_boot_record.Mbr_partition[i].Part_size[:])
								// Le quito los caracteres null
								s_part_s = strings.Trim(s_part_s, "\x00")
								i_part_s, _ = strconv.Atoi(s_part_s)

								// Obtengo el espacio utilizado
								s_part_start := string(master_boot_record.Mbr_partition[i].Part_start[:])
								// Le quito los caracteres null
								s_part_start = strings.Trim(s_part_start, "\x00")
								i_part_start, _ := strconv.Atoi(s_part_start)

								p1 := i_part_start + i_part_s

								mbr_empty_byte := Struct_a_bytes(mbr_empty)
								mbr_size := total + len(mbr_empty_byte)

								// Si esta libre
								if (mbr_size - p1) != 0 {
									libre := (float64(mbr_size) - float64(p1)) + float64(len(mbr_empty_byte))
									porcentaje_real = (float64(libre) * 100) / float64(total)
									porcentaje_aux = (porcentaje_real * 500) / 100
									graphDot += "     <td height='200' width='" + strconv.FormatFloat(porcentaje_aux, 'g', 3, 64) + "'>Libre<br/> " + strconv.FormatFloat(porcentaje_real, 'g', 3, 64) + " por ciento del Disco </td>\\n"
								}
							}
						} else {
							// Si es extendida
							graphDot += "     <td  height='200' width='" + strconv.FormatFloat(porcentaje_real, 'g', 3, 64) + "'>\\n     <table border='0'  height='200' WIDTH='" + strconv.FormatFloat(porcentaje_real, 'g', 3, 64) + "' cellborder='1'>\\n"
							graphDot += "     <tr>  <td height='60' colspan='15'>Extendida</td>  </tr>\\n     <tr>\\n"

							// Obtengo el espacio utilizado
							s_part_start := string(master_boot_record.Mbr_partition[i].Part_start[:])
							// Le quito los caracteres null
							s_part_start = strings.Trim(s_part_start, "\x00")
							i_part_start, _ := strconv.Atoi(s_part_start)

							f.Seek(int64(i_part_start), io.SeekStart)

							ebr_empty := beans.Ebr{}

							// Calculo del tamaño de struct en bytes
							ebr2 := Struct_a_bytes(ebr_empty)
							sstruct := len(ebr2)

							// Lectrura del archivo binario desde el inicio
							lectura := make([]byte, sstruct)
							f.Read(lectura)

							// Conversion de bytes a struct
							extended_boot_record := Bytes_a_struct_ebr(lectura)

							// Obtengo el espacio utilizado
							s_part_size := string(extended_boot_record.Part_size[:])
							// Le quito los caracteres null
							s_part_size = strings.Trim(s_part_size, "\x00")
							i_part_size, _ := strconv.Atoi(s_part_size)

							if i_part_size != 0 {
								// Obtengo el espacio utilizado
								s_part_start := string(master_boot_record.Mbr_partition[i].Part_start[:])
								// Le quito los caracteres null
								s_part_start = strings.Trim(s_part_start, "\x00")
								i_part_start, _ := strconv.Atoi(s_part_start)

								f.Seek(int64(i_part_start), io.SeekStart)

								band := true

								// Obtengo el espacio utilizado
								s_part_s = string(master_boot_record.Mbr_partition[i].Part_size[:])
								// Le quito los caracteres null
								s_part_s = strings.Trim(s_part_s, "\x00")
								i_part_s, _ = strconv.Atoi(s_part_s)

								// Obtengo el espacio utilizado
								s_part_start = string(master_boot_record.Mbr_partition[i].Part_start[:])
								// Le quito los caracteres null
								s_part_start = strings.Trim(s_part_start, "\x00")
								i_part_start, _ = strconv.Atoi(s_part_start)

								for band {
									// Calculo del tamaño de struct en bytes
									ebr2 := Struct_a_bytes(ebr_empty)
									sstruct := len(ebr2)

									// Lectrura del archivo binario desde el inicio
									lectura := make([]byte, sstruct)
									f.Seek(0, io.SeekStart)
									n, _ := f.Read(lectura)

									// Posicion actual en el archivo
									pos_actual, _ := f.Seek(0, io.SeekCurrent)

									if n != 0 && pos_actual < int64(i_part_start)+int64(i_part_s) {
										band = false
										break
									}

									// Obtengo el espacio utilizado
									s_part_s = string(extended_boot_record.Part_size[:])
									// Le quito los caracteres null
									s_part_s = strings.Trim(s_part_s, "\x00")
									i_part_s, _ = strconv.Atoi(s_part_s)

									parcial = i_part_start
									porcentaje_real = float64(parcial) * 100 / float64(total)

									if porcentaje_real != 0 {
										// Obtengo el espacio utilizado
										s_part_status = string(extended_boot_record.Part_status[:])
										// Le quito los caracteres null
										s_part_status = strings.Trim(s_part_status, "\x00")

										if s_part_status != "1" {
											graphDot += "     <td height='140'>EBR</td>\\n"
											graphDot += "     <td height='140'>Logica<br/> " + strconv.FormatFloat(porcentaje_real, 'g', 3, 64) + " por ciento del Disco </td>\\n"
										} else {
											// Espacio no asignado
											graphDot += "      <td height='150'>Libre 1 <br/> " + strconv.FormatFloat(porcentaje_real, 'g', 3, 64) + " por ciento del Disco </td>\\n"
										}

										// Obtengo el espacio utilizado
										s_part_next := string(extended_boot_record.Part_next[:])
										// Le quito los caracteres null
										s_part_next = strings.Trim(s_part_next, "\x00")
										i_part_next, _ := strconv.Atoi(s_part_next)

										if i_part_next == -1 {
											// Obtengo el espacio utilizado
											s_part_start := string(extended_boot_record.Part_start[:])
											// Le quito los caracteres null
											s_part_start = strings.Trim(s_part_start, "\x00")
											i_part_start, _ := strconv.Atoi(s_part_start)

											// Obtengo el espacio utilizado
											s_part_size := string(extended_boot_record.Part_size[:])
											// Le quito los caracteres null
											s_part_size = strings.Trim(s_part_size, "\x00")
											i_part_size, _ := strconv.Atoi(s_part_size)

											// Obtengo el espacio utilizado
											s_part_start_mbr := string(master_boot_record.Mbr_partition[i].Part_start[:])
											// Le quito los caracteres null
											s_part_start_mbr = strings.Trim(s_part_start_mbr, "\x00")
											i_part_start_mbr, _ := strconv.Atoi(s_part_start_mbr)

											// Obtengo el espacio utilizado
											s_part_s_mbr := string(master_boot_record.Mbr_partition[i].Part_size[:])
											// Le quito los caracteres null
											s_part_s_mbr = strings.Trim(s_part_s_mbr, "\x00")
											i_part_s_mbr, _ := strconv.Atoi(s_part_s_mbr)

											parcial = (i_part_start_mbr + i_part_s_mbr) - (i_part_size + i_part_start)
											porcentaje_real = (float64(parcial) * 100) / float64(total)

											if porcentaje_real != 0 {
												graphDot += "     <td height='150'>Libre 2<br/> " + strconv.FormatFloat(porcentaje_real, 'g', 3, 64) + " por ciento del Disco </td>\\n"
											}
											break

										} else {
											// Obtengo el espacio utilizado
											s_part_next := string(extended_boot_record.Part_next[:])
											// Le quito los caracteres null
											s_part_next = strings.Trim(s_part_next, "\x00")
											i_part_next, _ := strconv.Atoi(s_part_next)

											f.Seek(int64(i_part_next), io.SeekStart)
										}
									}

								}
							} else {
								graphDot += "     <td height='140'> " + strconv.FormatFloat(porcentaje_real, 'g', 3, 64) + " por ciento del Disco </td>\\n"
							}
							graphDot += "     </tr>\\n     </table>\\n     </td>\\n"

							// Verifica que no haya espacio fragemntado
							if i != 3 {
								// Obtengo el espacio utilizado
								s_part_s = string(master_boot_record.Mbr_partition[i].Part_size[:])
								// Le quito los caracteres null
								s_part_s = strings.Trim(s_part_s, "\x00")
								i_part_s, _ = strconv.Atoi(s_part_s)

								// Obtengo el espacio utilizado
								s_part_start := string(master_boot_record.Mbr_partition[i].Part_start[:])
								// Le quito los caracteres null
								s_part_start = strings.Trim(s_part_start, "\x00")
								i_part_start, _ := strconv.Atoi(s_part_start)

								p1 := i_part_start + i_part_s

								// Obtengo el espacio utilizado
								s_part_start = string(master_boot_record.Mbr_partition[i+1].Part_start[:])
								// Le quito los caracteres null
								s_part_start = strings.Trim(s_part_start, "\x00")
								i_part_start, _ = strconv.Atoi(s_part_start)

								p2 := i_part_start

								if s_part_start != "-1" {
									if (p2 - p1) != 0 {
										fragmentacion := p2 - p1
										porcentaje_real = float64(fragmentacion) * 100 / float64(total)
										porcentaje_aux = (porcentaje_real * 500) / 100

										graphDot += "     <td height='200' width='" + strconv.FormatFloat(porcentaje_aux, 'g', 3, 64) + "'>Libre<br/> " + strconv.FormatFloat(porcentaje_real, 'g', 3, 64) + " por ciento del Disco </td>\\n"
									}
								}
							} else {
								// Obtengo el espacio utilizado
								s_part_s = string(master_boot_record.Mbr_partition[i].Part_size[:])
								// Le quito los caracteres null
								s_part_s = strings.Trim(s_part_s, "\x00")
								i_part_s, _ = strconv.Atoi(s_part_s)

								// Obtengo el espacio utilizado
								s_part_start := string(master_boot_record.Mbr_partition[i].Part_start[:])
								// Le quito los caracteres null
								s_part_start = strings.Trim(s_part_start, "\x00")
								i_part_start, _ := strconv.Atoi(s_part_start)

								p1 := i_part_start + i_part_s

								mbr_empty_byte := Struct_a_bytes(mbr_empty)
								mbr_size := total + len(mbr_empty_byte)

								// Si esta libre
								if (mbr_size - p1) != 0 {
									libre := (float64(mbr_size) - float64(p1)) + float64(len(mbr_empty_byte))
									porcentaje_real = (float64(libre) * 100) / float64(total)
									porcentaje_aux = porcentaje_real * 500 / 100
									graphDot += "     <td height='200' width='" + strconv.FormatFloat(porcentaje_aux, 'g', 3, 64) + "'>Libre<br/> " + strconv.FormatFloat(porcentaje_real, 'g', 3, 64) + " por ciento del Disco </td>\\n"
								}
							}
						}
					} else {
						graphDot += "     <td height='200' width='" + strconv.FormatFloat(porcentaje_aux, 'g', 3, 64) + "'>Libre<br/> " + strconv.FormatFloat(porcentaje_real, 'g', 3, 64) + " por ciento del Disco </td>\\n"
					}
				}
			}

			graphDot += "     </tr> \\n     </table>        \\n>];\\n\\n}"
			return graphDot
		} else {
			return "[ERROR-REP] El disco no fue encontrado."
		}
	} else {
		return "[ERROR-REP] Disco vacio."
	}
}

/*--------------------------/Metodos o Funciones--------------------------*/
