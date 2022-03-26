class SessionsController < ApplicationController
  before_action :ensure_identifier, only: [:show, :update]

  def new
  end

  def create
    result = Sessions.authentication_methods(params[:identifier])
    if result.authentication_methods.any?
      redirect_to authentication_path(result.authentication_methods.first, identifier: result.username)
    else
      flash[:errors] = { identifier: "Could not find a user with that name or email" }
      redirect_to signin_path
    end
  end

  def show
  end

  def update
    response = Sessions.create(request, params[:identifier], password: params[:password])
    if response
      start_session(response)
      redirect_to "/"
    else
      flash[:errors] = { password: "The password is incorrect" }
      redirect_to authentication_path(params[:method], identifier: params[:identifier])
    end
  end

  private

  def ensure_identifier
    redirect_to signin_path unless params[:identifier].present?
  end
end
