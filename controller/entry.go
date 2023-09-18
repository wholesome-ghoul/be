package controller

import (
	"github/be/common"
	"github/be/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddEntry(ctx *gin.Context) {
	var entry model.Entry
	if err := ctx.ShouldBindJSON(&entry); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	user, err := common.CurrentUser(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user"})
		return
	}

	entry.UserID = user.ID

	savedEntry, err := entry.Save()
	if err != nil {
		message := "couldn't save entry"
		if _, ok := err.(*model.EntryError); ok {
			message = err.Error()
		}

		ctx.JSON(http.StatusBadRequest, gin.H{"error": message})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": savedEntry})
}

func GetAllEntries(ctx *gin.Context) {
	user, err := common.CurrentUser(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": user.Entries})
}
