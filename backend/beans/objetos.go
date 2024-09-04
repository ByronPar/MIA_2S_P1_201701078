package beans

/*--------------------------Structurass--------------------------*/

// Master Boot Record
type Mbr = struct {
	Mbr_tamano         [100]byte
	Mbr_fecha_creacion [100]byte
	Mbr_dsk_signature  [100]byte
	Dsk_fit            [100]byte
	Mbr_partition      [4]Partition
}

// Partition
type Partition = struct {
	Part_status      [100]byte
	Part_type        [100]byte
	Part_fit         [100]byte
	Part_start       [100]byte
	Part_size        [100]byte
	Part_name        [100]byte
	Part_correlative [100]byte
	Part_id          [100]byte
}

// Extended Boot Record
type Ebr = struct {
	Part_status [100]byte
	Part_fit    [100]byte
	Part_start  [100]byte
	Part_size   [100]byte
	Part_next   [100]byte
	Part_name   [100]byte
}

// Super Bloque
type Super_bloque = struct {
	S_filesystem_type   [100]byte
	S_inodes_count      [100]byte
	S_blocks_count      [100]byte
	S_free_blocks_count [100]byte
	S_free_inodes_count [100]byte
	S_mtime             [100]byte
	S_umtime            [100]byte
	S_mnt_count         [100]byte
	S_magic             [100]byte
	S_inode_size        [100]byte
	S_block_size        [100]byte
	S_firts_ino         [100]byte
	S_first_blo         [100]byte
	S_bm_inode_start    [100]byte
	S_bm_block_start    [100]byte
	S_inode_start       [100]byte
	S_block_start       [100]byte
}

// Tablas de Inodos
type Inodo = struct {
	I_uid   [100]byte
	I_gid   [100]byte
	I_size  [100]byte
	I_atime [100]byte
	I_ctime [100]byte
	I_mtime [100]byte
	I_block [100]byte
	I_type  [100]byte
	I_perm  [100]byte
}

// Bloques de Carpetas
type Bloque_carpeta = struct {
	B_content [4]Cotent
}

// Content
type Cotent = struct {
	B_name  [100]byte
	B_inodo [100]byte
}

// Bloques de Archivos
type Bloque_archivo = struct {
	B_content [100]byte
}

// Bloques de apuntadores
type Bloque_apuntadores = struct {
	B_pointers [100]byte
}

/*--------------------------/Structs--------------------------*/

// elementos para comando mount
type Nodo struct {
	id        string
	Direccion string
	Nombre    string
	Letra     string
	Num       int
	Siguiente *Nodo
}

type Lista struct {
	Primero *Nodo
	Ultimo  *Nodo
}

// Nuevo nodo con datos de la particion
func New_nodo(id string, direccion string, nombre string, letra string, num int) *Nodo {
	return &Nodo{id, direccion, nombre, letra, num, nil}
}

// Nueva lista para los nodos
func New_lista() *Lista {
	return &Lista{nil, nil}
}

// Inserta particion
func Insertar(nuevo *Nodo, lista *Lista) {
	// Si esta vacia
	if lista.Primero == nil {
		lista.Primero = nuevo
		lista.Ultimo = nuevo
	} else {
		lista.Ultimo.Siguiente = nuevo
		lista.Ultimo = lista.Ultimo.Siguiente
	}
}

// Elimina particion por nombre
func Eliminar(id string, lista *Lista) int {
	aux := lista.Primero

	// Si esta vacia
	if aux != nil {
		if aux.id == id {
			lista.Primero = aux.Siguiente
			return 1
		} else {
			var aux2 *Nodo = nil

			for aux != nil {
				if aux.id == id {
					aux2.Siguiente = aux.Siguiente
					return 1
				}
				aux2 = aux
				aux = aux.Siguiente
			}
		}
	}

	return 0
}

// Busca la particion por direccion y nombre
func Buscar_particion(direccion string, nombre string, lista *Lista) bool {
	aux := lista.Primero

	for aux != nil {
		if (direccion == aux.Direccion) && (nombre == aux.Nombre) {
			return true
		}
		aux = aux.Siguiente
	}

	return false
}

// Busca por direccion y obtiene el numero
func Buscar_letra(direccion string, lista *Lista) string {
	aux := lista.Primero
	num_ascii := 65

	for aux != nil {
		if (direccion == aux.Direccion) && (string(rune(num_ascii)) == aux.Letra) {
			num_ascii++
		}
		aux = aux.Siguiente
	}

	return string(rune(num_ascii))
}

// Busca por direccion y obtiene el numero
func Buscar_numero(direccion string, lista *Lista) int {
	aux := lista.Primero
	retorno := 1

	for aux != nil {
		if direccion == aux.Direccion {
			return retorno
		}

		aux = aux.Siguiente
		retorno++
	}

	return retorno
}

// Busca por direccion y nombre retorna un boolean
func Buscar_nodo(direccion string, nombre string, lista *Lista) bool {
	aux := lista.Primero

	for aux != nil {
		if (direccion == aux.Direccion) && (aux.Nombre == nombre) {
			return true
		}
		aux = aux.Siguiente
	}

	return false
}

// Busca por id y retorna un nodo
func Obtener_nodo(id string, lista *Lista) *Nodo {
	aux := lista.Primero

	for aux != nil {
		//201602530
		if id == aux.id {
			return aux
		}
		aux = aux.Siguiente
	}

	return nil
}
