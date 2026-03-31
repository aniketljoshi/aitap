package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/aniketljoshi/aitap/internal/export"
	"github.com/aniketljoshi/aitap/internal/model"
	"github.com/aniketljoshi/aitap/internal/tui"
)

var (
	version = "dev"
	commit  = "none"
)

func main() {
	port := flag.Int("port", 9119, "Proxy listen port")
	exportPath := flag.String("export", "", "Auto-export session to JSONL on exit")
	redactFlag := flag.Bool("redact", false, "Redact secrets in exports")
	filterProvider := flag.String("filter", "", "Only show calls to this provider (openai, anthropic, google, openrouter, ollama)")
	demoMode := flag.Bool("demo", false, "Launch the TUI with realistic sample traffic instead of starting the proxy")
	showVersion := flag.Bool("version", false, "Show version")
	flag.Parse()

	if *showVersion {
		fmt.Printf("aitap %s (%s)\n", version, commit)
		os.Exit(0)
	}

	session := model.NewSession()
	callChan := make(chan *model.Call, 100)

	if *demoMode {
		startDemoFeed(callChan, *filterProvider)
	} else {
		// Start proxy in background
		go func() {
			if err := startProxy(*port, callChan, *filterProvider); err != nil {
				log.Fatalf("Proxy error: %v", err)
			}
		}()
	}

	// Start TUI
	p := tea.NewProgram(
		tui.New(session, callChan, *port),
		tea.WithAltScreen(),
	)

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		if *exportPath != "" {
			_ = export.ToJSONL(session, *exportPath, *redactFlag)
			fmt.Fprintf(os.Stderr, "\nExported %d calls to %s\n", len(session.Calls), *exportPath)
		}
		p.Quit()
	}()

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}

	// Export on normal exit too
	if *exportPath != "" && len(session.Calls) > 0 {
		if err := export.ToJSONL(session, *exportPath, *redactFlag); err != nil {
			fmt.Fprintf(os.Stderr, "Export error: %v\n", err)
		} else {
			fmt.Fprintf(os.Stderr, "Exported %d calls to %s\n", len(session.Calls), *exportPath)
		}
	}
}
