package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ayofash1/netbox-checker/internal/netbox"
	"github.com/ayofash1/netbox-checker/internal/rules"
)

func main() {

	netboxURL := os.Getenv("NETBOX_URL")

	netboxToken := "f9558e5a2322369437ef631f606e84754d72bcef"

	ruleSet, err := rules.ParseRules("compliance-rules.yaml")
	if err != nil {
		log.Fatal(err)
	}

	for _, r := range ruleSet.Rules {
		fmt.Printf("Loaded rule: [%s] %s (%s)\n", r.ID, r.Description, r.Type)
	}

	client := netbox.NewClient(netboxURL, netboxToken)
	devices, err := client.FetchDevices()
	if err != nil {
		log.Fatalf("Failed to fetch NetBox devices: %v", err)
	}

	results := rules.Evaluate(devices, ruleSet)
	fmt.Printf("🔍 Evaluated compliance for all devices\n")

	// === Output Summary ===
	fmt.Println("\n📋 Compliance Report:")
	fmt.Println("======================")
	for _, r := range results {
		status := "✅"
		if !r.Compliant {
			status = "❌"
		}
		fmt.Printf("%s [%s] Device: %-15s | %s\n", status, r.RuleID, r.DeviceName, r.Message)
	}
}
