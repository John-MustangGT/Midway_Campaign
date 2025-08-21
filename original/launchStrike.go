package main

import (
	"fmt"
   "time"
   "strconv"
   "strings"
)

func launchStrike() {
	carrier := selectCarrier()
	if carrier == -1 {
		return
	}
	
   deckSbd := carriers[carrier].deckSbd
   if deckSbd >= 1000 {
   	deckSbd -= 1000
   }
   if deckSbd+carriers[carrier].deckTbd == 0 || carriers[carrier].deckSbd < 1000 {
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
	distance := getDistance(carrier, target)
	if distance > 200 {
		fmt.Printf("%.0f NAUTICAL MILES, OUT OF RANGE.\n", distance)
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
	
	// Launch strike
	flightTime := distance * 0.3
	strikes[strikeSlot].f4f = carriers[carrier].deckF4f
	strikes[strikeSlot].sbd = carriers[carrier].deckSbd
	strikes[strikeSlot].tbd = carriers[carrier].deckTbd
	strikes[strikeSlot].target = float64(target)
	strikes[strikeSlot].arrivalTime = gameTime + flightTime
	strikes[strikeSlot].returnTime = gameTime + flightTime + flightTime
	strikes[strikeSlot].launched = carrier
	
	// Clear deck
	carriers[carrier].deckF4f = 0
	carriers[carrier].deckSbd = 0
	carriers[carrier].deckTbd = 0
	
	fmt.Printf("%s STRIKE TAKING OFF.\n", getCarrierName(carrier))
	displayCarriers()
}
