class SessionsController < ApplicationController
  skip_before_action :ensure_authenticated
  before_action :prevent_duplicate_sessions, only: %i[show update]

  include Dry::Monads[:result]

  def new
  end

  def create
    result = Sessions::Identify.call(identifier: params[:identifier])
    if result.success?
      redirect_to authentication_for(result.value!, continue: params[:continue])
    else
      flash[:errors] = { identifier: "Invalid username or email" }
      redirect_to signin_path(continue: params[:continue])
    end
  end

  def show
    @user = User.find(params[:user_id])
  end

  def update
    case Sessions::Create.call(user_id: params[:user_id], password: params[:password], request:)
    in Success(session)
      start_session(session)
      signin_continue
    in Failure(NotFound)
      flash[:errors] = { password: "Wrong password. Please try again." }
      redirect_to authentication_path(method: :password, user_id: params[:user_id], continue: params[:continue])
    end
  end

  private

  def authentication_for(user, continue: nil)
    authentication_path(
      method: :password,
      user_id: user.id,
      continue:
    )
  end

  def prevent_duplicate_sessions
    return unless session_for?(params[:user_id])

    if current_user.id != params[:user_id]
      # TODO: reauthenticate if needed
      switch_session(session_token_for(params[:user_id]))
    end

    signin_continue
  end

  def signin_continue
    if params[:continue].present? && params[:continue].starts_with?("/")
      redirect_to params[:continue]
    else
      redirect_to root_path
    end
  end
end
