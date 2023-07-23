package LinusUser

type User struct {
	ChatID             int
	Name               string
	Description        string
	SkillsString       string
	YearsOfProgramming int
	Image              []byte

	UserSeen map[int]bool
}

// const (
// 	CmdStart = iota
// 	CmdMyProfile
// 	CmdChangeProfile
// 	CmdChangeProfilePhoto
// 	CmdChangeProfileText
// 	CmdChangeProfileName
// 	CmdSearching
// 	CmdShowMatches
// 	CmdStopShowingMyProfile
// )
//SkillsMap          map[string]bool
