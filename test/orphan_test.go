package main

/*
* Code to test the functionality of the orphan pool, which deals with agents
* that are not currently part of a team re-joining teams in subsequent turns.
 */

import (
	"SOMAS_Extended/agents"
	server "SOMAS_Extended/server"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert" // assert package, easier to
	"testing"                            // built-in go testing package
	"time"
)

/*
* Allocation as it occurs on the BasePlatform, where the VoteOnAgentEntry()
* function returns true for every candidate ID
 */
func TestBaseAllocation(t *testing.T) {
	// Default Test Configuration
	agentConfig := agents.AgentConfig{
		InitScore:    0,
		VerboseLevel: 10,
	}

	// Create a dummy server using the config
	serv := server.MakeEnvServer(2, 2, 3, 1000*time.Millisecond, 10, agentConfig)
	assert.NotNil(t, serv)

	// Ask all teams to accept any agent

}
