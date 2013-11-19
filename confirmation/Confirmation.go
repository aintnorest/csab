package Confirmation

import (
	"fmt"
	"regexp"
)

func ConfirmationPrompt(input string) (result bool, err error) {

	regYes, err := regexp.Compile("^([yY][eE][sS]|[Yy])$")
	regNo, err := regexp.Compile("^([nN][oO]|[nN])$")

	yesresult := regYes.MatchString(input)
	noresult := regNo.MatchString(input)

	if yesresult == true {
		return true, nil
	} 
	if noresult == true {
		return false, nil
	}

	fmt.Println("Please choose Yes or No.")
	var retry string
	fmt.Scanln(&retry)
	return ConfirmationPrompt(retry)

}
/*
func main() {
fmt.Println("Yes or No")
	var input string
	fmt.Scanln(&input)
	 yn, err := ConfirmationPrompt(input)
	 if err != nil {
	 	fmt.Println("Fucked Up")
	 }
	 fmt.Println(yn)
}
*/