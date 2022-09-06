package user

import (
	"strconv"
	"time"

	"github.com/zhiting-tech/smartassistant/modules/maintenance"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	"github.com/zhiting-tech/smartassistant/modules/api/utils/response"
	"github.com/zhiting-tech/smartassistant/modules/entity"
	"github.com/zhiting-tech/smartassistant/modules/types"
	"github.com/zhiting-tech/smartassistant/modules/types/status"
	"github.com/zhiting-tech/smartassistant/modules/utils/hash"
	"github.com/zhiting-tech/smartassistant/modules/utils/session"
	"github.com/zhiting-tech/smartassistant/pkg/rand"
	"github.com/zhiting-tech/smartassistant/pkg/wangpan"

	"github.com/zhiting-tech/smartassistant/pkg/errors"
)

// updateUserReq 修改用户接口请求参数
type updateUserReq struct {
	Nickname      *string `json:"nickname"`
	AccountName   *string `json:"account_name"`
	Password      *string `json:"password"`
	OldPassword   *string `json:"old_password"`
	AvatarID      *int    `json:"avatar_id"`
	RoleIds       []int   `json:"role_ids"`
	DepartmentIds []int   `json:"department_ids,omitempty"`
}

func (req *updateUserReq) Validate(updateUid, loginId int, areaType entity.AreaType, loginUserInfo entity.User) (updateUser entity.User, err error) {
	if len(req.RoleIds) != 0 {
		// 判断是否有修改角色权限
		if !entity.JudgePermit(loginId, types.AreaUpdateMemberRole) {
			err = errors.Wrap(err, status.Deny)
			return
		}
	}

	if entity.IsCompany(areaType) && req.DepartmentIds != nil {
		if !entity.JudgePermit(loginId, types.AreaUpdateMemberDepartment) {
			err = errors.Wrap(err, status.Deny)
			return
		}
	}

	// 自己才允许修改自己的用户名,密码和昵称
	if req.Nickname != nil || req.AccountName != nil || req.Password != nil || req.OldPassword != nil {
		if loginId != updateUid {
			err = errors.New(status.Deny)
			return
		}
	}

	if req.Nickname != nil {
		if err = checkNickname(*req.Nickname); err != nil {
			return
		} else {
			updateUser.Nickname = *req.Nickname
		}
	}
	if req.AccountName != nil {
		if err = checkAccountName(*req.AccountName); err != nil {
			return
		} else {
			updateUser.AccountName = *req.AccountName
		}
	}

	if req.AvatarID != nil {
		if _, err = entity.GetFileInfo(*req.AvatarID); err != nil {
			return
		}
		updateUser.AvatarID = *req.AvatarID
	}

	if req.Password != nil {
		if err = checkPassword(*req.Password); err != nil {
			return
		}

		// 通过密码是否设置过，判断是否是修改密码还是初始设置密码
		if loginUserInfo.Password == "" {
			salt := rand.String(rand.KindAll)
			updateUser.Salt = salt
			hashNewPassword := hash.GenerateHashedPassword(*req.Password, salt)
			updateUser.Password = hashNewPassword
		} else {
			rd := maintenance.DirectResetProPassword()
			if !rd {
				if loginUserInfo.Password != hash.GenerateHashedPassword(*req.OldPassword, loginUserInfo.Salt) {
					err = errors.New(status.OldPasswordErr)
					return
				}
			}
			updateUser.Password = hash.GenerateHashedPassword(*req.Password, loginUserInfo.Salt)
			updateUser.PasswordUpdateTime = time.Now()
		}
	}

	return
}

