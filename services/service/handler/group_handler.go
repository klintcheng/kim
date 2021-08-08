package handler

import (
	"errors"

	"github.com/kataras/iris/v12"
	"github.com/klintcheng/kim/services/service/database"
	"github.com/klintcheng/kim/wire/rpc"
	"gorm.io/gorm"
)

// var log = logger.WithField("module", "service.handler")

func (h *ServiceHandler) GroupCreate(c iris.Context) {
	app := c.Params().Get("app")
	var req rpc.CreateGroupReq
	if err := c.ReadBody(&req); err != nil {
		c.StopWithError(iris.StatusBadRequest, err)
		return
	}
	groupId := h.Idgen.Next()
	g := &database.Group{
		Model: database.Model{
			ID: groupId.Int64(),
		},
		App:          app,
		Group:        groupId.Base32(),
		Name:         req.Name,
		Avatar:       req.Avatar,
		Owner:        req.Owner,
		Introduction: req.Introduction,
	}
	members := make([]database.GroupMember, len(req.Members))
	for i, user := range req.Members {
		members[i] = database.GroupMember{
			Model: database.Model{
				ID: h.Idgen.Next().Int64(),
			},
			Account: user,
			Group:   groupId.Base32(),
		}
	}

	err := h.BaseDb.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(g).Error; err != nil {
			// return anywill rollback
			return err
		}
		if err := tx.Create(&members).Error; err != nil {
			return err
		}
		// return nil will commit the whole transaction
		return nil
	})
	if err != nil {
		c.StopWithError(iris.StatusInternalServerError, err)
		return
	}
	_, _ = c.Negotiate(&rpc.CreateGroupResp{
		GroupId: groupId.Base32(),
	})
}

func (h *ServiceHandler) GroupJoin(c iris.Context) {
	// app := c.Param("app")
	var req rpc.JoinGroupReq
	if err := c.ReadBody(&req); err != nil {
		c.StopWithError(iris.StatusBadRequest, err)
		return
	}
	gm := &database.GroupMember{
		Model: database.Model{
			ID: h.Idgen.Next().Int64(),
		},
		Account: req.Account,
		Group:   req.GroupId,
	}
	err := h.BaseDb.Create(gm).Error
	if err != nil {
		c.StopWithError(iris.StatusInternalServerError, err)
		return
	}
}

func (h *ServiceHandler) GroupQuit(c iris.Context) {
	// app := c.Param("app")
	var req rpc.QuitGroupReq
	if err := c.ReadBody(&req); err != nil {
		c.StopWithError(iris.StatusBadRequest, err)
		return
	}
	gm := &database.GroupMember{
		Account: req.Account,
		Group:   req.GroupId,
	}
	err := h.BaseDb.Delete(&database.GroupMember{}, gm).Error
	if err != nil {
		c.StopWithError(iris.StatusInternalServerError, err)
		return
	}
}

func (h *ServiceHandler) GroupMembers(c iris.Context) {
	group := c.Params().Get("group")
	if group == "" {
		c.StopWithError(iris.StatusBadRequest, errors.New("group is null"))
		return
	}
	var members []database.GroupMember
	err := h.BaseDb.Order("Updated_At asc").Find(&members, database.GroupMember{Group: group}).Error
	if err != nil {
		c.StopWithError(iris.StatusInternalServerError, err)
		return
	}
	var users = make([]*rpc.Member, len(members))
	for i, m := range members {
		users[i] = &rpc.Member{
			Account:  m.Account,
			Alias:    m.Alias,
			JoinTime: m.CreatedAt.Unix(),
		}
	}
	_, _ = c.Negotiate(&rpc.GroupMembersResp{
		Users: users,
	})
}
