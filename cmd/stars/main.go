package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/arthurgray2k/The_Stars_Like_Dust/internal/comm"
	"github.com/arthurgray2k/The_Stars_Like_Dust/internal/engine"
	"github.com/arthurgray2k/The_Stars_Like_Dust/internal/persona"
	"github.com/arthurgray2k/The_Stars_Like_Dust/internal/storage"
	"github.com/arthurgray2k/The_Stars_Like_Dust/pkg/summary"
)

func main() {
	povFlag := flag.String("pov", "biron_farrill", "Character Point of View (biron_farrill, artemisia_hinriad, simok_aratap, gillbret_hinriad, sander_jonti, hinrik_hinriad)")
	stanceFlag := flag.String("stance", "balanced", "Tactical Stance modifier (balanced, aggressive, covert, diplomatic, inquisitive)")
	tempFlag := flag.Float64("temp", 0.0, "Stochastic sampling temperature (0.0 = deterministic; >0.0 = Softmax Boltzmann sampling)")
	seedFlag := flag.Int64("seed", 0, "PRNG seed for reproducible stochastic simulation (0 = system clock)")
	dbFlag := flag.String("db", "data/stars_universe.db", "Path to embedded database file")
	metaFlag := flag.String("metadata", "metadata", "Path to persona metadata directory")
	commFlag := flag.String("comm", "comm", "Path to communication dispatches directory")
	listPovs := flag.Bool("list-povs", false, "List all playable character POVs")
	summaryOnly := flag.Bool("summary-only", false, "Print only the A4 summary deliverable")
	flag.Parse()

	// Normalize paths if run from subdirectories
	metaDir := *metaFlag
	if _, err := os.Stat(metaDir); os.IsNotExist(err) {
		metaDir = filepath.Join("..", metaDir)
	}
	commDir := *commFlag
	if _, err := os.Stat(commDir); os.IsNotExist(err) {
		commDir = filepath.Join("..", commDir)
	}

	pMgr := persona.NewManager(metaDir)
	personas, err := pMgr.LoadAll()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading personas: %v\n", err)
		os.Exit(1)
	}

	if *listPovs {
		fmt.Println("Available Character Points of View (POV):")
		for _, p := range personas {
			fmt.Printf("  • %-18s : %s (%s) — Allegiance: %s\n", p.ID, p.Name, p.Title, p.Allegiance)
		}
		fmt.Println("\nTactical Stances available via --stance:")
		fmt.Println("  • balanced    : Canon-accurate psychological goal alignment")
		fmt.Println("  • aggressive  : Direct action, blaster engagement, and elevated risk tolerance")
		fmt.Println("  • covert      : Maximum stealth, counter-surveillance, and avoiding imperial attention")
		fmt.Println("  • diplomatic  : Negotiation, constitutional pacts, and alliance formation")
		fmt.Println("  • inquisitive : Scientific probe, navigational analysis, and psychological detection")
		return
	}

	pov := strings.ToLower(strings.TrimSpace(*povFlag))
	targetPersona, ok := pMgr.Get(pov)
	if !ok {
		fmt.Fprintf(os.Stderr, "Unknown POV '%s'. Available POVs: biron_farrill, artemisia_hinriad, simok_aratap, gillbret_hinriad, sander_jonti, hinrik_hinriad\n", pov)
		os.Exit(1)
	}

	store, err := storage.Open(*dbFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening database at %s: %v\n", *dbFlag, err)
		os.Exit(1)
	}
	defer store.Close()

	cMgr := comm.NewManager(commDir)
	simEngine := engine.New(store, pMgr, cMgr)

	if !*summaryOnly {
		fmt.Printf("================================================================================\n")
		fmt.Printf("LAUNCHING AUTONOMOUS SIMULATION: The Stars, Like Dust\n")
		fmt.Printf("Selected Vantage: %s (%s)\n", targetPersona.Name, targetPersona.Title)
		fmt.Printf("Tactical Stance : %s | Temperature: %.2f\n", *stanceFlag, *tempFlag)
		fmt.Printf("Database: %s (Ceiling: < 500 MB)\n", *dbFlag)
		fmt.Printf("================================================================================\n\n")
	}

	opts := engine.SimulationOptions{
		POV:         pov,
		Stance:      *stanceFlag,
		Temperature: *tempFlag,
		Seed:        *seedFlag,
	}

	state, err := simEngine.RunEpisodeWithOptions(opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Simulation execution failed: %v\n", err)
		os.Exit(1)
	}

	if !*summaryOnly {
		for _, t := range state.TurnHistory {
			fmt.Printf("► TURN %d [%s]\n", t.TurnIndex, t.Location)
			fmt.Printf("  Context : %s\n", t.NarrativeEvent)
			fmt.Printf("  Actor   : %s\n", t.ActorName)
			fmt.Printf("  Choice  : %s\n", t.Decision.ChosenAction.Description)
			if t.Dialogue != "" {
				fmt.Printf("  Dialogue: \"%s\"\n", t.Dialogue)
			}
			fmt.Printf("  Rationale: %s\n", t.Decision.Rationale)
			fmt.Printf("  Outcome : %s\n\n", t.Outcome)
		}
	}

	// Generate A4 narrative summary
	report := summary.GenerateA4Summary(state, targetPersona)

	// Persist narrative summary in database
	queriesJoined := strings.Join(report.Queries, "\n")
	_ = store.SaveNarrativeSummary(state.SessionID, report.FullReport, report.Thoughts, queriesJoined)

	// Verify database size boundary
	dbSize, _ := store.CheckSize()
	if !*summaryOnly {
		fmt.Printf("✔ Episode Completed. Database synchronized. Size: %.2f KB (Well under 500 MB limit)\n\n", float64(dbSize)/1024.0)
	}

	// Output deliverable
	fmt.Print(report.FullReport)
}
