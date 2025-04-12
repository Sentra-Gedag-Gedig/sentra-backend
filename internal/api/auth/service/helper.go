package authService

import (
	"ProjectGolang/internal/api/auth"
	"ProjectGolang/internal/entity"
	"time"
)

func GetUserDifferenceData(DbUser entity.User, NewUser auth.UpdateUserRequest) (entity.User, error) {
	var result entity.User
	result.ID = DbUser.ID

	if NewUser.Name != "" && NewUser.Name != DbUser.Name {
		result.Name = NewUser.Name
	}

	if NewUser.Religion != "" && NewUser.Religion != DbUser.Religion {
		result.Religion = NewUser.Religion
	}

	if NewUser.MaritalStatus != "" && NewUser.MaritalStatus != DbUser.MaritalStatus {
		result.MaritalStatus = NewUser.MaritalStatus
	}

	if NewUser.Profession != "" && NewUser.Profession != DbUser.Profession {
		result.Profession = NewUser.Profession
	}

	if NewUser.District != "" && NewUser.District != DbUser.District {
		result.District = NewUser.District
	}

	if NewUser.Village != "" && NewUser.Village != DbUser.Village {
		result.Village = NewUser.Village
	}

	if NewUser.Address != "" && NewUser.Address != DbUser.Address {
		result.Address = NewUser.Address
	}

	if NewUser.BirthPlace != "" && NewUser.BirthPlace != DbUser.BirthPlace {
		result.BirthPlace = NewUser.BirthPlace
	}

	if NewUser.Gender != "" && NewUser.Gender != DbUser.Gender {
		result.Gender = NewUser.Gender
	}

	if NewUser.NationalIdentityNumber != "" && NewUser.NationalIdentityNumber != DbUser.NationalIdentityNumber {
		result.NationalIdentityNumber = NewUser.NationalIdentityNumber
	}

	zeroTime := time.Time{}
	if NewUser.CardValidUntil != zeroTime && NewUser.CardValidUntil != DbUser.CardValidUntil {
		result.CardValidUntil = NewUser.CardValidUntil
	}

	if NewUser.Citizenship != "" && NewUser.Citizenship != DbUser.Citizenship {
		result.Citizenship = NewUser.Citizenship
	}

	if NewUser.NeighborhoodCommunityUnit != "" && NewUser.NeighborhoodCommunityUnit != DbUser.NeighborhoodCommunityUnit {
		result.NeighborhoodCommunityUnit = NewUser.NeighborhoodCommunityUnit
	}

	if NewUser.BirthDate != zeroTime && NewUser.BirthDate != DbUser.BirthDate {
		result.BirthDate = NewUser.BirthDate
	}

	return result, nil
}

func MakeUserData(user entity.User) map[string]interface{} {
	return map[string]interface{}{
		"id":       user.ID,
		"email":    user.Email,
		"username": user.Name,
	}
}
