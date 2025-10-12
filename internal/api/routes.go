package api

import (
	_ "mime/multipart"
	_ "net/http"

	h "github.com/CarlosCaravanTsz/imgAI/internal/api/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	usersHandler := h.UsersRouteHandlers{}
	fotosHandler := h.FotosRouteHandlers{}
	albumesHandler := h.AlbumesRouteHandlers{}

	r.GET("/", func(c *gin.Context) {
		c.JSON(200,gin.H{"status": "Ok",})
		//c.Redirect(http.StatusMovedPermanently, "https:127.0.0.1/5173")
	}) 

	usuarios := r.Group("/api/users")
	{
		usuarios.POST("/auth/register", usersHandler.RegisterUser) // register
		usuarios.POST("/auth/login", usersHandler.LoginUser) // login
		usuarios.POST("/auth/logout", usersHandler.LogoutUser) // logut
	}

	fotos := r.Group("/api/fotos")
	{
		fotos.POST("/fotos", fotosHandler.SubirFotos)                               // subir una o muchas fotos
		fotos.GET("/fotos", fotosHandler.ListarFotos)                               // listar fotos order by timestamp
		fotos.GET("/fotos/:fotoid", fotosHandler.ListarUnaFoto)                         // obtener toda la info de una foto
		fotos.GET("/fotos/:fotoid/download", fotosHandler.DescargarFoto)                // descargar una foto
		fotos.DELETE("/fotos/:id", fotosHandler.EliminarFoto)                       // eliminar una foto
		fotos.PUT("/fotos/:fotoid/favoritos", fotosHandler.ToggleFavorito)              // agregar una foto a favoritos (o quitar)
		fotos.GET("/fotos/favoritos", fotosHandler.ListarFavoritos)                 // listar fotos favoritos
		fotos.POST("/fotos/:fotoid/album/:albumid", fotosHandler.AgregarFotoaAlbum) // agregar fotoid a albumid

	}

	albumes := r.Group("/api/albumes")
	{
		albumes.POST("/albumes", albumesHandler.CrearAlbum)                                 // crear un album
		albumes.GET("/albumes", albumesHandler.ListarAlbumes)                               // listar albumes
		albumes.GET("/albumes/:albumid/fotos", albumesHandler.ListarFotosAlbum)                  // listar las fotos del album id
		albumes.DELETE("/albumes/:albumid", albumesHandler.EliminarAlbum)                        // eliminar el album
		albumes.DELETE("/albumes/:albumid/fotos/:fotoid", albumesHandler.QuitarFotoDeAlbum) // quitar fotoid del album albumid
	}
}
