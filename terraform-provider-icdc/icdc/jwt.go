package icdc

import (
	"encoding/json"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

type Jwt struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type JwtClaims struct {
	External struct {
		Groups []string `json:"groups"`
		Locations map[string]string `json:"locations"`
	} `json:"external"`
}

func getJwt(username, password, ssoUrl, ssoRealm, ssoClientId string) (Jwt, diag.Diagnostics) {
	var diags diag.Diagnostics

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	reqUrl := fmt.Sprintf("https://%s/realms/%s/protocol/openid-connect/token", ssoUrl, ssoRealm)
	data := url.Values{}

	data.Set("username", username)
	data.Set("password", password)
	data.Set("client_id", ssoClientId)
	data.Set("grant_type", "password")

	encodedData := data.Encode()

	req, err := http.NewRequest("POST", reqUrl, strings.NewReader(encodedData))

	if err != nil {
		return Jwt{}, append(diags, diag.FromErr(err)...)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	response, err := client.Do(req)
	if err != nil {
		return Jwt{}, append(diags, diag.FromErr(err)...)
	}

	defer response.Body.Close()

	var jwt Jwt

	body, err := io.ReadAll(response.Body)

	if err != nil {
		return Jwt{}, append(diags, diag.FromErr(err)...)
	}

	err = json.Unmarshal([]byte(body), &jwt)

	if err != nil {
		return Jwt{}, append(diags, diag.FromErr(err)...)
	}

	return jwt, nil
}

func (j Jwt) Claims() (JwtClaims, diag.Diagnostics) {
	var diags diag.Diagnostics

	accessTokenBody := strings.Split(j.AccessToken, ".")[1]
	accessTokenBody += strings.Repeat("=", ((4 - len(accessTokenBody)%4) % 4))

	decodedBody, err := base64.StdEncoding.DecodeString(accessTokenBody)

	if err != nil {
		return JwtClaims{}, append(diags, diag.FromErr(err)...)
	}

	var claims JwtClaims
	err = json.Unmarshal(decodedBody, &claims)

	if err != nil {
		panic(err)
	}

	return claims, nil
}
