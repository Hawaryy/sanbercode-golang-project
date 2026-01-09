package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Bioskop struct {
	ID     int     `json:"id"`
	Nama   string  `json:"nama"`
	Lokasi string  `json:"lokasi"`
	Rating float64 `json:"rating"`
}

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("postgres", "user=postgres password=postgres dbname=bioskop_db host=localhost port=5432 sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	db.Exec(`CREATE TABLE IF NOT EXISTS bioskop (
		id SERIAL PRIMARY KEY,
		nama VARCHAR(255),
		lokasi VARCHAR(255),
		rating FLOAT
	)`)

	log.Println("✓ Database siap")
}

func postBioskop(c *gin.Context) {
	var bioskop Bioskop
	c.BindJSON(&bioskop)

	var id int
	db.QueryRow("INSERT INTO bioskop (nama, lokasi, rating) VALUES ($1, $2, $3) RETURNING id",
		bioskop.Nama, bioskop.Lokasi, bioskop.Rating).Scan(&id)

	bioskop.ID = id
	c.JSON(http.StatusCreated, bioskop)
}

func getBioskop(c *gin.Context) {
	rows, _ := db.Query("SELECT id, nama, lokasi, rating FROM bioskop")
	defer rows.Close()

	var bioskopList []Bioskop
	for rows.Next() {
		var b Bioskop
		rows.Scan(&b.ID, &b.Nama, &b.Lokasi, &b.Rating)
		bioskopList = append(bioskopList, b)
	}

	c.JSON(http.StatusOK, bioskopList)
}

func getBioskopByID(c *gin.Context) {
	id := c.Param("id")
	var b Bioskop

	db.QueryRow("SELECT id, nama, lokasi, rating FROM bioskop WHERE id = $1", id).
		Scan(&b.ID, &b.Nama, &b.Lokasi, &b.Rating)

	c.JSON(http.StatusOK, b)
}

func putBioskop(c *gin.Context) {
	id := c.Param("id")
	var b Bioskop
	c.BindJSON(&b)

	db.Exec("UPDATE bioskop SET nama = $1, lokasi = $2, rating = $3 WHERE id = $4",
		b.Nama, b.Lokasi, b.Rating, id)

	c.JSON(http.StatusOK, b)
}

func deleteBioskop(c *gin.Context) {
	id := c.Param("id")
	db.Exec("DELETE FROM bioskop WHERE id = $1", id)

	c.JSON(http.StatusOK, gin.H{"message": "Berhasil dihapus"})
}

func main() {
	router := gin.Default()

	router.POST("/bioskop", postBioskop)
	router.GET("/bioskop", getBioskop)
	router.GET("/bioskop/:id", getBioskopByID)
	router.PUT("/bioskop/:id", putBioskop)
	router.DELETE("/bioskop/:id", deleteBioskop)

	router.Run(":8080")
}