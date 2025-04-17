// Copyright (c) 2024 The Flokicoin developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package utils

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"golang.org/x/term"
)

func ReadPassword(prompt string, confirm bool) []byte {
	// Save the current terminal state
	oldState, err := term.GetState(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatalf("Failed to get terminal state: %v", err)
	}

	// Ensure the terminal state is restored to the original state
	defer func() {
		if err := term.Restore(int(os.Stdin.Fd()), oldState); err != nil {
			log.Fatalf("Failed to restore terminal state: %v", err)
		}
	}()

	// Prompt the user and read the password securely
	fmt.Print(prompt)
	password, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatalf("Failed to read password: %v", err)
	}
	fmt.Println() // Move to the next line after password input

	// If confirmation is required, prompt and validate
	if confirm {
		fmt.Print("Confirm password: ")
		confirmPassword, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			log.Fatalf("Failed to read confirmation password: %v", err)
		}
		fmt.Println() // Move to the next line after confirmation input

		if string(password) != string(confirmPassword) {
			log.Fatalf("Passwords do not match. Please try again.")
		}
	}

	return password
}

func ReadMnemonic() string {
	fmt.Println("Enter your mnemonic phrase:")
	fmt.Println("- If entering all words at once, type them and press Enter.")
	fmt.Println("- If entering word by word, type each word and press Enter. Continue until all words are entered.")
	fmt.Println("")

	reader := bufio.NewReader(os.Stdin)
	var words []string

	for {
		fmt.Print("Mnemonic (type 'done' when finished): ")
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Failed to read input: %v", err)
		}

		line = strings.TrimSpace(line)

		// Detect end of input
		if strings.ToLower(line) == "done" {
			break
		}

		// Split input into words
		inputWords := strings.Fields(line)
		if len(inputWords) > 1 && len(words) == 0 {
			// User entered multiple words at once
			words = append(words, inputWords...)
			break
		} else if len(inputWords) == 1 {
			// User entered one word
			words = append(words, inputWords[0])
		} else if len(inputWords) == 0 {
			// No valid input
			fmt.Println("Invalid input. Please enter a word or multiple words.")
			continue
		}
	}

	return strings.Join(words, " ")
}

func ReadLine(prompt string, validator func(string) error) (string, error) {
	for {
		fmt.Print(prompt)
		var input string
		if _, err := fmt.Scanln(&input); err != nil {
			// Could return the error or just continue looping
			// depending on how you want to handle read failures.
			log.Printf("Failed to read user input: %v\n", err)
			continue
		}

		// If user wants to exit/cancel
		if input == "exit" {
			// Return an error or handle it however you prefer:
			// return "", errors.New("user chose to exit")
			// or do: os.Exit(0)
			return "", errors.New("user chose to exit")
		}

		// Validate the input using the callback
		if err := validator(input); err != nil {
			fmt.Printf("Invalid input: %v. Try again or type 'exit' to cancel.\n", err)
			continue
		}

		// If everything is good, return the valid input
		return input, nil
	}
}

func ReadAmount() (float64, error) {
	// Validator function for FLC amount
	amountValidator := func(input string) error {
		if input == "" {
			return errors.New("amount cannot be empty")
		}
		_, err := strconv.ParseFloat(input, 64)
		if err != nil {
			return fmt.Errorf("cannot parse amount as float: %v", err)
		}
		return nil
	}

	amountStr, err := ReadLine("Enter the amount to send (in FLC): ", amountValidator)
	if err != nil {
		return 0, err
	}

	// Convert the string to float64 now that it’s validated
	amountFloat, _ := strconv.ParseFloat(amountStr, 64)
	return amountFloat, nil
}

func ReadAddress() (string, error) {
	addressValidator := func(input string) error {
		if input == "" {
			return errors.New("address cannot be empty")
		}
		// Optionally decode the address here to ensure it's valid.
		// e.g. _, err := chainutil.DecodeAddress(input, &chaincfg.MainNetParams)
		// if err != nil {
		//     return fmt.Errorf("invalid FLC address: %v", err)
		// }
		return nil
	}

	return ReadLine("Enter the destination address: ", addressValidator)
}
