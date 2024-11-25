package aoa

import "github.com/google/uuid"

type FixedAoA struct {}

func (f *FixedAoA) ResetAuditMap() {}

func (f *FixedAoA) GetExpectedContribution(agentId uuid.UUID, agentScore int) int {
	return agentScore
}

func (f *FixedAoA) SetContributionResult(agentId uuid.UUID, agentScore int, agentContribution int) {}

func (f *FixedAoA) GetExpectedWithdrawal(agentId uuid.UUID, agentScore int) int {
	return agentScore
}

func (f *FixedAoA) SetWithdrawalResult(agentId uuid.UUID, agentScore int, agentWithdrawal int) {}

func (f *FixedAoA) GetAuditCost(commonPool int) int {
	return 0
}

func (f *FixedAoA) GetVoteResult(votes []Vote) *uuid.UUID {
	return nil
}

func CreateFixedAoA() IArticlesOfAssociation {
	return &FixedAoA{}
}

