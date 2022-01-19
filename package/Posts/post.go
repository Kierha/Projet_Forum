package Post

import (
	"database/sql"
	"fmt"
	"html/template"
	"main/package/Structures"
	"main/package/Utils"
	"net/http"
)

func AllPosts(w http.ResponseWriter, r *http.Request, database *sql.DB) {
	tmpl := template.Must(template.ParseFiles("./templates/displayPost.html"))
	var allPost []Structures.ONEPOST
	var singlePost Structures.ONEPOST
	rows, _ := database.Query("SELECT id, title, description, content, category, image FROM posts")
	for rows.Next() {
		rows.Scan(&singlePost.ID, &singlePost.Title, &singlePost.Description, &singlePost.Content, &singlePost.Category, &singlePost.Image)
		allPost = append(allPost, singlePost)
	}
	tmpl.Execute(w, allPost)
}

func AddPost(w http.ResponseWriter, r *http.Request, database *sql.DB) {
	var post Structures.ONEPOST
	post.Title = r.FormValue("title")
	post.Description = r.FormValue("description")
	post.Content = r.FormValue("content")
	post.Category = r.FormValue("category")
	post.Image = r.FormValue("image")
	statement, _ := database.Prepare("INSERT INTO posts (title, description, content, category, image) VALUES (?, ?, ?, ?, ?)")
	statement.Exec(post.Title, post.Description, post.Content, post.Category, post.Image)
	http.Redirect(w, r, "/allPosts", http.StatusMovedPermanently)
}

func DeleteOnePost(w http.ResponseWriter, r *http.Request, database *sql.DB) {
	fmt.Println("--- DELETE PROCESS : RUNNING ---")
	r.ParseForm()
	id := r.FormValue("idPost")
	del, err := database.Prepare("DELETE FROM `posts` WHERE (`id` = ? );")
	if err != nil {
		panic(err)
	}
	defer del.Close()
	var res sql.Result
	res, err = del.Exec(id)
	rowsAff, _ := res.RowsAffected()
	if err != nil || rowsAff != 1 {
		fmt.Fprint(w, "Error deleting product")
		return
	}
	fmt.Println("err:", err)
	fmt.Println("--- DELETE PROCESS : SUCCESS ---")
	http.Redirect(w, r, "/allPosts", http.StatusMovedPermanently)
}

func EditPost(w http.ResponseWriter, r *http.Request, database *sql.DB) {
	fmt.Println("--- EDITING PROCESS : RUNNING ---")
	r.ParseForm()
	id := r.FormValue("idPost")
	row := database.QueryRow("SELECT * FROM `posts` WHERE Id = ?", id)
	var p Structures.ONEPOST
	err := row.Scan(&p.ID, &p.Title, &p.Description, &p.Content, &p.Category, &p.Image)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/allPosts", http.StatusTemporaryRedirect)
		fmt.Fprint(w, "Error editing post")
	}

	Utils.Tmpl.ExecuteTemplate(w, "modifyPost.html", p)
}

func NewPostEdited(w http.ResponseWriter, r *http.Request, database *sql.DB) {
	tmpl := template.Must(template.ParseFiles("./templates/modifyPost.html"))
	fmt.Println("--- REPLACE IN POST TABLE : RUNNING ---")
	r.ParseForm()
	id := r.FormValue("idPost")
	title := r.FormValue("title")
	description := r.FormValue("description")
	content := r.FormValue("content")
	updStatement := "UPDATE `posts` SET `title` = ?, `description` = ?, `content` = ? WHERE (`Id`= ?);"
	statement, err := database.Prepare(updStatement)
	if err != nil {
		fmt.Println("Statement failed")
		panic(err)
	}
	defer statement.Close()
	var res sql.Result
	res, err = statement.Exec(title, description, content, id)
	rowsAffect, _ := res.RowsAffected()
	if err != nil || rowsAffect != 1 {
		fmt.Println(err)
		tmpl.ExecuteTemplate(w, "/", "Update failed")
		return
	}
	fmt.Println("--- REPLACE IN POST TABLE : SUCCESS ---")
	http.Redirect(w, r, "/allPosts", http.StatusMovedPermanently)
}
