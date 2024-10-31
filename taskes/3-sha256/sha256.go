package main

import (
	"crypto/sha256"
	"fmt"
)

func main() {

	var testMsg string
	print("enter massage ): \n")
	fmt.Scan(&testMsg)
	hash := sha256.Sum256([]byte(testMsg))
	cnt := 0
	for {
		if hash[0] == 0 && hash[1] == 0 {
			break
		}
		hash = sha256.Sum256([]byte(testMsg + string(cnt)))
		cnt = cnt + 1
	}
	fmt.Printf("\n new hash  is : %x \n", hash[:])
	fmt.Printf("counter make  4 zeros : %d \n", cnt)

	return
}
