package main

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/sessions"
    "github.com/gin-contrib/sessions/cookie"
    "The-Paraco-Project/config"
    "The-Paraco-Project/controllers"
)

type Address struct {
    Street string `json:"street"`
    City   string `json:"city"`
    State  string `json:"state"`
    Zip    string `json:"zip"`
}

type User struct {
    Name    string  `json:"name"`
    Gender  string  `json:"gender"`
    Address Address `json:"address"`
    Balance float64 `json:"balance"`
    Used    int     `json:"used"`
}

func userHandler(c *gin.Context) {
    var user User
    if err := c.ShouldBindJSON(&user); err != nil { //[arses the json file]
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    // Calculate new balance
    newBalance := user.Balance - float64(user.Used)

    // Create a response object
    response := gin.H{
        "name":    user.Name,
        "gender":  user.Gender,
        "address": user.Address,
        "new_balance": newBalance,
    }

    c.JSON(http.StatusCreated, response) // Return the new balance and other user data
}


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

    // Define the GET /signup route to serve the signup page
    r.GET("/signup", func(c *gin.Context) {
        c.HTML(http.StatusOK, "signup.html", nil) // Serve the signup HTML form
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

    // Define the /users route for handling user data
    r.POST("/users", userHandler)

    // Start the server
    r.Run(":8080")
}
