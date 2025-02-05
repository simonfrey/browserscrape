package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"net/http"
	"net/url"
	"os"
	"server/pkg/html_to_md"
	"sync"
)

func main() {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Enable CORS
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			next.ServeHTTP(w, r)
		})
	})

	// SSE endpoint
	r.Get("/events/{uuid}", sseHandler)
	r.Post("/response/{response_uuid}", responseHandler)
	r.Get("/scrape/html/{uuid}/*", scrapeHandler(false))
	r.Get("/scrape/md/{uuid}/*", scrapeHandler(true))

	port := os.Getenv("PORT")
	fmt.Println("Server starting on :", port)
	http.ListenAndServe(":"+port, r)
}

var activeClients sync.Map
var responseQueue sync.Map

func init() {
	activeClients = sync.Map{}
	responseQueue = sync.Map{}
}

func scrapeHandler(convertToMarkdown bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract UUID from the URL
		reqUUID := chi.URLParam(r, "uuid")
		fmt.Println("Received scrape html request for UUID: ", reqUUID)

		// Check if client for uuid exists
		queue, ok := activeClients.Load(reqUUID)
		if !ok {
			http.Error(w, fmt.Sprintf("Client for uuid %q found", reqUUID), http.StatusNotFound)
			return
		}

		scrapeURL := chi.URLParam(r, "*")

		// Test if scrapeURL is a valid URL, by parsing it
		u, err := url.ParseRequestURI(scrapeURL)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid URL %q: %s", scrapeURL, err.Error()), http.StatusBadRequest)
		}
		if u.Scheme == "" {
			u.Scheme = "http"
		}
		scrapeURL = u.String()
		fmt.Println("Scraping URL: ", scrapeURL)

		responseUUID := uuid.New().String()
		fmt.Println("Response UUID: ", responseUUID)

		// Setup response queue
		responseChan := make(chan responseFormat)
		responseQueue.Store(responseUUID, responseChan)
		defer responseQueue.Delete(responseUUID)

		// Async wait for the response chan and write the byte to w until channel is closed
		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case response, ok := <-responseChan:
					if !ok {
						return
					}
					if convertToMarkdown {
						// Convert HTML to Markdown
						markdownResponse, err := html_to_md.HTMLToMD(response.Url, response.Html, true)
						if err != nil {
							http.Error(w, fmt.Sprintf("Error converting HTML to Markdown: %v", err), http.StatusInternalServerError)
							return
						}
						_, _ = w.Write([]byte(markdownResponse))
						return
					}

					_, _ = w.Write([]byte(response.Html))
					return
				}
			}
		}()

		// Send event to client
		queue.(chan SSEEvent) <- SSEEvent{
			Url:          scrapeURL,
			ResponseUuid: responseUUID,
		}

		wg.Wait()
	}
}

type responseFormat struct {
	Url   string `json:"url"`
	Title string `json:"title"`
	Html  string `json:"html"`
}

func responseHandler(w http.ResponseWriter, r *http.Request) {
	// Extract UUID from the URL
	responseUuid := chi.URLParam(r, "response_uuid")
	fmt.Printf("Received request for UUID: %s\n", responseUuid)

	// Check if responseUuid is in the responseQueue
	responseC, ok := responseQueue.Load(responseUuid)
	if !ok {
		http.Error(w, fmt.Sprintf("Response UUID %q not found", responseUuid), http.StatusNotFound)
		return
	}
	responseChan := responseC.(chan responseFormat)

	// Parse body into responseFormat
	resp := responseFormat{}
	err := json.NewDecoder(r.Body).Decode(&resp)
	if err != nil {
		http.Error(w, "Error parsing response", http.StatusBadRequest)
		return
	}

	fmt.Printf("Response: %+v\n", resp)

	// Write response to responseChan
	responseChan <- resp
	// Close responseChan
	close(responseChan)

}

type SSEEvent struct {
	Url          string `json:"url"`
	ResponseUuid string `json:"response_uuid"`
}

func sseHandler(w http.ResponseWriter, r *http.Request) {
	// Extract UUID from the URL
	reqUuuid := chi.URLParam(r, "uuid")
	fmt.Printf("Received request for UUID: %s\n", reqUuuid)

	// Set headers for SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Create channel to detect client disconnection
	clientGone := r.Context().Done()

	// Create flusher
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	// Add active client to map
	queue := make(chan SSEEvent)
	activeClients.Store(reqUuuid, queue)
	defer activeClients.Delete(reqUuuid)

	// Infinite loop to send events
	for {
		select {
		case <-clientGone:
			// Client disconnected
			fmt.Println("Client disconnected")
			return
		case request := <-queue:
			jsonEvent, err := json.Marshal(request)
			if err != nil {
				http.Error(w, "Error marshalling event", http.StatusInternalServerError)
				return
			}
			_, err = fmt.Fprintf(w, "data: %s\n\n", string(jsonEvent))
			if err != nil {
				fmt.Println("ERROR Fprintln", err)
				return
			}
			// Flush the data to the client
			flusher.Flush()
			fmt.Println("Send", string(jsonEvent))
		}
	}
}
