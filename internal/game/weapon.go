package game

import "math/rand/v2"

// Weapon describes one fighting style's damage profile.
type Weapon struct {
	Name       string
	MinDmg     int
	MaxDmg     int
	CritMin    int
	CritMax    int
	CritChance float64
	MissChance float64
}

var (
	Ax = Weapon{
		Name: "Ax", MinDmg: 20, MaxDmg: 50,
		CritMin: 45, CritMax: 75,
		CritChance: 0.05, MissChance: 0.15,
	}
	Sword = Weapon{
		Name: "Sword", MinDmg: 10, MaxDmg: 25,
		CritMin: 25, CritMax: 50,
		CritChance: 0.50, MissChance: 0.13,
	}
	Club = Weapon{
		Name: "Club", MinDmg: 15, MaxDmg: 40,
		CritMin: 40, CritMax: 60,
		CritChance: 0.15, MissChance: 0.10,
	}
)

// Weapons is a lookup table by menu number.
var Weapons = map[int]Weapon{
	1: Ax,
	2: Sword,
	3: Club,
}

// AttackResult reports what one attack roll did.
type AttackResult struct {
	Damage  int
	Outcome string // "hit", "crit", or "miss"
}

// Attack rolls a single attack for this weapon.
func (w Weapon) Attack() AttackResult {
	roll := rand.Float64()

	switch {
	case roll < w.MissChance:
		return AttackResult{Damage: 0, Outcome: "miss"}
	case roll < w.MissChance+w.CritChance:
		dmg := w.CritMin + rand.IntN(w.CritMax-w.CritMin+1)
		return AttackResult{Damage: dmg, Outcome: "crit"}
	default:
		dmg := w.MinDmg + rand.IntN(w.MaxDmg-w.MinDmg+1)
		return AttackResult{Damage: dmg, Outcome: "hit"}
	}
}
