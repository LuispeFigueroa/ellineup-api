package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/LuispeFigueroa/ellineup-api/handlers"
	"github.com/LuispeFigueroa/ellineup-api/middleware"
	"github.com/LuispeFigueroa/ellineup-api/router"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error cargando .env")
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error conectando a la base de datos:", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal("Ping fallido:", err)
	}

	log.Println("Conectado a PostgreSQL")

	handlers.DB = db

	r := gin.Default()
	r.Use(middleware.CORS())
	router.Setup(r)

	r.Run(":8080")
}
