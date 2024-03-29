package usecase

import (
	"fmt"
	"math"
	"sync"

	"github.com/jordyf15/thullo-api/board"
	"github.com/jordyf15/thullo-api/board_member"
	"github.com/jordyf15/thullo-api/custom_errors"
	"github.com/jordyf15/thullo-api/models"
	"github.com/jordyf15/thullo-api/storage"
	"github.com/jordyf15/thullo-api/unsplash"
	"github.com/jordyf15/thullo-api/user"
	"github.com/jordyf15/thullo-api/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type boardUsecase struct {
	boardRepo       board.Repository
	unsplashRepo    unsplash.Repository
	boardMemberRepo board_member.Repository
	userRepo        user.Repository
	storage         storage.Storage
}

func NewBoardUsecase(boardRepo board.Repository, unsplashRepo unsplash.Repository, boardMemberRepo board_member.Repository, userRepo user.Repository, storage storage.Storage) board.Usecase {
	return &boardUsecase{boardRepo: boardRepo, unsplashRepo: unsplashRepo, boardMemberRepo: boardMemberRepo, userRepo: userRepo, storage: storage}
}

func (usecase *boardUsecase) Create(userID primitive.ObjectID, title string, visibility string, cover map[string]interface{}) error {
	errors := []error{}

	source := cover["source"].(string)
	focalPointY := cover["fp_y"].(float64)
	photoID := cover["photo_id"].(string)

	focalPointY = math.Round(focalPointY*1000) / 1000

	if _, exist := board.BoardCoverSources[source]; !exist {
		errors = append(errors, custom_errors.ErrInvalidCoverSource)
	}

	if focalPointY > 1 {
		errors = append(errors, custom_errors.ErrUnsplashFocalPointYTooHigh)
	}

	if focalPointY < 0 {
		errors = append(errors, custom_errors.ErrUnsplashFocalPointYTooLow)
	}

	if title == "" {
		errors = append(errors, custom_errors.ErrBoardTitleEmpty)
	}

	if visibility != models.BoardVisibilityPrivate && visibility != models.BoardVisibilityPublic {
		errors = append(errors, custom_errors.ErrBoardInvalidVisibility)
	}

	if len(errors) > 0 {
		return &custom_errors.MultipleErrors{Errors: errors}
	}

	imageFiles, err := usecase.unsplashRepo.GetImagesForID(photoID, focalPointY)
	if err != nil {
		return err
	}

	boardCover := &models.BoardCover{
		PhotoID:     photoID,
		Source:      source,
		FocalPointY: focalPointY,
		Images:      make(models.Images, len(board.BoardCoverSizes)),
	}

	for i, width := range board.BoardCoverSizes {
		image := &models.Image{
			Width: width,
		}

		boardCover.Images[i] = image
	}

	var wg sync.WaitGroup
	uploadChannels := make(chan error, len(boardCover.Images))
	wg.Add(len(boardCover.Images))

	for idx, image := range boardCover.Images {
		imageFile := imageFiles[idx]
		name := utils.RandString(8)
		fileName := fmt.Sprintf("%s.%s", name, utils.GetFileExtension(imageFile.Name()))

		metaData := map[string]string{
			"name":        fileName,
			"title":       name,
			"description": fmt.Sprintf("cover picture of board %s with width of %v", title, image.Width),
		}

		go usecase.storage.UploadFile(uploadChannels, &wg, image, imageFile, metaData)
	}

	wg.Wait()
	close(uploadChannels)

	for err = range uploadChannels {
		if err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		return &custom_errors.MultipleErrors{Errors: errors}
	}

	<-uploadChannels

	_board := &models.Board{
		Title:   title,
		OwnerID: userID,
		Cover:   boardCover,
	}
	_board.SetVisibility(visibility)
	_board.EmptyImageURLs()

	err = usecase.boardRepo.Create(_board)
	if err != nil {
		return err
	}

	boardMember := &models.BoardMember{
		UserID:  userID,
		BoardID: _board.ID,
		Role:    models.MemberRoleAdmin,
	}

	err = usecase.boardMemberRepo.Create(boardMember)
	if err != nil {
		return err
	}

	return nil
}

func (usecase *boardUsecase) UpdateVisibility(requesterID, boardID primitive.ObjectID, visibility string) error {
	if visibility != models.BoardVisibilityPrivate && visibility != models.BoardVisibilityPublic {
		return custom_errors.ErrBoardInvalidVisibility
	}

	board, err := usecase.boardRepo.GetBoardByID(boardID)
	if err != nil {
		return err
	}

	boardMembers, err := usecase.boardMemberRepo.GetBoardMembers(boardID)
	if err != nil {
		return err
	}

	var requesterBoardMember *models.BoardMember
	for _, boardMember := range boardMembers {
		if boardMember.UserID == requesterID {
			requesterBoardMember = boardMember
			break
		}
	}

	if requesterBoardMember == nil || requesterBoardMember.Role == models.MemberRoleMember {
		return custom_errors.ErrNotAuthorized
	}

	board.Visibility = models.BoardVisibility(visibility)

	err = usecase.boardRepo.Update(board)
	if err != nil {
		return err
	}

	return nil
}

func (usecase *boardUsecase) UpdateTitle(requesterID, boardID primitive.ObjectID, title string) error {
	if title == "" {
		return custom_errors.ErrBoardTitleEmpty
	}

	board, err := usecase.boardRepo.GetBoardByID(boardID)
	if err != nil {
		return err
	}

	boardMembers, err := usecase.boardMemberRepo.GetBoardMembers(boardID)
	if err != nil {
		return err
	}

	var requesterBoardMember *models.BoardMember
	for _, boardMember := range boardMembers {
		if boardMember.UserID == requesterID {
			requesterBoardMember = boardMember
			break
		}
	}

	if requesterBoardMember == nil || requesterBoardMember.Role == models.MemberRoleMember {
		return custom_errors.ErrNotAuthorized
	}

	board.Title = title

	err = usecase.boardRepo.Update(board)
	if err != nil {
		return err
	}

	return nil
}

