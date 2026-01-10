package web

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"tmobile-stats/internal/analysis"
	"tmobile-stats/internal/charting"
)

const htmlTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta http-equiv="refresh" content="60">
    <title>Signal Sentry Live</title>
    <style>
        body { font-family: sans-serif; text-align: center; background: #f4f4f4; margin: 0; padding: 20px; }
        .nav { margin-bottom: 20px; }
        .nav a {
            display: inline-block;
            padding: 10px 20px;
            margin: 0 5px;
            background: #ddd;
            text-decoration: none;
            color: #333;
            border-radius: 5px;
            font-weight: bold;
        }
        .nav a.active { background: #333; color: #fff; }
        .nav a:hover { background: #bbb; }
        .chart-container { background: #fff; padding: 20px; border-radius: 10px; box-shadow: 0 0 10px rgba(0,0,0,0.1); display: inline-block; }
        img { max-width: 100%; height: auto; }
        .footer { margin-top: 20px; color: #777; font-size: 0.9em; }
    </style>
</head>
<body>
    <h1>Signal Sentry Live</h1>
    
    <div class="nav">
        {{range .Links}}
        <a href="/?range={{.Val}}" class="{{if .Active}}active{{end}}">{{.Label}}</a>
        {{end}}
    </div>

    <div class="chart-container">
        <img src="/chart.png?range={{.CurrentRange}}" alt="Signal Chart">
    </div>

    <div class="footer">
        Last updated: {{.LastUpdated}} | Auto-refreshing every 60s
    </div>
</body>
</html>
`

type Link struct {
	Label  string
	Val    string
	Active bool
}

type PageData struct {
	Links        []Link
	CurrentRange string
	Start        string
	End          string
	LastUpdated  string
}

func Run(port int, logFile string, quiet bool) error {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handleIndex(w, r, quiet)
	})

	mux.HandleFunc("/chart.png", func(w http.ResponseWriter, r *http.Request) {
		handleChart(w, r, logFile, quiet)
	})

	addr := fmt.Sprintf(":%d", port)
	if !quiet {
		log.Printf("Starting web server on http://localhost%s (Input: %s)", addr, logFile)
	}
	return http.ListenAndServe(addr, mux)
}

func handleIndex(w http.ResponseWriter, r *http.Request, quiet bool) {
	_, currentRange, err := parseTimeFilter(r)
	if err != nil {
		// Just log and continue with default if index fails parsing
		if !quiet {
			log.Printf("Filter parse error in index: %v", err)
		}
	}

	start := r.URL.Query().Get("start")
	end := r.URL.Query().Get("end")

	ranges := []struct {
		Label string
		Val   string
	}{
		{"10m", "10m"},
		{"1h", "1h"},
		{"6h", "6h"},
		{"24h", "24h"},
		{"Max", "0"},
	}

	links := make([]Link, len(ranges))
	for i, rng := range ranges {
		active := rng.Val == currentRange && start == "" && end == ""
		if currentRange == "0" && rng.Val == "0" && start == "" && end == "" {
			active = true
		}
		
		links[i] = Link{
			Label:  rng.Label,
			Val:    rng.Val,
			Active: active,
		}
	}

	data := PageData{
		Links:        links,
		CurrentRange: currentRange,
		Start:        start,
		End:          end,
		LastUpdated:  time.Now().Format("15:04:05"),
	}

	tmpl, err := template.New("index").Parse(htmlTemplate)
	if err != nil {
		http.Error(w, "Template Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	if err := tmpl.Execute(w, data); err != nil {
		if !quiet {
			log.Printf("Error executing template: %v", err)
		}
	}
}

func handleChart(w http.ResponseWriter, r *http.Request, logFile string, quiet bool) {
	filter, _, err := parseTimeFilter(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Filter Error: %v", err), http.StatusBadRequest)
		return
	}

	f, err := os.Open(logFile)
	if err != nil {
		http.Error(w, fmt.Sprintf("File Error: %v", err), http.StatusInternalServerError)
		return
	}
	defer f.Close()

	data, err := analysis.ParseLog(f, filter)
	if err != nil {
		http.Error(w, fmt.Sprintf("Parse Error: %v", err), http.StatusInternalServerError)
		return
	}

	if len(data) == 0 {
		// Create a blank image or return text?
		// Let's return text for now to debug
		// Or try to generate an empty chart? charting.GenerateToWriter handles len(0) check.
		http.Error(w, "No data available for this range", http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	if err := charting.GenerateToWriter(data, w); err != nil {
		if !quiet {
			log.Printf("Error generating chart: %v", err)
		}
		// Can't really write http.Error here if we started writing the image...
		// But GenerateToWriter should fail fast.
	}
}
