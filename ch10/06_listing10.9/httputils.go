package listing10_9

import (
	"fmt"
	"net"
	"os"
	"regexp"
)

var r, _ = regexp.Compile("GET (.+) HTTP/1.1\r\n")

func handleHttpRequest(conn net.Conn) {
	buff := make([]byte, 1024) // HTTPリクエストを保存するバッファを作成
	size, _ := conn.Read(buff) // 接続からバッファ文を読み込む
	if r.Match(buff[:size]) {  // リクエストが有効なら、資源ディレクトリから要求されたファイルを読みこむ
		file, err := os.ReadFile(
			fmt.Sprintf("../resources/%s", r.FindSubmatch(buff[:size])[1]))
		// ファイルが存在すれば、HTTPヘッダとファイルの内容でクライアントに応答する。
		if err == nil {
			conn.Write([]byte(fmt.Sprintf(
				"HTTP/1.1 200 OK\r\nContent-Length: %d\r\n\r\n", len(file))))
			conn.Write(file)
		} else { // ファイルが存在しなければ、エラーで応答する
			conn.Write([]byte(
				"HTTP/1.1 404 Not Found\r\n\r\n<html>Not Found</html>"))
		}
	} else { // HTTPリクエストが有効でなければ、エラーで応答する
		conn.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n\r\n"))
	}
	conn.Close() // リクエスト処理後、接続をクローズ
}

func StartHttpWorkers(n int, incomingConnections <-chan net.Conn) {
	for i := 0; i < n; i++ { // ゴルーチン開始
		go func() {
			// チャネルがクローズされるまで作業キューチャネルから接続を受け取る
			for c := range incomingConnections {
				// 受信したコネクションからHTTPリクエストを処理
				handleHttpRequest(c)
			}
		}()
	}
}
