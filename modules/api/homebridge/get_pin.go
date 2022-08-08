package homebridge

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"github.com/zhiting-tech/smartassistant/modules/api/utils/oauth"
	"github.com/zhiting-tech/smartassistant/modules/api/utils/response"
	"github.com/zhiting-tech/smartassistant/modules/config"
	"github.com/zhiting-tech/smartassistant/modules/entity"
	"github.com/zhiting-tech/smartassistant/modules/types"
	"github.com/zhiting-tech/smartassistant/modules/types/status"
	"github.com/zhiting-tech/smartassistant/modules/utils/session"
	"github.com/zhiting-tech/smartassistant/pkg/errors"
	"gopkg.in/oauth2.v3"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
)

type GetPinResponse struct {
	Pin string `json:"pin"`
}

func GetPin(c *gin.Context) {
	var (
		err  error
		resp GetPinResponse
	)

	defer func() {
		response.HandleResponse(c, err, resp)
	}()

	u := session.Get(c)
	if u == nil {
		err = errors.New(status.InvalidUserCredentials)
		return
	}

	// 创建homebridge 应用
	client, err := entity.CreateClient(oauth2.ClientCredentials.String(), types.WithScopes(types.ScopeDevice), u.AreaID)
	if err != nil {
		return
	}

	// 创建homebridge 所需token
	token, err := oauth.GetClientToken(client)
	if err != nil {
		return
	}

	param := map[string]interface{}{
		"token":   token,
		"area_id": strconv.FormatUint(u.AreaID, 10),
	}

	b, err := reqToHomeBridge(c, http.MethodPost, getHomeBridgeApi("pin"), param)
	if err != nil {
		err = errors.Wrap(err, status.GeneratePinError)
		return
	}

	if gjson.GetBytes(b, "status").Int() != 0 {
		err = errors.Wrap(err, status.GeneratePinError)
		return
	}

	resp.Pin = gjson.GetBytes(b, "data").Get("pin").String()
	return
}

func reqToHomeBridge(c *gin.Context, method, url string, param map[string]interface{}) (b []byte, err error) {

	var content []byte

	if len(param) != 0 {
		content, _ = json.Marshal(param)
	}

	cli := http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", getHomeBridgeSockAddr())
			},
		},
	}

	req, err := http.NewRequest(method, url, bytes.NewReader(content))
	if err != nil {
		return
	}

	req.Header.Set(types.SATokenKey, c.GetHeader(types.SATokenKey))
	rsp, err := cli.Do(req)
	if err != nil {
		return
	}

	defer rsp.Body.Close()

	b, err = ioutil.ReadAll(rsp.Body)
	if err != nil {
		return
	}
	return
}

func getHomeBridgeSockAddr() string {
	return fmt.Sprintf("%s/socks/homebridge.sock", config.GetConf().SmartAssistant.RuntimePath)
}

func getHomeBridgeApi(path string) string {
	// 使用 ~ 从根目录开始索引
	// 避免 https:///mnt/data/zt-smartassistant/ 无法识别
	return fmt.Sprintf("http://~%s/%s", getHomeBridgeSockAddr(), path)
}
