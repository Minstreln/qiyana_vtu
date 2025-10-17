package main

import (
	"crypto/tls"
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	mw "qiyana_vtu/internal/api/middlewares"
	"qiyana_vtu/internal/api/routers"
	"qiyana_vtu/internal/repositories/sqlconnect"
	"qiyana_vtu/pkg/utils"
	"time"

	"github.com/joho/godotenv"
)

//go:embed .env
var envFile embed.FS

func loadEnvFromEmbeddedFile() {
	// read the embedded .env
	content, err := envFile.ReadFile(".env")
	if err != nil {
		log.Fatalf("Error reading from .env file: %v", err)
	}

	// create a temp file to load the env variables
	tempFile, err := os.CreateTemp("", ".env")
	if err != nil {
		log.Fatalf("Error creating temp .env file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// write content in the temp file
	_, err = tempFile.Write(content)
	if err != nil {
		log.Fatalf("Error writting to temp .env file: %v", err)
	}

	err = tempFile.Close()
	if err != nil {
		log.Fatalf("Error closing temp .env file: %v", err)
	}

	// load env variables from temp file
	err = godotenv.Load(tempFile.Name())
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func main() {
	// err := godotenv.Load()
	// if err != nil {
	// 	return
	// }

	loadEnvFromEmbeddedFile()

	err := sqlconnect.ConnectDb()
	if err != nil {
		utils.Logger.Fatal("DB connection failed: ", err)
	}

	utils.InitLogger()

	port := os.Getenv("SERVER_PORT")

	cert := os.Getenv("CERT_FILE")
	key := os.Getenv("KEY_FILE")

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS10,
	}

	rl := mw.NewRateLimiter(5, time.Minute)

	hppOptions := mw.HPPOptions{
		CheckQuery:                  true,
		CheckBody:                   true,
		CheckBodyOnlyForContentType: "application/x-www-form-urlencoded",
		Whitelist:                   []string{"sortBy", "sortOrder", "name", "age", "class"},
	}

	router := routers.MainRouter()
	jwtMiddleware := mw.MiddlewaresExcludePaths(mw.JWTMiddleware, "/users/signup", "/users/login", "/users/confirmotp", "/users/resendotp", "/users/forgotpassword", "/execs/resetpassword/reset", "/wallet/webhook")

	secureMux := utils.ApplyMiddlewares(router, mw.SecurityHeaders, mw.Compression, mw.Hpp(hppOptions), jwtMiddleware, mw.ResponseTimeMiddleware, rl.Middleware, mw.Cors)
	// secureMux := jwtMiddleware(mw.SecurityHeaders(router))

	server := &http.Server{
		Addr:      port,
		Handler:   secureMux,
		TLSConfig: tlsConfig,
	}

	fmt.Println("Server is running on port", port)
	err = server.ListenAndServeTLS(cert, key)
	if err != nil {
		log.Fatalln("Error starting the server", err)
	}

}
