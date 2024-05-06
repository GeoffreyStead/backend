package Handler

import (
		"https://github.com/GeoffreyStead/backend.git"
		"net/http"
		"fmt"
)

func Handler(w http.ResponseWriter, r*http.Request){
	server := New()

	//Hello world
	server.get("/", func(context *Context)){
		context.JSON(200, H{
			"message":"hello from vercel!!"
		})
}
}