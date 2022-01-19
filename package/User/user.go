package User

import (
	"database/sql"
	"fmt"
	"main/package/Utils"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func InsertUserInfos(w http.ResponseWriter, r *http.Request, database *sql.DB) {
	fmt.Println("--- INFOINSERTION PROCESS : RUNNING ---")
	r.ParseForm()
	username := r.FormValue("register-username")

	statement := "SELECT Username FROM users WHERE username = ?"
	row := database.QueryRow(statement, username)
	var existingUsername string
	err := row.Scan(&existingUsername)
	if err != sql.ErrNoRows {
		fmt.Println("ERROR, This username is already taken, err : ", err)
		Utils.Tmpl.ExecuteTemplate(w, "index.html", "Username already exist")
		return
	}
	password := r.FormValue("register-password")

	// -------------- HASHED PASSWORD --------------
	var hashedPassword []byte
	hashedPassword, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("ERROR, encryption impossible err: ", err)
		Utils.Tmpl.ExecuteTemplate(w, "/", "Register failed")
	}
	addStatement, err := database.Prepare("INSERT INTO users (username, password) VALUES (?, ?);")
	if err != nil {
		fmt.Println("ERROR, Statement Failed err: ", err)
		Utils.Tmpl.ExecuteTemplate(w, "/", "Register failed")
	}
	defer addStatement.Close()
	var params sql.Result
	params, err = addStatement.Exec(username, hashedPassword)
	fmt.Println("--- INFOINSERTION PROCESS : SUCCESS ---")
	rowsAffected, _ := params.RowsAffected()
	lastInsertId, _ := params.LastInsertId()
	fmt.Println("rowsAffected:", rowsAffected)
	fmt.Println("lastIns:", lastInsertId)
	fmt.Println("err:", err)
	if err != nil {
		fmt.Println("ERROR, Insertion failed err: ", err)
		Utils.Tmpl.ExecuteTemplate(w, "/", "Impossible to register")
	}
	fmt.Println("--- REGISTER PROCESS : SUCCESS ---")
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func LoginUser(w http.ResponseWriter, r *http.Request, database *sql.DB) {
	fmt.Println("--- LOGIN PROCESS : RUNNING ---")
	r.ParseForm()
	username := r.FormValue("connexion-username")
	password := r.FormValue("connexion-password")
	fmt.Println("username: ", username, "password: ", password)
	var hashedPassword string
	statement := "SELECT password FROM users WHERE username = ?"
	row := database.QueryRow(statement, username)
	err := row.Scan(&hashedPassword)
	fmt.Println("hash from database: ", hashedPassword)
	if err != nil {
		fmt.Println("ERROR, Impossible to select Hashed Password from users database")
		Utils.Tmpl.ExecuteTemplate(w, "index.html", "Verify root to username and password")
	}
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err == nil {
		fmt.Println("--- LOGIN PROCESS : SUCCESS ---")
		Utils.Tmpl.ExecuteTemplate(w, "loggedIn.html", "Logged in")
		return
	}
	fmt.Println("--- LOGIN PROCESS : ERROR ---")
	Utils.Tmpl.ExecuteTemplate(w, "index.html", "Verify username and/or password")
}

func LoggedInHandler(w http.ResponseWriter, r *http.Request, database *sql.DB) {
	fmt.Println("--- LOGIN HANDLING : RUNNING ---")
	session, _ := Utils.Store.Get(r, "session")
	_, ok := session.Values["username"]
	fmt.Println("ok: ", ok)
	if !ok {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	fmt.Println("--- LOGIN HANDLING : WORKING ---")
	Utils.Tmpl.ExecuteTemplate(w, "loggedIn.html", "Connected")
}
