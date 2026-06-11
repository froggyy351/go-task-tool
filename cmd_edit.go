package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// runEdit は `tt edit <ID> [-d 日付] [-n "新しい名前"]` を処理する。
// -d と -n は両方指定することもできるが、少なくともどちらか1つが必要。
func runEdit(args []string) {
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "使い方: tt edit <ID> [-d 2026-06-25] [-n \"新しい名前\"]")
		os.Exit(1)
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Fprintln(os.Stderr, "IDは数値で指定してください:", args[0])
		os.Exit(1)
	}

	// フラグを順に解析する。
	flags := args[1:]
	var newDue *time.Time
	newName := ""

	for i := 0; i < len(flags); i++ {
		switch flags[i] {
		case "-d":
			if i+1 >= len(flags) {
				fmt.Fprintln(os.Stderr, "-d の後に日付を指定してください")
				os.Exit(1)
			}
			due, err := parseDueValue(flags[i+1])
			if err != nil {
				fmt.Fprintln(os.Stderr, "エラー:", err)
				os.Exit(1)
			}
			newDue = &due
			i++
		case "-n":
			if i+1 >= len(flags) {
				fmt.Fprintln(os.Stderr, "-n の後に名前を指定してください")
				os.Exit(1)
			}
			// タブはTSVを壊すのでスペースに置換する。
			newName = strings.ReplaceAll(flags[i+1], "\t", " ")
			i++
		}
	}

	if newDue == nil && newName == "" {
		fmt.Fprintln(os.Stderr, "-d または -n を少なくとも1つ指定してください")
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
		if newDue != nil {
			oldDue := tasks[i].Due.Format("2006-01-02")
			tasks[i].Due = *newDue
			fmt.Printf("[%d] %s の期限を %s → %s に変更しました\n",
				t.ID, t.Name, oldDue, newDue.Format("2006-01-02"))
		}
		if newName != "" {
			oldName := tasks[i].Name
			tasks[i].Name = newName
			fmt.Printf("[%d] タスク名を「%s」→「%s」に変更しました\n",
				t.ID, oldName, newName)
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
