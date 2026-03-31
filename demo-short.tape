# Short VHS tape — produces a tight ~12s GIF for social media.
#
# Usage:  vhs demo-short.tape
# Output: demo-short.gif

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

# Start demo mode (pre-built binary assumed)
Type "aitap --demo"
Enter

# Wait for all 5 calls to stream in
Sleep 12s

# Expand the first call
Type "g"
Sleep 300ms
Enter
Sleep 3s

# Quit for clean exit
Type "q"
Sleep 500ms
