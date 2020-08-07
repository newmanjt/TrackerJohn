package main

import (
	"fmt"
	// "github.com/DAddYE/vips"
	"github.com/daddye/vips"
	"github.com/google/uuid"
	"io/ioutil"
	"log"
	"net/http"
	// "net/url"
	"os"
	"strings"
)

// TrackerJohn
type Patient struct {

	//Personal Information
	FirstName string
	LastName  string
	Height    int
	Weight    int
	Location  string

	//Technical Information
	IPs        []string
	UserAgents []string

	//System Information
	ID uuid.UUID
}

func GoTo(s string, w http.ResponseWriter) {
	fmt.Fprintf(w, "<html><body><script>window.location.href='/"+s+"';</script></body></html>")
}

func loadFile(filename string) ([]byte, error) {
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return body, err
}

func checkUser(user string, password string) bool {
	body, err := ioutil.ReadFile("./logins/" + user)
	if err != nil {
		return false
	}
	log.Println(string(body))
	pdata := string(body)
	parts := strings.Split(pdata, ",")
	stored_pword := strings.Trim(parts[1], "\r\n")
	if stored_pword == password {
		return true
	} else {
		return false
	}
}

func getName(user string) string {
	body, err := ioutil.ReadFile("./logins/" + user)
	if err != nil {
		return "unknown"
	}
	log.Println(string(body))
	pdata := string(body)
	parts := strings.Split(pdata, ",")
	name := strings.Trim(parts[2], "\r\n")
	return name
}

func serveImage(name string, resize bool) (outBuf []byte) {
	modelPicFile, err := os.Open(name)
	if err != nil {
		return
	}
	defer modelPicFile.Close()

	options := vips.Options{
		Width:        500,
		Height:       300,
		Crop:         false,
		Extend:       vips.EXTEND_WHITE,
		Interpolator: vips.BILINEAR,
		Gravity:      vips.CENTRE,
		Quality:      95,
	}

	if resize {
		inBuf, _ := ioutil.ReadAll(modelPicFile)
		outBuf, err = vips.Resize(inBuf, options)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
	} else {
		outBuf, _ = ioutil.ReadAll(modelPicFile)
	}
	return
}

func getContentType(filename string) (contentType string) {
	if strings.HasSuffix(filename, ".eot") {
		contentType = "application/vnd.ms-fontobject"
	} else if strings.HasSuffix(filename, ".otf") {
		contentType = "application/x-font-opentype"
	} else if strings.HasSuffix(filename, ".svg") {
		contentType = "image/svg+xml"
	} else if strings.HasSuffix(filename, ".ttf") {
		contentType = "application/x-font-ttf"
	} else if strings.HasSuffix(filename, ".woff") {
		contentType = "application/font-woff"
	} else if strings.HasSuffix(filename, ".woff2") {
		contentType = "application/font-woff2"
	} else if strings.HasSuffix(filename, ".css") {
		// contentType = "text/css"
		contentType = ""
	} else if strings.HasSuffix(filename, ".png") {
		contentType = "image/png"
	} else if strings.HasSuffix(filename, ".jpg") {
		contentType = "image/jpeg"
	} else if strings.HasSuffix(filename, ".js") {
		contentType = "text/javascript"
	}
	return
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling visit")
	if r.URL.Path[1:] == "dashboard" {
		log.Println("Serving Dashboard")
		users, ok := r.URL.Query()["user"]
		if !ok || len(users[0]) < 1 {
			log.Println("Error, no user")
			fmt.Fprintf(w, "Error, Not a User")
			return
		}
		user := users[0]
		name := getName(user)
		dashboard, err := loadFile("dashboard.html")
		if err != nil {
			//TODO: change to error page
			fmt.Fprintf(w, "Error")
		}
		named_dashboard := strings.Replace(string(dashboard), "FIRSTNAME", name, 1)
		fmt.Fprintf(w, named_dashboard)
	} else if r.URL.Path[1:] == "login" {
		log.Println("Serving Login")
		login, err := loadFile("login.html")
		if err != nil {
			//TODO: change to error page
			fmt.Fprintf(w, "Error")
		}
		fmt.Fprintf(w, string(login))
	} else if r.URL.Path[1:] == "" {
		//base case, redirect to login
		log.Println("Redirecting to Login")
		GoTo("login", w)
	} else if r.URL.Path[1:] == "check" {
		//check if login info is correct
		// if so, redirect to dashboard with user info
		log.Println("Checking Login Information")
		users, ok := r.URL.Query()["id"]
		if !ok || len(users[0]) < 1 {
			log.Println("URL param user is missing")
			fmt.Fprintf(w, "Error logging in")
			return
		}
		user := users[0]
		passwords, ok := r.URL.Query()["password"]
		if !ok || len(passwords[0]) < 1 {
			log.Println("URL param password is missing")
			fmt.Fprintf(w, "Error logging in")
			return
		}
		password := passwords[0]
		log.Println("Logging in user " + user)
		log.Println("     with password " + password)
		is_user := checkUser(user, password)
		if is_user {
			//success logging in, redirect to dashboard
			GoTo("dashboard?user="+user, w)
		} else {
			fmt.Fprintf(w, "Not a user")
		}
	} else if strings.HasPrefix(r.URL.Path[1:], "dashboard_files/") {
		//serving necessary files for dashboard
		log.Println("Serving dashboard file: " + r.URL.Path[1:])
		if strings.HasSuffix(r.URL.Path[1:], ".png") || strings.HasSuffix(r.URL.Path[1:], "jpg") {
			w.Header().Set("Content-Type", getContentType(r.URL.Path[1:]))
			w.Write(serveImage(r.URL.Path[1:], false))
		} else if strings.HasSuffix(r.URL.Path[1:], "css") || strings.HasSuffix(r.URL.Path[1:], ".woff") || strings.HasSuffix(r.URL.Path[1:], "ttf") || strings.HasSuffix(r.URL.Path[1:], "eot") || strings.HasSuffix(r.URL.Path[1:], "woff2") || strings.HasSuffix(r.URL.Path[1:], "svg") || strings.HasSuffix(r.URL.Path[1:], "js") {
			data, err := loadFile(r.URL.Path[1:])
			if err != nil {
				fmt.Fprintf(w, "Error")
				log.Println(err)
				return
			}
			w.Header().Set("Content-Type", getContentType(r.URL.Path[1:]))
			fmt.Fprintf(w, string(data))
		}
	} else {
		//check to see if we can serve the file,
		//	otherwise, redirect to login
	}
}

func main() {
	log.Println("Starting TrackerJohn...")
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
