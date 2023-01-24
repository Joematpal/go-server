// @generated by protoc-gen-es v1.0.0 with parameter "target=ts"
// @generated from file pkg/streamer/v1/streamer.proto (package streamer.v1, syntax proto3)
/* eslint-disable */
// @ts-nocheck

import type { BinaryReadOptions, FieldList, JsonReadOptions, JsonValue, PartialMessage, PlainMessage } from "@bufbuild/protobuf";
import { Message, proto3 } from "@bufbuild/protobuf";

/**
 * Points are represented as latitude-longitude pairs in the E7 representation
 * (degrees multiplied by 10**7 and rounded to the nearest integer).
 * Latitudes should be in the range +/- 90 degrees and longitude should be in
 * the range +/- 180 degrees (inclusive).
 *
 * @generated from message streamer.v1.Point
 */
export class Point extends Message<Point> {
  /**
   * @generated from field: int32 latitude = 1;
   */
  latitude = 0;

  /**
   * @generated from field: int32 longitude = 2;
   */
  longitude = 0;

  constructor(data?: PartialMessage<Point>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime = proto3;
  static readonly typeName = "streamer.v1.Point";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "latitude", kind: "scalar", T: 5 /* ScalarType.INT32 */ },
    { no: 2, name: "longitude", kind: "scalar", T: 5 /* ScalarType.INT32 */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): Point {
    return new Point().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): Point {
    return new Point().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): Point {
    return new Point().fromJsonString(jsonString, options);
  }

  static equals(a: Point | PlainMessage<Point> | undefined, b: Point | PlainMessage<Point> | undefined): boolean {
    return proto3.util.equals(Point, a, b);
  }
}

/**
 * @generated from message streamer.v1.Status
 */
export class Status extends Message<Status> {
  /**
   * @generated from field: int32 status = 1;
   */
  status = 0;

  /**
   * @generated from field: string message = 2;
   */
  message = "";

  constructor(data?: PartialMessage<Status>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime = proto3;
  static readonly typeName = "streamer.v1.Status";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "status", kind: "scalar", T: 5 /* ScalarType.INT32 */ },
    { no: 2, name: "message", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): Status {
    return new Status().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): Status {
    return new Status().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): Status {
    return new Status().fromJsonString(jsonString, options);
  }

  static equals(a: Status | PlainMessage<Status> | undefined, b: Status | PlainMessage<Status> | undefined): boolean {
    return proto3.util.equals(Status, a, b);
  }
}

