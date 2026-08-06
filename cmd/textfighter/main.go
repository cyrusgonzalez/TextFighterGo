// Package main is the terminal entrypoint. It lives under cmd/textfighter
// instead of the repo root — that's the standard Go project layout: each
// buildable program gets its own folder under cmd/, and internal/ holds
// the packages that back it. Right now there's one binary (this terminal
// game); once we add a network server, it'll likely get its own
// cmd/server folder reusing the exact same internal/game package.
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"TextFighterGo/internal/game"
)

func main() {
	// bufio.Scanner reads stdin line by line. This replaces game.py's
	// input() calls — Go has no built-in input(); reading stdin is always
	// this explicit "wrap os.Stdin, then .Scan()/.Text()" dance.
	reader := bufio.NewScanner(os.Stdin)

	fmt.Println("Prepare to battle to the death!!!")
	p1 := game.NewPlayer(prompt(reader, "Player 1, enter your name: "))
	p2 := game.NewPlayer(prompt(reader, "Player 2, enter your name: "))

	for {
		playMatch(reader, p1, p2)

		answer := prompt(reader, "Fight again? (yes/no): ")
		if strings.ToLower(answer) != "yes" {
			fmt.Println("The End")
			return
		}
		p1.Reset()
		p2.Reset()
	}
}

// playMatch runs one fight to completion. game.py did this recursively
// (main() calling itself again to replay) — Go favors the plain for-loop
// in main() above instead; recursion here would just grow the call stack
// for no reason.
func playMatch(reader *bufio.Scanner, p1, p2 *game.Player) {
	m := game.NewMatch(p1, p2)
	fmt.Printf("\nFight between %s and %s! FIGHT!\n\n", p1.Name, p2.Name)

	for !m.Over() {
		attacker := p1
		if m.Turn() == 2 {
			attacker = p2
		}

		weapon := chooseWeapon(reader, attacker.Name)
		result := m.PlayTurn(weapon)

		switch result.Outcome {
		case "miss":
			fmt.Printf("%s attacked with the %s and missed!\n", result.Attacker.Name, weapon.Name)
		case "crit":
			fmt.Printf("%s landed a CRITICAL HIT with the %s for %d damage!\n",
				result.Attacker.Name, weapon.Name, result.Damage)
		default:
			fmt.Printf("%s hit %s with the %s for %d damage!\n",
				result.Attacker.Name, result.Defender.Name, weapon.Name, result.Damage)
		}
		fmt.Printf("%s has %d HP left.\n\n", result.Defender.Name, result.Defender.Health)

		if result.Winner != nil {
			fmt.Printf("We have a winner! %s WINS!\n", result.Winner.Name)
		}
	}
}

// chooseWeapon prompts until it gets a valid weapon number. game.py's
// equivalent (`int(input(...))`) just crashes with an uncaught exception on
// bad input. Go has no exceptions at all — strconv.Atoi returns (int,
// error) instead of throwing, and the `err != nil` check is how you handle
// "that wasn't valid" without a try/except anywhere in sight. You'll see
// this (value, error) return pattern constantly in Go.
func chooseWeapon(reader *bufio.Scanner, name string) game.Weapon {
	for {
		fmt.Printf("%s, choose your weapon (1=Ax, 2=Sword, 3=Club): ", name)
		if !reader.Scan() {
			os.Exit(0)
		}
		n, err := strconv.Atoi(strings.TrimSpace(reader.Text()))
		if err != nil {
			fmt.Println("Please enter a number 1-3.")
			continue
		}
		if w, ok := game.Weapons[n]; ok {
			return w
		}
		fmt.Println("Please choose 1, 2, or 3.")
	}
}

// prompt prints msg, reads one line of stdin, and returns it trimmed.
func prompt(reader *bufio.Scanner, msg string) string {
	fmt.Print(msg)
	reader.Scan()
	return strings.TrimSpace(reader.Text())
}
