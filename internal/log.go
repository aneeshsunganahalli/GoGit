package internal

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func FormatLog(currentHash string, authorLine string, message []string, isFirstCommit bool) {

	// 1. Print Commit Hash
	fmt.Printf("\033[33mcommit %s\033[0m", currentHash)

	branch, _ := getBranchName()

	if isFirstCommit {
		fmt.Printf(" \033[1;36m(\033[1;34mHEAD -> \033[1;32m%s\033[1;36m)\033[0m", branch)
	}
	fmt.Println()

	// 2. Parse Author & Date details
	parts := strings.Fields(authorLine)
	if len(parts) >= 3 {
		timezone := parts[len(parts)-1]
		timestampStr := parts[len(parts)-2]

		authorId := strings.Join(parts[:len(parts)-2], " ")

		sec, _ := strconv.ParseInt(timestampStr, 10, 64)
		t := time.Unix(sec, 0)
		// Mon Apr 13 22:38:00 2026
		dateStr := t.Format("Mon Jan 2 15:04:05 2006")

		fmt.Printf("Author: %s\n", authorId)
		fmt.Printf("Date:   %s %s\n", dateStr, timezone)

	}

	fmt.Println("")
	for _, m := range message {
		if m != "" {
			fmt.Printf("	%s\n", m)
		}
	}
	fmt.Println("")
}

func getBranchName() (string, error) {
	headContent, err := os.ReadFile(".gogit/HEAD")
	if err != nil {
		fmt.Println("Cannot read HEAD")
		return "", err
	}

	content := strings.TrimSpace(string(headContent))
	if content == "" {
		return "", nil
	}
	var branch string

	if strings.HasPrefix(content, "refs:refs/heads") {
		branch = strings.TrimPrefix(content, "refs:refs/heads/")
	} 

	return branch, err
}
