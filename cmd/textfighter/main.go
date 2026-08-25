package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"TextFighterGo/internal/game"
	"TextFighterGo/internal/session"
)

func main() {
	local := session.New(os.Stdin, os.Stdout)

	local.Println("Prepare to battle to the death!!!")
	p1 := game.NewPlayer(mustPrompt(local, "Player 1, enter your name: "))
	p2 := game.NewPlayer(mustPrompt(local, "Player 2, enter your name: "))

	for {
		playMatch(local, local, p1, p2)

		answer, ok := local.Prompt("Fight again? (yes/no): ")
		if !ok || strings.ToLower(answer) != "yes" {
			local.Println("The End")
			return
		}
		p1.Reset()
		p2.Reset()
	}
}

// main game loop
// on Match struct bool loop
// remove forever for with prompt - done
// s1/s2 are p1/p2's own sessions - same *Session today (shared terminal),
// different ones once networking exists.
func playMatch(s1, s2 *session.Session, p1, p2 *game.Player) {
	m := game.NewMatch(p1, p2)
	broadcast(s1, s2, fmt.Sprintf("\nFight between %s and %s! FIGHT!\n\n", p1.Name, p2.Name))

	for !m.Over() {
		attacker, sess := p1, s1
		if m.Turn() == 2 {
			attacker, sess = p2, s2
		}

		weapon, ok := chooseWeapon(sess, attacker.Name)
		if !ok {
			broadcast(s1, s2, fmt.Sprintf("%s disconnected. Match aborted.\n", attacker.Name))
			return
		}
		result := m.PlayTurn(weapon)

		switch result.Outcome {
		case "miss":
			broadcast(s1, s2, fmt.Sprintf("%s attacked with the %s and missed!\n", result.Attacker.Name, weapon.Name))
		case "crit":
			broadcast(s1, s2, fmt.Sprintf("%s landed a CRITICAL HIT with the %s for %d damage!\n",
				result.Attacker.Name, weapon.Name, result.Damage))
		default:
			broadcast(s1, s2, fmt.Sprintf("%s hit %s with the %s for %d damage!\n",
				result.Attacker.Name, result.Defender.Name, weapon.Name, result.Damage))
		}
		broadcast(s1, s2, fmt.Sprintf("%s has %d HP left.\n\n", result.Defender.Name, result.Defender.Health))

		if result.Winner != nil {
			broadcast(s1, s2, fmt.Sprintf("We have a winner! %s WINS!\n", result.Winner.Name))
		}
	}
}

// broadcast writes msg to both sessions. When s1 and s2 are the same
// *Session (today's shared terminal), it's written once - the pointer
// equality check below skips the duplicate. Once s1/s2 are two different
// network sessions, this is what keeps both players' screens in sync.
func broadcast(s1, s2 *session.Session, msg string) {
	s1.Printf("%s", msg)
	if s2 != s1 {
		s2.Printf("%s", msg)
	}
}

// Weapon selection function each turn
// pull from Struct Weapon
func chooseWeapon(sess *session.Session, name string) (game.Weapon, bool) {
	for {
		line, ok := sess.Prompt(fmt.Sprintf("%s, choose your weapon (1=Ax, 2=Sword, 3=Club): ", name))
		if !ok {
			return game.Weapon{}, false
		}
		n, err := strconv.Atoi(line)
		if err != nil {
			sess.Println("Please enter a number 1-3.")
			continue
		}
		if w, ok := game.Weapons[n]; ok {
			return w, true
		}
		sess.Println("Please choose 1, 2, or 3.")
	}
}

// mustPrompt is for the pre-match setup (name entry), where there's no
// in-progress match to abort - if input ends here, there's nothing to do
// but exit.
func mustPrompt(sess *session.Session, msg string) string {
	line, ok := sess.Prompt(msg)
	if !ok {
		os.Exit(0)
	}
	return line
}
