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
	
	// BASIC line 490: Check if strike already armed (special handling for SBD encoding)
	deckSbd := carriers[carrier].deckSbd
	if deckSbd >= 1000 {
		deckSbd -= 1000
	}
	
	// BASIC line 500: Special case - if SBD=1000 exactly, it means no strike (reset to 0)
	if carriers[carrier].deckF4f+carriers[carrier].deckTbd == 0 && carriers[carrier].deckSbd == 1000 {
		carriers[carrier].deckSbd = 0
	}
	
	// Check if any aircraft already on deck
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
	
	// BASIC lines 530-550: Limit to available aircraft
	if float64(f4f) > carriers[carrier].f4f {
		f4f = int(carriers[carrier].f4f)
	}
	if float64(sbd) > carriers[carrier].sbd {
		sbd = int(carriers[carrier].sbd)
	}
	if float64(tbd) > carriers[carrier].tbd {
		tbd = int(carriers[carrier].tbd)
	}
	
	// BASIC line 560: Move aircraft to deck (with special SBD encoding)
	carriers[carrier].deckF4f = float64(f4f)
	carriers[carrier].f4f -= float64(f4f)
	carriers[carrier].deckSbd = 1000 + float64(sbd) // Special encoding for armed strike
	carriers[carrier].sbd -= float64(sbd)
	carriers[carrier].deckTbd = float64(tbd)
	carriers[carrier].tbd -= float64(tbd)
	
	displayCarriers()
}
