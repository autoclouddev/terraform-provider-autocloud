package autocloud_sdk

import (
	"context"
	"encoding/json"
	"os"
)

type User struct {
	Me struct {
		ID    string `json:"ID"`
		Name  string `json:"Name"`
		Email string `json:"Email"`
	}
}

/*
	query {
		me {
			id
			name
			email
		}
	}
*/

func (c *Client) GetMe() (*User, error) {

	var q struct {
		Me struct {
			ID    string //graphql.ID
			Name  string
			Email string
		}
	}

	err := c.graphql.Query(context.Background(), &q, nil)
	if err != nil {
		return nil, err
	}

	// me := &User{}
	// //err = json.Unmarshal(q.Me, &me)
	// if err != nil {
	// 	return nil, err
	// }

	return (*User)(&q), nil
	//fmt.Println(q)
	//print(q)

	//return nil

}

// func UnmarshalGraphQL(data []byte, v interface{}) error {
// 	dec := json.NewDecoder(bytes.NewReader(data))
// 	dec.UseNumber()
// 	err := (&decoder{tokenizer: dec}).Decode(v)
// 	if err != nil {
// 		//ERROR IS THROWN HERE HERE
// 		return err
// 	}
// }

// print pretty prints v to stdout. It panics on any error.
func print(v interface{}) {
	w := json.NewEncoder(os.Stdout)
	w.SetIndent("", "\t")
	err := w.Encode(v)
	if err != nil {
		panic(err)
	}
}
