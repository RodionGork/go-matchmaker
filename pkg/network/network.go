package network

import (
    "encoding/json"
    "io"
    "net/http"
    "reflect"
)

func writeResult(w http.ResponseWriter, statusCode int, result bool, msg string) {
    w.WriteHeader(statusCode)
    m := map[string]any{"result": result}
    if msg != "" {
        m["msg"] = msg
    }
    bytes, _ := json.Marshal(m)
    w.Write(bytes)
}

func HandleJSON(path string, handlerFunc func(map[string]any)string, fieldsVerifier map[string]string) {
    http.HandleFunc(path, func(w http.ResponseWriter, req *http.Request) {
        defer req.Body.Close()
        data, err := io.ReadAll(req.Body)
        if err != nil {
            writeResult(w, http.StatusInternalServerError, false, err.Error())
            return
        }
        body := map[string]any{}
        err = json.Unmarshal(data, &body)
        if err != nil {
            writeResult(w, http.StatusBadRequest, false, "Seemingly, broken JSON")
            return
        }
        for k, v := range fieldsVerifier {
            if field, ok := body[k]; !ok {
                writeResult(w, http.StatusBadRequest, false, "Missing field: " + k)
                return
            } else if reflect.TypeOf(field).Name() != v {
                writeResult(w, http.StatusBadRequest, false, "Wrong field type: " + k)
                return
            }
        }
        res := handlerFunc(body)
        writeResult(w, http.StatusOK, true, res)
    })
}

func Start(port string) {
    if port == "" {
        port = "8080"
    }
    println("Starting to listen at " + port)
    err := http.ListenAndServe(":" + port, nil)
    println("Exiting due to ", err.Error())
}
