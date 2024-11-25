package aoa

//---------------- Articles of Association ---------------//

type FixedArticlesOfAssociation struct {
	contributionRule IFixedContributionRule
	withdrawalRule   IFixedWithdrawalRule
	auditCost        IFixedAuditCost
	punishment       IFixedPunishment
}

func CreateFixedArticlesOfAssociation(contributionRule IFixedContributionRule, withdrawalRule IFixedWithdrawalRule, auditCost IFixedAuditCost, punishment IFixedPunishment) *FixedArticlesOfAssociation {
	return &FixedArticlesOfAssociation{
		contributionRule: contributionRule,
		withdrawalRule:   withdrawalRule,
		auditCost:        auditCost,
		punishment:       punishment,
	}
}

//--------------- Contribution Strategies ---------------//

type IFixedContributionRule interface {
	GetExpectedContributionAmount(agentScore int) int
	SetContributionAmount(amount int)
}

type FixedContributionRule struct {
	contributionAmount int
}

func (f *FixedContributionRule) GetExpectedContributionAmount(agentScore int) int {
	// Agent score can be used if this were percentage based
	return f.contributionAmount
}

// Can be removed if we want to keep it fixed in future implementations
func (f *FixedContributionRule) SetContributionAmount(amount int) {
	f.contributionAmount = amount
}

func CreateFixedContributionRule(amount int) IFixedContributionRule {
	return &FixedContributionRule{
		contributionAmount: amount,
	}
}

//--------------- Withdrawal Strategies ---------------//

type IFixedWithdrawalRule interface {
	GetExpectedWithdrawalAmount(agentScore int) int
	SetWithdrawalAmount(amount int) // This can be removed or changed depending on future extensions
	// An extension could be to treat the withdrawal amount as a percentage of the agent score, could add the common pool to this as well maybe
}

type FixedWithdrawalRule struct {
	withdrawalAmount int
}

func (f *FixedWithdrawalRule) GetExpectedWithdrawalAmount(agentScore int) int {
	return f.withdrawalAmount
}

// Can be removed if we want to keep it fixed in future implementations
func (f *FixedWithdrawalRule) SetWithdrawalAmount(amount int) {
	f.withdrawalAmount = amount
}

func CreateFixedWithdrawalRule(amount int) IFixedWithdrawalRule {
	return &FixedWithdrawalRule{
		withdrawalAmount: amount,
	}
}

//--------------- Audit Strategies ---------------//

type IFixedAuditCost interface {
	GetAuditCost() int
	SetAuditCost(cost int) // This can be removed or changed depending on future extensions
}

type FixedAuditCost struct {
	auditCost int
}

func (f *FixedAuditCost) GetAuditCost() int {
	return f.auditCost
}

// Can be removed if we want to keep it fixed in future implementations
func (f *FixedAuditCost) SetAuditCost(cost int) {
	f.auditCost = cost
}

func CreateFixedAuditCost(cost int) IFixedAuditCost {
	return &FixedAuditCost{
		auditCost: cost,
	}
}

//--------------- Punishment Strategies ---------------//

type IFixedPunishment interface {
	GetPunishment() int
	SetPunishment(punishment int) // This can be removed or changed depending on future extensions
}

type FixedPunishment struct {
	punishment int
}

func (f *FixedPunishment) GetPunishment() int {
	return f.punishment
}

// Can be removed if we want to keep it fixed in future implementations
func (f *FixedPunishment) SetPunishment(punishment int) {
	f.punishment = punishment
}

func CreateFixedPunishment(punishment int) IFixedPunishment {
	return &FixedPunishment{
		punishment: punishment,
	}
}
