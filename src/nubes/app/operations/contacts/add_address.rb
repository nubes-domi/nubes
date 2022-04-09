module Contacts
  class AddAddress < BaseOperation
    around :transaction

    input_validation do
      required(:user).filled
      required(:contact_id).filled(:string)

      required(:attributes).schema do
        optional(:name).filled(:string)
        required(:kind).filled(:string)
        required(:data).filled(:string)
      end
    end

    fetch Contact, :contact_id
    authorize operation: :update?
    merge :build
    tap :save
    step :return_record

    def build(record:, attributes:, **)
      address = record.add_address(attributes:)
      Success({ record:, address: })
    end

    def return_record(address:, **)
      Success(address)
    end
  end
end
