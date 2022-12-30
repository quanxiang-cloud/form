package service

import (
	"context"
	id2 "github.com/quanxiang-cloud/cabin/id"
	"github.com/quanxiang-cloud/form/internal/models"
	"github.com/quanxiang-cloud/form/internal/models/mysql"
	config2 "github.com/quanxiang-cloud/form/pkg/misc/config"
	"gorm.io/gorm"
)

type Project interface {
	CreateProject(ctx context.Context, req *CreateProjectReq) (*CreateProjectResp, error)
	DeleteProject(ctx context.Context, req *DeleteProjectReq) (*DeleteProjectResp, error)
	ListProject(ctx context.Context, req *ListProjectReq) (*ListProjectResp, error)
	// AssignProjectUser  增加人
	AssignProjectUser(ctx context.Context, req *AssignProjectUserReq) (*AssignProjectUserResp, error)

	ListProjectUser(ctx context.Context, req *ListProjectUserReq) (*ListProjectUserResp, error)
}

type project struct {
	projectRepo    models.ProjectRepo
	projectUerRepo models.ProjectUserRepo
	db             *gorm.DB
}

func NewProject(conf *config2.Config) (Project, error) {
	db, err := CreateMysqlConn(conf)
	if err != nil {
		return nil, err
	}
	p := &project{
		projectRepo:    mysql.NewProjectRepo(),
		projectUerRepo: mysql.NewProjectUserRepo(),
		db:             db,
	}
	return p, nil
}

type AssignProjectUserReq struct {
	ProjectID   string   `json:"projectID"`
	ProjectName string   `json:"projectName"`
	Add         []*user  `json:"add"`
	Removes     []string `json:"removes"`
}

type user struct {
	UserName string `json:"userName"`
	UserID   string `json:"userID"`
}

type AssignProjectUserResp struct {
}

func (p *project) AssignProjectUser(ctx context.Context, req *AssignProjectUserReq) (*AssignProjectUserResp, error) {
	tx := p.db.Begin()
	projectUser := make([]*models.ProjectUser, len(req.Add))
	for index, value := range req.Add {
		projectUser[index] = &models.ProjectUser{
			ID:          id2.StringUUID(),
			ProjectID:   req.ProjectID,
			ProjectName: req.ProjectName,
			UserID:      value.UserID,
			UserName:    value.UserName,
		}
	}
	if len(projectUser) > 0 {
		err := p.projectUerRepo.BatchCreate(tx, projectUser...)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	if len(req.Removes) == 0 {
		tx.Commit()
		return &AssignProjectUserResp{}, nil
	}
	err := p.projectUerRepo.Delete(tx, &models.ProjectUserQuery{
		ProjectID: req.ProjectID,
		UserIDs:   req.Removes,
	})
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return &AssignProjectUserResp{}, nil
}

type ListProjectUserReq struct {
	Page      int    `json:"page"`
	Size      int    `json:"size"`
	ProjectID string `json:"projectID"`
	UserID    string `json:"userID "`
}

type ListProjectUserResp struct {
	List  []*projectUserVO `json:"list"`
	Total int64            `json:"total"`
}

type projectUserVO struct {
	ID          string `json:"id"`
	ProjectID   string `json:"projectID"`
	ProjectName string `json:"projectName"`
	UserID      string `json:"userID"`
	UserName    string `json:"userName"`
}

func (p *project) ListProjectUser(ctx context.Context, req *ListProjectUserReq) (*ListProjectUserResp, error) {
	list, total, err := p.projectUerRepo.List(p.db, &models.ProjectUserQuery{
		ProjectID: req.ProjectID,
	}, req.Page, req.Size)
	if err != nil {
		return nil, err
	}
	resp := &ListProjectUserResp{
		Total: total,
		List:  make([]*projectUserVO, len(list)),
	}
	for index, value := range list {
		resp.List[index] = &projectUserVO{
			ID:          value.ID,
			ProjectID:   value.ProjectID,
			ProjectName: value.ProjectName,
			UserID:      value.UserID,
			UserName:    value.UserName,
		}
	}
	return resp, nil
}

type CreateProjectReq struct {
	CreatedAt   int64
	CreatorID   string
	CreatorName string

	Name        string `json:"name"`
	Description string `json:"description"`
}

type CreateProjectResp struct {
	ID string `json:"id"`
}

func (p *project) CreateProject(ctx context.Context, req *CreateProjectReq) (*CreateProjectResp, error) {
	project1 := &models.Project{
		ID:          id2.StringUUID(),
		CreatedAt:   req.CreatedAt,
		CreatorID:   req.CreatorID,
		CreatorName: req.CreatorName,
		Name:        req.Name,
		Description: req.Description,
	}
	err := p.projectRepo.BatchCreate(p.db, project1)
	if err != nil {
		return nil, err
	}
	resp := &CreateProjectResp{
		ID: project1.ID,
	}
	return resp, nil
}

type DeleteProjectReq struct {
	ID string `json:"id"`
}

type DeleteProjectResp struct {
}

// DeleteProject 根据id 删除项目
func (p *project) DeleteProject(ctx context.Context, req *DeleteProjectReq) (*DeleteProjectResp, error) {
	project1 := &models.ProjectQuery{
		ID: req.ID,
	}
	err := p.projectRepo.Delete(p.db, project1)
	if err != nil {
		return nil, err
	}
	resp := &DeleteProjectResp{}
	return resp, nil
}

type ListProjectReq struct {
	Page int `json:"page"`
	Size int `json:"size"`
}

type ListProjectResp struct {
	List  []*projectVo `json:"list"`
	Total int64        `json:"total"`
}
type projectVo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// ListProject 分页查询项目
func (p *project) ListProject(ctx context.Context, req *ListProjectReq) (*ListProjectResp, error) {
	list, total, err := p.projectRepo.List(p.db, &models.ProjectQuery{}, req.Page, req.Size)
	if err != nil {
		return nil, err
	}
	resp := &ListProjectResp{
		Total: total,
		List:  make([]*projectVo, len(list)),
	}
	for index, value := range list {
		resp.List[index] = &projectVo{
			ID:          value.ID,
			Name:        value.Name,
			Description: value.Description,
		}
	}
	return resp, nil
}
