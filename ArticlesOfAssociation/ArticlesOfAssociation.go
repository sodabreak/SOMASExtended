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
	SetContributionResult(agentId uuid.UUID, agentScore int, agentContribution int)
	GetExpectedWithdrawal(agentId uuid.UUID, agentScore int) int
	SetWithdrawalResult(agentId uuid.UUID, agentScore int, agentWithdrawal int)
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
