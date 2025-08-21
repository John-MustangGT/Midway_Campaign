package main

func clearDeck() {
	carrier := selectCarrier()
	if carrier == -1 {
		return
	}
	
	// Handle the MOD 1000 operation from BASIC line 3120
	carriers[carrier].deckSbd = float64(int(carriers[carrier].deckSbd) % 1000)
	
	// Return deck aircraft to hangar
	carriers[carrier].f4f += carriers[carrier].deckF4f
	carriers[carrier].sbd += carriers[carrier].deckSbd
	carriers[carrier].tbd += carriers[carrier].deckTbd
	
	carriers[carrier].deckF4f = 0
	carriers[carrier].deckSbd = 0
	carriers[carrier].deckTbd = 0
	
	displayCarriers()
}
