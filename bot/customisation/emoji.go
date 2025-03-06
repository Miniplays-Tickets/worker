package customisation

import (
	"fmt"

	"github.com/rxdn/gdl/objects"
	"github.com/rxdn/gdl/objects/guild/emoji"
)

type CustomEmoji struct {
	Name     string
	Id       uint64
	Animated bool
}

func NewCustomEmoji(name string, id uint64, animated bool) CustomEmoji {
	return CustomEmoji{
		Name: name,
		Id:   id,
	}
}

func (e CustomEmoji) String() string {
	if e.Animated {
		return fmt.Sprintf("<a:%s:%d>", e.Name, e.Id)
	} else {
		return fmt.Sprintf("<:%s:%d>", e.Name, e.Id)
	}
}

func (e CustomEmoji) BuildEmoji() *emoji.Emoji {
	return &emoji.Emoji{
		Id:       objects.NewNullableSnowflake(e.Id),
		Name:     e.Name,
		Animated: e.Animated,
	}
}

var (
	EmojiId         = NewCustomEmoji("id", 1330141985008521239, false)
	EmojiOpen       = NewCustomEmoji("open", 1330141986883112980, false)
	EmojiOpenTime   = NewCustomEmoji("opentime", 1330141990435950592, false)
	EmojiClose      = NewCustomEmoji("close", 1330141988695179306, false)
	EmojiCloseTime  = NewCustomEmoji("closetime", 1330154075701907466, false)
	EmojiReason     = NewCustomEmoji("reason", 1330141993736732712, false)
	EmojiSubject    = NewCustomEmoji("subject", 1330154073810534494, false)
	EmojiTranscript = NewCustomEmoji("transcript", 1330154072048795729, false)
	EmojiClaim      = NewCustomEmoji("claim", 1330141992390361129, false)
	EmojiPanel      = NewCustomEmoji("panel", 1330141995129245706, false)
	EmojiRating     = NewCustomEmoji("rating", 1330154070433861653, false)
	EmojiStaff      = NewCustomEmoji("staff", 1330141996701978686, false)
	EmojiThread     = NewCustomEmoji("thread", 1330154069037420574, false)
	EmojiBulletLine = NewCustomEmoji("bulletline", 1330142310062882856, false)
	EmojiPatreon    = NewCustomEmoji("patreon", 1330154067552632952, false)
	EmojiDiscord    = NewCustomEmoji("discord", 1330154065921048667, false)
	//EmojiTime       = NewCustomEmoji("time", 974006684622159952, false)
)

// PrefixWithEmoji Useful for whitelabel bots
func PrefixWithEmoji(s string, emoji CustomEmoji, includeEmoji bool) string {
	if includeEmoji {
		return fmt.Sprintf("%s %s", emoji, s)
	} else {
		return s
	}
}
