package main

import (
	"fmt"

	"github.com/ayofash1/netbox-checker/internal"
)

func main() {
	checker := internal.NewChecker()
	if err := checker.Validate(); err != nil {
		fmt.Println("Validation error:", err)
		return
	}

	result := checker.Check()
	fmt.Println("Check result:", result)
}
