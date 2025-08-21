package main

import (
   "fmt"
   "time"
   "math"
   "math/rand"
)

func processJapaneseStrikes() {
	if gameTime > 1140 { // No strikes at night
		return
	}
	
	s9 := 0.0  // Target selection
	a9 := 0.0  // All-out attack flag
	a8 := 0.0  // Half-attack flag
	
	// Check if Japanese carriers have aircraft ready for strikes
	i := 0
	hasReadyAircraft := false
	for i = 0; i < 4; i++ {
		if carriers[i].deckF4f+carriers[i].deckSbd+carriers[i].deckTbd > 0 {
			hasReadyAircraft = true
			break
		}
	}
	
	if !hasReadyAircraft {
		s9 = 0
	} else {
		// Check for operational carriers (BASIC line 1130)
		i = 4
		for i = 4; i < 8; i++ {
			if carriers[i].damage < 60 {
				fleet := int(carriers[i].fleet)
				canReach := checkStrikeRange(fleet, 0)
				if canReach {
					break
				}
			}
		}
		
		if i >= 8 {
			// No operational US carriers in range, try attacking fleet directly
			i = 4
			for i = 4; i < 8; i++ {
				if carriers[i].damage < 100 {
					fleet := int(carriers[i].fleet)
					canReach := checkStrikeRange(fleet, 0)
					if canReach {
						break
					}
				}
			}
			
			if i >= 8 {
				// Try Midway
				canReach := checkStrikeRange(5, 0)
				if canReach {
					i = -7 // Special code for Midway
				}
			}
		}
		
		if i < 8 || i == -7 {
			if i == -7 {
				s9 = 5 // Midway
			} else {
				s9 = carriers[i].fleet
			}
		}
	}
	
	// Check if any strikes already targeting the same area
	if s9 >= 3 {
		for i := 0; i < 10; i++ {
			if strikes[i].target < 5 && strikes[i].launched != -1 && strikes[i].escort != -1 {
				s9 = 0 // Don't launch if strike already en route to carriers
				break
			}
		}
	}
	
	// Determine attack intensity
	if fleets[3].damage+fleets[4].damage > 0 { // US task forces spotted
		a9 = 1
	}
	
	// Check for Midway attack conditions
	distance := getDistance(5, 0)
	if distance <= 235 {
		flightTime := 60 * distance / 235
		if gameTime+flightTime >= 240 && gameTime+flightTime+flightTime <= 1140 {
			a8 = 1
			if carriers[3].sbd < 12 { // Hiryu has few dive bombers
				a9 = 1
			}
		}
	}
	
	if a9 == 1 {
		a8 = 0 // All-out overrides half-attack
	}
	
	// Launch strike if target selected
	if s9 >= 3 {
		strikeSlot := -1
		for j := 0; j < 10; j++ {
			if strikes[j].launched == -1 {
				strikeSlot = j
				break
			}
		}
		
		if strikeSlot != -1 {
			strikes[strikeSlot].target = s9
			strikes[strikeSlot].launched = 0 // Japanese strike
			
			distance := getDistance(int(s9), 0)
			flightTime := 60 * distance / 235
			strikes[strikeSlot].arrivalTime = gameTime + flightTime
			strikes[strikeSlot].returnTime = gameTime + flightTime + flightTime
			
			strikes[strikeSlot].f4f = 0
			strikes[strikeSlot].sbd = 0
			strikes[strikeSlot].tbd = 0
			
			// Collect aircraft from operational carriers
			for i := 0; i < 4; i++ {
				if carriers[i].damage < 60 {
					strikes[strikeSlot].f4f += carriers[i].deckF4f
					strikes[strikeSlot].sbd += carriers[i].deckSbd
					strikes[strikeSlot].tbd += carriers[i].deckTbd
					carriers[i].deckF4f = 0
					carriers[i].deckSbd = 0
					carriers[i].deckTbd = 0
				}
			}
			
			if strikes[strikeSlot].sbd+strikes[strikeSlot].tbd == 0 {
				strikes[strikeSlot].launched = -1 // Cancel if no bombers
			} else {
				strikes[strikeSlot].escort = 1
				strikes[strikeSlot].bomberType = 0
				if strikes[strikeSlot].sbd/(strikes[strikeSlot].sbd+strikes[strikeSlot].tbd) > rand.Float64() {
					strikes[strikeSlot].bomberType = -1
				}
			}
		}
	}
	
	// Prepare deck strikes for Japanese carriers
	for i := 0; i < 4; i++ {
		clearDeckForDamage(i) // Handle any damage first
		if carriers[i].damage < 60 {
			if a9 != 0 {
				// All-out attack - everything to deck
				carriers[i].deckF4f = carriers[i].f4f
				carriers[i].deckSbd = carriers[i].sbd
				carriers[i].deckTbd = carriers[i].tbd
				carriers[i].f4f = 0
				carriers[i].sbd = 0
				carriers[i].tbd = 0
			} else if a8 != 0 {
				// Half attack
				carriers[i].deckF4f = math.Floor(carriers[i].tbd / 2)
				carriers[i].deckSbd = math.Floor(carriers[i].sbd / 2)
				carriers[i].f4f -= carriers[i].deckF4f
				carriers[i].sbd -= carriers[i].deckSbd
				carriers[i].deckTbd = math.Floor(carriers[i].tbd / 2)
				carriers[i].tbd -= carriers[i].deckTbd
			}
			
			// If no strikes planned, maintain CAP
			if s9+a8+a9 == 0 {
				carriers[i].cap += carriers[i].f4f
				carriers[i].f4f = 0
			}
		}
	}
}

