package main

import (
    "database/sql"  // Import sql package
    "net/http"
    "log"
    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/sessions"
    "github.com/gin-contrib/sessions/cookie"
    "The-Paraco-Project/config"
    "The-Paraco-Project/controllers"
    "The-Paraco-Project/models"
)

type Address struct {
    Street string `json:"street"`
    City   string `json:"city"`
    State  string `json:"state"`
    Zip    string `json:"zip"`
}

type User struct {
    Name    string    `json:"name"`
    Gender  string    `json:"gender"`
    Address Address   `json:"address"`
    Used    float64   `json:"used"`
    Metadata Metadata `json:"metadata"`
}

type Metadata struct {
    Email string   `json:"email"`
    Platform uint8 `json:"platform"`
}

func userHandler(c *gin.Context, db *sql.DB) {
    var user User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if user.Metadata.Email == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "The 'email' field is required"})
        return
    }

    // Call UpdateBalance and handle errors
    if err := models.UpdateBalance(db, user.Metadata.Email, user.Used, user.Metadata.Platform); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"message": "Securly deducted balance"})
}

func main() {
    // Initialize the database
    if err := config.InitDB(); err != nil {
        log.Fatal("Failed to connect to the database:", err)
    }
    defer config.DB.Close() // Close the database connection when main exits

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

    r.GET("/signup", func(c *gin.Context) {
        c.HTML(http.StatusOK, "signup.html", nil)
    })

    r.POST("/login", func(c *gin.Context) {
        controllers.Login(c, config.DB) 
    })
    r.POST("/signup", func(c *gin.Context) {
        controllers.Signup(c, config.DB) 
    })
    r.POST("/logout", controllers.Logout)

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

    err := r.SetTrustedProxies([]string{"192.168.1.0/24"}) 
    if err != nil {
        panic(err)
    }

    // Define the /users route for handling user data
    r.POST("/users", func(c *gin.Context) {
        userHandler(c, config.DB)
    })

    // Start the server
    r.Run(":8080")
}
