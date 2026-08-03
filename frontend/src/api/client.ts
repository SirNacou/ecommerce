import { createClient } from "@connectrpc/connect"
import { createConnectTransport } from "@connectrpc/connect-web"
import { UserService } from "../gen/v1/user_pb"

// 1. Create transport pointing to your APISIX Gateway (Port 9080)
export const transport = createConnectTransport({
  baseUrl: "http://localhost:9080",
})

// 2. Instantiate type-safe client for UserService
export const userServiceClient = createClient(UserService, transport)
userServiceClient.register({ email: "test", password: "test" })