package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/mattn/go-runewidth"
)

const nameMaxWidth = 50 // タスク名の表示幅の上限（全角25文字相当）

func runLs(args []string) {
	if contains(args, "--done") {
		runLsDone()
		return
	}

	filterMy := contains(args, "--my")
	filterOther := contains(args, "--other")

	path, err := dataFilePath("tt.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "エラー:", err)
		os.Exit(1)
	}

	tasks, err := loadTasks(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ファイル読み込みエラー:", err)
		os.Exit(1)
	}

	if len(tasks) == 0 {
		fmt.Println("タスクはありません。")
		return
	}

	sort.SliceStable(tasks, func(i, j int) bool {
		return tasks[i].Due.Before(tasks[j].Due)
	})

	today := todayUTC()

	if !filterOther {
		printSection("MyBall", "my", tasks, today)
	}
	if !filterMy {
		printSection("OtherBall", "other", tasks, today)
	}
}

func printSection(title, ball string, tasks []Task, today time.Time) {
	var filtered []Task
	for _, t := range tasks {
		if t.Ball == ball {
			filtered = append(filtered, t)
		}
	}

	fmt.Println(green(fmt.Sprintf("=== %s (%d件) ===", title, len(filtered))))

	if len(filtered) == 0 {
		fmt.Println("  (なし)")
		fmt.Println()
		return
	}

	nameWidth := 0
	for _, t := range filtered {
		w := runewidth.StringWidth(t.Name)
		if w > nameWidth {
			nameWidth = w
		}
	}
	if nameWidth > nameMaxWidth {
		nameWidth = nameMaxWidth
	}

	for _, t := range filtered {
		line := formatTask(t, today, nameWidth)
		daysLeft := int(t.Due.Sub(today).Hours() / 24)
		if daysLeft <= 0 {
			line = red(line)
		}
		fmt.Println(line)
	}
	fmt.Println()
}

func formatTask(t Task, today time.Time, nameWidth int) string {
	daysLeft := int(t.Due.Sub(today).Hours() / 24)

	var dueMark string
	switch {
	case daysLeft < 0:
		dueMark = fmt.Sprintf("[!期限切れ %d日]", -daysLeft)
	case daysLeft == 0:
		dueMark = "[!今日]"
	case daysLeft <= 3:
		dueMark = fmt.Sprintf("[!あと%d日]", daysLeft)
	default:
		dueMark = fmt.Sprintf("あと%d日", daysLeft)
	}

	idPart := fmt.Sprintf("[%d]", t.ID)
	namePart := truncatePad(t.Name, nameWidth)
	duePart := t.Due.Format("2006-01-02")

	line := fmt.Sprintf("  %-5s  %s  %s  %s", idPart, namePart, duePart, dueMark)

	if t.Ball == "other" && t.Owner != "" {
		line += fmt.Sprintf("  担当: %s", t.Owner)
	}

	return line
}

func truncatePad(s string, width int) string {
	w := runewidth.StringWidth(s)
	if w <= width {
		return s + strings.Repeat(" ", width-w)
	}
	result := []rune{}
	current := 0
	for _, r := range s {
		rw := runewidth.RuneWidth(r)
		if current+rw > width-2 {
			break
		}
		result = append(result, r)
		current += rw
	}
	truncated := string(result) + "…"
	return truncated + strings.Repeat(" ", width-current-2)
}

func runLsDone() {
	path, err := dataFilePath("tt_done.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "エラー:", err)
		os.Exit(1)
	}

	tasks, err := loadDoneTasks(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ファイル読み込みエラー:", err)
		os.Exit(1)
	}

	fmt.Println(green(fmt.Sprintf("=== 完了済み (%d件) ===", len(tasks))))
	if len(tasks) == 0 {
		fmt.Println("  (なし)")
		return
	}

	sort.SliceStable(tasks, func(i, j int) bool {
		return tasks[i].DoneDate.After(tasks[j].DoneDate)
	})

	nameWidth := 0
	for _, t := range tasks {
		w := runewidth.StringWidth(t.Name)
		if w > nameWidth {
			nameWidth = w
		}
	}
	if nameWidth > nameMaxWidth {
		nameWidth = nameMaxWidth
	}

	for _, t := range tasks {
		namePart := truncatePad(t.Name, nameWidth)
		fmt.Printf("  [%d]  %s  完了: %s  (期限: %s)\n",
			t.ID, namePart,
			t.DoneDate.Format("2006-01-02"),
			t.Due.Format("2006-01-02"),
		)
	}
}

// todayUTC はローカルの今日の日付をUTC midnightとして返す。
// Due日付がUTC midnightでparseされるため、比較がずれないようにする。
func todayUTC() time.Time {
	y, m, d := time.Now().Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}

func contains(args []string, flag string) bool {
	for _, a := range args {
		if a == flag {
			return true
		}
	}
	return false
}
