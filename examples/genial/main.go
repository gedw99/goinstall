package main

import (
	"github.com/audrenbdb/goinstall"
	"golang.org/x/text/language"
)

func main() {
	installer := goinstall.New("App Géniale", language.French)

	installer.AddCondition("Condition 1", loremIpsumContent)

	installer.AddStep(func() error{ return nil }, "Lancement de géniales fusées...")
	installer.AddStep(func() error{ return nil }, "Créations de géniaux raccourcis...")
	installer.AddStep(func() error{ return nil }, "Connection avec de géniaux dauphins...")

	installer.StartWithGUI()
}

const (
loremIpsumContent = `Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.`
)