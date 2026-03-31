# VHS tape file — records a terminal GIF of aitap in demo mode.
#
# Install VHS:   go install github.com/charmbracelet/vhs@latest
# Record:        vhs demo.tape
# Output:        demo.gif (in current directory)
#
# The recording uses `aitap --demo` so no real API calls are made.
# It shows realistic sample traffic from multiple providers.

# --- Output settings ---
Output demo.gif
Output demo.webm

# --- Terminal settings ---
Set FontFamily "JetBrains Mono"
Set FontSize 14
Set Width 1200
Set Height 700
Set Padding 20
Set Theme "Dracula"
Set TypingSpeed 50ms
Set WindowBar Colorful
Set WindowBarSize 40
Set CursorBlink false
Set PlaybackSpeed 0.75

# --- Build and run ---
Type "go build -o bin/aitap ./cmd/aitap && clear"
Enter
Sleep 1s

# Show version
Type "bin/aitap --version"
Enter
Sleep 1.5s

# Clear and start demo mode
Type "clear"
Enter
Sleep 500ms

Type "bin/aitap --demo"
Enter

# Wait for the TUI to render and calls to stream in
Sleep 2s

# First call appears (Anthropic)
Sleep 2s

# Second call (OpenAI)
Sleep 2s

# Third call (Google)
Sleep 2s

# Fourth call (OpenRouter)
Sleep 2s

# Fifth call (Ollama)
Sleep 2s

# Navigate up to first call
Type "g"
Sleep 500ms

# Expand the Anthropic call detail
Enter
Sleep 3s

# Close detail
Enter
Sleep 500ms

# Move down to OpenAI call
Type "j"
Sleep 300ms

# Expand OpenAI call detail
Enter
Sleep 3s

# Close detail
Enter
Sleep 500ms

# Move to Google call
Type "j"
Sleep 300ms

Enter
Sleep 3s

# Close
Enter
Sleep 500ms

# Jump to last call (Ollama)
Type "G"
Sleep 300ms

Enter
Sleep 3s

# Close and pause on the full list view
Enter
Sleep 500ms

# Jump back to top for a clean final frame
Type "g"
Sleep 2s

# Quit
Type "q"
Sleep 1s
