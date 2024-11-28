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
var pool = make(OrphanPoolType)

// The percentage of agents that have to vote 'accept' in order for an orphan
// to be taken into a team
const MajorityVoteThreshold float32 = 0.7 

/* 
* Print the contents of the pool. Careful as this will not necessarily print
* the elements in the order that you added them. 
*/
func (pool OrphanPoolType) Print() {
    for i, v := range pool {
        // truncate the UUIDs to make it easier to read
        shortAgentId := i.String()[:8]
        shortTeamIds := make([]string, len(v))

        // go over all the teams in the wishlist and add to shortened IDs
        for _, teamID := range v {
            shortTeamIds = append(shortTeamIds, teamID.String()[:8])
        }

        fmt.Println(shortAgentId, " Wants to join : ", shortTeamIds)
    }
}

/* 
* Ask all the agents in a team if they would be willing to accept an orphan
* into the team. This function accepts a threshold, that is used to determine
* whether to grant entry or not. For example, a threshold of 0.7 means that at
* least 70% of agents in the team have to be willing to accept the orphan. 
*/
func (cs *EnvironmentServer) RequestOrphanEntry(orphanID, teamID uuid.UUID, entryThreshold float32) bool {
    // Get the team and the current number of team members
    team := cs.GetTeam(teamID)
    agent_map := cs.GetAgentMap()

    num_members := len(team.Agents) 
    total_votes := 0

    // For each agent in the team
    for _, agentID := range team.Agents {
        // Get their vote
        vote := agent_map[agentID].VoteOnAgentEntry(agentID)
        // increment the total votes if they vote 'yes'
        if vote { total_votes++ }
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
    for orphanID, teamsList := range pool {
        // for each team that orphan wants to join
        for _, teamID := range teamsList {
            accepted := cs.RequestOrphanEntry(orphanID, teamID, MajorityVoteThreshold) 
            // If the team has voted to accept the orphan
            if accepted {
                orphan := agent_map[orphanID]
                orphan.SetTeamID(teamID) // Update agent's knowledge of its team
                cs.AddAgentToTeam(orphanID, teamID) // Update team's knowledge of its agents
                delete(pool, orphanID) // remove the agent from the orphan pool
            }
            // Otherwise, continue to the next team in the preference list. 
        }

        unallocated[orphanID] = teamsList // add to unallocated
        fmt.Printf("%v remains in the orphan pool after allocation...\n", orphanID)
    }

    // Assign the unallocated pool as the new orphan pool.
    pool = unallocated
}
