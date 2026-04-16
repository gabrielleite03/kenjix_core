package service

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/gabrielleite03/kenjix_core/internal/dto"
	"github.com/gabrielleite03/kenjix_domain/model"
	"github.com/gabrielleite03/kenjix_persist/repository"
	"github.com/google/uuid"
)

const expenseAttachmentsS3Path = "uploads/expense/attachments/"

type ExpenseService interface {
	FindAll() ([]dto.ExpenseDTO, error)
	FindByID(id int64) (*dto.ExpenseDTO, error)
	Create(data *dto.ExpenseCreateUpdateDTO) error
	Update(id int64, data *dto.ExpenseCreateUpdateDTO) (*dto.ExpenseDTO, error)
	Delete(id int64) error
	FindAllExpenseCategories() ([]dto.ExpenseCategoryDTO, error)
}

type expenseService struct {
	dao repository.ExpenseDAO
}

func NewExpenseService(dao repository.ExpenseDAO) ExpenseService {
	return &expenseService{dao: dao}
}

func (s *expenseService) FindAll() ([]dto.ExpenseDTO, error) {
	expenses, err := s.dao.FindAll()
	if err != nil {
		return nil, err
	}

	var result []dto.ExpenseDTO
	for _, e := range expenses {
		result = append(result, mapToDTO(e))
	}

	return result, nil
}

func (s *expenseService) FindByID(id int64) (*dto.ExpenseDTO, error) {
	exp, err := s.dao.FindByID(id)
	if err != nil {
		return nil, err
	}

	dto := mapToDTO(*exp)
	return &dto, nil
}

func (s *expenseService) Create(data *dto.ExpenseCreateUpdateDTO) error {
	exp := &model.Expense{
		Description: data.Description,
		CategoryID:  data.CategoryID,
		Amount:      data.Amount,
		Date:        data.Date,
		Status:      string(data.Status),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := saveExpenseAttachmentsToS3(data)
	if err != nil {
		return err
	}

	exp.Attachments = []model.ExpenseAttachment{}
	for _, att := range data.Attachments {
		exp.Attachments = append(exp.Attachments, model.ExpenseAttachment{
			URL: att.URL,
		})
	}

	if err := s.dao.Create(exp); err != nil {
		return err
	}

	return nil
}

func (s *expenseService) Update(id int64, data *dto.ExpenseCreateUpdateDTO) (*dto.ExpenseDTO, error) {
	exp, err := s.dao.FindByID(id)
	if err != nil {
		return nil, err
	}

	exp.Description = data.Description
	exp.CategoryID = data.CategoryID
	exp.Amount = data.Amount
	exp.Date = data.Date
	exp.Status = string(data.Status)

	for _, att := range exp.Attachments {
		deleteFromS3(att.URL)
	}

	saveExpenseAttachmentsToS3(data)

	exp.Attachments = []model.ExpenseAttachment{}
	for _, att := range data.Attachments {
		exp.Attachments = append(exp.Attachments, model.ExpenseAttachment{
			URL: att.URL,
		})
	}

	if err := s.dao.Update(exp); err != nil {
		return nil, err
	}

	return s.FindByID(id)
}

func (s *expenseService) Delete(id int64) error {

	data, err := s.dao.FindByID(id)
	if err != nil {
		return err
	}

	for _, att := range data.Attachments {
		deleteFromS3(att.URL)
	}
	return s.dao.Delete(id)
}

func mapToDTO(e model.Expense) dto.ExpenseDTO {
	result := dto.ExpenseDTO{
		ID:          e.ID,
		Description: e.Description,
		CategoryID:  e.CategoryID,
		Amount:      e.Amount,
		Date:        e.Date,
		Status:      dto.ExpenseStatus(e.Status),
	}

	if e.Category != nil {
		result.Category = &dto.ExpenseCategoryDTO{
			ID:   e.Category.ID,
			Name: e.Category.Name,
		}
	}

	for _, a := range e.Attachments {
		result.Attachments = append(result.Attachments, dto.ExpenseAttachmentDTO{
			ID:  a.ID,
			URL: a.URL,
		})
	}

	return result
}

func saveExpenseAttachmentsToS3(data *dto.ExpenseCreateUpdateDTO) error {

	imageList := []dto.ExpenseAttachmentDTO{}
	files := data.Files
	for _, fileHeader := range files {

		file, err := fileHeader.Open()
		if err != nil {
			return err
		}

		ext := filepath.Ext(fileHeader.Filename)
		fileName := fmt.Sprintf("%s%s", uuid.New().String(), ext)

		err = uploadToS3(file, fileName, expenseAttachmentsS3Path)
		file.Close()

		if err != nil {
			return err
		}

		// adicionar na lista
		imageList = append(imageList, dto.ExpenseAttachmentDTO{
			URL: "https://aws-s3-site-kejipet.s3.us-east-1.amazonaws.com/" + expenseAttachmentsS3Path + fileName,
		})
	}
	data.Attachments = imageList
	return nil
}

func (s *expenseService) FindAllExpenseCategories() ([]dto.ExpenseCategoryDTO, error) {
	categories, err := s.dao.FindAllExpenseCategories()
	if err != nil {
		return nil, err
	}
	var result []dto.ExpenseCategoryDTO
	for _, c := range categories {
		result = append(result, dto.ExpenseCategoryDTO{
			ID:   c.ID,
			Name: c.Name,
		})
	}
	return result, nil
}
