module Contacts
  class Create < BaseOperation
    around :transaction

    input_validation do
      required(:user).filled

      required(:attributes).schema do
        optional(:given_name).filled(:string)
        optional(:middle_name).maybe(:string)
        optional(:family_name).maybe(:string)
        optional(:suffix).maybe(:string)
        optional(:nickname).maybe(:string)
        optional(:gender).maybe(:string)
        optional(:pronouns).maybe(:string, included_in?: %w[FEMININE MASCULINE NEUTRAL])
      end
    end

    authorize Contact
    merge :build
    step :create

    def build(user:, attributes:)
      record = Contact.new({ user_id: user.id, **attributes })
      Success({ record: })
    end
  end
end
