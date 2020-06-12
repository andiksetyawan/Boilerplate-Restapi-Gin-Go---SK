package main

import (
    "./config"
    "./routes"
    "./middleware"
    "github.com/gin-gonic/gin"
    "github.com/subosito/gotenv"
)

func main() {
    //set up database
    config.InitDB()
    defer config.DB.Close()
    gotenv.Load()

    //set up routing/router
    router := gin.Default()

    v1 := router.Group("/api/v1/")
    {
        v1.GET("/auth/:provider", routes.RedirectHandler)
        v1.GET("/auth/:provider/callback", routes.CallbackHandler)


        v1.GET("/profile", middleware.IsAuth(), routes.GetProfile)


        v1.GET("/article/:slug", routes.GetArticle)
        articles := v1.Group("/articles")
        {
            articles.GET("/", routes.GetHome)
            articles.POST("/", middleware.IsAuth(), routes.PostArticle)
            articles.GET("/tag/:tag", routes.GetArticleByTag)
            articles.PUT("/update/:id", middleware.IsAuth(), routes.UpdateArticle)
            articles.DELETE("/delete/:id", middleware.IsAdmin(), routes.DeleteArticle)
        }
    }

    router.Run()
}