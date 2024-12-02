package common

import (
	"container/list"
	"errors"
	"github.com/google/uuid"
	"math/rand"
	"sort"
	"log"
)

type Team1AoA struct {
	auditResult map[uuid.UUID]*list.List
	ranking     map[uuid.UUID]int
	threshold   int
}

func (t *Team1AoA) ResetAuditMap() {
	t.auditResult = make(map[uuid.UUID]*list.List)
}

func (t *Team1AoA) GetExpectedContribution(agentId uuid.UUID, agentScore int) int {
	return t.threshold // For now using threshold as minimum for all ranks, later have per rank minimums? But need to vote what is min?
}

func (t *Team1AoA) SetContributionAuditResult(agentId uuid.UUID, agentScore int, agentActualContribution int, agentStatedContribution int) {
	t.auditResult[agentId].PushBack((agentStatedContribution > agentActualContribution))

	// Just for our AoA we are updating rank based on stated
	t.ranking[agentId] += (agentStatedContribution / t.threshold) // Plus 1 rank every `threshold` points?

	// Cap agent rank at 5
	if t.ranking[agentId] > 5 {
		t.ranking[agentId] = 5
	}
	// Lower rank if insufficient contribution
	if agentActualContribution < t.threshold { // For now using threshold as minimum for all ranks, later have per rank minimums? But need to vote what is min?
		t.ranking[agentId] -= 1
	}
	if t.ranking[agentId] < 1 { // Ensure cant go below 1
		t.ranking[agentId] = 1
	}

}

func (t *Team1AoA) GetExpectedWithdrawal(agentId uuid.UUID, agentScore int, commonPool int) int {
	totalWeightedSum := 0
	for _, rank := range t.ranking {
		totalWeightedSum += rank
	}
	expectedWithdrawal := t.ranking[agentId] * (commonPool / (totalWeightedSum + 5)) // Weight 5 for pool?
	return expectedWithdrawal
}

func (t *Team1AoA) SetWithdrawalAuditResult(agentId uuid.UUID, agentScore int, agentActualWithdrawal int, agentStatedWithdrawal int, commonPool int) {
	t.auditResult[agentId].PushBack((agentActualWithdrawal > agentStatedWithdrawal) || (agentActualWithdrawal > t.GetExpectedWithdrawal(agentId, agentScore, commonPool)))
}

func (t *Team1AoA) GetAuditCost(commonPool int) int {
	// Need to get argument which agent being audited and then change cost?
	return 5
}

func (t *Team1AoA) GetVoteResult(votes []Vote) uuid.UUID {
	// Count total votes
	totalVotes := 0
	voteMap := make(map[uuid.UUID]int)
	highestVotes := -1
	highestVotedID := uuid.Nil
	for _, vote := range votes {
		totalVotes += vote.IsVote
		if vote.IsVote == 1 { // Should agents who didnt want to vote, get a vote if majority wants to?
			voteMap[vote.VotedForID]++
		}
		// Check if this ID has the highest votes
		if voteMap[vote.VotedForID] > highestVotes {
			highestVotedID = vote.VotedForID
			highestVotes = voteMap[vote.VotedForID]
		}
	}
	if totalVotes <= 0 {
		return uuid.Nil // Majority does not want to vote
	}
	return highestVotedID
}

func (t *Team1AoA) GetContributionAuditResult(agentId uuid.UUID) bool {
	return t.auditResult[agentId].Back().Value.(int) == 1
}

func (t *Team1AoA) GetWithdrawalAuditResult(agentId uuid.UUID) bool {
	return t.auditResult[agentId].Back().Value.(int) == 1
}

func (t *Team1AoA) GetWithdrawalOrder(agentIDs []uuid.UUID) []uuid.UUID {
	// Sort the agent based on their rank value in descending order
	sort.Slice(agentIDs, func(i, j int) bool {
		return t.ranking[agentIDs[i]] > t.ranking[agentIDs[j]]
	})
	return agentIDs
}

// WeightedRandomSelection selects one agent based on weights derived from ranks.
func (t *Team1AoA) WeightedRandomSelection(agentIds []uuid.UUID) uuid.UUID {
    if len(agentIds) == 0 {
        log.Fatal("No agents to select from")
    }

    totalWeight := 0
    for _, agentId := range agentIds {
        totalWeight += t.ranking[agentId]
    }
    if totalWeight == 0 {
        log.Fatal("All agents have 0 weight")
    }
    
    randomNumber := rand.Intn(totalWeight) + 1
    cumulativeWeight := 0
    for _, agentId := range agentIds {
        cumulativeWeight += t.ranking[agentId]
        if cumulativeWeight >= randomNumber {
            return agentId
        }
    }
    
    log.Fatal("Failed to select an agent")
    return uuid.Nil // This line will never be reached due to log.Fatal
}

// SelectNChairs selects n distinct agents to be chairs, with probability of selection based on rank.
func (t *Team1AoA) SelectNChairs(agentIds []uuid.UUID, n int) []uuid.UUID {
    if len(agentIds) < n {
        log.Fatal("not enough agents to select from")
    }

    selectedChairs := make([]uuid.UUID, 0, n)
    remainingAgents := make([]uuid.UUID, len(agentIds))
    copy(remainingAgents, agentIds)

    for i := 0; i < n; i++ {
        agent := t.WeightedRandomSelection(remainingAgents)
        selectedChairs = append(selectedChairs, agent)

        // Remove the selected agent from remainingAgents
        // Find the index of the selected agent
        index := -1
        for j, id := range remainingAgents {
            if id == agent {
                index = j
                break
            }
        }

        if index == -1 {
            log.Fatal("selected agent not found in remainingAgents")
        }

        // Remove the agent by swapping with the last element and truncating the slice
        remainingAgents[index] = remainingAgents[len(remainingAgents)-1]
        remainingAgents = remainingAgents[:len(remainingAgents)-1]
    }

    return selectedChairs
}

// TODO: move to AGENT
func (a *ExtendedAgent) ChairCountVotes(ENUM_THINGS_TO_VOTE_ON iota) T {
	switch (ENUM_THINGS_TO_VOTE_ON){
		case NEW_AGENT_RANK_PREFS:
			votes := 0
			for _, agent := range t.Agents {
				votes += agent.GetAgentRankPreferences()
			return votes
	}
}
}

func (t *Team1AoA) RunPostContributionAoaLogic(team *Team) {
	// Choose 2 chairs based on rank
	// call function for agents to vote on ranks
	// If the chairs decision do not match, then reduce rank by 1 of their score and give to common pool
	// Then repeat until two agents agree

	// Choose 2 chairs based on rank
	chairs := t.SelectNChairs(team.Agents, 2)

	listOfmaps := make([]map[uuid.UUID]int, 0)
	

	chair_1.ChairCountVotes(NEW_AGENT_RANK_PREFS)

	// get the pointers to chairs and call countVotes function to get the result
	
	



	return

}

func CreateTeam1AoA(team *Team) IArticlesOfAssociation {
	auditResult := make(map[uuid.UUID]*list.List)
	ranking := make(map[uuid.UUID]int)
	for _, agent := range team.Agents {
		auditResult[agent] = list.New()
		ranking[agent] = 1
	}

	return &Team1AoA{
		auditResult: auditResult,
		ranking:     ranking,
		threshold:   5,
	}
}
