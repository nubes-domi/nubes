module Contacts
  class UpdateAddress < BaseOperation
    around :transaction

    input_validation do
      required(:user).filled
      required(:id).filled(:string)
      required(:contact_id).filled(:string)

      required(:attributes).schema do
        optional(:name).filled(:string)
        optional(:kind).filled(:string)
        optional(:data).filled(:string)
      end
    end

    fetch Contact, :contact_id
    authorize operation: :update?
    step :update_address
    tap :save
    step :return_record

    def update_address(record:, id:, attributes:, **)
      address = record.update_address(address_id: id, attributes:)
      Success(record:, address:)
    end

    def return_record(address:, **)
      Success(address)
    end
  end
end
