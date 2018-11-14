package handy

type UserService struct {
}

func NewUserService() *UserService {
	return &UserService{}
}

func (s *UserService) CreateDeck(userID UUID, deck Deck) error {


	return nil
}