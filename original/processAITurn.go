package main

import (
   "math/rand"
   "math"
)

func processAITurn() {
	// Convert deck SBDs back from 1000+ encoding (BASIC line 880)
	for i := 4; i <= 7; i++ {
		carriers[i].deckSbd = float64(int(carriers[i].deckSbd) % 1000)
	}
	
	// AI fleet movement and attack decisions (BASIC lines 890-1040)
	processAIMovement()
	
	// Japanese CAP management (BASIC lines 1040-1080)  
	processJapaneseCAP()
	
	// Japanese strike planning (BASIC lines 1090-1430)
	processJapaneseStrikes()
	
	// Move fleets and advance time (BASIC lines 1510-1540)
	timeIncrement := 30 + int(30*rand.Float64())
	gameTime += float64(timeIncrement)
	
	// Handle day transitions
	if gameTime >= 1440 {
		gameDay--
		gameTime -= 1440
		if gameDay <= 0 {
			gameOver = true
			return
		}
	}
	
	// Move all fleets
	for i := 0; i < 5; i++ {
		fleets[i].x += float64(timeIncrement) * fleets[i].speed * math.Sin(fleets[i].course) / 60
		fleets[i].y += float64(timeIncrement) * fleets[i].speed * math.Cos(fleets[i].course) / 60
	}
	
	// Process reconnaissance (BASIC lines 1550-1720)
	processReconnaissance()
	
	// Process air strikes (BASIC lines 1730-2430)  
	processAirStrikes()
	
	// Process damage control (BASIC lines 2440-2490)
	processDamageControl()
	
	// Process returning strikes (BASIC lines 2500-2710)
	processReturningStrikes()
	
	// Check victory conditions (BASIC lines 2730-2770)
	checkVictoryConditions()
}
