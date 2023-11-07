package pwrgo

import (
    "log"
    "net/http"
	"io/ioutil"
	"bytes"
)

func get(url string) (response string) {
   resp, err := http.Get(url)
   if err != nil {
      log.Fatalln(err)
   }

   body, err := ioutil.ReadAll(resp.Body)
   if err != nil {
      log.Fatalln(err)
   }
   response = string(body)
   return
}

func post(url string, jsonStr string) string {
    var jsonBytes = []byte(jsonStr)
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBytes))
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    body, _ := ioutil.ReadAll(resp.Body)
    return string(body)
}