class ContactPostalAddressesController < ApplicationController
  def new
    @contact_postal_address = current_user.contact_postal_addresses.build
  end

  def create
    @contact_postal_address = current_user.contact_postal_addresses.build(contact_postal_address_params)
    if @contact_postal_address.save
      redirect_to contact_path(@contact_postal_address.contact)
    else
      render :new, status: :unprocessable_entity
    end
  end

  def show
    @contact_postal_address = current_user.contact_postal_addresses.find(params[:id])
  end

  def update
    @contact_postal_address = current_user.contact_postal_addresses.find(params[:id])
    if @contact_postal_address.update(contact_postal_address_params)
      redirect_to contact_path(@contact_postal_address.contact)
    else
      # render :edit, status: :unprocessable_entity
      abort
    end
  end

  def destroy
    @contact_postal_address = current_user.contact_postal_addresses.find(params[:id])
    @contact_postal_address.destroy!

    redirect_to contact_path(@contact_postal_address.contact)
  end

  private

  def contact_postal_address_params
    params.require(:contact_postal_address).permit(
      :name, :street, :locality, :region, :postal_code, :country, :preferred
    )
  end
end
