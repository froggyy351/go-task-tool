package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

func runLs(args []string) {
	// --done フラグがあれば完了済みタスクを表示して終了。
	if contains(args, "--done") {
		runLsDone()
		return
	}

	// --my / --other フラグを確認する。
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

	// 期限が近い順にソートする。
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].Due.Before(tasks[j].Due)
	})

	today := time.Now().Truncate(24 * time.Hour)

	// フィルタに応じてセクションを表示する。
	if !filterOther {
		printSection("MyBall", "my", tasks, today)
	}
	if !filterMy {
		printSection("OtherBall", "other", tasks, today)
	}
}

func printSection(title, ball string, tasks []Task, today time.Time) {
	// 対象セクションのタスクだけ抽出する。
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

	for _, t := range filtered {
		line := formatTask(t, today)
		daysLeft := int(t.Due.Sub(today).Hours() / 24)
		if daysLeft == 0 {
			line = red(line)
		}
		fmt.Println(line)
	}
	fmt.Println()
}

func formatTask(t Task, today time.Time) string {
	daysLeft := int(t.Due.Sub(today).Hours() / 24)

	// 期限の状態を判定する。
	var dueMark string
	switch {
	case daysLeft < 0:
		dueMark = fmt.Sprintf("[!期限切れ %d日]", -daysLeft)
	case daysLeft == 0:
		dueMark = "[!今日]"
	case daysLeft <= 3:
		dueMark = fmt.Sprintf("[!あと%d日]", daysLeft)
	default:
		dueMark = fmt.Sprintf("  あと%d日 ", daysLeft)
	}

	// ID と名前を固定幅にして揃える。
	idPart := fmt.Sprintf("[%d]", t.ID)
	namePart := padRight(t.Name, 20)
	duePart := t.Due.Format("2006-01-02")

	line := fmt.Sprintf("  %-4s %s %s %s", idPart, namePart, duePart, dueMark)

	// OtherBall は担当者を追加表示する。
	if t.Ball == "other" && t.Owner != "" {
		line += fmt.Sprintf("  担当: %s", t.Owner)
	}

	return line
}

// 文字列を指定の幅に右埋めする（日本語考慮なし・簡易版）。
func padRight(s string, width int) string {
	runes := []rune(s)
	if len(runes) >= width {
		return string(runes)
	}
	return s + strings.Repeat(" ", width-len(runes))
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

	// 完了日が新しい順に表示する。
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].DoneDate.After(tasks[j].DoneDate)
	})

	for _, t := range tasks {
		namePart := padRight(t.Name, 20)
		fmt.Printf("  [%d] %s 完了: %s  (期限: %s)\n",
			t.ID, namePart,
			t.DoneDate.Format("2006-01-02"),
			t.Due.Format("2006-01-02"),
		)
	}
}

func contains(args []string, flag string) bool {
	for _, a := range args {
		if a == flag {
			return true
		}
	}
	return false
}
