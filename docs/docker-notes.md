# Docker notes — zero to "why this project uses it"

You said you have no Docker experience, so this is written from scratch.
Skim the concepts once, then read the annotated `Dockerfile`/`docker-compose.yml`
next to it — that's where it'll actually click.

## The core idea

A normal program depends on whatever's already on your machine: your Go
version, your OS, libraries you happened to install. That's the classic
"works on my machine" problem. Docker's answer: package the program *with*
everything it needs into one unit that runs identically anywhere Docker
runs — your laptop, a teammate's laptop, a CI server, a cloud VM.

Three terms, used precisely:

- **Image** — a read-only *blueprint*: your app + its OS-level dependencies,
  built in layers, described by a `Dockerfile`. Think "class."
- **Container** — a *running instance* of an image. Think "object" — you
  can start several containers from the same image, like instantiating a
  class multiple times, each with its own process/memory but the same code.
- **Dockerfile** — the recipe: a list of steps ("start from this base,
  copy these files, run this build command") that produces an image.

That image/container split maps cleanly onto something you already know:
image is to container as struct type is to struct value, or class is to
instance. The image doesn't "run" — you run a container *from* it.

## Why this project uses it (and why the job listing cares)

1. **Reproducibility.** "I built it with Go 1.26 and it needs to keep
   working the same way" — a Docker image pins that, permanently.
2. **This is how you'll run Kafka locally.** You are not going to
   hand-install a Kafka broker on your machine. You'll run
   `docker compose up` and get a real broker in seconds, throw it away, and
   do it again. Same story for Flink later.
3. **It's how the pieces of a multi-service system (game server + Kafka +
   Flink) get wired together and started as one unit** — that's what
   `docker-compose.yml` is for, see below.
4. **Interview relevance.** A job listing that names Docker alongside
   Kafka/Flink almost certainly means "can you containerize a Go service
   and reason about a docker-compose stack," not "have you memorized every
   flag." What's here covers that.

## Dockerfile, annotated

Open `Dockerfile` in this repo alongside this. It uses a **multi-stage
build** — a pattern that matters *specifically* for Go, so it's worth
understanding why:

- **Stage 1 ("builder")**: starts from a full `golang` image (has the Go
  toolchain, ~800MB), copies your source in, and runs `go build`. Go
  compiles to a single **static binary** — no runtime, no interpreter, no
  dependency list to ship alongside it. That's a real advantage over, say,
  Python: there's no `pip install` step that has to happen inside the final
  image.
- **Stage 2 (final)**: starts from a *tiny* base image (`alpine`, ~5MB) and
  copies *only the compiled binary* out of stage 1. The entire Go toolchain
  and source tree get thrown away. The result: a final image that's tens of
  megabytes instead of nearly a gigabyte.

This two-stage shape — "build fat, ship thin" — is close to a default
expectation for a Go service in industry, so it's worth being able to
explain out loud, not just paste.

## docker-compose, annotated

A single `Dockerfile`/image describes *one* program. Real systems are
usually several programs that need to talk to each other — here, eventually:
the game server, a Kafka broker, maybe Flink. `docker-compose.yml` describes
a *group* of containers as one declarative file, and `docker compose up`
starts all of them, on a shared virtual network, with one command.

Right now the compose file only runs the game itself — the Kafka block is
present but commented out, ready to switch on once we get to that phase.
When we do, note we'll use Kafka's **KRaft mode** (no separate Zookeeper
container) — that's the current recommended setup, not the older
Kafka+Zookeeper pairing you'll see in a lot of outdated tutorials.

## Commands you'll actually use

```sh
# Build the image from the Dockerfile in this directory, tag it "textfighter"
docker build -t textfighter .

# Run a container from that image, attached to your terminal (-it),
# so stdin/stdout work for our interactive game
docker run -it --rm textfighter

# Start everything defined in docker-compose.yml
docker compose up

# Tear it down
docker compose down

# See what's running / peek at logs of a running container
docker ps
docker logs <container id or name>

# Get a shell inside a running container (handy for debugging)
docker exec -it <container id or name> sh
```

## Installing Docker

Docker isn't installed on this machine yet. Get **Docker Desktop**
(macOS) from docker.com — it gives you the `docker` CLI plus the
background engine everything above talks to. Once it's installed, `docker
build -t textfighter .` from the repo root is the thing to try first.
