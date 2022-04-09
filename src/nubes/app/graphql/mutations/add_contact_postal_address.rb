module Mutations
  class AddContactPostalAddress < Mutations::BaseMutation
    null true

    ContactPostalAddress.define_graphql_mutation(self, Types::Models::ContactPostalAddressType)

    def resolve(contact_id:, **attributes)
      result = Contacts::AddPostalAddress.call(user: context[:current_user], contact_id:, attributes:)
      handle_failures(result) do |contact_postal_address|
        { contact_postal_address: }
      end
    end
  end
end
