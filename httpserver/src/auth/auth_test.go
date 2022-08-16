package auth

import (
	"fmt"
	"os"
	"testing"

	"encoding/csv"
)

func TestGenToken(t *testing.T) {
	token, err := GenToken("Ding123")
	if err != nil {
		t.Error(err)
	}
	fmt.Println("Token: ", token)
}

func TestGenTokenList(t *testing.T) {
	csvFile, err := os.Create("/Users/yuan.ding/Desktop/code/entry_task/jmeter_test/Token_2000.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer csvFile.Close()
	writer := csv.NewWriter(csvFile)
	title := []string{"token"}
	writer.Write(title)
	for num:=0; num < 2000; num++ {
		name := "stress_test_" + fmt.Sprint(10000000 + num)
		token, err := GenToken(name)
		if err != nil {
			t.Error(err)
		}
		record := []string{token}
		err = writer.Write(record)
		if err != nil {
			t.Fatal(err)
		}
		writer.Flush()
	}
}