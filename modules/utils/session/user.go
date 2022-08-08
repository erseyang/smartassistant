package session

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/zhiting-tech/smartassistant/modules/api/utils/oauth"
	"github.com/zhiting-tech/smartassistant/modules/entity"
	"github.com/zhiting-tech/smartassistant/modules/types"
	"github.com/zhiting-tech/smartassistant/pkg/logger"
)

const sessionName = "user"

// User 表示请求的用户，可能是普通用户，也可能是Client
type User struct {
	UserID        int                    `json:"uid"`
	IsOwner       bool                   `json:"is_owner"`
	UserName      string                 `json:"user_name"`
	RoleID        int                    `json:"role_id"`
	Token         string                 `json:"token"`
	LoginAt       time.Time              `json:"login_at"`
	LoginDuration time.Duration          `json:"login_duration"`
	ExpiresAt     time.Time              `json:"expires_at"`
	AreaID        uint64                 `json:"area_id"`
	Option        map[string]interface{} `json:"option"`
	Key           string
	ClientID      string
	Scope         string
}

func (u User) BelongsToArea(areaID uint64) bool {
	return u.AreaID == areaID
}

// GetPermissions 获取用户的所有权限
func (u User) GetPermissions() (up entity.Permissions, err error) {
	if u.IsClient() {
		return entity.GetClientPermissions(u.ClientID, u.Scope)
	}
	return entity.GetUserPermissions(u.UserID)
}

func (u User) IsClient() bool {
	return u.ClientID != "" && u.UserID == 0
}

// Get 根据token或cookie获取用户数据
func Get(c *gin.Context) *User {
	if u, exists := c.Get("userInfo"); exists {
		return u.(*User)
	}
	var u *User
	u = getUserByToken(c)
	c.Set("userInfo", u)
	return u
}

func getUserByToken(c *gin.Context) *User {
	accessToken := c.GetHeader(types.SATokenKey)
	if accessToken == "" {
		return nil
	}
	ti, err := oauth.GetOauthServer().Manager.LoadAccessToken(accessToken)
	if err != nil {
		logger.Errorf("load access token err: %s", err)
		return nil
	}

	clientID := ti.GetClientID()
	cli, err := entity.GetClientByClientID(clientID)
	if err != nil {
		logger.Errorf("get client err: %s", err)
		return nil
	}
	uid, _ := strconv.Atoi(ti.GetUserID())
	user, _ := entity.GetUserByID(uid)

	area, err := entity.GetAreaByID(cli.AreaID)
	if err != nil {
		logger.Errorf("GetAreaByID err: %s", err)
		return nil
	}

	u := &User{
		UserID:   uid,
		UserName: user.AccountName,
		Token:    accessToken,
		AreaID:   cli.AreaID,
		IsOwner:  area.OwnerID == user.ID,
		Key:      user.Key,
		ClientID: clientID,
		Scope:    ti.GetScope(),
	}
	return u
}
