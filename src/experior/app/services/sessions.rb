module Sessions
  class << self
    def authentication_methods(identifier)
      SUM_SESSIONS.get_authentication_methods(
        Sum::GetAuthenticationMethodsRequest.new(identifier: identifier)
      )
    rescue GRPC::Unauthenticated
      []
    end

    def create(request, identifier, password:)
      SUM_SESSIONS.create(
        Sum::CreateSessionRequest.new(
          identifier: identifier, password: password,
          session: Sum::Session.new(
            user_agent: request.user_agent, ip_address: request.remote_addr
          )
        )
      )
    rescue GRPC::Unauthenticated
      false
    end
  end
end
