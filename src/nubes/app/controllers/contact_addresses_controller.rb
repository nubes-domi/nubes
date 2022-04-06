class ContactAddressesController < ApplicationController
  def new
    @contact = current_user.contacts.find(params[:contact_id])
    @contact_address = @contact.contact_addresses.build
  end

  def create
    @contact = current_user.contacts.find(params[:contact_id])
    @contact_address = @contact.contact_addresses.build(contact_address_params)
    if @contact_address.save
      redirect_to contact_path(@contact_address.contact)
    else
      render :new, status: :unprocessable_entity
    end
  end

  def show
    @contact_address = current_user.contact_addresses.find(params[:id])
  end

  def update
    @contact_address = current_user.contact_addresses.find(params[:id])
    if @contact_address.update(contact_address_params)
      redirect_to contact_path(@contact_address.contact)
    else
      # render :edit, status: :unprocessable_entity
      abort
    end
  end

  def destroy
    @contact_address = current_user.contact_addresses.find(params[:id])
    @contact_address.destroy!

    redirect_to contact_path(@contact_address.contact)
  end

  private

  def contact_address_params
    params.require(:contact_address).permit(
      :kind, :data, :name, :preferred
    )
  end
end
