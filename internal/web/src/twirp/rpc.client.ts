// @generated by protobuf-ts 2.9.3 with parameter generate_dependencies
// @generated from protobuf file "rpc.proto" (syntax proto3)
// tslint:disable
import { Admin } from "./rpc";
import type { DeleteDeviceReq } from "./rpc";
import type { UpdateDeviceReq } from "./rpc";
import type { GetDeviceResp } from "./rpc";
import type { GetDeviceReq } from "./rpc";
import type { CreateDeviceResp } from "./rpc";
import type { CreateDeviceReq } from "./rpc";
import type { SetGroupDisableReq } from "./rpc";
import type { DeleteGroupReq } from "./rpc";
import type { UpdateGroupReq } from "./rpc";
import type { GetGroupResp } from "./rpc";
import type { GetGroupReq } from "./rpc";
import type { CreateGroupResp } from "./rpc";
import type { CreateGroupReq } from "./rpc";
import type { SetUserAdminReq } from "./rpc";
import type { SetUserDisableReq } from "./rpc";
import type { GetAdminDevicesPageResp } from "./rpc";
import type { GetAdminDevicesPageReq } from "./rpc";
import type { GetAdminUsersPageResp } from "./rpc";
import type { GetAdminUsersPageReq } from "./rpc";
import type { GetAdminGroupIDPageResp } from "./rpc";
import type { GetAdminGroupIDPageReq } from "./rpc";
import type { GetAdminGroupsPageResp } from "./rpc";
import type { GetAdminGroupsPageReq } from "./rpc";
import { User } from "./rpc";
import type { RevokeAllMySessionsReq } from "./rpc";
import type { RevokeMySessionReq } from "./rpc";
import type { UpdateMyPasswordReq } from "./rpc";
import type { UpdateMyUsernameReq } from "./rpc";
import type { GetProfilePageResp } from "./rpc";
import type { GetHomePageResp } from "./rpc";
import type { Empty } from "./google/protobuf/empty";
import { Auth } from "./rpc";
import type { ForgotPasswordResp } from "./rpc";
import type { ForgotPasswordReq } from "./rpc";
import type { SignUpResp } from "./rpc";
import type { SignUpReq } from "./rpc";
import type { RpcTransport } from "@protobuf-ts/runtime-rpc";
import type { ServiceInfo } from "@protobuf-ts/runtime-rpc";
import { HelloWorld } from "./rpc";
import { stackIntercept } from "@protobuf-ts/runtime-rpc";
import type { HelloResp } from "./rpc";
import type { HelloReq } from "./rpc";
import type { UnaryCall } from "@protobuf-ts/runtime-rpc";
import type { RpcOptions } from "@protobuf-ts/runtime-rpc";
// ---------- HelloWorld

/**
 * @generated from protobuf service HelloWorld
 */
export interface IHelloWorldClient {
    /**
     * @generated from protobuf rpc: Hello(HelloReq) returns (HelloResp);
     */
    hello(input: HelloReq, options?: RpcOptions): UnaryCall<HelloReq, HelloResp>;
}
// ---------- HelloWorld

/**
 * @generated from protobuf service HelloWorld
 */
export class HelloWorldClient implements IHelloWorldClient, ServiceInfo {
    typeName = HelloWorld.typeName;
    methods = HelloWorld.methods;
    options = HelloWorld.options;
    constructor(private readonly _transport: RpcTransport) {
    }
    /**
     * @generated from protobuf rpc: Hello(HelloReq) returns (HelloResp);
     */
    hello(input: HelloReq, options?: RpcOptions): UnaryCall<HelloReq, HelloResp> {
        const method = this.methods[0], opt = this._transport.mergeOptions(options);
        return stackIntercept<HelloReq, HelloResp>("unary", this._transport, method, opt, input);
    }
}
// ---------- Auth

/**
 * @generated from protobuf service Auth
 */
export interface IAuthClient {
    /**
     * @generated from protobuf rpc: SignUp(SignUpReq) returns (SignUpResp);
     */
    signUp(input: SignUpReq, options?: RpcOptions): UnaryCall<SignUpReq, SignUpResp>;
    /**
     * @generated from protobuf rpc: ForgotPassword(ForgotPasswordReq) returns (ForgotPasswordResp);
     */
    forgotPassword(input: ForgotPasswordReq, options?: RpcOptions): UnaryCall<ForgotPasswordReq, ForgotPasswordResp>;
}
// ---------- Auth

/**
 * @generated from protobuf service Auth
 */
