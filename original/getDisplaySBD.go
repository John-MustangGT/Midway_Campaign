package main

//import (
//)

func getDisplaySBD(deckSbd float64) float64 {
	if deckSbd >= 1000 {
		return deckSbd - 1000
	}
	return deckSbd
}
