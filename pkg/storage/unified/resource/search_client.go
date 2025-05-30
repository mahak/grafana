package resource

import (
	"context"

	"google.golang.org/grpc"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/grafana/grafana/pkg/storage/unified/resourcepb"
)

type DualWriter interface {
	IsEnabled(schema.GroupResource) bool
	ReadFromUnified(context.Context, schema.GroupResource) (bool, error)
}

func NewSearchClient(dual DualWriter, gr schema.GroupResource, unifiedClient resourcepb.ResourceIndexClient,
	legacyClient resourcepb.ResourceIndexClient) resourcepb.ResourceIndexClient {
	if dual.IsEnabled(gr) {
		return &searchWrapper{
			dual:          dual,
			groupResource: gr,
			unifiedClient: unifiedClient,
			legacyClient:  legacyClient,
		}
	}
	//nolint:errcheck
	if ok, _ := dual.ReadFromUnified(context.Background(), gr); ok {
		return unifiedClient
	}
	return legacyClient
}

type searchWrapper struct {
	dual          DualWriter
	groupResource schema.GroupResource

	unifiedClient resourcepb.ResourceIndexClient
	legacyClient  resourcepb.ResourceIndexClient
}

func (s *searchWrapper) GetStats(ctx context.Context, in *resourcepb.ResourceStatsRequest,
	opts ...grpc.CallOption) (*resourcepb.ResourceStatsResponse, error) {
	client := s.legacyClient
	unified, err := s.dual.ReadFromUnified(ctx, s.groupResource)
	if err != nil {
		return nil, err
	}
	if unified {
		client = s.unifiedClient
	}
	return client.GetStats(ctx, in, opts...)
}

func (s *searchWrapper) Search(ctx context.Context, in *resourcepb.ResourceSearchRequest,
	opts ...grpc.CallOption) (*resourcepb.ResourceSearchResponse, error) {
	client := s.legacyClient
	unified, err := s.dual.ReadFromUnified(ctx, s.groupResource)
	if err != nil {
		return nil, err
	}
	if unified {
		client = s.unifiedClient
	}
	return client.Search(ctx, in, opts...)
}