func (usecase *boardUsecase) UpdateDescription(requesterID, boardID primitive.ObjectID, description string) error {
	board, err := usecase.boardRepo.GetBoardByID(boardID)
	if err != nil {
		return err
	}

	boardMembers, err := usecase.boardMemberRepo.GetBoardMembers(boardID)
	if err != nil {
		return err
	}

	var requesterBoardMember *models.BoardMember
	for _, boardMember := range boardMembers {
		if boardMember.UserID == requesterID {
			requesterBoardMember = boardMember
			break
		}
	}

	if requesterBoardMember == nil {
		return custom_errors.ErrNotAuthorized
	}

	board.Description = description

	err = usecase.boardRepo.Update(board)
	if err != nil {
		return err
	}

	return nil
}

func (usecase *boardUsecase) AddMember(requesterID, boardID, memberID primitive.ObjectID) error {
	board, err := usecase.boardRepo.GetBoardByID(boardID)
	if err != nil {
		return err
	}

	_, err = usecase.userRepo.GetByID(memberID)
	if err != nil {
		return err
	}

	boardMembers, err := usecase.boardMemberRepo.GetBoardMembers(board.ID)
	if err != nil {
		return err
	}

	// check whether requester and new member is already a board member
	isRequesterBoardMember := false
	isNewMemberBoardMember := false
	for _, boardMember := range boardMembers {
		if boardMember.UserID == requesterID {
			isRequesterBoardMember = true
		}
		if boardMember.UserID == memberID {
			isNewMemberBoardMember = true
		}
		if isNewMemberBoardMember && isRequesterBoardMember {
			break
		}
	}

	// if requester is not a board member that means he/she is not authorized to add member
	if !isRequesterBoardMember {
		return custom_errors.ErrNotAuthorized
	}
	// if the new member is already a board member than the operation should not proceed
	if isNewMemberBoardMember {
		return custom_errors.ErrUserIsAlreadyBoardMember
	}

	boardMember := &models.BoardMember{
		UserID:  memberID,
		BoardID: board.ID,
		Role:    models.MemberRoleMember,
	}

	err = usecase.boardMemberRepo.Create(boardMember)
	if err != nil {
		return err
	}

	return nil
}

func (usecase *boardUsecase) UpdateMemberRole(requesterID, boardID, memberID primitive.ObjectID, role string) error {
	if role != models.MemberRoleMember && role != models.MemberRoleAdmin {
		return custom_errors.ErrInvalidBoardMemberRole
	}

	boardMembers, err := usecase.boardMemberRepo.GetBoardMembers(boardID)
	if err != nil {
		return err
	}

	// if there are no board members it means there are no board
	// a board will always atleast have 1 member
	if len(boardMembers) == 0 {
		return custom_errors.ErrRecordNotFound
	}

	var requesterBoardMember *models.BoardMember
	var memberBoardMember *models.BoardMember
	adminCount := 0
	for _, boardMember := range boardMembers {
		if boardMember.UserID == requesterID {
			requesterBoardMember = boardMember
		}
		if boardMember.UserID == memberID {
			memberBoardMember = boardMember
		}
		if boardMember.Role == models.MemberRoleAdmin {
			adminCount += 1
		}
	}

	// if requester is not a member or his/her role is a member than he is not authorized
	// to perform this operation
	if requesterBoardMember == nil || requesterBoardMember.Role == models.MemberRoleMember {
		return custom_errors.ErrNotAuthorized
	}

	if memberBoardMember == nil {
		return custom_errors.ErrRecordNotFound
	}

	if adminCount == 1 && memberBoardMember.Role == models.MemberRoleAdmin && role == models.MemberRoleMember {
		return custom_errors.ErrBoardMustHaveAnAdmin
	}

	err = usecase.boardMemberRepo.UpdateBoardMemberRole(memberBoardMember.ID, models.MemberRole(role))
	if err != nil {
		return err
	}

	return nil
}

func (usecase *boardUsecase) DeleteMember(requesterID, boardID, memberID primitive.ObjectID) error {
	boardMembers, err := usecase.boardMemberRepo.GetBoardMembers(boardID)
	if err != nil {
		return err
	}

	// if there are no board members it means there are no board
	// a board will always atleast have 1 member
	if len(boardMembers) == 0 {
		return custom_errors.ErrRecordNotFound
	}

	var requesterBoardMember *models.BoardMember
	var memberBoardMember *models.BoardMember
	adminCount := 0
	for _, boardMember := range boardMembers {
		if boardMember.UserID == requesterID {
			requesterBoardMember = boardMember
		}
		if boardMember.UserID == memberID {
			memberBoardMember = boardMember
		}
		if boardMember.Role == models.MemberRoleAdmin {
			adminCount += 1
		}
	}

	// if requester is not a member or his/her role is a member than he is not authorized
	// to perform this operation
	if requesterBoardMember == nil || requesterBoardMember.Role == models.MemberRoleMember {
		return custom_errors.ErrNotAuthorized
	}

	if memberBoardMember == nil {
		return custom_errors.ErrRecordNotFound
	}

	if adminCount == 1 && memberBoardMember.Role == models.MemberRoleAdmin {
		return custom_errors.ErrBoardMustHaveAnAdmin
	}

	err = usecase.boardMemberRepo.DeleteBoardMemberByID(memberBoardMember.ID)
	if err != nil {
		return err
	}

	return nil
}
