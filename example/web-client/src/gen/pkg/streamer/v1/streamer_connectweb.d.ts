// @generated by protoc-gen-connect-web v0.6.0 with parameter "target=js+dts"
// @generated from file pkg/streamer/v1/streamer.proto (package streamer.v1, syntax proto3)
/* eslint-disable */
// @ts-nocheck

import { Point, Status } from "./streamer_pb.js";
import { MethodKind } from "@bufbuild/protobuf";

/**
 * Interface exported by the server.
 *
 * @generated from service streamer.v1.StreamerService
 */
export declare const StreamerService: {
  readonly typeName: "streamer.v1.StreamerService",
  readonly methods: {
    /**
     * @generated from rpc streamer.v1.StreamerService.StreamPoint
     */
    readonly streamPoint: {
      readonly name: "StreamPoint",
      readonly I: typeof Point,
      readonly O: typeof Status,
      readonly kind: MethodKind.ClientStreaming,
    },
  }
};
