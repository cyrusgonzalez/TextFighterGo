package game

import "math/rand/v2"

// Weapon describes one fighting style's damage profile. In game.py this was
// three separate copy-pasted functions (axAtt, swordAtt, clubAtt) with
// near-identical bodies. Go has no classes/inheritance to reach for here —
// instead a weapon is just data (a Weapon value), and every weapon shares
// one Attack method below. "Prefer a data table + shared behavior over
// duplicated per-case functions" is a very Go-flavored instinct.
type Weapon struct {
	Name       string
	MinDmg     int
	MaxDmg     int
	CritMin    int
	CritMax    int
	CritChance float64 // probability 0.0-1.0 of a critical hit
	MissChance float64 // probability 0.0-1.0 of a total miss
}

// Package-level weapon definitions, ported from game.py's header comments
// (e.g. "ax damage does 20-50 DPS (5% CRT/15% miss)"). A `var (...)` block
// is Go's way of grouping package-scoped values — there's no "static
// final"; exported identifiers (capitalized) are just visible to importers.
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

// Weapons is a lookup table standing in for game.py's weaponList + integer
// index (weapon = weaponList[0], etc). A map[int]Weapon is the direct
// analogue: menu code below indexes into it by the number the player types.
var Weapons = map[int]Weapon{
	1: Ax,
	2: Sword,
	3: Club,
}

// AttackResult reports what one attack roll did. game.py's damage functions
// both computed a number AND printed to the screen inside the same
// function. Go convention leans the other way: keep computation and I/O
// separate, so the same Attack() works whether the caller is printing to a
// local terminal (today) or writing to a network connection (later).
type AttackResult struct {
	Damage  int
	Outcome string // "hit", "crit", or "miss"
}

// Attack rolls a single attack for this weapon.
//
// Receiver note: `(w Weapon)` is a *value* receiver — Attack reads w's
// fields but never needs to modify the weapon itself, so Go convention says
// pass it by value (a cheap copy of a small struct). Compare this to
// Player.TakeDamage in player.go, which uses a *pointer* receiver because
// it must mutate the caller's actual Player.
//
// math/rand/v2 (stdlib since Go 1.22) replaces Python's `random` module:
// rand.Float64() gives [0,1), rand.IntN(n) gives [0,n).
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
