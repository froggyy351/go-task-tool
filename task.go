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

type Task struct {
	ID    int
	Name  string
	Due   time.Time
	Ball  string // "my" または "other"
	Owner string // Ball=="other" のときだけ使う
}

type DoneTask struct {
	Task
	DoneDate time.Time
}

func (t DoneTask) toLine() string {
	return fmt.Sprintf("%s\t%s", t.Task.toLine(), t.DoneDate.Format("2006-01-02"))
}

func dataFilePath(name string) (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Join(filepath.Dir(exe), name), nil
}

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

func (t Task) toLine() string {
	return fmt.Sprintf("%d\t%s\t%s\t%s\t%s",
		t.ID,
		t.Name,
		t.Due.Format("2006-01-02"),
		t.Ball,
		t.Owner,
	)
}

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

// saveTasks は一時ファイルに書き出してからリネームする（書き込み中クラッシュ対策）。
func saveTasks(path string, tasks []Task) error {
	tmp := path + ".tmp"
	f, err := os.Create(tmp)
	if err != nil {
		return err
	}
	w := bufio.NewWriter(f)
	for _, t := range tasks {
		fmt.Fprintln(w, t.toLine())
	}
	if err := w.Flush(); err != nil {
		f.Close()
		os.Remove(tmp)
		return err
	}
	if err := f.Close(); err != nil {
		os.Remove(tmp)
		return err
	}
	return os.Rename(tmp, path)
}

func parseDoneLine(line string) (DoneTask, error) {
	parts := strings.Split(line, "\t")
	if len(parts) < 6 {
		return DoneTask{}, fmt.Errorf("行の形式が不正: %s", line)
	}
	t, err := parseLine(strings.Join(parts[:5], "\t"))
	if err != nil {
		return DoneTask{}, err
	}
	doneDate, err := time.Parse("2006-01-02", parts[5])
	if err != nil {
		return DoneTask{}, fmt.Errorf("完了日の形式が不正: %s", parts[5])
	}
	return DoneTask{Task: t, DoneDate: doneDate}, nil
}

func loadDoneTasks(path string) ([]DoneTask, error) {
	f, err := os.Open(path)
	if os.IsNotExist(err) {
		return []DoneTask{}, nil
	}
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var tasks []DoneTask
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		t, err := parseDoneLine(line)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, scanner.Err()
}

// saveDoneTasks は完了ファイルを一時ファイル経由で上書き保存する。
func saveDoneTasks(path string, tasks []DoneTask) error {
	tmp := path + ".tmp"
	f, err := os.Create(tmp)
	if err != nil {
		return err
	}
	w := bufio.NewWriter(f)
	for _, t := range tasks {
		fmt.Fprintln(w, t.toLine())
	}
	if err := w.Flush(); err != nil {
		f.Close()
		os.Remove(tmp)
		return err
	}
	if err := f.Close(); err != nil {
		os.Remove(tmp)
		return err
	}
	return os.Rename(tmp, path)
}

// nextID は現役タスクと完了タスク両方の最大IDを考慮して次のIDを返す。
func nextID(active []Task, done []DoneTask) int {
	max := 0
	for _, t := range active {
		if t.ID > max {
			max = t.ID
		}
	}
	for _, t := range done {
		if t.ID > max {
			max = t.ID
		}
	}
	return max + 1
}
