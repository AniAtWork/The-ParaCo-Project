package controllers

import (
    "database/sql"
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/sessions"
    "golang.org/x/crypto/bcrypt"
    "log"
)

// Login handles the user login
func Login(c *gin.Context, db *sql.DB) {
    session := sessions.Default(c)
    username := c.PostForm("username")
    password := c.PostForm("password")

    var dbPassword string
    err := db.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&dbPassword)
    if err != nil {
        c.String(http.StatusUnauthorized, "Invalid login credentials")
        return
    }

    if err := bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(password)); err == nil {
        session.Set("username", username)
        session.Save()
        c.Redirect(http.StatusFound, "/landing")
    } else {
        c.String(http.StatusUnauthorized, "Invalid login credentials")
    }
}

// Signup handles user registration
func Signup(c *gin.Context, db *sql.DB) {
    username := c.PostForm("username")
    email := c.PostForm("email")
    password := c.PostForm("password")

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        c.String(http.StatusInternalServerError, "Error creating account")
        return
    }

    _, err = db.Exec("INSERT INTO users (username, email, password) VALUES (?, ?, ?)", username, email, string(hashedPassword))
    if err != nil {
        log.Fatal(err)
    }

    c.Redirect(http.StatusFound, "/")
}

// Logout handles user logout
func Logout(c *gin.Context) {
    session := sessions.Default(c)
    session.Clear()
    session.Save()
    c.Redirect(http.StatusFound, "/")
}
