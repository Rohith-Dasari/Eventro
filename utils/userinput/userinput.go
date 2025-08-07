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
		return 0, fmt.Errorf("error reading input: %w", err)
	}
	input = strings.TrimSpace(input)
	choice, err := strconv.Atoi(input)
	if err != nil {
		return 0, fmt.Errorf("invalid input: %w", err)
	}
	return choice, nil
}
