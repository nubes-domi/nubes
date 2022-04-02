class SessionsController < ApplicationController
  def new
  end

  def create
    @user = User.identify(params[:identifier])
    if @user
      redirect_to authentication_for(@user)
    else
      flash[:errors] = { identifier: "Invalid username or email" }
      redirect_to signin_path
    end
  end

  def show
    @user = User.find(params[:user_id])
  end

  def update
    @user = User.find(params[:user_id])

    if @user.authenticate(params[:password])
      redirect_to root_path
    else
      flash[:errors] = { password: "Wrong password" }
      redirect_to authentication_for(@user)
    end
  end

  private

  def authentication_for(user)
    authentication_path(
      method: :password,
      user_id: user.id
    )
  end
end
