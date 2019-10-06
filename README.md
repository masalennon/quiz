# quiz
## 概要
csvに書いてある問題を読み込んで、コマンドラインからの入力を回答として受け取り、最後に正解数を表示する時間制限付きクイズです。

![image](https://user-images.githubusercontent.com/17466118/66254304-83521b00-e7af-11e9-9bba-d77f535c6ace.png)

## 使い方
- limitというオプションで制限時間（秒）を指定してプログラムを実行すると、クイズが開始されます。
- 解答を入力していき、指定した制限時間が来るか全ての問題を解ききると終了になります。
- 不正解だった場合には正しい回答が表示され、次の問題へと進みます。
- プログラムの終了時に何問中何問正解したかが表示されます。

## コードのポイント
### CSVを読み込んで、CSVファイルの値を独自に定義した構造体型のスライスにセットする
#### CSVの読み込み
```go
csvFilename := flag.String("csv", "problems.csv", "'問題, 答え'のフォーマットのCSVファイルを指定します。")
timeLimit := flag.Int("limit", 7, "クイズの制限時間を秒で指定します。")
flag.Parse()
```
上記のコードで、プログラムにコマンドライン引数を渡すことができます。ちなみに、flag.Parse()がないと引数で値を指定してもデフォルト値が渡されます。

#### CSVファイルの値を独自に定義した構造体型のスライスにセットする
```go
func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines))
	for i, line := range lines {
		ret[i] = problem{
			q: line[0],
			a: strings.TrimSpace(line[1]), //スペースがcsvに入っているときのため
		}
	}
	return ret
}
```
読み込むcsvファイルの名前を引数で指定したら、csvパッケージのReaderを使ってcsvファイルを読み込み、stringの二次元配列を独自で定義したproblemという構造体型のスライスに格納します。こうすることで、forループで問題を順番に出す処理を単純化できます。

### チャネルとゴルーチンを使って、制限時間が過ぎるとユーザの入力の受付を停止しプログラムを終了させる
### ゴルーチンとチャネル

```go
timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

problemloop:
	for i, p := range problems {
		fmt.Printf("第%d問: %s = ", i+1, p.q)
		answerCh := make(chan string)
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer) //時間が過ぎても止め続ける
			answerCh <- answer
		}()
		select {
		case <-timer.C:
			break problemloop
		case answer := <-answerCh:
			if answer == p.a {
				fmt.Println("正解です。\n")
				correct++
			} else {
				fmt.Printf("不正解です。正解は%sです。\n\n", p.a)
			}
		}
	}
```

上記のコードには一つのゴルーチンがあります。ゴルーチンを使って並列処理をしている理由は、			
`fmt.Scanf("%s\n", &answer)` この行でユーザからの入力が来るまでプログラムが止まってしまい、それ以下の行へと処理が進んでいかないため設定した制限時間を過ぎてもプログラムが終了しないということを防ぐためです。

ゴルーチン内でユーザからの入力を受け付けながら、同時に制限時間が来たことを知らせるチャネル、ユーザからの入力を受け取るチャネルもselect文を使って待機させることで、ユーザからの入力があった時と制限時間が尽きた時でそれぞれ別の処理を実行させることができます。

チャネルと並列処理（ゴルーチン）の使い所と、他言語では複雑・長くなってしまう処理も、Go言語の並列処理とチャネルを使うと短く簡潔に書くことができることが分かるいい例です。
