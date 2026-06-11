package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func runAdd(args []string) {
	if len(args) < 3 {
		fmt.Fprintln(os.Stderr, "使い方: tt add \"タスク名\" -d 2026-06-20 [--to 田中]")
		os.Exit(1)
	}

	// タブはTSVを壊すのでスペースに置換する。
	name := strings.ReplaceAll(args[0], "\t", " ")

	due, err := parseDueFromArgs(args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, "エラー:", err)
		os.Exit(1)
	}

	// --to フラグがあれば OtherBall で登録する。
	owner := ""
	for i, a := range args[1:] {
		if a == "--to" && i+1 < len(args[1:]) {
			owner = args[1:][i+1]
		}
	}

	activePath, err := dataFilePath("tt.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "エラー:", err)
		os.Exit(1)
	}
	donePath, err := dataFilePath("tt_done.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "エラー:", err)
		os.Exit(1)
	}

	active, err := loadTasks(activePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ファイル読み込みエラー:", err)
		os.Exit(1)
	}
	done, err := loadDoneTasks(donePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ファイル読み込みエラー:", err)
		os.Exit(1)
	}

	ball := "my"
	if owner != "" {
		ball = "other"
	}

	newTask := Task{
		ID:    nextID(active, done),
		Name:  name,
		Due:   due,
		Ball:  ball,
		Owner: owner,
	}

	active = append(active, newTask)

	if err := saveTasks(activePath, active); err != nil {
		fmt.Fprintln(os.Stderr, "ファイル書き込みエラー:", err)
		os.Exit(1)
	}

	if owner != "" {
		fmt.Printf("[%d] %s を追加しました（期限: %s, OtherBall: %s）\n",
			newTask.ID, newTask.Name, newTask.Due.Format("2006-01-02"), owner)
	} else {
		fmt.Printf("[%d] %s を追加しました（期限: %s）\n",
			newTask.ID, newTask.Name, newTask.Due.Format("2006-01-02"))
	}
}

// parseDueFromArgs は引数リストから -d フラグの値を取り出して time.Time に変換する。
func parseDueFromArgs(args []string) (time.Time, error) {
	for i, a := range args {
		if a == "-d" && i+1 < len(args) {
			return parseDueValue(args[i+1])
		}
	}
	return time.Time{}, fmt.Errorf("-d で期限を指定してください（例: -d 2026-06-20 または -d +3）")
}

// parseDueValue は日付文字列を time.Time に変換する。
// "+3" のような相対指定（N日後）と "2026-06-20" の絶対指定に対応する。
func parseDueValue(s string) (time.Time, error) {
	if strings.HasPrefix(s, "+") {
		n, err := strconv.Atoi(s[1:])
		if err != nil {
			return time.Time{}, fmt.Errorf("相対日付の形式が不正: %s（例: +3）", s)
		}
		y, m, d := time.Now().AddDate(0, 0, n).Date()
		return time.Date(y, m, d, 0, 0, 0, 0, time.UTC), nil
	}
	return time.Parse("2006-01-02", s)
}
