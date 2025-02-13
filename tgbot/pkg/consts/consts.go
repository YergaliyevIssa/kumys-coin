package consts

import "time"

const (
	StateInSectionMain     = "in_main"
	StateInSectionAnalysis = "in_analysis"
	StateInSectionProfile  = "in_profile"
	StateChangingProfile   = "changing_profile"
)

const (
	UserSessionTTL = 100 * time.Hour * 24
)
