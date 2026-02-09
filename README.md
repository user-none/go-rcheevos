# go-rcheevos

Go bindings for [rcheevos](https://github.com/RetroAchievements/rcheevos), the
RetroAchievements client library. This package enables Go-based emulators to
integrate with [RetroAchievements](https://retroachievements.org) for
achievement and leaderboard support.

## Features

- **Authentication** - Login with password or token
- **Game Identification** - Hash ROMs and identify games in the RetroAchievements database
- **Achievements** - Track, unlock, and display achievement progress
- **Leaderboards** - Submit scores and fetch leaderboard entries
- **Rich Presence** - Generate dynamic status messages
- **Multi-disc Support** - Handle disc changes for multi-disc games
- **Progress Serialization** - Save/restore achievement progress for savestates
- **Subsets** - Support for bonus and DLC achievement sets

## Requirements

- Go 1.20+
- CGO enabled
- C compiler (gcc, clang, etc.)

## Installation

```bash
go get github.com/user-none/go-rcheevos
```

The rcheevos C library is included as a git submodule and compiled automatically via CGO.

## Quick Start

```go
package main

import (
    "fmt"
    "net/http"
    "io"

    rcheevos "github.com/user-none/go-rcheevos"
)

func main() {
    // Create client with memory read and HTTP callbacks
    client := rcheevos.NewClient(
        // Memory read callback - return bytes read from emulator memory
        func(address uint32, buffer []byte) uint32 {
            // Read from your emulator's memory here
            return 0
        },
        // Server callback - make HTTP requests to RetroAchievements API
        func(req *rcheevos.ServerRequest) {
            go func() {
                resp, err := http.Get(req.URL)
                if err != nil {
                    req.Respond(nil, 0)
                    return
                }
                defer resp.Body.Close()
                body, _ := io.ReadAll(resp.Body)
                req.Respond(body, resp.StatusCode)
            }()
        },
    )
    defer client.Destroy()

    // Set up event handler for achievement notifications
    client.SetEventHandler(func(event *rcheevos.Event) {
        switch event.Type {
        case rcheevos.EventAchievementTriggered:
            fmt.Printf("Achievement Unlocked: %s\n", event.Achievement.Title)
        case rcheevos.EventGameCompleted:
            fmt.Println("All achievements unlocked!")
        }
    })

    // Login
    done := make(chan error, 1)
    client.LoginWithPassword("username", "password", func(result int, errMsg string) {
        if result != rcheevos.OK {
            done <- fmt.Errorf("login failed: %s", errMsg)
            return
        }
        done <- nil
    })
    if err := <-done; err != nil {
        panic(err)
    }

    // Hash and load game
    hash := rcheevos.HashFromFile(rcheevos.ConsoleNES, "/path/to/rom.nes")
    client.LoadGame(hash, func(result int, errMsg string) {
        // Game loaded, achievements are now active
    })

    // In your emulator's main loop, call DoFrame() every frame
    // client.DoFrame()
}
```

## API Overview

### Client Lifecycle

```go
client := rcheevos.NewClient(readMemoryFunc, serverCallFunc)
defer client.Destroy()
```

### Authentication

```go
client.LoginWithPassword(username, password, callback)
client.LoginWithToken(username, token, callback)  // Use stored token for re-login
client.Logout()
user := client.GetUser()  // Get current user info
```

### Game Loading

```go
// Hash a ROM file
hash := rcheevos.HashFromFile(consoleID, "/path/to/rom")
hash := rcheevos.HashFromBuffer(consoleID, romData)

// Load by hash
client.LoadGame(hash, callback)

// Or identify and load in one step
client.IdentifyAndLoadGame(consoleID, filePath, data, callback)

client.UnloadGame()
game := client.GetGame()  // Get current game info
```

### Achievements

```go
client.HasAchievements()
client.GetAchievement(id)
list := client.CreateAchievementList(category, grouping)
defer list.Destroy()
summary := client.GetUserGameSummary()
```

### Leaderboards

```go
client.HasLeaderboards()
client.GetLeaderboard(id)
list := client.CreateLeaderboardList(grouping)
defer list.Destroy()
client.FetchLeaderboardEntries(id, firstEntry, count, callback)
client.FetchLeaderboardEntriesAroundUser(id, count, callback)
```

### Frame Processing

```go
// Call every frame during gameplay
client.DoFrame()

// Call when emulator is paused
client.Idle()

// Reset achievement state (e.g., when user resets the game)
client.Reset()
```

### Rich Presence

```go
if client.HasRichPresence() {
    message := client.GetRichPresenceMessage()
}
```

### Progress Serialization

```go
// Save progress (include in savestates)
data, result := client.SerializeProgress()

// Restore progress (when loading savestates)
result := client.DeserializeProgress(data)
```

### Configuration

```go
client.SetHardcoreEnabled(true)
client.EnableLogging(rcheevos.LogLevelInfo, logFunc)
client.SetReadMemoryFunc(readFunc)
```

## Supported Consoles

The library supports all consoles available in RetroAchievements, including:

- Nintendo: NES, SNES, N64, Game Boy, GBC, GBA
- Sega: Master System, Genesis/Mega Drive, Game Gear, Sega CD, 32X
- Sony: PlayStation
- NEC: PC Engine/TurboGrafx-16
- And many more

Use the console constants (e.g., `rcheevos.ConsoleNES`, `rcheevos.ConsoleSNES`) when hashing ROMs.

## Building

```bash
# Build the library
make build

# Run tests
make test

# Build the example
make example

# Update rcheevos submodule to latest
make update-rcheevos
```

## Example

See the [example](example/main.go) directory for a complete working example that demonstrates login, game loading, and achievement display.

## API Coverage

This library wraps the `rc_client` API from rcheevos, which is the recommended
integration path for emulators. The bindings cover approximately **98%** of the
rc_client functionality.

### Not Supported

The following features are not currently exposed:

- **Custom hash callbacks** (`rc_client_set_hash_callbacks`) - For custom file
  format handling (ZIP, CHD, encrypted ROMs). The default file/CD readers work
  for standard ROM formats.
- **Achievement iterator** (`rc_client_get_next_achievement_info`) - For
  iterating achievements by bucket. Use `CreateAchievementList()` instead.
- **RAIntegration** - Windows-only developer tools for achievement creation.
  Not applicable to cross-platform Go applications.

These are advanced features that most emulator integrations don't require. The
core achievement/leaderboard gameplay loop is fully supported.

## Related Projects

- [rcheevos](https://github.com/RetroAchievements/rcheevos) - The underlying C library
- [RetroAchievements](https://retroachievements.org) - The achievement platform

## License

MIT License. See [LICENSE](LICENSE) for details.

The underlying [rcheevos](rcheevos/LICENSE) C library has its own license terms.
