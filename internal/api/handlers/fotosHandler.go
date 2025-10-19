package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"sync"
"gorm.io/datatypes"

	ai "github.com/CarlosCaravanTsz/imgAI/internal/ai"
	"github.com/CarlosCaravanTsz/imgAI/internal/database"
	log "github.com/CarlosCaravanTsz/imgAI/internal/logger"
	s3 "github.com/CarlosCaravanTsz/imgAI/internal/storage"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type FotosArray struct {
	URL string `json:"url"`
}

type FotosRouteHandlers struct{}

func (h *FotosRouteHandlers) SubirFotos(c *gin.Context) {
	usuario, exists := c.Get("user")
	if !exists {
		log.LogError("Error loading multipart form", logrus.Fields{})
		c.JSON(http.StatusBadRequest, gin.H{"status": "error: User not set in req"})
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		log.LogError("Error loading multipart form", logrus.Fields{"error": err})
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	files := form.File["images[]"]

	// Inicia Cliente s3:
	s3Client, err := s3.NewS3Client()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize S3 client"})
		return
	}

	// Inicializa conn a DB
	db := database.GetConnection()

	// Creamos channel con tipo
	type uploadResult struct {
		URL      string
		Filename string
		Error    error
	}
	resultsChan := make(chan uploadResult, len(files))
	var wg sync.WaitGroup

	for _, fileHeader := range files {

		wg.Add(1)
		go func(fh *multipart.FileHeader) {
			defer wg.Done()

			file, err := fh.Open()
			if err != nil {
				resultsChan <- uploadResult{Error: fmt.Errorf("Error reading %s: %v", fh.Filename, err)}
				return
			}
			defer file.Close()

			fileBytes, err := io.ReadAll(file)
			if err != nil {
				resultsChan <- uploadResult{Error: fmt.Errorf("Error reading %s: %v", fh.Filename, err)}
				return
			}

			// Guardar la Foto en S3 en el bucket y carpeta /{user}/filename

			foto := s3.FotoUpload{
				Filename: fileHeader.Filename,
				Path:     usuario.(database.Usuario).Email,
				Buffer:   fileBytes,
			}

			url, err := s3Client.Upload(foto)
			if err != nil {
				resultsChan <- uploadResult{Error: fmt.Errorf("upload failed for %s: %v", fh.Filename, err)}
				return
			}

			log.LogInfo("Uploaded to S3", logrus.Fields{"url": url})

			analysis, err := ai.ObtainDescription(url)
			if err != nil {
				resultsChan <- uploadResult{Error: fmt.Errorf("AI analysis failed: %v", err)}
				return
			}
			fmt.Print("TAGS RETURNED BEFORE CASTING TO JSON chatAI: ", analysis.Tags)

			tagsJSON, err := json.Marshal(analysis.Tags)
			if err != nil {
				resultsChan <- uploadResult{Error: fmt.Errorf("Error reading tags")}
			}

			//  Obtener la URL de S3 y guardar la informacion de la Foto en la BD
			dbFoto := database.Foto{
				UsuarioID:   usuario.(database.Usuario).ID,
				Nombre:      fileHeader.Filename,
				Descripcion: analysis.Description,
				URLArchivo:  url,
				TamanoBytes: int64(len(fileBytes)),
				Etiquetas:   datatypes.JSON(tagsJSON),
				Formato:     "image",
			}

			if err := db.Create(&dbFoto).Error; err != nil {
				resultsChan <- uploadResult{Error: fmt.Errorf("db inserted failed for %s: %v", fh.Filename, err)}
				return
			}

			resultsChan <- uploadResult{URL: url, Filename: fh.Filename, Error: nil}
		}(fileHeader)
	}

	wg.Wait()
	close(resultsChan)

	var uploaded []string
	var failed []string

	for r := range resultsChan {
		if r.Error != nil {
			log.LogError("Upload error ocurred", logrus.Fields{"file": r.Filename, "error": r.Error})
			fmt.Print("ERROR IN AI: ",r.Error)
			failed = append(failed, r.Filename)
		} else {
			uploaded = append(uploaded, r.URL)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"uploaded": uploaded,
		"failed":   failed,
	})
}

func (h *FotosRouteHandlers) ListarFotos(c *gin.Context) {
	usuario, exists := c.Get("user")
	if !exists {
		log.LogError("Error loading multipart form", logrus.Fields{})
		c.JSON(http.StatusBadRequest, gin.H{"status": "error: User not set in req"})
		return
	}

	var fotos []database.Foto

	db := database.GetConnection()

	// Query only photos for this user
	err := db.Joins("JOIN usuarios ON usuarios.id = fotos.usuario_id").
		Where("usuarios.email = ?", usuario.(database.Usuario).Email).
		Find(&fotos).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Build array of URLs
	urls := make([]FotosArray, len(fotos))

	for i, f := range fotos {
		urls[i] = FotosArray{URL: f.URLArchivo}
	}

	c.JSON(http.StatusOK, urls)
}

