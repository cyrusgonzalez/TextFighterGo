package game

// Player is Go's stand-in for game.py's `class Player`. There's no class
// keyword and no __init__ — a struct just declares the fields (data only,
// no behavior attached yet), and behavior is added separately as methods
// below.
type Player struct {
	Name   string
	Health int
	Wins   int // running stat; meant to live host-side once networking exists
}

const startingHealth = 250

// NewPlayer is the constructor-equivalent. Go convention: name it
// New<Type>, have it return the struct (here, a pointer to it) fully
// initialized. This is doing the same job __init__ did in game.py, just as
// a plain function instead of magic tied to the type.
func NewPlayer(name string) *Player {
	return &Player{
		Name:   name,
		Health: startingHealth,
	}
}

// TakeDamage mutates the player, so it needs a *pointer* receiver
// (*Player). If this were `func (p Player) TakeDamage(...)` — a value
// receiver — Go would silently operate on a COPY of the player, and the
// health change would vanish the instant the method returns. This trips up
// almost everyone coming from Python/Java, where self/this always refers to
// the real object; Go makes you choose, explicitly, via the receiver type.
func (p *Player) TakeDamage(dmg int) {
	p.Health -= dmg
	if p.Health < 0 {
		p.Health = 0
	}
}

// Alive reports whether the player can still fight. A value receiver is
// fine here since Alive only reads state and returns a bool — no mutation,
// no need for a pointer.
func (p Player) Alive() bool {
	return p.Health > 0
}

// Reset puts a player back to full health for a rematch (game.py restarts
// this by recursively calling main() again — Go favors an explicit reset
// over recursive re-entry into the game loop).
func (p *Player) Reset() {
	p.Health = startingHealth
}
