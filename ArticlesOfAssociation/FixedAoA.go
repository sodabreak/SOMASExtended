package aoa

import (
	"github.com/google/uuid"
	"math/rand"
	"time"
)

type FixedAoA struct{}

func (f *FixedAoA) ResetAuditMap() {}

func (f *FixedAoA) GetExpectedContribution(agentId uuid.UUID, agentScore int) int {
	return agentScore
}

func (f *FixedAoA) SetContributionAuditResult(agentId uuid.UUID, agentScore int, agentActualContribution int, agentStatedContribution int) {
}

func (f *FixedAoA) GetExpectedWithdrawal(agentId uuid.UUID, agentScore int) int {
	return 2
}

func (f *FixedAoA) SetWithdrawalAuditResult(agentId uuid.UUID, agentScore int, agentActualWithdrawal int, agentStatedWithdrawal int) {
}

func (f *FixedAoA) GetAuditCost(commonPool int) int {
	return 0
}

func (f *FixedAoA) GetVoteResult(votes []Vote) *uuid.UUID {
	return nil
}

func (t *FixedAoA) GetWithdrawalOrder(agentIDs []uuid.UUID) []uuid.UUID {
	// Seed the random number generator to ensure different shuffles each time
	rand.Seed(time.Now().UnixNano())

	// Create a copy of the agentIDs to avoid modifying the original list
	shuffledAgents := make([]uuid.UUID, len(agentIDs))
	copy(shuffledAgents, agentIDs)

	// Shuffle the agent list
	rand.Shuffle(len(shuffledAgents), func(i, j int) {
		shuffledAgents[i], shuffledAgents[j] = shuffledAgents[j], shuffledAgents[i]
	})

	return shuffledAgents
}

func CreateFixedAoA() IArticlesOfAssociation {
	return &FixedAoA{}
}
