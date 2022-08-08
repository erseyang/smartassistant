package entity

import (
	"github.com/twinj/uuid"
	"github.com/zhiting-tech/smartassistant/modules/types"
	"github.com/zhiting-tech/smartassistant/pkg/errors"
	"github.com/zhiting-tech/smartassistant/pkg/rand"
	"gopkg.in/oauth2.v3"
	"gorm.io/gorm"
)

const (
	// AreaClient 默认Client
	AreaClient = iota + 1
	// SCClient SC
	SCClient
)

type Client struct {
	ID           int
	AreaID       uint64
	ClientID     string
	ClientSecret string
	GrantType    string
	AllowScope   string // 允许客户端申请的权限
	Type         int    // 0普通Client 1初始默认Client 2 SC Client
}

func (c Client) TableName() string {
	return "clients"
}

// CreateClient 创建应用
func CreateClient(grantType, allowScope string, areaID uint64) (client Client, err error) {
	client = Client{
		GrantType:  getAllowGrantType(grantType),
		AllowScope: allowScope,
		AreaID:     areaID,
	}

	if err = GetDB().Create(&client).Error; err != nil {
		err = errors.Wrap(err, errors.InternalServerErr)
		return
	}
	return
}

// CreateTypeClient 创建预定义应用
func CreateTypeClient(grantType, allowScope string, clientType int, areaID uint64) (client Client, err error) {
	client = Client{
		GrantType:  getAllowGrantType(grantType),
		AllowScope: allowScope,
		Type:       clientType,
		AreaID:     areaID,
	}

	if err = GetDB().Model(&client).Where(client).FirstOrCreate(&client).Error; err != nil {
		err = errors.Wrap(err, errors.InternalServerErr)
		return
	}

	return
}

// GetClientByClientID 根据ClientID获取Client信息
func GetClientByClientID(clientID string) (client Client, err error) {
	if err = GetDB().Where("client_id=?", clientID).First(&client).Error; err != nil {
		return
	}
	return
}

// getAllowGrantType 获取Client允许的授权类型
func getAllowGrantType(grantType string) string {
	if grantType != string(oauth2.Implicit) || grantType != string(oauth2.ClientCredentials) {
		grantType = grantType + "," + string(oauth2.Refreshing)
	}
	return grantType
}

// InitClient 初始化Client
func InitClient(areaID uint64) (err error) {
	var clients = make([]Client, 0)
	saClient := Client{
		GrantType:  string(oauth2.Implicit) + "," + string(oauth2.PasswordCredentials) + "," + string(oauth2.Refreshing),
		AllowScope: types.WithScopes(types.ScopeAll),
		AreaID:     areaID,
		Type:       AreaClient,
	}

	scClient := Client{
		GrantType:  string(oauth2.ClientCredentials),
		Type:       SCClient,
		AreaID:     areaID,
		AllowScope: types.WithScopes(types.ScopeGetTokenBySC, types.ScopeDevice),
	}

	clients = append(clients, saClient, scClient)

	for _, client := range clients {
		if err = GetDB().Create(&client).Error; err != nil {
			err = errors.Wrap(err, errors.InternalServerErr)
			return
		}
	}
	return
}

// GetSAClient 获取SAClient
func GetSAClient(areaID uint64) (client Client, err error) {
	if err = GetDB().First(&client, "type=? and area_id=?", AreaClient, areaID).Error; err != nil {
		err = errors.Wrap(err, errors.InternalServerErr)
		return
	}
	return
}

// GetSCClient 获取SCClient
func GetSCClient(areaID uint64) (client Client, err error) {
	if err = GetDB().First(&client, "type=? and area_id=?", SCClient, areaID).Error; err != nil {
		err = errors.Wrap(err, errors.InternalServerErr)
		return
	}
	return
}

func GetClientByAreaID(areaID uint64) (client []Client, err error) {
	if err = GetDB().Where("area_id=?", areaID).Find(&client).Error; err != nil {
		err = errors.Wrap(err, errors.InternalServerErr)
		return
	}
	return
}

func (c *Client) BeforeCreate(tx *gorm.DB) error {
	c.ClientID = uuid.NewV4().String()
	c.ClientSecret = rand.StringK(32, rand.KindAll)
	return nil
}
