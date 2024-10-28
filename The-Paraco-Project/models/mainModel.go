package models

import (
    "database/sql"
    "errors"
    "time"
)

// UpdateBalance updates the user's balance and logs the transaction
func UpdateBalance(db *sql.DB, email string, used float64, platform uint8) error {
    // Define platform types
    platformTypes := map[uint8]string{
        1: "SMS",
        2: "RCS",
        3: "WhatsApp",
        4: "Email",
    }

    // Get the platform name from the map
    platformName, exists := platformTypes[platform]
    if !exists {
        return errors.New("invalid platform type")
    }

    // Start a transaction
    tx, err := db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    // Get the current balance
    var balance float64
    err = tx.QueryRow("SELECT balance FROM users WHERE email = ?", email).Scan(&balance)
    if err == sql.ErrNoRows {
        return errors.New("user not found")
    } else if err != nil {
        return err
    }

    // Check if there is enough balance
    if balance < used {
        return errors.New("insufficient balance")
    }

    // Update the balance
    newBalance := balance - used
    _, err = tx.Exec("UPDATE users SET balance = ? WHERE email = ?", newBalance, email)
    if err != nil {
        return err
    }

    // Insert the transaction log
    _, err = tx.Exec(
        "INSERT INTO logs (email, platform, spends, log_date) VALUES (?, ?, ?, ?)",
        email, platformName, used, time.Now(),
    )
    if err != nil {
        return err
    }

    // Commit the transaction
    if err := tx.Commit(); err != nil {
        return err
    }

    return nil
}
