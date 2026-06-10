package main

import (
	"fmt"
	"os"
)


func main() {
	// os.Args[0] はプログラム名(tt)自身。実際の引数は [1:] から。
	args := os.Args[1:]

	// サブコマンドが無ければ使い方を表示して終了。
	if len(args) == 0 {
		printUsage()
		return
	}

	// 最初の引数をサブコマンドとして振り分ける。
	switch args[0] {
	case "add":
		runAdd(args[1:])
	case "ls":
		runLs(args[1:])
	case "mv":
		runMv(args[1:])
	case "done":
		runDone(args[1:])
	case "edit":
		runEdit(args[1:])
	case "rm":
		runRm(args[1:])
	default:
		fmt.Printf("不明なコマンド: %s\n\n", args[0])
		printUsage()
	}
}

func printUsage() {
	fmt.Println(`tt - タスク管理ツール

使い方:
  tt add "タスク名" -d 2026-06-20   タスクを追加（MyBall）
  tt ls                            一覧を期限が近い順に表示
  tt mv <ID> <相手>                ボールを相手へ（me で自分に引き取り）
  tt done <ID>                     完了にする
  tt edit <ID> -d 2026-06-25       期限を変更
  tt rm <ID>                       削除`)
}
