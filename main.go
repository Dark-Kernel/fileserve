package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"golang.org/x/net/html"
    "io"
    "os"
    "path"
    "math/rand"
)
     

func randomFileName() string {
    // Generate a random string of length 4 using the letters a-z and A-Z.
    // https://stackoverflow.com/a/22892986

    letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
    b := make([]rune, 4)
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))]
    }
    return string(b)
}

func renameFile(filePath string) string {
    // fileName := strings.Split(filePath, ".")[0]
    filePathDir := strings.Split(filePath, "/")[0]
    extension := strings.Split(filePath, ".")[1]
    newFilePath := filePathDir + "/" + randomFileName() + "." + extension
    os.Rename(filePath, newFilePath)
    return newFilePath
}


func main() {
    var name string
    fmt.Printf("Enter the number: ")
    fmt.Scan(&name)
    fmt.Printf("hello %v \nServer is runninggggg......\n", name)


    resp, err := http.Get("https://google.com")
    if err != nil {
        fmt.Println("Some error occurred: ", err)
    }

    defer resp.Body.Close()
    

    http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
    })

    http.HandleFunc("/cmd", func(w http.ResponseWriter, r *http.Request) {
        inp := html.EscapeString(r.FormValue("inp"))
        fmt.Println(strings.Split(inp, " "))
        spt := strings.Split(inp, " ")
        cmd := exec.Command(spt[0], spt[1])
        cmd.Dir = "../"
        var out bytes.Buffer
        cmd.Stdout = &out 
        err := cmd.Run()
        if err != nil {
            log.Fatal("Error occured: ",err)
        }
        fmt.Fprint(w, "Output : \n", out.String())
    })

    http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
        r.ParseMultipartForm(10 << 20)
        file, handler, err := r.FormFile("file")
        if err != nil {
            w.WriteHeader(http.StatusBadRequest)
            fmt.Println("Error retrieving file from form-data", err)
            return
        }
        defer file.Close()
        fileBytes, err := io.ReadAll(file)
        if err != nil {
            w.WriteHeader(http.StatusInternalServerError)
            fmt.Println("Error reading file", err)
            return
        }
        newFilePath := path.Join("files", handler.Filename)
        err = os.WriteFile(newFilePath, fileBytes, 0644)

        renameFilePath := renameFile(newFilePath)
        
        if err != nil {
            w.WriteHeader(http.StatusInternalServerError)
            fmt.Println("Error writing file to disk", err)
            return
        }
        w.WriteHeader(http.StatusOK)
        fmt.Fprintf(w, "File uploaded", renameFilePath)
        fmt.Println("File uploaded successfully")
    })

    
    log.Fatal(http.ListenAndServe(":8080", nil))

}
