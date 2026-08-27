# TextFighterGo

A text-based fighting game, originally a Python passion project from high
school, now being rewritten in Go as a hands-on way to learn the language —
structs and methods instead of classes, explicit error handling instead of
exceptions, goroutines instead of threads — with an eye toward real
gameplay depth from here.

## Version 1.0

**Requirements:** 2 players, offline, either sharing one CLI session or
connected over TCP; also runnable in a Docker container.

What works today:

- Two players fight it out, round by round, both choosing an action
  simultaneously each round — neither sees the other's pick before
  committing to their own.
- Each player starts at 250 HP.
- Attack with one of three weapons — Ax, Sword, Club — each with its own
  damage range, crit chance, and miss chance.
- Or play defense instead of attacking: Block/Armor up (flat damage
  reduction), Dodge (chance to fully avoid the hit), or Counter (reflect
  the incoming attack back at your opponent if they attacked into it).
- Full match flow: live HP tracking, win detection, draws (both players
  falling the same round), and a "fight again?" prompt to replay.
- Runs directly with the Go toolchain, or as a Docker container.
- Network play: one person hosts on a port, a second player connects in
  over TCP as the opponent — no custom client needed, any plain TCP
  client (e.g. `nc`) works as the joining side.

Not built yet:
- A dedicated join client / real SSH transport (right now "joining" means
  pointing a generic TCP client like `nc` at the host).
- Classes/backgrounds (player archetypes affecting available actions).
- Skill tree/campaign backstory log, PVE bossing/routes.

See **Roadmap** below for what's next.

## Running it

**Local play** (two players, one terminal) — this is the default, no flags
needed:

With Go installed:

```sh
go run ./cmd/textfighter
```

*PREFERRED:*
With Docker instead (see `docs/docker-notes.md` for a from-scratch
explanation of what this is doing):

```sh
docker build -t textfighter .
docker run -it --rm textfighter
```

**Host a match over the network** instead, with `-listen`:

```sh
go run ./cmd/textfighter -listen=:9000
```

The host plays locally as Player 1 and waits (up to 2 minutes) for an
opponent to connect. A second player joins as Player 2 from any machine
with a plain TCP client — no custom client exists yet:

```sh
nc <host-address> 9000
```

## Project layout

| Path | What it is |
|---|---|
| `cmd/textfighter/` | The terminal entrypoint — wires the game logic up to stdin/stdout or a network connection, and hosts a match when `-listen` is passed. |
| `internal/game/` | Core game rules: `Player`, `Weapon`, `Action`, `Match`. No I/O, no networking — just the logic of a fight. |
| `internal/session/` | Wraps a player's input/output behind `io.Reader`/`io.Writer`, so the same game logic works over a local terminal or a network connection. |
| `docs/docker-notes.md` | A Docker primer, plus notes on this repo's `Dockerfile`/`docker-compose.yml`. |
| `Dockerfile`, `docker-compose.yml` | Build and run the game in a container. |
| `game.py` | The original Python prototype this is being ported from. Kept locally for reference; not part of the Go build. |

## Roadmap — what's planned next

Built in phases, roughly in this order:

1. ~~**Networking.**~~ ✅ Done — one person hosts over TCP, a second player
   connects in as the opponent, running match stats stay in the host's
   process.
2. ~~**Simultaneous combat.**~~ ✅ Done — both players choose an action
   each round at the same time (attack or a defensive option), resolved
   together instead of one-at-a-time.
3. **Classes/backgrounds (stretch).** Player archetypes that affect which
   actions are available or how they perform.
4. **SSH transport (stretch).** Swap the network layer for a real SSH
   server, so a second player can join with a plain `ssh` command instead
   of a generic TCP client like `nc`.

## Why this exists

A learning project: rewriting an old Python game in Go to build real,
hands-on fluency in the language — networking, concurrency, and Docker
along the way — ahead of a job interview where Go is a listed requirement.
