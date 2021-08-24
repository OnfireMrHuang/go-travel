package routers

import (
	"go-travel/blog-service/global"
	"go-travel/blog-service/internal/middleware"
	"go-travel/blog-service/internal/routers/api"
	v1 "go-travel/blog-service/internal/routers/api/v1"
	"go-travel/blog-service/pkg/limiter"
	"net/http"
	"time"

	_ "go-travel/blog-service/docs"

	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	"github.com/gin-gonic/gin"
)

var methodLimiters = limiter.NewMethodLimiter().AddBuckets(limiter.LimiterBucketRule{
	Key:          "/auth",
	FillInterval: time.Second,
	Capacity:     10,
	Quantum:      10,
})

func NewRouter() *gin.Engine {
	// 引入服务
	article := v1.NewArticle()
	tag := v1.NewTag()
	upload := NewUpload()

	r := gin.New()
	r.Use(gin.Logger())
	if global.ServerSetting.RunMode == "debug" {
		r.Use(gin.Recovery())
	} else {
		r.Use(middleware.Recovery())
	}
	r.Use(middleware.AccessLog())
	r.Use(middleware.ContextTimeout(60 * time.Second))
	r.Use(middleware.Translations())
	r.Use(middleware.AccessLog())
	r.Use(middleware.Tracing())

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.POST("/upload/file", upload.UploadFile)
	r.StaticFS("/static", http.Dir(global.AppSetting.UploadSavePath))

	r.POST("/auth", api.GetAuth)

	apiV1 := r.Group("/api/v1")
	apiV1.Use(middleware.JWT())
	{
		apiV1.POST("/tags", tag.Create)
		apiV1.DELETE("/tags/:id", tag.Delete)
		apiV1.PUT("/tags/:id", tag.Update)
		apiV1.PATCH("/tags/:id/state", tag.Update)
		apiV1.GET("/tags", tag.List)

		apiV1.POST("/articles", article.Create)
		apiV1.DELETE("/articles/:id", article.Delete)
		apiV1.PUT("/articles/:id", article.Update)
		apiV1.PATCH("/articles/:id/state", article.Update)
		apiV1.GET("/articles/:id", article.Get)
		apiV1.GET("/articles", article.List)
	}
	return r
}
