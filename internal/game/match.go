package game

import "math/rand/v2"

// Match owns one fight between two players.
type Match struct {
	P1, P2 *Player
}

// NewMatch sets up a fight between two players.
func NewMatch(p1, p2 *Player) *Match {
	return &Match{P1: p1, P2: p2}
}

// PlayerRoundOutcome reports what happened to one player during a round.
type PlayerRoundOutcome struct {
	Player      *Player
	Action      Action
	AttackRoll  *AttackResult // non-nil only if Action.Kind == ActionAttack
	DamageTaken int
	Dodged      bool // true if this player attempted Dodge and it succeeded
	Countered   bool // true if this player's Counter connected
}

// RoundResult reports what happened to both players in one round.
type RoundResult struct {
	P1, P2 PlayerRoundOutcome
	Winner *Player // nil if the match continues, or on a draw
	Draw   bool
}

// PlayRound resolves one round given both players' simultaneous choices.
func (m *Match) PlayRound(a1, a2 Action) RoundResult {
	var roll1, roll2 *AttackResult
	if a1.Kind == ActionAttack {
		r := a1.Weapon.Attack()
		roll1 = &r
	}
	if a2.Kind == ActionAttack {
		r := a2.Weapon.Attack()
		roll2 = &r
	}

	dmgToP2, p2Dodged, p2Countered := resolveIncoming(roll1, a2)
	dmgToP1, p1Dodged, p1Countered := resolveIncoming(roll2, a1)

	if p2Countered {
		dmgToP1 += int(float64(roll1.Damage) * counterReflect)
	}
	if p1Countered {
		dmgToP2 += int(float64(roll2.Damage) * counterReflect)
	}

	m.P1.TakeDamage(dmgToP1)
	m.P2.TakeDamage(dmgToP2)

	res := RoundResult{
		P1: PlayerRoundOutcome{Player: m.P1, Action: a1, AttackRoll: roll1, DamageTaken: dmgToP1, Dodged: p1Dodged, Countered: p1Countered},
		P2: PlayerRoundOutcome{Player: m.P2, Action: a2, AttackRoll: roll2, DamageTaken: dmgToP2, Dodged: p2Dodged, Countered: p2Countered},
	}

	p1Dead, p2Dead := !m.P1.Alive(), !m.P2.Alive()
	switch {
	case p1Dead && p2Dead:
		res.Draw = true
	case p2Dead:
		res.Winner = m.P1
		m.P1.Wins++
	case p1Dead:
		res.Winner = m.P2
		m.P2.Wins++
	}
	return res
}

// Over reports whether the match has ended (someone, or both, at 0 HP).
func (m *Match) Over() bool {
	return !m.P1.Alive() || !m.P2.Alive()
}

// resolveIncoming computes how much damage lands on a defender, given the
// attacker's roll (nil if the attacker didn't attack this round) and the
// defender's own chosen action.
func resolveIncoming(attackerRoll *AttackResult, defenderAction Action) (dmg int, dodged bool, countered bool) {
	if attackerRoll == nil {
		return 0, false, false
	}
	switch defenderAction.Kind {
	case ActionBlock:
		dmg = attackerRoll.Damage - blockReduction
		if dmg < 0 {
			dmg = 0
		}
		return dmg, false, false
	case ActionDodge:
		if rand.Float64() < dodgeChance {
			return 0, true, false
		}
		return attackerRoll.Damage, false, false
	case ActionCounter:
		return 0, false, true
	default: // ActionAttack: trading blows, no mitigation
		return attackerRoll.Damage, false, false
	}
}
