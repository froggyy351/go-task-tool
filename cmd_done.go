package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

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

	activePath, err := dataFilePath("tt.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "エラー:", err)
		os.Exit(1)
	}

	tasks, err := loadTasks(activePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ファイル読み込みエラー:", err)
		os.Exit(1)
	}

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

	donePath, err := dataFilePath("tt_done.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "エラー:", err)
		os.Exit(1)
	}

	doneTasks, err := loadDoneTasks(donePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "完了ファイル読み込みエラー:", err)
		os.Exit(1)
	}

	doneTasks = append(doneTasks, DoneTask{Task: *target, DoneDate: time.Now()})

	if err := saveDoneTasks(donePath, doneTasks); err != nil {
		fmt.Fprintln(os.Stderr, "完了ファイル書き込みエラー:", err)
		os.Exit(1)
	}

	if err := saveTasks(activePath, remaining); err != nil {
		fmt.Fprintln(os.Stderr, "ファイル書き込みエラー:", err)
		os.Exit(1)
	}

	fmt.Printf("[%d] %s を完了しました\n", target.ID, target.Name)
}
