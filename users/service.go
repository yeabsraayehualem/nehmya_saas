package users


type UserService struct {
	repository UserRepository
}

func NewService(r UserRepository) *UserService{
	return &UserService{
		repository: r,
	}
}

func (s *UserService) CreateUser(user User) error{
	return s.repository.Create(user)
}