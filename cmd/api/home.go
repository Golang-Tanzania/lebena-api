package main

import (
	"fmt"
	"net/http"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintln(w, "Thank you,to use this Lebena Sports API made specifically \nfor Tanzania premier league (NBC premier League) \nKindly visit https://sports-doc.eadevs.com/")

}
