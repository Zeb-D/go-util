package pulsar

import (
	"context"
	"github.com/TuyaInc/pulsar-client-go/pkg/api"

	"github.com/tetratelabs/go2sky"
	"github.com/tetratelabs/go2sky/reporter/grpc/common"
)

var (
	skyTrace *go2sky.Tracer
)

func SetTracer(t *go2sky.Tracer) {
	skyTrace = t
}

type traceHook struct {
	topic string
}

func (t *traceHook) OnSend(ctx context.Context, msg *api.MessageMetadata, payload []byte) {
	if skyTrace == nil {
		return
	}
	k := TraceKey
	kv := api.KeyValue{
		Key: &k,
	}
	span, err := skyTrace.CreateExitSpan(ctx, getOperationName(t.topic), "ferrari", func(header string) error {
		kv.Value = &header
		return nil
	})
	if err != nil {
		return
	}
	msg.Properties = append(msg.Properties, &kv)
	defer span.End()
	span.SetSpanLayer(common.SpanLayer_MQ)
	span.Tag(TagPulsarMQTopic, t.topic)
	span.Tag(TagPulsarMQPayload, string(payload))
}

func getOperationName(topic string) string {
	return "PULSAR_PUBLISH_" + topic
}

func getConsumerOperationName(topic string) string {
	return "PULSAR_CONSUMER_" + topic
}
