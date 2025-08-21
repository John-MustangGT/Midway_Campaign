package main

import (
   "math/rand"
   "math"
   "fmt"
   "time"
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

func processAIMovement() {
	// BASIC lines 890-990: Check for proximity to US forces
	for i := 3; i <= 4; i++ {
		distance := getDistance(i, 2) // Distance to cruiser group
		if distance < 50 {
			c5 = 10
		}
	}
	
	// Japanese movement logic
	skip = false
	distance := getDistance(1, 5) // Transport to Midway
	if distance < 15 {
		fleets[1].speed = 0
	}
	
	// Retreat conditions
	if skip || c5 > 9 {
		fleets[2].speed = 25 + 15 // Cruisers retreat
		if c7 > 255 {
			fleets[2].speed += 15
		}
		fleets[2].course = 270 * pi
	}
	
	if c5 > 9 {
		// Process cruiser bombardment
		distance = getDistance(2, 5)
		if distance <= 15 {
			fmt.Print("CRUISERS BOMBARD ")
			fmt.Print("MIDWAY")
			fleets[2].fleetType = 2
			c5++
			if skip == false && c7 <= 255 {
				fleets[2].speed = 0
			}
			if c6 > 0 {
				gameOver = false // Force update
				c6 = 0
			}
			
			// Calculate bombardment damage
			hits := 0
			nearMisses := 0
			for k := int(c7); k <= 255; k += 4 {
				r := rand.Float64()
				if r < 0.05 {
					hits++
				}
				if r < 0.1 {
					nearMisses++
				}
			}
			nearMisses -= hits
			
			// Apply damage to Midway
			applyDamage(7, hits, nearMisses, 24)
			displayCarriers()
		}
	}
	
	// Carrier group movement
	distance = getDistance(5, 0) // Midway to Japanese carriers
	if distance > 250 {
		angle := getAngle(0, 5)
		fleets[0].course = angle
	}
	if distance < 100 {
		angle := getAngle(5, 0)
		fleets[0].course = angle
	}
	
	// Move carriers toward spotted US forces
	for k := 6; k >= 4; k-- {
		fleetIndex := int(carriers[k].fleet)
		if fleets[fleetIndex].damage > 0 && carriers[k].damage < 100 {
			angle := getAngle(0, fleetIndex)
			fleets[0].course = angle
		}
	}
	
	if skip {
		fleets[0].course = 270 * pi
	}
}

func processJapaneseCAP() {
	// BASIC lines 1040-1080: Japanese CAP management
	for i := 0; i <= 3; i++ {
		if carriers[i].cap >= 5 || carriers[i].damage >= 60 {
			continue
		}
		
		// Move hangar fighters to CAP
		carriers[i].cap += carriers[i].f4f
		carriers[i].f4f = 0
		
		if carriers[i].cap < 5 {
			// Add deck fighters if needed
			carriers[i].cap += carriers[i].deckF4f
			carriers[i].deckF4f = 0
			if carriers[i].cap > 5 {
				carriers[i].deckF4f = carriers[i].cap - 5
				carriers[i].cap = 5
			}
		} else {
			carriers[i].f4f = carriers[i].cap - 5
			carriers[i].cap = 5
		}
	}
}

func applyDamage(target int, hits int, nearMisses int, damageType int) {
	if target >= 5 { // Fleet target
		fmt.Printf("%s!", getFleetName(target-5))
		fmt.Printf("%s", getFleetName(target-5))
		printHitMessage(hits, nearMisses)
		damage := damageType * (hits + nearMisses/3)
		victory1 += float64(damage)
		fmt.Printf("%d VICTORY POINTS AWARDED.", damage)
		
		if target == 2 { // Cruiser group
			c7 += float64(damage)
			if c7 >= 255 && c7-float64(damage) < 255 {
				time.Sleep(1 * time.Second)
				fmt.Println("CRUISERS SEVERELY CRIPPLED.")
				fleets[2].aa *= 3
			}
			fleets[2].speed = 10
			c5 = 10
			if c7 >= 512 {
				time.Sleep(1 * time.Second)
				fmt.Println("ALL CRUISERS ARE SUNK!")
				victory0 -= c7 - 512
				c7 = 512
				fleets[2].fleetType = 0
				fleets[2].speed = 0
				fleets[2].x = -1000
			}
		}
		time.Sleep(1 * time.Second)
	} else { // Carrier target
		fmt.Printf("%s!", getCarrierName(target))
		fmt.Printf("%s", getCarrierName(target))
		printHitMessage(hits, nearMisses)
		
		hasSecondary := (target == 2 || target == 7) && 
			(carriers[target].deckF4f+carriers[target].deckSbd+carriers[target].deckTbd > 0)
		if hasSecondary && hits+nearMisses > 0 {
			fmt.Print("SECONDARY EXPLOSIONS!")
		}
		
		finalDamage := damageType * (1 + 0)
		if hasSecondary {
			finalDamage *= 2
		}
		
		// Apply hit damage
		for h := 0; h < hits; h++ {
			damage := float64(finalDamage) * rand.Float64()
			if target == 7 {
				damage /= 3
			}
			if target == 8 {
				damage *= 2
			}
			
			// Damage aircraft if not torpedo attack on non-Midway target
			if !(damageType == 4 && target != 7) {
				for aircraft := 1; aircraft <= 6; aircraft++ {
					if gameTime >= 240 && gameTime <= 1140 {
						continue // Daytime, aircraft might be flying
					}
					// Damage aircraft based on probability
					switch aircraft {
					case 1:
						for plane := 0; plane < int(carriers[target].f4f); plane++ {
							if rand.Float64()*100 < damage {
								carriers[target].f4f++
							}
						}
					case 2:
						for plane := 0; plane < int(carriers[target].sbd); plane++ {
							if rand.Float64()*100 < damage {
								carriers[target].sbd++
							}
						}
					// ... similar for other aircraft types
					}
				}
			}
			
			carriers[target].damage += damage
			if carriers[target].damage >= 60 {
				// Return deck aircraft to hangar when heavily damaged
				clearDeckForDamage(target)
				if carriers[target].damage >= 100 {
					gameOver = false // Force display update
					carriers[target].damage = 100
					time.Sleep(1 * time.Second)
					if target == 7 {
						fmt.Printf("%s AIRBASE DESTROYED!", getCarrierName(target))
					} else {
						fmt.Printf("%s BLOWS UP AND SINKS!", getCarrierName(target))
					}
					// Destroy all aircraft
					for aircraft := 1; aircraft <= 7; aircraft++ {
						switch aircraft {
						case 1:
							carriers[target].f4f = 0
						case 2:
							carriers[target].sbd = 0
						case 3:
							carriers[target].tbd = 0
						case 4:
							carriers[target].deckF4f = 0
						case 5:
							carriers[target].deckSbd = 0
						case 6:
							carriers[target].deckTbd = 0
						case 7:
							carriers[target].cap = 0
						}
					}
				}
			}
		}
		
		// Apply near miss damage (1/3 effect)
		nearMissDamage := float64(finalDamage) / 3
		for n := 0; n < nearMisses; n++ {
			damage := nearMissDamage * rand.Float64()
			if target == 7 {
				damage /= 3
			}
			if target == 8 {
				damage *= 2
			}
			carriers[target].damage += damage
		}
	}
	
	time.Sleep(1 * time.Second)
}

func clearDeckForDamage(carrier int) {
	carriers[carrier].deckSbd = float64(int(carriers[carrier].deckSbd) % 1000)
	carriers[carrier].f4f += carriers[carrier].deckF4f
	carriers[carrier].sbd += carriers[carrier].deckSbd
	carriers[carrier].tbd += carriers[carrier].deckTbd
	carriers[carrier].deckF4f = 0
	carriers[carrier].deckSbd = 0
	carriers[carrier].deckTbd = 0
}

func printHitMessage(hits int, nearMisses int) {
	fmt.Printf(" TAKES %d HITS", hits)
	if nearMisses > 0 {
		fmt.Printf("\nAND %d NEAR MISSES", nearMisses)
	}
	fmt.Println(".")
}
