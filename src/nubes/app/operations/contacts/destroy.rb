require "step_adapters/validate"

module Contacts
  class Destroy < BaseOperation
    # around :transaction

    input_validation do
      required(:user).filled
      required(:id).filled(:string)
    end

    fetch Contact
    authorize
    step :destroy
  end
end
