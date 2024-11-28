package main

/*
* Code to test the functionality of the orphan pool, which deals with agents
* that are not currently part of a team re-joining teams in subsequent turns.
 */

import (
	"SOMAS_Extended/agents"
	"github.com/google/uuid"
	server "SOMAS_Extended/server"
	"github.com/stretchr/testify/assert" // assert package, easier to
	"testing"                            // built-in go testing package
	"time"
)

/*
* Allocation as it occurs on the BasePlatform, where the VoteOnAgentEntry()
* function returns true for every candidate ID
 */
func TestBaseAllocation(t *testing.T) {
	// Default Test Configuration
	agentConfig := agents.AgentConfig{
		InitScore:    0,
		VerboseLevel: 10,
	}

	// Create a dummy server using the config
	serv := server.MakeEnvServer(2, 2, 3, 1000*time.Millisecond, 10, agentConfig)
    // Extract the list of agent IDs
    agentIDs := make([]uuid.UUID, 0)
    for id := range serv.GetAgentMap() {
        agentIDs = append(agentIDs, id) 
    }

    teamID := serv.CreateAndInitTeamWithAgents(agentIDs)
    agents := serv.GetAgentsInTeam(teamID)
    assert.Equal(t, agentIDs, agents)

	// Ask every team if it will accept every agent
    for agentID, agent := range serv.GetAgentMap() {
        teamID := agent.GetTeamID() 
        accepted := serv.RequestOrphanEntry(agentID, teamID, 1.00)
        assert.Equal(t, true, accepted)
    }
}
