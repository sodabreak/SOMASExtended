package environmentServer

import (
	"fmt"
	"github.com/google/uuid"
)

/* Declare the orphan pool for keeping track of agents that are not currently
* part of a team. This maps agentID -> slice of teamIDs that agent wants to
* join. Note that the slice of teams is processed in order, so the agent should
* put the team it most wants to join at the start of the slice. */
type orphanPool map[uuid.UUID][]uuid.UUID
var pool = make(orphanPool)

// The percentage of agents that have to vote 'accept' in order for an orphan
// to be taken into a team
const MajorityVoteThreshold float32 = 0.7 

/* 
* Print the contents of the pool. Careful as this will not necessarily print
* the elements in the order that you added them. 
*/
func (pool orphanPool) Print() {
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
* Ask all the agents in a team if they would be willing to accept an orphan into the team
*/
func (cs *EnvironmentServer) RequestOrphanEntry(orphanID, teamID uuid.UUID, entryThreshold) bool {
    // todo
    // todo
}

/*
* Return a function (closure) for counting up the total votes in a team. 
*/
func voteAdder() func(bool) int {
    // An instance of this 'votes' var will be created for each closure
    // created. Whenever that closure is called, it will increment its own copy
    // of votes! No need to keep data in another data structure. 
    votes := 0
    return func(agent_vote bool) int {
        // increment the total number of votes only if the agent returned
        // true, i.e. 'yes I will accept this member into the team'
        if agent_vote { votes += 1 }
        return votes
    }
}

/*
* Go through the pool and attempt to allocate each of the orphans to a team,
* based on the preference they have expressed.
*/
func (bds * BaseDiceServer) AllocateOrphans() {
    // for each orphan currently in the pool / shelter
    for orphan, teamsList := range pool {
        // for each team that orphan wants to join
        for _, teamID := range teamsList {
            // get the members of that team, if the team exists
            team, ok := bds.teams[teamID]
            if !ok { 
                // if the orphan is trying to join a team that does not exist,
                // make a not of this but do not kill the program. 
                fmt.Print("Orphan ", orphan, " tried to join team ", teamID, " which does not exist")
                continue 
            }
            // otherwise, add the vote of each member
            adder := voteAdder()
            members := team.GetTeamMembers()
            var votes int
            for _, agent := range members {
                votes = adder(GetAgentVoteFromId(orphan, agent))
            }

            // calculate the percentage of agents that voted yes
            acceptancePercentage := float32(votes) / float32(len(members))
            
            // add agent to team only if acceptance is above percentage
            if acceptancePercentage >= MajorityVoteThreshold {
                // add the orphan to the team
                team.AddMember(orphan) 
                // tell the orphan what team it now belongs to
                bds.GetAgentMap()[orphan].SetTeam(team) 
                // move onto the next orphan in the pool
                delete(pool, orphan)
                break
            }
        }
    }
}

// test code ---------------------------------------

func RunTests() {

    pool = orphanPool{
        uuid.New(): {uuid.New(), uuid.New(), uuid.New()},
        uuid.New(): {uuid.New(), uuid.New(), uuid.New(), uuid.New(), uuid.New()},
        uuid.New(): {}, // outcast, does not want to join any team. 
        uuid.New(): {uuid.New()}, 
    }

    pool.Print()
}
