module Mutations
  class UpdateContactPostalAddress < Mutations::BaseMutation
    null true

    argument :contact_id, ID
    ContactPostalAddress.define_graphql_mutation(self, Types::Models::ContactPostalAddressType, type: :update)

    def resolve(id:, contact_id:, **attributes)
      result = Contacts::UpdatePostalAddress.call(user: context[:current_user], id:, contact_id:, attributes:)
      handle_failures(result) do |contact_postal_address|
        { contact_postal_address: }
      end
    end
  end
end
