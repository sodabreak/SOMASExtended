package aoa

import "github.com/google/uuid"

type FixedAoA struct {
	AuditMap map[uuid.UUID]
}

func (f *FixedAoA) ResetAuditMap() {}

func (f *FixedAoA) GetExpectedContribution(agentId uuid.UUID, agentScore int) int {
	return agentScore
}

func (f *FixedAoA) SetContributionAuditResult(agentId uuid.UUID, agentScore int, agentActualContribution int, agentStatedContribution int) {}

func (f *FixedAoA) GetExpectedWithdrawal(agentId uuid.UUID, agentScore int) int {
	return 2
}

func (f *FixedAoA) SetWithdrawalAuditResult(agentId uuid.UUID, agentScore int, agentActualWithdrawal int, agentStatedWithdrawal int) {}

func (f *FixedAoA) GetAuditCost(commonPool int) int {
	return 0
}

func (f *FixedAoA) GetVoteResult(votes []Vote) *uuid.UUID {
	return nil
}

func CreateFixedAoA() IArticlesOfAssociation {
	return &FixedAoA{}
}

