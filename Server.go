package main

import (
	"fmt"
	"net/http"
	"log"
	"html/template"
	"math/rand"
	"time"
	"github.com/skratchdot/open-golang/open"
	"strconv"
	"net"
	"github.com/ttacon/chalk"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"io/ioutil"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

func getUserName(r *http.Request) (username string) {
	if cookie, err := r.Cookie("session"); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
			username = cookieValue["username"]
		}
	}
	return username
}

func setSession(username string, w http.ResponseWriter) {
	value := map[string]string {
		"username": username,
	}
	if encoded, err := cookieHandler.Encode("session", value); err == nil {
		cookie := &http.Cookie{
			Name: "session",
			Value: encoded,
			Path: "/",
		}
		http.SetCookie(w, cookie)
	}
}

func clearSession(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name: "session",
		Value: "",
		Path: "/",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	redirectTarget := "/"
	fmt.Println("\nLogin Form: ", r.Form)
	if checkDBUser(username, password) {
		setSession(username, w)
		redirectTarget = "/internal"
		fmt.Println(chalk.Green, "\nUser successful logged in", chalk.Reset)
	}
	http.Redirect(w, r, redirectTarget, 302)
}

func signinHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")

	redirectTarget := "/"
	fmt.Println("\nSignIn Form: ", r.Form)

	b := writeDBUser(username, password, email)
	if b == true {
		fmt.Println(chalk.Green, "User successfull created", chalk.Reset)
	} else {
		fmt.Println(chalk.Red, "Error creating user", chalk.Reset)
	}
	http.Redirect(w, r, redirectTarget, 302)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	clearSession(w)
	fmt.Println(chalk.Red, "\nUser logged out", chalk.Reset)
	http.Redirect(w, r, "/", 302)
}

func tokenHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\nGenerate Token:")
	access_token := generateToken()
	refresh_token := generateToken()

	username := getUserName(r)
	if checkDBTokens(username) == false {
		setDBTokens(access_token, refresh_token, username)
	}
	clearSession(w)
	http.Redirect(w, r, "/", 302)
}

func indexPageHandler(w http.ResponseWriter, r *http.Request) {
	tmp, err := template.ParseFiles("index.html")
	if err != nil {
		log.Fatal("HTML Template Error: ", err)
	}
	tmp.Execute(w, nil)
}

func signinformHandler(w http.ResponseWriter, r *http.Request) {
	tmp, err := template.ParseFiles("signin.html")
	if err != nil {
		log.Fatal("HTML Template Error: ", err)
	}
	tmp.Execute(w, nil)
}

func internalPageHandler(w http.ResponseWriter, r *http.Request) {
	username := getUserName(r)

	if username != "" {
		file, err := ioutil.ReadFile("content.html") // just pass the file name
		if err != nil {
			fmt.Print(err)
		}
		internalPage := string(file) // convert content to a 'string'
		fmt.Fprintf(w, internalPage, username)
	}
	http.Redirect(w, r, "/", 302)
}

func checkDBErr(err error) {
	if err != nil {
		panic(err)
	}
}

func checkDBTokens(usr string) bool {
	db, err := sql.Open("mysql", "admin:mg545Iop@/goLoginApp")
	checkDBErr(err)


	accessstmt, err := db.Prepare("SELECT access_token FROM users WHERE username=?")
	checkDBErr(err)
	accessRow, err := accessstmt.Query(usr)
	checkDBErr(err)

	refreshstmt, err := db.Prepare("SELECT refresh_token FROM users WHERE username=?")
	checkDBErr(err)
	refreshRow, err := refreshstmt.Query(usr)
	checkDBErr(err)

	var access_token sql.NullString
	var refresh_token sql.NullString

	for accessRow.Next() {
		err = accessRow.Scan(&access_token)
		checkDBErr(err)
		fmt.Println("access_token: ", access_token)
	}
	for refreshRow.Next() {
		err = refreshRow.Scan(&refresh_token)
		checkDBErr(err)
		fmt.Println("refresh_token: ", refresh_token)
	}

	if access_token.Valid && refresh_token.Valid {
		db.Close()
		return true
	}

	db.Close()
	return false
}

