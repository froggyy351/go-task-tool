package main

import (
	"fmt"
	"os"
	"time"
)

// runAdd は `tt add "タスク名" -d 2026-06-20` を処理する。
// args には "add" より後ろの部分が渡される（例: ["資料作成", "-d", "2026-06-20"]）。
func runAdd(args []string) {
	// 引数が足りない場合はエラーを出して終了。
	if len(args) < 3 {
		fmt.Fprintln(os.Stderr, "使い方: tt add \"タスク名\" -d 2026-06-20")
		os.Exit(1)
	}

	name := args[0]
	due, err := parseDueFromArgs(args[1:])
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

	newTask := Task{
		ID:   nextID(tasks),
		Name: name,
		Due:  due,
		Ball: "my",
	}

	tasks = append(tasks, newTask)

	if err := saveTasks(path, tasks); err != nil {
		fmt.Fprintln(os.Stderr, "ファイル書き込みエラー:", err)
		os.Exit(1)
	}

	fmt.Printf("[%d] %s を追加しました（期限: %s）\n",
		newTask.ID, newTask.Name, newTask.Due.Format("2006-01-02"))
}

// "-d 2026-06-20" 形式の引数から time.Time を取り出す。
func parseDueFromArgs(args []string) (time.Time, error) {
	for i, a := range args {
		if a == "-d" && i+1 < len(args) {
			return time.Parse("2006-01-02", args[i+1])
		}
	}
	return time.Time{}, fmt.Errorf("-d で期限を指定してください（例: -d 2026-06-20）")
}
