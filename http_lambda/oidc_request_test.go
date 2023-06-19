package main

import "testing"

func TestDecodeState1(t *testing.T) {
	testState := "ogGkAXBHb1ZncGMzRG9BTUNMclE9AgEDeCFBV1NfODY4MDI0ODk5NTMxX2llc2F3c25hLXNhbmRib3gEZTEuMS4yAlhHMEUCIQDbpP-DGD014NUrrcZGAxnJNeKegSL8aWs-7CBIqFvgywIgHmG-akscs3z-HiYMzhX-M3GynG5U0gIxLvw3GtOQXGw"
	si := &SignedInformation{}
	err := si.DecodeFromString(testState)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%+v", si)
}