func setDBTokens(access_token string, refresh_token string, usr string) {
	db, err := sql.Open("mysql", "admin:mg545Iop@/goLoginApp")
	checkDBErr(err)

	stmt, err := db.Prepare("update users set access_token=?,refresh_token=? where username=?")
	checkDBErr(err)

	res, err := stmt.Exec(access_token, refresh_token, usr)
	checkDBErr(err)

	affect, err := res.RowsAffected()
	checkDBErr(err)

	fmt.Println(affect)

	db.Close()
}

func writeDBUser(usr string, psw string, email string) bool {
	db, err := sql.Open("mysql", "admin:mg545Iop@/goLoginApp")
	checkDBErr(err)

	stmt, err := db.Prepare("INSERT users SET username=?,email=?,password=?")
	checkDBErr(err)

	res, err := stmt.Exec(usr, email, psw)
	checkDBErr(err)

	id, err := res.LastInsertId()
	checkDBErr(err)

	if id != -1 {
		db.Close()
		return true
	}

	db.Close()
	return false
}

func checkDBUser(usr string, psw string) bool {
	db, err := sql.Open("mysql", "admin:mg545Iop@/goLoginApp")
	checkDBErr(err)

	usernameRow, err := db.Query("SELECT username FROM users")
	checkDBErr(err)

	passwordRow, err := db.Query("SELECT password FROM users")
	checkDBErr(err)

	for usernameRow.Next(){
		var username string
		err = usernameRow.Scan(&username)
		checkDBErr(err)
		if usr == username {
			for passwordRow.Next(){
				var password string
				err = passwordRow.Scan(&password)
				checkDBErr(err)
				if psw == password {
					db.Close()
					return true
				}
			}
		}
	}

	db.Close()
	return false
}

func myDatabase() {
	db, err := sql.Open("mysql", "admin:mg545Iop@/goLoginApp")
	checkDBErr(err)

	rows, err := db.Query("SELECT id, username, email, password FROM users")
	checkDBErr(err)

	for rows.Next() {
		var id int
		var username string
		var email string
		var password string
		err = rows.Scan(&id, &username, &email, &password)
		checkDBErr(err)
		fmt.Println(id)
		fmt.Println(username)
		fmt.Println(email)
		fmt.Println(password)
	}

	db.Close()

}

func main() {

	fmt.Println("Database:")
	myDatabase()
	fmt.Println("")

	_, err := initServer(false)
	if err != nil {
		fmt.Println("Main error")
		log.Fatal("Server Error: ", err)
	}

}

var router = mux.NewRouter()

func initServer(isTesting bool) (bool, error) {
	//init
	rand.Seed(time.Now().UnixNano())

	router.HandleFunc("/", indexPageHandler)
	router.HandleFunc("/signinform", signinformHandler)
	router.HandleFunc("/internal", internalPageHandler)

	router.HandleFunc("/login", loginHandler).Methods("POST")
	router.HandleFunc("/logout", logoutHandler).Methods("POST")
	router.HandleFunc("/signin", signinHandler).Methods("POST")


	router.HandleFunc("/token", tokenHandler).Methods("POST")

	http.Handle("/", router)
	http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("resources"))))
	ports := generatePortArray()
	var error error;
	var isCreated bool;

	// ERROR DEBUG
	//ports[0] = ":80"

	//start server
	if ports[0] != "" {
		cs, err := createServer(ports, isTesting)
		isCreated = cs
		if err != nil {
			error := err;
			return isCreated, error
		} else {
			return isCreated, nil
		}
	} else {
		log.Fatal("No port list created!")
	}

	return isCreated, error

}

func openBrowser(port string){
	fmt.Println("Server started at http://localhost", port, "\n")
	url := "http://localhost" + port
	open.Run(url)
}

func createServer(ports [5]string, isTesting bool) (bool, error) {

	var error error;

	for _, p := range ports {
		li, err := net.Listen("tcp", p)
		error = err
		if error != nil {
			fmt.Println("Cannot use this Port", p, "retry...")
		} else {
			if (isTesting == false){
				openBrowser(p)
				http.Serve(li, nil)
			} else {
				return true, nil
			}
		}
	}
	return false, error
}

func generatePortArray() [5]string {

	var ports [5]string

	for i := 0; i < len(ports); i++{
		ports[i] = ":" + strconv.Itoa(randomInt(1025, 60000));
	}

	return ports
}

func generateToken() string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

	myToken := make([]byte, 25)
	for i := range myToken {
		myToken[i] = letters[rand.Intn(len(letters))]
	}
	return string(myToken)
}

func randomInt(min, max int) int {
	return rand.Intn(max - min) + min
}
