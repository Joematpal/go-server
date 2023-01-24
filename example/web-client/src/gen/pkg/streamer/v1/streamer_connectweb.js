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
export const StreamerService = {
  typeName: "streamer.v1.StreamerService",
  methods: {
    /**
     * @generated from rpc streamer.v1.StreamerService.StreamPoint
     */
    streamPoint: {
      name: "StreamPoint",
      I: Point,
      O: Status,
      kind: MethodKind.ClientStreaming,
    },
  }
};

