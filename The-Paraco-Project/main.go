package main

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/sessions"
    "github.com/gin-contrib/sessions/cookie"
    "The-Paraco-Project/config"
    "The-Paraco-Project/controllers"
)

func main() {
    // Initialize the database
    config.InitDB()

    // Initialize the Gin router
    r := gin.Default()

    // Load HTML templates
    r.LoadHTMLGlob("views/*.html")

    // Setup session store
    store := cookie.NewStore([]byte("secret"))
    r.Use(sessions.Sessions("mysession", store))

    // Serve static files (CSS, JS)
    r.Static("/static", "./static")

    // Define routes
    r.GET("/", func(c *gin.Context) {
        c.HTML(http.StatusOK, "welcome.html", nil)
    })
    r.POST("/login", controllers.Login)
    r.POST("/signup", controllers.Signup)
    r.POST("/logout", controllers.Logout)

    // Serve the landing page
    r.GET("/landing", func(c *gin.Context) {
        session := sessions.Default(c)
        username := session.Get("username")
        if username == nil {
            c.Redirect(http.StatusFound, "/")
            return
        }
        c.HTML(http.StatusOK, "landing.html", gin.H{
            "username": username,
        })
    })

    // Start the server
    r.Run(":8080")
}
