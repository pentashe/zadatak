package saxpy

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"sync"
	"zadatak/auth"

	"gonum.org/v1/gonum/blas/blas32"
)

const N uint = 3

var dataComplete sync.WaitGroup
var saxpyComplete sync.WaitGroup

var sharedMap = struct {
	sync.RWMutex
	m map[string]map[uint]float32
}{m: make(map[string]map[uint]float32)}

var result string

func Init() {
	dataComplete.Add(int(2*N + 1))
	saxpyComplete.Add(1)

	dataComplete.Wait()
	result = saxpy(&sharedMap.m)
	saxpyComplete.Done()
}

func Handler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "text/plain")
	if request.Method == http.MethodPost && auth.BasicAuth(&writer, request) {
		var streamCopy bytes.Buffer
		stream := io.TeeReader(request.Body, &streamCopy)

		// verify received json
		var fetchedStruct JSON
		decoder := json.NewDecoder(stream)
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&fetchedStruct); err != nil || fetchedStruct.invalidA() {
			http.Error(writer, "JSON Format Error", http.StatusBadRequest)
			return
		} else if fetchedStruct.indexGreaterThanN() {
			http.Error(writer, "X or Y index greater than N", http.StatusBadRequest)
			return
		}

		// decode verified json into map object
		fetchedMap := make(map[string]map[uint]float32)
		_ = json.NewDecoder(&streamCopy).Decode(&fetchedMap)

		// copy fetched map into shared map while incrementing dataComplete semaphore
		// in case of newly added (not overwritten) data
		for axy, values := range fetchedMap {
			sharedMap.RLock()
			_, exist := sharedMap.m[axy]
			sharedMap.RUnlock()
			if !exist {
				sharedMap.Lock()
				sharedMap.m[axy] = make(map[uint]float32)
				sharedMap.Unlock()
			}
			for i := range values {
				sharedMap.RLock()
				_, exist := sharedMap.m[axy][i]
				sharedMap.RUnlock()
				if !exist {
					sharedMap.Lock()
					sharedMap.m[axy][i] = fetchedMap[axy][i]
					sharedMap.Unlock()
					dataComplete.Done()
				}
			}
		}
		saxpyComplete.Wait()
		_, _ = io.WriteString(writer, result)
	}
}

func saxpy(data *map[string]map[uint]float32) string {
	x := make([]float32, N)
	y := make([]float32, N)
	var i uint
	for i = 0; i < N; i++ {
		x[i] = (*data)["x"][i]
		y[i] = (*data)["y"][i]
	}

	vectorX := blas32.Vector{N: int(N), Inc: 1, Data: x}
	vectorY := blas32.Vector{N: int(N), Inc: 1, Data: y}
	blas32.Axpy((*data)["a"][0], vectorX, vectorY)

	result := "\t[\n\t  "
	for _, v := range vectorY.Data {
		entry := float64(v)
		if entry >= 0 {
			result += " "
		}
		result += strconv.FormatFloat(entry, 'G', 4, 32) + "\n\t  "
	}
	result = result[:len(result)-4]
	result += "\n\t]"
	return result
}

type JSON struct {
	A map[uint]float32 `json:"a"`
	X map[uint]float32 `json:"x"`
	Y map[uint]float32 `json:"y"`
}

func (j JSON) invalidA() bool {
	for i := range j.A {
		if i != 0 {
			return true
		}
	}
	return false
}

func (j JSON) indexGreaterThanN() bool {
	for i := range j.X {
		if i >= N {
			return true
		}
	}
	for i := range j.Y {
		if i >= N {
			return true
		}
	}
	return false
}
