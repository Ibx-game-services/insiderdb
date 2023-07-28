package routers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/bwmarrin/snowflake"
	// "github.com/ibx34/furry/middleware"
	"insider.db/m/models"
	"insider.db/m/internal"
)

type App struct {
	Database *models.Database
	Config   internal.Config
	SnowflakeGeneratorNode *snowflake.Node
}

// func (app *App) SessionMiddleware(ctx *gin.Context) {
// 	middleware.InnerSessionMiddleware(ctx, app.Redis, app.Database)
// }

func CORSMiddleware(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-Password")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "PATCH, POST, OPTIONS, GET, PUT")

	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(204)
		return
	}

	c.Next()
}

func Init(database *models.Database, config internal.Config, snowflake *snowflake.Node) {
	app := App{Database: database, Config: config, SnowflakeGeneratorNode: snowflake}
	if config.Mode != nil && *config.Mode == "prod" || *config.Mode == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(CORSMiddleware)
	r.MaxMultipartMemory = 10 << 50
	r.Use(gin.Recovery())
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s (%s) %d --> %s\n",
			param.Method,
			param.ClientIP,
			param.StatusCode,
			param.Path,
		)
	}))
	// Routes here do NOT require authentication and can be
	// accessed by anyone. Some may have weird reqiurements and
	// will require authentication in some cases and not in others.
	// One example is `GET /users/:id`: This route requires authentication
	// when the requested user is marked as Private, or NSFW.

	// users.Init(r, app.Database, app.Config, app.S3, app.Mg, app.Redis)

	// r.GET("/posts", app.GetPosts)
	// // r.Use(app.SessionMiddleware)
	// r.GET("/users/me/posts", app.GetMyPosts)
	r.POST("/insider", app.SubmitInsider)
	// r.GET("/qr", app.GetQrCode)
	// r.GET("/qr/validate", app.VerifyCode)
	// r.POST("/verify/captcha", app.CheckCaptcha)
	// r.POST("/verify/email", app.VerifyUserEmail)
	r.Run()
}