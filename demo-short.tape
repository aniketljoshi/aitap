# Short VHS tape - produces a tight GIF for social media.
#
# Usage:  vhs demo-short.tape
# Output: demo-short.gif
#
# This tape runs from a fresh checkout and does not require a prebuilt binary.

Output demo-short.gif

Set FontFamily "JetBrains Mono"
Set FontSize 15
Set Width 1100
Set Height 600
Set Padding 16
Set Theme "Dracula"
Set WindowBar Colorful
Set WindowBarSize 40
Set CursorBlink false
Set PlaybackSpeed 1

# Start demo mode directly.
Type "go run ./cmd/aitap --demo"
Enter

# Wait for the sample calls to stream in.
Sleep 11s

# Expand the first call.
Type "g"
Sleep 300ms
Enter
Sleep 3s

# Quit cleanly.
Type "q"
Sleep 500ms
