package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"

	"github.com/gin-gonic/gin"
	cp "github.com/otiai10/copy"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/drive/v3"
)

type body struct {
	Name   string   `json:"name"`
	Skills []string `json:"skills"`
}

func getDirectoryContents(path string, skills []string) []string {
	entries, err := os.ReadDir(path)
	i := 0
	var skillset []string
	handleerror(err)
	for _, entry := range entries {
		if i < len(skills) && entry.Name() == skills[i] {
			//fmt.Println("=>" + skills[i])
			skillset = append(skillset, skills[i])
			i = i + 1
		}
	}
	fmt.Println(skillset)
	return skillset
}
func createDirectory(username string, skillset []string) {
	err := os.Mkdir(username, 0755)
	handleerror(err)
	for _, skill := range skillset {
		src := "./icons/" + skill
		err := cp.Copy(src, username)
		handleerror(err)
	}
}
func handleerror(err error) {
	if err != nil {
		panic(err)
	}
}

func getMyicons(c *gin.Context) {
	var user body
	err := c.BindJSON(&user)
	handleerror(err)
	var skills = user.Skills
	sort.Strings(skills)
	path := "./icons"
	skillset := getDirectoryContents(path, skills)
	createDirectory(user.Name, skillset)
	handleerror(err)
	files := uploadDir(user.Name)

	c.IndentedJSON(http.StatusCreated, files)
}

// Use Service account
func ServiceAccount(secretFile string) *http.Client {
	b, err := ioutil.ReadFile(secretFile)
	if err != nil {
		log.Fatal("error while reading the credential file", err)
	}
	var s = struct {
		Email      string `json:"client_email"`
		PrivateKey string `json:"private_key"`
	}{}
	json.Unmarshal(b, &s)
	config := &jwt.Config{
		Email:      s.Email,
		PrivateKey: []byte(s.PrivateKey),
		Scopes: []string{
			drive.DriveScope,
		},
		TokenURL: google.JWTTokenURL,
	}
	client := config.Client(context.Background())
	return client
}

func createFile(service *drive.Service, name string, mimeType string, content io.Reader, parentId string) (*drive.File, error) {
	f := &drive.File{
		MimeType: mimeType,
		Name:     name,
		Parents:  []string{parentId},
	}
	file, err := service.Files.Create(f).Media(content).Do()

	if err != nil {
		log.Println("Could not create file: " + err.Error())
		return nil, err
	}

	return file, nil
}
func createFolder(service *drive.Service, name string, parentId string) (*drive.File, error) {
	d := &drive.File{
		Name:     name,
		MimeType: "application/vnd.google-apps.folder",
		Parents:  []string{parentId},
	}

	file, err := service.Files.Create(d).Do()

	if err != nil {
		log.Println("Could not create dir: " + err.Error())
		return nil, err
	}

	return file, nil
}
func uploadDir(username string) []string {
	var dirPath string = "./" + username
	client := ServiceAccount("client_secret.json")
	ParentfolderId := "1ehOv0dTV8Rz2RRArwpXsgVstd7_iPrW-"
	srv, err := drive.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve drive Client %v", err)
	}
	dir, err := createFolder(srv, username, ParentfolderId)
	dirId := dir.Id
	handleerror(err)
	entries, err := os.ReadDir(dirPath)
	handleerror(err)
	var files []string
	for _, entry := range entries {
		go uploadfile(entry.Name(), username, dirId)
		files = append(files, entry.Name())
	}
	return files
}

// Step 1: Open  file
func uploadfile(filename string, username string, folderId string) {
	filepath := "./" + username + "/" + filename
	f, err := os.Open(filepath)
	if err != nil {
		panic(fmt.Sprintf("cannot open file: %v", err))
	}
	defer f.Close()
	// Step 2: Get the Google Drive service
	client := ServiceAccount("client_secret.json")
	srv, err := drive.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve drive Client %v", err)
	}
	// Step 4: create the file and upload
	file, err := createFile(srv, filename, "application/octet-stream", f, folderId)
	if err != nil {
		panic(fmt.Sprintf("Could not create file: %v\n", err))
	}
	fmt.Printf("File '%s' successfully uploaded", file.Name)
	fmt.Printf("\nFile Id: '%s' ", file.Id)
}
func getskills(c *gin.Context) {
	var path string
	path = "./icons"
	entries, err := os.ReadDir(path)
	handleerror(err)
	var list []string
	for _, entry := range entries {
		list = append(list, entry.Name())
	}
	c.IndentedJSON(http.StatusCreated, list)

}
func main() {

	router := gin.Default()
	router.POST("/getMyicons", getMyicons)
	router.GET("/getskills", getskills)

	router.Run(":8080")

}
