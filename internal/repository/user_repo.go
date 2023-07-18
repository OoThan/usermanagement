package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/OoThan/usermanagement/internal/model"
	"github.com/OoThan/usermanagement/pkg/dto"
	"github.com/OoThan/usermanagement/pkg/utils"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type userRepository struct {
	db  *gorm.DB
	rdb *redis.Client
	mdb *mongo.Client
}

func newUserRespository(rConfig *RepoConfig) *userRepository {
	return &userRepository{
		db:  rConfig.DS.DB,
		rdb: rConfig.DS.RDB,
		mdb: rConfig.DS.MDB,
	}
}

func (r *userRepository) FindByField(ctx context.Context, field, value any) (*model.User, error) {
	db := r.db.WithContext(ctx).Debug().Model(&model.User{})
	user := &model.User{}
	err := db.First(&user, fmt.Sprintf("BINARY %s = ?", field), value).Error
	return user, err
}

func (r *userRepository) FindOrByField(ctx context.Context, field1, field2, value any) (*model.User, error) {
	db := r.db.WithContext(ctx).Debug().Model(&model.User{})
	user := &model.User{}
	err := db.First(&user, fmt.Sprintf("%s = ? OR %s = ?", field1, field2), value, value).Error
	return user, err
}

func (r *userRepository) UpdateByFields(ctx context.Context, updateFields *model.UpdateFields) (*model.User, error) {
	db := r.db.WithContext(ctx).Debug().Model(&model.User{})
	db.Where(updateFields.Field, updateFields.Value)
	err := db.Updates(updateFields.Data).Error
	return nil, err
}

func (r *userRepository) List(ctx context.Context, req *dto.UserListReq) ([]*model.User, int64, error) {
	list := make([]*model.User, 0)
	db := r.db.WithContext(ctx).Debug().Model(&model.User{})

	if req.ID != 0 {
		db.Where("id", req.ID)
	}

	if req.Name != "" {
		db.Where("name LIKE ?", "%"+req.Name+"%")
	}

	if req.Email != "" {
		db.Where("email LIKE ?", "%"+req.Email+"%")
	}

	var total int64
	db.Count(&total)
	if err := db.Scopes(utils.Paginate(req.Page, req.PageSize)).Find(&list).Error; err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	db := r.db.WithContext(ctx).Debug().Model(&model.User{})
	var count int64
	oldUser := &model.User{}
	db.Unscoped().Where("name = ? AND deleted_at IS NOT NULL", user.Name).Count(&count).First(&oldUser)
	if count > 0 {
		user.Id = oldUser.Id
		user.CreatedAt = time.Now()
		user.UpdatedAt = time.Now()
		user.DeletedAt = gorm.DeletedAt{
			Time:  time.Time{},
			Valid: false,
		}
		return db.Save(&user).Error
	}
	tb := r.db.WithContext(ctx).Debug().Model(&model.User{})
	return tb.Create(&user).Error
}

func (r *userRepository) Update(ctx context.Context, updateFields *model.UpdateFields) error {
	db := r.db.WithContext(ctx).Debug().Model(&model.User{})
	db.Where(updateFields.Field, updateFields.Value)
	return db.Updates(&updateFields.Data).Error
}

func (r *userRepository) Delete(ctx context.Context, ids string) error {
	db := r.db.WithContext(ctx).Debug().Model(&model.User{})
	db.Where(fmt.Sprintf("id in (%s)", ids))
	return db.Delete(&model.User{}).Error
}
