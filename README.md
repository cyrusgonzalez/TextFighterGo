# TextFighterGo

A text-based, turn-based fighting game, originally a Python passion project
from high school, now being rewritten in Go as a hands-on way to learn the
language — structs and methods instead of classes, explicit error handling
instead of exceptions, goroutines instead of threads — with an eye toward a
larger, more realistic system (networked play, an event pipeline, containerization).

## Version 1.0

**Requirements:** 2 players, local, in one CLI session; also runnable in a
Docker container.

What works today:

- Two players take turns fighting at the same terminal.
- Each player starts at 250 HP.
- Three weapons — Ax, Sword, Club — each with its own damage range, crit
  chance, and miss chance.
- A full match: random starting turn, turn-by-turn attacks, live HP
  tracking, win detection, and a "fight again?" prompt to replay.
- Runs directly with the Go toolchain, or as a Docker container.

Not built yet: 
- Any of it over a network. Right now both players share one
keyboard. See **Roadmap** below.
- _Armor/Block under construction_
- Skill tree/campaign backstory log
- PVE bossing/routes


## Running it

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

## Project layout

| Path | What it is |
|---|---|
| `cmd/textfighter/` | The terminal entrypoint — wires the game logic up to stdin/stdout. |
| `internal/game/` | Core game rules: `Player`, `Weapon`, `Match`. No I/O, no networking — just the logic of a fight. |
| `internal/producer/`, `internal/consumer/` | Placeholders for a future Kafka event pipeline. Not wired up yet. |
| `docs/docker-notes.md` | A Docker primer, plus notes on this repo's `Dockerfile`/`docker-compose.yml`. |
| `Dockerfile`, `docker-compose.yml` | Build and run the game in a container. |
| `game.py` | The original Python prototype this is being ported from. Kept locally for reference; not part of the Go build. |

## Roadmap — what's planned next

Built in phases, roughly in this order:

1. **Networking.** Turn the shared-terminal loop into a real client/server
   game — one person hosts, a second player connects over the network, and
   running match stats are kept on the host's side.
2. **SSH transport (stretch).** Swap the network layer for a real SSH
   server, so a second player can join with a plain `ssh` command instead
   of a custom client.
3. **Kafka.** Publish real game events (attacks, hits, match results) to a
   Kafka topic instead of just printing them, and consume them for
   stats/logging via the `internal/producer` / `internal/consumer`
   packages.
4. **Flink.** Stream-process those Kafka events for things like live
   aggregated stats across matches.
5. **Full container stack.** Game server + Kafka + (eventually) Flink,
   brought up together with a single `docker compose up`.

None of items 1–5 exist yet — `internal/producer` and `internal/consumer`
are currently empty package stubs, and there is no networking of any kind.

## Why this exists

A learning project: rewriting an old Python game in Go to build real,
hands-on fluency in the language — and in Kafka, Flink, and Docker along
the way — ahead of a job interview where those are listed skills.
