package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/dotabuff/yasha"
	"github.com/dotabuff/yasha/dota"
)

var pp = spew.Dump

var whitelist = map[uint64]bool{
// 76561197985166697: true
}

var steamHero = map[uint64]string{}
var suffixSteam = map[string]uint64{} // 0005

type playerGold struct {
	suffix string
	minute int
	gold   uint
	teamID int
}

var teams = []struct {
	name string
	id   int
}{
	{"Radiant", 2},
	{"Dire", 3},
	{"Spectator", 5},
}

type heroState struct {
	gold uint
	posX int
	posY int
	posZ int
}

func main() {
	if len(os.Args) < 2 {
		spew.Println("Expected a .dem file as argument")
	}

	for _, path := range os.Args[1:] {

		visualizer := &GifVisualizer{}
		visualizer.Setup(path)
		fmt.Println("path", path)

		parser := yasha.ParserFromFile(path)

		players := []*playerGold{}
		for _, p := range playersSuffix() {
			players = append(players, &playerGold{suffix: p})
		}

		parser.OnFileInfo = func(fileinfo *dota.CDemoFileInfo) {
			for _, pls := range fileinfo.GetGameInfo().GetDota().GetPlayerInfo() {
				steamHero[pls.GetSteamid()] = strings.TrimPrefix(pls.GetHeroName(), "npc_dota_hero_")
			}
		}

		//	playersGold := make(map[string]uint)
		// minuteGold := make(map[int]map[string]uint)

		heroMinuteState := make(map[int]map[string]*heroState)

		heroPos := make(map[string]heroState)

		parser.OnEntityPreserved = func(e *yasha.PacketEntity) {
			if e == nil {
				return
			}
			minute := e.Tick / 1800

			if e.Name == "DT_DOTA_BaseNPC" {
				// Position tracking
				jsonify(e.Values)
				for _, p := range players {
					fmt.Println("POS", e.Values[e.Name+".m_cellX."+p.suffix])
					if posX, ok := e.Values[e.Name+".m_cellX."+p.suffix].(int); ok {
						if posY, ok := e.Values[e.Name+".m_cellY."+p.suffix].(int); ok {
							if posZ, ok := e.Values[e.Name+".m_cellZ."+p.suffix].(int); ok {
								heroPos[p.suffix] = heroState{
									posX: posX,
									posY: posY,
									posZ: posZ,
								}
								fmt.Println("woot", heroPos[p.suffix])
								os.Exit(0)
							}
						}
					}
				}
			}

			// GOLD tracking.
			if e.Name != "DT_DOTA_PlayerResource" {
				return
			}
			playersGold := make(map[string]uint)
			for _, p := range players {
				g, ok := e.Values["m_iTotalEarnedGold."+p.suffix].(uint)
				if !ok {
					fmt.Println("not uint")
					return
				}
				if t, ok := e.Values["m_iPlayerTeams."+p.suffix].(int); ok {
					p.teamID = t
				}
				suffixSteam[p.suffix] = e.Values["m_iPlayerSteamIDs."+p.suffix].(uint64)

				if minute != p.minute {
					delta := g - p.gold // removes the background ones
					if delta == 0 {
						return
					}
					//	p.totalEarned = append(p.totalEarned, struct{ X, Y float64 }{float64(e.Tick / 1800), float64(delta) - 100})
					//	p.totalGained = append(p.totalEarned, struct{ X, Y float64 }{float64(e.Tick / 1800), float64(delta)})
					p.gold = g
					p.minute = minute
					playersGold[p.suffix] = delta // p.gold
				}
			}

			if _, ok := heroMinuteState[minute]; !ok {
				heroMinuteState[minute] = make(map[string]*heroState)
			}
			for _, p := range playersSuffix() {
				if heroMinuteState[minute][p] == nil {
					heroMinuteState[minute][p] = new(heroState)
				}
				if g, ok := playersGold[p]; ok {
					heroMinuteState[minute][p].gold = g
				}
				if pos, ok := heroPos[p]; ok {
					heroMinuteState[minute][p].posX = pos.posX
					heroMinuteState[minute][p].posY = pos.posY
					heroMinuteState[minute][p].posZ = pos.posZ
				}

			}

			//			if _, ok := minuteGold[minute]; !ok && len(playersGold) > 0 {
			//			minuteGold[minute] = playersGold
			//		}
		}
		parser.Parse() // this can panic :-/
		for _, team := range teams {
			for _, pl := range players {
				if len(whitelist) > 0 {
					if _, ok := whitelist[suffixSteam[pl.suffix]]; !ok {
						continue
					}
				}
				if pl.teamID != team.id {
					continue
				}
				//	steamHero[suffixSteam[pl.suffix]]
			}
			//	fmt.Sprintf("gold-%v-%v.png", path, team.name)); err != nil {
		}

		for _, state := range heroMinuteState {
			visualizer.AddFrame(state)
			fmt.Println("frame added", state)
			//	minuteGold[minute] = true
		}
		visualizer.Complete()
	}
}

func jsonify(in interface{}) {
	data, err := json.MarshalIndent(in, "", "  ")
	if err != nil {
		panic(err)
	}
	spew.Println(string(data))
}

func playersSuffix() []string {
	players := make([]string, 0, 10)
	for i := 0; i < 10; i++ {
		players = append(players, fmt.Sprintf("%04d", i))
	}
	return players
}
