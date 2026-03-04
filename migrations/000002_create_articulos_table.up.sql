CREATE TABLE IF NOT EXISTS articulos (
    id VARCHAR(36) PRIMARY KEY,
    titulo VARCHAR(200) NOT NULL,
    contenido TEXT NOT NULL,
    autor_id VARCHAR(36) NOT NULL,
    estado ENUM('borrador', 'publicado', 'archivado') DEFAULT 'borrador',
    fecha_publicacion TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    -- Índices para búsquedas rápidas
    INDEX idx_articulos_autor (autor_id),
    INDEX idx_articulos_estado (estado),
    INDEX idx_articulos_fecha_publicacion (fecha_publicacion DESC),
    
    -- Relación con autores
    CONSTRAINT fk_articulos_autor FOREIGN KEY (autor_id) REFERENCES autores(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;