package main

/*
* Code to test the functionality of the orphan pool, which deals with agents
* that are not currently part of a team re-joining teams in subsequent turns.
*/

import (
	"SOMAS_Extended/agents"
	server "SOMAS_Extended/server"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert" // assert package, easier to
	"testing"                            // built-in go testing package
	"time"
)

/*
* Return a test server environment
*/
func CreateTestServer() *server.EnvironmentServer {
    // Default test config
    agentConfig := agents.AgentConfig{
		InitScore:    0,
		VerboseLevel: 10,
	}

	// Create a dummy server using the config
	serv := server.MakeEnvServer(2, 2, 3, 1000*time.Millisecond, 10, agentConfig)
    return serv
}

/* As it stands it would be difficult to simulate different scenarios where
* agents are accepted / rejected. To solve this problem we can introduce a mock
* agent, which will allow us to override specific functions for testing
* purposes. */
type MockAgent struct { 
    agents.ExtendedAgent 
    /* We introduce the concept of a blacklist, such that the agent can
    * deterministically accept / reject new agents from entering the team. This
    * is done for testing purposes only. The VoteOnAgentEntry function will
    * vote 'no' if the UUID is in the blacklist, and vote 'yes' otherwise. */
    Blacklist []uuid.UUID
}

/* 
* Override the VoteOnAgentEntry function to only accept if the agent is not in
* the blacklist
*/
func (mi *MockAgent) VoteOnAgentEntry(candidateID uuid.UUID) bool {
    // Iterate over all agents in blacklist and return false if UUID found
    for _, blacklistedID := range mi.Blacklist {
        if candidateID == blacklistedID {
            return false
        }
    }
    // return true otherwise
    return true
}

/*
* Allocation as it occurs on the BasePlatform, where the VoteOnAgentEntry()
* function returns true for every candidate ID
 */
func TestBaseAllocation(t *testing.T) {
	// Default Test Configuration
    serv := CreateTestServer()

    // Extract the list of agent IDs
	agentIDs := make([]uuid.UUID, 0)
	for id := range serv.GetAgentMap() {
		agentIDs = append(agentIDs, id)
	}

	teamID := serv.CreateAndInitTeamWithAgents(agentIDs)
	agents := serv.GetAgentsInTeam(teamID)
	assert.Equal(t, agentIDs, agents)

    // Go through all agents, and ask the team it is already in if it will
    // accept it. Yes this should technically never be asked, this is a base
    // case. 
	for agentID, agent := range serv.GetAgentMap() {
		teamID := agent.GetTeamID()
		accepted := serv.RequestOrphanEntry(agentID, teamID, 1.00)
		assert.Equal(t, true, accepted)
	}
}
