package users

// UserService contains the business logic on top of the Repository.
type UserService struct {
	repository Repository
}

// NewService wires a user service to a repository.
func NewService(r Repository) *UserService {
	return &UserService{
		repository: r,
	}
}

func (s *UserService) CreateUser(user *User) error {
	return s.repository.Create(user)
}

func (s *UserService) Authenticate(login, password string) (*User, error) {
	return s.repository.Authenticate(login, password)
}

func (s *UserService) GetByID(id uint) (*User, error) {
	return s.repository.GetByID(id)
}
