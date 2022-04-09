module Contacts
  class AddPostalAddress < BaseOperation
    around :transaction

    input_validation do
      required(:user).filled
      required(:contact_id).filled(:string)

      required(:attributes).schema do
        required(:name).filled(:string)
        optional(:street).filled(:string)
        optional(:locality).filled(:string)
        optional(:region).filled(:string)
        optional(:postal_code).filled(:string)
        optional(:country).filled(:string)
      end
    end

    fetch Contact, :contact_id
    authorize operation: :update?
    merge :build
    tap :save
    step :return_record

    def build(record:, attributes:, **)
      postal_address = record.add_postal_address(attributes:)
      Success({ record:, postal_address: })
    end

    def return_record(postal_address:, **)
      Success(postal_address)
    end
  end
end
