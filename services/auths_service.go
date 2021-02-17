package services

import (
	"kf_server/models"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

// AuthsRepositoryInterface interface
type AuthsRepositoryInterface interface {
	GetAdminAuthInfo(token string) *models.Auths
	GetAdminAuthInfoWithTypeAndUID(authType int64, uid int64) *models.Auths
	GetAdminOnlineCount(uid int64) int64
	Delete(id int64) (int64, error)
	Add(id *models.Auths) (int64, error)
	Update(id int64, params orm.Params) (int64, error)
}

// AuthsRepository struct
type AuthsRepository struct {
	BaseRepository
}

// GetAuthsRepositoryInstance get instance
func GetAuthsRepositoryInstance() *AuthsRepository {
	instance := new(AuthsRepository)
	instance.Init(new(models.Auths))
	return instance
}

// GetAdminAuthInfo get a auth info
func (r *AuthsRepository) GetAdminAuthInfo(token string) *models.Auths {
	var auth models.Auths
	if err := r.q.Filter("token", token).One(&auth); err != nil {
		logs.Warn("GetAdminAuthInfo get a auth info------------", err)
		return nil
	}
	return &auth
}

// GetAdminAuthInfoWithTypeAndUID get a auth info with type and uid
func (r *AuthsRepository) GetAdminAuthInfoWithTypeAndUID(authType int64, uid int64) *models.Auths {
	var auth models.Auths
	if err := r.q.Filter("auth_type", authType).Filter("uid", uid).One(&auth); err != nil {
		logs.Warn("GetAdminAuthInfoWithTypeAndUID get a auth info with type and uid------------", err)
		return nil
	}
	return &auth
}

// GetAdminOnlineCount get admin anth count
func (r *AuthsRepository) GetAdminOnlineCount(uid int64) int64 {
	count, err := r.q.Filter("uid", uid).Count()
	if err != nil {
		logs.Warn(" GetAdminOnlineCount get admin anth count------------", err)
		return 0
	}
	return count
}

// Delete delete a auth
func (r *AuthsRepository) Delete(id int64) (int64, error) {
	row, err := r.q.Filter("id", id).Delete()
	if err != nil {
		logs.Warn("Delete delete a auth------------", err)
		return 0, err
	}
	return row, nil
}

// Update admin
func (r *AuthsRepository) Update(id int64, params orm.Params) (int64, error) {
	index, err := r.q.Filter("id", id).Update(params)
	if err != nil {
		logs.Warn("Update admin------------", err)
	}
	return index, err
}

// Add add auth
func (r *AuthsRepository) Add(auth *models.Auths) (int64, error) {
	index, err := r.o.Insert(auth)
	if err != nil {
		logs.Warn("Add add auth------------", err)
	}
	return index, err
}
