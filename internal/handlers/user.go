package handlers

import (
	"bytes"
	"cippus-backend/internal/services"
	"io"
	"log/slog"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service *services.UserService
	minio   *services.MinioService
}

func NewUserHandler(service *services.UserService, minio *services.MinioService) *UserHandler {
	return &UserHandler{service: service, minio: minio}
}

func (h *UserHandler) PatchMe(ctx *gin.Context) {
	userID := uint(ctx.GetFloat64("userID"))
	if userID == 0 {
		ctx.JSON(400, gin.H{"error": "Invalid Request"})
		return
	}

	req := &PatchMeRequest{}
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid Request"})
		return
	}

	userMap := map[string]interface{}{}
	if req.Username != nil {
		userMap["username"] = *req.Username
	}
	if req.Bio != nil {
		userMap["bio"] = *req.Bio
	}

	if req.NotifReply != nil {
		userMap["notif_reply"] = *req.NotifReply
	}
	if req.NotifFollows != nil {
		userMap["notif_follow"] = *req.NotifFollows
	}
	if req.NotifLike != nil {
		userMap["notif_like"] = *req.NotifLike
	}
	if req.NotifMention != nil {
		userMap["notif_mention"] = *req.NotifMention
	}
	if req.NotifDM != nil {
		userMap["notif_dm"] = *req.NotifDM
	}
	if req.NotifPostApproved != nil {
		userMap["notif_post_approved"] = *req.NotifPostApproved
	}

	if len(userMap) == 0 {
		ctx.JSON(400, gin.H{"error": "Invalid Request"})
		return
	}

	err = h.service.PatchMe(userID, userMap)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Cannot update user"})
		return
	}

	ctx.JSON(200, gin.H{})
}

func (h *UserHandler) GetMe(ctx *gin.Context) {
	userID := uint(ctx.GetFloat64("userID"))
	if userID == 0 {
		ctx.JSON(400, gin.H{"error": "user not found"})
		return
	}

	user, err := h.service.GetUserByID(userID)
	if err != nil {
		if err.Error() == "user not found" {
			ctx.JSON(404, gin.H{"error": "user does not exists"})
			return
		}
		ctx.JSON(500, gin.H{"error": "cannot find user"})
		return
	}

	ctx.JSON(200, gin.H{"user": toUserResponse(user)})
}

func (h *UserHandler) GetUser(ctx *gin.Context) {
	strID := ctx.Param("userId")
	id, err := strconv.ParseUint(strID, 10, 64)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "could not find user id"})
		return
	}

	user, err := h.service.GetUserByID(uint(id))
	if err != nil {
		if err.Error() == "user not found" {
			ctx.JSON(404, gin.H{"error": "user does not exists"})
			return
		}
		ctx.JSON(500, gin.H{"error": "cannot find user"})
		return
	}

	ctx.JSON(200, gin.H{"user": toUserResponse(user)})
}

func (h *UserHandler) UploadAvatar(ctx *gin.Context) {
	userID := uint(ctx.GetFloat64("userID"))
	if userID == 0 {
		ctx.JSON(400, gin.H{"error": "cannot find user id"})
		return
	}

	user, err := h.service.GetUserByID(uint(userID))
	if err != nil {
		if err.Error() == "user not found" {
			ctx.JSON(404, gin.H{"error": "user does not exists"})
			return
		}
		ctx.JSON(500, gin.H{"error": "cannot find user"})
		return
	}

	file, err := ctx.FormFile("avatar")

	if err != nil {
		ctx.JSON(400, gin.H{"error": "could not upload avatar"})
		return
	}
	if file.Size > (20 * 1024 * 1024) {
		ctx.JSON(400, gin.H{"error": "file is too big"})
		return
	}

	opened, err := file.Open()
	if err != nil {
		ctx.JSON(400, gin.H{"error": "could not upload avatar"})
		return
	}
	buf, err := validateImage(opened)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "could not upload avatar"})
		return
	}

	fullReader := io.MultiReader(bytes.NewReader(buf), opened)
	defer opened.Close()

	if user.AvatarURL != "" {
		err := h.minio.DeleteFile(user.AvatarURL)
		if err != nil {
			slog.Warn("could not delete old avatar")
		}
	}

	url, err := h.minio.UploadFile(fullReader, file.Size, file.Filename)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "could not upload avatar"})
		return
	}

	err = h.service.UpdateAvatarURL(userID, url)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "could not upload avatar"})
		return
	}

	ctx.JSON(200, gin.H{"avatarUrl": url})
}

func (h *UserHandler) DeleteAvatar(ctx *gin.Context) {
	userID := uint(ctx.GetFloat64("userID"))
	if userID == 0 {
		ctx.JSON(400, gin.H{"error": "cannot find user id"})
		return
	}

	user, err := h.service.GetUserByID(uint(userID))
	if err != nil {
		if err.Error() == "user not found" {
			ctx.JSON(404, gin.H{"error": "user does not exists"})
			return
		}
		ctx.JSON(500, gin.H{"error": "cannot find user"})
		return
	}

	if user.AvatarURL == "" {
		ctx.JSON(400, gin.H{"error": "no avatar to delete"})
		return
	}

	err = h.minio.DeleteFile(user.AvatarURL)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "could not delete avatar"})
		return
	}

	err = h.service.UpdateAvatarURL(userID, "")
	if err != nil {
		ctx.JSON(500, gin.H{"error": "could not delete avatar"})
		return
	}

	ctx.JSON(200, gin.H{})
}
