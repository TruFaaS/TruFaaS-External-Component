package fission

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func Create(respWriter http.ResponseWriter, req *http.Request) {

	var function Function

	// get the json value and convert to struct
	err := json.NewDecoder(req.Body).Decode(&function)
	if err != nil {
		fmt.Println(err)
		return
	}

	//send a json response back
	marshal, err := json.Marshal(function)
	if err != nil {
		return
	}
	respWriter.Header().Set("Content-Type", "application/json")
	respWriter.WriteHeader(http.StatusOK)
	_, err = respWriter.Write(marshal)
	if err != nil {
		return
	}

	//fmt.Printf("%#v", function)

}
