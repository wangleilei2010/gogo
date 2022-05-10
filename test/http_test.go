package test

import (
	"fmt"
	"github.com/wangleilei2010/gogo/http"
	"testing"
	"time"
)

var Json map[string]interface{}

func TestHttpClient(t *testing.T) {
	url := "http://172.18.1.69/beetle/api/coreresource/i18n/getLangItems/v1?languageCode=zh_cn"
	client := http.NewSession(time.Second * 20)
	resp, _ := http.Get(client, url)
	fmt.Println(resp.StatusCode)

	resp2, _ := http.GetAndUnmarshal[map[string]interface{}](client, url)
	fmt.Println(resp2.StatusCode)
}
