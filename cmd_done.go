package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// runDone は `tt done <ID>` を処理する。
// タスクを tt.txt から削除し、tt_done.txt に完了日付きで追記する。
func runDone(args []string) {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "使い方: tt done <ID>")
		os.Exit(1)
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Fprintln(os.Stderr, "IDは数値で指定してください:", args[0])
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

	// IDが一致するタスクを探し、残りと分離する。
	var target *Task
	var remaining []Task
	for i, t := range tasks {
		if t.ID == id {
			target = &tasks[i]
		} else {
			remaining = append(remaining, t)
		}
	}

	if target == nil {
		fmt.Fprintf(os.Stderr, "ID %d のタスクが見つかりません\n", id)
		os.Exit(1)
	}

	// tt_done.txt に「元のTSV行 + タブ + 完了日」を追記する。
	donePath, err := dataFilePath("tt_done.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "エラー:", err)
		os.Exit(1)
	}

	doneFile, err := os.OpenFile(donePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintln(os.Stderr, "完了ファイル書き込みエラー:", err)
		os.Exit(1)
	}
	defer doneFile.Close()

	doneDate := time.Now().Format("2006-01-02")
	fmt.Fprintf(doneFile, "%s\t%s\n", target.toLine(), doneDate)

	// tt.txt から該当タスクを除いて上書き保存する。
	if err := saveTasks(path, remaining); err != nil {
		fmt.Fprintln(os.Stderr, "ファイル書き込みエラー:", err)
		os.Exit(1)
	}

	fmt.Printf("[%d] %s を完了しました\n", target.ID, target.Name)
}
