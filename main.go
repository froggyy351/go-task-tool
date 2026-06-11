package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]

	// 引数なしは ls と同じ動作にする。
	if len(args) == 0 {
		runLs([]string{})
		return
	}

	switch args[0] {
	case "add":
		runAdd(args[1:])
	case "ls":
		runLs(args[1:])
	case "mv":
		runMv(args[1:])
	case "done":
		runDone(args[1:])
	case "undone":
		runUndone(args[1:])
	case "edit":
		runEdit(args[1:])
	case "rm":
		runRm(args[1:])
	case "help":
		printUsage()
	default:
		fmt.Printf("不明なコマンド: %s\n\n", args[0])
		printUsage()
	}
}

func printUsage() {
	fmt.Println(`tt - タスク管理ツール

使い方:
  tt add "タスク名" -d <日付>         MyBallにタスクを追加
  tt add "タスク名" -d <日付> --to 田中  OtherBallで登録
  tt ls                              一覧を期限が近い順に表示
  tt ls --my / --other               MyBall/OtherBallだけ表示
  tt ls --done                       完了済みタスクを表示
  tt mv <ID> <相手>                  ボールを相手へ（me で自分に引き取り）
  tt done <ID>                       完了にする
  tt undone <ID>                     完了を取り消して現役に戻す
  tt edit <ID> [-d <日付>] [-n "名前"]  期限・名前を変更
  tt rm <ID>                         削除
  tt help                            この使い方を表示

日付指定:
  2026-06-20   絶対指定
  +3           今日から3日後`)
}
