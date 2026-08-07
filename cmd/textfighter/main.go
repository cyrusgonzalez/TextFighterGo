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

// main game loop
// on Match struct bool loop
// remove forever for with prompt - done
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


// Weapon selection function each turn
// pull from Struct Weapon
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
// hit message on attack
func prompt(reader *bufio.Scanner, msg string) string {
	fmt.Print(msg)
	reader.Scan()
	return strings.TrimSpace(reader.Text())
}
