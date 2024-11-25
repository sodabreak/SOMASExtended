package aoa

import "github.com/google/uuid"

type Vote struct {
	IsVote bool
	VoterID uuid.UUID
	VotedForID uuid.UUID
}

type IArticlesOfAssociation interface {
	ResetAuditMap()
	SetContributionResult(agentId uuid.UUID, agentScore int, agentContribution int)
	SetWithdrawalResult(agentId uuid.UUID, agentScore int, agentWithdrawal int)
	GetAuditCost(commonPool int) int // to be removed from the common pool, called if successful
	GetVoteResult(votes []Vote) *uuid.UUID // nullable, if this isn't nil, then the team has votes for an agent
}

func CreateVote(isVote bool, voterId uuid.UUID, votedForId uuid.UUID) Vote {
	return Vote{
		IsVote:  isVote,
		VoterID: voterId,
		VotedForID: votedForId,
	}
}
