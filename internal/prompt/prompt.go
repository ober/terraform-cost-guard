package prompt

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// ConfirmApply prompts the user to confirm applying the terraform plan
func ConfirmApply(monthlyCostChange float64) (bool, error) {
	var message string

	if monthlyCostChange > 0 {
		message = fmt.Sprintf("\n\033[1;33mHey, these changes will cost an additional $%.2f/month. Proceed? [y/N]\033[0m ", monthlyCostChange)
	} else if monthlyCostChange < 0 {
		message = fmt.Sprintf("\n\033[1;32mThese changes will save $%.2f/month. Proceed? [y/N]\033[0m ", -monthlyCostChange)
	} else {
		message = "\n\033[1;34mNo significant cost change detected. Proceed? [y/N]\033[0m "
	}

	fmt.Print(message)

	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return false, fmt.Errorf("failed to read response: %w", err)
	}

	response = strings.TrimSpace(strings.ToLower(response))

	return response == "y" || response == "yes", nil
}

// ConfirmWithThreshold prompts only if cost exceeds threshold
func ConfirmWithThreshold(monthlyCostChange float64, threshold float64) (bool, error) {
	if monthlyCostChange <= threshold {
		fmt.Printf("\033[1;32mCost change ($%.2f/month) is within threshold ($%.2f). Proceeding...\033[0m\n",
			monthlyCostChange, threshold)
		return true, nil
	}

	return ConfirmApply(monthlyCostChange)
}

// PrintCostSummary prints a detailed cost summary
func PrintCostSummary(totalChange float64, created, destroyed, updated int, unsupportedTypes []string) {
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("                    COST ESTIMATE SUMMARY")
	fmt.Println(strings.Repeat("=", 60))

	fmt.Printf("\n  Resources to be created:   %d\n", created)
	fmt.Printf("  Resources to be destroyed: %d\n", destroyed)
	fmt.Printf("  Resources to be updated:   %d\n", updated)

	fmt.Println("\n" + strings.Repeat("-", 60))

	if totalChange > 0 {
		fmt.Printf("\n  \033[1;33mEstimated Monthly Cost Increase: +$%.2f\033[0m\n", totalChange)
	} else if totalChange < 0 {
		fmt.Printf("\n  \033[1;32mEstimated Monthly Cost Savings: -$%.2f\033[0m\n", -totalChange)
	} else {
		fmt.Printf("\n  \033[1;34mNo significant cost change\033[0m\n")
	}

	if len(unsupportedTypes) > 0 {
		fmt.Println("\n  Note: The following resource types are not yet supported")
		fmt.Println("  for cost estimation (estimated as $0):")
		for _, t := range unsupportedTypes {
			fmt.Printf("    - %s\n", t)
		}
	}

	fmt.Println("\n" + strings.Repeat("=", 60))
}
