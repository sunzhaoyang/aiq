package tool

import (
	"fmt"

	"github.com/aiq/aiq/internal/chart"
	"github.com/aiq/aiq/internal/db"
)

// RenderChartString renders query results as a chart string
// This function does NOT print anything - it only returns the chart output
func RenderChartString(result *db.QueryResult, chartTypeStr string) (string, error) {
	var chartType chart.ChartType
	switch chartTypeStr {
	case "bar":
		chartType = chart.ChartTypeBar
	case "line":
		chartType = chart.ChartTypeLine
	case "pie":
		chartType = chart.ChartTypePie
	case "scatter":
		chartType = chart.ChartTypeScatter
	default:
		return "", fmt.Errorf("unknown chart type: %s", chartTypeStr)
	}

	// Render chart
	config := chart.DefaultConfig()
	config.Width = 80
	config.Height = 20
	config.Title = fmt.Sprintf("Chart (%d rows)", len(result.Rows))

	chartOutput, err := chart.RenderChart(result, chartType, config)
	if err != nil {
		return "", fmt.Errorf("chart rendering failed: %w", err)
	}

	return chartOutput, nil
}
