package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
   "github.com/joho/godotenv"
   "os"
   "slices"
) 

// Response defines the structure for API responses
type Response struct {
	Status  int    `json:"status"`  // Status code
	Message string `json:"message"` // Response message
}

func main() {
	http.HandleFunc("/text", textHandler)

	log.Println("Server started on port :3000")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatal("Server error:", err)
	}
}

// jsonResponse sends a JSON response with status code in the body and header
func jsonResponse(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status) // Set the HTTP status code in the response header

	// Send JSON response body including the status code and message
	json.NewEncoder(w).Encode(Response{
		Status:  status,
		Message: message,
	})
}

func textHandler(w http.ResponseWriter, r *http.Request) {
	// Check Content-Type
	contentType := r.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		jsonResponse(w, http.StatusUnsupportedMediaType, "Content-Type must be application/json")
		return
	}

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, "Error reading request body")
		return
	}
	defer r.Body.Close()

	// Unmarshal the JSON body into a struct
	var requestBody Response
	if err := json.Unmarshal(body, &requestBody); err != nil {
		jsonResponse(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}
    
   if requestBody.Message == ""  {
     jsonResponse(w, http.StatusBadRequest, "message missing from the body")
     return
   }

	// Log the received message
	log.Println("Received Message:", requestBody.Message)

   if len(requestBody.Message) > 1000 {
      jsonResponse(w, http.StatusRequestEntityTooLarge , "Due to temporary Cloudflare limits, a message can only be upto 1000 characters.")
      return
   }

   
   // Create a slice of the message
   wordChunks := strings.Fields(requestBody.Message)
   
   // Store the semantic chunks in this variable


   semanticChunks := createSemanticChunks(wordChunks) 

   for _, chunk := range semanticChunks {
      log.Println(chunk)
   }
	// Respond with success
	jsonResponse(w, http.StatusOK, "Request Accepted")
}


func goDotEnvVariable(key string) string {

  // load .env file
  err := godotenv.Load(".env")

  if err != nil {
    log.Fatalf("Error loading .env file")
  }

  return os.Getenv(key)
}

func createSemanticChunks(chunks []string) []string{
   var totalLength int = 0
   var semanticChunk []string
   var result []string
   for i := 0 ; i < len(chunks) ; i++ {
      if totalLength < 25 {
         semanticChunk = append(semanticChunk, chunks[i])
         totalLength += len(chunks[i]) + 1 
      } else {
         var rl int = 0
         var rp []string

         for i, _ := range semanticChunk  {
            if rl < 8 {
               rl += len(semanticChunk[len(semanticChunk) - i - 1]) + 1
               rp = append(rp, semanticChunk[len(semanticChunk) - 1 - i])
               i = i - 1
            }
         }
         i = i - 1 
         a := strings.Join(semanticChunk, " ")
         result = append(result, a)
         semanticChunk = []string{}
         slices.Reverse(rp)
         semanticChunk = rp
         totalLength = rl
      }
   }
   a := strings.Join(semanticChunk, " ")
   result = append(result, a)

   return result
}
