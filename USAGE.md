# Usage Guide: The Stars, Like Dust

This guide covers common CLI commands, configuration options, and autonomous simulation workflows.

---

## 1. CLI Commands & Options

### Listing Available Character POVs
To see all available characters that can be simulated:
```bash
./stars --list-povs
```

**Output**:
```text
Available Character Points of View (POV):
  • biron_farrill      : Biron Farrill (Heir of Widemos) — Allegiance: Independent / Anti-Tyrann
  • artemisia_hinriad  : Artemisia Hinriad (Lady of Rhodia) — Allegiance: House of Hinriad / Covert Subversive
  • gillbret_hinriad   : Gillbret Hinriad (Cousin of the Directorate) — Allegiance: House of Hinriad / Romantic Visionary
  • sander_jonti       : Sander Jonti (Autarch of Lingane) — Allegiance: Autarchy of Lingane / Personal Ambition
  • simok_aratap       : Simok Aratap (Commissioner of the Tyranni Empire) — Allegiance: The Khanate of Tyrann
  • hinrik_hinriad     : Hinrik Hinriad (Director of Rhodia) — Allegiance: House of Hinriad / Deep Subversive Network
```

---

## 2. Running Autonomous Gameplay Episodes

### Playing from Biron Farrill's Perspective
```bash
./stars --pov biron_farrill
```
Runs the 5-turn dramatic sequence:
1. Escape from the dormitory trap on Earth
2. Infiltrating the Directorate of Rhodia
3. Seizure of the fast cruiser *Remembrance*
4. Rendezvous with the Autarch at Lingane
5. Confrontation in the Horsehead Nebula and revelation of the ancient Earth Constitution

### Playing from Commissioner Simok Aratap's Perspective
```bash
./stars --pov simok_aratap
```
Experience imperial counter-insurgency:
1. Strategic evaluation of the Widemos execution
2. Interrogation of Director Hinrik on Rhodia
3. Affixing the sub-ether tracking transponder to *Remembrance*
4. Trailing the conspirators into the Horsehead Nebula
5. Ambush, unmasking of Jonti's betrayal, and philosophical confrontation with constitutional democracy

### Summary-Only Output Mode
To omit scene logs and output strictly the 1-to-2 A4 page report deliverable:
```bash
./stars --pov artemisia_hinriad --summary-only
```

---

## 3. Dynamic Narrative Branching & Stances

You can modify a character's tactical posture to explore divergent storylines and alternate endings:

```bash
# Run with an aggressive posture (direct confrontation, blaster engagement)
./stars --pov biron_farrill --stance aggressive

# Run with a covert posture (stealth evasion, subterranean exfiltration)
./stars --pov biron_farrill --stance covert

# Run with a diplomatic posture (constitutional negotiation)
./stars --pov biron_farrill --stance diplomatic

# Run with an inquisitive posture (investigative probe and psychological detection)
./stars --pov biron_farrill --stance inquisitive
```

### Stance Modifiers
- `balanced` (default): Follows canon-accurate psychological weighting.
- `aggressive`: Prioritizes high-risk, direct-action options (+40 score boost), unlocking violent naval skirmishes and armed ceasefires.
- `covert`: Prioritizes stealth codes and sensor evasion (+40 score boost), unlocking subterranean archival exfiltration.
- `diplomatic`: Prioritizes negotiation and ideological pacts (+40 score boost).
- `inquisitive`: Prioritizes scientific astrogation and psychological probe chords (+40 score boost).

---

## 4. Stochastic Temperature Sampling (`--temp` & `--seed`)

To introduce probabilistic variability into a character's choices rather than pure deterministic optimization:

```bash
# Run with stochastic Boltzmann sampling (temperature = 0.8)
./stars --pov biron_farrill --temp 0.8

# Run with a reproducible PRNG seed for stochastic replay
./stars --pov biron_farrill --temp 0.8 --seed 4242
```
- At `--temp 0.0` (default): The decision engine is strictly deterministic (classic argmax).
- At `--temp 0.5 - 1.2`: High-probability actions are preferred, but second-best or desperate alternatives may be triggered.
- At `--temp > 1.5`: Decisions become highly volatile, testing unpredictable narrative branches.

---

## 5. Database Persistence & Inspection

The simulation state is automatically persisted in `data/stars_universe.db`.

### Key Tables
- `characters`: Character profiles, allegiance, and home world.
- `metaopinions`: Canonical and exploratory opinions with certainty scores.
- `game_sessions`: Session UUID, POV character, status, and outcome.
- `session_turns`: Turn-by-turn logs, choices, dialogues, and outcomes.
- `temp_decision_evaluations`: Temporary decision evaluations (purged after each run to prevent growth).
- `narrative_summaries`: Persistent A4 reports, thoughts, and queries.

### Size Limit Compliance (< 500 MB)
The embedded database employs WAL mode and automatic compaction (`PRAGMA vacuum`). Typical storage consumption across 100 runs is under 2 MB, strictly adhering to the 500 MB ceiling.

---

## 4. Hyper-Comm Protocol (HCP) DSL

Subspace transmissions in `comm/` follow the HCP syntax:

```hcp
DISPATCH HCP-820-004
ENCRYPTION: TYRANN-CIPHER-7
ORIGIN: Commissioner-Simok-Aratap / Flagship-Core
TARGET: Imperial-Security-Council
STARDATE: 820-G.E.138
PRIORITY: STRATEGIC
STATUS: EYES-ONLY
PAYLOAD:
  EVENT: "Vessel Remembrance successfully breached Rhodia planetary blockade."
  INTEL: "Miniaturized sub-ether transponder activated on auxiliary engine hull."
  DIRECTIVE: "Do not fire upon Remembrance. Allow conspirators unhindered passage."
AUTHENTICATION: SIG-ARATAP-COMMISSIONER-991
END_DISPATCH
```

New dispatches can be created in `comm/` and will be loaded and indexed automatically by the engine.
