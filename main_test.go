package main

import (
	"Leagly/config"
	"Leagly/query"
	"testing"
)

type TestCase struct {
	playerName string
	region     string
	region2    string
}

func TestLiveCommand(t *testing.T) {
	err := config.ReadConfig()

	if err != nil {
		t.Error(err)
	}

	testCase := &TestCase{playerName: "sosnek", region: "NA1"}
	_, err = query.IsInGame(testCase.playerName, testCase.region)
	if err != nil {
		t.Error(err)
	}
}

func TestLastMatchCommand(t *testing.T) {
	err := config.ReadConfig()

	if err != nil {
		t.Error(err)
	}

	testCase := &TestCase{playerName: "TFCtrikz", region: "NA1", region2: "americas"}
	_, err = query.GetLastMatch(testCase.playerName, testCase.region, testCase.region2)
	if err != nil {
		t.Error(err)
	}
}

func TestLookupCommand(t *testing.T) {
	err := config.ReadConfig()

	if err != nil {
		t.Error(err)
	}

	testCase := &TestCase{playerName: "Lets", region: "NA1", region2: "americas"}
	_, err = query.LookupPlayer(testCase.playerName, testCase.region, testCase.region2)
	if err != nil {
		t.Error(err)
	}
}

func TestMasteryCommand(t *testing.T) {
	err := config.ReadConfig()

	if err != nil {
		t.Error(err)
	}

	testCase := &TestCase{playerName: "MrBlameyy", region: "NA1"}
	_, err = query.MasteryPlayer(testCase.playerName, testCase.region)
	if err != nil {
		t.Error(err)
	}
}
