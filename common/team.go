package common

import (
	aoa "SOMAS_Extended/ArticlesOfAssociation"

	"github.com/google/uuid"
)

type Team struct {
	TeamID     	uuid.UUID
	CommonPool 	int
	Agents     	[]uuid.UUID
	AuditResult map[uuid.UUID]bool // Default is false, which means if false then there is no deferral
	TeamAoA 	aoa.IArticlesOfAssociation
}

// constructor: NewTeam creates a new Team with a unique TeamID and initializes other fields as blank.
func NewTeam() Team {
	teamAoA := aoa.CreateFixedAoA()
	return Team{
		TeamID:     	uuid.New(),             // Generate a unique TeamID
		CommonPool: 	0,                      // Initialize commonPool to 0
		Agents:     	[]uuid.UUID{},          // Initialize an empty slice of agent UUIDs
		AuditResult:	map[uuid.UUID]bool{},  // Initialize an empty map of agentID -> bool
		TeamAoA: teamAoA,   // Initialize strategy as 0
	}
}