// UpdateUser 用于处理修改用户接口的请求
func UpdateUser(c *gin.Context) {
	var (
		err           error
		req           updateUserReq
		updateUser    entity.User
		sessionUser   *session.User
		userID        int
		curArea       entity.Area
		loginUserInfo entity.User
	)
	defer func() {
		response.HandleResponse(c, err, nil)
	}()

	if userID, err = strconv.Atoi(c.Param("id")); err != nil {
		err = errors.Wrap(err, errors.BadRequest)
		return
	}

	sessionUser = session.Get(c)
	if sessionUser == nil {
		err = errors.Wrap(err, status.AccountNotExistErr)
		return
	}

	err = c.BindJSON(&req)
	if err != nil {
		err = errors.Wrap(err, errors.BadRequest)
		return
	}

	if curArea, err = entity.GetAreaByID(sessionUser.AreaID); err != nil {
		return
	}

	if loginUserInfo, err = entity.GetUserByID(sessionUser.UserID); err != nil {
		return
	}

	if updateUser, err = req.Validate(userID, sessionUser.UserID, curArea.AreaType, loginUserInfo); err != nil {
		return
	}

	if len(req.RoleIds) != 0 {
		if entity.IsOwnerOfArea(userID, sessionUser.AreaID) {
			err = errors.New(status.NotAllowModifyRoleOfTheOwner)
			return
		}
		err = entity.GetDB().Transaction(func(tx *gorm.DB) error {
			// 删除用户原有角色
			err = tx.Unscoped().Where("user_id=?", userID).Delete(&entity.UserRole{}).Error
			if err != nil {
				return errors.Wrap(err, errors.InternalServerErr)
			}
			// 添加新的用户角色
			roles := wrapURoles(userID, req.RoleIds)
			if err = tx.Create(&roles).Error; err != nil {
				return errors.Wrap(err, errors.InternalServerErr)
			}
			return nil
		})
		if err != nil {
			return
		}
	}

	if entity.IsCompany(curArea.AreaType) && req.DepartmentIds != nil {
		if err = CheckDepartmentsManager(userID, req.DepartmentIds, curArea.ID); err != nil {
			return
		}
		if err = entity.CreateDepartmentUser(entity.WrapDepUsersOfUId(userID, req.DepartmentIds)); err != nil {
			return
		}
	}

	if err = req.updateSmbConf(c, loginUserInfo); err != nil {
		return
	}

	if err = entity.EditUser(userID, updateUser); err != nil {
		return
	}

	return
}

// getDelManagerDepartments 获取需要重置主管的部门
func getDelManagerDepartments(oldDepartment []entity.Department, newDepartmentIDs []int) (departmentIDs []int) {
	if len(newDepartmentIDs) == 0 {
		for _, department := range oldDepartment {
			departmentIDs = append(departmentIDs, department.ID)
		}
		return
	}

	for _, department := range oldDepartment {
		isOld := true
		for _, departmentID := range newDepartmentIDs {
			if departmentID == department.ID {
				isOld = false
				break
			}
		}
		if isOld {
			departmentIDs = append(departmentIDs, department.ID)
		}
	}
	return
}

// CheckDepartmentsManager 检查多个部门的主管是否被删除，并重置它
// TODO 尝试在删除的hook中做删除主管的操作，但要注意，原先的删除/更新逻辑的影响
func CheckDepartmentsManager(userID int, departmentIds []int, areaID uint64) (err error) {
	var departments []entity.Department
	if departments, err = entity.GetManagerDepartments(areaID, userID); err != nil {
		return
	}
	delManagerDepartmentIDs := getDelManagerDepartments(departments, departmentIds)
	if err = entity.ResetDepartmentManager(areaID, delManagerDepartmentIDs...); err != nil {
		return
	}

	if err = entity.UnScopedDelUserDepartments(userID); err != nil {
		return
	}
	return
}
func (req *updateUserReq) updateSmbConf(c *gin.Context, loginUserInfo entity.User) error {
	if req.AccountName == nil && req.Password == nil {
		return nil
	}
	if req.AccountName == nil {
		req.AccountName = new(string)
	}
	if req.Password == nil {
		req.Password = new(string)
	}
	smb := wangpan.NewSmbMountStr(loginUserInfo.AccountName, *req.AccountName, *req.Password)
	if err := smb.SetMountPath(""); err != nil {
		err = errors.Wrap(err, errors.InternalServerErr)
		return err
	}
	if err := smb.Exec(); err != nil {
		err = errors.Wrap(err, errors.InternalServerErr)
		return err
	}
	return nil
}
