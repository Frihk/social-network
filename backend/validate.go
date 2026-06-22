package main

import (
	"fmt"
	"os"
)

// Simple validation script to check code structure
func main() {
	checks := []struct {
		name string
		path string
	}{
		{"Database initialization", "pkg/db/sqlite/sqlite.go"},
		{"User model", "models/user.go"},
		{"Auth handlers", "pkg/handlers/auth.go"},
		{"User handlers", "pkg/handlers/user.go"},
		{"Auth middleware", "queries/middleware/auth.go"},
		{"CORS middleware", "queries/middleware/cors.go"},
		{"Server setup", "server/server.go"},
		{"Main entry", "cmd/app/main.go"},
		{"Users table migration up", "pkg/db/migration/sqlite/000001_create_users_table.up.sql"},
		{"Users table migration down", "pkg/db/migration/sqlite/000001_create_users_table.down.sql"},
	}

	fmt.Println("🔍 Validating Backend Structure...")
	fmt.Println()

	allValid := true
	for _, check := range checks {
		if _, err := os.Stat(check.path); err == nil {
			fmt.Printf("✅ %s: %s\n", check.name, check.path)
		} else {
			fmt.Printf("❌ %s: %s (NOT FOUND)\n", check.name, check.path)
			allValid = false
		}
	}

	fmt.Println()
	if allValid {
		fmt.Println("✅ All files present!")
		fmt.Println()
		fmt.Println("📝 Next Steps:")
		fmt.Println("1. Install GCC/MinGW-w64 for SQLite support")
		fmt.Println("2. Run: $env:CGO_ENABLED=1; go build ./cmd/app")
		fmt.Println("3. Run: go test -v ./models/...")
		fmt.Println("4. Run: ./app.exe")
		fmt.Println("5. Test the API endpoints")
	} else {
		fmt.Println("❌ Some files are missing!")
		os.Exit(1)
	}
}
