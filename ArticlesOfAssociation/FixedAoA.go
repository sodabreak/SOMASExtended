package aoa

import "github.com/google/uuid"

type FixedAoA struct {
	ContributionAuditMap map[uuid.UUID]bool
	WithdrawalAuditMap   map[uuid.UUID]bool
}

func (f *FixedAoA) ResetAuditMap() {}

func (f *FixedAoA) GetExpectedContribution(agentId uuid.UUID, agentScore int) int {
	return agentScore
}

func (f *FixedAoA) SetContributionAuditResult(agentId uuid.UUID, agentScore int, agentActualContribution int, agentStatedContribution int) {
}

func (f *FixedAoA) GetContributionAuditResult(agentId uuid.UUID) bool {
	// true means agent failed the audit (cheated)
	return f.ContributionAuditMap[agentId]
}

func (f *FixedAoA) GetExpectedWithdrawal(agentId uuid.UUID, agentScore int) int {
	return 2
}

func (f *FixedAoA) SetWithdrawalAuditResult(agentId uuid.UUID, agentScore int, agentActualWithdrawal int, agentStatedWithdrawal int) {
}

func (f *FixedAoA) GetWithdrawalAuditResult(agentId uuid.UUID) bool {
	// true means agent failed the audit (cheated)
	return f.WithdrawalAuditMap[agentId]
}

func (f *FixedAoA) GetAuditCost(commonPool int) int {
	return 0
}

// MUST return UUID nil if audit should not be executed
// Otherwise, implement a voting mechanism to determine the agent to be audited
// and return its UUID
func (f *FixedAoA) GetVoteResult(votes []Vote) uuid.UUID {
	return uuid.Nil
}

func CreateFixedAoA() IArticlesOfAssociation {
	return &FixedAoA{}
}
