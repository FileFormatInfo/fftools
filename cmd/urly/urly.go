package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/spf13/pflag"
)

var (
	BUILDER = "unknown"
	COMMIT  = "(local)"
	LASTMOD = "(local)"
	VERSION = "internal"
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
	Hostname string              `json:"hostname"`
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

func toJson(theUrl *url.URL, pretty bool) string {
	var theHost string
	if theUrl.Port() != "" {
		theHost = fmt.Sprintf("%s:%s", theUrl.Hostname(), theUrl.Port())
	} else {
		theHost = theUrl.Hostname()
	}
	urlJson := UrlJson{
		Scheme:   theUrl.Scheme,
		Hostname: theUrl.Hostname(),
		Host:     theHost,
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
	var jsonStr []byte
	var err error
	if pretty {
		jsonStr, err = json.MarshalIndent(urlJson, "", "  ")
	} else {
		jsonStr, err = json.Marshal(urlJson)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: unable to marshal URL to JSON: %v\n", err)
		os.Exit(1)
	}
	return string(jsonStr)
}

func aliasFunc(f *pflag.FlagSet, name string) pflag.NormalizedName {
	switch name {
	case "user":
		name = "username"
		break
	}
	return pflag.NormalizedName(name)
}

// Detailed help text
var helpText = `urly: A URL parsing and processing tool.`

func main() {

	var scheme = pflag.String("scheme", "", "Set the URL scheme")
	var noScheme = pflag.Bool("no-scheme", false, "Remove the URL scheme")

	var envUsername = pflag.String("username-env", "", "Environment variable containing the username for URL processing")
	var textUsername = pflag.String("username", "", "Username for URL processing")
	var noUsername = pflag.Bool("no-username", false, "Remove the username from the URL")

	var envPassword = pflag.String("password-env", "", "Environment variable containing the password for URL processing")
	var stdinPassword = pflag.Bool("password-stdin", false, "Read password from standard input")
	var noPassword = pflag.Bool("no-password", false, "Remove the password from the URL")

	var hostname = pflag.String("hostname", "", "Set the URL hostname")
	var noHostname = pflag.Bool("no-hostname", false, "Remove the URL hostname")
	var port = pflag.String("port", "", "Set the URL port")
	var noPort = pflag.Bool("no-port", false, "Remove the URL port")
	var path = pflag.String("path", "", "Set the URL path")
	var noPath = pflag.Bool("no-path", false, "Remove the URL path")
	var query = pflag.String("query", "", "Set the URL query")
	var noQuery = pflag.Bool("no-query", false, "Remove the URL query")
	var fragment = pflag.String("fragment", "", "Set the URL fragment")
	var noFragment = pflag.Bool("no-fragment", false, "Remove the URL fragment")

	var addParams = pflag.StringArray("addparam", []string{}, "Add a query parameter (key=value)")
	var setParams = pflag.StringArray("setparam", []string{}, "Set a query parameter (key=value)")
	var delParams = pflag.StringArray("delparam", []string{}, "Delete a query parameter (key)")

	var envUrl = pflag.String("url-env", "", "Environment variable containing the URL to process")

	var output = pflag.String("output", "url", "Output type: url, scheme, host, port, path, query, fragment, userinfo, username, password, segments, segment[N]")
	var newline = pflag.Bool("newline", false, "Append newline to output")

	var help = pflag.Bool("help", false, "Detailed help")
	var version = pflag.Bool("version", false, "Version info")

	pflag.CommandLine.SetNormalizeFunc(aliasFunc)
	pflag.Parse()

	if *version {
		fmt.Printf("urly version %s (built on %s from %s by %s)\n", VERSION, LASTMOD, COMMIT, BUILDER)
		return
	}

	if *help {
		fmt.Println("urly - manipulate URLs")
		pflag.PrintDefaults()
		fmt.Println()
		fmt.Println("Use `man urly` for detailed help.")
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

	if *noScheme {
		theUrl.Scheme = ""
	} else if *scheme != "" {
		theUrl.Scheme = *scheme
	}

	if *noUsername {
		if theUrl.User != nil {
			thePassword, hasPassword := theUrl.User.Password()
			if hasPassword {
				theUrl.User = url.UserPassword("", thePassword)
			} else {
				theUrl.User = nil
			}
		}
	} else if *envUsername != "" {
		theUrl.User = setUserName(theUrl.User, os.Getenv(*envUsername))
	} else if *textUsername != "" {
		theUrl.User = setUserName(theUrl.User, *textUsername)
	}

	if *noPassword {
		if theUrl.User != nil {
			theUsername := theUrl.User.Username()
			theUrl.User = url.User(theUsername)
		} else {
			theUrl.User = nil
		}
	} else if *envPassword != "" {
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

	if *noScheme {
		theUrl.Scheme = ""
	} else if *scheme != "" {
		theUrl.Scheme = *scheme
	}

	if *noHostname {
		if theUrl.Port() != "" {
			theUrl.Host = ":" + theUrl.Port()
		} else {
			theUrl.Host = ""
		}
	} else if *hostname != "" {
		if theUrl.Port() != "" {
			theUrl.Host = *hostname + ":" + theUrl.Port()
		} else {
			theUrl.Host = *hostname
		}
	}

	if *noPort {
		theUrl.Host = theUrl.Hostname()
	} else if *port != "" {
		theUrl.Host = fmt.Sprintf("%s:%s", theUrl.Hostname(), *port)
	}

	if *noPath {
		theUrl.Path = ""
	} else if *path != "" {
		theUrl.Path = *path
	}

	if *noQuery {
		theUrl.RawQuery = ""
	} else if *query != "" {
		theUrl.RawQuery = *query
	}

	if len(*setParams) > 0 {
		queryValues := theUrl.Query()
		for _, param := range *setParams {
			kv := strings.SplitN(param, "=", 2)
			if kv[0] != "" {
				queryValues.Set(kv[0], kv[1])
			}
		}
		theUrl.RawQuery = queryValues.Encode()
	}

	if len(*delParams) > 0 {
		queryValues := theUrl.Query()
		for _, key := range *delParams {
			queryValues.Del(key)
		}
		theUrl.RawQuery = queryValues.Encode()
	}

	if len(*addParams) > 0 {
		queryValues := theUrl.Query()
		for _, param := range *addParams {
			kv := strings.SplitN(param, "=", 2)
			if kv[0] != "" {
				if len(kv) > 1 {
					queryValues.Add(kv[0], kv[1])
				} else {
					queryValues.Add(kv[0], "")
				}
			}
		}
		theUrl.RawQuery = queryValues.Encode()
	}

	if *noFragment {
		theUrl.Fragment = ""
	} else if *fragment != "" {
		theUrl.Fragment = *fragment
	}

	switch *output {
	case "url":
		fmt.Print(theUrl.String())
	case "scheme":
		fmt.Print(theUrl.Scheme)
	case "username", "user":
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
	case "json":
		fmt.Print(toJson(theUrl, true))
	case "jsonl":
		fmt.Print(toJson(theUrl, false))
	case "segments":
		fmt.Print(theUrl.Path[1:]) // remove leading slash
	default:
		if strings.HasPrefix(*output, "segment[") && strings.HasSuffix(*output, "]") {
			indexStr := (*output)[len("segment[") : len(*output)-1]
			var index int
			_, err := fmt.Sscanf(indexStr, "%d", &index)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: Invalid segment index: %s\n", indexStr)
				os.Exit(1)
			}
			segments := strings.Split(theUrl.Path, "/")
			if index < 0 || index >= len(segments) {
				fmt.Fprintf(os.Stderr, "ERROR: Segment index out of range: %d\n", index)
				os.Exit(1)
			}
			fmt.Print(segments[index])
		} else {
			fmt.Fprintf(os.Stderr, "ERROR: Unknown output type: %s\n", *output)
			os.Exit(1)
		}
	}
	if *newline {
		fmt.Println()
	}
}