export class AuthClient implements IAuthClient, ServiceInfo {
    typeName = Auth.typeName;
    methods = Auth.methods;
    options = Auth.options;
    constructor(private readonly _transport: RpcTransport) {
    }
    /**
     * @generated from protobuf rpc: SignUp(SignUpReq) returns (SignUpResp);
     */
    signUp(input: SignUpReq, options?: RpcOptions): UnaryCall<SignUpReq, SignUpResp> {
        const method = this.methods[0], opt = this._transport.mergeOptions(options);
        return stackIntercept<SignUpReq, SignUpResp>("unary", this._transport, method, opt, input);
    }
    /**
     * @generated from protobuf rpc: ForgotPassword(ForgotPasswordReq) returns (ForgotPasswordResp);
     */
    forgotPassword(input: ForgotPasswordReq, options?: RpcOptions): UnaryCall<ForgotPasswordReq, ForgotPasswordResp> {
        const method = this.methods[1], opt = this._transport.mergeOptions(options);
        return stackIntercept<ForgotPasswordReq, ForgotPasswordResp>("unary", this._transport, method, opt, input);
    }
}
// ---------- User

/**
 * @generated from protobuf service User
 */
export interface IUserClient {
    /**
     * Pages
     *
     * @generated from protobuf rpc: GetHomePage(google.protobuf.Empty) returns (GetHomePageResp);
     */
    getHomePage(input: Empty, options?: RpcOptions): UnaryCall<Empty, GetHomePageResp>;
    /**
     * @generated from protobuf rpc: GetProfilePage(google.protobuf.Empty) returns (GetProfilePageResp);
     */
    getProfilePage(input: Empty, options?: RpcOptions): UnaryCall<Empty, GetProfilePageResp>;
    /**
     * User
     *
     * @generated from protobuf rpc: UpdateMyUsername(UpdateMyUsernameReq) returns (google.protobuf.Empty);
     */
    updateMyUsername(input: UpdateMyUsernameReq, options?: RpcOptions): UnaryCall<UpdateMyUsernameReq, Empty>;
    /**
     * @generated from protobuf rpc: UpdateMyPassword(UpdateMyPasswordReq) returns (google.protobuf.Empty);
     */
    updateMyPassword(input: UpdateMyPasswordReq, options?: RpcOptions): UnaryCall<UpdateMyPasswordReq, Empty>;
    /**
     * @generated from protobuf rpc: RevokeMySession(RevokeMySessionReq) returns (google.protobuf.Empty);
     */
    revokeMySession(input: RevokeMySessionReq, options?: RpcOptions): UnaryCall<RevokeMySessionReq, Empty>;
    /**
     * @generated from protobuf rpc: RevokeAllMySessions(RevokeAllMySessionsReq) returns (google.protobuf.Empty);
     */
    revokeAllMySessions(input: RevokeAllMySessionsReq, options?: RpcOptions): UnaryCall<RevokeAllMySessionsReq, Empty>;
}
// ---------- User

/**
 * @generated from protobuf service User
 */
export class UserClient implements IUserClient, ServiceInfo {
    typeName = User.typeName;
    methods = User.methods;
    options = User.options;
    constructor(private readonly _transport: RpcTransport) {
    }
    /**
     * Pages
     *
     * @generated from protobuf rpc: GetHomePage(google.protobuf.Empty) returns (GetHomePageResp);
     */
    getHomePage(input: Empty, options?: RpcOptions): UnaryCall<Empty, GetHomePageResp> {
        const method = this.methods[0], opt = this._transport.mergeOptions(options);
        return stackIntercept<Empty, GetHomePageResp>("unary", this._transport, method, opt, input);
    }
    /**
     * @generated from protobuf rpc: GetProfilePage(google.protobuf.Empty) returns (GetProfilePageResp);
     */
    getProfilePage(input: Empty, options?: RpcOptions): UnaryCall<Empty, GetProfilePageResp> {
        const method = this.methods[1], opt = this._transport.mergeOptions(options);
        return stackIntercept<Empty, GetProfilePageResp>("unary", this._transport, method, opt, input);
    }
    /**
     * User
     *
     * @generated from protobuf rpc: UpdateMyUsername(UpdateMyUsernameReq) returns (google.protobuf.Empty);
     */
    updateMyUsername(input: UpdateMyUsernameReq, options?: RpcOptions): UnaryCall<UpdateMyUsernameReq, Empty> {
        const method = this.methods[2], opt = this._transport.mergeOptions(options);
        return stackIntercept<UpdateMyUsernameReq, Empty>("unary", this._transport, method, opt, input);
    }
    /**
     * @generated from protobuf rpc: UpdateMyPassword(UpdateMyPasswordReq) returns (google.protobuf.Empty);
     */
    updateMyPassword(input: UpdateMyPasswordReq, options?: RpcOptions): UnaryCall<UpdateMyPasswordReq, Empty> {
        const method = this.methods[3], opt = this._transport.mergeOptions(options);
        return stackIntercept<UpdateMyPasswordReq, Empty>("unary", this._transport, method, opt, input);
    }
    /**
     * @generated from protobuf rpc: RevokeMySession(RevokeMySessionReq) returns (google.protobuf.Empty);
     */
    revokeMySession(input: RevokeMySessionReq, options?: RpcOptions): UnaryCall<RevokeMySessionReq, Empty> {
        const method = this.methods[4], opt = this._transport.mergeOptions(options);
        return stackIntercept<RevokeMySessionReq, Empty>("unary", this._transport, method, opt, input);
    }
    /**
     * @generated from protobuf rpc: RevokeAllMySessions(RevokeAllMySessionsReq) returns (google.protobuf.Empty);
     */
    revokeAllMySessions(input: RevokeAllMySessionsReq, options?: RpcOptions): UnaryCall<RevokeAllMySessionsReq, Empty> {
        const method = this.methods[5], opt = this._transport.mergeOptions(options);
        return stackIntercept<RevokeAllMySessionsReq, Empty>("unary", this._transport, method, opt, input);
    }
}
// ---------- Admin

