module Mutations
  class DestroyContactPostalAddress < Mutations::BaseMutation
    null true

    argument :id, ID, description: "ID of the contact postal address to be destroyed"
    argument :contact_id, ID, description: "ID of the contact that owns the postal address"

    field :errors, [Types::Error]

    def resolve(id:, contact_id:)
      result = Contacts::DestroyPostalAddress.call(id:, contact_id:, user: context[:current_user])
      handle_failures(result) do
        {}
      end
    end
  end
end
