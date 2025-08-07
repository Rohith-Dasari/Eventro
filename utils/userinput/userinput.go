package utils

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func TakeUserInput() (int, error) {
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return 0, err
	}
	input = strings.TrimSpace(input)
	choice, err := strconv.Atoi(input)
	if err != nil {
		fmt.Println("Invalid input. Please enter a number.")
		return 0, err
	}

	return choice, nil
}
