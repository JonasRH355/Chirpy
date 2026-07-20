package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"github.com/JonasRH355/Chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/google/uuid"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries      database.Queries
	platform       string
}

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uuid.UUID `json:"user_id"`
	Body      string    `json:"body"`
}

func (c *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.fileserverHits.Store(c.fileserverHits.Add(1))
		next.ServeHTTP(w, r)
	})
}

func (c *apiConfig) handlerHits(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("charset", "utf-8")
	w.WriteHeader(http.StatusOK)

	str := fmt.Sprintf(`<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, c.fileserverHits.Load())

	w.Write([]byte((str)))
}

func (c *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if c.platform != "dev" {
		w.WriteHeader(403)
		w.Write([]byte("403 Forbidden"))
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("charset", "utf-8")
	w.WriteHeader(http.StatusOK)

	c.dbQueries.DeleteUsers(r.Context())
	c.fileserverHits.Store(0)
	str := fmt.Sprintf("Hits: %v\n", c.fileserverHits.Load())
	w.Write([]byte((str)))
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("charset", "utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(" OK"))
}

func removeBadWords(sentence string) string {
	words := strings.Split(sentence, " ")
	var newSentence []string

	for _, word := range words {
		switch strings.ToLower(word) {
		case "kerfuffle":
			newSentence = append(newSentence, "****")
		case "sharbert":
			newSentence = append(newSentence, "****")
		case "fornax":
			newSentence = append(newSentence, "****")
		default:
			newSentence = append(newSentence, word)
		}
	}

	return strings.Join(newSentence, " ")
}

func (c *apiConfig) chirpHandler(w http.ResponseWriter, r *http.Request) {
	sttsCode := 200
	type parameters struct {
		Id uuid.UUID `json:"id"`
		Body  string `json:"body"`
		UserId uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		sttsCode = 500
		return
	}

	fmt.Println("Decoded request")

	newChirp := database.CreateChirpParams{}
	newChirp.UserID.UUID = params.UserId
	newChirp.UserID.Valid = true

	newChirp.Body = removeBadWords(params.Body)

	if len(newChirp.Body) <= 140 {
		sttsCode = 201
	} else {
		sttsCode = 400
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(sttsCode)
	}
	
	fmt.Printf("%v Removed bad words:\nBefore: %v\nAfter: %v\n",sttsCode,params.Body,newChirp.Body)

	resBody,err := c.dbQueries.CreateChirp(r.Context(),newChirp)
	if err != nil {
		sttsCode = 400
		w.Header().Set("Content-Type", "application/json")
		log.Println(err)
		w.WriteHeader(sttsCode)
		return
	}

	fmt.Printf("Created new chirp: %v \n", resBody)


	dat, err := json.Marshal(parameters{
		Id: resBody.ID,
		Body: resBody.Body,
		UserId: resBody.UserID.UUID,
	})
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		sttsCode = 400
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(sttsCode)
		return
	}

	w.WriteHeader(sttsCode)
	w.Header().Set("Content-Type", "application/json")
	w.Write(dat)
	
}

func (c *apiConfig) addUser(w http.ResponseWriter, r *http.Request) {
	type User struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	reqBody := User{}
	err := decoder.Decode(&reqBody)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}

	newUser, err := c.dbQueries.CreateUser(r.Context(), reqBody.Email)
	if err != nil {
		log.Printf("Error on creating user: %s", err)
		w.WriteHeader(400)
		return
	}

	dat, err := json.Marshal(User{
		ID: newUser.ID,
		CreatedAt: newUser.CreatedAt,
		UpdatedAt: newUser.UpdatedAt,
		Email: newUser.Email,
	})
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(201)
	w.Header().Set("Content-Type", "application/json")
	w.Write(dat)
}

func (c *apiConfig) getChirps(w http.ResponseWriter, r *http.Request){
	chirps, err := c.dbQueries.GetChirps(r.Context())
	if err != nil {
		w.WriteHeader(500)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("Error on DB"))
	}

	respBody := []Chirp{}
	for _,i := range chirps {
		respBody = append(respBody, Chirp{
			ID: i.ID,
			CreatedAt: i.CreatedAt,
			UpdatedAt: i.UpdatedAt,
			UserID: i.UserID.UUID,
			Body: i.Body,
		})
	}


	response, err := json.Marshal(respBody)
	if err != nil {
		w.WriteHeader(500)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("Error on marshal data"))
	}

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func (c *apiConfig) getChirp(w http.ResponseWriter, r *http.Request) {
	path := r.PathValue("chirpID")
	uuidVal, err := uuid.Parse(path)
	if err != nil {
		log.Println(err)
		log.Println()
		log.Println(path)
		w.WriteHeader(500)
		return
	}

	respBody, err := c.dbQueries.GetChirp(r.Context(),uuidVal)
	if err != nil {
		if fmt.Sprint(err) == "sql: no rows in result set"{
			w.WriteHeader(404)
			return
		}
		log.Println(err)
		
		w.WriteHeader(500)
		return 
	}

	if respBody.ID.String() == " " {
		w.WriteHeader(404)
		return	
	}

	response, err := json.Marshal(Chirp{
		ID: respBody.ID,
		CreatedAt: respBody.CreatedAt,
		UpdatedAt: respBody.UpdatedAt,
		UserID: respBody.UserID.UUID,
		Body: respBody.Body,
	})
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
	}

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	dbPLATFORM := os.Getenv("PLATFORM")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	dbQueries := database.New(db)

	const filepathRoot = "."
	const port = "8080"
	apiCfg := apiConfig{
		dbQueries: *dbQueries,
		platform:  dbPLATFORM,
	}
	apiCfg.fileserverHits.Store(0)

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir(filepathRoot))

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", fileServer)))
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerHits)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	mux.HandleFunc("GET /api/healthz", handler)
	mux.HandleFunc("POST /api/chirps", apiCfg.chirpHandler)
	mux.HandleFunc("GET /api/chirps", apiCfg.getChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.getChirp)
	mux.HandleFunc("POST /api/users", apiCfg.addUser)

	serv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(serv.ListenAndServe())

}
