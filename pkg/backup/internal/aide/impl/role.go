package impl

import (
	"context"
	"fmt"

	id2 "github.com/quanxiang-cloud/cabin/id"
	time2 "github.com/quanxiang-cloud/cabin/time"
	"github.com/quanxiang-cloud/form/internal/models"
	"github.com/quanxiang-cloud/form/pkg/backup/internal/aide"
)

var (
	exportRoleURL = "%s/api/v1/form/%s/internal/backup/export/role"
	importRoleURL = "%s/api/v1/form/%s/internal/backup/import/role"

	exportPermitURL = "%s/api/v1/form/%s/internal/backup/export/permit"
	importPermitURL = "%s/api/v1/form/%s/internal/backup/import/permit"
)

// Role role.
type Role struct{}

func (r *Role) roleTag() string {
	return "roles"
}

func (r *Role) permitTag() string {
	return "permits"
}

// Export export.
// NOTE: roles are related to permissions and are handled in roles.
func (r *Role) Export(ctx context.Context, opts *aide.ExportOption) (map[string]aide.Object, error) {
	roleObj, err := r.exportRole(ctx, opts)
	if err != nil {
		return nil, err
	}

	permitObj, err := r.exportPermit(ctx, opts)
	if err != nil {
		return nil, err
	}

	return map[string]aide.Object{
		r.roleTag():   roleObj,
		r.permitTag(): permitObj,
	}, nil
}

func (r *Role) exportRole(ctx context.Context, opts *aide.ExportOption) (aide.Object, error) {
	url := fmt.Sprintf(exportRoleURL, opts.Host, opts.AppID)

	obj, err := aide.ExportObject(ctx, url, opts)
	if err != nil {
		return nil, err
	}

	return obj, nil
}

func (r *Role) exportPermit(ctx context.Context, opts *aide.ExportOption) (aide.Object, error) {
	url := fmt.Sprintf(exportPermitURL, opts.Host, opts.AppID)

	obj, err := aide.ExportObject(ctx, url, opts)
	if err != nil {
		return nil, err
	}

	return obj, nil
}

// Import import.
func (r *Role) Import(ctx context.Context, objs map[string]aide.Object, opts *aide.ImportOption) (map[string]string, error) {
	ids := make(map[string]string)

	roleObj := objs[r.roleTag()]

	roleIDs, err := r.importRole(ctx, roleObj, opts)
	if err != nil {
		return nil, err
	}

	for k, v := range roleIDs {
		ids[k] = v
	}

	permitObj := objs[r.permitTag()]

	permitIDs, err := r.importPermit(ctx, permitObj, roleIDs, opts)
	if err != nil {
		return nil, err
	}

	for k, v := range permitIDs {
		ids[k] = v
	}

	return ids, nil
}

func (r *Role) importRole(ctx context.Context, obj aide.Object, opts *aide.ImportOption) (map[string]string, error) {
	var roles []*models.Role
	err := aide.Serialize(obj, &roles)
	if err != nil {
		return nil, err
	}

	ids := r.replaceRoleParam(roles, opts)

	data := make(aide.Object, len(obj))
	for i := 0; i < len(obj); i++ {
		data[i] = roles[i]
	}

	url := fmt.Sprintf(importRoleURL, opts.Host, opts.AppID)

	err = aide.ImportObject(ctx, url, data, opts)
	if err != nil {
		return nil, err
	}

	return ids, nil
}

func (r *Role) replaceRoleParam(roles []*models.Role, opts *aide.ImportOption) map[string]string {
	ids := make(map[string]string)

	for i := 0; i < len(roles); i++ {
		id := id2.StringUUID()
		ids[roles[i].ID] = id

		roles[i].ID = id
		roles[i].AppID = opts.AppID
		roles[i].CreatorID = opts.UserID
		roles[i].CreatorName = opts.UserName
		roles[i].CreatedAt = time2.NowUnix()
	}

	return ids
}

func (r *Role) importPermit(ctx context.Context, obj aide.Object, roleIDs map[string]string, opts *aide.ImportOption) (map[string]string, error) {
	var permits []*models.Permit
	err := aide.Serialize(obj, &permits)
	if err != nil {
		return nil, err
	}

	ids := r.replacePermitParam(permits, roleIDs, opts)

	data := make(aide.Object, len(obj))
	for i := 0; i < len(obj); i++ {
		data[i] = permits[i]
	}

	url := fmt.Sprintf(importPermitURL, opts.Host, opts.AppID)

	err = aide.ImportObject(ctx, url, data, opts)
	if err != nil {
		return nil, err
	}

	return ids, nil
}

func (r *Role) replacePermitParam(permits []*models.Permit, roleIDs map[string]string, opts *aide.ImportOption) map[string]string {
	ids := make(map[string]string)

	for i := 0; i < len(permits); i++ {
		id := id2.StringUUID()
		ids[permits[i].ID] = id

		permits[i].RoleID = roleIDs[permits[i].RoleID]
		permits[i].ID = id
		permits[i].CreatorID = opts.UserID
		permits[i].CreatorName = opts.UserName
		permits[i].CreatedAt = time2.NowUnix()
	}

	return ids
}
