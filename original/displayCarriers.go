package main

import (
   "fmt"
)

func displayCarriers() {
	for i := 4; i <= 7; i++ {
		// Print carrier name first
		fmt.Printf("%-11s", getCarrierName(i))
		
		if carriers[i].damage >= 60 {
			if carriers[i].damage >= 100 {
				if i == 7 {
					fmt.Printf("** AIRBASE DESTROYED **   \n")
				} else {
					fmt.Printf("** SUNK **                \n")
				}
			} else {
				fmt.Printf("HEAVY DAMAGE      ")
				// Show hangar aircraft only when heavily damaged
				fmt.Printf("    %3.0f %3.0f %3.0f\n", carriers[i].f4f, carriers[i].sbd, carriers[i].tbd)
			}
		} else {
			// Normal display: CAP, then deck aircraft (with MOD 1000 for SBDs), then hangar
			displaySbd := int(carriers[i].deckSbd) % 1000
         fmt.Printf("%3.0f %3.0f %3.0f %3.0f %3.0f %3.0f %3.0f\n",
         	carriers[i].cap,           // CAP F4F
         	carriers[i].deckF4f,       // DECK F4F  
         	float64(displaySbd),       // DECK SBD
         	carriers[i].deckTbd,       // DECK TBD
         	carriers[i].f4f,           // HANGAR F4F
         	carriers[i].sbd,           // HANGAR SBD
         	carriers[i].tbd)           // HANGAR TBD - THIS WAS MISSING!
		}
	}
	fmt.Println()
}
