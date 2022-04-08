class ContactsController < ApplicationController
  def index
    @contacts = current_user.contacts.all
  end

  def new
    @contact = current_user.contacts.build
  end

  def create
    op = Contacts::Operations::Create.call(current_user:, params: contact_params)
    if op.success?
      redirect_to contact_path(op["model"].id)
    else
      @contact = op["model"]
      render :new, status: :unprocessable_entity
    end
  end

  def show
    @contact = current_user.contacts.find(params[:id])
  end

  def update
    op = Contacts::Operations::Create.call(current_user:, contact_id: params[:id], params: contact_params)
    @contact = op["model"]
    if op.success?
      redirect_to contact_path(op["model"].id)
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
