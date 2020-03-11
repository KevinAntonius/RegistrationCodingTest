package model

import (
    "time"
    null "gopkg.in/guregu/null.v3"
)

type User struct{
    ID              int             `json:"id,omitempty" gorm:"primary_key;AUTO_INCREMENT"`
    Mobile          null.String     `json:"mobile,omitempty" gorm:"NOT NULL;UNIQUE"`
    FirstName       null.String     `json:"fname,omitempty" gorm:"NOT NULL"`
    LastName		null.String		`json:"lname,omitempty" gorm:"NOT NULL"`
    BirthDate		*time.Time	    `json:"birthdate,omitempty"`
    Gender          null.String     `json:"gender,omitempty"`
    Email           null.String     `json:"email,omitempty" gorm:"NOT NULL;UNIQUE"`
}

func (User) TableName() string{
    return "users"
}