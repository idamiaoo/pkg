package util

import "testing"

func TestFAlias_CreateFile(t *testing.T) {
	f, err := File.CreateFile("can.txt")
	if err != nil {
		t.Error(err)
	}

	t.Log(f)
}

func TestFAlias_WriteFile(t *testing.T) {
	if err := File.WriteFile("can.txt", []byte("hello word")); err != nil {
		t.Error(err)
	}
}

func TestFAlias_DeleteFile(t *testing.T) {
	if err := File.DeleteFile("can.txt"); err != nil {
		t.Error(err)
	}
}
