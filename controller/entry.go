package controller

import (
	"github/be/common"
	"github/be/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetAllEntriesResult struct {
	Data string `json:"data"`
}

//	@BasePath	/api/v1

// AddEntry handles adding a new entry
//
//	@Summary	Add a new entry
//	@Accept		json
//	@Produce	json
//	@Param		input	body		model.Entry	true	"Entry details"
//	@Success	201		{object}	model.Entry
//	@Router		/entry [post]
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

// GetAllEntries handles getting all entries
//
//	@Summary	Get all entries
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	GetAllEntriesResult
//	@Router		/entry [get]
func GetAllEntries(ctx *gin.Context) {
	user, err := common.CurrentUser(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user"})
		return
	}

	user.Entries, err = model.GetAllEntries(user.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "couldn't get entries"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": user.Entries})
}
