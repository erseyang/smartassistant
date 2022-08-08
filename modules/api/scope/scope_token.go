package scope

import (
	"gopkg.in/oauth2.v3"
	"gopkg.in/oauth2.v3/server"
	"strconv"
	"strings"
	"time"

	"github.com/zhiting-tech/smartassistant/modules/api/utils/oauth"
	"github.com/zhiting-tech/smartassistant/modules/entity"
	"github.com/zhiting-tech/smartassistant/modules/types/status"
	"github.com/zhiting-tech/smartassistant/pkg/cache"
	"github.com/zhiting-tech/smartassistant/pkg/logger"

	"github.com/gin-gonic/gin"

	"github.com/zhiting-tech/smartassistant/modules/api/utils/response"
	"github.com/zhiting-tech/smartassistant/modules/types"
	"github.com/zhiting-tech/smartassistant/modules/utils/session"
	"github.com/zhiting-tech/smartassistant/pkg/errors"
)

type token struct {
	Token     string `json:"token"`
	ExpiresIn int    `json:"expires_in"`
}

type scopeTokenResp struct {
	ScopeToken token `json:"scope_token"`
}

var (
	ExpiresIn = time.Hour * 24 * 30
)

type scopeTokenReq struct {
	AreaIDStr string `json:"area_id"`
	areaID    uint64
	Scopes    []string `json:"scopes"`
}

func (req *scopeTokenReq) validateRequest(c *gin.Context) (err error) {
	if err = c.BindJSON(&req); err != nil {
		err = errors.Wrap(err, errors.BadRequest)
		return
	}

	if len(req.Scopes) == 0 {
		err = errors.Wrap(err, errors.BadRequest)
		return
	}
	// 必须是允许范围内的scope
	for _, scope := range req.Scopes {
		if _, ok := types.Scopes[scope]; !ok {
			err = errors.New(errors.BadRequest)
			return
		}
	}
	req.AreaIDStr = c.GetHeader("Area-ID")
	if req.AreaIDStr != "" {
		req.areaID, _ = strconv.ParseUint(req.AreaIDStr, 10, 64)
	}
	return
}

// 根据用户选择，使用用户的token作为生成 JWT
func scopeToken(c *gin.Context) {
	var (
		req       scopeTokenReq
		resp      scopeTokenResp
		err       error
		tokenInfo oauth2.TokenInfo
	)

	defer func() {
		response.HandleResponse(c, err, resp)
	}()

	if err = req.validateRequest(c); err != nil {
		return
	}

	// 当前只有两个授权对象:SC和网盘，有临时密码的则是授权给SC
	code := c.GetHeader(types.VerificationKey)
	if code != "" {
		_, err = cache.Get(code)
		if err != nil {
			return
		}
		var scClient entity.Client
		scClient, err = entity.GetSCClient(req.areaID)
		if err != nil {
			return
		}

		tgr := &oauth2.TokenGenerateRequest{
			Request:      c.Request,
			ClientID:     scClient.ClientID,
			Scope:        scClient.AllowScope,
			ClientSecret: scClient.ClientSecret,
		}

		tokenInfo, err = oauth.GetOauthServer().GetAccessToken(oauth2.ClientCredentials, tgr)
		if err != nil {
			logger.Errorf("get oauth2 token error %s", err.Error())
			err = errors.Wrap(err, errors.BadRequest)
			return
		}
	} else { // 授权给网盘
		sessionUser := session.Get(c)
		if sessionUser == nil {
			err = errors.New(status.InvalidUserCredentials)
			return
		}
		var saClient entity.Client
		saClient, err = entity.GetSAClient(sessionUser.AreaID)
		if err != nil {
			return
		}

		tgr := &server.AuthorizeRequest{
			Request: c.Request,
		}
		tgr.ResponseType = oauth2.Token
		tgr.ClientID = saClient.ClientID
		tgr.UserID = strconv.Itoa(sessionUser.UserID)
		tgr.AccessTokenExp = ExpiresIn
		tgr.Scope = strings.Join(req.Scopes, ",")
		// TODO 使用oauth2生成scope_token，后续需要与前端联调去除
		tokenInfo, err = oauth.GetOauthServer().GetAuthorizeToken(tgr)
		if err != nil {
			logger.Errorf("get oauth2 token error %s", err.Error())
			err = errors.Wrap(err, errors.BadRequest)
			return
		}
	}

	resp.ScopeToken.Token = tokenInfo.GetAccess()
	resp.ScopeToken.ExpiresIn = int(tokenInfo.GetAccessExpiresIn() / time.Second)

	if code != "" {
		// 验证成功后删除code
		cache.Delete(code)
	}
}
