package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

type HandlePathData struct {
	Path         string
	TemplateFile string
}

const (
	UPLOADPATH = "/new/"
	EDITPATH   = "/edit/"
	VIEWPATH   = "/view/"
	IMAGEPATH  = "/img/"
)

var templates = template.Must(template.ParseFiles(
	"templates/new.html",
	"templates/edit.html",
	"templates/view.html",
))

// BUG(spencer): Users should not be able to POST edits to other's photos.

// newHandler serves the page where users upload their initial photo. It then
// redirects them to the editing page where they can finish their creation.
func newHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		err := templates.ExecuteTemplate(w, "new.html", "fake_token")
		if err != nil {
			fmt.Fprintf(w, "An error occurred. TODO: 500 this")
			log.Print("Error: ", err)
			return
		}
	} else if r.Method == "POST" {
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			log.Println("Could not get form file:", err)
			return
		}
		defer file.Close()

		// TODO make sure that it's an image (handler.Header)
		_ = handler
		rawBytes, err := ioutil.ReadAll(file)
		if err != nil {
			log.Println("Could not read the image:", err)
			return
		}

		// Crop it
		bytes, err := cropBytesToBytes(rawBytes)
		if err != nil {
			log.Println("Could not crop the image:", err)
			return
		}

		id, err := uploadImage(bytes)
		if err != nil {
			log.Println("Could not upload image:", err)
			return
		}

		http.Redirect(w, r, EDITPATH+id.Hex(), http.StatusFound)
	}
}

// editHandler serves the page on which users write code to filter the image they
// previously uploaded. It allows them to submit code matching their image.
func editHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		idStr := r.URL.Path[len(EDITPATH):]
		err := templates.ExecuteTemplate(w, "edit.html", idStr)
		if err != nil {
			fmt.Fprintf(w, "An error occurred. TODO: 500 this")
			log.Print("Error: ", err)
			return
		}
	} else if r.Method == "POST" {
		id := r.FormValue("id")
		code := r.FormValue("code")
		err := updateCode(id, code)
		if err != nil {
			log.Println("Could not update code:", err)
			return
		}

		http.Redirect(w, r, VIEWPATH+id, http.StatusFound)
	}
}

// viewHandler displays a finished image+code combo.
func viewHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len(VIEWPATH):]
	post, err := getPost(idStr)
	if err != nil {
		fmt.Fprintf(w, "An error occurred. TODO: 500 this")
		log.Print("Error: ", err)
		return
	}
	err = templates.ExecuteTemplate(w, "view.html", post)
	if err != nil {
		fmt.Fprintf(w, "An error occurred. TODO: 500 this")
		log.Print("Error: ", err)
		return
	}
}

func imageHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len(IMAGEPATH):]
	img, err := getImage(idStr)
	if err != nil {
		log.Println("Could not fetch image:", err)
		return
	}
	w.Write(img)
}

func main() {
	// Register handlers
	http.HandleFunc(UPLOADPATH, newHandler)
	http.HandleFunc(EDITPATH, editHandler)
	http.HandleFunc(VIEWPATH, viewHandler)
	http.HandleFunc(IMAGEPATH, imageHandler)

	// Serve forever
	log.Fatal(http.ListenAndServe(":8000", nil))
}
