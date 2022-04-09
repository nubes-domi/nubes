module Contacts
  class UpdatePostalAddress < BaseOperation
    around :transaction

    input_validation do
      required(:user).filled
      required(:id).filled(:string)
      required(:contact_id).filled(:string)

      required(:attributes).schema do
        optional(:name).filled(:string)
        optional(:street).filled(:string)
        optional(:locality).filled(:string)
        optional(:region).filled(:string)
        optional(:postal_code).filled(:string)
        optional(:country).filled(:string)
      end
    end

    fetch Contact, :contact_id
    authorize operation: :update?
    step :update_postal_address
    tap :save
    step :return_record

    def update_postal_address(record:, id:, attributes:, **)
      postal_address = record.update_postal_address(postal_address_id: id, attributes:)
      Success(record:, postal_address:)
    end

    def return_record(postal_address:, **)
      Success(postal_address)
    end
  end
end
