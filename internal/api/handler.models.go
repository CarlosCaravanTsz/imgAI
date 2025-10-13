package api

import (
	"github.com/gin-gonic/gin"
)

type UsersRouteHandlers struct{}

func (h *UsersRouteHandlers) RegisterUser(c *gin.Context)
func (h *UsersRouteHandlers) LoginUser(c *gin.Context)
func (h *UsersRouteHandlers) LogoutUser(c *gin.Context)

type FotosRouteHandlers struct{}

func (h *FotosRouteHandlers) SubirFotos(c *gin.Context)
func (h *FotosRouteHandlers) ListarFotos(c *gin.Context)
func (h *FotosRouteHandlers) ListarUnaFoto(c *gin.Context)
func (h *FotosRouteHandlers) DescargarFoto(c *gin.Context)
func (h *FotosRouteHandlers) EliminarFoto(c *gin.Context)
func (h *FotosRouteHandlers) ToggleFavorito(c *gin.Context)
func (h *FotosRouteHandlers) ListarFavoritos(c *gin.Context)
func (h *FotosRouteHandlers) AgregarFotoaAlbum(c *gin.Context)

type AlbumesRouteHandlers struct{}

func (h *FotosRouteHandlers) CrearAlbum(c *gin.Context)
func (h *FotosRouteHandlers) ListarAlbumes(c *gin.Context)
func (h *FotosRouteHandlers) ListarFotosAlbum(c *gin.Context)
func (h *FotosRouteHandlers) EliminarAlbum(c *gin.Context)
func (h *FotosRouteHandlers) QuitarFotoDeAlbum(c *gin.Context)
