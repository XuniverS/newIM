package queue

// MessageQueue 消息队列接口
type MessageQueue interface {
	Enqueue(message interface{}) error
	Dequeue() (interface{}, error)
	GetQueueLength() (int, error)
	Clear() error
}

// PendingMessage 待发送消息
type PendingMessage struct {
	ID        string
	Sender    string
	Receiver  string
	Content   string
	Timestamp int64
	Retries   int
}

// QueueConsumer 队列消费者
type QueueConsumer struct {
	queue MessageQueue
}

func NewQueueConsumer(queue MessageQueue) *QueueConsumer {
	return &QueueConsumer{queue: queue}
}

// ConsumeMessage 消费消息
func (qc *QueueConsumer) ConsumeMessage() (interface{}, error) {
	return qc.queue.Dequeue()
}

// ProduceMessage 生产消息
func (qc *QueueConsumer) ProduceMessage(msg interface{}) error {
	return qc.queue.Enqueue(msg)
}
