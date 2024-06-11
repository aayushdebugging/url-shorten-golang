# url-shorten-golang
# URL Shortener Service

![Go](https://img.shields.io/badge/Go-1.18-blue.svg)
![License](https://img.shields.io/badge/license-MIT-green.svg)

A simple URL shortener service written in Go, which provides functionality to shorten URLs and redirect shortened URLs to their original destinations.

## Features

- **URL Shortening**: Converts long URLs into short, easy-to-share links.
- **Redirection**: Redirects short URLs to the original URLs.
- **Simple API**: RESTful endpoints to create short URLs and redirect them.
- **In-Memory Storage**: Uses an in-memory map to store URL mappings.

## Getting Started

### Prerequisites

Ensure you have Go installed on your machine. You can download it from [https://golang.org/dl/](https://golang.org/dl/).

### Installation

1. Clone the repository:

    ```bash
    git clone https://github.com/yourusername/url-shortener.git
    cd url-shortener
    ```

2. Build the project:

    ```bash
    go build -o url-shortener
    ```

3. Run the application:

    ```bash
    ./url-shortener
    ```

    The server will start on `http://localhost:8080`.

## Usage

### Endpoints

- **Root Page**

    ```
    GET /
    ```

    Responds with a simple "Hello World" message.

- **Shorten URL**

    ```
    POST /shorten
    ```

    Request Body:
    ```json
    {
        "url": "https://www.example.com"
    }
    ```

    Response Body:
    ```json
    {
        "short_url": "abc12345"
    }
    ```

- **Redirect URL**

    ```
    GET /redirect/{short_url}
    ```

    Redirects to the original URL associated with the `short_url`.

### Example Usage

1. **Shorten a URL**

    ```bash
    curl -X POST http://localhost:8080/shorten -H "Content-Type: application/json" -d '{"url":"https://www.example.com"}'
    ```

    This will return a shortened URL in the response.

2. **Redirect to Original URL**

    Navigate to `http://localhost:8080/redirect/{short_url}` in your browser or use a tool like curl.

    ```bash
    curl -L http://localhost:8080/redirect/{short_url}
    ```

## Code Overview

### Main Components

- **`main` function**: Sets up the HTTP server and routes.

    ```go
    func main() {
        http.HandleFunc("/", RootPageURL)
        http.HandleFunc("/shorten", ShortURLHandler)
        http.HandleFunc("/redirect/", redirectURLHandler)

        fmt.Println("Starting server on PORT 8080....")
        error := http.ListenAndServe(":8080", nil)
        if error != nil {
            fmt.Println("Error on starting the Server:", error)
        }
    }
    ```

- **`generateShortURL` function**: Generates a short URL from the original URL using MD5 hashing.

    ```go
    func generateShortURL(OriginalURL string) string {
        hasher := md5.New()
        hasher.Write([]byte(OriginalURL))
        data := hasher.Sum(nil)
        hash := hex.EncodeToString(data)
        return hash[:8]
    }
    ```

- **`createURL` function**: Creates a new URL record and stores it in the in-memory database.

    ```go
    func createURL(OriginalURL string) string {
        shortURL := generateShortURL(OriginalURL)
        newURL := URL{
            ID:           shortURL,
            OriginalURL:  OriginalURL,
            ShortURL:     shortURL,
            CreationDate: time.Now(),
        }
        urlDB[newURL.ID] = newURL
        return newURL.ShortURL
    }
    ```

- **`getURL` function**: Retrieves the original URL from the short URL.

    ```go
    func getURL(id string) (URL, error) {
        url, ok := urlDB[id]
        if !ok {
            return URL{}, errors.New("URL not Found")
        }
        return url, nil
    }
    ```

- **`ShortURLHandler` function**: Handles the creation of short URLs via POST requests.

    ```go
    func ShortURLHandler(w http.ResponseWriter , r *http.Request){
        var data struct{
            URL string `json:"url"`
        }
        err :=json.NewDecoder(r.Body).Decode(&data)
        if err!= nil{
            http.Error(w,"Invalid request Body",http.StatusBadRequest)
            return 
        }
        shortURL := createURL(data.URL)
        response := struct{
            ShortURL string `json:"short_url"`
        }{ShortURL: shortURL}

        w.Header().Set("Content-Type","application/json")
        json.NewEncoder(w).Encode(response)
    }
    ```

- **`redirectURLHandler` function**: Redirects the short URL to the original URL.

    ```go
    func redirectURLHandler(w http.ResponseWriter ,r *http.Request){
        id := r.URL.Path[len("/redirect/"):]
        url, err := getURL(id)
        if err != nil {
            http.Error(w,"Invalid Response",http.StatusNotFound)
            return
        }
        http.Redirect(w, r, url.OriginalURL, http.StatusFound)
    }
    ```

### Project Directory Structure

- `main.go`: Contains the main application logic.
- `README.md`: Documentation for the project.
- `.gitignore`: Specifies which files and directories to ignore in Git.

### Future Improvements

- **Persistent Storage**: Implement a database to store URLs permanently.
- **Authentication**: Add user authentication for creating and managing short URLs.
- **Custom URL Aliases**: Allow users to choose custom aliases for their URLs.
- **Analytics**: Track and display statistics for URL usage.

## License

This project is licensed under the MIT License.



