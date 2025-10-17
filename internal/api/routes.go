package api

import (
	_ "mime/multipart"
	_ "net/http"

	h "github.com/CarlosCaravanTsz/imgAI/internal/api/handlers"
	m "github.com/CarlosCaravanTsz/imgAI/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	usersHandler := h.UsersRouteHandlers{}
	fotosHandler := h.FotosRouteHandlers{}
	albumesHandler := h.AlbumesRouteHandlers{}

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "Ok"})
		// c.Redirect(http.StatusMovedPermanently, "https:127.0.0.1/5173")
	})

	usuarios := r.Group("/api/users")
	//usuarios.Use(m.RequireAuthLogin)
	{
		usuarios.POST("/auth/register", usersHandler.RegisterUser) // register
		usuarios.POST("/auth/login", usersHandler.LoginUser)       // login
		// usuarios.POST("/auth/logout", usersHandler.LogoutUser) // logut
	}

	fotos := r.Group("/api/fotos")
	fotos.Use(m.RequireAuth)
	{
		fotos.POST("/", fotosHandler.SubirFotos)        // subir una o muchas fotos
		fotos.GET("/", fotosHandler.ListarFotos)        // listar fotos order by timestamp - ?email=email
		fotos.GET("/foto/", fotosHandler.ListarUnaFoto) // obtener toda la info de una foto - ?email=email&id=1
		// fotos.GET("/:fotoid/download", fotosHandler.DescargarFoto)                // descargar una foto
		fotos.DELETE("/:id", fotosHandler.EliminarFoto)                       // eliminar una foto
		fotos.PUT("/:id/favoritos", fotosHandler.ToggleFavorito)              // agregar una foto a favoritos (o quitar)
		fotos.GET("/favoritos", fotosHandler.ListarFavoritos)                 // listar fotos favoritos
		fotos.POST("/:fotoid/album/:albumid", fotosHandler.AgregarFotoaAlbum) // agregar fotoid a albumid

	}

	albumes := r.Group("/api/albumes")
	albumes.Use(m.RequireAuth)
	{
		albumes.POST("/", albumesHandler.CrearAlbum)                                // crear un album
		albumes.GET("/", albumesHandler.ListarAlbumes)                              // ?email                       // listar albumes
		albumes.GET("/:albumid/fotos", albumesHandler.ListarFotosAlbum)             // ?email     // listar las fotos del album id
		albumes.DELETE("/:albumid", albumesHandler.EliminarAlbum)                   // eliminar el album
		albumes.DELETE("/:albumid/fotos/:fotoid", albumesHandler.QuitarFotoDeAlbum) // quitar fotoid del album albumid
	}
}
