class ContactsController < ApplicationController
  def index
    @contacts = current_user.contacts.all
  end

  def new
    @contact = current_user.contacts.build
  end

  def create
    @contact = current_user.contacts.build(contact_params)
    if @contact.save
      redirect_to contact_path(@contact)
    else
      render :new, status: :unprocessable_entity
    end
  end

  def show
    @contact = current_user.contacts.find(params[:id])
  end

  def update
    @contact = current_user.contacts.find(params[:id])
    if @contact.update(contact_params)
      redirect_to contact_path(@contact)
    else
      # render :edit, status: :unprocessable_entity
      abort
    end
  end

  def destroy
    @contact = current_user.contacts.find(params[:id])
    @contact.destroy!

    redirect_to contacts_path
  end

  private

  def contact_params
    params.require(:contact).permit(:given_name, :middle_name, :family_name)
  end
end
