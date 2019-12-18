package transaction

import "strings"

type validateUpdateDto struct {
	dto   updateDto
	error map[string]string
}

func newValidateUpdateDto(dto *updateDto) *validateUpdateDto {
	return &validateUpdateDto{
		dto:   *dto,
		error: make(map[string]string),
	}
}

func (vid *validateUpdateDto) validate() bool {
	if vid.dto.Amount <= 0 {
		vid.error["amount"] = "Amount must be greater than 0"
	}

	if len(vid.dto.Description) > 250 {
		vid.error["description"] = "Description length is cannot be greater than 255 characters"
	}

	if strings.TrimSpace(vid.dto.Category) == "" {
		vid.error["category"] = "Category is cannot be empty"
	}

	if vid.dto.Type != "INCOME" && vid.dto.Type != "EXPENSE" {
		vid.error["type"] = "Type is must be INCOME or EXPENSE"
	}

	if strings.TrimSpace(vid.dto.Type) == "" {
		vid.error["type"] = "Type is cannot be empty"
	}

	if len(vid.error) > 0 {
		return false
	}

	return true
}
