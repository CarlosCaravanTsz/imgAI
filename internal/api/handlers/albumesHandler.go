package handlers

import (
	"net/http"
	_ "strconv"

	"github.com/CarlosCaravanTsz/imgAI/internal/database"
	log "github.com/CarlosCaravanTsz/imgAI/internal/logger"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type AlbumParams struct {
	Nombre string `form:"nombre" json:"nombre"`
	Email  string `form:"email" json:"email"`
}

type AlbumArray struct {
	Nombre      string
	Descripcion string
}

type AlbumesRouteHandlers struct{}

func (h *AlbumesRouteHandlers) CrearAlbum(c *gin.Context) {
	var params AlbumParams

	if err := c.ShouldBind(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db, err := database.GetConnection()
	if err != nil {
		log.LogInfo("Error getting the db conn", logrus.Fields{})
	}

	type userInfo struct {
		ID     uint
		Nombre string
	}

	var usuario userInfo
	// Buscar el user con el email y obtener ID
	err = db.Model(&database.Usuario{}).Select("id", "nombre").Where("email = ?", params.Email).First(&usuario).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "User not found"})
		return
	}

	dbAlbum := database.Album{
		UsuarioID:   usuario.ID,
		Nombre:      params.Nombre,
		Descripcion: "",
	}

	results := db.Create(&dbAlbum)
	if results.Error != nil {
		c.JSON(500, gin.H{"error": results.Error})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Ok, Created Album"})
}

func (h *AlbumesRouteHandlers) ListarAlbumes(c *gin.Context) {
	email := c.Query("email")

	var albums []database.Album

	db, err := database.GetConnection()
	if err != nil {
		log.LogInfo("Error getting the db conn", logrus.Fields{})
	}

	type userInfo struct {
		ID     uint
		Nombre string
	}

	var usuario userInfo
	// Buscar el user con el email y obtener ID
	err = db.Model(&database.Usuario{}).Select("id", "nombre").Where("email = ?", email).First(&usuario).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "User not found"})
		return
	}

	results := db.Joins("JOIN usuarios ON usuarios.id = albums.usuario_id").
		Where("usuarios.email = ?", email).
		Find(&albums)
	if results.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	urls := make([]AlbumArray, len(albums))

	for i, a := range albums {
		urls[i] = AlbumArray{Nombre: a.Nombre, Descripcion: a.Descripcion}
	}

	c.JSON(http.StatusOK, urls)
}

func (h *AlbumesRouteHandlers) ListarFotosAlbum(c *gin.Context) {
	email := c.Query("email")
	albumid := c.Param("albumid")

	db, err := database.GetConnection()
	if err != nil {
		log.LogInfo("Error getting the db conn", logrus.Fields{})
	}

	type userInfo struct {
		ID     uint
		Nombre string
	}

	var usuario userInfo
	// Buscar el user con el email y obtener ID
	err = db.Model(&database.Usuario{}).Select("id", "nombre").Where("email = ?", email).First(&usuario).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "User not found"})
		return
	}

	var album database.Album
	results := db.Preload("Fotos").Where("id = ? AND usuario_id = ?", albumid, usuario.ID).First(&album)

	if results.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Álbum no encontrado o no pertenece al usuario"})
		return
	}

	c.JSON(http.StatusOK, album.Fotos)
}

func (h *AlbumesRouteHandlers) EliminarAlbum(c *gin.Context) {
	email := c.Query("email")
	albumID := c.Param("albumid")
	type userInfo struct {
		ID     uint
		Nombre string
	}

	db, err := database.GetConnection()
	if err != nil {
		log.LogInfo("Error getting the db conn", logrus.Fields{})
	}

	var usuario userInfo
	// Buscar el user con el email y obtener ID
	err = db.Model(&database.Usuario{}).Select("id", "nombre").Where("email = ?", email).First(&usuario).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "User not found"})
		return
	}

	var album database.Album
	if err := db.Where("id = ? AND usuario_id = ?", albumID, usuario.ID).First(&album).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Álbum no encontrado o no pertenece al usuario"})
		return
	}

	if err := db.Model(&album).Association("Fotos").Clear(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudieron eliminar las relaciones del álbum"})
		return
	}

	if err := db.Delete(&album).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar el álbum"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Álbum eliminado correctamente"})
}

func (h *AlbumesRouteHandlers) QuitarFotoDeAlbum(c *gin.Context) {
	email := c.Query("email")
	albumID := c.Param("albumid")
	fotoID := c.Param("fotoid")

	db, err := database.GetConnection()
	if err != nil {
		log.LogInfo("Error getting the db conn", logrus.Fields{})
	}

	type userInfo struct {
		ID     uint
		Nombre string
	}

	var usuario userInfo
	// Buscar el user con el email y obtener ID
	err = db.Model(&database.Usuario{}).Select("id", "nombre").Where("email = ?", email).First(&usuario).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "User not found"})
		return
	}

	var album database.Album
	if err := db.Where("id = ? AND usuario_id = ?", albumID, usuario.ID).First(&album).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Álbum no encontrado o no pertenece al usuario"})
		return
	}

	var foto database.Foto
	if err := db.First(&foto, fotoID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Foto no encontrada"})
		return
	}

	// 4️⃣ Remove association
	if err := db.Model(&album).Association("Fotos").Delete(&foto); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo eliminar la foto del álbum"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Foto eliminada del álbum correctamente"})
}