/**
 * @generated from protobuf service Admin
 */
export interface IAdminClient {
    /**
     * Pages
     *
     * @generated from protobuf rpc: GetAdminGroupsPage(GetAdminGroupsPageReq) returns (GetAdminGroupsPageResp);
     */
    getAdminGroupsPage(input: GetAdminGroupsPageReq, options?: RpcOptions): UnaryCall<GetAdminGroupsPageReq, GetAdminGroupsPageResp>;
    /**
     * @generated from protobuf rpc: GetAdminGroupIDPage(GetAdminGroupIDPageReq) returns (GetAdminGroupIDPageResp);
     */
    getAdminGroupIDPage(input: GetAdminGroupIDPageReq, options?: RpcOptions): UnaryCall<GetAdminGroupIDPageReq, GetAdminGroupIDPageResp>;
    /**
     * @generated from protobuf rpc: GetAdminUsersPage(GetAdminUsersPageReq) returns (GetAdminUsersPageResp);
     */
    getAdminUsersPage(input: GetAdminUsersPageReq, options?: RpcOptions): UnaryCall<GetAdminUsersPageReq, GetAdminUsersPageResp>;
    /**
     * @generated from protobuf rpc: GetAdminDevicesPage(GetAdminDevicesPageReq) returns (GetAdminDevicesPageResp);
     */
    getAdminDevicesPage(input: GetAdminDevicesPageReq, options?: RpcOptions): UnaryCall<GetAdminDevicesPageReq, GetAdminDevicesPageResp>;
    /**
     * User
     *
     * @generated from protobuf rpc: SetUserDisable(SetUserDisableReq) returns (google.protobuf.Empty);
     */
    setUserDisable(input: SetUserDisableReq, options?: RpcOptions): UnaryCall<SetUserDisableReq, Empty>;
    /**
     * @generated from protobuf rpc: SetUserAdmin(SetUserAdminReq) returns (google.protobuf.Empty);
     */
    setUserAdmin(input: SetUserAdminReq, options?: RpcOptions): UnaryCall<SetUserAdminReq, Empty>;
    /**
     * Group
     *
     * @generated from protobuf rpc: CreateGroup(CreateGroupReq) returns (CreateGroupResp);
     */
    createGroup(input: CreateGroupReq, options?: RpcOptions): UnaryCall<CreateGroupReq, CreateGroupResp>;
    /**
     * @generated from protobuf rpc: GetGroup(GetGroupReq) returns (GetGroupResp);
     */
    getGroup(input: GetGroupReq, options?: RpcOptions): UnaryCall<GetGroupReq, GetGroupResp>;
    /**
     * @generated from protobuf rpc: UpdateGroup(UpdateGroupReq) returns (google.protobuf.Empty);
     */
    updateGroup(input: UpdateGroupReq, options?: RpcOptions): UnaryCall<UpdateGroupReq, Empty>;
    /**
     * @generated from protobuf rpc: DeleteGroup(DeleteGroupReq) returns (google.protobuf.Empty);
     */
    deleteGroup(input: DeleteGroupReq, options?: RpcOptions): UnaryCall<DeleteGroupReq, Empty>;
    /**
     * @generated from protobuf rpc: SetGroupDisable(SetGroupDisableReq) returns (google.protobuf.Empty);
     */
    setGroupDisable(input: SetGroupDisableReq, options?: RpcOptions): UnaryCall<SetGroupDisableReq, Empty>;
    /**
     * Device
     *
     * @generated from protobuf rpc: CreateDevice(CreateDeviceReq) returns (CreateDeviceResp);
     */
    createDevice(input: CreateDeviceReq, options?: RpcOptions): UnaryCall<CreateDeviceReq, CreateDeviceResp>;
    /**
     * @generated from protobuf rpc: GetDevice(GetDeviceReq) returns (GetDeviceResp);
     */
    getDevice(input: GetDeviceReq, options?: RpcOptions): UnaryCall<GetDeviceReq, GetDeviceResp>;
    /**
     * @generated from protobuf rpc: UpdateDevice(UpdateDeviceReq) returns (google.protobuf.Empty);
     */
    updateDevice(input: UpdateDeviceReq, options?: RpcOptions): UnaryCall<UpdateDeviceReq, Empty>;
    /**
     * @generated from protobuf rpc: DeleteDevice(DeleteDeviceReq) returns (google.protobuf.Empty);
     */
    deleteDevice(input: DeleteDeviceReq, options?: RpcOptions): UnaryCall<DeleteDeviceReq, Empty>;
}
// ---------- Admin

