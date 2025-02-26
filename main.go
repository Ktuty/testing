package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

func Sum(v1, v2 int) int {
	return v1 + v2
}

func FindError(v int) error {
	if v%2 == 0 {
		return fmt.Errorf("%d is an error", v)
	}

	return nil
}

func CountWords(str string) map[string]int {
	re := regexp.MustCompile(`[.,!?;:'*/"` + "`" + `-]+`)
	str = re.ReplaceAllString(str, " ")

	words := strings.Fields(strings.ToLower(str))

	wordCount := make(map[string]int, len(words))

	for _, word := range words {
		wordCount[word]++
	}

	return wordCount
}

func MergeAndSortSlices(arr1, arr2 []int) []int {

	if len(arr1) == 0 && len(arr2) == 0 {
		return []int{}
	}

	arr := append(arr1, arr2...)

	sort.Ints(arr)

	uniqueArr := make([]int, 0, len(arr))
	uniqueArr = append(uniqueArr, arr[0])

	for i := 1; i < len(arr); i++ {
		if arr[i] != arr[i-1] {
			uniqueArr = append(uniqueArr, arr[i])
		}
	}

	return uniqueArr
}

func HandlerExample(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/example" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello World!"))
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not Found!"))
	}
}

// для работы со внешним апи
type Client interface {
	GetData(url string) (string, error)
}

type HttpClient struct {
	Client *http.Client
}

func (c *HttpClient) GetData(url string) (string, error) {
	req, err := c.Client.Get(url)
	if err != nil {
		return "", err
	}
	defer req.Body.Close()

	if req.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ошибка: неверный статус кода %d", req.StatusCode)
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return "", err
	}

	var result map[string]interface{}

	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	data, ok := result["data"].(string)
	if !ok {
		return "", fmt.Errorf("ошибка: неверный формат ответа")
	}

	return data, nil
}

func StartAndStop() time.Duration {
	start := time.Now()

	time.Sleep(10 * time.Second)

	return time.Since(start)
}

func Count(n int) int {

	if n <= 0 {
		return 0
	}

	count := 0
	wg := new(sync.WaitGroup)
	wg.Add(n)

	for i := 0; i < n; i++ {
		go func(int) {
			defer wg.Done()
			count++
		}(count)
	}

	wg.Wait()
	return count
}

func CheckPassword(pswrd string) error {
	if pswrd == "" {
		return fmt.Errorf("empty string")
	}

	if utf8.RuneCountInString(pswrd) < 5 {
		return fmt.Errorf("too short password")
	}

	if utf8.RuneCountInString(pswrd) > 15 {
		return fmt.Errorf("too long password")
	}

	re := regexp.MustCompile("[^a-zA-Z0-9]+")
	if re.MatchString(pswrd) {
		return fmt.Errorf("password must contain only alphanumeric characters")
	}

	return nil
}

// _______________ Task HTTP Server ____________________

func MethodGetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}

	fmt.Fprintln(w, "Hello World!")
}

type Message struct {
	Text string `json:"text"`
}

func MethodPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var msg Message
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(msg)
}

//func main() {
//	http.HandleFunc("/", MethodGetHandler)
//	http.HandleFunc("/example", MethodPostHandler)
//
//	http.ListenAndServe(":8080", nil)
//}
