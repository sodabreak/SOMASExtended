package common

import "github.com/google/uuid"

type AuditRecord struct {
	auditMap map[uuid.UUID][]int
	duration int
	cost     int
	// reliability float64
}

// Getters
func (a *AuditRecord) GetAuditMap() map[uuid.UUID][]int {
	return a.auditMap
}

func (a *AuditRecord) GetAuditDuration() int {
	return a.duration
}

func (a *AuditRecord) GetAuditCost() int {
	return a.cost
}

// Get the number of infractions in the last n rounds, given by the quality of the audit
func (a *AuditRecord) GetAgentHistory(agentId uuid.UUID) int {
	infractions := 0
	records := a.auditMap[agentId]

	history := len(records)
	if history > a.duration {
		history = a.duration
	}

	for _, infraction := range records[len(records)-history:] {
		infractions += infraction
	}

	return infractions
}

// After a successful audit, clear the history of the agent so that there are
// no repeated warnings for the same infraction (this may change if using probabilistic auditing as well)
func (a *AuditRecord) ClearAgentHistory(agentId uuid.UUID) {
	a.auditMap[agentId] = []int{}
}

// After an agent's contribution, add a new record to the audit map
func (a *AuditRecord) AddAgentRecord(agentId uuid.UUID, roundInfractions int) {
	if _, ok := a.auditMap[agentId]; !ok {
		a.auditMap[agentId] = []int{}
	}

	a.auditMap[agentId] = append(a.auditMap[agentId], roundInfractions)
}

// After the agent's withdrawal, which is after the contribution, update the last record instead of adding a new one
func (a *AuditRecord) UpdateLastRecord(agentId uuid.UUID, extraInfractions int) {
	if _, ok := a.auditMap[agentId]; !ok {
		a.auditMap[agentId] = []int{}
	}

	records := a.auditMap[agentId]
	records[len(records)-1] += extraInfractions
}

func NewAuditRecord(duration int) *AuditRecord {
	cost := duration // For now, this can change depending on the tiering system

	return &AuditRecord{
		auditMap: make(map[uuid.UUID][]int),
		duration: duration,
		cost:     cost,
	}
}
