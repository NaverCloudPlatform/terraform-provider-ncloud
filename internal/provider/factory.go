package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/provider/fwprovider"
)

// ProtoV6ProviderServerFactory returns a muxed terraform-plugin-go protocol v6 provider factory function.
// This factory function is suitable for use with the terraform-plugin-go Serve function.
// The primary (Plugin SDK) provider server is also returned (useful for testing).
func ProtoV6ProviderServerFactory(ctx context.Context) (func() tfprotov6.ProviderServer, *schema.Provider, error) {
	primary := New(ctx)

	upgradedSdkProvider, err := tf5to6server.UpgradeServer(
		ctx,
		primary.GRPCProvider,
	)
	if err != nil {
		return nil, nil, err
	}

	servers := []func() tfprotov6.ProviderServer{
		func() tfprotov6.ProviderServer {
			return upgradedSdkProvider
		},
		providerserver.NewProtocol6(fwprovider.New(primary)),
	}

	muxServer, err := tf6muxserver.NewMuxServer(ctx, servers...)

	if err != nil {
		return nil, nil, err
	}

	return muxServer.ProviderServer, primary, nil
}
