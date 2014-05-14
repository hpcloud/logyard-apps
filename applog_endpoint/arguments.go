package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type Arguments struct {
	Token string
	GUID  string
	Num   int
}

func (args *Arguments) Validate() error {
	if args.Token == "" {
		return fmt.Errorf("missing token")
	}
	if args.GUID == "" {
		return fmt.Errorf("missing app guid")
	}
	if args.Num < 0 {
		return fmt.Errorf("Num argument is negative")
	}
	return nil
}

func ParseArguments(r *http.Request) (*Arguments, error) {
	var err error
	args := new(Arguments)
	vars := mux.Vars(r)

	args.Token = r.Header.Get("Authorization")
	args.GUID = vars["guid"]
	if r.FormValue("num") == "" {
		args.Num = 0
	} else {
		args.Num, err = strconv.Atoi(r.FormValue("num"))
		if err != nil {
			return nil, err
		}
	}

	return args, args.Validate()
}
