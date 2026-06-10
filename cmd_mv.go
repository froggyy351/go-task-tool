package main

import (
	"fmt"
	"os"
	"strconv"
)

// runMv は `tt mv <ID> <相手>` を処理する。
// 相手が "me" なら MyBall に引き取り、それ以外は OtherBall に移して担当者を設定する。
func runMv(args []string) {
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "使い方: tt mv <ID> <相手>  （自分に引き取る場合は me）")
		os.Exit(1)
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Fprintln(os.Stderr, "IDは数値で指定してください:", args[0])
		os.Exit(1)
	}
	to := args[1]

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

	// IDが一致するタスクを探してボールを更新する。
	found := false
	for i, t := range tasks {
		if t.ID != id {
			continue
		}
		found = true
		if to == "me" {
			tasks[i].Ball = "my"
			tasks[i].Owner = ""
			fmt.Printf("[%d] %s を引き取りました（MyBall）\n", t.ID, t.Name)
		} else {
			tasks[i].Ball = "other"
			tasks[i].Owner = to
			fmt.Printf("[%d] %s のボールを %s へ渡しました（OtherBall）\n", t.ID, t.Name, to)
		}
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
