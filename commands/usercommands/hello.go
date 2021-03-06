package usercommands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"

	"github.com/starshine-sys/covebotnt/crouter"
)

const pkAPIversion = 1

type pkAPIResponse struct {
	Member struct {
		Name string `json:"name"`
	} `json:"member"`
}

var (
	greetings []string = []string{"Hello", "Heyo", "Hiya", "Heya"}
)

// Hello says hello to the user invoking the command
func Hello(ctx *crouter.Ctx) (err error) {
	var apiResponse pkAPIResponse
	var person string

	time.Sleep(1 * time.Second)

	resp, err := http.Get(fmt.Sprintf("https://api.pluralkit.me/v%v/msg/%v", pkAPIversion, ctx.Message.ID))
	if err != nil {
		ctx.CommandError(err)
		return nil
	}
	defer resp.Body.Close()
	if resp.StatusCode == 404 {
		person = ctx.Author.Mention()
	} else {
		apiRespBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			ctx.CommandError(err)
			return nil
		}

		json.Unmarshal(apiRespBytes, &apiResponse)
		if apiResponse.Member.Name == "" {
			person = ctx.Author.Mention()
		} else {
			person = apiResponse.Member.Name
		}
	}

	rand.Seed(time.Now().Unix())

	_, err = ctx.Send(fmt.Sprintf("%v, %v!", greetings[rand.Intn(len(greetings))], person))
	return
}
