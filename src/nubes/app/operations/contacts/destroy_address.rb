module Contacts
  class DestroyAddress < BaseOperation
    around :transaction

    input_validation do
      required(:user).filled
      required(:id).filled(:string)
      required(:contact_id).filled(:string)
    end

    fetch Contact, :contact_id
    authorize operation: :update?
    step :destroy_address

    def destroy_address(record:, id:, **)
      record.destroy_address(address_id: id)
      record.save!
      Success({ record: })
    end
  end
end
