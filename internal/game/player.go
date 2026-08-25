package game

// Player is Go's stand-in for game.py's `class Player`.
type Player struct {
	Name   string
	Health int
	Wins   int
}

const startingHealth = 250

// NewPlayer is the constructor-equivalent.
func NewPlayer(name string) *Player {
	return &Player{
		Name:   name,
		Health: startingHealth,
	}
}

// TakeDamage mutates the player, so it needs a pointer receiver.
func (p *Player) TakeDamage(dmg int) {
	p.Health -= dmg
	if p.Health < 0 {
		p.Health = 0
	}
}

// Alive reports whether the player can still fight.
func (p Player) Alive() bool {
	return p.Health > 0
}

// Reset puts a player back to full health for a rematch.
func (p *Player) Reset() {
	p.Health = startingHealth
}
