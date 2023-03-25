package LinusUser

type TheCommandID int
type Img []byte

type User struct {
	ChatID             int
	Name               string
	Description        string
	SkillsString       string
	SkillsMap          map[string]bool
	YearsOfProgramming int
	Image              Img
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
