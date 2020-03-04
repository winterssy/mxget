package server

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/winterssy/mxget/pkg/provider"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

func searchSongs(c *gin.Context) {
	platform := c.Param("platform")
	client, err := provider.GetClient(platform)
	if err != nil {
		c.JSON(400, Response{
			Code: 400,
			Msg:  err.Error(),
		})
		return
	}

	data, err := client.SearchSongs(context.Background(), c.Param("keyword"))
	if err != nil {
		c.JSON(500, Response{
			Code: 500,
			Msg:  err.Error(),
		})
		return
	}

	c.JSON(200, Response{
		Code: 200,
		Data: data,
	})
}

func getSong(c *gin.Context) {
	platform := c.Param("platform")
	client, err := provider.GetClient(platform)
	if err != nil {
		c.JSON(400, Response{
			Code: 400,
			Msg:  err.Error(),
		})
		return
	}

	data, err := client.GetSong(context.Background(), c.Param("id"))
	if err != nil {
		c.JSON(500, Response{
			Code: 500,
			Msg:  err.Error(),
		})
		return
	}

	c.JSON(200, Response{
		Code: 200,
		Data: data,
	})
}

func getArtist(c *gin.Context) {
	platform := c.Param("platform")
	client, err := provider.GetClient(platform)
	if err != nil {
		c.JSON(400, Response{
			Code: 400,
			Msg:  err.Error(),
		})
		return
	}

	data, err := client.GetArtist(context.Background(), c.Param("id"))
	if err != nil {
		c.JSON(500, Response{
			Code: 500,
			Msg:  err.Error(),
		})
		return
	}

	c.JSON(200, Response{
		Code: 200,
		Data: data,
	})
}

func getAlbum(c *gin.Context) {
	platform := c.Param("platform")
	client, err := provider.GetClient(platform)
	if err != nil {
		c.JSON(400, Response{
			Code: 400,
			Msg:  err.Error(),
		})
		return
	}

	data, err := client.GetAlbum(context.Background(), c.Param("id"))
	if err != nil {
		c.JSON(500, Response{
			Code: 500,
			Msg:  err.Error(),
		})
		return
	}

	c.JSON(200, Response{
		Code: 200,
		Data: data,
	})
}

func getPlaylist(c *gin.Context) {
	platform := c.Param("platform")
	client, err := provider.GetClient(platform)
	if err != nil {
		c.JSON(400, Response{
			Code: 400,
			Msg:  err.Error(),
		})
		return
	}

	data, err := client.GetPlaylist(context.Background(), c.Param("id"))
	if err != nil {
		c.JSON(500, Response{
			Code: 500,
			Msg:  err.Error(),
		})
		return
	}

	c.JSON(200, Response{
		Code: 200,
		Data: data,
	})
}

func Init(router *gin.Engine) {
	r := router.Group("/api")

	r.GET("/:platform/search/:keyword", searchSongs)
	r.GET("/:platform/song/:id", getSong)
	r.GET("/:platform/artist/:id", getArtist)
	r.GET("/:platform/album/:id", getAlbum)
	r.GET("/:platform/playlist/:id", getPlaylist)
}
