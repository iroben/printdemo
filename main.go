package main

import (
  "crypto/md5"
  "encoding/json"
  "fmt"
  "io/ioutil"
  "log"
  "net/http"
  "net/url"
  "os"
  "strconv"
  "time"
)

type Setting struct {
  AccessKey string
  SecretKey string
}

func (s Setting) Sign(unixTime int, printerId string, templateId string) string {
  return fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%s-%d-%s-%s-%s",
    s.AccessKey, unixTime, printerId, templateId, s.SecretKey))))
}

const (
  DOMAIN = "https://dayin-api.78900c.com"
)

type Resp struct {
  StatusCode int    `json:"statusCode"`
  Msg        string `json:"msg"`
}

func (r Resp) Success() bool {
  return r.StatusCode == 10000
}

func main() {
  fp, err := os.Open("data.json")
  if err != nil {
    log.Fatalln("数据文件不存在：", err)
    return
  }
  defer fp.Close()
  bt, _ := ioutil.ReadAll(fp)
  log.Println("post data: ", string(bt))
  setting := Setting{
    AccessKey: "ak2a37182803a1413a3695236c9c15b8ca",
    SecretKey: "sk63557a40f2860a7cd474c954df6dd124",
  }
  values := url.Values{}
  unixTime := int(time.Now().Unix())
  printerId := "19"
  templateId := "96"
  values.Add("access_key", setting.AccessKey)
  values.Add("time", strconv.Itoa(unixTime))
  values.Add("print_id", printerId)     //打印机ID
  values.Add("template_id", templateId) //模板ID
  values.Add("sign", setting.Sign(unixTime, printerId, templateId))
  values.Add("data", string(bt))

  resp, err := http.PostForm(DOMAIN+"/api/print?", values)
  if err != nil {
    log.Fatalln("请求失败：", err)
    return
  }
  defer resp.Body.Close()
  bt, _ = ioutil.ReadAll(resp.Body)
  log.Println("响应数据：", string(bt))
  var respBody Resp
  if err := json.Unmarshal(bt, &respBody); err != nil {
    log.Println("err：", err)
    return
  }
  if !respBody.Success() {
    log.Println("打印失败：", respBody.Msg)
    return
  }
  log.Println("打印成功...")
}
