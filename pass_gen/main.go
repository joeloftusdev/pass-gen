package main

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"os"
	"strings"
)

type Alphabet struct {
	pool string
}

const (
	up    = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	low   = "abcdefghijklmnopqrstuvwxyz"
	nums  = "1234567890"
	symbs = "!@#$%^&*()-_=+\\/~?"
)

func NewAlphabet(upInc, lowInc, numbsInc, symbolsInc bool) *Alphabet {
	pool := ""
	if upInc {
		pool += up
	}
	if lowInc {
		pool += low
	}
	if numbsInc {
		pool += nums
	}
	if symbolsInc {
		pool += symbs
	}
	return &Alphabet{pool: pool}
}

func (a *Alphabet) GetAlphabet() string {
	return a.pool
}

type Generator struct {
	alphabet *Alphabet
}

func NewGenerator(includeUpper, includeLower, includeNum, includeSym bool) *Generator {
	return &Generator{alphabet: NewAlphabet(includeUpper, includeLower, includeNum, includeSym)}
}

func (g *Generator) GeneratePassword(length int) string {

	password := make([]byte, length)
	alphabetLength := len(g.alphabet.GetAlphabet())
	randomBytes := make([]byte, length)

	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}

	for i, randomByte := range randomBytes {
		index := int(randomByte) % alphabetLength
		password[i] = g.alphabet.GetAlphabet()[index]
	}
	return string(password)
}

func mainLoop() {
	fmt.Println("Welcome to your Password Generator")
	printMenu()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		userOption := strings.TrimSpace(scanner.Text())

		switch userOption {
		case "1":
			requestPassword()
			printMenu()
		case "2":
			checkPassword()
			printMenu()
		case "3":
			printQuitMessage()
			os.Exit(0)
		default:
			fmt.Println("Select one of the available commands!")
			printMenu()
		}
	}
}

func requestPassword() {
	fmt.Println("\nAnswer the following questions with either yes or no")

	var Upper, Lower, Num, Sym bool
	scanner := bufio.NewScanner(os.Stdin)

	includeOptions := map[string]*bool{
		"Do you want to use Lowercase letters? ": &Lower,
		"Do you want to use Uppercase letters? ": &Upper,
		"Do you want to use Numbers? ":           &Num,
		"Do you want to use Symbols? ":           &Sym,
	}

	for question, include := range includeOptions {
		yesOrNo := false
		for !yesOrNo {
			fmt.Print(question)
			scanner.Scan()
			input := strings.TrimSpace(scanner.Text())
			if input == "yes" {
				*include = true
				yesOrNo = true
			} else if input == "no" {
				*include = false
				yesOrNo = true
			} else {
				fmt.Println("Please answer with either yes or no")
			}
		}
	}

	if !Upper && !Lower && !Num && !Sym {
		fmt.Println("You have selected no characters")
		return
	}

	fmt.Print("How long do you want your password to be? ")
	var length int
	_, err := fmt.Scanln(&length)
	if err != nil {
		fmt.Printf("Error reading password length: %v\n", err)
		return
	}

	if length < 1 {
		fmt.Println("Password length must be a positive number. Please try again.")
		return
	}

	generator := NewGenerator(Upper, Lower, Num, Sym)
	password := generator.GeneratePassword(length)

	fmt.Println("Your generated password is:", password)
}

func checkPassword() {
	fmt.Print("\nEnter your password: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := scanner.Text()

	p := NewPassword(input)

	fmt.Println(p.CalculateScore())
}

func printMenu() {
	fmt.Println("\nEnter 1 - Generate Password")
	fmt.Println("Enter 2 - Password Strength Check")
	fmt.Println("Enter 3 - Quit")
	fmt.Print("\nChoice: ")
}

func printQuitMessage() {
	fmt.Println("Closing the program..")
}

type Password struct {
	Value  string
	Length int
}

func NewPassword(s string) *Password {
	return &Password{Value: s, Length: len(s)}
}

func (p *Password) CharType(C byte) int {
	value := 0
	if C >= 'A' && C <= 'Z' {
		value = 1
	} else if C >= 'a' && C <= 'z' {
		value = 2
	} else if C >= '0' && C <= '9' {
		value = 3
	} else {
		value = 4
	}
	return value
}

func (p *Password) PasswordStrength() int {
	s := p.Value
	IsUpper := false
	IsLower := false
	IsNum := false
	IsSymbol := false
	var typedValue int
	Score := 0

	for i := 0; i < len(s); i++ {
		c := s[i]
		typedValue = p.CharType(c)

		if typedValue == 1 {
			IsUpper = true
		}
		if typedValue == 2 {
			IsLower = true
		}
		if typedValue == 3 {
			IsNum = true
		}
		if typedValue == 4 {
			IsSymbol = true
		}
	}

	if IsUpper {
		Score += 1
	}
	if IsLower {
		Score += 1
	}
	if IsNum {
		Score += 1
	}
	if IsSymbol {
		Score += 1
	}

	if len(s) >= 8 {
		Score += 1
	}
	if len(s) >= 16 {
		Score += 1
	}

	return Score
}

func (p *Password) CalculateScore() string {
	Score := p.PasswordStrength()

	switch {
	case Score == 6:
		return "This is a very good password!"
	case Score >= 4:
		return "Good password, but you can still do better!"
	case Score >= 3:
		return "Medium password. Try making it better!"
	default:
		return "This is a weak password. Generate a new one!"
	}
}

func main() {
	mainLoop()
}
