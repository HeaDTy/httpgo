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
	"os"
)

func getData(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		tmp, err := template.ParseFiles("text.html")
		if err != nil {
			log.Fatal("HTML Template Error: ", err)
		}
		tmp.Execute(w, nil)
	} else {
		fmt.Println("\nGenerated Tokens")
		fmt.Println("access token: ", generateToken())
		fmt.Println("refresh token: ", generateToken())
		os.Exit(1)
	}
}

func userContent(w http.ResponseWriter, r *http.Request) {

	usr := ""
	psw := ""

	if r.Method == "POST" {
		r.ParseForm()
		fmt.Println("Request map: ", r.Form)
		for key, val := range r.Form {
			for _, val := range val {
				fmt.Println("KEY: ", key, "VALUE: ", val)
				switch key {
				case "username":
					usr = val
				case "password":
					psw = val
				}
			}
		}
	}

	// Check if login correct and when true display content.html
	if usr != "" && psw != "" {
		if checkLogin(usr, psw) == true {
			fmt.Println("\nUser successful logged in")
			tmp, err := template.ParseFiles("content.html")
			if err != nil {
				log.Fatal("HTML Template Error: ", err)
			}
			tmp.Execute(w, nil)
		}else{
			log.Fatal("Cannot Login")
		}
	}else {
		log.Fatal("No input")
	}
}

func checkLogin(usr string, psw string) bool {

	if usr == "test" && psw == "1234" {
		return true
	} else {
		fmt.Println("Wrong login!")
		return false
	}

}

func main() {

	_, err := initServer(false)
	if err != nil {
		fmt.Println("Main error")
		log.Fatal("Server Error: ", err)
	}

}

func initServer(isTesting bool) (bool, error) {
	//init
	rand.Seed(time.Now().UnixNano())
	http.HandleFunc("/", getData)
	http.HandleFunc("/userContent", userContent)
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
	//log.Fatal("Server Error: ", error)
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
