# Router

Router package wraps a standard `http.ServeMux` to register method handler (GET / PUT / PATCH / POST / DELETE).

## Usage

```

func main() {
    r := router.New()
    r.GET("/foo", handleFoo)
    
    http.ListenAndServe(":8080", r)
}


func handleFoo(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "bar")
}
```