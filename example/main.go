package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"

	rcheevos "github.com/user-none/go-rcheevos"
)

// Example SMS emulator memory (in a real emulator this would be the actual memory)
type SMSMemory struct {
	ram  [8192]byte  // 8KB main RAM
	vram [16384]byte // 16KB VRAM
	mu   sync.RWMutex
}

func (m *SMSMemory) Read(address uint32, buffer []byte) uint32 {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// SMS memory map for RetroAchievements:
	// 0x0000-0x1FFF: Main RAM (8KB, mirrored at 0xC000-0xDFFF)
	// For SMS, rcheevos maps system RAM starting at address 0
	if address < 0x2000 {
		n := copy(buffer, m.ram[address:])
		return uint32(n)
	}

	return 0
}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: example <username> <password> <rom_path>")
		os.Exit(1)
	}

	username := os.Args[1]
	password := os.Args[2]
	romPath := os.Args[3]

	// Create emulator memory
	memory := &SMSMemory{}

	// Create the rcheevos client
	client := rcheevos.NewClient(
		// Memory read callback
		func(address uint32, buffer []byte) uint32 {
			return memory.Read(address, buffer)
		},
		// Server call callback - make HTTP requests to RetroAchievements API
		func(req *rcheevos.ServerRequest) {
			go func() {
				var resp *http.Response
				var err error

				if req.PostData != "" {
					resp, err = http.Post(req.URL, req.ContentType, bytes.NewBufferString(req.PostData))
				} else {
					resp, err = http.Get(req.URL)
				}

				if err != nil {
					log.Printf("HTTP error: %v", err)
					req.Respond(nil, 0)
					return
				}
				defer resp.Body.Close()

				body, err := io.ReadAll(resp.Body)
				if err != nil {
					log.Printf("Read error: %v", err)
					req.Respond(nil, 0)
					return
				}

				req.Respond(body, resp.StatusCode)
			}()
		},
	)
	defer client.Destroy()

	// Enable logging
	client.EnableLogging(rcheevos.LogLevelInfo, func(message string) {
		log.Printf("[rcheevos] %s", message)
	})

	// Set event handler for achievements
	client.SetEventHandler(func(event *rcheevos.Event) {
		switch event.Type {
		case rcheevos.EventAchievementTriggered:
			if event.Achievement != nil {
				fmt.Printf("🏆 Achievement Unlocked: %s (%d points)\n",
					event.Achievement.Title, event.Achievement.Points)
				fmt.Printf("   %s\n", event.Achievement.Description)
			}
		case rcheevos.EventLeaderboardStarted:
			if event.Leaderboard != nil {
				fmt.Printf("📊 Leaderboard Started: %s\n", event.Leaderboard.Title)
			}
		case rcheevos.EventLeaderboardSubmitted:
			if event.Leaderboard != nil {
				fmt.Printf("📊 Leaderboard Submitted: %s - %s\n",
					event.Leaderboard.Title, event.Leaderboard.TrackerValue)
			}
		case rcheevos.EventGameCompleted:
			fmt.Println("🎉 Game Mastered! All achievements unlocked!")
		case rcheevos.EventReset:
			fmt.Println("⚠️  System reset requested (hardcore mode enabled)")
		}
	})

	// Login
	var loginDone = make(chan error, 1)
	client.LoginWithPassword(username, password, func(result int, errorMessage string) {
		if result != rcheevos.OK {
			loginDone <- fmt.Errorf("login failed: %s", errorMessage)
			return
		}
		loginDone <- nil
	})

	if err := <-loginDone; err != nil {
		log.Fatal(err)
	}

	user := client.GetUser()
	if user != nil {
		fmt.Printf("Logged in as: %s (Score: %d)\n", user.DisplayName, user.Score)
	}

	// Hash the ROM
	hash := rcheevos.HashFromFile(rcheevos.ConsoleMasterSystem, romPath)
	if hash == "" {
		log.Fatal("Failed to hash ROM")
	}
	fmt.Printf("ROM Hash: %s\n", hash)

	// Load the game
	var loadDone = make(chan error, 1)
	client.LoadGame(hash, func(result int, errorMessage string) {
		if result != rcheevos.OK {
			loadDone <- fmt.Errorf("load game failed: %s", errorMessage)
			return
		}
		loadDone <- nil
	})

	if err := <-loadDone; err != nil {
		log.Fatal(err)
	}

	game := client.GetGame()
	if game != nil {
		fmt.Printf("Game: %s (ID: %d)\n", game.Title, game.ID)
		fmt.Printf("Console: %s\n", rcheevos.ConsoleName(game.ConsoleID))
	}

	// Display game summary
	summary := client.GetUserGameSummary()
	if summary != nil {
		fmt.Printf("Progress: %d/%d achievements (%d/%d points)\n",
			summary.NumUnlockedAchievements, summary.NumCoreAchievements,
			summary.PointsUnlocked, summary.PointsCore)
	}

	// List all achievements
	if client.HasAchievements() {
		fmt.Println("\n--- Achievements ---")
		achievementList := client.CreateAchievementList(rcheevos.AchievementCategoryCore, rcheevos.AchievementListGroupingLockState)
		if achievementList != nil {
			for _, bucket := range achievementList.Buckets {
				fmt.Printf("\n[%s] (%d achievements)\n", bucket.Label, len(bucket.Achievements))
				for _, ach := range bucket.Achievements {
					status := "🔒"
					if ach.State == rcheevos.AchievementStateUnlocked {
						status = "✅"
					}
					fmt.Printf("  %s %s (%d pts) - %s\n", status, ach.Title, ach.Points, ach.Description)
					if ach.MeasuredProgress != "" {
						fmt.Printf("      Progress: %s (%.1f%%)\n", ach.MeasuredProgress, ach.MeasuredPercent)
					}
				}
			}
			achievementList.Destroy()
		}
	}

	// List leaderboards
	if client.HasLeaderboards() {
		fmt.Println("\n--- Leaderboards ---")
		leaderboardList := client.CreateLeaderboardList(rcheevos.LeaderboardListGroupingNone)
		if leaderboardList != nil {
			for _, bucket := range leaderboardList.Buckets {
				for _, lb := range bucket.Leaderboards {
					fmt.Printf("  📊 %s - %s\n", lb.Title, lb.Description)
				}
			}
			leaderboardList.Destroy()
		}
	}

	if client.HasRichPresence() {
		fmt.Printf("\nRich Presence: %s\n", client.GetRichPresenceMessage())
	}

	// In a real emulator, you would call client.DoFrame() every frame
	// For this example, we'll just do a few frames
	fmt.Println("\nSimulating frames...")
	for i := 0; i < 60; i++ {
		client.DoFrame()
	}

	// Demonstrate progress serialization (for savestates)
	progressData, result := client.SerializeProgress()
	if result == rcheevos.OK && progressData != nil {
		fmt.Printf("\nSerialized progress: %d bytes (include this in savestates)\n", len(progressData))

		// Restore progress (normally done when loading a savestate)
		result = client.DeserializeProgress(progressData)
		if result == rcheevos.OK {
			fmt.Println("Progress restored successfully")
		}
	}

	fmt.Println("\nDone!")
}
