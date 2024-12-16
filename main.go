package main

import (
  "net/http"
  "log"
)

func main() {

   http.HandleFunc("/text", text)


   log.Println("Server started :3000")
   http.ListenAndServe(":3000", nil)
}

func text (w http.ResponseWriter, r *http.Request) {
   content_type := r.Header.Get("Content-Type");

   if content_type != "application/json" {
      log.Println("wrong content type")
   }
}
