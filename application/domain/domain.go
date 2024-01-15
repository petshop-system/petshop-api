package domain

import "time"

type CustomerDomain struct {
	ID          int64
	Name        string
	DateCreated time.Time
}

type PhoneDomain struct {
	ID        int64
	Number    string
	CodeArea  string
	PhoneType string
}

type SpeciesDomain struct {
	ID   int64
	Name string
}

type BreedDomain struct {
	ID   int64
	Name string
}

type AddressDomain struct {
	ID     int64
	Street string
	Number string
}

type PersonDomain struct {
	ID         int64
	Document   string
	PersonType string
}

type ScheduleMessage struct {
	Booking                    string
	PetId                      int
	ServiceEmployeeAttentionId int
}
