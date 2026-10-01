package main

import (
	"log"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/dynamicpb"
)

const (
	entitlementsAddress = "127.0.0.1:8083"
	entitlementsRPC     = "/entitlements.v1.EntitlementsService/CheckEntitlement"
)

func main() {
	method := entitlementMethodDescriptor()

	http.HandleFunc("/check", func(w http.ResponseWriter, r *http.Request) {
		connection, err := grpc.NewClient(
			entitlementsAddress,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			http.Error(w, "dependency request failed", http.StatusBadGateway)
			return
		}
		defer connection.Close()

		request := dynamicpb.NewMessage(method.Input())
		request.Set(
			method.Input().Fields().ByName("customer_id"),
			protoreflect.ValueOfString("customer-2"),
		)

		response := dynamicpb.NewMessage(method.Output())

		const testExecutionIDKey = "e2engine-test-execution-id"

		callCtx := metadata.AppendToOutgoingContext(
			r.Context(),
			testExecutionIDKey,
			r.Header.Get("E2Engine-Test-Execution-ID"),
		)

		if err := connection.Invoke(
			callCtx,
			entitlementsRPC,
			request,
			response,
		); err != nil {
			http.Error(
				w,
				"dependency returned non-ok status",
				http.StatusBadGateway,
			)
			return
		}

		entitled := response.Get(
			method.Output().Fields().ByName("entitled"),
		).Bool()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if entitled {
			_, _ = w.Write([]byte(`{"entitled":true}`))
			return
		}

		_, _ = w.Write([]byte(`{"entitled":false}`))
	})

	if err := http.ListenAndServe("127.0.0.1:9000", nil); err != nil {
		log.Fatal(err)
	}
}

func entitlementMethodDescriptor() protoreflect.MethodDescriptor {
	file, err := protodesc.NewFile(
		&descriptorpb.FileDescriptorProto{
			Name:    proto.String("entitlements.proto"),
			Package: proto.String("entitlements.v1"),
			Syntax:  proto.String("proto3"),
			MessageType: []*descriptorpb.DescriptorProto{
				{
					Name: proto.String("CheckEntitlementRequest"),
					Field: []*descriptorpb.FieldDescriptorProto{
						{
							Name:   proto.String("customer_id"),
							Number: proto.Int32(1),
							Label: descriptorpb.
								FieldDescriptorProto_LABEL_OPTIONAL.
								Enum(),
							Type: descriptorpb.
								FieldDescriptorProto_TYPE_STRING.
								Enum(),
						},
					},
				},
				{
					Name: proto.String("CheckEntitlementResponse"),
					Field: []*descriptorpb.FieldDescriptorProto{
						{
							Name:   proto.String("entitled"),
							Number: proto.Int32(1),
							Label: descriptorpb.
								FieldDescriptorProto_LABEL_OPTIONAL.
								Enum(),
							Type: descriptorpb.
								FieldDescriptorProto_TYPE_BOOL.
								Enum(),
						},
					},
				},
			},
			Service: []*descriptorpb.ServiceDescriptorProto{
				{
					Name: proto.String("EntitlementsService"),
					Method: []*descriptorpb.MethodDescriptorProto{
						{
							Name: proto.String(
								"CheckEntitlement",
							),
							InputType: proto.String(
								".entitlements.v1.CheckEntitlementRequest",
							),
							OutputType: proto.String(
								".entitlements.v1.CheckEntitlementResponse",
							),
						},
					},
				},
			},
		},
		nil,
	)
	if err != nil {
		panic(err)
	}

	service := file.Services().ByName("EntitlementsService")
	return service.Methods().ByName("CheckEntitlement")
}
