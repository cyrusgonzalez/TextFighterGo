package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"TextFighterGo/internal/game"
	"TextFighterGo/internal/session"
)

const acceptTimeout = 2 * time.Minute

func main() {
	listenAddr := flag.String("listen", "", "host a match: TCP address to listen on for an opponent, e.g. :9000")
	flag.Parse()

	local := session.New(os.Stdin, os.Stdout)
	p2Sess := local // default: shared terminal, same as before -listen existed

	if *listenAddr != "" {
		conn, err := hostAndWaitForOpponent(local, *listenAddr)
		if err != nil {
			local.Printf("Failed to host: %v\n", err)
			os.Exit(1)
		}
		defer conn.Close()
		p2Sess = session.New(conn, conn)
	}

	local.Println("Prepare to battle to the death!!!")
	p1 := game.NewPlayer(mustPrompt(local, "Player 1, enter your name: "))

	p2Prompt := "Player 2, enter your name: "
	if p2Sess != local {
		p2Prompt = "Enter your name: "
	}
	p2 := game.NewPlayer(mustPrompt(p2Sess, p2Prompt))

	for {
		playMatch(local, p2Sess, p1, p2)

		answer, ok := local.Prompt("Fight again? (yes/[no]): ")
		if !ok || strings.ToLower(answer) != "yes" {
			local.Println("The End")
			return
		}
		p1.Reset()
		p2.Reset()
	}
}

// hostAndWaitForOpponent listens on addr and waits (up to acceptTimeout)
// for one incoming connection, treated as player 2. Any TCP client works,
// e.g. `nc host port`.
func hostAndWaitForOpponent(local *session.Session, addr string) (net.Conn, error) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	defer ln.Close()

	local.Printf("Waiting for an opponent to connect on %s...\n", addr)

	connCh := make(chan net.Conn, 1)
	errCh := make(chan error, 1)
	go func() {
		conn, err := ln.Accept()
		if err != nil {
			errCh <- err
			return
		}
		connCh <- conn
	}()

	select {
	case conn := <-connCh:
		local.Println("Opponent connected!")
		return conn, nil
	case err := <-errCh:
		return nil, err
	case <-time.After(acceptTimeout):
		return nil, fmt.Errorf("timed out after %s waiting for an opponent", acceptTimeout)
	}
}

// playMatch runs one fight to completion, one simultaneous round at a time.
// s1/s2 are p1/p2's own sessions - same *Session for local play,
// or one local + one network-backed *Session once -listen is used.
func playMatch(s1, s2 *session.Session, p1, p2 *game.Player) {
	m := game.NewMatch(p1, p2)
	broadcast(s1, s2, fmt.Sprintf("\nFight between %s and %s! FIGHT!\n\n", p1.Name, p2.Name))

	for !m.Over() {
		a1, a2, ok := gatherActions(s1, s2, p1, p2)
		if !ok {
			broadcast(s1, s2, "A player disconnected. Match aborted.\n")
			return
		}

		result := m.PlayRound(a1, a2)
		broadcast(s1, s2, narrate(result))

		switch {
		case result.Draw:
			broadcast(s1, s2, "Both fighters fall! It's a DRAW!\n")
		case result.Winner != nil:
			broadcast(s1, s2, fmt.Sprintf("We have a winner! %s WINS!\n", result.Winner.Name))
		}
	}
}

// gatherActions collects both players' choices for one round.
//
// When s1 and s2 are the same *Session (local shared-terminal play), the
// two prompts are read sequentially - a single bufio.Scanner isn't safe
// for two goroutines to read at once (nothing stops both from racing to
// read the same next line), so there's no safe way to make that case
// genuinely concurrent. The round is still resolved together afterward
// either way, so neither player's outcome is revealed until both have
// committed.
//
// When s1 and s2 are different sessions (a network opponent has their own
// net.Conn), each has its own independent scanner, so there's no shared
// state to race on - reading both at once via goroutines is both safe and
// meaningful: neither player can see the other's choice before making
// their own.
func gatherActions(s1, s2 *session.Session, p1, p2 *game.Player) (a1, a2 game.Action, ok bool) {
	if s1 == s2 {
		a1, ok1 := chooseAction(s1, p1.Name)
		a2, ok2 := chooseAction(s2, p2.Name)
		return a1, a2, ok1 && ok2
	}

	type result struct {
		action game.Action
		ok     bool
	}
	ch1 := make(chan result, 1)
	ch2 := make(chan result, 1)
	go func() {
		a, ok := chooseAction(s1, p1.Name)
		ch1 <- result{a, ok}
	}()
	go func() {
		a, ok := chooseAction(s2, p2.Name)
		ch2 <- result{a, ok}
	}()
	r1, r2 := <-ch1, <-ch2
	return r1.action, r2.action, r1.ok && r2.ok
}

