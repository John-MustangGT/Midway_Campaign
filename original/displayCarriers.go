package main

import (
   "fmt"
)

// Claude.ai chanegs v1
func displayCarriers() {
	fmt.Printf("CAP - ON DECK - -- BELOW --\n")
	fmt.Printf("    F4F SBD TBD F4F SBD TBD\n")
	
	for i := 4; i <= 7; i++ {
		if carriers[i].damage >= 100 {
			if i == 7 {
				fmt.Printf("** AIRBASE DESTROYED **               \n")
			} else {
				fmt.Printf("** SUNK **                            \n")
			}
		} else if carriers[i].damage >= 60 {
			fmt.Printf("HEAVY DAMAGE      ")
			fmt.Printf("    %3.0f %3.0f %3.0f\n", carriers[i].f4f, carriers[i].sbd, carriers[i].tbd)
		} else {
			displaySbd := getDisplaySBD(carriers[i].deckSbd)
			fmt.Printf("%3.0f %3.0f %3.0f %3.0f %3.0f %3.0f\n",
				carriers[i].cap,
				carriers[i].deckF4f,
				displaySbd,
				carriers[i].deckTbd,
				carriers[i].f4f,
				carriers[i].sbd)
		}
	}
	fmt.Println()
}
