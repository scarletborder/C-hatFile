package blogs

import (
	"os"
	"path/filepath"
	"sort"
)

// 重新初始化
func loadFiles(directory string) error {
	fileInfo, err := os.ReadDir(directory)
	if err != nil {
		return err
	}

	fileIndexMap = make(map[string]int)
	files = make([]FileInfo, 0)

	for _, file := range fileInfo {
		if !file.IsDir() {
			filePath := filepath.Join(directory, file.Name())
			fileStat, err := os.Stat(filePath)
			if err != nil {
				return err
			}
			files = append(files, FileInfo{
				Name:       file.Name(),
				CreateTime: fileStat.ModTime(),
			})
		}
	}

	// Sort files by creation time, newest first
	sort.Slice(files, func(i, j int) bool {
		return files[i].CreateTime.After(files[j].CreateTime)
	})

	leng := len(files)
	for i := 0; i < leng; i += 1 {
		files[i].ID = i
		fileIndexMap[files[i].Name] = i
	}

	return nil
}

// 读取预览信息
// return files

// 读取某个特定文件
func GetBlogContent(file_name string) ([]byte, error) {
	fileID := fileIndexMap[file_name]
	filePath := filepath.Join(directory, files[fileID].Name)
	return os.ReadFile(filePath)
}
