# 🎯 Random Pétanque Draw

A small Go project to generate a smart random draw for pétanque tournaments.  
The goal is to maximize the variety of encounters between players across 4 games, minimizing **collisions** (i.e. players facing or teaming up with the same people too often).

---

## 🚀 How It Works

### Requirements

- Go ≥ 1.24.3  
- (Optional) [Nix](https://nixos.org/) for a reproducible development environment

### Run a Simulation

Edit the parameters in `cmd/draw-cli/main.go`, then run:

```sh
go run cmd/draw-cli/main.go
```

### Output

The program generates a draw for 4 rounds with a variable number of players (e.g., 17, 19, 23).  
Each draw can be exported as a PDF

---

## 🏁 Context

On the first Friday of each month, a pétanque tournament is held in the village where I grew up.  
Players register individually and are assigned to **4 games** in the evening. Teams are formed randomly for each game using different formats:

- **Doublette**: 2 vs 2  
- **Triplette**: 3 vs 3  
- **Mixed 3v2**: 3 vs 2, with adjusted number of balls per player

Each player is assigned a number. The draw produces the teams and matches for all 4 games.  
💬 Over time, players noticed that they often end up playing with or against the same people multiple times.  
This project aims to reduce such occurrences and improve the tournament experience.

---

## 🧠 Theoretical Analysis

This is a classic case of a **combinatorial optimization** problem.

### Goal

Minimize **collisions**, meaning:

- A player teams up with the **same teammate** more than once
- A player faces the **same opponent** multiple times

### Constraints

- Total number of players `N` is known
- Each player must play **4 games**
- Each game includes **4 to 6 players**
- Only one match per round may use a **3v2 format**
- Preference is given to **2v2 games** for pacing

### Two Implemented Approaches

#### 🪓 Branch and Bound

- Exhaustive approach: explores all possible combinations
- Can find the optimal solution (0 collisions), but…
- Has **factorial complexity**, which becomes infeasible for `N > 12`
- Optimizations include pruning worse-than-current branches

#### ⚡ Greedy Algorithm

- Heuristic approach: builds lineups randomly several times, trying to minimize collisions at each step
- Much faster
- Does **not guarantee the best** solution, but yields good results in practice
- This is the default method used in the project

## 🏌️ Similarity to the Social Golfer Problem

This project is closely related to the well-known **Social Golfer Problem** in combinatorial optimization.

### What is the Social Golfer Problem?

The Social Golfer Problem asks:

> "Given *g* groups of *s* golfers who play together once a week for *w* weeks, can you arrange the groups so that no two golfers play in the same group more than once?"

This is strikingly similar to the goals of this pétanque draw:

- Players participate in **multiple games**
- We want to **minimize repeated encounters**
- Group sizes are **limited and varied**
- It’s a challenge of **fair scheduling** and **collision avoidance**

Although the constraints differ (pétanque allows flexible team formats like 3v2 or 3v3), the core idea of **optimally distributing players across matches** is shared.

Understanding and referencing the Social Golfer Problem helps when researching algorithms or benchmarking this solution.

---

## 📦 Project Structure

```
.
├── cmd/
│   └── draw-cli/         # CLI entrypoint
├── draw/                 # Core draw logic
├── fonts/                # PDF fonts
├── tournament/           # Data models and structures
├── utils/                # Utility functions
├── go.mod
└── README.md
```

## 📚 Ressources

- [Pétanque](https://en.wikipedia.org/wiki/P%C3%A9tanque)
- [Combinatorial optimization](https://en.wikipedia.org/wiki/Combinatorial_optimization)
- [Social golfer problem](https://en.wikipedia.org/wiki/Social_golfer_problem)

