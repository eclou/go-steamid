#### GoSteamId
go version of [https://github.com/DoctorMcKay/node-steamid](https://github.com/DoctorMcKay/node-steamid)

```
go get github.com/eclou/go-steamid
```

##### Init from trade url

```
steamId := steamid.SteamId{}

err = steamId.FromTradeUrl("https://steamcommunity.com/tradeoffer/new/?partner=474816273&token=RZbHYcrV")
if err != nil {
    t.Fatalf("invalid trade url: %v", err)
}
```