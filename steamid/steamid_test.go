package steamid_test

import (
	"testing"

	"github.com/eclou/go-steamid/steamid"
)

func Test_SteamId(t *testing.T) {
	steamId, err := steamid.NewSteamId("76561198435082001")
	if err != nil {
		t.Fatalf("parse steam id failed: %v", err)
	}

	steam2Id := steamId.GetSteam2RenderedID()

	if steam2Id != "STEAM_0:1:237408136" {
		t.Error("expected steam2 id STEAM_0:1:237408136")
	}

	steamId2 := steamid.SteamId{}

	err = steamId2.FromTradeUrl("https://steamcommunity.com/tradeoffer/new/?partner=474816273&token=RZbHYcrV")

	if err != nil {
		t.Fatalf("parse trade url failed: %v", err)
	}

	if steamId2.GetSteamID64() != "76561198435082001" {
		t.Error("expected steam2 id 76561198435082001")
	}
}
