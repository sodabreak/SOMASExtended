package aoa

import "github.com/google/uuid"

type Vote struct {
	IsVote bool
	VoterID uuid.UUID
	VotedForID uuid.UUID
}

type IArticlesOfAssociation interface {
	ResetAuditMap()
	GetExpectedContribution(agentId uuid.UUID, agentScore int) int
	SetContributionAuditResult(agentId uuid.UUID, agentScore int, agentActualContribution int, agentStatedContribution int)
	GetExpectedWithdrawal(agentId uuid.UUID, agentScore int) int
	SetWithdrawalAuditResult(agentId uuid.UUID, agentScore int, agentActualWithdrawal int, agentStatedWithdrawal int)
	GetAuditCost(commonPool int) int
	GetVoteResult(votes []Vote) *uuid.UUID
}

func CreateVote(isVote bool, voterId uuid.UUID, votedForId uuid.UUID) Vote {
	return Vote{
		IsVote:  isVote,
		VoterID: voterId,
		VotedForID: votedForId,
	}
}
