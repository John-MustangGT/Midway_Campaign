package main

import (
	"fmt"
   "math/rand"
   "time"
   "strconv"
   "strings"
)

func launchStrike() {
	carrier := selectCarrier()
	if carrier == -1 {
		return
	}
	
	// BASIC line 700: Check if strike is ready (SBD must be 1000+ AND have bombers)
	deckSbd := carriers[carrier].deckSbd
	if deckSbd < 1000 || (deckSbd-1000)+carriers[carrier].deckTbd == 0 {
		fmt.Printf("%s HAS NO STRIKE READY.\n", getCarrierName(carrier))
		time.Sleep(1 * time.Second)
		return
	}
	
	// Find available targets
	targets := []int{}
	for i := 0; i < 3; i++ {
		if fleets[i].damage > 0 {
			targets = append(targets, i)
		}
	}
	
	if len(targets) == 0 {
		fmt.Println("NO TARGETS.")
		time.Sleep(1 * time.Second)
		return
	}
	
	var target int
	if len(targets) > 1 {
		fmt.Print("TARGET CONTACT: ")
		if !scanner.Scan() {
			return
		}
		targetNum, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
		if err != nil || targetNum < 1 || targetNum > len(targets) {
			return
		}
		target = targets[targetNum-1]
	} else {
		target = targets[0]
	}
	
	// Check range
	carrierFleet := int(carriers[carrier].fleet)
	distance := getDistance(carrierFleet, target)
	if distance > 200 {
		fmt.Printf("%.0f NAUTICAL MILES, OUT OF RANGE.\n", distance)
		time.Sleep(1 * time.Second)
		return
	}
	
	// Check timing constraints (BASIC lines 750-780)
	flightTime := distance * 0.3
	if carrier != 7 && (gameTime+flightTime+flightTime > 240 && gameTime+flightTime+flightTime <= 1140) {
		// No night carrier landings
	} else if carrier != 7 {
		fmt.Println("NO NIGHT CARRIER LANDINGS.")
		time.Sleep(1 * time.Second)
		return
	}
	
	if gameTime+flightTime < 240 || gameTime+flightTime > 1140 {
		fmt.Println("NO NIGHT ATTACKS.")
		time.Sleep(1 * time.Second)
		return
	}
	
	// Find available strike slot
	strikeSlot := -1
	for i := 0; i < 10; i++ {
		if strikes[i].launched == -1 {
			strikeSlot = i
			break
		}
	}
	
	if strikeSlot == -1 {
		fmt.Println("TOO MANY STRIKES ALOFT.")
		time.Sleep(1 * time.Second)
		return
	}
	
	// Launch strike (BASIC line 830)
	strikes[strikeSlot].f4f = carriers[carrier].deckF4f
	strikes[strikeSlot].sbd = carriers[carrier].deckSbd  // Keep the 1000+ encoding
	strikes[strikeSlot].tbd = carriers[carrier].deckTbd
	strikes[strikeSlot].target = float64(target)
	strikes[strikeSlot].arrivalTime = gameTime + flightTime
	strikes[strikeSlot].returnTime = gameTime + flightTime + flightTime
	strikes[strikeSlot].launched = carrier
	strikes[strikeSlot].escort = 1
	strikes[strikeSlot].bomberType = 0
	
	// BASIC line 850: Determine if SBDs or TBDs lead attack
	sbdCount := carriers[carrier].deckSbd
	if sbdCount >= 1000 {
		sbdCount -= 1000
	}
	if sbdCount/(sbdCount+carriers[carrier].deckTbd) > rand.Float64() {
		strikes[strikeSlot].bomberType = -1
	}
	
	// Clear deck
	carriers[carrier].deckF4f = 0
	carriers[carrier].deckSbd = 0
	carriers[carrier].deckTbd = 0
	
	fmt.Printf("%s STRIKE TAKING OFF.\n", getCarrierName(carrier))
	displayCarriers()
}
