package game

// ActionKind is one of the choices a player can make each round.
type ActionKind int

const (
	ActionAttack ActionKind = iota
	ActionBlock
	ActionDodge
	ActionCounter
)

// Action is one player's choice for a round. Weapon only matters when Kind
// is ActionAttack.
type Action struct {
	Kind   ActionKind
	Weapon Weapon
}

// Tuning constants for the defensive options.
const (
	blockReduction = 20   // flat damage absorbed by Block/Armor
	dodgeChance    = 0.35 // chance a Dodge succeeds
	counterReflect = 1.0  // fraction of the incoming attack reflected back by a successful Counter
)
