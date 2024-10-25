package initializers

import (
	"jwt/models"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DBConn *gorm.DB

func InitialierEnvVariable() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Unable to load env variables")
	}
	log.Println("Loaded Environment Variables")
}

func InitiazeDB() {
	var err error
	dsn := os.Getenv("DSN")
	DBConn, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println("Unable to connect to DataDase")
	}
	log.Println("Connected to DB..!")
}

func MigrateDB() {
	DBConn.AutoMigrate(&models.User{}, &models.Group{}, &models.Role{}, &models.RSAKeyPair{})
	log.Println("Finished AutoMigration..!")
}

// SeedRoles seeds the database with initial roles, including Admin
func SeedRoles() {
	adminRole := models.Role{Name: "admin"}
	userRole := models.Role{Name: "user"}

	// Check if roles exist, if not, create them
	if err := DBConn.FirstOrCreate(&adminRole, models.Role{Name: "admin"}).Error; err != nil {
		log.Fatal("Failed to seed roles: ", err)
	}
	if err := DBConn.FirstOrCreate(&userRole, models.Role{Name: "user"}).Error; err != nil {
		log.Fatal("Failed to seed roles: ", err)
	}

	adminGroup := models.Group{Name: "admin"}
	userGroup := models.Group{Name: "user"}

	// Check if Group exist, if not, create them
	if err := DBConn.FirstOrCreate(&adminGroup, models.Group{Name: "admin"}).Error; err != nil {
		log.Fatal("Failed to seed roles: ", err)
	}
	if err := DBConn.FirstOrCreate(&userGroup, models.Group{Name: "user"}).Error; err != nil {
		log.Fatal("Failed to seed roles: ", err)
	}
}
