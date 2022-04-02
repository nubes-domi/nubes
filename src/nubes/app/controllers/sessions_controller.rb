class SessionsController < ApplicationController
  before_action :prevent_duplicate_sessions, only: %i[create show update]

  def new
  end

  def create
    run Sessions::Operation::Identify do |result|
      return redirect_to authentication_for(result["user"], continue: params[:continue])
    end

    flash[:errors] = { identifier: "Invalid username or email" }
    redirect_to signin_path(continue: params[:continue])
  end

  def show
    @user = User.find(params[:user_id])
  end

  def update
    ctx = run Sessions::Operation::Start, request: request do |result|
      start_session(result["session"])
      return signin_continue
    end

    flash[:errors] = { password: "Wrong password. Please try again." }
    redirect_to authentication_for(ctx["user"])
  end

  private

  def authentication_for(user, continue: nil)
    authentication_path(
      method: :password,
      user_id: user.id,
      continue: continue
    )
  end

  def prevent_duplicate_sessions
    return unless session_for?(params[:user_id])

    if current_user.id != params[:user_id]
      # TODO: reauthenticate if needed
      switch_session(session_for(params[:user_id]))
    end

    signin_continue
  end

  def signin_continue
    if params[:continue].starts_with?("/")
      redirect_to params[:continue]
    else
      redirect_to root_path
    end
  end
end
