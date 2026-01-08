package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Bioskop struct {
	ID       int     `json:"id"`
	Nama     string  `json:"nama"`
	Lokasi   string  `json:"lokasi"`
	Rating   float64 `json:"rating"`
}

var db *sql.DB

func init() {
	dsn := "user=postgres password=postgres dbname=bioskop_db host=localhost port=5432 sslmode=disable"
	
	var err error
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Gagal membuka koneksi database: %v", err)
	}

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS bioskop (
		id SERIAL PRIMARY KEY,
		nama VARCHAR(255) NOT NULL,
		lokasi VARCHAR(255) NOT NULL,
		rating FLOAT DEFAULT 0
	);
	`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Gagal membuat tabel: %v", err)
	}

	log.Println("✓ Tabel bioskop siap")
}

func postBioskop(c *gin.Context) {
	var bioskop Bioskop

	if err := c.ShouldBindJSON(&bioskop); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format JSON tidak valid"})
		return
	}

	if bioskop.Nama == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama tidak boleh kosong"})
		return
	}

	if bioskop.Lokasi == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Lokasi tidak boleh kosong"})
		return
	}

	sqlStatement := `
	INSERT INTO bioskop (nama, lokasi, rating)
	VALUES ($1, $2, $3)
	RETURNING id;
	`
	id := 0
	err := db.QueryRow(sqlStatement, bioskop.Nama, bioskop.Lokasi, bioskop.Rating).Scan(&id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menambah bioskop"})
		return
	}

	bioskop.ID = id
	c.JSON(http.StatusCreated, gin.H{
		"message": "Bioskop berhasil ditambahkan",
		"data":    bioskop,
	})
}

func getBioskop(c *gin.Context) {
	rows, err := db.Query("SELECT id, nama, lokasi, rating FROM bioskop")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data"})
		return
	}
	defer rows.Close()

	var bioskopList []Bioskop
	for rows.Next() {
		var bioskop Bioskop
		err := rows.Scan(&bioskop.ID, &bioskop.Nama, &bioskop.Lokasi, &bioskop.Rating)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membaca data"})
			return
		}
		bioskopList = append(bioskopList, bioskop)
	}

	c.JSON(http.StatusOK, gin.H{"data": bioskopList})
}

func main() {
	router := gin.Default()

	router.POST("/bioskop", postBioskop)
	router.GET("/bioskop", getBioskop)

	router.Run(":8080")
}