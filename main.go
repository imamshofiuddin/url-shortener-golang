package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

type URLShortener struct {
	urls map[string]string
}

func main() {
	shortener := &URLShortener{
		urls: make(map[string]string),
	}

	http.HandleFunc("/", shortener.handler)
	http.HandleFunc("/shorten", shortener.HandleShorten)

	fmt.Println("Server is running on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func (us *URLShortener) handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(us.urls)
	var html string = ""
	if r.URL.Path[len("/"):] == "" {
		w.Header().Set("Content-Type", "text/html")
		html = `<h2>URL Shortener</h2>
        <form method="post" action="/shorten">
            <input type="text" name="url" size=100 placeholder="Enter a URL"></br>
            <input type="submit" value="Shorten">
        </form>
		`
		fmt.Fprintf(w, html)
	} else {
		originalURL, found := us.urls[r.URL.Path[len("/"):]]
		if !found {
			http.Error(w, "Shortened key not found", http.StatusNotFound)
			return
		}
		http.Redirect(w, r, originalURL, http.StatusMovedPermanently)
	}
}

func (us *URLShortener) HandleShorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	originalURL := r.FormValue("url")
	if originalURL == "" {
		http.Error(w, "URL parameter is missing", http.StatusBadRequest)
		return
	}

	shortKey := generateShortKey()
	us.urls[shortKey] = originalURL

	shortenedURL := fmt.Sprintf("http://localhost:8080/%s", shortKey)

	w.Header().Set("Content-Type", "text/html")
	responseHTML := fmt.Sprintf(`
        <h2>URL Shortener</h2>
        <p>Original URL: %s</p>
        <p>Shortened URL: <a href="%s">%s</a></p>
    `, originalURL, shortenedURL, shortenedURL)
	fmt.Fprintf(w, responseHTML)
}

func generateShortKey() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const keyLength = 6

	rand.Seed(time.Now().UnixNano())
	shortKey := make([]byte, keyLength)
	for i := range shortKey {
		shortKey[i] = charset[rand.Intn(len(charset))]
	}
	return string(shortKey)
}
