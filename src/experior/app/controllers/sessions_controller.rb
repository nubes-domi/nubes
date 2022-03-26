class SessionsController < ApplicationController
  def index
  end

  def new
  end

  def create
    result = SUM_CLIENT.create(
      Sum::CreateSessionRequest.new(
        session: Sum::Session.new(
          user_agent: request.user_agent,
          ip_address: request.remote_addr
        ),
        username: params[:username],
        password: params[:password]
      )
    )

    cookies.permanent[:current_session] = result.access_token
    cookies.permanent[:sessions] = [(cookies[:sessions] || ""), result.access_token].join("|")
    redirect_to "/"
  rescue GRPC::Unauthenticated
    flash[:error] = "Invalid username or password"
    redirect_to signin_path
  end
end
