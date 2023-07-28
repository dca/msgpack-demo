# msgpack-demo

這個套件庫可以簡單地將 json 轉換為 messagepack

## 使用

這個包裡的主要功能是 Marshal 方法，可以將給定的 JSON 對象轉換為 MessagePack 格式。

以下是一個使用的範例：

```golang
    package main

    import (
        "fmt"
        msgpack "github.com/dca/msgpack-demo"
    )

    func main() {
        data := map[string]interface{}{
            "name":  "John Doe",
            "age":   35,
            "email": "john@doe.com",
        }

        mp := NewMsgpack()

        m1, err := mp.Marshal(jsonObj)
        
        if err != nil {
            fmt.Println("Error while marshalling:", err)
            return
        }
        fmt.Println("Marshaled data:", byteData)
    }
```
