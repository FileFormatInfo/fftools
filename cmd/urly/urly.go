package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"

	"github.com/FileFormatInfo/fftools/internal"
	"github.com/spf13/pflag"
)

func setUserName(userInfo *url.Userinfo, username string) *url.Userinfo {
	if userInfo != nil {
		password, hasPassword := userInfo.Password()
		if hasPassword {
			return url.UserPassword(username, password)
		} else {
			return url.User(username)
		}
	} else {
		return url.User(username)
	}
}

func setPassword(userInfo *url.Userinfo, password string) *url.Userinfo {
	if userInfo != nil {
		username := userInfo.Username()
		return url.UserPassword(username, password)
	} else {
		return url.UserPassword("", password)
	}
}

type UrlJson struct {
	Scheme   string              `json:"scheme"`
	Host     string              `json:"host"`
	Port     string              `json:"port"`
	Path     string              `json:"path"`
	Query    string              `json:"query"`
	Fragment string              `json:"fragment"`
	Username string              `json:"username"`
	Password string              `json:"password"`
	Url      string              `json:"url"`
	Params   map[string][]string `json:"params,omitempty"`
}

func toJson(theUrl *url.URL) string {
	urlJson := UrlJson{
		Scheme:   theUrl.Scheme,
		Host:     theUrl.Hostname(),
		Port:     theUrl.Port(),
		Path:     theUrl.Path,
		Query:    theUrl.RawQuery,
		Fragment: theUrl.Fragment,
		Url:      theUrl.String(),
		Params:   theUrl.Query(),
	}
	if theUrl.User != nil {
		urlJson.Username = theUrl.User.Username()
		password, hasPassword := theUrl.User.Password()
		if hasPassword {
			urlJson.Password = password
		}
	}
	jsonStr, err := json.Marshal(urlJson)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: unable to marshal URL to JSON: %v\n", err)
		os.Exit(1)
	}
	return string(jsonStr)
}

// Detailed help text
var helpText = `urly: A URL parsing and processing tool.`

func main() {

	var scheme = pflag.String("scheme", "", "Set the URL scheme")
	var envPassword = pflag.String("password-env", "", "Environment variable containing the password for URL processing")
	var stdinPassword = pflag.Bool("password-stdin", false, "Read password from standard input")
	var envUrl = pflag.String("url-env", "", "Environment variable containing the URL to process")
	var envUsername = pflag.String("username-env", "", "Environment variable containing the username for URL processing")
	var textUsername = pflag.String("username", "", "Username for URL processing")
	//LATER: var format = pflag.String("format", "text", "Output format: text or json")
	var output = pflag.String("output", "url", "Output type: url, scheme, host, port, path, query, fragment, userinfo, username, password")
	var newline = pflag.Bool("newline", false, "Append newline to output")
	var help = pflag.Bool("help", false, "Detailed help")
	var version = pflag.Bool("version", false, "Version info")

	pflag.Parse()

	if *version {
		internal.PrintVersion("urly")
		return
	}

	if *help {
		fmt.Printf("%s\n", helpText)
		return
	}

	var theUrl *url.URL
	var parseErr error

	if *envUrl != "" {
		textUrl := os.Getenv(*envUrl)
		theUrl, parseErr = url.Parse(textUrl)
		if parseErr != nil {
			fmt.Fprintf(os.Stderr, "ERROR: unable to parse URL from environment variable %s: %v\n", *envUrl, parseErr)
			os.Exit(1)
		}
	} else {
		args := pflag.Args()
		if len(args) > 0 {
			theUrl, parseErr = url.Parse(args[0])
			if parseErr != nil {
				fmt.Fprintf(os.Stderr, "ERROR: Unable to parse URL from argument: %v\n", parseErr)
				os.Exit(1)
			}
			if len(args) > 1 {
				fmt.Fprintf(os.Stderr, "WARNING: Ignoring extra arguments (count=%d)\n", len(args)-1)
			}
		}
	}

	if theUrl == nil {
		theUrl = &url.URL{}
	}

	if *envUsername != "" {
		theUrl.User = setUserName(theUrl.User, os.Getenv(*envUsername))
	} else if *textUsername != "" {
		theUrl.User = setUserName(theUrl.User, *textUsername)
	}

	if *envPassword != "" {
		theUrl.User = setPassword(theUrl.User, os.Getenv(*envPassword))
	} else if *stdinPassword {
		var password string
		_, err := fmt.Fscanln(os.Stdin, &password)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: Unable to read password from stdin: %v\n", err)
			os.Exit(1)
		}
		theUrl.User = setPassword(theUrl.User, password)
	}

	if *scheme != "" {
		theUrl.Scheme = *scheme
	}

	switch *output {
	case "url":
		fmt.Println(theUrl.String())
	case "scheme":
		fmt.Print(theUrl.Scheme)
	case "host":
		fmt.Print(theUrl.Hostname())
	case "port":
		fmt.Print(theUrl.Port())
	case "path":
		fmt.Print(theUrl.Path)
	case "query":
		fmt.Print(theUrl.RawQuery)
	case "fragment":
		fmt.Print(theUrl.Fragment)
	case "userinfo":
		if theUrl.User != nil {
			fmt.Print(theUrl.User.Username())
			password, isSet := theUrl.User.Password()
			if isSet {
				fmt.Print(":")
				fmt.Print(password)
			}
		}
	case "username":
		if theUrl.User != nil {
			fmt.Print(theUrl.User.Username())
		}
	case "password":
		if theUrl.User != nil {
			password, hasPassword := theUrl.User.Password()
			if hasPassword {
				fmt.Print(password)
			}
		}
	case "json":
		fmt.Print(toJson(theUrl))
	default:
		fmt.Fprintf(os.Stderr, "ERROR: Unknown output type: %s\n", *output)
		os.Exit(1)
	}
	if *newline {
		fmt.Println()
	}
}
