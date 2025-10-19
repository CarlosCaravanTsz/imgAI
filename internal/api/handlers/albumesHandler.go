package handlers

import (
	"net/http"
	_ "strconv"

	"github.com/CarlosCaravanTsz/imgAI/internal/database"
	log "github.com/CarlosCaravanTsz/imgAI/internal/logger"
	_ "github.com/aws/aws-sdk-go-v2/aws/retry"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type AlbumArray struct {
		Nombre      string `form:"name" json:"name" `
		Descripcion string `form:"description" json:"description" `
}

type AlbumesRouteHandlers struct{}

func (h *AlbumesRouteHandlers) CrearAlbum(c *gin.Context) {
	var body AlbumArray

	if err := c.ShouldBind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	usuario, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error: User not set in req"})
		return
	}

	db := database.GetConnection()

	dbAlbum := database.Album{
		UsuarioID:   usuario.(database.Usuario).ID,
		Nombre:      body.Nombre,
		Descripcion: body.Descripcion,
	}

	results := db.Create(&dbAlbum)
	if results.Error != nil {
		c.JSON(500, gin.H{"error": results.Error})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Ok, Created Album"})
}

func (h *AlbumesRouteHandlers) ListarAlbumes(c *gin.Context) {
	var albums []database.Album

	usuario, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error: User not set in req"})
		return
	}

	db := database.GetConnection()


	results := db.Joins("JOIN usuarios ON usuarios.id = albums.usuario_id").
		Where("usuarios.email = ?", usuario.(database.Usuario).Email).
		Find(&albums)
	if results.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting the association"})
		return
	}


	type AlbumStruct struct {
		ID uint
		Nombre string
		Descripcion string
	}

	urls := make([]AlbumStruct, len(albums))

	for i, a := range albums {
		urls[i] = AlbumStruct{ID: a.ID, Nombre: a.Nombre, Descripcion: a.Descripcion}
	}

	c.JSON(http.StatusOK, urls)
}

func (h *AlbumesRouteHandlers) ListarFotosAlbum(c *gin.Context) {
	albumid := c.Param("albumid")

	usuario, exists := c.Get("user")
	if !exists {
		log.LogError("Error loading multipart form", logrus.Fields{})
		c.JSON(http.StatusBadRequest, gin.H{"status": "error: User not set in req"})
		return
	}

	db := database.GetConnection()

	var album database.Album

	results := db.Preload("Fotos").Where("id = ? AND usuario_id = ?", albumid, usuario.(database.Usuario).ID).First(&album)

	if results.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Álbum no encontrado o no pertenece al usuario"})
		return
	}

	c.JSON(http.StatusOK, album.Fotos)
}

func (h *AlbumesRouteHandlers) EliminarAlbum(c *gin.Context) {
	albumID := c.Param("albumid")

	usuario, exists := c.Get("user")
	if !exists {
		log.LogError("Error loading multipart form", logrus.Fields{})
		c.JSON(http.StatusBadRequest, gin.H{"status": "error: User not set in req"})
		return
	}

	db := database.GetConnection()

	var album database.Album
	if err := db.Where("id = ? AND usuario_id = ?", albumID, usuario.(database.Usuario).ID).First(&album).Error; err != nil {
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
	albumID := c.Param("albumid")
	fotoID := c.Param("fotoid")

	usuario, exists := c.Get("user")
	if !exists {
		log.LogError("Error loading multipart form", logrus.Fields{})
		c.JSON(http.StatusBadRequest, gin.H{"status": "error: User not set in req"})
		return
	}

	db := database.GetConnection()

	var album database.Album
	if err := db.Where("id = ? AND usuario_id = ?", albumID, usuario.(database.Usuario).ID).First(&album).Error; err != nil {
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
