// @generated
/// Generated client implementations.
pub mod service_client {
    #![allow(
        unused_variables,
        dead_code,
        missing_docs,
        clippy::wildcard_imports,
        clippy::let_unit_value,
    )]
    use tonic::codegen::*;
    use tonic::codegen::http::Uri;
    #[derive(Debug, Clone)]
    pub struct ServiceClient<T> {
        inner: tonic::client::Grpc<T>,
    }
    impl ServiceClient<tonic::transport::Channel> {
        /// Attempt to create a new client by connecting to a given endpoint.
        pub async fn connect<D>(dst: D) -> Result<Self, tonic::transport::Error>
        where
            D: TryInto<tonic::transport::Endpoint>,
            D::Error: Into<StdError>,
        {
            let conn = tonic::transport::Endpoint::new(dst)?.connect().await?;
            Ok(Self::new(conn))
        }
    }
    impl<T> ServiceClient<T>
    where
        T: tonic::client::GrpcService<tonic::body::BoxBody>,
        T::Error: Into<StdError>,
        T::ResponseBody: Body<Data = Bytes> + std::marker::Send + 'static,
        <T::ResponseBody as Body>::Error: Into<StdError> + std::marker::Send,
    {
        pub fn new(inner: T) -> Self {
            let inner = tonic::client::Grpc::new(inner);
            Self { inner }
        }
        pub fn with_origin(inner: T, origin: Uri) -> Self {
            let inner = tonic::client::Grpc::with_origin(inner, origin);
            Self { inner }
        }
        pub fn with_interceptor<F>(
            inner: T,
            interceptor: F,
        ) -> ServiceClient<InterceptedService<T, F>>
        where
            F: tonic::service::Interceptor,
            T::ResponseBody: Default,
            T: tonic::codegen::Service<
                http::Request<tonic::body::BoxBody>,
                Response = http::Response<
                    <T as tonic::client::GrpcService<tonic::body::BoxBody>>::ResponseBody,
                >,
            >,
            <T as tonic::codegen::Service<
                http::Request<tonic::body::BoxBody>,
            >>::Error: Into<StdError> + std::marker::Send + std::marker::Sync,
        {
            ServiceClient::new(InterceptedService::new(inner, interceptor))
        }
        /// Compress requests with the given encoding.
        ///
        /// This requires the server to support it otherwise it might respond with an
        /// error.
        #[must_use]
        pub fn send_compressed(mut self, encoding: CompressionEncoding) -> Self {
            self.inner = self.inner.send_compressed(encoding);
            self
        }
        /// Enable decompressing responses.
        #[must_use]
        pub fn accept_compressed(mut self, encoding: CompressionEncoding) -> Self {
            self.inner = self.inner.accept_compressed(encoding);
            self
        }
        /// Limits the maximum size of a decoded message.
        ///
        /// Default: `4MB`
        #[must_use]
        pub fn max_decoding_message_size(mut self, limit: usize) -> Self {
            self.inner = self.inner.max_decoding_message_size(limit);
            self
        }
        /// Limits the maximum size of an encoded message.
        ///
        /// Default: `usize::MAX`
        #[must_use]
        pub fn max_encoding_message_size(mut self, limit: usize) -> Self {
            self.inner = self.inner.max_encoding_message_size(limit);
            self
        }
        pub async fn get_token_auth_info(
            &mut self,
            request: impl tonic::IntoRequest<super::TokenRequest>,
        ) -> std::result::Result<tonic::Response<super::TokenResponse>, tonic::Status> {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::unknown(
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/message.Service/GetTokenAuthInfo",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("message.Service", "GetTokenAuthInfo"));
            self.inner.unary(req, path, codec).await
        }
        pub async fn renew_token(
            &mut self,
            request: impl tonic::IntoRequest<super::TokenRequest>,
        ) -> std::result::Result<tonic::Response<super::StatusResponse>, tonic::Status> {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::unknown(
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/message.Service/RenewToken",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("message.Service", "RenewToken"));
            self.inner.unary(req, path, codec).await
        }
        pub async fn create_session(
            &mut self,
            request: impl tonic::IntoRequest<super::SessionCreateRequest>,
        ) -> std::result::Result<
            tonic::Response<super::SessionCreateResponse>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::unknown(
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/message.Service/CreateSession",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("message.Service", "CreateSession"));
            self.inner.unary(req, path, codec).await
        }
        pub async fn finish_session(
            &mut self,
            request: impl tonic::IntoRequest<super::SessionFinishRequest>,
        ) -> std::result::Result<
            tonic::Response<super::SessionFinishResp>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::unknown(
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/message.Service/FinishSession",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("message.Service", "FinishSession"));
            self.inner.unary(req, path, codec).await
        }
        pub async fn upload_replay_file(
            &mut self,
            request: impl tonic::IntoRequest<super::ReplayRequest>,
        ) -> std::result::Result<tonic::Response<super::ReplayResponse>, tonic::Status> {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::unknown(
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/message.Service/UploadReplayFile",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("message.Service", "UploadReplayFile"));
            self.inner.unary(req, path, codec).await
        }
        pub async fn upload_command(
            &mut self,
            request: impl tonic::IntoRequest<super::CommandRequest>,
        ) -> std::result::Result<
            tonic::Response<super::CommandResponse>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::unknown(
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/message.Service/UploadCommand",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("message.Service", "UploadCommand"));
            self.inner.unary(req, path, codec).await
        }
        pub async fn dispatch_task(
            &mut self,
            request: impl tonic::IntoStreamingRequest<
                Message = super::FinishedTaskRequest,
            >,
        ) -> std::result::Result<
            tonic::Response<tonic::codec::Streaming<super::TaskResponse>>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::unknown(
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/message.Service/DispatchTask",
            );
            let mut req = request.into_streaming_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("message.Service", "DispatchTask"));
            self.inner.streaming(req, path, codec).await
        }
        pub async fn scan_remain_replays(
            &mut self,
            request: impl tonic::IntoRequest<super::RemainReplayRequest>,
        ) -> std::result::Result<
            tonic::Response<super::RemainReplayResponse>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::unknown(
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/message.Service/ScanRemainReplays",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("message.Service", "ScanRemainReplays"));
            self.inner.unary(req, path, codec).await
        }
        pub async fn create_command_ticket(
            &mut self,
            request: impl tonic::IntoRequest<super::CommandConfirmRequest>,
        ) -> std::result::Result<
            tonic::Response<super::CommandConfirmResponse>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::unknown(
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/message.Service/CreateCommandTicket",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("message.Service", "CreateCommandTicket"));
            self.inner.unary(req, path, codec).await
        }
        pub async fn check_or_create_asset_login_ticket(
            &mut self,
            request: impl tonic::IntoRequest<super::AssetLoginTicketRequest>,
        ) -> std::result::Result<
            tonic::Response<super::AssetLoginTicketResponse>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::unknown(
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/message.Service/CheckOrCreateAssetLoginTicket",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(
                    GrpcMethod::new("message.Service", "CheckOrCreateAssetLoginTicket"),
                );
            self.inner.unary(req, path, codec).await
        }
        pub async fn check_ticket_state(
            &mut self,
            request: impl tonic::IntoRequest<super::TicketRequest>,
        ) -> std::result::Result<
            tonic::Response<super::TicketStateResponse>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::unknown(
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/message.Service/CheckTicketState",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("message.Service", "CheckTicketState"));
            self.inner.unary(req, path, codec).await
        }
        pub async fn cancel_ticket(
            &mut self,
            request: impl tonic::IntoRequest<super::TicketRequest>,
        ) -> std::result::Result<tonic::Response<super::StatusResponse>, tonic::Status> {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::unknown(
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/message.Service/CancelTicket",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("message.Service", "CancelTicket"));
            self.inner.unary(req, path, codec).await
        }
        pub async fn create_forward(
            &mut self,
            request: impl tonic::IntoRequest<super::ForwardRequest>,
        ) -> std::result::Result<
            tonic::Response<super::ForwardResponse>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::unknown(
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/message.Service/CreateForward",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("message.Service", "CreateForward"));
            self.inner.unary(req, path, codec).await
        }
        pub async fn delete_forward(
            &mut self,
            request: impl tonic::IntoRequest<super::ForwardDeleteRequest>,
        ) -> std::result::Result<tonic::Response<super::StatusResponse>, tonic::Status> {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::unknown(
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/message.Service/DeleteForward",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("message.Service", "DeleteForward"));
            self.inner.unary(req, path, codec).await
        }
        pub async fn get_public_setting(
            &mut self,
            request: impl tonic::IntoRequest<super::Empty>,
        ) -> std::result::Result<
            tonic::Response<super::PublicSettingResponse>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::unknown(
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/message.Service/GetPublicSetting",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("message.Service", "GetPublicSetting"));
            self.inner.unary(req, path, codec).await
        }
        pub async fn get_listen_ports(
            &mut self,
            request: impl tonic::IntoRequest<super::Empty>,
        ) -> std::result::Result<
            tonic::Response<super::ListenPortResponse>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::unknown(
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/message.Service/GetListenPorts",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("message.Service", "GetListenPorts"));
            self.inner.unary(req, path, codec).await
        }
        pub async fn get_port_info(
            &mut self,
            request: impl tonic::IntoRequest<super::PortInfoRequest>,
        ) -> std::result::Result<
            tonic::Response<super::PortInfoResponse>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::unknown(
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/message.Service/GetPortInfo",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("message.Service", "GetPortInfo"));
            self.inner.unary(req, path, codec).await
        }
        pub async fn handle_port_failure(
            &mut self,
            request: impl tonic::IntoRequest<super::PortFailureRequest>,
        ) -> std::result::Result<tonic::Response<super::StatusResponse>, tonic::Status> {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::unknown(
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/message.Service/HandlePortFailure",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("message.Service", "HandlePortFailure"));
            self.inner.unary(req, path, codec).await
        }
        pub async fn check_user_by_cookies(
            &mut self,
            request: impl tonic::IntoRequest<super::CookiesRequest>,
        ) -> std::result::Result<tonic::Response<super::UserResponse>, tonic::Status> {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::unknown(
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/message.Service/CheckUserByCookies",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("message.Service", "CheckUserByCookies"));
            self.inner.unary(req, path, codec).await
        }
        pub async fn record_session_lifecycle_log(
            &mut self,
            request: impl tonic::IntoRequest<super::SessionLifecycleLogRequest>,
        ) -> std::result::Result<tonic::Response<super::StatusResponse>, tonic::Status> {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::unknown(
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/message.Service/RecordSessionLifecycleLog",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("message.Service", "RecordSessionLifecycleLog"));
            self.inner.unary(req, path, codec).await
        }
        pub async fn face_recognition_callback(
            &mut self,
            request: impl tonic::IntoRequest<super::FaceRecognitionCallbackRequest>,
        ) -> std::result::Result<
            tonic::Response<super::FaceRecognitionCallbackResponse>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::unknown(
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/message.Service/FaceRecognitionCallback",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("message.Service", "FaceRecognitionCallback"));
            self.inner.unary(req, path, codec).await
        }
        pub async fn face_monitor_callback(
            &mut self,
            request: impl tonic::IntoRequest<super::FaceMonitorCallbackRequest>,
        ) -> std::result::Result<
            tonic::Response<super::FaceMonitorCallbackResponse>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::unknown(
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/message.Service/FaceMonitorCallback",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("message.Service", "FaceMonitorCallback"));
            self.inner.unary(req, path, codec).await
        }
        pub async fn join_face_monitor(
            &mut self,
            request: impl tonic::IntoRequest<super::JoinFaceMonitorRequest>,
        ) -> std::result::Result<
            tonic::Response<super::JoinFaceMonitorResponse>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::unknown(
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/message.Service/JoinFaceMonitor",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("message.Service", "JoinFaceMonitor"));
            self.inner.unary(req, path, codec).await
        }
    }
}
/// Generated server implementations.
pub mod service_server {
    #![allow(
        unused_variables,
        dead_code,
        missing_docs,
        clippy::wildcard_imports,
        clippy::let_unit_value,
    )]
    use tonic::codegen::*;
    /// Generated trait containing gRPC methods that should be implemented for use with ServiceServer.
    #[async_trait]
    pub trait Service: std::marker::Send + std::marker::Sync + 'static {
        async fn get_token_auth_info(
            &self,
            request: tonic::Request<super::TokenRequest>,
        ) -> std::result::Result<tonic::Response<super::TokenResponse>, tonic::Status>;
        async fn renew_token(
            &self,
            request: tonic::Request<super::TokenRequest>,
        ) -> std::result::Result<tonic::Response<super::StatusResponse>, tonic::Status>;
        async fn create_session(
            &self,
            request: tonic::Request<super::SessionCreateRequest>,
        ) -> std::result::Result<
            tonic::Response<super::SessionCreateResponse>,
            tonic::Status,
        >;
        async fn finish_session(
            &self,
            request: tonic::Request<super::SessionFinishRequest>,
        ) -> std::result::Result<
            tonic::Response<super::SessionFinishResp>,
            tonic::Status,
        >;
        async fn upload_replay_file(
            &self,
            request: tonic::Request<super::ReplayRequest>,
        ) -> std::result::Result<tonic::Response<super::ReplayResponse>, tonic::Status>;
        async fn upload_command(
            &self,
            request: tonic::Request<super::CommandRequest>,
        ) -> std::result::Result<tonic::Response<super::CommandResponse>, tonic::Status>;
        /// Server streaming response type for the DispatchTask method.
        type DispatchTaskStream: tonic::codegen::tokio_stream::Stream<
                Item = std::result::Result<super::TaskResponse, tonic::Status>,
            >
            + std::marker::Send
            + 'static;
        async fn dispatch_task(
            &self,
            request: tonic::Request<tonic::Streaming<super::FinishedTaskRequest>>,
        ) -> std::result::Result<
            tonic::Response<Self::DispatchTaskStream>,
            tonic::Status,
        >;
        async fn scan_remain_replays(
            &self,
            request: tonic::Request<super::RemainReplayRequest>,
        ) -> std::result::Result<
            tonic::Response<super::RemainReplayResponse>,
            tonic::Status,
        >;
        async fn create_command_ticket(
            &self,
            request: tonic::Request<super::CommandConfirmRequest>,
        ) -> std::result::Result<
            tonic::Response<super::CommandConfirmResponse>,
            tonic::Status,
        >;
        async fn check_or_create_asset_login_ticket(
            &self,
            request: tonic::Request<super::AssetLoginTicketRequest>,
        ) -> std::result::Result<
            tonic::Response<super::AssetLoginTicketResponse>,
            tonic::Status,
        >;
        async fn check_ticket_state(
            &self,
            request: tonic::Request<super::TicketRequest>,
        ) -> std::result::Result<
            tonic::Response<super::TicketStateResponse>,
            tonic::Status,
        >;
        async fn cancel_ticket(
            &self,
            request: tonic::Request<super::TicketRequest>,
        ) -> std::result::Result<tonic::Response<super::StatusResponse>, tonic::Status>;
        async fn create_forward(
            &self,
            request: tonic::Request<super::ForwardRequest>,
        ) -> std::result::Result<tonic::Response<super::ForwardResponse>, tonic::Status>;
        async fn delete_forward(
            &self,
            request: tonic::Request<super::ForwardDeleteRequest>,
        ) -> std::result::Result<tonic::Response<super::StatusResponse>, tonic::Status>;
        async fn get_public_setting(
            &self,
            request: tonic::Request<super::Empty>,
        ) -> std::result::Result<
            tonic::Response<super::PublicSettingResponse>,
            tonic::Status,
        >;
        async fn get_listen_ports(
            &self,
            request: tonic::Request<super::Empty>,
        ) -> std::result::Result<
            tonic::Response<super::ListenPortResponse>,
            tonic::Status,
        >;
        async fn get_port_info(
            &self,
            request: tonic::Request<super::PortInfoRequest>,
        ) -> std::result::Result<
            tonic::Response<super::PortInfoResponse>,
            tonic::Status,
        >;
        async fn handle_port_failure(
            &self,
            request: tonic::Request<super::PortFailureRequest>,
        ) -> std::result::Result<tonic::Response<super::StatusResponse>, tonic::Status>;
        async fn check_user_by_cookies(
            &self,
            request: tonic::Request<super::CookiesRequest>,
        ) -> std::result::Result<tonic::Response<super::UserResponse>, tonic::Status>;
        async fn record_session_lifecycle_log(
            &self,
            request: tonic::Request<super::SessionLifecycleLogRequest>,
        ) -> std::result::Result<tonic::Response<super::StatusResponse>, tonic::Status>;
        async fn face_recognition_callback(
            &self,
            request: tonic::Request<super::FaceRecognitionCallbackRequest>,
        ) -> std::result::Result<
            tonic::Response<super::FaceRecognitionCallbackResponse>,
            tonic::Status,
        >;
        async fn face_monitor_callback(
            &self,
            request: tonic::Request<super::FaceMonitorCallbackRequest>,
        ) -> std::result::Result<
            tonic::Response<super::FaceMonitorCallbackResponse>,
            tonic::Status,
        >;
        async fn join_face_monitor(
            &self,
            request: tonic::Request<super::JoinFaceMonitorRequest>,
        ) -> std::result::Result<
            tonic::Response<super::JoinFaceMonitorResponse>,
            tonic::Status,
        >;
    }
    #[derive(Debug)]
    pub struct ServiceServer<T> {
        inner: Arc<T>,
        accept_compression_encodings: EnabledCompressionEncodings,
        send_compression_encodings: EnabledCompressionEncodings,
        max_decoding_message_size: Option<usize>,
        max_encoding_message_size: Option<usize>,
    }
    impl<T> ServiceServer<T> {
        pub fn new(inner: T) -> Self {
            Self::from_arc(Arc::new(inner))
        }
        pub fn from_arc(inner: Arc<T>) -> Self {
            Self {
                inner,
                accept_compression_encodings: Default::default(),
                send_compression_encodings: Default::default(),
                max_decoding_message_size: None,
                max_encoding_message_size: None,
            }
        }
        pub fn with_interceptor<F>(
            inner: T,
            interceptor: F,
        ) -> InterceptedService<Self, F>
        where
            F: tonic::service::Interceptor,
        {
            InterceptedService::new(Self::new(inner), interceptor)
        }
        /// Enable decompressing requests with the given encoding.
        #[must_use]
        pub fn accept_compressed(mut self, encoding: CompressionEncoding) -> Self {
            self.accept_compression_encodings.enable(encoding);
            self
        }
        /// Compress responses with the given encoding, if the client supports it.
        #[must_use]
        pub fn send_compressed(mut self, encoding: CompressionEncoding) -> Self {
            self.send_compression_encodings.enable(encoding);
            self
        }
        /// Limits the maximum size of a decoded message.
        ///
        /// Default: `4MB`
        #[must_use]
        pub fn max_decoding_message_size(mut self, limit: usize) -> Self {
            self.max_decoding_message_size = Some(limit);
            self
        }
        /// Limits the maximum size of an encoded message.
        ///
        /// Default: `usize::MAX`
        #[must_use]
        pub fn max_encoding_message_size(mut self, limit: usize) -> Self {
            self.max_encoding_message_size = Some(limit);
            self
        }
    }
    impl<T, B> tonic::codegen::Service<http::Request<B>> for ServiceServer<T>
    where
        T: Service,
        B: Body + std::marker::Send + 'static,
        B::Error: Into<StdError> + std::marker::Send + 'static,
    {
        type Response = http::Response<tonic::body::BoxBody>;
        type Error = std::convert::Infallible;
        type Future = BoxFuture<Self::Response, Self::Error>;
        fn poll_ready(
            &mut self,
            _cx: &mut Context<'_>,
        ) -> Poll<std::result::Result<(), Self::Error>> {
            Poll::Ready(Ok(()))
        }
        fn call(&mut self, req: http::Request<B>) -> Self::Future {
            match req.uri().path() {
                "/message.Service/GetTokenAuthInfo" => {
                    #[allow(non_camel_case_types)]
                    struct GetTokenAuthInfoSvc<T: Service>(pub Arc<T>);
                    impl<T: Service> tonic::server::UnaryService<super::TokenRequest>
                    for GetTokenAuthInfoSvc<T> {
                        type Response = super::TokenResponse;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<super::TokenRequest>,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as Service>::get_token_auth_info(&inner, request).await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = GetTokenAuthInfoSvc(inner);
                        let codec = tonic::codec::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/message.Service/RenewToken" => {
                    #[allow(non_camel_case_types)]
                    struct RenewTokenSvc<T: Service>(pub Arc<T>);
                    impl<T: Service> tonic::server::UnaryService<super::TokenRequest>
                    for RenewTokenSvc<T> {
                        type Response = super::StatusResponse;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<super::TokenRequest>,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as Service>::renew_token(&inner, request).await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = RenewTokenSvc(inner);
                        let codec = tonic::codec::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/message.Service/CreateSession" => {
                    #[allow(non_camel_case_types)]
                    struct CreateSessionSvc<T: Service>(pub Arc<T>);
                    impl<
                        T: Service,
                    > tonic::server::UnaryService<super::SessionCreateRequest>
                    for CreateSessionSvc<T> {
                        type Response = super::SessionCreateResponse;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<super::SessionCreateRequest>,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as Service>::create_session(&inner, request).await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = CreateSessionSvc(inner);
                        let codec = tonic::codec::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/message.Service/FinishSession" => {
                    #[allow(non_camel_case_types)]
                    struct FinishSessionSvc<T: Service>(pub Arc<T>);
                    impl<
                        T: Service,
                    > tonic::server::UnaryService<super::SessionFinishRequest>
                    for FinishSessionSvc<T> {
                        type Response = super::SessionFinishResp;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<super::SessionFinishRequest>,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as Service>::finish_session(&inner, request).await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = FinishSessionSvc(inner);
                        let codec = tonic::codec::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/message.Service/UploadReplayFile" => {
                    #[allow(non_camel_case_types)]
                    struct UploadReplayFileSvc<T: Service>(pub Arc<T>);
                    impl<T: Service> tonic::server::UnaryService<super::ReplayRequest>
                    for UploadReplayFileSvc<T> {
                        type Response = super::ReplayResponse;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<super::ReplayRequest>,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as Service>::upload_replay_file(&inner, request).await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = UploadReplayFileSvc(inner);
                        let codec = tonic::codec::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/message.Service/UploadCommand" => {
                    #[allow(non_camel_case_types)]
                    struct UploadCommandSvc<T: Service>(pub Arc<T>);
                    impl<T: Service> tonic::server::UnaryService<super::CommandRequest>
                    for UploadCommandSvc<T> {
                        type Response = super::CommandResponse;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<super::CommandRequest>,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as Service>::upload_command(&inner, request).await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = UploadCommandSvc(inner);
                        let codec = tonic::codec::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/message.Service/DispatchTask" => {
                    #[allow(non_camel_case_types)]
                    struct DispatchTaskSvc<T: Service>(pub Arc<T>);
                    impl<
                        T: Service,
                    > tonic::server::StreamingService<super::FinishedTaskRequest>
                    for DispatchTaskSvc<T> {
                        type Response = super::TaskResponse;
                        type ResponseStream = T::DispatchTaskStream;
                        type Future = BoxFuture<
                            tonic::Response<Self::ResponseStream>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<
                                tonic::Streaming<super::FinishedTaskRequest>,
                            >,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as Service>::dispatch_task(&inner, request).await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = DispatchTaskSvc(inner);
                        let codec = tonic::codec::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.streaming(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/message.Service/ScanRemainReplays" => {
                    #[allow(non_camel_case_types)]
                    struct ScanRemainReplaysSvc<T: Service>(pub Arc<T>);
                    impl<
                        T: Service,
                    > tonic::server::UnaryService<super::RemainReplayRequest>
                    for ScanRemainReplaysSvc<T> {
                        type Response = super::RemainReplayResponse;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<super::RemainReplayRequest>,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as Service>::scan_remain_replays(&inner, request).await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = ScanRemainReplaysSvc(inner);
                        let codec = tonic::codec::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/message.Service/CreateCommandTicket" => {
                    #[allow(non_camel_case_types)]
                    struct CreateCommandTicketSvc<T: Service>(pub Arc<T>);
                    impl<
                        T: Service,
                    > tonic::server::UnaryService<super::CommandConfirmRequest>
                    for CreateCommandTicketSvc<T> {
                        type Response = super::CommandConfirmResponse;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<super::CommandConfirmRequest>,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as Service>::create_command_ticket(&inner, request).await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = CreateCommandTicketSvc(inner);
                        let codec = tonic::codec::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/message.Service/CheckOrCreateAssetLoginTicket" => {
                    #[allow(non_camel_case_types)]
                    struct CheckOrCreateAssetLoginTicketSvc<T: Service>(pub Arc<T>);
                    impl<
                        T: Service,
                    > tonic::server::UnaryService<super::AssetLoginTicketRequest>
                    for CheckOrCreateAssetLoginTicketSvc<T> {
                        type Response = super::AssetLoginTicketResponse;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<super::AssetLoginTicketRequest>,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as Service>::check_or_create_asset_login_ticket(
                                        &inner,
                                        request,
                                    )
                                    .await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = CheckOrCreateAssetLoginTicketSvc(inner);
                        let codec = tonic::codec::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/message.Service/CheckTicketState" => {
                    #[allow(non_camel_case_types)]
                    struct CheckTicketStateSvc<T: Service>(pub Arc<T>);
                    impl<T: Service> tonic::server::UnaryService<super::TicketRequest>
                    for CheckTicketStateSvc<T> {
                        type Response = super::TicketStateResponse;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<super::TicketRequest>,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as Service>::check_ticket_state(&inner, request).await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = CheckTicketStateSvc(inner);
                        let codec = tonic::codec::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/message.Service/CancelTicket" => {
                    #[allow(non_camel_case_types)]
                    struct CancelTicketSvc<T: Service>(pub Arc<T>);
                    impl<T: Service> tonic::server::UnaryService<super::TicketRequest>
                    for CancelTicketSvc<T> {
                        type Response = super::StatusResponse;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<super::TicketRequest>,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as Service>::cancel_ticket(&inner, request).await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = CancelTicketSvc(inner);
                        let codec = tonic::codec::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/message.Service/CreateForward" => {
                    #[allow(non_camel_case_types)]
                    struct CreateForwardSvc<T: Service>(pub Arc<T>);
                    impl<T: Service> tonic::server::UnaryService<super::ForwardRequest>
                    for CreateForwardSvc<T> {
                        type Response = super::ForwardResponse;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<super::ForwardRequest>,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as Service>::create_forward(&inner, request).await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = CreateForwardSvc(inner);
                        let codec = tonic::codec::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/message.Service/DeleteForward" => {
                    #[allow(non_camel_case_types)]
                    struct DeleteForwardSvc<T: Service>(pub Arc<T>);
                    impl<
                        T: Service,
                    > tonic::server::UnaryService<super::ForwardDeleteRequest>
                    for DeleteForwardSvc<T> {
                        type Response = super::StatusResponse;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<super::ForwardDeleteRequest>,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as Service>::delete_forward(&inner, request).await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = DeleteForwardSvc(inner);
                        let codec = tonic::codec::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/message.Service/GetPublicSetting" => {
                    #[allow(non_camel_case_types)]
                    struct GetPublicSettingSvc<T: Service>(pub Arc<T>);
                    impl<T: Service> tonic::server::UnaryService<super::Empty>
                    for GetPublicSettingSvc<T> {
                        type Response = super::PublicSettingResponse;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<super::Empty>,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as Service>::get_public_setting(&inner, request).await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = GetPublicSettingSvc(inner);
                        let codec = tonic::codec::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/message.Service/GetListenPorts" => {
                    #[allow(non_camel_case_types)]
                    struct GetListenPortsSvc<T: Service>(pub Arc<T>);
                    impl<T: Service> tonic::server::UnaryService<super::Empty>
                    for GetListenPortsSvc<T> {
                        type Response = super::ListenPortResponse;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<super::Empty>,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as Service>::get_listen_ports(&inner, request).await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = GetListenPortsSvc(inner);
                        let codec = tonic::codec::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/message.Service/GetPortInfo" => {
                    #[allow(non_camel_case_types)]
                    struct GetPortInfoSvc<T: Service>(pub Arc<T>);
                    impl<T: Service> tonic::server::UnaryService<super::PortInfoRequest>
                    for GetPortInfoSvc<T> {
                        type Response = super::PortInfoResponse;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<super::PortInfoRequest>,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as Service>::get_port_info(&inner, request).await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = GetPortInfoSvc(inner);
                        let codec = tonic::codec::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/message.Service/HandlePortFailure" => {
                    #[allow(non_camel_case_types)]
                    struct HandlePortFailureSvc<T: Service>(pub Arc<T>);
                    impl<
                        T: Service,
                    > tonic::server::UnaryService<super::PortFailureRequest>
                    for HandlePortFailureSvc<T> {
                        type Response = super::StatusResponse;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<super::PortFailureRequest>,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as Service>::handle_port_failure(&inner, request).await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = HandlePortFailureSvc(inner);
                        let codec = tonic::codec::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/message.Service/CheckUserByCookies" => {
                    #[allow(non_camel_case_types)]
                    struct CheckUserByCookiesSvc<T: Service>(pub Arc<T>);
                    impl<T: Service> tonic::server::UnaryService<super::CookiesRequest>
                    for CheckUserByCookiesSvc<T> {
                        type Response = super::UserResponse;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<super::CookiesRequest>,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as Service>::check_user_by_cookies(&inner, request).await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = CheckUserByCookiesSvc(inner);
                        let codec = tonic::codec::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/message.Service/RecordSessionLifecycleLog" => {
                    #[allow(non_camel_case_types)]
                    struct RecordSessionLifecycleLogSvc<T: Service>(pub Arc<T>);
                    impl<
                        T: Service,
                    > tonic::server::UnaryService<super::SessionLifecycleLogRequest>
                    for RecordSessionLifecycleLogSvc<T> {
                        type Response = super::StatusResponse;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<super::SessionLifecycleLogRequest>,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as Service>::record_session_lifecycle_log(
                                        &inner,
                                        request,
                                    )
                                    .await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = RecordSessionLifecycleLogSvc(inner);
                        let codec = tonic::codec::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/message.Service/FaceRecognitionCallback" => {
                    #[allow(non_camel_case_types)]
                    struct FaceRecognitionCallbackSvc<T: Service>(pub Arc<T>);
                    impl<
                        T: Service,
                    > tonic::server::UnaryService<super::FaceRecognitionCallbackRequest>
                    for FaceRecognitionCallbackSvc<T> {
                        type Response = super::FaceRecognitionCallbackResponse;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<
                                super::FaceRecognitionCallbackRequest,
                            >,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as Service>::face_recognition_callback(&inner, request)
                                    .await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = FaceRecognitionCallbackSvc(inner);
                        let codec = tonic::codec::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/message.Service/FaceMonitorCallback" => {
                    #[allow(non_camel_case_types)]
                    struct FaceMonitorCallbackSvc<T: Service>(pub Arc<T>);
                    impl<
                        T: Service,
                    > tonic::server::UnaryService<super::FaceMonitorCallbackRequest>
                    for FaceMonitorCallbackSvc<T> {
                        type Response = super::FaceMonitorCallbackResponse;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<super::FaceMonitorCallbackRequest>,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as Service>::face_monitor_callback(&inner, request).await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = FaceMonitorCallbackSvc(inner);
                        let codec = tonic::codec::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/message.Service/JoinFaceMonitor" => {
                    #[allow(non_camel_case_types)]
                    struct JoinFaceMonitorSvc<T: Service>(pub Arc<T>);
                    impl<
                        T: Service,
                    > tonic::server::UnaryService<super::JoinFaceMonitorRequest>
                    for JoinFaceMonitorSvc<T> {
                        type Response = super::JoinFaceMonitorResponse;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<super::JoinFaceMonitorRequest>,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as Service>::join_face_monitor(&inner, request).await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = JoinFaceMonitorSvc(inner);
                        let codec = tonic::codec::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                _ => {
                    Box::pin(async move {
                        let mut response = http::Response::new(empty_body());
                        let headers = response.headers_mut();
                        headers
                            .insert(
                                tonic::Status::GRPC_STATUS,
                                (tonic::Code::Unimplemented as i32).into(),
                            );
                        headers
                            .insert(
                                http::header::CONTENT_TYPE,
                                tonic::metadata::GRPC_CONTENT_TYPE,
                            );
                        Ok(response)
                    })
                }
            }
        }
    }
    impl<T> Clone for ServiceServer<T> {
        fn clone(&self) -> Self {
            let inner = self.inner.clone();
            Self {
                inner,
                accept_compression_encodings: self.accept_compression_encodings,
                send_compression_encodings: self.send_compression_encodings,
                max_decoding_message_size: self.max_decoding_message_size,
                max_encoding_message_size: self.max_encoding_message_size,
            }
        }
    }
    /// Generated gRPC service name
    pub const SERVICE_NAME: &str = "message.Service";
    impl<T> tonic::server::NamedService for ServiceServer<T> {
        const NAME: &'static str = SERVICE_NAME;
    }
}