func selectJapaneseTarget() int {
	// Prefer spotted US task forces, then Midway
	for i := 3; i <= 4; i++ {
		if fleets[i].damage >= 2 { // Spotted and identified
			return i
		}
	}
	
	// Attack Midway if in range
	distance := getDistance(0, 5)
	if distance <= 235 {
		return 5
	}
	
	return -1
}

func checkStrikeRange(target int, attacker int) bool {
	distance := getDistance(target, attacker)
	if fleets[attacker].fleetType == 0 || distance > 235 {
		return false
	}
	
	flightTime := distance * 60 / 235
	if gameTime+flightTime < 240 || gameTime+flightTime+flightTime > 1140 {
		return false
	}
	
	return true
}

func processAirStrikes() {
	// Process all active air strikes (BASIC lines 1730-2430)
	for i := 0; i < 10; i++ {
		if strikes[i].launched == -1 || strikes[i].arrivalTime > gameTime || strikes[i].escort == -1 {
			continue
		}
		
		// Strike arrives at target
		targetIndex := int(strikes[i].target)
		isJapanese := 1
		if strikes[i].launched >= 4 {
			isJapanese = 0 // US strike
		}
		
		// Use the variables to avoid compiler warnings
		_ = targetIndex
		_ = isJapanese
		
		// Check for aborts due to weather/navigation
		for k := 0; k <= 4; k += 2 {
			var aircraftCount float64
			if k == 0 {
				aircraftCount = strikes[i].f4f
			} else if k == 2 {
				aircraftCount = strikes[i].sbd
			} else {
				aircraftCount = strikes[i].tbd
			}
			
			if aircraftCount == 0 {
				continue
			}
			
			// Navigation check
			if rand.Float64() > (strikes[i].returnTime-strikes[i].arrivalTime-20)/100 {
				if k == 0 {
					strikes[i].f4f = -1 // Mark as aborted
				} else if k == 2 {
					strikes[i].escort = -1
				} else {
					strikes[i].bomberType = -1
				}
			}
		}
		
		// Check if strike was aborted
		if strikes[i].escort == -1 {
			continue
		}
		
		// Handle fighter cover assignment
		if strikes[i].bomberType != -1 {
			if strikes[i].escort == -1 {
				strikes[i].bomberType = 1 - strikes[i].bomberType
				if strikes[i].escort == -1 {
					strikes[i].bomberType = -1
				}
			}
		}
		
		// Combat resolution would continue here...
		fmt.Printf("AIR STRIKE ATTACKING TARGET!\n")
		time.Sleep(2 * time.Second)
		
		// Mark escorts as used
		for k := 1; k <= 5; k += 2 {
			if k == 1 {
				strikes[i].f4f = -1
			} else if k == 3 {
				strikes[i].escort = -1
			} else {
				strikes[i].bomberType = -1
			}
		}
		
		displayCarriers()
	}
}

