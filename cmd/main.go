package main

import (
	"database/sql"

	"github.com/viniferr33/img-processor/internal/config"
	"github.com/viniferr33/img-processor/internal/image"
	"github.com/viniferr33/img-processor/internal/jwt"
	"github.com/viniferr33/img-processor/internal/postgres"
	"github.com/viniferr33/img-processor/internal/restful"
	"github.com/viniferr33/img-processor/internal/user"
	"github.com/viniferr33/img-processor/pkg/database"
	"github.com/viniferr33/img-processor/pkg/logger"
	"github.com/viniferr33/img-processor/pkg/minio"
	"github.com/viniferr33/img-processor/pkg/server"
)

func main() {
	logger.Init(logger.Config{
		Development: true,
		Level:       "debug",
		Color:       true,
	})
	defer logger.Sync()

	// Setup infrastructure (e.g., database, cache)
	db := setupDatabase()
	defer db.Close()

	// Setup Repositories
	userRepo := postgres.NewUserRepository(db)
	imgRepo := postgres.NewImageRepository(db)
	objStorage := minio.NewMinIO(
		config.MinIOEndpoint,
		config.MinIOAccessKeyID,
		config.MinIOSecretAccessKey,
		config.MinIOUseSSL,
	)

	// Setup Services
	userService := user.NewUserService(userRepo)
	jwtService := jwt.NewJwtTokenService(
		config.JWTSecretKey,
		config.JWTIssuer,
		config.JWTExpirationSec,
	)

	imgService := image.NewImageService(imgRepo, objStorage, config.MinIODefaultBucket)

	// Setup Handlers
	authHandler := restful.NewAuthHandler(*userService, *jwtService)
	imageHandler := restful.NewImageHandler(*imgService)

	// Setup and start the HTTP server
	router := restful.NewRouter(authHandler, imageHandler)
	server.Init(server.Config{
		Host:    config.ServerHost,
		Port:    config.ServerPort,
		Handler: router,
	})

	server.Start()
}

func setupDatabase() *sql.DB {
	db, err := database.NewConnection(&database.Config{
		Host:     config.DBHost,
		Port:     config.DBPort,
		User:     config.DBUser,
		Password: config.DBPassword,
		DBName:   config.DBName,
		SSLMode:  config.DBSSLMode,
	})
	if err != nil {
		panic(err)
	}

	if err := database.RunMigrations(db, config.MigrationsPath); err != nil {
		panic(err)
	}

	return db
}
