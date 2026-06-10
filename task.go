package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// Task は1件のタスクを表す構造体（Goでいうデータの型定義）。
type Task struct {
	ID    int
	Name  string
	Due   time.Time
	Ball  string // "my" または "other"
	Owner string // Ball=="other" のときだけ使う
}

// ファイルパスを実行ファイルと同じフォルダに解決する。
func dataFilePath(name string) (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Join(filepath.Dir(exe), name), nil
}

// TSV の1行を Task に変換する。
// 形式: ID\tName\tDue\tBall\tOwner
func parseLine(line string) (Task, error) {
	parts := strings.Split(line, "\t")
	if len(parts) < 5 {
		return Task{}, fmt.Errorf("行の形式が不正: %s", line)
	}

	id, err := strconv.Atoi(parts[0])
	if err != nil {
		return Task{}, fmt.Errorf("IDが数値でない: %s", parts[0])
	}

	due, err := time.Parse("2006-01-02", parts[2])
	if err != nil {
		return Task{}, fmt.Errorf("日付の形式が不正: %s", parts[2])
	}

	return Task{
		ID:    id,
		Name:  parts[1],
		Due:   due,
		Ball:  parts[3],
		Owner: parts[4],
	}, nil
}

// Task を TSV の1行に変換する。
func (t Task) toLine() string {
	return fmt.Sprintf("%d\t%s\t%s\t%s\t%s",
		t.ID,
		t.Name,
		t.Due.Format("2006-01-02"),
		t.Ball,
		t.Owner,
	)
}

// ファイルから全タスクを読み込む。ファイルが存在しない場合は空スライスを返す。
func loadTasks(path string) ([]Task, error) {
	f, err := os.Open(path)
	if os.IsNotExist(err) {
		return []Task{}, nil
	}
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var tasks []Task
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		t, err := parseLine(line)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, scanner.Err()
}

// タスク一覧をファイルに書き出す（上書き）。
func saveTasks(path string, tasks []Task) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	for _, t := range tasks {
		fmt.Fprintln(w, t.toLine())
	}
	return w.Flush()
}

// 現在の最大IDを求めて +1 を返す。タスクが0件なら 1 を返す。
func nextID(tasks []Task) int {
	max := 0
	for _, t := range tasks {
		if t.ID > max {
			max = t.ID
		}
	}
	return max + 1
}
