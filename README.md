# The Stars, Like Dust (Autonomous Simulation Game)

An exploratory, autonomous narrative simulation game in Go based on Isaac Asimov's novel *The Stars, Like Dust* (Empire Series prelude).

The game is designed for **autonomous execution** by an AI agent—modeling the complex interstellar intrigue between the autocratic Tyranni Empire, the Directorate of Rhodia, the militant Autarchy of Lingane, and the rebellious worlds of the Nebular Kingdoms.

---

## Key Features

1. **Multi-POV Character Agency**:
   Play or simulate from the Point of View (POV) of any major novel character:
   - **Biron Farrill**: Son of the executed Rancher of Widemos, seeking justice and the secret rebellion world.
   - **Artemisia Hinriad**: Rhodian noblewoman fighting against arranged Tyranni marriage and subjugation.
   - **Gillbret Hinriad**: Eccentric inventor/visisonor virtuoso driven to vindicate his vision of the Horsehead Nebula.
   - **Sander Jonti (Autarch of Lingane)**: Cunning conspirator seeking to crown himself ruler of the Nebular Kingdoms.
   - **Simok Aratap**: The quiet, calculating Tyranni Commissioner hunting down dissidents via sub-ether tracking.
   - **Hinrik of Rhodia**: The seemingly foolish Director masking a 30-year subterranean resistance strategy.

2. **Exploratory "Metaopinions"**:
   Beyond canonical events, the simulation tracks philosophical and political "metaopinions" regarding galactic governance, feudalism vs democracy, and the enduring power of constitutional ideas.

3. **Subspace Communication DSL (`comm/`)**:
   Includes an invented human-readable Domain Specific Language (**HCP** - Hyper-Comm Protocol) for interstellar dispatches and orders, with a Go parser, AST, and formatter in `pkg/dsl`.

4. **Embedded Relational Storage (`data/`)**:
   Persistent state is stored in an embedded SQLite database (`data/stars_universe.db`) via `modernc.org/sqlite`. Features temporary evaluation tables and permanent summarized sessions, guaranteed to stay strictly under 500 MB.

5. **Episodic Bounded Execution**:
   Executes in finite, bounded dramatic turns per episode, eliminating infinite loops and preventing token exhaustion.

6. **1-to-2 A4 Page Summary Deliverable**:
   At the end of each simulation episode, outputs a structured 1-to-2 A4 page narrative report containing psychological profiles, turn breakdowns, metaopinions manifested, reflective thoughts, and targeted queries for the user.

7. **Multi-Agent Sub-Agent Delegation (`--subagents`)**:
   Delegates secondary characters encountered during the simulation to autonomous sub-agents with randomized or archetype-tailored model tiers (`flash_lite`, `flash`, `pro`) and cognitive effort levels, conducting dynamic 1-round HCP dispatch standoffs in `comm/`.

---

## Project Structure

```text
The_Stars_Like_Dust/
├── cmd/
│   └── stars/
│       └── main.go                  # CLI entrypoint supporting POV selection and episode play
├── internal/
│   ├── engine/                      # Game turn director, state machine, episodic resolution
│   ├── persona/                     # Character loader, metaopinion evaluator, goal weighting
│   ├── comm/                        # Subspace transmission manager, dispatch queue
│   └── storage/                     # Embedded pure-Go SQLite storage (data/stars_universe.db)
├── pkg/
│   ├── dsl/                         # HCP DSL parser, tokenizer, AST, and formatter
│   └── summary/                     # 1-to-2 page A4 narrative report generator
├── metadata/                        # Human-readable JSON persona and sector definitions
│   ├── biron_farrill.json
│   ├── artemisia_hinriad.json
│   ├── gillbret_hinriad.json
│   ├── sander_jonti.json
│   ├── simok_aratap.json
│   ├── hinrik_hinriad.json
│   └── sector_map.json
├── comm/                            # Interstellar subspace dispatches in .hcp format
│   ├── dispatch_001_widemos_execution.hcp
│   ├── dispatch_002_earth_radiation_probe.hcp
│   ├── dispatch_003_lingane_shadow_directive.hcp
│   ├── dispatch_004_aratap_surveillance.hcp
│   └── dispatch_005_rebellion_coordinates.hcp
├── data/                            # Persistent embedded database (capped < 500 MB)
│   └── stars_universe.db
├── Makefile
├── README.md
├── USAGE.md
├── ARCHITECTURE.md
├── WORKFLOW.md
├── brief.md
├── prospective.md
└── LICENSE
```

---

## Documentation

- **[brief.md](brief.md)**: Literary synopsis of Asimov's novel, timeline, political factions, and thematic analysis.
- **[prospective.md](prospective.md)**: Software design perspective, engineering philosophy, multi-POV agency, and capacity control.
- **[USAGE.md](USAGE.md)**: CLI flags, character POV listing, execution examples, and output modes.
- **[ARCHITECTURE.md](ARCHITECTURE.md)**: Detailed system architecture, component diagrams, state machine, and data flow.
- **[WORKFLOW.md](WORKFLOW.md)**: Development lifecycle, Plane project synchronization, adding characters/dispatches, and verification rules.

---

## Quick Start

### Build & Run
```bash
# Build the binary
make build

# List available character POVs
./stars --list-povs

# Run an autonomous episode from Biron Farrill's POV
./stars --pov biron_farrill

# Run an episode from Commissioner Aratap's POV (summary only)
./stars --pov simok_aratap --summary-only
```

### Running Tests
```bash
make test-coverage
make vet
```
