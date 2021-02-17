package services

import (
	"strconv"
	"kf_server/models"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

// AdminRepositoryInterface interface
type AdminRepositoryInterface interface {
	GetAdmin(id int64) *models.Admin
	GetAdminWithUserName(userName string) *models.Admin
	Update(id int64, params *orm.Params) (int64, error)
	Add(admin *models.Admin, col1 string) (bool, int64, error)
	Delete(id int64) (int64, error)
	GetAdmins(request *AdminPaginationDto) (*AdminPaginationDto, error)
	GetOnlineAdmins() []models.Admin
	CheckAdminsLoginTimeOutAndSetOffline(lastMessageUnixTimer int64) int64
}

// AdminPaginationDto  a struct
type AdminPaginationDto struct {
	PageSize int         `json:"page_size"`
	PageOn   int         `json:"page_on"`
	Keyword  string      `json:"keyword"`
	Total    int64       `json:"total"`
	Online   int         `json:"online"`
	List     interface{} `json:"list"`
}

// UpdatePasswordRequest admin password
type UpdatePasswordRequest struct {
	OldPassword   string `json:"old_password"`
	NewPassword   string `json:"new_password"`
	EnterPassword string `json:"enter_password"`
}

// AdminRepository struct
type AdminRepository struct {
	BaseRepository
}

// GetAdminRepositoryInstance get instance
func GetAdminRepositoryInstance() *AdminRepository {
	instance := new(AdminRepository)
	instance.Init(new(models.Admin))
	return instance
}

// CheckAdminsLoginTimeOutAndSetOffline  Check if user login timeout
func (r *AdminRepository) CheckAdminsLoginTimeOutAndSetOffline(lastMessageUnixTimer int64) int64 {
	count, _ := r.q.Filter("online__in", 1, 2).Filter("last_activity__lte", lastMessageUnixTimer).Update(orm.Params{
		"online":           0,
		"current_con_user": 0,
	})
	return count
}

// GetAdmin get one admin with id
func (r *AdminRepository) GetAdmin(id int64) *models.Admin {
	var admin models.Admin
	if err := r.q.Filter("id", id).One(&admin); err != nil {
		logs.Warn("GetAdmin get one admin with id------------" + strconv.FormatInt(id,10)  , err)
		return nil
	}
	return &admin
}

// GetAdminWithUserName get one admin with username
func (r *AdminRepository) GetAdminWithUserName(userName string) *models.Admin {
	var admin models.Admin
	if err := r.q.Filter("UserName", userName).One(&admin); err != nil {
		logs.Warn("GetAdminWithUserName get one admin with username------------", err)
		return nil
	}
	return &admin
}

// Update admin
func (r *AdminRepository) Update(id int64, params orm.Params) (int64, error) {
	index, err := r.q.Filter("id", id).Update(params)
	if err != nil {
		logs.Warn("Update admin------------", err)
	}
	return index, err
}

// GetOnlineAdmins online all admin
func (r *AdminRepository) GetOnlineAdmins() []models.Admin {
	var admins []models.Admin
	_, err := r.o.Raw("SELECT a.*, IFNULL(c.count,0) AS `count` FROM admin as a LEFT  JOIN (SELECT to_account,COUNT(*) AS count FROM `contact` WHERE is_session_end = 0 GROUP BY to_account) c ON a.id = c.to_account WHERE a.`online` = 1 ORDER BY c.count").QueryRows(&admins)
	if err != nil {
		logs.Warn("GetOnlineAdmins online all admin------------", err)
		return []models.Admin{}
	}
	return admins
}

// Add create a admin
func (r *AdminRepository) Add(admin *models.Admin, col1 string) (bool, int64, error) {
	_bool, index, err := r.o.ReadOrCreate(admin, col1)
	if err != nil {
		logs.Warn("Add create a admin------------", err)
	}
	return _bool, index, err
}

// Delete delete a admin
func (r *AdminRepository) Delete(id int64) (int64, error) {
	index, err := r.q.Filter("id", id).Delete()
	if err != nil {
		logs.Warn("Delete delete a admin------------", err)
	}
	return index, err
}

// GetAdmins get admin list
func (r *AdminRepository) GetAdmins(request *AdminPaginationDto) (*AdminPaginationDto, error) {

	var lists []models.Admin

	qs := r.q.Filter("nickname__icontains", request.Keyword)
	if request.Online == 0 {
		qs = qs.Filter("online", 0)
	}
	if request.Online == 1 {
		qs = qs.Filter("online", 1)
	}
	if request.PageSize < MinPageSize {
		request.PageSize = MinPageSize
	}
	if request.PageSize > MaxPageSize {
		request.PageSize = MaxPageSize
	}
	if _, err := qs.OrderBy("-root", "id").Limit(request.PageSize, (request.PageOn-1)*request.PageSize).All(&lists); err != nil {
		logs.Warn("GetAdmins get admin list------------", err)
		return nil, err
	}
	total, _ := qs.Count()
	for index := range lists {
		lists[index].Password = "******"
	}
	request.Total = total
	request.List = &lists

	return request, nil
}
