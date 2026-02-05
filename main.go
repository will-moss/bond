package main

import (
	"fmt"
	"net/http"
	"strconv"
	"os"
	"strings"
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	qrcode "github.com/skip2/go-qrcode"
)

// Alias for os.GetEnv, with support for fallback value, and boolean normalization
func getEnv(key string, fallback ...string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		if len(fallback) > 0 {
			value = fallback[0]
		} else {
			value = ""
		}
	} else {
		// Quotes removal
		value = strings.Trim(value, "\"")

		// Boolean normalization
		mapping := map[string]string{
			"0":     "FALSE",
			"off":   "FALSE",
			"false": "FALSE",
			"1":     "TRUE",
			"on":    "TRUE",
			"true":  "TRUE",
		}
		normalized, isBool := mapping[strings.ToLower(value)]
		if isBool {
			value = normalized
		}
	}

	return value
}

// Sends an HTTP request to the /health endpoint on the running server
// and returns an exit code. This assumes that a first "./bond" program runs,
// and this was called via "./bond --healthcheck" on the exact same machine / container
func runHealthcheck() int {
	// Reconstruct the URL to reach the server locally
	protocol := "http"
	if getEnv("SSL") == "TRUE" {
		protocol = "https"
	}

	healthURL := fmt.Sprintf("%s://localhost:%s/health", protocol, getEnv("PORT"))

	// Call the health check URL
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(healthURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Health check failed: %v\n", err)
		return 1
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "Health check failed: status %d\n", resp.StatusCode)
		return 1
	}

	fmt.Println("OK")
	return 0
}

// Entrypoint
func main() {
	godotenv.Load("default.env")

	// Load custom settings via .env file
	err := godotenv.Overload(".env")
	if err != nil {
		log.Print("No .env file provided, will continue with system env")
	}

	// Run a simple health check when the "--healthcheck" arg is provided on the commandline
	if len(os.Args) > 1 && os.Args[1] == "--healthcheck" {
		exitCode := runHealthcheck()
		os.Exit(exitCode)
	}

	// Define the keyword-int association for QR code recovery levels
	recoveryLevels := map[string]qrcode.RecoveryLevel{
		"LOW": qrcode.Low,
		"MEDIUM": qrcode.Medium,
		"HIGH": qrcode.High,
		"HIGHEST": qrcode.Highest,
	}

	// Instantiate server
	app := chi.NewRouter()

	// Set up basic middleware
	if getEnv("ENABLE_LOGS") == "TRUE" {
		app.Use(middleware.Logger)
	}
	app.Use(middleware.Recoverer)

	// CORS-specific
	app.Options("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET")
	})

	// GET /
	app.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET")

		secret := r.URL.Query().Get("secret")
		content := r.URL.Query().Get("content")
		size := r.URL.Query().Get("size")

		// Ensure secret is provided and matches the one in store
		if secret != getEnv("SECRET") {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		// Ensure Size and Content are provided
		if content == "" || size == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		realSize, err := strconv.Atoi(size)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		maxSize, _ := strconv.Atoi(getEnv("MAX_SIZE"))
		if realSize <= 0 || realSize > maxSize {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Generate the QR code
		png, err := qrcode.Encode(content, recoveryLevels[strings.ToUpper(getEnv("RECOVERY_LEVEL"))], realSize)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Output the QR code
		_, err = w.Write(png)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Print(err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/octet-stream")
	})

	// GET /health
	app.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET")
		w.WriteHeader(http.StatusOK)
	})

	log.Printf("Server starting on port %s", getEnv("PORT"))
	if getEnv("SSL") == "TRUE" {
		err = http.ListenAndServeTLS(fmt.Sprintf(":%s", getEnv("PORT")), "certificate.pem", "key.pem", app)
	} else {
		err = http.ListenAndServe(fmt.Sprintf(":%s", getEnv("PORT")), app)
	}

	if err != nil {
		log.Print(err.Error())
	}
}
