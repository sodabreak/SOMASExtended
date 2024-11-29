package main

/*
* Code to test the functionality of the orphan pool, which deals with agents
* that are not currently part of a team re-joining teams in subsequent turns.
 */

import (
	"SOMAS_Extended/agents"
	server "SOMAS_Extended/server"
	"testing" // built-in go testing package
	"time"
    "reflect"
    "fmt"

	"bou.ke/monkey"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert" // assert package, easier to
)

/*
* Return a test server environment
 */
func CreateTestServer() (*server.EnvironmentServer, []uuid.UUID) {
    // Default test config
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

    return serv, agentIDs
}

/* Define Mock functions for the VoteOnAgentEntry function. These will override
* the base implementation to test different voting logic. Yes, monkeypatching
* is not particularly safe practice with Go but seeing as there is literally no
* way to override the base platform during the game instantiation without
* putting test code inside the base platform, this is the best I can think of.
* The fact that it is so hard to override the base platfrom from outside the
* base platform is an issue we will probably encounter later. */

/* 
* Behaviour 1 - Override the VoteOnAgentEntry function to always vote 'no'
*/
func mockVoteAlwaysFalse(mi *agents.ExtendedAgent, candidateID uuid.UUID) bool {
    return false
}


/*
* Allocation as it occurs on the BasePlatform, where the VoteOnAgentEntry()
* function returns true for every candidate ID
 */
func TestBaseAllocation(t *testing.T) {
	// Default Test Configuration
    serv, agentIDs := CreateTestServer()
    
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

/*
* Test if all the orphans are rejected when all agents vote 'false' to accept
*/
func TestAlwaysReject(t *testing.T) {
    // Create a test server and put all the agents in the same team
    serv, agentIDs := CreateTestServer()
	serv.CreateAndInitTeamWithAgents(agentIDs)

    // Monkey-path the voting method to always reject
    monkey.PatchInstanceMethod(reflect.TypeOf(&agents.ExtendedAgent{}), "VoteOnAgentEntry", mockVoteAlwaysFalse)
    defer monkey.UnpatchAll()

    // Go through all the agents, every agent should be rejected. 
    for agentID, agent := range serv.GetAgentMap() {
		teamID := agent.GetTeamID()
		accepted := serv.RequestOrphanEntry(agentID, teamID, 1.00)
		assert.Equal(t, false, accepted)
	}
}

/*
* Test basic team-based voting logic. In this example, agents will vote to
* 'accept' an orphan if that 'orphan' is already in that team. In theory this
* would not happen since an orphan should not be in that team. Just a basic
* unit test
*/
func TestAcceptIfInCurrentTeam(t *testing.T) {
    // Create a test server and put all the agents in the same team
    serv, agentIDs := CreateTestServer()
    
    // Create two sub-teams, each with half the agents. 
    numAgents := len(agentIDs)
    team1ID := serv.CreateAndInitTeamWithAgents(agentIDs[:(numAgents / 2)])
    team2ID := serv.CreateAndInitTeamWithAgents(agentIDs[(numAgents / 2):])

    /* 
    * Behaviour 2 - Override the VoteOnAgentEntry function to vote yes only if the
    * agent is already in the current team. In order to be able to check what
    * team we're in, we define this as a closure to have access to the server
    * (you can't access the agent's private members inside the mock function
    * here...
    */
    mockVoteBasedOnTeam := func (mi *agents.ExtendedAgent, candidateID uuid.UUID) bool {
        myTeam := mi.GetTeamID()
        fmt.Printf("%v (votee) is in team %v\n", mi.GetID(), myTeam)
        candidateTeam := serv.GetTeam(candidateID).TeamID
        fmt.Printf("%v is in team %v\n", candidateID, candidateTeam)
        return (candidateTeam == myTeam)
    }

    // Monkey-path the voting method to only accept an agent if it is already in that team
    monkey.PatchInstanceMethod(reflect.TypeOf(&agents.ExtendedAgent{}), "VoteOnAgentEntry", mockVoteBasedOnTeam)
    defer monkey.UnpatchAll()

    // Go through all the agents. The agent should be accepted if it is already
    // in that team, otherwise it should be rejected. 
    for agentID, agent := range serv.GetAgentMap() {
		teamID := agent.GetTeamID()
        fmt.Printf("Inspecting %v\n", teamID)

        // Try to enter team1
		accepted := serv.RequestOrphanEntry(agentID, team1ID, 1.00)
		assert.Equal(t, (teamID == team1ID), accepted)

        // Try to enter team2
		accepted = serv.RequestOrphanEntry(agentID, team2ID, 1.00)
		assert.Equal(t, (teamID == team2ID), accepted)
	}
}
