package handlers

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"sync"

	"github.com/CarlosCaravanTsz/imgAI/internal/database"
	log "github.com/CarlosCaravanTsz/imgAI/internal/logger"
	s3 "github.com/CarlosCaravanTsz/imgAI/internal/storage"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type FotoParams struct { // send as email the current user email
	Email string `form:"email" json:"email"`
}

type FotosArray struct {
	URL string `json:"url"`
}

type FotosRouteHandlers struct{}

func (h *FotosRouteHandlers) SubirFotos(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		log.LogError("Error loading multipart form", logrus.Fields{"error": err})
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	files := form.File["images[]"]

	var params FotoParams
	if err := c.ShouldBind(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Inicia Cliente s3:
	s3Client, err := s3.NewS3Client()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize S3 client"})
		return
	}

	// Inicializa conn a DB
	db, err := database.GetConnection()
	if err != nil {
		log.LogInfo("Error getting the db conn", logrus.Fields{})
	}

	// Valida que el usuario exista: Buscar el user con el email y obtener ID
	var usuario struct {
		ID     uint
		Nombre string
	}

	err = db.Model(&database.Usuario{}).Select("id", "nombre").Where("email = ?", params.Email).First(&usuario).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "User not found"})
		return
	}

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
				Path:     params.Email,
				Buffer:   fileBytes,
			}

			url, err := s3Client.Upload(foto)
			if err != nil {
				resultsChan <- uploadResult{Error: fmt.Errorf("upload failed for %s: %v", fh.Filename, err)}
				return
			}

			log.LogInfo("Uploaded to S3", logrus.Fields{"url": url})

			//  Obtener la URL de S3 y guardar la informacion de la Foto en la BD
			dbFoto := database.Foto{
				UsuarioID:   usuario.ID,
				Nombre:      fileHeader.Filename,
				Descripcion: "",
				URLArchivo:  url,
				TamanoBytes: int64(len(fileBytes)),
				Formato:     "image",
			}

			if err := db.Create(&dbFoto).Error; err != nil {
				resultsChan <- uploadResult{Error: fmt.Errorf("db inserted failed for %s: %v", fh.Filename, err)}
				return
			}

			resultsChan <- uploadResult{URL: url, Filename: fh.Filename, Error: nil}

			c.JSON(http.StatusOK, gin.H{"status": "Ok, Uploaded to S3"})
		}(fileHeader)
	}

	wg.Wait()
	close(resultsChan)

	var uploaded []string
	var failed []string

	for r := range resultsChan {
		if r.Error != nil {
			log.LogError("Upload error ocurred", logrus.Fields{"file": r.Filename, "error": r.Error})
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
	email := c.Query("email")

	var fotos []database.Foto

	db, err := database.GetConnection()
	if err != nil {
		log.LogInfo("Error getting the db conn", logrus.Fields{})
	}

	// Query only photos for this user
	err = db.Joins("JOIN usuarios ON usuarios.id = fotos.usuario_id").
		Where("usuarios.email = ?", email).
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
	email := c.Query("email")
	fotoid := c.Query("id")

	var foto database.Foto

	type userInfo struct {
		ID     uint
		Nombre string
	}
	var usuario userInfo
	// Buscar el user con el email y obtener ID

	db, err := database.GetConnection()
	if err != nil {
		log.LogInfo("Error getting the db conn", logrus.Fields{})
	}

	err = db.Model(&database.Usuario{}).Select("id", "nombre").Where("email = ?", email).First(&usuario).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "User not found"})
		return
	}

	err = db.Where("usuario_id = ? AND id = ?", usuario.ID, fotoid).First(&foto).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Foto not found"})
		return
	}

	c.JSON(http.StatusOK, foto)
}

func (h *FotosRouteHandlers) EliminarFoto(c *gin.Context) {
	fotoid := c.Param("id")

	var params FotoParams
	if err := c.ShouldBind(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	type userInfo struct {
		ID     uint
		Nombre string
	}
	var usuario userInfo

	db, err := database.GetConnection()
	if err != nil {
		log.LogInfo("Error getting the db conn", logrus.Fields{})
	}

	err = db.Model(&database.Usuario{}).Select("id", "nombre").Where("email = ?", params.Email).First(&usuario).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "User not found"})
		return
	}

	result := db.Where("usuario_id = ? AND id = ?", usuario.ID, fotoid).Delete(&database.Foto{})
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

	db, err := database.GetConnection()
	if err != nil {
		log.LogInfo("Error getting the db conn", logrus.Fields{})
	}

	type fotoFavorito struct {
		Favorito bool
	}
	var fav fotoFavorito

	result := db.Model(&database.Foto{}).Where("id = ?", fotoid).Select("Favorito").Find(&fav)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Didnt find foto with that id"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"status": "No foto found"})
		return
	}

	err = db.Model(&database.Foto{}).Where("id = ?", fotoid).Update("Favorito", !fav.Favorito).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Couldnt toggle favorites 2"})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"status": "Ok, toggled favorites"})
}

func (h *FotosRouteHandlers) ListarFavoritos(c *gin.Context) {
	email := c.Query("email")

	var fotos []database.Foto

	db, err := database.GetConnection()
	if err != nil {
		log.LogInfo("Error getting the db conn", logrus.Fields{})
	}

	// Query only photos for this user
	err = db.Joins("JOIN usuarios ON usuarios.id = fotos.usuario_id").
		Where("usuarios.email = ? AND fotos.favorito=1", email).
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

	db, err := database.GetConnection()
	if err != nil {
		log.LogInfo("Error getting the db conn", logrus.Fields{})
	}

	// Get userID for validations
	var params FotoParams
	type userInfo struct {
		ID uint
	}
	var usuario userInfo

	if err := c.ShouldBind(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Buscar el user con el email y obtener ID
	err = db.Model(&database.Usuario{}).Select("id").Where("email = ?", params.Email).First(&usuario).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "User not found"})
		return
	}

	// Find album and foto
	var album database.Album
	if err := db.Where("id = ? AND usuario_id = ?", albumID, usuario.ID).First(&album).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Album not found or does not belong to user"})
		return
	}

	var foto database.Foto
	if err := db.Where("id = ? AND usuario_id = ?", fotoID, usuario.ID).First(&foto).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Foto not found or does not belong to user"})
		return
	}

	if db.Model(&album).Where("id = ?", foto.ID).Association("Fotos").Find(&foto) == nil {
		c.JSON(http.StatusOK, gin.H{"status": "Foto already in album"})
		return
	}

	// Associate
	if err := db.Model(&album).Association("Fotos").Append(&foto); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add foto to album"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Ok, foto added to album successfully"})
}