func processDamageControl() {
	// BASIC lines 2440-2490: Handle secondary explosions and damage control
	for i := 0; i < 9; i++ {
		if carriers[i].damage < 10 || carriers[i].damage >= 100 {
			continue
		}
		
		// Random secondary explosions
		if rand.Float64() <= 0.05*(1-0) { // Reduced chance for US carriers
			if i < 4 { // Japanese carriers more prone to explosions
				fmt.Printf("EXPLOSION ON %s!", getCarrierName(i))
				applyDamage(i, 1, 0, 12)
				displayCarriers()
			}
		}
		
		// Damage control - slow repair
		if rand.Float64() <= 0.2*(1-0) { // Reduced chance for US carriers  
			if i > 3 && i < 8 { // US carriers have better damage control
				carriers[i].damage -= 5 * rand.Float64()
				if carriers[i].damage < 0 {
					carriers[i].damage = 0
				}
			}
		}
	}
}

func processReturningStrikes() {
	// BASIC lines 2500-2710: Handle strikes returning to carriers
	for j := 0; j < 10; j++ {
		if strikes[j].launched == -1 || gameTime < strikes[j].returnTime {
			continue
		}
		
		carrier := strikes[j].launched
		if carrier < 4 { // Japanese strike returning
			// Find operational Japanese carrier to land on
			for i := 0; i < 4; i++ {
				if carriers[i].damage <= 60 {
					clearDeckForDamage(i)
					carriers[i].f4f += strikes[j].f4f
					carriers[i].sbd += strikes[j].sbd
					carriers[i].tbd += strikes[j].tbd
					displayCarriers()
					break
				}
			}
		} else { // US strike returning
			gameOver = false // Force display update
			
			if carriers[carrier].damage <= 60 {
				fmt.Printf("STRIKE LANDING ON %s.\n", getCarrierName(carrier))
				clearDeckForDamage(carrier)
				carriers[carrier].f4f += strikes[j].f4f
				carriers[carrier].sbd += strikes[j].sbd
				carriers[carrier].tbd += strikes[j].tbd
				displayCarriers()
			} else {
				// Try to divert to other carriers
				foundAlternate := false
				if carrier <= 5 || (carriers[4].damage > 60 && carriers[5].damage > 60) {
					// Try other US carriers
					for k := 3; k < 7; k++ {
						if k == carrier || carriers[k].damage > 60 {
							continue
						}
						
						// Check if alternate carrier is reachable
						distance := getDistance(int(carriers[carrier].fleet), int(carriers[k].fleet))
						if distance/100 < rand.Float64() {
							fmt.Printf("%s STRIKE DIVERTED TO %s\n", 
								getCarrierName(carrier), getCarrierName(k))
							clearDeckForDamage(k)
							carriers[k].f4f += strikes[j].f4f
							carriers[k].sbd += strikes[j].sbd
							carriers[k].tbd += strikes[j].tbd
							displayCarriers()
							foundAlternate = true
							break
						}
					}
				}
				
				if !foundAlternate {
					fmt.Printf("%s STRIKE SPLASHES!\n", getCarrierName(carrier))
					time.Sleep(1 * time.Second)
				}
			}
		}
		
		strikes[j].launched = -1
	}
}

func checkVictoryConditions() {
	// BASIC lines 2730-2770: Check if game should end
	v2 := 0
	for j := 0; j < 10; j++ {
		if strikes[j].launched != -1 {
			v2 = 1
			break
		}
	}
	
	// Check if any US carriers operational
	operationalUS := 0
	for i := 0; i < 4; i++ {
		if carriers[i].damage <= 60 {
			operationalUS++
		}
	}
	
	skip = operationalUS == 0
	
	if v2 == 1 {
		return // Strikes still active, continue game
	}
	
	if skip && fleets[0].x < 0 {
		gameOver = true
		return
	}
	
	// Check if all carriers heavily damaged
	totalUSCarrierDamage := 0.0
	for i := 0; i < 4; i++ {
		totalUSCarrierDamage += carriers[i].damage
	}
	if totalUSCarrierDamage >= 400 {
		gameOver = true
		return
	}
	
	// Check if Japanese retreating
	if fleets[3].x >= 1150 || fleets[4].x >= 1150 {
		gameOver = true
		return
	}
	
	// Check if all Japanese carriers heavily damaged
	totalJapCarrierDamage := 0.0
	for i := 4; i < 8; i++ {
		totalJapCarrierDamage += carriers[i].damage
	}
	if totalJapCarrierDamage < 400 {
		return // Continue game
	}
	
	gameOver = true
}
