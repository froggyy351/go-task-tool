package main

import (
	"fmt"
	"os"
	"strconv"
)

// runEdit は `tt edit <ID> -d 2026-06-25` を処理する。
func runEdit(args []string) {
	if len(args) < 3 {
		fmt.Fprintln(os.Stderr, "使い方: tt edit <ID> -d 2026-06-25")
		os.Exit(1)
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Fprintln(os.Stderr, "IDは数値で指定してください:", args[0])
		os.Exit(1)
	}

	newDue, err := parseDueFromArgs(args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, "エラー:", err)
		os.Exit(1)
	}

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

	found := false
	for i, t := range tasks {
		if t.ID != id {
			continue
		}
		found = true
		oldDue := tasks[i].Due.Format("2006-01-02")
		tasks[i].Due = newDue
		fmt.Printf("[%d] %s の期限を %s → %s に変更しました\n",
			t.ID, t.Name, oldDue, newDue.Format("2006-01-02"))
		break
	}

	if !found {
		fmt.Fprintf(os.Stderr, "ID %d のタスクが見つかりません\n", id)
		os.Exit(1)
	}

	if err := saveTasks(path, tasks); err != nil {
		fmt.Fprintln(os.Stderr, "ファイル書き込みエラー:", err)
		os.Exit(1)
	}
}
