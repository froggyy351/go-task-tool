package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// runRm は `tt rm <ID>` を処理する。
// 削除前に確認を求める。
func runRm(args []string) {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "使い方: tt rm <ID>")
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

	// 誤削除防止のため確認を求める。
	fmt.Printf("[%d] %s（期限: %s）を削除しますか？ [y/N]: ",
		target.ID, target.Name, target.Due.Format("2006-01-02"))

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	answer := strings.TrimSpace(scanner.Text())

	if answer != "y" && answer != "Y" {
		fmt.Println("キャンセルしました。")
		return
	}

	if err := saveTasks(path, remaining); err != nil {
		fmt.Fprintln(os.Stderr, "ファイル書き込みエラー:", err)
		os.Exit(1)
	}

	fmt.Printf("[%d] %s を削除しました\n", target.ID, target.Name)
}
