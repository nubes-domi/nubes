module Contacts
  class DestroyPostalAddress < BaseOperation
    around :transaction

    input_validation do
      required(:user).filled
      required(:id).filled(:string)
      required(:contact_id).filled(:string)
    end

    fetch Contact, :contact_id
    authorize operation: :update?
    step :destroy_postal_address

    def destroy_postal_address(record:, id:, **)
      record.destroy_postal_address(postal_address_id: id)
      record.save!
      Success({ record: })
    end
  end
end
