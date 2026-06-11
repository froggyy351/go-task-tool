package main

import (
	"fmt"
	"os"
	"strconv"
)

// runUndone は `tt undone <ID>` を処理する。
// 完了済みタスクを現役に戻す。IDが現役と衝突する場合は新IDを採番する。
func runUndone(args []string) {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "使い方: tt undone <ID>")
		os.Exit(1)
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Fprintln(os.Stderr, "IDは数値で指定してください:", args[0])
		os.Exit(1)
	}

	donePath, err := dataFilePath("tt_done.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "エラー:", err)
		os.Exit(1)
	}

	doneTasks, err := loadDoneTasks(donePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ファイル読み込みエラー:", err)
		os.Exit(1)
	}

	var target *DoneTask
	var remainingDone []DoneTask
	for i, t := range doneTasks {
		if t.ID == id {
			target = &doneTasks[i]
		} else {
			remainingDone = append(remainingDone, t)
		}
	}

	if target == nil {
		fmt.Fprintf(os.Stderr, "完了済みタスクに ID %d が見つかりません\n", id)
		os.Exit(1)
	}

	activePath, err := dataFilePath("tt.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "エラー:", err)
		os.Exit(1)
	}

	active, err := loadTasks(activePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ファイル読み込みエラー:", err)
		os.Exit(1)
	}

	// 現役タスクとIDが衝突する場合は新IDを採番する。
	restored := target.Task
	for _, t := range active {
		if t.ID == restored.ID {
			newID := nextID(active, doneTasks)
			fmt.Printf("IDが重複したため %d → %d に変更しました\n", restored.ID, newID)
			restored.ID = newID
			break
		}
	}

	active = append(active, restored)

	if err := saveTasks(activePath, active); err != nil {
		fmt.Fprintln(os.Stderr, "ファイル書き込みエラー:", err)
		os.Exit(1)
	}

	if err := saveDoneTasks(donePath, remainingDone); err != nil {
		fmt.Fprintln(os.Stderr, "完了ファイル書き込みエラー:", err)
		os.Exit(1)
	}

	fmt.Printf("[%d] %s を現役に戻しました\n", restored.ID, restored.Name)
}
