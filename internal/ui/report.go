package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/jasonsoprovich/shellclock/internal/model"
	"github.com/jasonsoprovich/shellclock/internal/util"
)

// ReportModel renders the summary table.
type ReportModel struct {
	store  *model.Store
	keys   KeyMap
	width  int
	height int

	SwitchToTree bool
}

func NewReportModel(store *model.Store, keys KeyMap) ReportModel {
	return ReportModel{store: store, keys: keys}
}

func (m ReportModel) Update(msg tea.Msg) (ReportModel, tea.Cmd) {
	if ws, ok := msg.(tea.WindowSizeMsg); ok {
		m.width = ws.Width
		m.height = ws.Height
		return m, nil
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q", "R":
			m.SwitchToTree = true
		}
	}
	return m, nil
}

func (m ReportModel) View() string {
	var sb strings.Builder

	sb.WriteString(StyleTitle.Render("Report") + "\n\n")

	colProject := 24
	colTask := 24
	colTime := 12

	header := fmt.Sprintf("%-*s  %-*s  %*s", colProject, "PROJECT", colTask, "TASK", colTime, "DURATION")
	sb.WriteString(StyleDimmed.Render(header) + "\n")
	sb.WriteString(StyleDimmed.Render(strings.Repeat("─", colProject+colTask+colTime+4)) + "\n")

	for _, p := range m.store.Projects {
		for _, t := range p.Tasks {
			dur := util.FormatDuration(t.TotalSeconds())
			line := fmt.Sprintf("%-*s  %-*s  %*s",
				colProject, truncate(p.Name, colProject),
				colTask, truncate(t.Name, colTask),
				colTime, dur,
			)
			sb.WriteString(line + "\n")
		}
		// Project total
		total := fmt.Sprintf("%-*s  %-*s  %*s",
			colProject, truncate(p.Name, colProject),
			colTask, "(total)",
			colTime, util.FormatDuration(p.TotalSeconds()),
		)
		sb.WriteString(StyleDuration.Render(total) + "\n")
		sb.WriteString("\n")
	}

	sb.WriteString("\n" + StyleDimmed.Render("esc/q back"))
	return StylePanel.Render(sb.String())
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-1] + "…"
}
