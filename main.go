package main

import (
	"fmt"
	"time"

	"SOMAS_Extended/agents"
	envServer "SOMAS_Extended/server"
)

func main() {
	fmt.Println("main function started.")

	// agent configurations
	agentConfig := agents.AgentConfig{
		InitScore:    0,
		VerboseLevel: 10,
	}

	// note: the zero turn is used for team forming
	serv := envServer.MakeEnvServer(2,
		2,                    // agent num PER TEAM
		10,                   // turns
		3,                    // after this many turns, apply threshold
		100*time.Millisecond, // max duration
		10,                   // max thread
		agentConfig)

	//serv.ReportMessagingDiagnostics()
	serv.Start()

	// custom function to see agent result
	serv.LogAgentStatus()
	serv.LogTeamStatus()

	// // record data
	serv.DataRecorder.GamePlaybackSummary()
}
