package v4

type Config struct {
	Endpoint         string `json:"endpoint"`
	AccessKey        string `json:"accessKey"`
	SecretKey        string `json:"secretKey"`
	TopicNormal      string `json:"topicNormal"`
	TopicDelay       string `json:"topicDelay"`
	TopicGlobalFifo  string `json:"topicGlobalFifo"`
	TopicTransaction string `json:"topicTransaction"`
	TopicRegionFifo  string `json:"topicRegionFifo"`
	InstanceId       string `json:"instanceId"`
	GroupId          string `json:"groupId"`
}
