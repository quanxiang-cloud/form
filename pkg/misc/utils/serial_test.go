package utils

import (
	"fmt"
	"math/big"
	"testing"
)

func TestParseTemplate(t *testing.T) {
	//ParseTemplate("ER.date{yyyyMMdd}.incr[name]{5,1}.step[name]{1}.XXXX")

	i, b := big.NewInt(0).SetString("10", 10)
	fmt.Println(i, b)

	i, b = big.NewInt(0).SetString("", 20)
	fmt.Println(i, b)

	//fmt.Printf("%s", fmt.Sprintf("%%0%[1]s%[2]s", "5", "3"))

	s := fmt.Sprintf("%05s", "123")
	fmt.Println(s)

	format := "%0" + "5" + "s"
	fmt.Printf("fmt.Sprintf(format, \"123\"): %v\n", fmt.Sprintf(format, "123456999"))

	bit := big.NewInt(12220)
	fmt.Printf("bit.Add(bit, big.NewInt(1)).String(): %v\n", bit.Add(bit, big.NewInt(1)).String())

}
