package game

import "math/rand/v2"

// Match owns one fight between two players. game.py handled this with a
// wall of nested while-loops and a bare `playerTurn` int floating in
// main(). Bundling that state into a type with methods is the Go-flavored
// fix: no globals, and the turn-taking logic lives in one place instead of
// being duplicated per `if playerTurn == 1 / elif playerTurn == 2` branch.
type Match struct {
	P1, P2 *Player
	turn   int // 1 or 2, whose turn it is. Unexported (lowercase): only
	// this package can change it directly — callers go through
	// PlayTurn/Turn() instead, same idea as a "private" field.
}

// NewMatch sets up a fight and picks a random starting turn, replacing
// game.py's `playerTurn = random.randint(1,2)`.
func NewMatch(p1, p2 *Player) *Match {
	return &Match{
		P1:   p1,
		P2:   p2,
		turn: rand.IntN(2) + 1,
	}
}

// TurnResult reports what happened on one turn, for the caller to display —
// a local terminal print today, a message written to a network session
// later.
//
// The embedded AttackResult (a field with no name, just the type) means
// TurnResult.Damage and TurnResult.Outcome work directly — Go's
// lightweight stand-in for "has-a, but let callers reach through it like it
// were part of me." It's not inheritance (no is-a relationship, no
// polymorphism), just field/method promotion.
type TurnResult struct {
	Attacker, Defender *Player
	AttackResult
	Winner *Player // nil unless this turn ended the match
}

// PlayTurn resolves one attack with the given weapon and swaps whose turn
// it is next. This one method replaces both the `if playerTurn == 1` and
// `elif playerTurn == 2` blocks from game.py's main() — they were doing the
// identical thing to swapped players, which is exactly the duplication a
// method like this exists to remove.
func (m *Match) PlayTurn(weapon Weapon) TurnResult {
	attacker, defender := m.P1, m.P2
	if m.turn == 2 {
		attacker, defender = m.P2, m.P1
	}

	res := weapon.Attack()
	defender.TakeDamage(res.Damage)

	tr := TurnResult{Attacker: attacker, Defender: defender, AttackResult: res}
	if !defender.Alive() {
		tr.Winner = attacker
		attacker.Wins++
	}

	if m.turn == 1 {
		m.turn = 2
	} else {
		m.turn = 1
	}
	return tr
}

// Over reports whether the match has a winner yet.
func (m *Match) Over() bool {
	return !m.P1.Alive() || !m.P2.Alive()
}

// Turn returns whose turn it currently is (1 or 2) so a caller (terminal
// loop, or later a network handler) knows who to prompt for a move.
func (m *Match) Turn() int {
	return m.turn
}
