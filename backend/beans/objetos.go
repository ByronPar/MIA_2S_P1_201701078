package beans

// Definición de la estructura de datos
/*type Profesor struct {
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
}*/

/*--------------------------/Import--------------------------*/

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
	Part_mount [100]byte
	Part_fit   [100]byte
	Part_start [100]byte
	Part_size  [100]byte
	Part_next  [100]byte
	Part_name  [100]byte
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
