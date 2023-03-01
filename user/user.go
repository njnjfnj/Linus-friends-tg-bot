package user

type TheCommandID int
type Img []byte

type User struct {
	ChatID             int
	Name               string
	Description        string
	Skills             string
	YearsOfProgramming int

	LastCommand TheCommandID
	IsImportant bool

	Image Img
}

const (
	CmdStart = iota
	CmdMyProfile
	CmdChangeProfile
	CmdChangeProfilePhoto
	CmdChangeProfileText
	CmdChangeProfileName
	CmdSearching
	CmdShowMatches
	CmdStopShowingMyProfile
)
