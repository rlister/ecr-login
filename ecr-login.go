package main

import (
	"encoding/base64"
	// "github.com/rlister/ecr-login/Godeps/_workspace/src/github.com/aws/aws-sdk-go/aws"
	"github.com/rlister/ecr-login/Godeps/_workspace/src/github.com/aws/aws-sdk-go/aws/session"
	"github.com/rlister/ecr-login/Godeps/_workspace/src/github.com/aws/aws-sdk-go/service/ecr"
	"os"
	"strings"
	"text/template"
	"time"
)

type Auth struct {
	Token         string
	User          string
	Pass          string
	ProxyEndpoint string
	ExpiresAt     time.Time
}

// error handler
func check(e error) {
	if e != nil {
		panic(e.Error())
	}
}

// default template prints docker login command
const DEFAULT_TEMPLATE = `{{range .}}docker login -u {{.User}} -p {{.Pass}} -e none {{.ProxyEndpoint}}
{{end}}`

// load template from file or use default
func getTemplate() *template.Template {
	var tmpl *template.Template
	var err error

	file, exists := os.LookupEnv("TEMPLATE")

	if exists {
		tmpl, err = template.ParseFiles(file)
	} else {
		tmpl, err = template.New("default").Parse(DEFAULT_TEMPLATE)
	}

	check(err)
	return tmpl
}

func main() {
	svc := ecr.New(session.New())

	// this would be how to get tokens for multiple registries
	// params := &ecr.GetAuthorizationTokenInput{
	// 	RegistryIds: []*string{
	// 		aws.String("123"),
	// 		aws.String("456"),
	// 	},
	// }
	resp, err := svc.GetAuthorizationToken(nil)
	check(err)

	// fields to send to template
	fields := make([]Auth, len(resp.AuthorizationData))
	for i, auth := range resp.AuthorizationData {

		// extract base64 token
		data, err := base64.StdEncoding.DecodeString(*auth.AuthorizationToken)
		check(err)

		// extract username and password
		token := strings.SplitN(string(data), ":", 2)

		// object to pass to template
		fields[i] = Auth{
			Token:         *auth.AuthorizationToken,
			User:          token[0],
			Pass:          token[1],
			ProxyEndpoint: *(auth.ProxyEndpoint),
			ExpiresAt:     *(auth.ExpiresAt),
		}
	}

	// run the template
	err = getTemplate().Execute(os.Stdout, fields)
	check(err)
}
