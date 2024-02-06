package topic

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
)

type TopicRepository interface {
	CreateTopic(ctx context.Context, topic *Topic) error
	GetTopics(ctx context.Context) ([]*Topic, error)
	UpdateTopic(ctx context.Context, id string, topic *Topic) error
	DeleteTopic(ctx context.Context, id string) error
	IsTopicNameUnique(ctx context.Context, name string) (bool, error)
}

type RedisTopicRepository struct {
	client *redis.Client
}

func NewRedisTopicRepository(client *redis.Client) TopicRepository {
	return &RedisTopicRepository{client: client}
}

func (r *RedisTopicRepository) CreateTopic(ctx context.Context, topic *Topic) error {
	unique, err := r.IsTopicNameUnique(ctx, topic.Name)
	if err != nil {
		return err
	}
	if !unique {
		return fmt.Errorf("topic name '%s' already exists", topic.Name)
	}

	id, err := r.client.Incr(ctx, "topic_id").Result()
	if err != nil {
		return err
	}
	topic.ID = fmt.Sprintf("%d", id)

	data, err := json.Marshal(topic)
	if err != nil {
		return err
	}

	err = r.client.Set(ctx, getTopicKey(topic.ID), data, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *RedisTopicRepository) GetTopics(ctx context.Context) ([]*Topic, error) {
	var topics []*Topic

	keys := r.client.Keys(ctx, "topic:*").Val()
	for _, key := range keys {
		data, err := r.client.Get(ctx, key).Bytes()
		if err != nil {
			return nil, err
		}

		var topic Topic
		err = json.Unmarshal(data, &topic)
		if err != nil {
			return nil, err
		}

		topics = append(topics, &topic)
	}

	return topics, nil
}

func (r *RedisTopicRepository) UpdateTopic(ctx context.Context, id string, topic *Topic) error {
	unique, err := r.IsTopicNameUnique(ctx, topic.Name)
	if err != nil {
		return err
	}
	if !unique {
		return fmt.Errorf("topic name '%s' already exists", topic.Name)
	}

	data, err := json.Marshal(topic)
	if err != nil {
		return err
	}

	err = r.client.Set(ctx, getTopicKey(id), data, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *RedisTopicRepository) DeleteTopic(ctx context.Context, id string) error {
	err := r.client.Del(ctx, getTopicKey(id)).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *RedisTopicRepository) IsTopicNameUnique(ctx context.Context, name string) (bool, error) {
	keys := r.client.Keys(ctx, "topic:*").Val()
	for _, key := range keys {
		data, err := r.client.Get(ctx, key).Bytes()
		if err != nil {
			return false, err
		}

		var topic Topic
		err = json.Unmarshal(data, &topic)
		if err != nil {
			return false, err
		}

		if topic.Name == name {
			return false, nil
		}
	}
	return true, nil
}
