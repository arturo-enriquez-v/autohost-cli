// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (c) 2026 Arturo (AutoHost)

package prompt

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

const (
	yellow = "\033[33m"
	reset  = "\033[0m"
)

// AskInput shows a prompt and reads user input, falling back to def if empty.
func AskInput(reader *bufio.Reader, message, def string) string {
	fmt.Printf("%s [%s%s%s]: ", message, yellow, def, reset)
	val, _ := reader.ReadString('\n')
	val = strings.TrimSpace(val)
	if val == "" {
		return def
	}
	return val
}

// Confirm asks a yes/no question, returning true only if user types 'y' or 'yes'.
func Confirm(message string) bool {
	fmt.Print(message)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.ToLower(strings.TrimSpace(input))
	return input == "y" || input == "yes"
}

// AskOption prompts user with a numbered list of options and returns the chosen string.
func AskOption(message string, options []string) string {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println(message)
		for i, opt := range options {
			fmt.Printf("[%d] %s\n", i+1, opt)
		}
		fmt.Print("Elige una opción: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if i, err := strconv.Atoi(input); err == nil && i >= 1 && i <= len(options) {
			return options[i-1]
		}
		fmt.Println("❌ Opción inválida, intenta de nuevo.")
	}
}

// ValidPort checks if portStr is a valid integer and not currently bound.
func ValidPort(portStr string) (int, error) {
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return 0, fmt.Errorf("puerto inválido: %s", portStr)
	}

	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return 0, fmt.Errorf("el puerto %d ya está en uso", port)
	}
	ln.Close()

	return port, nil
}

// AskAppPort repeatedly prompts for a port until a valid, unbound port is entered.
func AskAppPort(reader *bufio.Reader, message string, port string) string {
	for {
		port = AskInput(reader, message, port)
		if _, err := ValidPort(port); err == nil {
			break
		} else {
			fmt.Println("❌", err)
		}
	}
	return port
}

