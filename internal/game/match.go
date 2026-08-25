package game

import "math/rand/v2"

// Match owns one fight between two players.
type Match struct {
	P1, P2 *Player
	turn   int // 1 or 2, whose turn it is
}

// NewMatch sets up a fight and picks a random starting turn.
func NewMatch(p1, p2 *Player) *Match {
	return &Match{
		P1:   p1,
		P2:   p2,
		turn: rand.IntN(2) + 1,
	}
}

// TurnResult reports what happened on one turn.
type TurnResult struct {
	Attacker, Defender *Player
	AttackResult
	Winner *Player // nil unless this turn ended the match
}

// PlayTurn resolves one attack with the given weapon and swaps whose turn
// it is next.
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

// Turn returns whose turn it currently is (1 or 2).
func (m *Match) Turn() int {
	return m.turn
}
