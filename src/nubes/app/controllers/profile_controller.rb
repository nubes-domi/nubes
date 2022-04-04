class ProfileController < ApplicationController
  def show
    @user = current_user
  end

  def name
    @user = current_user
  end

  def birthdate
    @user = current_user
  end

  def gender
    @user = current_user
  end

  def update
    @user = current_user

    if @user.update(user_params)
      redirect_to me_path
    else
      render :edit
    end
  end

  private

  def user_params
    params.require(:user).permit(
      :username, :given_name, :middle_name, :family_name, :birthdate, :email,
      :phone_number, :gender, :pronouns, :address_street, :address_locality,
      :address_region, :address_postal_code, :address_country
    )
  end
end
