package oauth

import (
	"encoding/json"
	"fmt"
	"os"

	"net/http"
	"strconv"

	"github.com/sijanstha/common-utils/src/client/rest"
	"github.com/sijanstha/common-utils/src/utils/errors"
)

const (
	headerXPublic      = "X-PUBLIC"
	headerXClientId    = "X-CLIENT-ID"
	headerXCallerId    = "X-CALLER-ID"
	headerXAccessToken = "X-ACCESS-TOKEN"
)

var (
	oauthAPIURL = os.Getenv("oauth_ms_url")
)

type accessToken struct {
	Id       string `json:"id"`
	UserId   int64  `json:"user_id"`
	ClientId int64  `json:"client_id"`
}

func IsPublic(request *http.Request) bool {
	if request == nil {
		return true
	}
	return request.Header.Get(headerXPublic) == "true"
}

func GetCallerId(request *http.Request) int64 {
	if request == nil {
		return 0
	}
	callerId, err := strconv.ParseInt(request.Header.Get(headerXCallerId), 10, 64)
	if err != nil {
		return 0
	}
	return callerId
}

func GetClientId(request *http.Request) int64 {
	if request == nil {
		return 0
	}
	clientId, err := strconv.ParseInt(request.Header.Get(headerXClientId), 10, 64)
	if err != nil {
		return 0
	}
	return clientId
}

func AuthenticateRequest(request *http.Request) *errors.RestErr {
	if request == nil {
		return nil
	}

	cleanRequest(request)

	accessTokenId := request.Header.Get(headerXAccessToken)
	if accessTokenId == "" {
		return errors.NewBadRequestError("Invalid access token")
	}

	at, err := getAccessToken(accessTokenId)
	if err != nil {
		return err
	}

	request.Header.Add(headerXCallerId, fmt.Sprintf("%v", at.UserId))
	request.Header.Add(headerXClientId, fmt.Sprintf("%v", at.ClientId))

	return nil
}

func cleanRequest(request *http.Request) {
	if request == nil {
		return
	}
	request.Header.Del(headerXClientId)
	request.Header.Del(headerXCallerId)
}

func getAccessToken(accessTokenId string) (*accessToken, *errors.RestErr) {
	response, err := rest.RestClient.R().Get(oauthAPIURL + "/" + accessTokenId)
	if err != nil {
		return nil, errors.NewInternalServerError(fmt.Sprintf("rest client error: %s", err.Error()))
	}

	if response == nil || response.Body() == nil {
		return nil, errors.NewInternalServerError("Invalid rest client response when try to get access token")
	}

	if response.StatusCode() > 299 {
		var restErr errors.RestErr
		err := json.Unmarshal(response.Body(), &restErr)
		if err != nil {
			return nil, errors.NewInternalServerError("Invalid error interface when trying to get access token")
		}
		return nil, &restErr
	}

	var token accessToken
	if err := json.Unmarshal(response.Body(), &token); err != nil {
		return nil, errors.NewInternalServerError("error when trying to unmarshal access token response")
	}

	return &token, nil
}
