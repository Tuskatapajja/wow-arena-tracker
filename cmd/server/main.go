package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"

    _ "github.com/lib/pq"
)

type Match struct {
    Comp string `json:"comp"`
    Win  bool   `json:"win"`
    MMR  int    `json:"mmr"`
}

type Stats struct {
    TotalGames int     `json:"total_games"`
    Wins       int     `json:"wins"`
    Losses     int     `json:"losses"`
    WinRate    float64 `json:"win_rate"`
}

var db *sql.DB

func main() {
    // Настройки подключения к БД
    host := getEnv("DB_HOST", "localhost")
    port := getEnv("DB_PORT", "5432")
    user := getEnv("DB_USER", "wow")
    password := getEnv("DB_PASSWORD", "wow123")
    dbname := getEnv("DB_NAME", "arena")

    connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        host, port, user, password, dbname)

    var err error
    db, err = sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal("Open error:", err)
    }
    defer db.Close()

    err = db.Ping()
    if err != nil {
        log.Fatal("Ping error:", err)
    }
    log.Println("Connected to PostgreSQL!")

    // Создаём таблицу
    _, err = db.Exec(`CREATE TABLE IF NOT EXISTS matches (
        id SERIAL PRIMARY KEY,
        comp TEXT,
        win BOOLEAN,
        mmr INT,
        created_at TIMESTAMP DEFAULT NOW()
    )`)
    if err != nil {
        log.Fatal("Table error:", err)
    }

    // Маршруты
    http.HandleFunc("/api/match", handleMatch)
    http.HandleFunc("/api/stats", handleStats)

    log.Println("Server starting on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleMatch(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        w.WriteHeader(http.StatusMethodNotAllowed)
        return
    }

    var match Match
    if err := json.NewDecoder(r.Body).Decode(&match); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    _, err := db.Exec("INSERT INTO matches (comp, win, mmr) VALUES ($1, $2, $3)",
        match.Comp, match.Win, match.MMR)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{"status": "match saved"})
}

func handleStats(w http.ResponseWriter, r *http.Request) {
    var total, wins int
    row := db.QueryRow("SELECT COUNT(*), COALESCE(SUM(CASE WHEN win THEN 1 ELSE 0 END), 0) FROM matches")
    err := row.Scan(&total, &wins)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    losses := total - wins
    winRate := 0.0
    if total > 0 {
        winRate = float64(wins) / float64(total) * 100
    }

    stats := Stats{
        TotalGames: total,
        Wins:       wins,
        Losses:     losses,
        WinRate:    winRate,
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(stats)
}

func getEnv(key, defaultVal string) string {
    if val := os.Getenv(key); val != "" {
        return val
    }
    return defaultVal
}
