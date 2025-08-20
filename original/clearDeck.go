package main

func clearDeck() {
	carrier := selectCarrier()
	if carrier == -1 {
		return
	}
	
	// Convert any deck SBDs that are coded as 1000+ back to normal
	if carriers[carrier].deckSbd >= 1000 {
		carriers[carrier].deckSbd -= 1000
	}
	
	// Return deck aircraft to hangar
	carriers[carrier].f4f += carriers[carrier].deckF4f
	carriers[carrier].sbd += carriers[carrier].deckSbd
	carriers[carrier].tbd += carriers[carrier].deckTbd
	
	carriers[carrier].deckF4f = 0
	carriers[carrier].deckSbd = 0
	carriers[carrier].deckTbd = 0
	
	displayCarriers()
}
