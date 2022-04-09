module Contacts
  class Update < BaseOperation
    # around :transaction

    input_validation do
      required(:user).filled
      required(:id).filled(:string)

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

    fetch Contact
    authorize
    step :update
  end
end
