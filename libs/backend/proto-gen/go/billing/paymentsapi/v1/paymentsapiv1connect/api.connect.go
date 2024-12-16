// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: billing/paymentsapi/v1/api.proto

package paymentsapiv1connect

import (
	connect "connectrpc.com/connect"
	context "context"
	errors "errors"
	v1 "libs/backend/proto-gen/go/billing/paymentsapi/v1"
	http "net/http"
	strings "strings"
)

// This is a compile-time assertion to ensure that this generated file and the connect package are
// compatible. If you get a compiler error that this constant is not defined, this code was
// generated with a version of connect newer than the one compiled into your binary. You can fix the
// problem by either regenerating this code with an older version of connect or updating the connect
// version compiled into your binary.
const _ = connect.IsAtLeastVersion1_13_0

const (
	// PaymentsAPIName is the fully-qualified name of the PaymentsAPI service.
	PaymentsAPIName = "billing.paymentsapi.v1.PaymentsAPI"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// PaymentsAPIGetPaymentMethodProcedure is the fully-qualified name of the PaymentsAPI's
	// GetPaymentMethod RPC.
	PaymentsAPIGetPaymentMethodProcedure = "/billing.paymentsapi.v1.PaymentsAPI/GetPaymentMethod"
)

// These variables are the protoreflect.Descriptor objects for the RPCs defined in this package.
var (
	paymentsAPIServiceDescriptor                = v1.File_billing_paymentsapi_v1_api_proto.Services().ByName("PaymentsAPI")
	paymentsAPIGetPaymentMethodMethodDescriptor = paymentsAPIServiceDescriptor.Methods().ByName("GetPaymentMethod")
)

// PaymentsAPIClient is a client for the billing.paymentsapi.v1.PaymentsAPI service.
type PaymentsAPIClient interface {
	GetPaymentMethod(context.Context, *connect.Request[v1.GetPaymentMethodRequest]) (*connect.Response[v1.GetPaymentMethodResponse], error)
}

// NewPaymentsAPIClient constructs a client for the billing.paymentsapi.v1.PaymentsAPI service. By
// default, it uses the Connect protocol with the binary Protobuf Codec, asks for gzipped responses,
// and sends uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the
// connect.WithGRPC() or connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewPaymentsAPIClient(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) PaymentsAPIClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &paymentsAPIClient{
		getPaymentMethod: connect.NewClient[v1.GetPaymentMethodRequest, v1.GetPaymentMethodResponse](
			httpClient,
			baseURL+PaymentsAPIGetPaymentMethodProcedure,
			connect.WithSchema(paymentsAPIGetPaymentMethodMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
	}
}

// paymentsAPIClient implements PaymentsAPIClient.
type paymentsAPIClient struct {
	getPaymentMethod *connect.Client[v1.GetPaymentMethodRequest, v1.GetPaymentMethodResponse]
}

// GetPaymentMethod calls billing.paymentsapi.v1.PaymentsAPI.GetPaymentMethod.
func (c *paymentsAPIClient) GetPaymentMethod(ctx context.Context, req *connect.Request[v1.GetPaymentMethodRequest]) (*connect.Response[v1.GetPaymentMethodResponse], error) {
	return c.getPaymentMethod.CallUnary(ctx, req)
}

// PaymentsAPIHandler is an implementation of the billing.paymentsapi.v1.PaymentsAPI service.
type PaymentsAPIHandler interface {
	GetPaymentMethod(context.Context, *connect.Request[v1.GetPaymentMethodRequest]) (*connect.Response[v1.GetPaymentMethodResponse], error)
}

// NewPaymentsAPIHandler builds an HTTP handler from the service implementation. It returns the path
// on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewPaymentsAPIHandler(svc PaymentsAPIHandler, opts ...connect.HandlerOption) (string, http.Handler) {
	paymentsAPIGetPaymentMethodHandler := connect.NewUnaryHandler(
		PaymentsAPIGetPaymentMethodProcedure,
		svc.GetPaymentMethod,
		connect.WithSchema(paymentsAPIGetPaymentMethodMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	return "/billing.paymentsapi.v1.PaymentsAPI/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case PaymentsAPIGetPaymentMethodProcedure:
			paymentsAPIGetPaymentMethodHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedPaymentsAPIHandler returns CodeUnimplemented from all methods.
type UnimplementedPaymentsAPIHandler struct{}

func (UnimplementedPaymentsAPIHandler) GetPaymentMethod(context.Context, *connect.Request[v1.GetPaymentMethodRequest]) (*connect.Response[v1.GetPaymentMethodResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("billing.paymentsapi.v1.PaymentsAPI.GetPaymentMethod is not implemented"))
}
