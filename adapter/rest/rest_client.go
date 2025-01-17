package rest

import (
	"bytes"
	"craftgate-go-client/model"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func SendRequest(req *http.Request, v interface{}, opts model.RequestOptions) error {
	client := &http.Client{
		Timeout: time.Minute,
	}

	body := ""
	if req.Body != nil {
		var buf bytes.Buffer
		tee := io.TeeReader(req.Body, &buf)
		req.Body = ioutil.NopCloser(&buf)
		bodyBytes, bodyErr := ioutil.ReadAll(tee)
		if bodyErr == nil {
			body = fmt.Sprintf("%s", bodyBytes)
		}
	}
	randomStr := GenerateRandomString()
	hashStr := GenerateHash(req.URL.String(), opts.ApiKey, opts.SecretKey, randomStr, body)
	fmt.Println(req.URL.String())

	req.Header.Set(model.ApiKeyHeaderName, opts.ApiKey)
	req.Header.Set(model.RandomHeaderName, randomStr)
	req.Header.Set(model.AuthVersionHeaderName, model.AuthVersion)
	req.Header.Set(model.ClientVersionHeaderName, model.ClientVersion)
	req.Header.Set(model.SignatureHeaderName, hashStr)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json; charset=utf-8")

	res, err := client.Do(req)
	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		var errRes model.Response[any]
		if err = json.NewDecoder(res.Body).Decode(&errRes); err == nil {
			return errors.New(errRes.Errors.ErrorDescription)
		}

		return fmt.Errorf("unknown error, status code: %d", res.StatusCode)
	}

	if _, ok := v.(*model.Void); ok {
		return nil
	}
	if err = json.NewDecoder(res.Body).Decode(&v); err != nil {
		return err
	}

	return nil
}

func GenerateHash(url, apiKey, secretKey, randomString, body string) string {
	hashStr := strings.Join([]string{url, apiKey, secretKey, randomString, body}, "")
	hash := sha256.New()
	hash.Write([]byte(hashStr))
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}

func GenerateRandomString() string {
	s := strconv.FormatInt(time.Now().UnixNano(), 16)
	fmt.Println(s[8:])
	return s[8:]
}
