package gameRecorder

import (
	"fmt"
	"os"
	"sort"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
)

// CreatePlaybackHTML generates visualizations for the recorded game data
func CreatePlaybackHTML(recorder *ServerDataRecorder) {
	createScorePlots(recorder)
}

// createScorePlots creates line charts showing agent scores over time for each iteration
func createScorePlots(recorder *ServerDataRecorder) {
	// Add safety check at the start
	if len(recorder.TurnRecords) == 0 {
		fmt.Println("Warning: No turn records to visualize")
		return
	}

	// Group turn records by iteration
	iterationMap := make(map[int][]TurnRecord)
	for _, record := range recorder.TurnRecords {
		iterationMap[record.IterationNumber] = append(iterationMap[record.IterationNumber], record)
	}

	// Create a page to hold all iteration charts
	page := components.NewPage()
	page.PageTitle = "Agent Scores Per Iteration"

	// For each iteration, create a line chart
	for iteration, turns := range iterationMap {
		// Sort turns by turn number to ensure correct order
		sort.Slice(turns, func(i, j int) bool {
			return turns[i].TurnNumber < turns[j].TurnNumber
		})

		// Find first turn with agent records to initialize our agent map
		var initialAgentRecords []AgentRecord
		for _, turn := range turns {
			if len(turn.AgentRecords) > 0 {
				initialAgentRecords = turn.AgentRecords
				break
			}
		}

		if len(initialAgentRecords) == 0 {
			fmt.Printf("Warning: No agent records found in iteration %d\n", iteration)
			continue
		}

		// Create a new line chart with adjusted layout
		line := charts.NewLine()
		line.SetGlobalOptions(
			charts.WithTitleOpts(opts.Title{
				Title: fmt.Sprintf("Iteration %d - Agent Scores over Time", iteration),
				Top:   "5%", // Move title to top with some padding
			}),
			charts.WithTooltipOpts(opts.Tooltip{
				Show: opts.Bool(true),
			}),
			charts.WithLegendOpts(opts.Legend{
				Show: opts.Bool(true),
				Top:  "15%", // Move legend below title
			}),
			charts.WithXAxisOpts(opts.XAxis{
				Name:    "Turn Number",
				NameGap: 30, // Add gap between axis and name
				AxisLabel: &opts.AxisLabel{
					Show:   opts.Bool(true),
					Margin: 20, // Add margin to labels
				},
			}),
			charts.WithYAxisOpts(opts.YAxis{
				Name:    "Score",
				NameGap: 30, // Add gap between axis and name
			}),
			charts.WithGridOpts(opts.Grid{
				Top:    "40%", // Add more space at top for title and legend
				Right:  "5%",
				Left:   "10%", // Add more space for Y-axis labels
				Bottom: "15%", // Add more space for X-axis labels
			}),
		)

		// Get turn numbers for X-axis
		xAxis := make([]int, len(turns))
		for i, turn := range turns {
			xAxis[i] = turn.TurnNumber
		}

		// Create a map of agent scores over turns
		agentScores := make(map[string][]float64)

		// Initialize the map with empty slices using the first found agents
		for _, agent := range initialAgentRecords {
			agentID := agent.AgentID.String()
			agentScores[agentID] = make([]float64, len(turns))
		}

		// Fill in the scores
		for turnIdx, turn := range turns {
			for _, agent := range turn.AgentRecords {
				agentID := agent.AgentID.String()
				agentScores[agentID][turnIdx] = float64(agent.Score)
			}
		}

		// Add each agent's data series to the chart
		for agentID, scores := range agentScores {
			line.AddSeries(agentID, generateLineItems(xAxis, scores))
		}

		// Set X-axis data
		line.SetXAxis(xAxis)

		// Add the chart to the page
		page.AddCharts(line)
	}

	// Create the output file
	f, err := os.Create("agent_scores.html")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// Render the page
	err = page.Render(f)
	if err != nil {
		panic(err)
	}
}

// Helper function to generate line chart items
func generateLineItems(xAxis []int, yAxis []float64) []opts.LineData {
	items := make([]opts.LineData, len(xAxis))
	for i := 0; i < len(xAxis); i++ {
		items[i] = opts.LineData{Value: yAxis[i]}
	}
	return items
}
