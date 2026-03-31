# VHS tape file - records a terminal GIF of aitap in demo mode.
#
# Install VHS:   go install github.com/charmbracelet/vhs@latest
# Record:        vhs demo.tape
# Output:        demo.gif and demo.webm
#
# The recording uses `go run ./cmd/aitap --demo` so no real API calls are made.
# It also works from a fresh checkout without assuming a prebuilt binary on PATH.

Output demo.gif
Output demo.webm

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

# Warm the Go build cache and show the current version.
Type "go run ./cmd/aitap --version"
Enter
Sleep 2s

# Start demo mode.
Type "go run ./cmd/aitap --demo"
Enter

# Wait for the TUI to render and all sample calls to arrive.
Sleep 11s

# Inspect the first call.
Type "g"
Sleep 400ms
Enter
Sleep 3s
Enter
Sleep 500ms

# Move to the second call.
Type "j"
Sleep 300ms
Enter
Sleep 3s
Enter
Sleep 500ms

# Move to the third call.
Type "j"
Sleep 300ms
Enter
Sleep 3s
Enter
Sleep 500ms

# Jump to the last call.
Type "G"
Sleep 300ms
Enter
Sleep 3s
Enter
Sleep 500ms

# Return to the top for a clean last frame.
Type "g"
Sleep 2s

# Quit cleanly.
Type "q"
Sleep 1s
