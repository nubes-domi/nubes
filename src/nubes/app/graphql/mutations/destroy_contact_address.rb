module Mutations
  class DestroyContactAddress < Mutations::BaseMutation
    null true

    field :contact_address, Types::Models::ContactType
    field :errors, [String]

    def resolve(id:)
      contact_address = context[:current_user].addresses.find(id)
      if contact_address.destroy
        # Successful creation, return the created object with no errors
        {}
      else
        # Failed save, return the errors to the client
        {
          errors: contact_address.errors.full_messages
        }
      end
    end
  end
end
