package main

import (
   "fmt"
   "time"
   "strings"
   "strconv"
)

func armStrike() {
	carrier := selectCarrier()
	if carrier == -1 {
		return
	}
	
	// Check if strike already armed (including the 1000+ encoding for SBDs)
	deckSbd := carriers[carrier].deckSbd
	if deckSbd >= 1000 {
		deckSbd -= 1000
	}
	
	if carriers[carrier].deckF4f+deckSbd+carriers[carrier].deckTbd > 0 {
		fmt.Printf("%s STRIKE ALREADY ON DECK.\n", getCarrierName(carrier))
		time.Sleep(1 * time.Second)
		return
	}
	
	fmt.Printf("BRING AIRCRAFT TO %s DECK.\n", getCarrierName(carrier))
	fmt.Print("F4F,SBD,TBD: ")
	
	if !scanner.Scan() {
		return
	}
	
	parts := strings.Fields(scanner.Text())
	if len(parts) < 3 {
		return
	}
	
	f4f, _ := strconv.Atoi(parts[0])
	sbd, _ := strconv.Atoi(parts[1])
	tbd, _ := strconv.Atoi(parts[2])
	
	// Limit to available aircraft
	if float64(f4f) > carriers[carrier].f4f {
		f4f = int(carriers[carrier].f4f)
	}
	if float64(sbd) > carriers[carrier].sbd {
		sbd = int(carriers[carrier].sbd)
	}
	if float64(tbd) > carriers[carrier].tbd {
		tbd = int(carriers[carrier].tbd)
	}
	
	// Move aircraft to deck (with special SBD encoding)
	carriers[carrier].deckF4f = float64(f4f)
	carriers[carrier].f4f -= float64(f4f)
	carriers[carrier].deckSbd = 1000 + float64(sbd) // Special encoding
	carriers[carrier].sbd -= float64(sbd)
	carriers[carrier].deckTbd = float64(tbd)
	carriers[carrier].tbd -= float64(tbd)
	
	displayCarriers()
}
