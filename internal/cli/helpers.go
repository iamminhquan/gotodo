package cli

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	humanize "github.com/dustin/go-humanize"
	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"

	"github.com/iamminhquan/gotodo/internal/task"
)

// timePtr is an alias so add.go can declare a local *time.Time conveniently.
type timePtr = time.Time

// ── Terminal width detection ─────────────────────────────────────────────────

const (
	defaultTermWidth = 80  // safe fallback for narrow or unknown terminals
	maxTableWidth    = 140 // never let the table go wider than this
)

// terminalWidth tries to detect the terminal column count using only the
// standard library.  It checks the COLUMNS environment variable (set by most
// shells on resize) and falls back to defaultTermWidth.
func terminalWidth() int {
	if cols := os.Getenv("COLUMNS"); cols != "" {
		if n, err := strconv.Atoi(cols); err == nil && n > 0 {
			return n
		}
	}
	return defaultTermWidth
}

// clamp returns v clamped to [lo, hi].
func clamp(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

// shortID returns the first 8 characters of a UUID – long enough to be unique
// in a personal task list while staying readable.
func shortID(id string) string {
	if len(id) >= 8 {
		return id[:8]
	}
	return id
}

// humanDate formats a time.Time as a relative string ("in 3 days",
// "2 hours ago") using go-humanize, with a fallback to an absolute date.
func humanDate(t time.Time) string {
	return humanize.Time(t)
}

// formatDue returns a coloured relative due-date string, or an empty string
// if the task has no due date.
func formatDue(t *task.Task) string {
	if t.Due == nil {
		return ""
	}
	rel := humanDate(*t.Due)
	if t.IsOverdue() {
		return color.RedString("⚠ " + rel)
	}
	if t.IsDueToday() {
		return color.YellowString("🔔 " + rel)
	}
	return color.CyanString("📅 " + rel)
}

// formatTags returns a comma-joined tag string, or empty.
func formatTags(tags []string) string {
	if len(tags) == 0 {
		return ""
	}
	return "🏷 " + strings.Join(tags, ", ")
}

// priorityColor returns a coloured priority string.
func priorityColor(p string) string {
	switch p {
	case task.PriorityHigh:
		return color.RedString(p)
	case task.PriorityMedium:
		return color.YellowString(p)
	case task.PriorityLow:
		return color.CyanString(p)
	default:
		return p
	}
}

func priorityLabel(t *task.Task, useColor bool) string {
	word := t.Priority
	if useColor {
		word = priorityColor(t.Priority)
	}
	return t.PriorityEmoji() + " " + word
}

func buildTable(tasks []*task.Task, useColor bool) table.Writer {
	tw := terminalWidth()
	width := clamp(tw, 40, maxTableWidth)

	t := table.NewWriter()
	t.SetStyle(table.StyleRounded)
	t.SetAllowedRowLength(width)

	const (
		colOverhead = 3 // per column: "│" (1) + left pad (1) + right pad (1)
		borderExtra = 1 // trailing "│" on the rightmost column
		statusW     = 4 + colOverhead
		idW         = 8 + colOverhead
		priorityW   = 15 + colOverhead // "🟡 medium" needs ~15 display chars
		dueW        = 22 + colOverhead
		tagsW       = 22 + colOverhead
		createdW    = 16 + colOverhead
	)

	// Determine which optional columns fit at this width.
	showPriority := width >= 50
	showDue := width >= 65
	showTags := width >= 95
	showCreated := width >= 115

	// Compute remaining width for the Title column.
	fixedUsed := borderExtra + statusW + idW
	if showPriority {
		fixedUsed += priorityW
	}
	if showDue {
		fixedUsed += dueW
	}
	if showTags {
		fixedUsed += tagsW
	}
	if showCreated {
		fixedUsed += createdW
	}

	titleWidth := width - fixedUsed - colOverhead // subtract Title's own overhead
	maxTitleWidth := (width * 3) / 4              // never wider than 3/4 of the terminal
	if titleWidth < 12 {
		titleWidth = 12
	}
	if titleWidth > maxTitleWidth {
		titleWidth = maxTitleWidth
	}

	cols := []table.ColumnConfig{
		{Number: 1, Align: text.AlignCenter, AlignHeader: text.AlignCenter, WidthMin: 4, WidthMax: 4},
		{Number: 2, Align: text.AlignLeft, AlignHeader: text.AlignCenter, WidthMin: 8, WidthMax: 10},
		{Number: 3, Align: text.AlignLeft, AlignHeader: text.AlignCenter, WidthMax: titleWidth, WidthMaxEnforcer: text.WrapSoft},
		{Number: 4, Align: text.AlignLeft, AlignHeader: text.AlignLeft, WidthMin: 13, WidthMax: 15, Hidden: !showPriority},
		{Number: 5, Align: text.AlignLeft, AlignHeader: text.AlignCenter, WidthMax: 22, WidthMaxEnforcer: text.WrapSoft, Hidden: !showDue},
		{Number: 6, Align: text.AlignLeft, AlignHeader: text.AlignCenter, WidthMax: 22, WidthMaxEnforcer: text.WrapSoft, Hidden: !showTags},
		{Number: 7, Align: text.AlignLeft, AlignHeader: text.AlignCenter, WidthMax: 16, Hidden: !showCreated},
	}

	t.SetColumnConfigs(cols)

	t.AppendHeader(table.Row{"", "ID", "Title", "Priority", "Due", "Tags", "Created"})

	for i, tk := range tasks {
		title := tk.Title
		if tk.Done && useColor {
			title = color.New(color.Faint).Sprint(title)
		}

		due := formatDue(tk)
		if !useColor && tk.Due != nil {
			due = fmt.Sprintf("📅 %s", tk.Due.Format("2006-01-02"))
		}

		t.AppendRow(table.Row{
			tk.StatusEmoji(),
			shortID(tk.ID),
			title,
			priorityLabel(tk, useColor),
			due,
			formatTags(tk.Tags),
			humanDate(tk.CreatedAt),
		})

		// Separator after every row (including the last) so each task is
		// visually framed: the table closing border acts as the bottom of
		// the final separator, making all rows look uniform.
		_ = i // i is kept for any future per-row logic
		t.AppendSeparator()
	}

	return t
}

// printNoTasks prints a friendly message when nothing matches the filter.
func printNoTasks(msg string) {
	if msg == "" {
		msg = "No tasks found."
	}
	fmt.Println()
	color.New(color.FgCyan).Printf("  %s\n", msg)
	fmt.Printf("  Run %s to create one!\n\n", color.GreenString("gotodo add \"your task\""))
}

// confirmPrompt asks the user for y/N confirmation and returns true on "y".
func confirmPrompt(msg string) bool {
	fmt.Printf("%s [y/N] ", msg)
	var ans string
	_, _ = fmt.Scanln(&ans)
	return strings.ToLower(strings.TrimSpace(ans)) == "y"
}
