package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	Post "main/package/Posts"
	"main/package/User"
	"main/package/Utils"
	"net/http"

	"github.com/gorilla/context"
	_ "github.com/mattn/go-sqlite3"
)

func indexpage(w http.ResponseWriter, r *http.Request) {

	tmpl := template.Must(template.ParseFiles("./templates/index.html"))
	if r.Method != http.MethodPost {
		tmpl.Execute(w, nil)
		return
	}
}

func main() {
	// PARSING TEMPLATES --> TO VARIABLE TMPL
	Utils.Tmpl, _ = template.ParseGlob("templates/*.html")
	// -------------- DATABASE / TABLE CREATION --------------
	database, _ := sql.Open("sqlite3", "./database/forum.db")
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS posts (id INTEGER PRIMARY KEY, title TEXT, description TEXT, content TEXT, category TEXT, image TEXT)")
	statement.Exec()
	statement, _ = database.Prepare("CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, username TEXT, password TEXT, mail TEXT)")
	statement.Exec()

	//-------------- ROOT ----------------
	http.HandleFunc("/", indexpage)

	// ---------------- POST ----------------
	http.HandleFunc("/allPosts", func(w http.ResponseWriter, r *http.Request) {
		Post.AllPosts(w, r, database)
	})
	http.HandleFunc("/addPost", func(w http.ResponseWriter, r *http.Request) {
		Post.AddPost(w, r, database)
	})
	http.HandleFunc("/delete", func(w http.ResponseWriter, r *http.Request) {
		Post.DeleteOnePost(w, r, database)
	})
	http.HandleFunc("/insertEditionPost/", func(w http.ResponseWriter, r *http.Request) {
		Post.NewPostEdited(w, r, database)
	})
	http.HandleFunc("/edit/", func(w http.ResponseWriter, r *http.Request) {
		Post.EditPost(w, r, database)
	})
	// ---------------- USER ----------------
	http.HandleFunc("/insertUserInfos", func(w http.ResponseWriter, r *http.Request) {
		User.InsertUserInfos(w, r, database)
	})
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		User.LoginUser(w, r, database)
	})

	//-------------- Integration CSS + IMG + JS --------------
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("img"))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("js"))))

	//-------------- SERVEUR ----------------
	fmt.Printf("Démarage du serveur Go sur le port 8080 --> Projet Forum")
	if err := http.ListenAndServe(":8080", context.ClearHandler(http.DefaultServeMux)); err != nil {
		log.Fatal(err)
	}
}
