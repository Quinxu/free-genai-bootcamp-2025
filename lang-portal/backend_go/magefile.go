//+build mage

package main

import (
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Default target to run when none is specified
var Default = Build

// Build builds the application
func Build() error {
	return sh.Run("go", "build", "./cmd/server")
}

// Run runs the application
func Run() error {
	return sh.Run("go", "run", "./cmd/server")
}

// Migrate runs database migrations
func Migrate() error {
	// TODO: Implement migration logic
	return nil
}

// Seed seeds the database with initial data
func Seed() error {
	// TODO: Implement seeding logic
	return nil
}

// Test runs tests
func Test() error {
	return sh.Run("go", "test", "./...")
}

// Clean cleans build artifacts
func Clean() error {
	return sh.Run("go", "clean")
}
