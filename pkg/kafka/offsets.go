package kafka

import (
	"fmt"

	"github.com/IBM/sarama"
	"go.uber.org/multierr"
)

// getCommittedOffsets returns offsets committed by the group for topics/partitions.
// Not every [sarama.Broker] can be used, but group coordinator only.
// Negative value is returned, if no offset committed yet.
// Unfortunately [sarama.Client] doesn't expose such method, so have to write low level code.
// It's similar to [sarama.clusterAdmin.ListConsumerGroupOffsets].
func getCommittedOffsets(
	coordinator *sarama.Broker,
	groupID string,
	topicPartitions map[string]map[int32]struct{},
) (map[string]map[int32]int64, error) {
	req := new(sarama.OffsetFetchRequest)
	req.Version = 2 // like in sarama.clusterAdmin.ListConsumerGroupOffsets
	req.ConsumerGroup = groupID
	for topic, partitions := range topicPartitions {
		for partition := range partitions {
			req.AddPartition(topic, partition)
		}
	}

	resp, err := coordinator.FetchOffset(req)
	if err != nil {
		return nil, fmt.Errorf("can't fetch offset: %w", err)
	}
	if resp.Err != sarama.ErrNoError {
		return nil, fmt.Errorf("response with error: %w", resp.Err)
	}

	res := make(map[string]map[int32]int64)
	var errs error
	for topic, partitions := range topicPartitions {
		topicMap := make(map[int32]int64)
		for partition := range partitions {
			offset, errX := readOffsetFetchResponse(resp, topic, partition)
			errs = multierr.Append(errs, errX)
			topicMap[partition] = offset
		}
		res[topic] = topicMap
	}
	if errs != nil {
		return nil, errs
	}
	return res, nil
}

func readOffsetFetchResponse(resp *sarama.OffsetFetchResponse, topic string, partition int32) (int64, error) {
	block := resp.GetBlock(topic, partition)
	if block == nil {
		return 0, fmt.Errorf("block not found (topic: %s, partition %d)", topic, partition)
	}
	if block.Err != sarama.ErrNoError {
		return 0, fmt.Errorf("block with error (topic: %s, partition %d): %w", topic, partition, block.Err)
	}
	return block.Offset, nil
}
