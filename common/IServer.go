package common

import (
	"github.com/MattSScott/basePlatformSOMAS/v2/pkg/agent"
	"github.com/google/uuid"
)

type IServer interface {
	agent.IExposedServerFunctions[IExtendedAgent]
	// Team management functions
	CreateTeam()
	AddAgentToTeam(agentID uuid.UUID, teamID uuid.UUID)
	GetAgentsInTeam(teamID uuid.UUID) []uuid.UUID
	CheckAgentAlreadyInTeam(agentID uuid.UUID) bool
	CreateAndInitTeamWithAgents(agentIDs []uuid.UUID) uuid.UUID
	UpdateAndGetAgentExposedInfo() []ExposedAgentInfo
	StartAgentTeamForming()
	UpdateCommonPools()

	GetTeam(agentID uuid.UUID) Team
	ProcessWithdrawals() // Added this method

	// Debug functions
	LogAgentStatus()
}