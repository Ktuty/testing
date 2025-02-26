package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sort"
	"testing"
	"testing/quick"
	"time"
	//"testing/synctest"
)

func Test_Sum(t *testing.T) {
	type args struct {
		v1 int
		v2 int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Test 1",
			args: args{
				v1: 1,
				v2: 2,
			},
			want: 3,
		},
		{
			name: "Test 2",
			args: args{
				v1: 5,
				v2: 0,
			},
			want: 5,
		},
		{
			name: "Test 3",
			args: args{
				v1: -4,
				v2: 2,
			},
			want: -2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Sum(tt.args.v1, tt.args.v2); got != tt.want {
				t.Errorf("sum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFindError(t *testing.T) {
	type args struct {
		v int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test 1: odd value",
			args: args{
				v: 1,
			},
			wantErr: false,
		},
		{
			name: "Test 2: even value",
			args: args{
				v: 2,
			},
			wantErr: true,
		},
		{
			name: "Test 3: negative odd value",
			args: args{
				v: -1,
			},
			wantErr: false,
		},
		{
			name: "Test 4: negative even value",
			args: args{
				v: -2,
			},
			wantErr: true,
		},
		{
			name: "Test 4: nil",
			args: args{
				v: 0,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := FindError(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("FindError() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCountWords(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want map[string]int
	}{
		{
			name: "Test 1",
			args: args{
				str: "Привет, мир! Привет мир. Привет, \"мир\"! Привет - мир.",
			},
			want: map[string]int{
				"привет": 4,
				"мир":    4,
			},
		},
		{
			name: "Test 2",
			args: args{
				str: "Привет, мир! Привет мир. Привет, \"мир\"! Привет - мир.    мир рим рим РИМ    .....``/*/*/*/*/*/```/*-    мир     ; :      ",
			},
			want: map[string]int{
				"привет": 4,
				"мир":    6,
				"рим":    3,
			},
		},
		{
			name: "Test 3",
			args: args{
				str: "Привет, мир! Привет мир. Привет,\"мир\"!Привет - мир.    мир рим рим РИМ    .....``/*/*/*/*/*/```/*-    мир     ; :      ",
			},
			want: map[string]int{
				"привет": 4,
				"мир":    6,
				"рим":    3,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CountWords(tt.args.str); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CountWords() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMergeAndSortSlices(t *testing.T) {
	type args struct {
		arr1 []int
		arr2 []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "Test 1",
			args: args{
				arr1: []int{1, 3, 5, 7},
				arr2: []int{3, 5, 6, 8},
			},
			want: []int{1, 3, 5, 6, 7, 8},
		},
		{
			name: "Test 2",
			args: args{
				arr1: []int{1, 3, 5, 7, 3, 5, 6, 8},
				arr2: []int{3, 5, 6, 8, 1, 3, 5, 7},
			},
			want: []int{1, 3, 5, 6, 7, 8},
		},
		{
			name: "Test 3",
			args: args{
				arr1: []int{1, 1, 1, 1},
				arr2: []int{1, 1, 1, 1},
			},
			want: []int{1},
		},
		{
			name: "Test 4",
			args: args{
				arr1: []int{},
				arr2: []int{},
			},
			want: []int{},
		},
		{
			name: "Test 5",
			args: args{
				arr1: []int{},
				arr2: []int{1, 1, 1, 1},
			},
			want: []int{1},
		},
		{
			name: "Test 6",
			args: args{
				arr1: []int{1, 2, 3},
				arr2: []int{},
			},
			want: []int{1, 2, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MergeAndSortSlices(tt.args.arr1, tt.args.arr2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MergeAndSortSlices() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMergeAndSortSlicesIdempotence(t *testing.T) {
	fn := func(arr1, arr2 []int) bool {
		res := MergeAndSortSlices(arr1, arr2)
		return equal(res, MergeAndSortSlices(res, []int{}))
	}
	if err := quick.Check(fn, nil); err != nil {
		t.Error("Свойство коммутативности не выполняется:", err)
	}
}

func TestMergeAndSortSlicesCommutativity(t *testing.T) {
	fn := func(arr1, arr2 []int) bool {
		return equal(MergeAndSortSlices(arr1, arr2), MergeAndSortSlices(arr2, arr1))
	}

	if err := quick.Check(fn, nil); err != nil {
		t.Error("Свойство коммутативности не выполняется:", err)
	}
}

func TestMergeAndSortSlicesOrder(t *testing.T) {
	fn := func(arr1, arr2 []int) bool {
		return sort.IntsAreSorted(MergeAndSortSlices(arr1, arr2))
	}

	if err := quick.Check(fn, nil); err != nil {
		t.Error("Свойство сохранения порядка не выполняется:", err)
	}
}

func equal(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func TestHandlerExample(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/example", nil)
	w := httptest.NewRecorder()

	HandlerExample(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("неверный статус ответа, ожидание: %d реальность: %d\n", http.StatusOK, w.Code)
	}

	expected := "Hello World!"
	if w.Body.String() != expected {
		t.Errorf("неверный статус ответа, ожидание: %#v реальность: %#v\n", expected, w.Body.String())
	}
}

func TestHandlerExampleNotFound(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/unknown", nil)
	w := httptest.NewRecorder()

	HandlerExample(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("неверный статус ответа, ожидание: %d реальность: %d\n", http.StatusNotFound, w.Code)
	}

	expected := "Not Found!"
	if w.Body.String() != expected {
		t.Errorf("неверный статус ответа, ожидание: %#v реальность: %#v\n", expected, w.Body.String())
	}
}

type MockClient struct{}

func (m *MockClient) GetData(url string) (string, error) {
	return "моковые данные", nil
}

func TestGetDataMock(t *testing.T) {
	client := &MockClient{}

	data, err := client.GetData("http://test.api")
	if err != nil {
		t.Fatalf("ожидался nil, но получена ошибка: %v\n", err)
	}

	expected := "моковые данные"
	if data != expected {
		t.Errorf("ожидалось %s, но нолучено %s\n", expected, data)
	}
}

func TestGetDataHttp(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data": "тестовые данные"}`))
	}))
	defer mockServer.Close()

	client := &HttpClient{Client: &http.Client{}}

	data, err := client.GetData(mockServer.URL)
	if err != nil {
		t.Fatalf("ошибка при запросе: %v\n", err)
	}

	expected := "тестовые данные"
	if data != expected {
		t.Errorf("ожидалось '%s', но получено '%s'", expected, data)
	}
}

func TestStartAndStop(t *testing.T) {
	duration := StartAndStop() / time.Second

	if duration != 10 {
		t.Errorf("не достигнуто 10 сек, получено %d", duration)
	}
}

//func TestStartAndStop_2_24(t *testing.T) {
//	synctest.Run(func() {
//		duration := StartAndStop()
//
//		if duration < 10 {
//			t.Errorf("не достигнуто 10 сек, получено %d", duration)
//		}
//	})
//}

func TestCount(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Test 1",
			args: args{
				n: 5,
			},
			want: 5,
		},
		{
			name: "Test 2",
			args: args{
				n: 0,
			},
			want: 0,
		},
		{
			name: "Test 3",
			args: args{
				n: -5,
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Count(tt.args.n); got != tt.want {
				t.Errorf("Count() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckPassword(t *testing.T) {
	type args struct {
		pswrd string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		errIs   string
	}{
		{
			name: "Test 1",
			args: args{
				pswrd: "qwerty123",
			},
			wantErr: false,
			errIs:   "",
		},
		{
			name: "Test 2",
			args: args{
				pswrd: "",
			},
			wantErr: true,
			errIs:   "empty string",
		},
		{
			name: "Test 3",
			args: args{
				pswrd: "qwer",
			},
			wantErr: true,
			errIs:   "too short password",
		},
		{
			name: "Test 4",
			args: args{
				pswrd: "qwertyuiophjgsdfjghsdkjfhgsjdhgfkjfhgfjhgfdjg",
			},
			wantErr: true,
			errIs:   "too long password",
		},
		{
			name: "Test 5",
			args: args{
				pswrd: "qwerty@123",
			},
			wantErr: true,
			errIs:   "password must contain only alphanumeric characters",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CheckPassword(tt.args.pswrd); (err != nil) != tt.wantErr && err.Error() != tt.errIs {
				t.Errorf("CheckPassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMethodGetHandler(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatalf("error creating request: %v", err.Error())
	}
	w := httptest.NewRecorder()
	handler := http.HandlerFunc(MethodGetHandler)

	handler.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := "Hello World!\n"
	if w.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", w.Body.String(), expected)
	}
}

func TestMethodPostHandler(t *testing.T) {
	expected := "Test message"

	msg := Message{Text: expected}
	body, _ := json.Marshal(msg)

	req, err := http.NewRequest(http.MethodPost, "/example", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	w := httptest.NewRecorder()
	handler := http.HandlerFunc(MethodPostHandler)

	handler.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {

		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var respMsg Message
	if err := json.Unmarshal(w.Body.Bytes(), &respMsg); err != nil {
		t.Fatalf("error creating request: %v", err)
	}

	if respMsg.Text != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", respMsg.Text, expected)
	}
}

func TestMethodGetHandlerInvalid(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "/", nil)
	if err != nil {
		t.Fatalf("error creating request: %v", err)
	}

	w := httptest.NewRecorder()
	handler := http.HandlerFunc(MethodGetHandler)

	handler.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusMethodNotAllowed)
	}
}

func TestMethodPostHandlerInvalid(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/example", nil)
	if err != nil {
		t.Fatalf("error creating request: %v", err)
	}

	w := httptest.NewRecorder()
	handler := http.HandlerFunc(MethodPostHandler)

	handler.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusMethodNotAllowed)
	}
}

func TestMethodPostHandlerInvalidBody(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "/example", bytes.NewBufferString("invalid json"))
	if err != nil {
		t.Fatalf("error creating request: %v", err)
	}

	w := httptest.NewRecorder()
	handler := http.HandlerFunc(MethodPostHandler)

	handler.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}
