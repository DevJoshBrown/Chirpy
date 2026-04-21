package main

//BUILD :  go build -o out && ./out
import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync/atomic"

	"github.com/DevJoshBrown/Chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Server struct {
}

// used to hold any stateful in-memory data that we want to keep track of.
type apiConfig struct {
	fileserverHits atomic.Int32
	database       *database.Queries
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1) // this processes before the handler runs
		next.ServeHTTP(w, r)      // run the wrapped handler
		// you could add processes to occurr after the handler like logging time to complete or errors etc.
	})
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(fmt.Appendf(nil,
		`
		<html>
  			<body>
			    <h1>Welcome, Chirpy Admin</h1>
			    <p>Chirpy has been visited %d times!</p>
			</body>
		</html>
		`, cfg.fileserverHits.Load()))
}

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type chirpRequest struct {
		Body string `json:"body"`
	}

	type ChirpResponse struct {
		//Valid       bool   `json:"valid"`
		CleanedBody string `json:"cleaned_body"`
	}

	type errorResponse struct {
		Error string `json:"error"`
	}

	w.Header().Set("Content-Type", "application/json")

	//	1. Decode:
	var req chirpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{Error: "Something went wrong"})
		return
	}

	//	2. Validate:
	if len(req.Body) > 140 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{Error: "Chirp is too long"})
		return
	}

	// 2.5 AntiProfanity
	badWords := []string{"kerfuffle", "sharbert", "fornax"}

	words := strings.Split(req.Body, " ")
	for i := range words {
		words[i] = ProfanFilter(words[i], badWords)
	}
	cleanWords := strings.Join(words, " ")

	// 	3. Success:
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ChirpResponse{CleanedBody: cleanWords})
}

func ProfanFilter(word string, badWords []string) string {
	for _, badWord := range badWords {
		if strings.ToLower(word) == badWord {
			return "****"
		}
	}
	return word
}

func (cfg *apiConfig) handlerResetCounter(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("failed to open postgres database at %s", dbURL)
	}
	dbQueries := database.New(db)

	fmt.Println("Running main.go")
	apiCfg := &apiConfig{
		database: dbQueries,
	}
	mux := http.NewServeMux()

	server := &http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	fileServerHandler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(fileServerHandler))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerResetCounter)

	server.ListenAndServe()
}
