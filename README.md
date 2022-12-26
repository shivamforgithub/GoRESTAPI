# GoRESTAPI
How to use this project step by step guide :

1. Check the env variables with "go env" command in powershell 
2. Check the GoPath and go into that directory . It's your WORKSPACE
3. RUN 
$ mkdir web-service-gin
$ cd web-service-gin
4. Run the go mod init command, giving it the path of the module your code will be in.
$ go mod init example/web-service-gin
5. Run the code 
$ go get .
$ go run .


In this project we will create a REST APIs using gin-gonic/gin package and will learn 

1. How to route GET/POST/PUT/DELETE/PATCH/HEAD request 
2. "http"  package and it's StatusCodes 
3. How to deal with data inside requests and send JSON response
