package mocks

import (
	"github.com/stretchr/testify/mock"
)

// MockKafkaProducer is a mock implementation of the KafkaProducer interface
type MockKafkaProducer struct {
	mock.Mock
}

// ProduceMessage provides a mock function for sending messages to Kafka
func (m *MockKafkaProducer) ProduceMessage(data []byte, topic string) error {
	args := m.Called(data, topic)
	return args.Error(0)
}
