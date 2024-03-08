package types

type Type uint64

const (
	TYPE_INVALID Type = iota
	TYPE_INDIVIDUAL
	TYPE_MULTISEAT
	TYPE_GAMESERVER
	TYPE_ANON_GAMESERVER
	TYPE_PENDING
	TYPE_CONTENT_SERVER
	TYPE_CLAN
	TYPE_CHAT
	TYPE_P2P_SUPER_SEEDER
	TYPE_ANON_USER
)

var TypeChars = map[Type]string{
	TYPE_INVALID:         "I",
	TYPE_INDIVIDUAL:      "U",
	TYPE_MULTISEAT:       "M",
	TYPE_GAMESERVER:      "G",
	TYPE_ANON_GAMESERVER: "A",
	TYPE_PENDING:         "P",
	TYPE_CONTENT_SERVER:  "C",
	TYPE_CLAN:            "g",
	TYPE_CHAT:            "T",
	TYPE_ANON_USER:       "a",
}

func (t Type) String() string {
	char, ok := TypeChars[t]

	if !ok {
		char = "i"
	}

	return char
}

func Parse(char string) Type {
	for k, v := range TypeChars {
		if v == char {
			return k
		}
	}

	return TYPE_INVALID
}
