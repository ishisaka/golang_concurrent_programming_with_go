// filehash.go はファイルのハッシュコードを計算するための関数を提供する。

package listing10_1

import (
	"crypto/sha256"
	"io"
	"os"
)

// FHash は指定されたファイルパスのSHA-256ハッシュを計算し、バイト配列として返します。
func FHash(filepath string) []byte {
	file, _ := os.Open(filepath) // ファイルをオープン
	defer file.Close()

	sha := sha256.New() // crypto/sha256ライブラリを使ってハッシュコードを計算
	_, err := io.Copy(sha, file)
	if err != nil {
		panic(err)
	}
	return sha.Sum(nil) // ハッシュの結果を返す。
}
