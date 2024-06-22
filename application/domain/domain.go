package domain

type CustomerDomain struct {
	ID         int64
	Name       string
	Email      string
	Document   string
	PersonType string
	ContractID int64
	AddressID  int64
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
	ID           int64
	Street       string
	Number       string
	Complement   string
	Block        string
	Neighborhood string
	ZipCode      string
	City         string
	State        string
	Country      string
}

type ScheduleMessage struct {
	Booking                    string
	PetId                      int
	ServiceEmployeeAttentionId int
}
