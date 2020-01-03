package transaction

import (
	"strings"

	"github.com/novaladip/geldstroom-api-go/core/validator"
)

type CreateDto struct {
	Amount      int64  `form:"amount"`
	Description string `form:"description"`
	Category    string `form:"category"`
	Type        string `form:"type"`
}

func (dto CreateDto) Validate() validator.Validate {
	v := validator.New()

	if dto.Amount <= 0 {
		v.Error["amount"] = "Amount must be greater than 0"
	}

	if len(dto.Description) > 256 {
		v.Error["desciprition"] = "Description length is cannot be greater than 256 character"
	}

	if strings.TrimSpace(dto.Category) == "" {
		v.Error["category"] = "Category is cannot be empty"
	}

	if strings.ToUpper(dto.Type) != "INCOME" && strings.ToUpper(dto.Type) != "EXPENSE" {
		v.Error["type"] = "Type must be INCOME or EXPENSE"
	}

	if len(v.Error) > 0 {
		v.IsValid = false
	}

	return v
}

type UpdateDto struct {
	Amount      int64  `form:"amount"`
	Description string `form:"description"`
	Category    string `form:"category"`
	Type        string `form:"type"`
}

func (dto UpdateDto) Validate() validator.Validate {
	v := validator.New()

	if dto.Amount <= 0 {
		v.Error["amount"] = "Amount must be greater than 0"
	}

	if len(dto.Description) > 256 {
		v.Error["desciprition"] = "Description length is cannot be greater than 256 character"
	}

	if strings.TrimSpace(dto.Category) == "" {
		v.Error["category"] = "Category is cannot be empty"
	}

	if strings.ToUpper(dto.Type) != "INCOME" && strings.ToUpper(dto.Type) != "EXPENSE" {
		v.Error["type"] = "Type must be INCOME or EXPENSE"
	}

	if len(v.Error) > 0 {
		v.IsValid = false
	}

	return v
}
