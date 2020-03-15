package swyftx

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"math"
	"strconv"
	"strings"
	"time"
)

func buildString(strs ...string) string {
	var builtStr strings.Builder
	for _, str := range strs {
		builtStr.WriteString(str)
	}

	return builtStr.String()
}

func isEmptyStr(str string) bool {
	return str == ""
}

func copyReadCloser(rc io.ReadCloser) (*bytes.Buffer, error) {
	var body bytes.Buffer
	if _, err := io.Copy(&body, rc); err != nil {
		return nil, err
	}
	rc = ioutil.NopCloser(bytes.NewReader(body.Bytes()))

	return &body, nil
}

func decodeJSON(r io.Reader, v interface{}) error {
	if err := json.NewDecoder(r).Decode(v); err != nil {
		return err
	}
	return nil
}

type SwyftxTime struct {
	time.Time
}

func (s *SwyftxTime) UnmarshalJSON(b []byte) (err error) {
	if bytes.ContainsAny(b, "\"") {
		b = bytes.Trim(b, "\"")
	}

	var floatTime float64
	floatTime, err = strconv.ParseFloat(string(b), 32)
	if err != nil {
		return err
	}

	// https://stackoverflow.com/questions/37628254
	sec, dec := math.Modf(floatTime)
	s.Time = time.Unix(int64(sec), int64(dec*(1e9)))

	return nil
}
