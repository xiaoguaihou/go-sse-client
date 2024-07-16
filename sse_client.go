package client

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type SSE_Callback func(error, *map[string]*bytes.Buffer)

var gSseClient *http.Client

func init() {
	gSseClient = &http.Client{
		Timeout: 30 * time.Second,
	}
}

var sseKeys = []string{"id:", "data:", "event:", "meta:"}

type SSE_CLIENT_STATE int

const (
	SSE_STATE_SCAN_TAG SSE_CLIENT_STATE = iota
	SSE_STATE_SCAN_VALUE
)

func getSessionMap() map[string]*bytes.Buffer {
	var result = map[string]*bytes.Buffer{}
	for _, key := range sseKeys {
		result[key] = &bytes.Buffer{}
	}
	return result
}

func getTag(line []byte, keys []string) (string, int) {
	if len(line) < len(keys) {
		return "", 0
	}

	for index, key := range keys {
		byteKey := []byte(key)
		if bytes.HasPrefix(line, byteKey) {
			return keys[index], len(byteKey)
		}
	}
	return "", 0
}

func doCallback(session *map[string]*bytes.Buffer, callback SSE_Callback) {
	callback(nil, session)
}
func GetSSE(url string, header map[string]string, request interface{}, stream SSE_Callback) error {
	var requestBody []byte
	if request != nil {
		d, err := json.Marshal(request)
		if err != nil {
			return err
		}
		requestBody = d
	}

	httpReq, _ := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	httpReq.Header.Set("Content-type", "application/json")

	for k, v := range header {
		httpReq.Header.Set(k, v)
	}

	httpReq.Header.Add("Accept", "text/event-stream")

	resp, err := gSseClient.Do(httpReq)
	if err != nil {
		stream(err, nil)
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		msg := fmt.Sprintf("http error:%d, %s", resp.StatusCode, resp.Status)
		stream(fmt.Errorf(msg), nil)
		return fmt.Errorf(msg)

	} else {
		reader := bufio.NewReader(resp.Body)
		sessionMap := getSessionMap()

		for {
			line, err := reader.ReadBytes('\n')
			if err != nil {
				if err == io.EOF {
					doCallback(&sessionMap, stream)
				}
				break
			}
			line = line[:len(line)-1]
			if len(line) >= 1 && line[len(line)-1] == '\r' {
				line = line[:len(line)-1]
			}
			if len(line) > 0 {
				tag, l := getTag(line, sseKeys)
				if l > 0 {

					buffer, ok := sessionMap[tag]
					if ok {
						if buffer.Len() > 0 {
							buffer.WriteString("\n")
						}
						buffer.Write(line[l:])
					} else {
						break
					}
				}
			} else {
				doCallback(&sessionMap, stream)
				sessionMap = getSessionMap()
			}
		}
	}
	return nil
}
