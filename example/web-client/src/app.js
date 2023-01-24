import { createPromiseClient, createConnectTransport } from "@bufbuild/connect-web";

import { StreamerService } from "./gen/pkg/streamer/v1/streamer_connectweb.js"
import { Point } from "./gen/pkg/streamer/v1/streamer_pb.js";


const main = async () => {
    const logger = (next) => async (req) => {
        console.log(`sending message to ${req.url}`);
        return await next(req);
    };
    
    const transport = createConnectTransport({
        baseUrl: "http://localhost:8080",
        interceptors: [logger],
    });

    let point = new Point({latitude:44, longitude:55})

    const streamService = createPromiseClient(StreamerService, transport);

    // stream is only one direction. not client -> server
    var status = await streamService.streamPoint(point)
    console.log("status:", status)
}

main()