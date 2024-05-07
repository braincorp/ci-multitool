package jira

import (
	gojira "github.com/andygrunwald/go-jira"
)

type CommonArgs struct {
	InstanceUrl string
	Project     string
	User        string
	Password    string
	Board       int
}

func GetClient(args *CommonArgs) (*gojira.Client, error) {
	tp := gojira.BasicAuthTransport{
		Username: args.User,
		Password: args.Password,
	}
	return gojira.NewClient(tp.Client(), args.InstanceUrl)
}
