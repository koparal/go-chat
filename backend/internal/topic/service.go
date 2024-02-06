package topic

import "context"

type TopicService struct {
	Repository TopicRepository
}

func NewTopicService(repository TopicRepository) *TopicService {
	return &TopicService{
		Repository: repository,
	}
}

func (s *TopicService) CreateTopic(c context.Context, topic *Topic) error {
	err := s.Repository.CreateTopic(c, topic)
	if err != nil {
		return err
	}

	return nil
}

func (s *TopicService) GetTopics(c context.Context) ([]*Topic, error) {
	topics, err := s.Repository.GetTopics(c)
	if err != nil {
		return nil, err
	}

	return topics, nil
}

func (s *TopicService) UpdateTopic(c context.Context, id string, topic *Topic) error {
	err := s.Repository.UpdateTopic(c, id, topic)
	if err != nil {
		return err
	}

	return nil
}

func (s *TopicService) DeleteTopic(c context.Context, id string) error {
	err := s.Repository.DeleteTopic(c, id)
	if err != nil {
		return err
	}

	return nil
}
