package listing10_13

import (
	"fmt"
	"time"
)

const (

	// ovenTime はオーブンでの調理時間を秒単位で指定する定数です。
	ovenTime = 5

	// everyThingElseTime は、その他の処理にかかる時間を秒単位で指定する定数です。
	everyThingElseTime = 2
)

// PrepareTray は指定されたトレイ番号の空のトレイを準備し、その番号を文字列で返します。
// trayNumberは準備するトレイの番号を指定します。
// トレイの準備中に少しの待機時間があります。
func PrepareTray(trayNumber int) string {
	fmt.Println("Preparing empty tray", trayNumber)
	time.Sleep(everyThingElseTime * time.Second)
	return fmt.Sprintf("tray number %d", trayNumber)
}

// Mixture は、指定されたトレイにカップケーキの混合物を注ぎ、処理完了後に結果を文字列で返します。
func Mixture(tray string) string {
	fmt.Println("Pouring cupcake Mixture in", tray)
	time.Sleep(everyThingElseTime * time.Second)
	return fmt.Sprintf("cupcake in %s", tray)
}

// Bake は指定された材料をオーブンで焼き上げ、焼き上がりの結果を文字列として返します。
func Bake(mixture string) string {
	fmt.Println("Baking", mixture)
	time.Sleep(ovenTime * time.Second)
	return fmt.Sprintf("baked %s", mixture)
}

// AddToppings は、指定された焼き上がったカップケーキにトッピングを追加し、その結果を文字列で返します。
func AddToppings(bakedCupCake string) string {
	fmt.Println("Adding topping to", bakedCupCake)
	time.Sleep(everyThingElseTime * time.Second)
	return fmt.Sprintf("topping on %s", bakedCupCake)
}

// Box は完成したカップケーキを箱詰めし、"～ boxed"形式の文字列を返します。
func Box(finishedCupCake string) string {
	fmt.Println("Boxing", finishedCupCake)
	time.Sleep(everyThingElseTime * time.Second)
	return fmt.Sprintf("%s boxed", finishedCupCake)
}
