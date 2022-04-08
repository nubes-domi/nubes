module Mutations
  class DestroyContactPostalAddress < Mutations::BaseMutation
    null true

    argument :id, ID

    field :contact_postal_address, Types::Models::ContactType
    field :errors, [String]

    def resolve(id:)
      contact_postal_address = context[:current_user].contact_postal_addresses.find(id)
      if contact_postal_address.destroy
        # Successful creation, return the created object with no errors
        {}
      else
        # Failed save, return the errors to the client
        {
          errors: contact_postal_address.errors.full_messages
        }
      end
    end
  end
end
