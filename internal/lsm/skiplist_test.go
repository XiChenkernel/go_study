package lsm

import (
	"SQL/internal/database"
	"SQL/internal/storage"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"testing"
	"time"
)

// TestDataInfoGenerationAndWrite 测试 DataInfo 生成和写入
func TestDataInfoGenerationAndWrite(t *testing.T) {
	// 并发写入的 goroutine 数量
	concurrency := 100

	// 创建跳表实例
	sl := NewSkipList(16)

	// 生成测试数据并插入跳表
	var testData []DataInfo
	for i := 0; i < concurrency; i++ {
		data := DataInfo{
			DataMeta: database.DataMeta{
				Key:       []byte(generateRandomKey()),
				Value:     []byte(fmt.Sprintf("value%d", i)),
				Extra:     []byte(fmt.Sprintf("extra%d", i)),
				KeySize:   uint32(len(fmt.Sprintf("key%d", i))),
				ValueSize: uint32(len(fmt.Sprintf("value%d", i))),
				ExtraSize: uint32(len(fmt.Sprintf("extra%d", i))),
				TTL:       time.Duration(rand.Intn(3600)) * time.Second, // 随机生成 TTL
			},
			StorageLocation: storage.StorageLocation{
				FileName: []byte("data.txt"),
				Offset:   int64(i * 100), // 假设每条数据占用 100 字节
				Size:     100,
			},
		}
		testData = append(testData, data)
		sl.Insert(data.Key, &data)
	}

	// 并发写入数据到文件
	var wg sync.WaitGroup
	wg.Add(concurrency)
	for _, data := range testData {
		go func(d DataInfo) {
			defer wg.Done()
			writeDataToFile(d)
		}(data)
	}
	wg.Wait()

	// 将跳表内容写入文件
	file, err := os.OpenFile("../../data/testdata/skiplist/skiplist_content.txt", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("failed to open file: %v\n", err)
		return
	}
	defer file.Close()

	node := sl.Head.Next[0]
	for node != nil {
		line := fmt.Sprintf("Key: %s, Value: %s, Extra: %s, KeySize: %d, ValueSize: %d, ExtraSize: %d, TTL: %s, FileName: %s, Offset: %d, Size: %d\n",
			node.Key, node.DataInfo.Value, node.DataInfo.Extra, node.DataInfo.KeySize, node.DataInfo.ValueSize, node.DataInfo.ExtraSize, node.DataInfo.TTL, node.DataInfo.FileName, node.DataInfo.Offset, node.DataInfo.Size)
		_, err := file.WriteString(line)
		if err != nil {
			fmt.Printf("failed to write to file: %v\n", err)
			return
		}
		node = node.Next[0]
	}
}

// writeDataToFile 将 DataInfo 写入文件
func writeDataToFile(data DataInfo) {
	// 打开文件，如果不存在则创建
	file, err := os.OpenFile(string("../../data/testdata/skiplist/test1.txt"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("failed to open file: %v\n", err)
		return
	}
	defer file.Close()

	// 将 DataInfo 格式化为字符串
	line := fmt.Sprintf("Key: %s, Value: %s, Extra: %s, KeySize: %d, ValueSize: %d, ExtraSize: %d, TTL: %s, FileName: %s, Offset: %d, Size: %d\n",
		data.Key, data.Value, data.Extra, data.KeySize, data.ValueSize, data.ExtraSize, data.TTL, data.FileName, data.Offset, data.Size)

	// 写入数据到文件
	_, err = file.WriteString(line)
	if err != nil {
		fmt.Printf("failed to write to file: %v\n", err)
	}
}

// go test -run=^$ -bench=. -benchmem
func BenchmarkDataInfoGenerationAndWrite(b *testing.B) {
	for i := 0; i < b.N; i++ {
		TestDataInfoGenerationAndWrite(nil)
	}
}

// generateRandomKey 生成随机键值
func generateRandomKey() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	keyLen := 10
	b := make([]byte, keyLen)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}