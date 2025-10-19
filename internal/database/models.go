package database

import (
	"time"

	"gorm.io/datatypes"
)

type Usuario struct {
	ID            uint      `gorm:"primaryKey"`
	Nombre        string    `gorm:"size:100;not null"`
	Email         string    `gorm:"size:150;unique;not null"`
	PasswordHash  string    `gorm:"size:255;not null"`
	FechaCreacion time.Time `gorm:"autoCreateTime"`
	Fotos         []Foto    `gorm:"foreignKey:UsuarioID;constraint:OnDelete:CASCADE"`
	Token					string 		`gorm:"size:255;not null"`
}

type Foto struct {
	ID          uint    `gorm:"primaryKey"`
	UsuarioID   uint     `gorm:"not null"`
	Usuario     Usuario  `gorm:"constraint:OnDelete:CASCADE;"`
	Nombre      string	`gorm:"size:255;not null;uniqueIndex:idx_usuario_nombre"`
	Descripcion string
	URLArchivo  string    `gorm:"not null"`
	FechaSubida time.Time `gorm:"autoCreateTime"`
	Etiquetas   datatypes.JSON
	Favorito    bool `gorm:"default:false"`
	TamanoBytes int64
	Formato     string
	Albums      []Album   `gorm:"many2many:album_fotos"`
}

type Album struct {
	ID            uint   `gorm:"primaryKey"`
	UsuarioID     uint   `gorm:"not null"`
	Nombre        string `gorm:"size:255;not null"`
	Descripcion   string
	Tipo          string    `gorm:"default:'normal'"`
	FechaCreacion time.Time `gorm:"autoCreateTime"`
	Fotos         []Foto    `gorm:"many2many:album_fotos;constraint:OnDelete:CASCADE;"`
}
