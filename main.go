package main 


// REQUIREMENTS:
// auth token matches --> OK
// POST data to the /upload --> OK
// should write the received file data to a temporary file -->  OK
// content type of the uploaded file is an image --> OK
// Images larger than 8 megabytes should also be rejected --> OK
// If the submission is bad, please return a 403 HTTP error code --> OK
// Write the image metadata (content type, size, etc) to a database of your choice, including all relevant HTTP information. --> OK


// imports here
// nosurf : library for auth token generation & validation 

import (
	"fmt"
	"github.com/justinas/nosurf"
	"html/template"
	"net/http"
        "io/ioutil"
	"strings"
	"database/sql"
      _ "github.com/go-sql-driver/mysql"
         "time"
	 "strconv"
)

//mysl dtbse configurtion here
func dbConn() (db *sql.DB) {
    dbDriver := "mysql"
    dbUser := "teguh"
    dbPass := "sys.admin3"
    dbName := "go_db"
    db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
    if err != nil {
        panic(err.Error())
    }
    return db
}


// HTML template
var templateString string = `
<!doctype html>
<html>
<body>
<form action="/upload" method="POST" enctype="multipart/form-data">
<input type="file" name="data">
<input type="hidden" name="csrf_token" value="{{ .token }}">
<input type="submit" value="Send">
</form>
</body>
</html>
`
var templ = template.Must(template.New("t1").Parse(templateString))

func myFunc(w http.ResponseWriter, r *http.Request) {

	context := make(map[string]string)
	context["token"] = nosurf.Token(r)
	if r.Method == "POST" {

	r.ParseMultipartForm(8 << 20)

	// Get handler for filename, size and headers
	file, handler, err := r.FormFile("data")
	if err != nil {
	    w.WriteHeader(http.StatusForbidden)
            w.Write([]byte("403 - Status Forbidden!"))
	    return
	}

	defer file.Close()

	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)
        contentType := handler.Header.Get("Content-type")
	fmt.Printf("contentType: %+v\n", contentType)

        // check if uploaded file is image
        res := strings.Contains(contentType, "image/")
        if (res != true) {
	    w.WriteHeader(http.StatusForbidden)
            w.Write([]byte("403 - Status Forbidden!"))
	    return
        }

        const MAX_UPLOAD_SIZE = 1024 * 1024 * 8
        // filesize blocking (max 8MB)
        r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)
	if err := r.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil {
	    w.WriteHeader(http.StatusForbidden)
            w.Write([]byte("403 - Status Forbidden!"))
	    return
	}


    // Create a temporary file within our temp-images directory that follows a particular naming pattern

    t := time.Now()
    t.String()
    tUnixMicro := int64(time.Nanosecond) * t.UnixNano() / int64(time.Microsecond)
    str_tUnixMicro := strconv.FormatInt(tUnixMicro, 10)
    var fileNm = "image-" + str_tUnixMicro + ".png"

    tempFile, err := ioutil.TempFile("temp-images", "image-*.png")
    if err != nil {
            w.WriteHeader(http.StatusForbidden)
            w.Write([]byte("403 - Status Forbidden!"))
	    return
    }
    defer tempFile.Close()

    // read all of the contents of our uploaded file into a byte array
    fileBytes, err := ioutil.ReadAll(file)
    if err != nil {
            w.WriteHeader(http.StatusForbidden)
            w.Write([]byte("403 - Status Forbidden!"))
	    return
    }
    // write this byte array to our temporary file
    tempFile.Write(fileBytes)


    // return that we have successfully uploaded our file!
    fmt.Fprintf(w, "Successfully Uploaded File\n")



              // Write the image metadata (content type, size, etc) to a database, including all relevant HTTP information. 
              db := dbConn()
              insForm, err := db.Prepare("INSERT INTO tbl_image (real_filename, saved_filename, content_type, file_size) VALUES(?,?,?,?)")
              if err != nil {
                  panic(err.Error())
	          w.WriteHeader(http.StatusForbidden)
                  w.Write([]byte("403 - Status Forbidden!"))
	          return
              }
              insForm.Exec(handler.Filename, fileNm, contentType, handler.Size)
	}
	
	templ.Execute(w, context)

}


func main() {
	myHandler := http.HandlerFunc(myFunc)
	fmt.Println("Listening on http://127.0.0.1:8000/")
	http.ListenAndServe(":8000", nosurf.New(myHandler))
}
