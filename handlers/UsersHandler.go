package handlers

import (
	"net/http"
)

var (
	UsersQuery = `
        SELECT 
			username
        FROM
            User			
`
)

func UsersHandler(w http.ResponseWriter, r *http.Request) {

}
