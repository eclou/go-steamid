package steamid

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"

	"github.com/eclou/go-steamid/types"
)

const (
	UNIVERSE_INVALID int64 = iota
	UNIVERSE_PUBLIC
	UNIVERSE_BETA
	UNIVERSE_INTERNAL
	UNIVERSE_DEV
)

const (
	INSTANCE_ALL int64 = iota
	INSTANCE_DESKTOP
	INSTANCE_CONSOLE
	INSTANCE_WEB
)

const ACCOUNT_ID_MASK int64 = 0xFFFFFFFF
const ACCOUNT_INSTANCE_MASK int64 = 0x000FFFFF

const CHAT_INSTANCE_FLAGS_CLAN = (ACCOUNT_INSTANCE_MASK + 1) >> 1
const CHAT_INSTANCE_FLAGS_LOBBY = (ACCOUNT_INSTANCE_MASK + 1) >> 2
const CHAT_INSTANCE_FLAGS_MMSLOBBY = (ACCOUNT_INSTANCE_MASK + 1) >> 3

type SteamId struct {
	Universe  int64
	Type      types.Type
	Instance  int64
	AccountId int64
	Token     string
}

func NewSteamId(input string) (*SteamId, error) {
	if input == "" {
		return nil, errors.New("empty string")
	}

	id := SteamId{
		Universe:  UNIVERSE_INVALID,
		Type:      types.TYPE_INVALID,
		Instance:  INSTANCE_ALL,
		AccountId: 0,
		Token:     "",
	}

	if ok, _ := regexp.MatchString(`^\d+$`, input); ok {
		num, err := strconv.ParseInt(input, 10, 64)
		if err != nil {
			return nil, err
		}

		id.AccountId = num & ACCOUNT_ID_MASK
		id.Instance = (num >> 32) & ACCOUNT_INSTANCE_MASK
		id.Type = types.Type((num >> 52) & 0xF)
		id.Universe = num >> 56
	} else if matches := regexp.MustCompile(`^STEAM_([0-5]):([0-1]):([0-9]+)$`).FindStringSubmatch(input); matches != nil {
		universe, _ := strconv.ParseInt(matches[1], 10, 64)
		mod, _ := strconv.ParseInt(matches[2], 10, 64)
		accountid, _ := strconv.ParseInt(matches[3], 10, 64)
		if universe == 0 {
			universe = UNIVERSE_PUBLIC
		}
		id.Universe = universe
		id.Type = types.TYPE_INDIVIDUAL
		id.Instance = INSTANCE_DESKTOP
		id.AccountId = accountid*2 + mod
	} else if matches := regexp.MustCompile(`^\[([a-zA-Z]):([0-5]):([0-9]+)(:[0-9]+)?]$`).FindStringSubmatch(input); matches != nil {
		typeChar := matches[1]
		universe, _ := strconv.ParseInt(matches[2], 10, 64)
		accountid, _ := strconv.ParseInt(matches[3], 10, 64)
		instanceid := matches[4]

		id.Universe = universe
		id.AccountId = accountid

		if instanceid != "" {
			id.Instance, _ = strconv.ParseInt(instanceid[1:], 10, 64)
		}

		switch typeChar {
		case "U":
			id.Type = types.TYPE_INDIVIDUAL
			if instanceid != "" {
				id.Instance = INSTANCE_DESKTOP
			}
		case "c":
			id.Instance |= CHAT_INSTANCE_FLAGS_CLAN
			id.Type = types.TYPE_CLAN
		case "L":
			id.Instance |= CHAT_INSTANCE_FLAGS_LOBBY
			id.Type = types.TYPE_CHAT
		default:
			id.Type = types.Parse(typeChar)
		}
	} else {
		return nil, fmt.Errorf("unknown SteamID input format %s", input)
	}

	return &id, nil
}

func (id *SteamId) FromIndividualAccountID(accountid int64) *SteamId {
	id.Universe = UNIVERSE_PUBLIC
	id.Type = types.TYPE_INDIVIDUAL
	id.Instance = INSTANCE_DESKTOP
	id.AccountId = accountid
	id.Token = ""

	return id
}