// chooseAction prompts until it gets a valid move.
func chooseAction(sess *session.Session, name string) (game.Action, bool) {
	const menu = "%s, choose your move:\n" +
		"  1) Attack - Ax\n" +
		"  2) Attack - Sword\n" +
		"  3) Attack - Club\n" +
		"  4) Block/Armor up\n" +
		"  5) Dodge\n" +
		"  6) Counter\n" +
		"> "
	for {
		line, ok := sess.Prompt(fmt.Sprintf(menu, name))
		if !ok {
			return game.Action{}, false
		}
		n, err := strconv.Atoi(strings.TrimSpace(line))
		if err != nil {
			sess.Println("Please enter a number 1-6.")
			continue
		}
		switch n {
		case 1:
			return game.Action{Kind: game.ActionAttack, Weapon: game.Ax}, true
		case 2:
			return game.Action{Kind: game.ActionAttack, Weapon: game.Sword}, true
		case 3:
			return game.Action{Kind: game.ActionAttack, Weapon: game.Club}, true
		case 4:
			return game.Action{Kind: game.ActionBlock}, true
		case 5:
			return game.Action{Kind: game.ActionDodge}, true
		case 6:
			return game.Action{Kind: game.ActionCounter}, true
		default:
			sess.Println("Please choose 1-6.")
		}
	}
}

// narrate turns a RoundResult into the text shown to both players.
func narrate(r game.RoundResult) string {
	var b strings.Builder
	b.WriteString(narrateOutcome(r.P1, r.P2))
	b.WriteString(narrateOutcome(r.P2, r.P1))
	b.WriteString(fmt.Sprintf("%s: %d HP | %s: %d HP\n\n",
		r.P1.Player.Name, r.P1.Player.Health, r.P2.Player.Name, r.P2.Player.Health))
	return b.String()
}

// narrateOutcome describes what `you` did this round. `opp` is needed too:
// e.g. telling "you dodged nothing" apart from "you dodged a real hit"
// depends on whether the opponent actually attacked.
func narrateOutcome(you, opp game.PlayerRoundOutcome) string {
	name := you.Player.Name

	// oppMissed is true when the opponent DID attack but the roll itself
	// whiffed (0 damage) - distinct from opp.AttackRoll == nil, where the
	// opponent didn't attack at all. Both cases mean "your defense didn't
	// actually need to do anything," so Block/Dodge/Counter all check this
	// first rather than claiming credit (or blame) for a hit that was
	// never coming.
	oppMissed := opp.AttackRoll != nil && opp.AttackRoll.Outcome == "miss"

	switch you.Action.Kind {
	case game.ActionAttack:
		var msg string
		switch you.AttackRoll.Outcome {
		case "miss":
			msg = fmt.Sprintf("%s attacked with the %s and missed!\n", name, you.Action.Weapon.Name)
		case "crit":
			msg = fmt.Sprintf("%s landed a CRITICAL HIT with the %s for %d damage!\n",
				name, you.Action.Weapon.Name, you.AttackRoll.Damage)
		default:
			msg = fmt.Sprintf("%s hit with the %s for %d damage!\n", name, you.Action.Weapon.Name, you.AttackRoll.Damage)
		}
		if opp.Countered && you.DamageTaken > 0 {
			msg += fmt.Sprintf("...but %s countered! %s takes %d damage right back!\n", opp.Player.Name, name, you.DamageTaken)
		}
		return msg

	case game.ActionBlock:
		switch {
		case opp.AttackRoll == nil:
			return fmt.Sprintf("%s braced behind their armor, but nothing came.\n", name)
		case oppMissed:
			return fmt.Sprintf("%s's attack missed anyway - the block wasn't even needed.\n", opp.Player.Name)
		case you.DamageTaken == 0:
			return fmt.Sprintf("%s blocked, fully absorbing the hit!\n", name)
		default:
			return fmt.Sprintf("%s blocked, taking %d damage through their armor.\n", name, you.DamageTaken)
		}

	case game.ActionDodge:
		switch {
		case opp.AttackRoll == nil:
			return fmt.Sprintf("%s stayed light on their feet - no attack to dodge.\n", name)
		case oppMissed:
			return fmt.Sprintf("%s's attack missed anyway - no need to dodge.\n", opp.Player.Name)
		case you.Dodged:
			return fmt.Sprintf("%s dodged the attack completely!\n", name)
		default:
			return fmt.Sprintf("%s tried to dodge and failed, taking %d damage!\n", name, you.DamageTaken)
		}

	case game.ActionCounter:
		switch {
		case opp.AttackRoll == nil:
			return fmt.Sprintf("%s waited to counter, but %s never attacked - wasted turn!\n", name, opp.Player.Name)
		case oppMissed:
			return fmt.Sprintf("%s anticipated an attack from %s that missed anyway - nothing to counter.\n", name, opp.Player.Name)
		default:
			return fmt.Sprintf("%s perfectly countered, reflecting %d damage back at %s!\n", name, opp.DamageTaken, opp.Player.Name)
		}
	}
	return ""
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