/**
 * @generated from protobuf service Admin
 */
export class AdminClient implements IAdminClient, ServiceInfo {
    typeName = Admin.typeName;
    methods = Admin.methods;
    options = Admin.options;
    constructor(private readonly _transport: RpcTransport) {
    }
    /**
     * Pages
     *
     * @generated from protobuf rpc: GetAdminGroupsPage(GetAdminGroupsPageReq) returns (GetAdminGroupsPageResp);
     */
    getAdminGroupsPage(input: GetAdminGroupsPageReq, options?: RpcOptions): UnaryCall<GetAdminGroupsPageReq, GetAdminGroupsPageResp> {
        const method = this.methods[0], opt = this._transport.mergeOptions(options);
        return stackIntercept<GetAdminGroupsPageReq, GetAdminGroupsPageResp>("unary", this._transport, method, opt, input);
    }
    /**
     * @generated from protobuf rpc: GetAdminGroupIDPage(GetAdminGroupIDPageReq) returns (GetAdminGroupIDPageResp);
     */
    getAdminGroupIDPage(input: GetAdminGroupIDPageReq, options?: RpcOptions): UnaryCall<GetAdminGroupIDPageReq, GetAdminGroupIDPageResp> {
        const method = this.methods[1], opt = this._transport.mergeOptions(options);
        return stackIntercept<GetAdminGroupIDPageReq, GetAdminGroupIDPageResp>("unary", this._transport, method, opt, input);
    }
    /**
     * @generated from protobuf rpc: GetAdminUsersPage(GetAdminUsersPageReq) returns (GetAdminUsersPageResp);
     */
    getAdminUsersPage(input: GetAdminUsersPageReq, options?: RpcOptions): UnaryCall<GetAdminUsersPageReq, GetAdminUsersPageResp> {
        const method = this.methods[2], opt = this._transport.mergeOptions(options);
        return stackIntercept<GetAdminUsersPageReq, GetAdminUsersPageResp>("unary", this._transport, method, opt, input);
    }
    /**
     * @generated from protobuf rpc: GetAdminDevicesPage(GetAdminDevicesPageReq) returns (GetAdminDevicesPageResp);
     */
    getAdminDevicesPage(input: GetAdminDevicesPageReq, options?: RpcOptions): UnaryCall<GetAdminDevicesPageReq, GetAdminDevicesPageResp> {
        const method = this.methods[3], opt = this._transport.mergeOptions(options);
        return stackIntercept<GetAdminDevicesPageReq, GetAdminDevicesPageResp>("unary", this._transport, method, opt, input);
    }
    /**
     * User
     *
     * @generated from protobuf rpc: SetUserDisable(SetUserDisableReq) returns (google.protobuf.Empty);
     */
    setUserDisable(input: SetUserDisableReq, options?: RpcOptions): UnaryCall<SetUserDisableReq, Empty> {
        const method = this.methods[4], opt = this._transport.mergeOptions(options);
        return stackIntercept<SetUserDisableReq, Empty>("unary", this._transport, method, opt, input);
    }
    /**
     * @generated from protobuf rpc: SetUserAdmin(SetUserAdminReq) returns (google.protobuf.Empty);
     */
    setUserAdmin(input: SetUserAdminReq, options?: RpcOptions): UnaryCall<SetUserAdminReq, Empty> {
        const method = this.methods[5], opt = this._transport.mergeOptions(options);
        return stackIntercept<SetUserAdminReq, Empty>("unary", this._transport, method, opt, input);
    }
    /**
     * Group
     *
     * @generated from protobuf rpc: CreateGroup(CreateGroupReq) returns (CreateGroupResp);
     */
    createGroup(input: CreateGroupReq, options?: RpcOptions): UnaryCall<CreateGroupReq, CreateGroupResp> {
        const method = this.methods[6], opt = this._transport.mergeOptions(options);
        return stackIntercept<CreateGroupReq, CreateGroupResp>("unary", this._transport, method, opt, input);
    }
    /**
     * @generated from protobuf rpc: GetGroup(GetGroupReq) returns (GetGroupResp);
     */
    getGroup(input: GetGroupReq, options?: RpcOptions): UnaryCall<GetGroupReq, GetGroupResp> {
        const method = this.methods[7], opt = this._transport.mergeOptions(options);
        return stackIntercept<GetGroupReq, GetGroupResp>("unary", this._transport, method, opt, input);
    }
    /**
     * @generated from protobuf rpc: UpdateGroup(UpdateGroupReq) returns (google.protobuf.Empty);
     */
    updateGroup(input: UpdateGroupReq, options?: RpcOptions): UnaryCall<UpdateGroupReq, Empty> {
        const method = this.methods[8], opt = this._transport.mergeOptions(options);
        return stackIntercept<UpdateGroupReq, Empty>("unary", this._transport, method, opt, input);
    }
    /**
     * @generated from protobuf rpc: DeleteGroup(DeleteGroupReq) returns (google.protobuf.Empty);
     */
    deleteGroup(input: DeleteGroupReq, options?: RpcOptions): UnaryCall<DeleteGroupReq, Empty> {
        const method = this.methods[9], opt = this._transport.mergeOptions(options);
        return stackIntercept<DeleteGroupReq, Empty>("unary", this._transport, method, opt, input);
    }
    /**
     * @generated from protobuf rpc: SetGroupDisable(SetGroupDisableReq) returns (google.protobuf.Empty);
     */
    setGroupDisable(input: SetGroupDisableReq, options?: RpcOptions): UnaryCall<SetGroupDisableReq, Empty> {
        const method = this.methods[10], opt = this._transport.mergeOptions(options);
        return stackIntercept<SetGroupDisableReq, Empty>("unary", this._transport, method, opt, input);
    }
    /**
     * Device
     *
     * @generated from protobuf rpc: CreateDevice(CreateDeviceReq) returns (CreateDeviceResp);
     */
    createDevice(input: CreateDeviceReq, options?: RpcOptions): UnaryCall<CreateDeviceReq, CreateDeviceResp> {
        const method = this.methods[11], opt = this._transport.mergeOptions(options);
        return stackIntercept<CreateDeviceReq, CreateDeviceResp>("unary", this._transport, method, opt, input);
    }
    /**
     * @generated from protobuf rpc: GetDevice(GetDeviceReq) returns (GetDeviceResp);
     */
    getDevice(input: GetDeviceReq, options?: RpcOptions): UnaryCall<GetDeviceReq, GetDeviceResp> {
        const method = this.methods[12], opt = this._transport.mergeOptions(options);
        return stackIntercept<GetDeviceReq, GetDeviceResp>("unary", this._transport, method, opt, input);
    }
    /**
     * @generated from protobuf rpc: UpdateDevice(UpdateDeviceReq) returns (google.protobuf.Empty);
     */
    updateDevice(input: UpdateDeviceReq, options?: RpcOptions): UnaryCall<UpdateDeviceReq, Empty> {
        const method = this.methods[13], opt = this._transport.mergeOptions(options);
        return stackIntercept<UpdateDeviceReq, Empty>("unary", this._transport, method, opt, input);
    }
    /**
     * @generated from protobuf rpc: DeleteDevice(DeleteDeviceReq) returns (google.protobuf.Empty);
     */
    deleteDevice(input: DeleteDeviceReq, options?: RpcOptions): UnaryCall<DeleteDeviceReq, Empty> {
        const method = this.methods[14], opt = this._transport.mergeOptions(options);
        return stackIntercept<DeleteDeviceReq, Empty>("unary", this._transport, method, opt, input);
    }
}
