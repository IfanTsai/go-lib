package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type Message struct {
	Topic string
	Value []byte
}

type Client struct {
	producer sarama.SyncProducer
	ch       chan *Message
	started  bool
}

func NewClient(addrs []string, username, password string) (*Client, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true
	config.Net.SASL.Enable = true
	config.Net.SASL.Mechanism = sarama.SASLTypePlaintext
	config.Net.SASL.User = username
	config.Net.SASL.Password = password

	producer, err := sarama.NewSyncProducer(addrs, config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to new kafka client")
	}

	return &Client{
		producer: producer,
	}, nil
}

func (c *Client) SendMessageSync(msg *Message) error {
	_, _, err := c.producer.SendMessage(&sarama.ProducerMessage{
		Topic: msg.Topic,
		Value: sarama.ByteEncoder(msg.Value),
	})

	return errors.Wrapf(err, "failed to send message(Topic: %s) to kafka", msg.Topic)
}

func (c *Client) SendMessageAsync(msg *Message) {
	if c.ch == nil || !c.started {
		return
	}

	c.ch <- msg
}

func (c *Client) StartSendingMessage(sendGoroutineCount, channelSize int, logger *zap.Logger) {
	if c.started {
		return
	}

	if c.ch == nil {
		c.ch = make(chan *Message, channelSize)
	}

	c.started = true

	for i := 0; i < sendGoroutineCount; i++ {
		go func() {
			for {
				msg, ok := <-c.ch
				if !ok {
					return
				}

				if err := c.SendMessageSync(msg); err != nil {
					logger.Error("failed to send message", zap.String("error", err.Error()))
				}
			}
		}()
	}
}

func (c *Client) GetChannelStatus() (int, int) {
	return len(c.ch), cap(c.ch)
}

func (c *Client) StopSendingMessage() {
	c.started = false
	close(c.ch)
}
