package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {

	rand.Seed(time.Now().UnixNano())

	game()

}

func game(){

	r := rand.Intn(100)
	userNumber := 101
	fmt.Println(r)
	playAgain := ""

	for userNumber != r {

		fmt.Print("What is my Number: ")
		fmt.Scan(&userNumber)
		fmt.Println(userNumber)
		if(userNumber < r) {
			fmt.Println("My Number is bigger!")
		}else if(userNumber > r){

			fmt.Println("My Number is smaller!")
		}else{
			fmt.Println("Thats right\n")
			fmt.Print("Play Again? (y/n)")
			fmt.Scan(&playAgain)
			if playAgain == "y" {
				game()
			}
		}
	}
}
