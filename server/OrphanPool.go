package environmentServer

import (
	"fmt"
	"github.com/google/uuid"
)

/* Declare the orphan pool for keeping track of agents that are not currently
* part of a team. This maps agentID -> slice of teamIDs that agent wants to
* join. Note that the slice of teams is processed in order, so the agent should
* put the team it most wants to join at the start of the slice. */
type OrphanPoolType map[uuid.UUID][]uuid.UUID

// The percentage of agents that have to vote 'accept' in order for an orphan
// to be taken into a team
const MajorityVoteThreshold float32 = 0.7

/*
* Ask all the agents in a team if they would be willing to accept an orphan
* into the team. This function accepts a threshold, that is used to determine
* whether to grant entry or not. For example, a threshold of 0.7 means that at
* least 70% of agents in the team have to be willing to accept the orphan.
*
* There is no logic in this function to check for the case where the agent is
* already in the team, this is not the responsibility of this function. It
* should not happen if the orphan pool is correctly managed.
 */
func (cs *EnvironmentServer) RequestOrphanEntry(orphanID, teamID uuid.UUID, entryThreshold float32) bool {
	// Get the team and the current number of team members
	team := cs.GetTeamFromTeamID(teamID)
	agent_map := cs.GetAgentMap()

	num_members := len(team.Agents)
	total_votes := 0

	// For each agent in the team
	for _, agentID := range team.Agents {
		// Get their vote
		vote := agent_map[agentID].VoteOnAgentEntry(orphanID)
		// increment the total votes if they vote 'yes'
		if vote {
			total_votes++
		}
	}

	// Calculate the acceptance ratio and return 'yes' only if enough of the
	// team has voted to accept.
	acceptance := float32(total_votes) / float32(num_members)
	return (acceptance >= entryThreshold)
}

/*
* Go through the pool and attempt to allocate each of the orphans to a team,
* based on the preference they have expressed.
 */
func (cs *EnvironmentServer) AllocateOrphans() {
	agent_map := cs.GetAgentMap()

	// Create pool for keeping track of which orphans have not been allocated.
	// This is because we want to only keep the unallocated orphans in the
	// pool. This is the safer alternative to deleting from the pool as we are
	// iterating through it.
	unallocated := make(OrphanPoolType)

	// for each orphan currently in the pool / shelter
	for orphanID, teamsList := range cs.orphanPool {
		// for each team that orphan wants to join
		for _, teamID := range teamsList {
			accepted := cs.RequestOrphanEntry(orphanID, teamID, MajorityVoteThreshold)
			// If the team has voted to accept the orphan
			if accepted {
				orphan := agent_map[orphanID]
				orphan.SetTeamID(teamID)            // Update agent's knowledge of its team
				cs.AddAgentToTeam(orphanID, teamID) // Update team's knowledge of its agents
				delete(cs.orphanPool, orphanID)     // remove the agent from the orphan pool
			}
			// Otherwise, continue to the next team in the preference list.
		}

		unallocated[orphanID] = teamsList // add to unallocated
		fmt.Printf("%v remains in the orphan pool after allocation...\n", orphanID)
	}

	// Assign the unallocated pool as the new orphan pool.
	cs.orphanPool = unallocated
}

/*
* Go over all the agents in the agent map. If there is an agent that is not
* part of a team, then add it to the orphan pool. This allows the server to
* actively pick up agents that have been removed from a team or that have left.
* This prevents agents from having to tell the server to 'make me an orphan'. 
*
* This will not break for dead agents, because dead agents should be in a
* separate map. (deadAgents)
*/
func (cs *EnvironmentServer) PickUpOrphans() {
    // sweep over all the agents in the server's agent map
    for agentID := range cs.GetAgentMap() {
        // if the agent does not belong to a team, and is not in the orphan
        // pool already, then add it to the orphan pool. 
        _, exists := cs.orphanPool[agentID]

        if !exists {
		    fmt.Printf("%v was added to the orphan pool \n", agentID)
            cs.orphanPool[agentID] = make([]uuid.UUID, 0) 
        }
    }
}

