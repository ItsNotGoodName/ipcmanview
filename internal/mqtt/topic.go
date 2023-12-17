package mqtt

func TopicServerState(topic Topic) string {
	return topic.Join("server", "state")
}