func (h *FotosRouteHandlers) ListarUnaFoto(c *gin.Context) {
	fotoid := c.Query("id")

	usuario, exists := c.Get("user")
	if !exists {
		log.LogError("Error loading multipart form", logrus.Fields{})
		c.JSON(http.StatusBadRequest, gin.H{"status": "error: User not set in req"})
		return
	}

	var foto database.Foto

	db := database.GetConnection()

	err := db.Where("usuario_id = ? AND id = ?", usuario.(database.Usuario).ID, fotoid).First(&foto).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"status": "error: Foto does not belong to the user"})
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error: Foto not found"})
			return
		}
	}

	c.JSON(http.StatusOK, foto)
}

func (h *FotosRouteHandlers) EliminarFoto(c *gin.Context) {
	fotoid := c.Param("id")

	usuario, exists := c.Get("user")
	if !exists {
		log.LogError("Error loading multipart form", logrus.Fields{})
		c.JSON(http.StatusBadRequest, gin.H{"status": "error: User not set in req"})
		return
	}

	db := database.GetConnection()

	result := db.Where("usuario_id = ? AND id = ?", usuario.(database.Usuario).ID, fotoid).Delete(&database.Foto{})
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Foto not found with that userid"})
		return

	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"status": "No record found"})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"status": "Ok, Deleted"})
}

func (h *FotosRouteHandlers) ToggleFavorito(c *gin.Context) {
	fotoid := c.Param("id")

	usuario, exists := c.Get("user")
	if !exists {
		log.LogError("Error loading multipart form", logrus.Fields{})
		c.JSON(http.StatusBadRequest, gin.H{"status": "error: User not set in req"})
		return
	}

	db := database.GetConnection()

result := db.Model(&database.Foto{}).
    Where("usuario_id = ? AND id = ?", usuario.(database.Usuario).ID, fotoid).
    Update("favorito", gorm.Expr("NOT favorito"))

if result.Error != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"status": "Database error", "error": result.Error.Error()})
    return
}

if result.RowsAffected == 0 {
    c.JSON(http.StatusNotFound, gin.H{"status": "Foto not found"})
    return
}

	c.JSON(http.StatusAccepted, gin.H{"status": "Favorito toggled successfully"})

}

func (h *FotosRouteHandlers) ListarFavoritos(c *gin.Context) {
	var fotos []database.Foto

	usuario, exists := c.Get("user")
	if !exists {
		log.LogError("Error loading multipart form", logrus.Fields{})
		c.JSON(http.StatusBadRequest, gin.H{"status": "error: User not set in req"})
	}

	db := database.GetConnection()

	// Query only photos for this user
	err := db.Joins("JOIN usuarios ON usuarios.id = fotos.usuario_id").
		Where("usuarios.email = ? AND fotos.favorito=1", usuario.(database.Usuario).Email).
		Find(&fotos).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Build array of URLs
	urls := make([]FotosArray, len(fotos))

	for i, f := range fotos {
		urls[i] = FotosArray{URL: f.URLArchivo}
	}

	c.JSON(http.StatusOK, urls)
}

func (h *FotosRouteHandlers) AgregarFotoaAlbum(c *gin.Context) {
	fotoid := c.Param("fotoid")
	albumid := c.Param("albumid")

	fotoID, err := strconv.ParseUint(fotoid, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid foto ID"})
		return
	}

	albumID, err := strconv.ParseUint(albumid, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid album ID"})
		return
	}

	usuario, ok := c.Get("user")
	if !ok {
		log.LogError("Error loading multipart form", logrus.Fields{})
		c.JSON(http.StatusBadRequest, gin.H{"status": "error: User not set in req"})
		return
	}

	db := database.GetConnection()

	// Find album and foto
	var album database.Album
	if err := db.Where("id = ? AND usuario_id = ?", albumID, usuario.(database.Usuario).ID).First(&album).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Album not found or does not belong to user"})
		return
	}

	var foto database.Foto
	if err := db.Where("id = ? AND usuario_id = ?", fotoID, usuario.(database.Usuario).ID).First(&foto).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Foto not found or does not belong to user"})
		return
	}

// count := db.Model(&album).Where("id = ?", foto.ID).Association("Fotos").Count()
// if count > 0 {
//     c.JSON(http.StatusOK, gin.H{"status": "Foto already in album"})
//     return
// }

var exists bool
err = db.Raw(`
    SELECT EXISTS(
        SELECT 1 FROM album_fotos WHERE album_id = ? AND foto_id = ?
    )
`, albumID, fotoID).Scan(&exists).Error

if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Database check failed"})
    return
}

if exists {
    c.JSON(http.StatusConflict, gin.H{"status": "Foto already in album"})
    return
}

	// Associate
	if err := db.Model(&album).Association("Fotos").Append(&foto); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add foto to album"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Ok, foto added to album successfully"})
}