func (id *SteamId) FromTradeUrl(url string) error {
	if matches := regexp.MustCompile(`^https://steamcommunity\.com/tradeoffer/new/\?partner=(\d+)&token=([\w-]+)$`).FindStringSubmatch(url); matches != nil {
		accountid, _ := strconv.ParseInt(matches[1], 10, 64)
		token := matches[2]
		id.FromIndividualAccountID(accountid)
		id.Token = token
		return nil
	} else {
		return errors.New("invalid trade url")
	}
}

/**
 * Returns whether Steam would consider a given ID to be "valid".
 * This does not check whether the given ID belongs to a real account, nor does it check that the given ID is for
 * an individual account or in the public universe.
 * @returns {boolean}
 */
func (id SteamId) IsValid() bool {
	if id.Type <= types.TYPE_INVALID || id.Type > types.TYPE_ANON_USER {
		return false
	}

	if id.Universe <= UNIVERSE_INVALID || id.Universe > UNIVERSE_DEV {
		return false
	}

	if id.Type == types.TYPE_INDIVIDUAL && (id.AccountId == 0 || id.Instance > INSTANCE_WEB) {
		return false
	}

	if id.Type == types.TYPE_CLAN && (id.AccountId == 0 || id.Instance != INSTANCE_ALL) {
		return false
	}

	if id.Type == types.TYPE_GAMESERVER && id.AccountId == 0 {
		return false
	}

	return true
}

/**
 * Returns whether this SteamID is valid and belongs to an individual user in the public universe with a desktop instance.
 * This is what most people think of when they think of a SteamID. Does not check whether the account actually exists.
 * @returns {boolean}
 */
func (id SteamId) IsValidIndividual() bool {
	return id.Universe == UNIVERSE_PUBLIC &&
		id.Type == types.TYPE_INDIVIDUAL &&
		id.Instance == INSTANCE_DESKTOP && id.IsValid()
}

/**
 * Checks whether the given ID is for a legacy group chat.
 * @returns {boolean}
 */
func (id SteamId) IsGroupChat() bool {
	return id.Type == types.TYPE_CHAT && (id.Instance&CHAT_INSTANCE_FLAGS_CLAN) > 0
}

/**
 * Check whether the given Id is for a game lobby.
 * @returns {boolean}
 */
func (id SteamId) IsLobby() bool {
	return id.Type == types.TYPE_CHAT && (id.Instance&CHAT_INSTANCE_FLAGS_LOBBY > 0 || id.Instance&CHAT_INSTANCE_FLAGS_MMSLOBBY > 0)
}

func (id SteamId) GetSteam2RenderedID(newFormats ...bool) string {
	var newFormat bool

	if len(newFormats) > 0 {
		newFormat = newFormats[0]
	} else {
		newFormat = false
	}

	if id.Type != types.TYPE_INDIVIDUAL {
		return ""
	} else {
		universe := id.Universe
		if !newFormat && universe == 1 {
			universe = 0
		}

		return fmt.Sprintf("STEAM_%d:%d:%d", universe, id.AccountId&1, int64(math.Floor(float64(id.AccountId)/2)))
	}
}

func (id SteamId) GetSteam3RenderedID() string {
	typeChar := id.Type.String()

	if id.Instance&CHAT_INSTANCE_FLAGS_CLAN > 0 {
		typeChar = "c"
	} else if id.Instance&CHAT_INSTANCE_FLAGS_LOBBY > 0 {
		typeChar = "L"
	}

	shouldRenderInstance := id.Type == types.TYPE_ANON_GAMESERVER || id.Type == types.TYPE_MULTISEAT || (id.Type == types.TYPE_INDIVIDUAL && id.Instance != INSTANCE_DESKTOP)

	if shouldRenderInstance {
		return fmt.Sprintf("[%s:%d:%d:%d]", typeChar, id.Universe, id.AccountId, id.Instance)
	} else {
		return fmt.Sprintf("[%s:%d:%d]", typeChar, id.Universe, id.AccountId)
	}
}

func (id SteamId) GetBigIntId() int64 {
	universe := id.Universe << 56
	typeTmp := int64(id.Type) << 52
	instance := id.Instance << 32
	accountId := id.AccountId

	return universe | typeTmp | instance | accountId
}

func (id SteamId) String() string {
	return strconv.FormatInt(id.GetBigIntId(), 10)
}
func (id SteamId) GetSteamID64() string {
	return id.String()
}
