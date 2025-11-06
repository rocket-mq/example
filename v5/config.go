package v5

type Config struct {
	Debug            bool   `json:"debug"`
	Endpoint         string `json:"endpoint"`
	AccessKey        string `json:"accessKey"`
	SecretKey        string `json:"secretKey"`
	NameSpace        string `json:"nameSpace"`
	ConsumeGroup     string `json:"consumeGroup"`
	TopicNormal      string `json:"topicNormal"`
	TopicDelay       string `json:"topicDelay"`
	TopicFifo        string `json:"topicFifo"`
	TopicTransaction string `json:"topicTransaction"`
}
