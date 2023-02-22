package user

type TheCommandID int
type TheLastPosition int
type Img []byte

const (
	// commands
	Start = iota
	ChangeProfile
	Searching
	ShowMatches
)

type User struct {
	ChatID      int
	Name        string
	Description string

	LastCommand TheCommandID
	LastPos     TheLastPosition

	Image Img
}
