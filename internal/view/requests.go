package view

import (
	"errors"
	"fmt"

	"github.com/caos/logging"
	"github.com/jinzhu/gorm"

	caos_errs "github.com/caos/zitadel/internal/errors"
)

func PrepareGetByKey(table string, key ColumnKey, id string) func(db *gorm.DB, res interface{}) error {
	return func(db *gorm.DB, res interface{}) error {
		err := db.Table(table).
			Where(fmt.Sprintf("%s = ?", key.ToColumnName()), id).
			Take(res).
			Error
		if err == nil {
			return nil
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return caos_errs.ThrowNotFound(err, "VIEW-XRI9c", "object not found")
		}
		logging.LogWithFields("VIEW-xVShS", "AggregateID", id).WithError(err).Warn("get from view error")
		return caos_errs.ThrowInternal(err, "VIEW-J92Td", "view error")
	}
}

func PrepareGetByQuery(table string, queries ...SearchQuery) func(db *gorm.DB, res interface{}) error {
	return func(db *gorm.DB, res interface{}) error {
		query := db.Table(table)
		for _, q := range queries {
			var err error
			query, err = SetQuery(query, q.GetKey(), q.GetValue(), q.GetMethod())
			if err != nil {
				return caos_errs.ThrowInvalidArgument(err, "VIEW-KaGue", "query is invalid")
			}
		}

		err := query.Take(res).Error
		if err == nil {
			return nil
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return caos_errs.ThrowNotFound(err, "VIEW-hodc6", "object not found")
		}
		logging.LogWithFields("VIEW-Mg6la", "table ", table).WithError(err).Warn("get from cache error")
		return caos_errs.ThrowInternal(err, "VIEW-qJBg9", "cache error")
	}
}

func PrepareSave(table string) func(db *gorm.DB, object interface{}) error {
	return func(db *gorm.DB, object interface{}) error {
		err := db.Table(table).Save(object).Error
		if err != nil {
			return caos_errs.ThrowInternal(err, "VIEW-AfC7G", "unable to put object to view")
		}
		return nil
	}
}

func PrepareDeleteByKey(table string, key ColumnKey, id string) func(db *gorm.DB) error {
	return func(db *gorm.DB) error {
		err := db.Table(table).
			Where(fmt.Sprintf("%s = ?", key.ToColumnName()), id).
			Delete(nil).
			Error
		if err != nil {
			return caos_errs.ThrowInternal(err, "VIEW-die73", "could not delete object")
		}
		return nil
	}
}

type Key struct {
	Key   ColumnKey
	Value string
}

func PrepareDeleteByKeys(table string, keys ...Key) func(db *gorm.DB) error {
	return func(db *gorm.DB) error {
		for _, key := range keys {
			db = db.Table(table).
				Where(fmt.Sprintf("%s = ?", key.Key.ToColumnName()), key.Value)
		}
		err := db.
			Delete(nil).
			Error
		if err != nil {
			return caos_errs.ThrowInternal(err, "VIEW-die73", "could not delete object")
		}
		return nil
	}
}

func PrepareDeleteByObject(table string, object interface{}) func(db *gorm.DB) error {
	return func(db *gorm.DB) error {
		err := db.Table(table).
			Delete(object).
			Error
		if err != nil {
			return caos_errs.ThrowInternal(err, "VIEW-lso9w", "could not delete object")
		}
		return nil
	}
}